// Copyright 2022 Listware

package errors

import (
	"fmt"
)

var (
	ErrPayloadNil     = fmt.Errorf("payload can't be empty")
	ErrBadCmdbContext = fmt.Errorf("bad cmdb context")
)

func ErrFunctionId(id string) error {
	return fmt.Errorf("unsupported function id: %s", id)
}
