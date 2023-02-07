// Copyright 2023 NJWS, Inc.
// Copyright 2022 Listware

package system

import (
	"encoding/json"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/vertex/types"
	"git.fg-tech.ru/listware/go-core/pkg/client/system/errors"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"google.golang.org/protobuf/proto"
)

func prepareType(id string) (functionContext *pbtypes.FunctionContext) {
	functionType := &pbtypes.FunctionType{
		Namespace: namespace,
		Type:      typesType,
	}

	functionContext = &pbtypes.FunctionContext{
		Id:           id,
		FunctionType: functionType,
	}
	return
}

func CreateTypeMessage(pt *types.Type) (typeMessage *pbcmdb.TypeMessage, err error) {
	if pt == nil {
		return nil, errors.ErrPayloadNil
	}

	if pt.Schema == nil {
		return nil, errors.ErrPayloadNil
	}

	typeMessage = &pbcmdb.TypeMessage{
		Method: pbcmdb.Method_CREATE,
		Name:   pt.Schema.Title,
	}

	if typeMessage.Payload, err = json.Marshal(pt); err != nil {
		return
	}

	return
}

func CreateType(pt *types.Type) (functionContext *pbtypes.FunctionContext, err error) {
	typeMessage, err := CreateTypeMessage(pt)
	if err != nil {
		return
	}

	functionContext = prepareType("system/types")

	if functionContext.Value, err = proto.Marshal(typeMessage); err != nil {
		return
	}

	return
}

func RegisterType(pt *types.Type, async bool) (registerTypeMessage *pbcmdb.RegisterTypeMessage, err error) {
	typeMessage, err := CreateTypeMessage(pt)
	if err != nil {
		return
	}

	registerTypeMessage = &pbcmdb.RegisterTypeMessage{
		Id:          "system/types",
		TypeMessage: typeMessage,
		Async:       async,
	}
	return
}

func UpdateTypeMessage(pt *types.Type) (typeMessage *pbcmdb.TypeMessage, err error) {
	if pt == nil {
		return nil, errors.ErrPayloadNil
	}

	if pt.Schema == nil {
		return nil, errors.ErrPayloadNil
	}

	typeMessage = &pbcmdb.TypeMessage{
		Method: pbcmdb.Method_UPDATE,
		Name:   pt.Schema.Title,
	}

	if typeMessage.Payload, err = json.Marshal(pt); err != nil {
		return
	}

	return
}

func UpdateType(pt *types.Type) (functionContext *pbtypes.FunctionContext, err error) {
	typeMessage, err := UpdateTypeMessage(pt)
	if err != nil {
		return
	}

	functionContext = prepareType("types/" + pt.Schema.Title)

	if functionContext.Value, err = proto.Marshal(typeMessage); err != nil {
		return
	}

	return
}

func ReplaceTypeMessage(pt *types.Type) (typeMessage *pbcmdb.TypeMessage, err error) {
	if pt == nil {
		return nil, errors.ErrPayloadNil
	}

	if pt.Schema == nil {
		return nil, errors.ErrPayloadNil
	}

	typeMessage = &pbcmdb.TypeMessage{
		Method: pbcmdb.Method_REPLACE,
		Name:   pt.Schema.Title,
	}

	if typeMessage.Payload, err = json.Marshal(pt); err != nil {
		return
	}

	return
}

func ReplaceType(pt *types.Type) (functionContext *pbtypes.FunctionContext, err error) {
	typeMessage, err := ReplaceTypeMessage(pt)
	if err != nil {
		return
	}

	functionContext = prepareType("types/" + pt.Schema.Title)

	if functionContext.Value, err = proto.Marshal(typeMessage); err != nil {
		return
	}

	return
}

func DeleteTypeMessage(moType string) (typeMessage *pbcmdb.TypeMessage, err error) {
	typeMessage = &pbcmdb.TypeMessage{
		Method: pbcmdb.Method_DELETE,
		Name:   moType,
	}

	return
}

func DeleteType(moType string) (functionContext *pbtypes.FunctionContext, err error) {
	typeMessage, err := DeleteTypeMessage(moType)
	if err != nil {
		return
	}

	functionContext = prepareType(moType)

	if functionContext.Value, err = proto.Marshal(typeMessage); err != nil {
		return
	}
	return
}

func CreateObjectMessage(moType string, payload any, functions ...*pbtypes.FunctionMessage) (typeMessage *pbcmdb.TypeMessage, err error) {
	if payload == nil {
		return nil, errors.ErrPayloadNil
	}

	payloadRaw, ok := payload.([]byte)
	if !ok {
		if payloadRaw, err = json.Marshal(payload); err != nil {
			return
		}
	}

	typeMessage = &pbcmdb.TypeMessage{
		Method:    pbcmdb.Method_CREATE_CHILD,
		Name:      moType,
		Functions: functions,
		Payload:   payloadRaw,
	}
	return
}

func CreateObject(moType string, payload any, functions ...*pbtypes.FunctionMessage) (functionContext *pbtypes.FunctionContext, err error) {
	typeMessage, err := CreateObjectMessage(moType, payload, functions...)
	if err != nil {
		return
	}

	functionContext = prepareType(moType)

	if functionContext.Value, err = proto.Marshal(typeMessage); err != nil {
		return
	}
	return
}

func AddTriggerMessage(moType string, trigger *pbcmdb.Trigger) (typeMessage *pbcmdb.TypeMessage, err error) {
	if trigger == nil {
		return nil, errors.ErrPayloadNil
	}

	typeMessage = &pbcmdb.TypeMessage{
		Method: pbcmdb.Method_CREATE_TRIGGER,
		Name:   moType,
	}

	if typeMessage.Payload, err = proto.Marshal(trigger); err != nil {
		return
	}
	return
}

func AddTrigger(moType string, trigger *pbcmdb.Trigger) (functionContext *pbtypes.FunctionContext, err error) {
	typeMessage, err := AddTriggerMessage(moType, trigger)
	if err != nil {
		return
	}

	functionContext = prepareType(moType)

	if functionContext.Value, err = proto.Marshal(typeMessage); err != nil {
		return
	}
	return
}

func DeleteTriggerMessage(moType string, trigger *pbcmdb.Trigger) (typeMessage *pbcmdb.TypeMessage, err error) {
	if trigger == nil {
		return nil, errors.ErrPayloadNil
	}

	typeMessage = &pbcmdb.TypeMessage{
		Method: pbcmdb.Method_DELETE_TRIGGER,
		Name:   moType,
	}

	if typeMessage.Payload, err = proto.Marshal(trigger); err != nil {
		return
	}
	return
}
func DeleteTrigger(moType string, trigger *pbcmdb.Trigger) (functionContext *pbtypes.FunctionContext, err error) {
	typeMessage, err := DeleteTriggerMessage(moType, trigger)
	if err != nil {
		return
	}

	functionContext = prepareType(moType)

	if functionContext.Value, err = proto.Marshal(typeMessage); err != nil {
		return
	}

	return
}
