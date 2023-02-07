// Copyright 2023 NJWS, Inc.
// Copyright 2022 Listware

package system

import (
	"encoding/json"

	"git.fg-tech.ru/listware/go-core/pkg/client/system/errors"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"google.golang.org/protobuf/proto"
)

func prepareObject(id string) (functionContext *pbtypes.FunctionContext) {
	functionType := &pbtypes.FunctionType{
		Namespace: namespace,
		Type:      objectsType,
	}

	functionContext = &pbtypes.FunctionContext{
		Id:           id,
		FunctionType: functionType,
	}
	return
}

func CreateChildMessage(moType, name string, payload any, functions ...*pbtypes.FunctionMessage) (objectMessage *pbcmdb.ObjectMessage, err error) {
	if payload == nil {
		return nil, errors.ErrPayloadNil
	}

	payloadRaw, ok := payload.([]byte)
	if !ok {
		if payloadRaw, err = json.Marshal(payload); err != nil {
			return
		}
	}

	objectMessage = &pbcmdb.ObjectMessage{
		Method:    pbcmdb.Method_CREATE_CHILD,
		Type:      moType,
		Name:      name,
		Payload:   payloadRaw,
		Functions: functions,
	}

	return
}

func CreateChild(from, moType, name string, payload any, functions ...*pbtypes.FunctionMessage) (functionContext *pbtypes.FunctionContext, err error) {
	objectMessage, err := CreateChildMessage(moType, name, payload, functions...)
	if err != nil {
		return
	}

	functionContext = prepareObject(from)

	if functionContext.Value, err = proto.Marshal(objectMessage); err != nil {
		return
	}
	return
}

func RegisterObject(from, moType, name string, payload any, async, router bool, functions ...*pbtypes.FunctionMessage) (registerObjectMessage *pbcmdb.RegisterObjectMessage, err error) {
	objectMessage, err := CreateChildMessage(moType, name, payload, functions...)
	if err != nil {
		return
	}

	registerObjectMessage = &pbcmdb.RegisterObjectMessage{
		Id:            from,
		ObjectMessage: objectMessage,
		Async:         async,
		Router:        router,
	}
	return
}

func UpdateObjectMessage(payload any) (objectMessage *pbcmdb.ObjectMessage, err error) {
	if payload == nil {
		return nil, errors.ErrPayloadNil
	}

	payloadRaw, ok := payload.([]byte)
	if !ok {
		if payloadRaw, err = json.Marshal(payload); err != nil {
			return
		}
	}

	objectMessage = &pbcmdb.ObjectMessage{
		Method:  pbcmdb.Method_UPDATE,
		Payload: payloadRaw,
	}
	return
}

func UpdateObject(id string, payload any) (functionContext *pbtypes.FunctionContext, err error) {
	functionContext = prepareObject(id)

	objectMessage, err := UpdateObjectMessage(payload)
	if err != nil {
		return
	}

	if functionContext.Value, err = proto.Marshal(objectMessage); err != nil {
		return
	}
	return
}

func ReplaceObjectMessage(payload any) (objectMessage *pbcmdb.ObjectMessage, err error) {
	if payload == nil {
		return nil, errors.ErrPayloadNil
	}

	payloadRaw, ok := payload.([]byte)
	if !ok {
		if payloadRaw, err = json.Marshal(payload); err != nil {
			return
		}
	}

	objectMessage = &pbcmdb.ObjectMessage{
		Method:  pbcmdb.Method_REPLACE,
		Payload: payloadRaw,
	}
	return
}

func ReplaceObject(id string, payload any) (functionContext *pbtypes.FunctionContext, err error) {
	functionContext = prepareObject(id)

	objectMessage, err := ReplaceObjectMessage(payload)
	if err != nil {
		return
	}

	if functionContext.Value, err = proto.Marshal(objectMessage); err != nil {
		return
	}
	return
}

func DeleteObjectMessage() (objectMessage *pbcmdb.ObjectMessage, err error) {
	objectMessage = &pbcmdb.ObjectMessage{
		Method: pbcmdb.Method_DELETE,
	}
	return
}

func DeleteObject(id string) (functionContext *pbtypes.FunctionContext, err error) {
	functionContext = prepareObject(id)

	objectMessage, err := DeleteObjectMessage()
	if err != nil {
		return
	}

	if functionContext.Value, err = proto.Marshal(objectMessage); err != nil {
		return
	}
	return
}
