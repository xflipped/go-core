// Copyright 2022 Listware

package main

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/vertex/types"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
)

type Node struct {
	Hostname string `json:"hostname,omitempty"`
	Domain   string `json:"domain"`
	Model    string `json:"model"`
}

func createNodeType(ctx context.Context) (err error) {
	node := Node{
		Hostname: "sky01",
	}

	pt := types.ReflectType(node)
	//	pt.Triggers["create"] = append(pt.Triggers["create"], logTriger)

	message, err := system.CreateType(pt)
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}

func updateNodeType(ctx context.Context) (err error) {
	node := Node{
		Hostname: "sky01",
	}

	pt := types.ReflectType(node)
	//	pt.Triggers["create"] = append(pt.Triggers["create"], logTriger)

	message, err := system.UpdateType(pt)
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}

func deleteNodeType(ctx context.Context) (err error) {
	message, err := system.DeleteType("node")
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}

func createNodeObject(ctx context.Context) (err error) {
	node := Node{
		Hostname: "sky01",
	}

	message, err := system.CreateObject("node", node)
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}

func createChild(ctx context.Context) (err error) {
	disk := Disk{
		Name: "nvme4n1",
	}

	message, err := system.CreateChild("objects/64f21ca8-c2c7-4fc1-bccc-de04842dbd49", "types/disk", disk.Name, disk)
	if err != nil {
		return
	}

	return exec.ExecSync(ctx, message)
}

func updateObject(ctx context.Context) (err error) {
	node := Node{
		Hostname: "sky0asd112",
	}

	message, err := system.UpdateObject("58d12e3e-63f7-4c90-9d28-3812aebf81ce", node)
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}

func createNodeTypeTrigger(ctx context.Context) (err error) {
	message, err := system.AddTrigger("types/node", logTrigger)
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}
func deleteNodeTypeTrigger(ctx context.Context) (err error) {
	message, err := system.DeleteTrigger("types/node", logTrigger)
	if err != nil {
		return
	}

	return exec.ExecAsync(ctx, message)
}
