// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package module

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/vertex"
	"git.fg-tech.ru/listware/go-core/pkg/log"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
	"github.com/google/uuid"
)

// StatefulFunction function interface
type StatefulFunction func(Context) error

// functionMiddleware adapter implemets middleware
type functionMiddleware struct {
	syncTableSpec statefun.ValueSpec
	syncTable     pbtypes.SyncTable
	// key           string

	m *module
	f StatefulFunction
}

// Invoke statefun.StatefulFunction implementation
func (fm *functionMiddleware) Invoke(ctx statefun.Context, message statefun.Message) (err error) {
	ctx.Storage().Get(fm.syncTableSpec, &fm.syncTable)

	if ctx.Caller().FunctionType != nil {
		// log.Debug("caller", ctx.Caller())

		var functionResult pbtypes.FunctionResult
		if message.Is(functionResultType) {
			if err = message.As(functionResultType, &functionResult); err != nil {
				return
			}
			fm.onResult(ctx, &functionResult)
			return
		}
	}

	start := time.Now()
	defer func() {
		log.Debug(ctx.Self(), "took", time.Since(start))
	}()
	log.Debug(ctx.Self())

	var functionContext pbtypes.FunctionContext
	if err = message.As(functionContextType, &functionContext); err != nil {
		return
	}

	key := fm.onInit(ctx, &functionContext)

	if err = fm.invoke(ctx, &functionContext, key); err != nil {
		log.Debug(err)
		fm.onError(ctx, key, err)
	}

	fm.onReply(ctx, key)

	return errors.Unwrap(err)
}

func (fm *functionMiddleware) onInit(ctx statefun.Context, functionContext *pbtypes.FunctionContext) (key string) {
	fm.syncTable.ErrorsTable = make(map[string]*pbtypes.ErrorContainer)
	fm.syncTable.ReplyTable = make(map[string]*pbtypes.ReplyResult)
	fm.syncTable.ResultTable = make(map[string]string)

	if replyResult := functionContext.GetReplyResult(); replyResult != nil {
		key = replyResult.GetKey()

		if ctx.Caller().FunctionType != nil {
			replyResult.Egress = false
		} else {
			replyResult.Egress = true
		}

		fm.syncTable.ReplyTable[key] = replyResult

	} else {
		key = uuid.New().String()
	}
	fm.syncTable.ErrorsTable[key] = &pbtypes.ErrorContainer{Complete: true}
	ctx.Storage().Set(fm.syncTableSpec, &fm.syncTable)
	return
}

func (fm *functionMiddleware) invoke(ctx statefun.Context, functionContext *pbtypes.FunctionContext, key string) (err error) {
	collectionKey := strings.Split(ctx.Self().Id, "/")
	if len(collectionKey) != 2 {
		return fmt.Errorf("unknown id: %+v", ctx.Self().Id)
	}

	response, err := vertex.Read(ctx, collectionKey[1], collectionKey[0])
	if err != nil {
		return
	}

	context := &contextAdapter{
		fm:          fm,
		key:         key,
		Context:     ctx,
		message:     functionContext.Value,
		cmdbContext: response.GetPayload(),
	}

	return fm.f(context)
}

func (fm *functionMiddleware) onError(ctx statefun.Context, key string, err error) {
	if errorContainer, ok := fm.syncTable.ErrorsTable[key]; ok {
		errorContainer.Complete = false
		errorContainer.Errors = append(errorContainer.Errors, err.Error())
	}
}

func (fm *functionMiddleware) onResult(ctx statefun.Context, functionResult *pbtypes.FunctionResult) {
	replyKey := functionResult.GetReplyEgress().GetKey()
	key := fm.syncTable.ResultTable[replyKey]
	delete(fm.syncTable.ResultTable, replyKey)

	if !functionResult.Complete {
		if errorContainer, ok := fm.syncTable.ErrorsTable[key]; ok {
			errorContainer.Complete = false
			errorContainer.Errors = append(errorContainer.Errors, functionResult.Errors...)
		}
	}

	ctx.Storage().Set(fm.syncTableSpec, &fm.syncTable)
	fm.onReply(ctx, key)
}

func (fm *functionMiddleware) onReply(ctx statefun.Context, key string) {
	errorContainer, ok := fm.syncTable.ErrorsTable[key]
	if !ok {
		return
	}

	for _, value := range fm.syncTable.ResultTable {
		if value == key {
			return
		}
	}

	delete(fm.syncTable.ErrorsTable, key)
	ctx.Storage().Set(fm.syncTableSpec, &fm.syncTable)

	replyResult, ok := fm.syncTable.ReplyTable[key]
	if !ok {
		return
	}
	delete(fm.syncTable.ReplyTable, key)
	ctx.Storage().Set(fm.syncTableSpec, &fm.syncTable)

	fm.reply(ctx, replyResult, errorContainer)
}

func (fm *functionMiddleware) reply(ctx statefun.Context, replyResult *pbtypes.ReplyResult, errorContainer *pbtypes.ErrorContainer) {
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
		ctx.Send(message)
	} else {
		kafkaEgressBuilder := statefun.KafkaEgressBuilder{
			Target:    kafkaEgressTypeName,
			Topic:     replyResult.Topic,
			Key:       replyResult.Key,
			Value:     functionResult,
			ValueType: functionResultType,
		}
		ctx.SendEgress(kafkaEgressBuilder)
	}

}
