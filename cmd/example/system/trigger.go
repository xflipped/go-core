// Copyright 2022 Listware

package main

import (
	"context"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
)

var logTrigger = &pbcmdb.Trigger{
	Type: "create",
	FunctionType: &pbtypes.FunctionType{
		Namespace: "system",
		Type:      "log.system.functions.root",
	},
}
var initTrigger = &pbcmdb.Trigger{
	Type: "update",
	FunctionType: &pbtypes.FunctionType{
		Namespace: "proxy",
		Type:      "init.inventory.functions.root",
	},
}

func createInitTrigger(ctx context.Context) (err error) {
	message, err := system.AddLinkTrigger("types/node", "types/function", initTrigger)
	if err != nil {
		return
	}

	return exec.ExecSync(ctx, message)
}

func deleteInitTrigger(ctx context.Context) (err error) {
	message, err := system.DeleteLinkTrigger("types/node", "types/function", initTrigger)
	if err != nil {
		return
	}

	return exec.ExecSync(ctx, message)
}
