// Copyright 2023 NJWS Inc.

package module

import (
	"testing"

	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
)

// TestWithValueSpec do not use reserved names
func TestWithValueSpec(t *testing.T) {
	var m functionMiddleware

	if err := WithValueSpec(syncTableSpec)(&m); err != ErrReservedName {
		t.Fatalf("expected %+v, but have %+v", ErrReservedName, err)
	}

	if err := WithValueSpec(statefun.ValueSpec{})(&m); err != ErrEmptyName {
		t.Fatalf("expected %+v, but have %+v", ErrEmptyName, err)
	}

	if err := WithValueSpec(statefun.ValueSpec{Name: "name"})(&m); err != nil {
		t.Fatalf("expected %+v, but have %+v", nil, err)
	}

}
