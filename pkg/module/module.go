// Copyright 2022 Listware

package module

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"

	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
)

type Module interface {
	// Bind new function to module
	Bind(string, StatefulFunction, ...statefun.ValueSpec) error

	// RegisterAndListen register module and listen port
	RegisterAndListen(context.Context) error

	// listener port
	Port() int
	// addr for register
	Addr() string
}

type module struct {
	builder statefun.StatefulFunctions

	namespace string
	port      int

	states map[string][]statefun.ValueSpec

	listener net.Listener
}

// New module
func newModule(namespace string) *module {
	return &module{
		builder:   statefun.StatefulFunctionsBuilder(),
		states:    make(map[string][]statefun.ValueSpec),
		namespace: namespace,
	}
}

func New(namespace string, opts ...Opt) Module {
	m := newModule(namespace)

	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Bind new function to module
func (m *module) Bind(methodType string, f StatefulFunction, states ...statefun.ValueSpec) error {
	typeName := statefun.TypeNameFrom(path.Join(m.namespace, methodType))
	m.states[typeName.String()] = states
	return m.builder.WithSpec(statefun.StatefulFunctionSpec{
		FunctionType: typeName,
		States:       states,
		Function:     &functionMiddleware{parent: m, f: f},
	})
}

// RegisterAndListen register module and listen port
func (m *module) RegisterAndListen(ctx context.Context) (err error) {
	mux := http.NewServeMux()
	mux.Handle("/statefun", m.builder.AsHandler())
	srv := &http.Server{Handler: mux}

	if m.listener, err = net.Listen("tcp4", fmt.Sprintf(":%d", m.port)); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if err = srv.Serve(m.listener); err != nil {
			cancel()
		}
	}()

	<-ctx.Done()

	if err != nil {
		return
	}

	return srv.Shutdown(ctx)
}

func (m *module) Port() int {
	if m.listener != nil {
		if addr, ok := m.listener.Addr().(*net.TCPAddr); ok {
			return addr.Port
		}
	}
	return m.port
}

func (m *module) Addr() string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("http://%s:%d/statefun", hostname, m.Port())
}
