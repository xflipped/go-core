// Copyright 2023 NJWS Inc.

package module

import (
	"encoding/json"
	"io"

	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
)

var GenericJsonType = &genericJsonType{}

type genericJsonType struct{}

func (g *genericJsonType) GetTypeName() statefun.TypeName {
	return statefun.TypeNameFrom("generic/json")
}

func (g *genericJsonType) Deserialize(r io.Reader, receiver interface{}) error {
	return json.NewDecoder(r).Decode(receiver)
}

func (g *genericJsonType) Serialize(w io.Writer, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
