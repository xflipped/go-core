// Copyright 2022 Listware

package system

import (
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"google.golang.org/protobuf/proto"
)

func prepareRegister(id string) (functionContext *pbtypes.FunctionContext) {
	functionType := &pbtypes.FunctionType{
		Namespace: namespace,
		Type:      registerType,
	}

	functionContext = &pbtypes.FunctionContext{
		Id:           id,
		FunctionType: functionType,
	}
	return
}

func Register(name string, types []*pbcmdb.RegisterTypeMessage, objects []*pbcmdb.RegisterObjectMessage, links []*pbcmdb.RegisterLinkMessage) (functionContext *pbtypes.FunctionContext, err error) {
	functionContext = prepareRegister(name)

	registerMessage := &pbcmdb.RegisterMessage{
		TypeMessages:   types,
		ObjectMessages: objects,
		LinkMessages:   links,
	}

	if functionContext.Value, err = proto.Marshal(registerMessage); err != nil {
		return
	}

	return
}
