// Copyright 2022 Listware

package module

import (
	"errors"
	"strings"
	"time"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/vertex"
	"git.fg-tech.ru/listware/go-core/pkg/log"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
)

var (
	functionContextTypeName = statefun.MakeProtobufTypeWithTypeName(statefun.TypeNameFrom("type.googleapis.com/FunctionContext"))
)

// StatefulFunction function interface
type StatefulFunction func(Context) error

// functionMiddleware adapter implemets middleware
type functionMiddleware struct {
	parent *module
	f      StatefulFunction
}

// Invoke statefun.StatefulFunction implementation
func (fm *functionMiddleware) Invoke(ctx statefun.Context, message statefun.Message) (err error) {
	start := time.Now()
	defer func() {
		log.Debug(ctx.Self(), "took", time.Since(start))
	}()
	log.Debug(ctx.Self())

	// if ctx.Caller().FunctionType != nil {
	// 	log.Debug("caller", ctx.Caller())
	// }
	//

	var functionContext pbtypes.FunctionContext
	if err = message.As(functionContextTypeName, &functionContext); err != nil {
		log.Debug(ctx.Self(), "proto", err)
		// do not retry idempotent call
		return nil
	}

	context := &contextAdapter{
		Context: ctx,
		message: functionContext.Value,
		states:  fm.parent.states[ctx.Self().FunctionType.String()],
	}

	collectionKey := strings.Split(ctx.Self().Id, "/")
	if len(collectionKey) != 2 {
		log.Debug(ctx.Self(), "bad id")
		return nil
	}

	response, err := vertex.Read(ctx, collectionKey[1], collectionKey[0])
	if err != nil {
		log.Debug(err)
		return nil
	}
	context.cmdbContext = response.GetPayload()

	if context.cmdbContext == nil {
		log.Debug(ctx.Self(), "bad context")
		return nil
	}

	if err = fm.f(context); err != nil {
		log.Debug(err)
		// print all errors
	}

	if err = errors.Unwrap(err); err != nil {
		// retry only if wrapped error
		return
	}

	return nil
}
