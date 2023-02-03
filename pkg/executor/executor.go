// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"git.fg-tech.ru/listware/proto/sdk/pbflink"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/Shopify/sarama"
	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
	"github.com/google/uuid"
)

const (
	defaultBroker = "127.0.0.1:9092"

	inputTopic = "router.system"
)

type Executor interface {
	ExecAsync(context.Context, ...*pbtypes.FunctionContext) error

	ExecSync(context.Context, *pbtypes.FunctionContext) error

	Close()
}

type executor struct {
	p sarama.SyncProducer
	c sarama.Consumer
}

func New(brokers ...string) (Executor, error) {
	if addr := os.Getenv("KAFKA_ADDR"); addr != "" {
		brokers = append(brokers, addr)
	}

	if len(brokers) == 0 {
		brokers = append(brokers, defaultBroker)
	}

	e := &executor{}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	var err error
	if e.p, err = sarama.NewSyncProducer(brokers, config); err != nil {
		return nil, err
	}

	if e.c, err = sarama.NewConsumer(brokers, config); err != nil {
		return nil, err
	}

	return e, err
}

func (k *executor) Close() {
	k.p.Close()
	k.c.Close()
}

func (k *executor) ExecAsync(ctx context.Context, msgs ...*pbtypes.FunctionContext) (err error) {
	for _, msg := range msgs {
		var buffer bytes.Buffer
		if err = statefun.MakeProtobufType(msg).Serialize(&buffer, msg); err != nil {
			return
		}

		message := &sarama.ProducerMessage{
			Topic: inputTopic,
			// Key:   sarama.StringEncoder(msg.GetId()),
			Value: sarama.ByteEncoder(buffer.Bytes()),
		}
		if _, _, err = k.p.SendMessage(message); err != nil {
			return
		}
	}
	return
}

func (k *executor) ExecSync(ctx context.Context, msg *pbtypes.FunctionContext) (err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	replyResult := &pbtypes.ReplyResult{
		Namespace: msg.FunctionType.Namespace,
		Topic:     msg.FunctionType.Type,
		Key:       uuid.New().String(),
	}

	msg.ReplyResult = replyResult

	var buffer bytes.Buffer
	if err = statefun.MakeProtobufType(msg).Serialize(&buffer, msg); err != nil {
		return
	}

	message := &sarama.ProducerMessage{
		Topic: inputTopic,
		// Key:   sarama.StringEncoder(msg.GetId()),
		Value: sarama.ByteEncoder(buffer.Bytes()),
	}

	consumer, err := k.c.ConsumePartition(replyResult.Topic, 0, sarama.OffsetOldest)
	if err != nil {
		return
	}
	defer consumer.Close()
	inputChan := consumer.Messages()
	errorsInput := consumer.Errors()
	if _, _, err = k.p.SendMessage(message); err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case m := <-inputChan:
			if string(m.Key) == replyResult.Key {
				var typedValue pbflink.TypedValue
				if err = statefun.MakeProtobufType(&typedValue).Deserialize(bytes.NewReader(m.Value), &typedValue); err != nil {
					return err
				}

				var functionResult pbtypes.FunctionResult
				if err = statefun.MakeProtobufType(&functionResult).Deserialize(bytes.NewReader(typedValue.Value), &functionResult); err != nil {
					return err
				}
				if !functionResult.Complete {
					err = fmt.Errorf(strings.Join(functionResult.Errors, ";"))
				}
				return err
			}
		case err = <-errorsInput:
			return
		}

	}
}
