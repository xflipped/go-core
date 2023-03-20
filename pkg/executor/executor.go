// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package executor

import (
	"bytes"
	"context"
	"os"
	"time"

	"git.fg-tech.ru/listware/proto/sdk/pbflink"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/Shopify/sarama"
	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	defaultBroker = "kafka:9092"

	defaultTopic = "router.system"

	defaultBrokerEnv = "KAFKA_ADDR"

	defaultTimeout = time.Second * 10
)

type Executor interface {
	Configure(...Opt) error

	ExecAsync(context.Context, ...*pbtypes.FunctionContext) error

	ExecSync(context.Context, *pbtypes.FunctionContext) error

	Close()
}

type executor struct {
	topic   string
	brokers []string
	timeout time.Duration

	ap sarama.AsyncProducer
	p  sarama.SyncProducer
}

func New(opts ...Opt) (Executor, error) {
	e := &executor{}
	return e, e.configure(opts...)
}

func (e *executor) Configure(opts ...Opt) (err error) {
	for _, opt := range opts {
		if err = opt(e); err != nil {
			return
		}
	}
	return
}

func asyncConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	return config
}

func syncConfig() *sarama.Config {
	config := asyncConfig()
	config.Producer.Return.Successes = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	return config
}

func (e *executor) configure(opts ...Opt) (err error) {
	if addr, ok := os.LookupEnv(defaultBrokerEnv); ok {
		e.brokers = append(e.brokers, addr)
	}

	if err = e.Configure(opts...); err != nil {
		return
	}

	if e.timeout == 0 {
		e.timeout = defaultTimeout
	}

	if e.topic == "" {
		e.topic = defaultTopic
	}

	if len(e.brokers) == 0 {
		e.brokers = append(e.brokers, defaultBroker)
	}

	if e.ap, err = sarama.NewAsyncProducer(e.brokers, asyncConfig()); err != nil {
		return
	}

	if e.p, err = sarama.NewSyncProducer(e.brokers, syncConfig()); err != nil {
		return
	}

	return
}

func (e *executor) Close() {
	e.ap.Close()
	e.p.Close()
}

func (e *executor) ExecAsync(ctx context.Context, msgs ...*pbtypes.FunctionContext) (err error) {
	for _, msg := range msgs {
		var buffer bytes.Buffer
		if err = statefun.MakeProtobufType(msg).Serialize(&buffer, msg); err != nil {
			return
		}

		message := &sarama.ProducerMessage{
			Topic: e.topic,
			Value: sarama.ByteEncoder(buffer.Bytes()),
		}

		e.ap.Input() <- message
	}
	return
}

func (e *executor) ExecSync(ctx context.Context, msg *pbtypes.FunctionContext) (err error) {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	msg.ReplyResult = &pbtypes.ReplyResult{
		Namespace: msg.FunctionType.Namespace,
		Topic:     msg.FunctionType.Type,
		Key:       uuid.New().String(),
	}

	var buffer bytes.Buffer
	if err = statefun.MakeProtobufType(msg).Serialize(&buffer, msg); err != nil {
		return
	}

	consumerGroup, err := sarama.NewConsumerGroup(e.brokers, msg.ReplyResult.Key, syncConfig())
	if err != nil {
		return
	}
	defer consumerGroup.Close()

	message := &sarama.ProducerMessage{
		Topic: e.topic,
		Value: sarama.ByteEncoder(buffer.Bytes()),
	}

	c := &Consumer{
		ready:  make(chan bool),
		key:    msg.ReplyResult.Key,
		output: make(chan error),
	}

	go func() {
		if err := consumerGroup.Consume(ctx, []string{msg.ReplyResult.Topic}, c); err != nil {
			cancel()
			return
		}
	}()

	<-c.ready

	if _, _, err = e.p.SendMessage(message); err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case err = <-c.output:
			return err
		}
	}

}

type Consumer struct {
	ready chan bool

	key string

	output chan error
}

func (consumer *Consumer) Setup(s sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
	for {
		select {
		case m := <-claim.Messages():
			session.MarkMessage(m, "")

			if string(m.Key) == consumer.key {
				consumer.output <- readResult(m.Value)
				return
			}

		case <-session.Context().Done():
			return session.Context().Err()
		}
	}
}

func readResult(value []byte) (err error) {
	var typedValue pbflink.TypedValue
	if err = statefun.MakeProtobufType(&typedValue).Deserialize(bytes.NewReader(value), &typedValue); err != nil {
		return
	}

	var kafkaProducerRecord pbflink.KafkaProducerRecord
	if err = statefun.MakeProtobufType(&kafkaProducerRecord).Deserialize(bytes.NewReader(typedValue.Value), &kafkaProducerRecord); err != nil {
		return
	}

	var functionResult pbtypes.FunctionResult
	if err = statefun.MakeProtobufType(&functionResult).Deserialize(bytes.NewReader(kafkaProducerRecord.GetValueBytes()), &functionResult); err != nil {
		return
	}

	if !functionResult.Complete {
		for _, errorMessage := range functionResult.Errors {
			if err == nil {
				err = errors.New(errorMessage)
			} else {
				err = errors.Wrap(err, errorMessage)
			}
		}
	}

	return
}
