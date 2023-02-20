// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package module

import (
	"encoding/json"
	"path"

	"git.fg-tech.ru/listware/go-core/pkg/log"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
	"github.com/google/uuid"
)

type Context interface {
	statefun.Context
	Message() json.RawMessage
	States() []statefun.ValueSpec

	CmdbContext() json.RawMessage
	GetReplyResult() *pbtypes.ReplyResult
	AddError(error)
}

type contextAdapter struct {
	statefun.Context
	message     json.RawMessage
	cmdbContext json.RawMessage

	key       string
	syncTable *pbtypes.SyncTable

	states []statefun.ValueSpec

	OnReply  ReplyFunction
	OnResult ResultFunction
}

func (c *contextAdapter) Message() json.RawMessage {
	return c.message
}

func (c *contextAdapter) States() []statefun.ValueSpec {
	return c.states
}

func (c *contextAdapter) CmdbContext() json.RawMessage {
	return c.cmdbContext
}

func (c *contextAdapter) updateSyncTable() {
	c.Storage().Set(syncTableSpec, c.syncTable)
}

func (c *contextAdapter) GetReplyResult() *pbtypes.ReplyResult {
	key := uuid.New().String()
	c.syncTable.ResultTable[key] = c.key
	c.updateSyncTable()
	return &pbtypes.ReplyResult{
		Key:       key,
		Namespace: c.Self().FunctionType.GetNamespace(),
		Id:        c.Self().Id,
		Topic:     c.Self().FunctionType.GetType(),
	}
}

func (c *contextAdapter) checkReply() {
	errorContainer, ok := c.syncTable.ErrorsTable[c.key]
	if !ok {
		return
	}

	for _, value := range c.syncTable.ResultTable {
		if value == c.key {
			return
		}
	}

	delete(c.syncTable.ErrorsTable, c.key)
	c.updateSyncTable()

	replyResult, ok := c.syncTable.ReplyTable[c.key]
	if !ok {
		return
	}
	delete(c.syncTable.ReplyTable, c.key)
	c.updateSyncTable()

	if c.OnReply != nil {
		if !c.OnReply(c, replyResult, errorContainer) {
			return
		}
	}

	c.onReply(replyResult, errorContainer)
}

func (c *contextAdapter) onReply(replyResult *pbtypes.ReplyResult, errorContainer *pbtypes.ErrorContainer) {
	functionResult := &pbtypes.FunctionResult{
		ReplyEgress: replyResult,
		Complete:    errorContainer.Complete,
		Errors:      errorContainer.Errors,
	}

	if !replyResult.Egress {
		message := statefun.MessageBuilder{
			Target: statefun.Address{
				Id:           replyResult.GetId(),
				FunctionType: statefun.TypeNameFrom(path.Join(replyResult.GetNamespace(), replyResult.GetTopic())),
			},
			Value:     functionResult,
			ValueType: functionResultType,
		}
		c.Send(message)
	} else {
		kafkaEgressBuilder := statefun.KafkaEgressBuilder{
			Target:    kafkaEgressTypeName,
			Topic:     replyResult.Topic,
			Key:       replyResult.Key,
			Value:     functionResult,
			ValueType: functionResultType,
		}
		c.SendEgress(kafkaEgressBuilder)
	}
}

func (c *contextAdapter) AddError(err error) {
	if err == nil {
		return
	}

	log.Debug(err)

	if errorContainer, ok := c.syncTable.ErrorsTable[c.key]; ok {
		errorContainer.Complete = false
		errorContainer.Errors = append(errorContainer.Errors, err.Error())
	}
	c.updateSyncTable()
}
func (c *contextAdapter) onResult(functionResult *pbtypes.FunctionResult) {
	if c.OnResult != nil {
		c.OnResult(c, functionResult)
	}
}
