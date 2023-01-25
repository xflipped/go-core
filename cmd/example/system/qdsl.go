// Copyright 2022 Listware

package main

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/vertex/types"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"google.golang.org/protobuf/proto"
)

func qdsl(ctx context.Context) (err error) {
	disk := Disk{
		Name: "nvme0n1",
	}

	ffc, err := system.CreateChild("", "disk", disk.Name, disk)
	if err != nil {
		return
	}

	message, err := system.QdslRouter("*.node.types", ffc)
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}

func qdsl2(ctx context.Context) (err error) {
	disk := Disk{
		Name: "nvme0n1",
	}

	pt := types.ReflectType(disk)
	ffc, err := system.CreateType(pt)
	if err != nil {
		return
	}

	/////

	ft := &pbtypes.FunctionType{
		Namespace: "dev0.office",
		Type:      "worker",
	}

	goFuncCtx := &pbtypes.FunctionContext{
		FunctionType: ft,
	}

	if goFuncCtx.Value, err = proto.Marshal(ffc); err != nil {
		return
	}

	_ = goFuncCtx

	_ = ffc
	message, err := system.QdslRouter("nodes.root", ffc)
	if err != nil {
		return
	}
	_ = message
	return exec.ExecAsync(ctx, message)
}
