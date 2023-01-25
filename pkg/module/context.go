// Copyright 2022 Listware

package module

import (
	"encoding/json"

	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
)

type Context interface {
	statefun.Context
	Message() json.RawMessage
	States() []statefun.ValueSpec

	CmdbContext() json.RawMessage
}

type contextAdapter struct {
	statefun.Context
	message json.RawMessage
	states  []statefun.ValueSpec

	cmdbContext json.RawMessage
}

func (c *contextAdapter) Message() json.RawMessage {
	return c.message
}

func (c *contextAdapter) States() []statefun.ValueSpec {
	return c.states
}

func (c *contextAdapter) CmdbContext() json.RawMessage {
	return c.cmdbContext
}
