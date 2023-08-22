package golog

import (
	"sync/atomic"
)

var loggerID atomic.Uintptr

func NewLoggerID() uintptr {
	return loggerID.Add(1)
}
