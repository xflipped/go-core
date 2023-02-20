// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package module

import (
	"encoding/json"
	"errors"
	"fmt"
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

type ResultFunction func(Context, *pbtypes.FunctionResult)

type ReplyFunction func(Context, *pbtypes.ReplyResult, *pbtypes.ErrorContainer) bool

// functionMiddleware adapter implemets middleware
type functionMiddleware struct {
	f StatefulFunction

	states []statefun.ValueSpec

	OnReply  ReplyFunction
	OnResult ResultFunction
}

// Invoke statefun.StatefulFunction implementation
func (f *functionMiddleware) Invoke(ctx statefun.Context, message statefun.Message) (err error) {

	// get state of sync table
	var syncTable pbtypes.SyncTable
	ctx.Storage().Get(syncTableSpec, &syncTable)

	// if caller
	if ctx.Caller().FunctionType != nil {
		// if functionResult message type
		var functionResult pbtypes.FunctionResult
		if message.Is(functionResultType) {
			if err = message.As(functionResultType, &functionResult); err != nil {
				return
			}
			// on result
			return f.onResult(ctx, &functionResult, &syncTable)
		}

		log.Debug("caller", ctx.Caller())
	}

	// log time of function only
	start := time.Now()
	defer func() {
		log.Debug(ctx.Self(), "took", time.Since(start))
	}()
	log.Debug(ctx.Self())

	var functionContext pbtypes.FunctionContext
	if err = message.As(functionContextType, &functionContext); err != nil {
		return
	}

	// prepare contextAdapter
	context, err := f.onInit(ctx, &functionContext, &syncTable)
	if err != nil {
		return
	}

	// exec user function code
	if err = f.f(context); err != nil {
		context.AddError(err)
	}

	// if need reply to caller
	context.checkReply()

	return errors.Unwrap(err)
}

func (f *functionMiddleware) onInit(ctx statefun.Context, functionContext *pbtypes.FunctionContext, syncTable *pbtypes.SyncTable) (context *contextAdapter, err error) {
	var key string
	syncTable.ErrorsTable = make(map[string]*pbtypes.ErrorContainer)
	syncTable.ReplyTable = make(map[string]*pbtypes.ReplyResult)
	syncTable.ResultTable = make(map[string]string)

	if replyResult := functionContext.GetReplyResult(); replyResult != nil {
		key = replyResult.GetKey()

		if ctx.Caller().FunctionType != nil {
			replyResult.Egress = false
		} else {
			replyResult.Egress = true
		}

		syncTable.ReplyTable[key] = replyResult

	} else {
		key = uuid.New().String()
	}

	if context, err = f.getContext(ctx, functionContext.Value, key, syncTable); err != nil {
		return
	}

	// update only if no error
	syncTable.ErrorsTable[key] = &pbtypes.ErrorContainer{Complete: true}
	context.updateSyncTable()

	return
}

func (f *functionMiddleware) getContext(ctx statefun.Context, message json.RawMessage, key string, syncTable *pbtypes.SyncTable) (context *contextAdapter, err error) {
	collectionKey := strings.Split(ctx.Self().Id, "/")
	if len(collectionKey) != 2 {
		err = fmt.Errorf("unknown id: %+v", ctx.Self().Id)
		return
	}

	response, err := vertex.Read(ctx, collectionKey[1], collectionKey[0])
	if err != nil {
		return
	}

	context = &contextAdapter{
		Context:     ctx,
		states:      f.states,
		OnReply:     f.OnReply,
		OnResult:    f.OnResult,
		key:         key,
		message:     message,
		cmdbContext: response.GetPayload(),
		syncTable:   syncTable,
	}
	return
}

func (f *functionMiddleware) onResult(ctx statefun.Context, functionResult *pbtypes.FunctionResult, syncTable *pbtypes.SyncTable) (err error) {
	replyKey := functionResult.GetReplyEgress().GetKey()

	key, ok := syncTable.ResultTable[replyKey]
	if !ok {
		return
	}

	delete(syncTable.ResultTable, replyKey)

	if !functionResult.Complete {
		if errorContainer, ok := syncTable.ErrorsTable[key]; ok {
			errorContainer.Complete = false
			errorContainer.Errors = append(errorContainer.Errors, functionResult.Errors...)
		}
	}

	ctx.Storage().Set(syncTableSpec, syncTable)

	context, err := f.getContext(ctx, nil, key, syncTable)
	if err != nil {
		return
	}
	context.onResult(functionResult)
	context.checkReply()
	return
}
