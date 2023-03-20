// Copyright 2022 Listware

package main

import (
	"context"
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/vertex/types"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
)

var (
	exec executor.Executor
)

func main() {
	var err error
	exec, err = executor.New(executor.WithBroker("127.0.0.1:9092"))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer exec.Close()

	ctx := context.Background()
	if err = register(ctx); err != nil {
		fmt.Println(err)
		return
	}
}

type NodeContainer struct{}

func register(ctx context.Context) (err error) {
	var registerTypes = []*pbcmdb.RegisterTypeMessage{}

	nodeContainerType, err := system.RegisterType(types.ReflectType(NodeContainer{}), true)
	if err != nil {
		return err
	}

	nodeType, err := system.RegisterType(types.ReflectType(Node{}), true)
	if err != nil {
		return err
	}
	_ = nodeType

	// nolint
	registerTypes = append(registerTypes, nodeContainerType, nodeType)

	var registerObjects = []*pbcmdb.RegisterObjectMessage{}

	nodeContainer, err := system.RegisterObject("system/root", "types/node-container", "nodes", NodeContainer{}, false, false)
	if err != nil {
		return
	}

	nodeContainer1, err := system.RegisterObject("nodes.root", "types/node", "node1", NodeContainer{}, false, true)
	if err != nil {
		return
	}
	registerObjects = append(registerObjects, nodeContainer, nodeContainer1)

	message, err := system.Register("appname", nil, registerObjects, nil)
	if err != nil {
		return
	}

	if err = exec.ExecSync(ctx, message); err != nil {
		return err
	}
	return
}
