// Copyright 2022 Listware

package module

import (
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
)

func ToMessage(functionContext *pbtypes.FunctionContext) (message statefun.MessageBuilder, err error) {
	message = statefun.MessageBuilder{
		Target: statefun.Address{
			Id: functionContext.GetId(),
		},
		Value: functionContext,
	}

	if message.Target.FunctionType, err = statefun.TypeNameFromParts(functionContext.GetFunctionType().GetNamespace(), functionContext.GetFunctionType().GetType()); err != nil {
		return
	}
	message.ValueType = statefun.MakeProtobufType(functionContext)
	return
}
