package golog

import (
	"strings"
)

type DefaultLevel int

const (
	DEBUG DefaultLevel = iota
	CONFIG
	INFO
	WARNING
	ERROR
	MISUSE
	FATAL
)

func(level DefaultLevel) Numerical() int {
	return int(level)
}

func(level DefaultLevel) HumanReadable(adjust Adjustment) string {
	switch adjust {
		case ADJ_LEFT:
			switch level {
				case DEBUG:
					return "DEBUG  "
				case CONFIG:
					return "CONFIG "
				case INFO:
					return "INFO   "
				case WARNING:
					return "WARNING"
				case ERROR:
					return "ERROR  "
				case MISUSE:
					return "MISUSE "
				case FATAL:
					return "FATAL  "
				default:
					return "???    "
			}
		case ADJ_RIGHT:
			switch level {
				case DEBUG:
					return "  DEBUG"
				case CONFIG:
					return " CONFIG"
				case INFO:
					return "   INFO"
				case WARNING:
					return "WARNING"
				case ERROR:
					return "  ERROR"
				case MISUSE:
					return " MISUSE"
				case FATAL:
					return "  FATAL"
				default:
					return "    ???"
			}
		default:
			switch level {
				case DEBUG:
					return "DEBUG"
				case CONFIG:
					return "CONFIG"
				case INFO:
					return "INFO"
				case WARNING:
					return "WARNING"
				case ERROR:
					return "ERROR"
				case MISUSE:
					return "MISUSE"
				case FATAL:
					return "FATAL"
				default:
					return "???"
			}
	}
}

func(level DefaultLevel) IsNominal() bool {
	return level < WARNING
}

func ParseDefaultLevel(spec string) Level {
	switch strings.ToUpper(strings.TrimSpace(spec)) {
		case "DEBUG":
			return DEBUG
		case "CONFIG":
			return CONFIG
		case "INFO":
			return INFO
		case "WARNING":
			return WARNING
		case "ERROR":
			return ERROR
		case "MISUSE":
			return MISUSE
		case "FATAL":
			return FATAL
		default:
			return nil
	}
}

type StringMessage struct {
	Text []string
	Details Structure
}

func(msg *StringMessage) Lines() []string {
	return msg.Text
}

func(msg *StringMessage) PutStruct(sink StructSink) {
	if msg.Details != nil {
		msg.Details.PutStruct(sink)
	}
}

type DefaultSource struct {
	Module string
	Type string
	Function string
}

func(src *DefaultSource) StringSource() string {
	if src == nil {
		return ""
	}
	var builder strings.Builder
	var needDot bool
	if len(src.Module) > 0 {
		builder.WriteString(src.Module)
		needDot = true
	}
	if len(src.Type) > 0 {
		if needDot {
			builder.WriteRune('.')
		} else {
			needDot = true
		}
		builder.WriteString(src.Type)
	}
	if len(src.Function) > 0 {
		if needDot {
			builder.WriteRune('.')
		}
		builder.WriteString(src.Function)
	}
	return builder.String()
}

var _ Level = INFO
var _ Message = &StringMessage{}
var _ Source = &DefaultSource{}
