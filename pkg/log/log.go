// Copyright 2022 Listware

package log

import (
	"fmt"
)

// чтобы потом в одном месте сменить реализацию
func Debug(a ...any) {
	fmt.Println(a...)
}

func Debugf(format string, a ...any) {
	fmt.Printf(format, a...)
}
