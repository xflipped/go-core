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

func prepareLink(id string) (fc *pbtypes.FunctionContext) {
	ft := &pbtypes.FunctionType{
		Namespace: namespace,
		Type:      linksType,
	}

	fc = &pbtypes.FunctionContext{
		Id:           id,
		FunctionType: ft,
	}
	return
}

func CreateLinkMessage(to, name, linktype string, payload any) (linkMessage *pbcmdb.LinkMessage, err error) {
	if payload == nil {
		return nil, errors.ErrPayloadNil
	}

	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return
	}

	linkMessage = &pbcmdb.LinkMessage{
		Method:  pbcmdb.Method_CREATE,
		Name:    name,
		To:      to,
		Type:    linktype,
		Payload: payloadRaw,
	}
	return
}

// CreateLink create link 'from' -> 'to' with 'name'
// from will be 'root', 'node', '17136214'
// to will be 'root', 'node', '17136214'
func CreateLink(from, to, name, linktype string, payload any) (functionContext *pbtypes.FunctionContext, err error) {
	linkMessage, err := CreateLinkMessage(to, name, linktype, payload)
	if err != nil {
		return
	}

	functionContext = prepareLink(from)

	if functionContext.Value, err = proto.Marshal(linkMessage); err != nil {
		return
	}
	return
}

func UpdateLink(from, name string, payload any) (fc *pbtypes.FunctionContext, err error) {
	if payload == nil {
		return nil, errors.ErrPayloadNil
	}

	payloadRaw, ok := payload.([]byte)
	if !ok {
		if payloadRaw, err = json.Marshal(payload); err != nil {
			return
		}
	}

	fc = prepareLink(from)

	om := &pbcmdb.LinkMessage{
		Method:  pbcmdb.Method_UPDATE,
		Name:    name,
		Payload: payloadRaw,
	}

	if fc.Value, err = proto.Marshal(om); err != nil {
		return
	}
	return
}

func ReplaceLink(from, name string, payload any) (fc *pbtypes.FunctionContext, err error) {
	if payload == nil {
		return nil, errors.ErrPayloadNil
	}

	payloadRaw, ok := payload.([]byte)
	if !ok {
		if payloadRaw, err = json.Marshal(payload); err != nil {
			return
		}
	}

	fc = prepareLink(from)

	om := &pbcmdb.LinkMessage{
		Method:  pbcmdb.Method_REPLACE,
		Name:    name,
		Payload: payloadRaw,
	}

	if fc.Value, err = proto.Marshal(om); err != nil {
		return
	}
	return
}

func DeleteLink(from, name string) (fc *pbtypes.FunctionContext, err error) {
	fc = prepareLink(from)

	om := &pbcmdb.LinkMessage{
		Method: pbcmdb.Method_DELETE,
		Name:   name,
	}

	if fc.Value, err = proto.Marshal(om); err != nil {
		return
	}
	return
}

func AddLinkTriggerMessage(to string, trigger *pbcmdb.Trigger) (linkMessage *pbcmdb.LinkMessage, err error) {
	if trigger == nil {
		return nil, errors.ErrPayloadNil
	}

	linkMessage = &pbcmdb.LinkMessage{
		Method: pbcmdb.Method_CREATE_TRIGGER,
		To:     to,
	}

	if linkMessage.Payload, err = proto.Marshal(trigger); err != nil {
		return
	}
	return
}

func AddLinkTrigger(from, to string, trigger *pbcmdb.Trigger) (functionContext *pbtypes.FunctionContext, err error) {
	linkMessage, err := AddLinkTriggerMessage(to, trigger)
	if err != nil {
		return
	}

	functionContext = prepareLink(from)

	if functionContext.Value, err = proto.Marshal(linkMessage); err != nil {
		return
	}
	return
}

func DeleteLinkTriggerMessage(to string, trigger *pbcmdb.Trigger) (linkMessage *pbcmdb.LinkMessage, err error) {
	if trigger == nil {
		return nil, errors.ErrPayloadNil
	}

	linkMessage = &pbcmdb.LinkMessage{
		Method: pbcmdb.Method_DELETE_TRIGGER,
		To:     to,
	}

	if linkMessage.Payload, err = proto.Marshal(trigger); err != nil {
		return
	}
	return
}

func DeleteLinkTrigger(from, to string, trigger *pbcmdb.Trigger) (functionContext *pbtypes.FunctionContext, err error) {
	linkMessage, err := DeleteLinkTriggerMessage(to, trigger)
	if err != nil {
		return
	}

	functionContext = prepareLink(from)

	if functionContext.Value, err = proto.Marshal(linkMessage); err != nil {
		return
	}

	return
}

func RegisterLink(from, to, name, linktype string, payload any, async bool) (registerLinkMessage *pbcmdb.RegisterLinkMessage, err error) {
	linkMessage, err := CreateLinkMessage(to, name, linktype, payload)
	if err != nil {
		return
	}

	registerLinkMessage = &pbcmdb.RegisterLinkMessage{
		Id:          from,
		LinkMessage: linkMessage,
		Async:       async,
	}
	return
}

func RegisterLinkTrigger(from, to string, trigger *pbcmdb.Trigger, async bool) (registerLinkMessage *pbcmdb.RegisterLinkMessage, err error) {
	linkMessage, err := AddLinkTriggerMessage(to, trigger)
	if err != nil {
		return
	}

	registerLinkMessage = &pbcmdb.RegisterLinkMessage{
		Id:          from,
		LinkMessage: linkMessage,
		Async:       async,
	}
	return
}
