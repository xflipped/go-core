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

func prepareAdvancedLink(id string) (fc *pbtypes.FunctionContext) {
	ft := &pbtypes.FunctionType{
		Namespace: namespace,
		Type:      advancedLinksType,
	}

	fc = &pbtypes.FunctionContext{
		Id:           id,
		FunctionType: ft,
	}
	return
}

func UpdateAdvancedLink(id string, payload any) (fc *pbtypes.FunctionContext, err error) {
	if payload == nil {
		return nil, errors.ErrPayloadNil
	}

	payloadRaw, ok := payload.([]byte)
	if !ok {
		if payloadRaw, err = json.Marshal(payload); err != nil {
			return
		}
	}

	fc = prepareAdvancedLink(id)

	om := &pbcmdb.LinkMessage{
		Method:  pbcmdb.Method_UPDATE,
		Payload: payloadRaw,
	}

	if fc.Value, err = proto.Marshal(om); err != nil {
		return
	}
	return
}

func ReplaceAdvancedLink(id string, payload any) (fc *pbtypes.FunctionContext, err error) {
	if payload == nil {
		return nil, errors.ErrPayloadNil
	}

	payloadRaw, ok := payload.([]byte)
	if !ok {
		if payloadRaw, err = json.Marshal(payload); err != nil {
			return
		}
	}

	fc = prepareAdvancedLink(id)

	om := &pbcmdb.LinkMessage{
		Method:  pbcmdb.Method_REPLACE,
		Payload: payloadRaw,
	}

	if fc.Value, err = proto.Marshal(om); err != nil {
		return
	}
	return
}

func DeleteAdvancedLink(id string) (fc *pbtypes.FunctionContext, err error) {
	fc = prepareLink(id)

	om := &pbcmdb.LinkMessage{
		Method: pbcmdb.Method_DELETE,
	}

	if fc.Value, err = proto.Marshal(om); err != nil {
		return
	}
	return
}
