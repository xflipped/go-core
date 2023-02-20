// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package module

import (
	"fmt"

	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
)

var (
	ErrReservedName = fmt.Errorf("reserved name")
	ErrEmptyName    = fmt.Errorf("empty name")
)

type Opt func(*module)

func WithStatefulFunctions(builder statefun.StatefulFunctions) Opt {
	return func(m *module) {
		m.builder = builder
	}
}

func WithPort(port int) Opt {
	return func(m *module) {
		m.port = port
	}
}

type BindOpt func(*functionMiddleware) error

func WithValueSpec(state statefun.ValueSpec) BindOpt {
	return func(f *functionMiddleware) (err error) {
		if state.Name == syncTableSpec.Name {
			return ErrReservedName
		}

		if state.Name == "" {
			return ErrEmptyName
		}

		f.states = append(f.states, state)
		return
	}
}

func WithOnReply(onReply ReplyFunction) BindOpt {
	return func(f *functionMiddleware) (err error) {
		f.OnReply = onReply
		return
	}
}

func WithOnResult(onResult ResultFunction) BindOpt {
	return func(f *functionMiddleware) (err error) {
		f.OnResult = onResult
		return
	}
}
