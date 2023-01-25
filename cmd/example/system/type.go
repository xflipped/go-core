// Copyright 2022 Listware

package main

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/vertex/types"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
)

type Str struct{}

func createTypes(ctx context.Context) (err error) {
	var names = []string{"group", "node", "cpu", "os", "baseboard", "bios", "mem", "netlink", "temp"}
	str := Str{}
	pt := types.ReflectType(str)

	for _, name := range names {
		pt.Schema.Title = name
		message, err := system.CreateType(pt)
		if err != nil {
			return err
		}

		if err = exec.ExecAsync(ctx, message); err != nil {
			return err
		}
	}

	return
}

func deleteType(ctx context.Context, ptName string) (err error) {
	message, err := system.DeleteType(ptName)
	if err != nil {
		return
	}

	return exec.ExecSync(ctx, message)
}

func createType(ctx context.Context) (err error) {
	str := Str{}
	pt := types.ReflectType(str)

	pt.Schema.Title = "temp"

	message, err := system.CreateType(pt)
	if err != nil {
		return err
	}

	if err = exec.ExecAsync(ctx, message); err != nil {
		return err
	}

	return
}

func updateType(ctx context.Context) (err error) {
	str := Str{}
	pt := types.ReflectType(str)

	pt.Schema.Title = "temp"

	message, err := system.UpdateType(pt)
	if err != nil {
		return
	}

	return exec.ExecSync(ctx, message)
}
