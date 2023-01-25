// Copyright 2022 Listware

package module

import (
	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
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
