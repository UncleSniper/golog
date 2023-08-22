package golog

import (
	"time"
)

type Adjustment uint

const (
	ADJ_NONE Adjustment = iota
	ADJ_LEFT
	ADJ_RIGHT
)

type Level interface {
	Numerical() int
	HumanReadable(Adjustment) string
	IsNominal() bool
}

type Message interface {
	Structure
	Lines() []string
}

type Source interface {
	StringSource() string
}

type Packet struct {
	Level Level
	Message Message
	Source Source
	Timestamp time.Time
}
