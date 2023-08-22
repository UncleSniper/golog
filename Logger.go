package golog

import (
	"fmt"
	"time"
)

type Logger interface {
	Log(*Packet)
	Close()
	SubLoggers() []Logger
	Identity() uintptr
}

type Log struct {
	Logger Logger
}

func(log *Log) Logp(packet *Packet) {
	if log.Logger == nil {
		return
	}
	if packet != nil && packet.Timestamp.IsZero() {
		packet.Timestamp = time.Now()
	}
	log.Logger.Log(packet)
}

func(log *Log) Log(level Level, src Source, msg Message) {
	if log.Logger != nil {
		log.Logger.Log(&Packet {
			Level: level,
			Message: msg,
			Source: src,
			Timestamp: time.Now(),
		})
	}
}

func(log *Log) Logv(level Level, src Source, details Structure, args ...any) {
	log.Log(level, src, &StringMessage {
		Text: []string { fmt.Sprint(args...) },
		Details: details,
	})
}

func(log *Log) Logf(level Level, src Source, details Structure, format string, args ...any) {
	log.Log(level, src, &StringMessage {
		Text: []string { fmt.Sprintf(format, args...) },
		Details: details,
	})
}

func(log *Log) Debug(src Source, msg Message) {
	log.Log(DEBUG, src, msg)
}

func(log *Log) Debugv(src Source, details Structure, args ...any) {
	log.Logv(DEBUG, src, details, args...)
}

func(log *Log) Debugf(src Source, details Structure, format string, args ...any) {
	log.Logf(DEBUG, src, details, format, args...)
}

func(log *Log) Config(src Source, msg Message) {
	log.Log(CONFIG, src, msg)
}

func(log *Log) Configv(src Source, details Structure, args ...any) {
	log.Logv(CONFIG, src, details, args...)
}

func(log *Log) Configf(src Source, details Structure, format string, args ...any) {
	log.Logf(CONFIG, src, details, format, args...)
}

func(log *Log) Info(src Source, msg Message) {
	log.Log(INFO, src, msg)
}

func(log *Log) Infov(src Source, details Structure, args ...any) {
	log.Logv(INFO, src, details, args...)
}

func(log *Log) Infof(src Source, details Structure, format string, args ...any) {
	log.Logf(INFO, src, details, format, args...)
}

func(log *Log) Warn(src Source, msg Message) {
	log.Log(WARNING, src, msg)
}

func(log *Log) Warnv(src Source, details Structure, args ...any) {
	log.Logv(WARNING, src, details, args...)
}

func(log *Log) Warnf(src Source, details Structure, format string, args ...any) {
	log.Logf(WARNING, src, details, format, args...)
}

func(log *Log) Error(src Source, msg Message) {
	log.Log(ERROR, src, msg)
}

func(log *Log) Errorv(src Source, details Structure, args ...any) {
	log.Logv(ERROR, src, details, args...)
}

func(log *Log) Errorf(src Source, details Structure, format string, args ...any) {
	log.Logf(ERROR, src, details, format, args...)
}

func(log *Log) Misuse(src Source, msg Message) {
	log.Log(MISUSE, src, msg)
}

func(log *Log) Misusev(src Source, details Structure, args ...any) {
	log.Logv(MISUSE, src, details, args...)
}

func(log *Log) Misusef(src Source, details Structure, format string, args ...any) {
	log.Logf(MISUSE, src, details, format, args...)
}

func(log *Log) Fatal(src Source, msg Message) {
	log.Log(FATAL, src, msg)
}

func(log *Log) Fatalv(src Source, details Structure, args ...any) {
	log.Logv(FATAL, src, details, args...)
}

func(log *Log) Fatalf(src Source, details Structure, format string, args ...any) {
	log.Logf(FATAL, src, details, format, args...)
}

type BoundLog struct {
	Logger Logger
	Source Source
}

func(log *BoundLog) Logp(packet *Packet) {
	if log.Logger == nil {
		return
	}
	if packet != nil {
		if packet.Source == nil && log.Source != nil {
			packet.Source = log.Source
		}
		if packet.Timestamp.IsZero() {
			packet.Timestamp = time.Now()
		}
	}
	log.Logger.Log(packet)
}

func(log *BoundLog) Log(level Level, msg Message) {
	if log.Logger != nil {
		log.Logger.Log(&Packet {
			Level: level,
			Message: msg,
			Source: log.Source,
			Timestamp: time.Now(),
		})
	}
}

func(log *BoundLog) Logv(level Level, details Structure, args ...any) {
	log.Log(level, &StringMessage {
		Text: []string { fmt.Sprint(args...) },
		Details: details,
	})
}

func(log *BoundLog) Logf(level Level, details Structure, format string, args ...any) {
	log.Log(level, &StringMessage {
		Text: []string { fmt.Sprintf(format, args...) },
		Details: details,
	})
}

func(log *BoundLog) Debug(msg Message) {
	log.Log(DEBUG, msg)
}

func(log *BoundLog) Debugv(details Structure, args ...any) {
	log.Logv(DEBUG, details, args...)
}

func(log *BoundLog) Debugf(details Structure, format string, args ...any) {
	log.Logf(DEBUG, details, format, args...)
}

func(log *BoundLog) Config(msg Message) {
	log.Log(CONFIG, msg)
}

func(log *BoundLog) Configv(details Structure, args ...any) {
	log.Logv(CONFIG, details, args...)
}

func(log *BoundLog) Configf(details Structure, format string, args ...any) {
	log.Logf(CONFIG, details, format, args...)
}

func(log *BoundLog) Info(msg Message) {
	log.Log(INFO, msg)
}

func(log *BoundLog) Infov(details Structure, args ...any) {
	log.Logv(INFO, details, args...)
}

func(log *BoundLog) Infof(details Structure, format string, args ...any) {
	log.Logf(INFO, details, format, args...)
}

func(log *BoundLog) Warn(msg Message) {
	log.Log(WARNING, msg)
}

func(log *BoundLog) Warnv(details Structure, args ...any) {
	log.Logv(WARNING, details, args...)
}

func(log *BoundLog) Warnf(details Structure, format string, args ...any) {
	log.Logf(WARNING, details, format, args...)
}

func(log *BoundLog) Error(msg Message) {
	log.Log(ERROR, msg)
}

func(log *BoundLog) Errorv(details Structure, args ...any) {
	log.Logv(ERROR, details, args...)
}

func(log *BoundLog) Errorf(details Structure, format string, args ...any) {
	log.Logf(ERROR, details, format, args...)
}

func(log *BoundLog) Misuse(msg Message) {
	log.Log(MISUSE, msg)
}

func(log *BoundLog) Misusev(details Structure, args ...any) {
	log.Logv(MISUSE, details, args...)
}

func(log *BoundLog) Misusef(details Structure, format string, args ...any) {
	log.Logf(MISUSE, details, format, args...)
}

func(log *BoundLog) Fatal(msg Message) {
	log.Log(FATAL, msg)
}

func(log *BoundLog) Fatalv(details Structure, args ...any) {
	log.Logv(FATAL, details, args...)
}

func(log *BoundLog) Fatalf(details Structure, format string, args ...any) {
	log.Logf(FATAL, details, format, args...)
}

type LoggerWalker interface {
	EnterLogger(Logger, int, uint, bool) bool
	LeaveLogger(Logger, int, uint, bool)
	SkipLogger(logger Logger, index int, depth uint, leaf bool, seen bool)
}

func WalkLoggers(root Logger, walker LoggerWalker) {
	if walker == nil {
		return
	}
	seen := make(map[uintptr]bool)
	walkLoggersRec(root, walker, -1, 0, seen)
}

func walkLoggersRec(logger Logger, walker LoggerWalker, index int, depth uint, seen map[uintptr]bool) {
	var id uintptr
	if logger != nil {
		id = logger.Identity()
	}
	var children []Logger
	if logger != nil {
		children = logger.SubLoggers()
	}
	leaf := len(children) == 0
	if id != 0 && seen[id] {
		walker.SkipLogger(logger, index, depth, leaf, true)
		return
	}
	skip := walker.EnterLogger(logger, index, depth, leaf)
	if skip {
		walker.SkipLogger(logger, index, depth, leaf, false)
		return
	}
	if id != 0 {
		seen[id] = true
	}
	for childIndex, child := range children {
		walkLoggersRec(child, walker, childIndex, depth + 1, seen)
	}
	walker.LeaveLogger(logger, index, depth, leaf)
}
