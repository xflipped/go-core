// Copyright 2022 Listware

package main

import (
	"context"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
)

func deleteObject(ctx context.Context) (err error) {
	message, err := system.DeleteObject("17177459")
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}
