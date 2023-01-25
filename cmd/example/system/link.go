// Copyright 2022 Listware

package main

import (
	"context"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
)

func createLink(ctx context.Context) (err error) {
	node := Node{
		Hostname: "sky",
	}

	message, err := system.CreateLink("objects/a0c17af0-45aa-4cdb-ace4-f60440583c8a", "objects/64f21ca8-c2c7-4fc1-bccc-de04842dbd49", "parent", "node", node)
	if err != nil {
		return
	}

	return exec.ExecSync(ctx, message)
}

func updateLink(ctx context.Context) (err error) {
	disk := Disk{}
	message, err := system.UpdateLink("objects/a0c17af0-45aa-4cdb-ace4-f60440583c8a", "parent", disk)
	if err != nil {
		return
	}

	return exec.ExecSync(ctx, message)
}
func deleteLinkById(ctx context.Context) (err error) {
	message, err := system.DeleteLink("21668", "")
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}
func deleteLinkByFromName(ctx context.Context) (err error) {
	message, err := system.DeleteLink("2aac7fcf-a3c5-4089-b873-11c0f52c488c", "nvme0n1")
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}
