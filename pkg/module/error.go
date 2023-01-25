// Copyright 2022 Listware

package module

import (
	"fmt"
)

type errorRetry struct {
	err error
}

func (e *errorRetry) Error() string {
	return fmt.Sprintf("retry: %+v", e.err)
}

func (e *errorRetry) Unwrap() error {
	return e.err
}

func NewError(err error) error {
	return &errorRetry{err}
}
