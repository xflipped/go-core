// Copyright 2022 Listware

package main

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/vertex/types"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
)

type Disk struct {
	Name string `json:"name,omitempty"`
}

func createDiskType(ctx context.Context) (err error) {
	disk := Disk{
		Name: "nvme0n1",
	}

	pt := types.ReflectType(disk)

	message, err := system.CreateType(pt)
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}

func updateDiskType(ctx context.Context) (err error) {
	disk := Disk{
		Name: "nvme0n1",
	}

	pt := types.ReflectType(disk)

	//pt.Triggers["create"] = append(pt.Triggers["create"], logTriger)

	message, err := system.UpdateType(pt)
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}

func deleteDiskType(ctx context.Context) (err error) {
	message, err := system.DeleteType("disk")
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}

func createDiskObject(ctx context.Context) (err error) {
	disk := Disk{
		Name: "nvme0n1",
	}

	typeName := "disk"

	message, err := system.CreateObject(typeName, disk)
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}
