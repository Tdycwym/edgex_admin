package utils

import (
	"runtime/debug"

	"github.com/tydcwym/edgex_admin/logs"
)

// RecoverPanic ...
func RecoverPanic() {
	if x := recover(); x != nil {
		logs.Error("runtime panic: %v\n%v", x, string(debug.Stack()))
	}
}
