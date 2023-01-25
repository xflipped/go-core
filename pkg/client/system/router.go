// Copyright 2022 Listware

package system

import (
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"google.golang.org/protobuf/proto"
)

func prepareRouter(qdsl string) (fc *pbtypes.FunctionContext) {
	ft := &pbtypes.FunctionType{
		Namespace: namespace,
		Type:      routerType,
	}

	fc = &pbtypes.FunctionContext{
		Id:           qdsl,
		FunctionType: ft,
	}
	return
}

func QdslRouter(qdsl string, ffc *pbtypes.FunctionContext) (fc *pbtypes.FunctionContext, err error) {
	fc = prepareRouter(qdsl)

	if fc.Value, err = proto.Marshal(ffc); err != nil {
		return
	}
	return
}
