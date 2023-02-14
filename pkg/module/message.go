// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package module

import (
	"path"

	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
)

func ToMessage(functionContext *pbtypes.FunctionContext) (message statefun.MessageBuilder, err error) {
	message = statefun.MessageBuilder{
		Target: statefun.Address{
			Id:           functionContext.GetId(),
			FunctionType: statefun.TypeNameFrom(path.Join(functionContext.GetFunctionType().GetNamespace(), functionContext.GetFunctionType().GetType())),
		},
		Value:     functionContext,
		ValueType: statefun.MakeProtobufType(functionContext),
	}

	return
}
