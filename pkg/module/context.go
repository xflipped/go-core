// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package module

import (
	"encoding/json"

	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
	"github.com/google/uuid"
)

type Context interface {
	statefun.Context
	Message() json.RawMessage
	States() []statefun.ValueSpec

	CmdbContext() json.RawMessage
	GetReplyResult() *pbtypes.ReplyResult
}

type contextAdapter struct {
	statefun.Context
	message     json.RawMessage
	cmdbContext json.RawMessage
	fm          *functionMiddleware
	key         string
}

func (c *contextAdapter) Message() json.RawMessage {
	return c.message
}

func (c *contextAdapter) States() []statefun.ValueSpec {
	return c.fm.m.states[c.Self().FunctionType.String()]
}

func (c *contextAdapter) CmdbContext() json.RawMessage {
	return c.cmdbContext
}

func (c *contextAdapter) GetReplyResult() *pbtypes.ReplyResult {
	key := uuid.New().String()
	c.fm.syncTable.ResultTable[key] = c.key
	c.Storage().Set(c.fm.syncTableSpec, &c.fm.syncTable)
	return &pbtypes.ReplyResult{
		Key:       key,
		Namespace: c.Self().FunctionType.GetNamespace(),
		Id:        c.Self().Id,
		Topic:     c.Self().FunctionType.GetType(),
	}
}
