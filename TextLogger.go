package golog

import (
	"os"
	"fmt"
	"sync"
)

type TextLogger struct {
	ID uintptr
	WriteInfo func(string)
	WriteError func(string)
	CloseStream func()
	Formatter TextFormatter
	mutex sync.Mutex
}

func(logger *TextLogger) Log(packet *Packet) {
	if packet == nil || packet.Message == nil {
		return
	}
	var lines []string
	if logger.Formatter == nil {
		lines = packet.Message.Lines()
	} else {
		lines = logger.Formatter.PacketToText(packet)
	}
	if len(lines) == 0 {
		return
	}
	var writer func(string)
	if packet.Level == nil || packet.Level.IsNominal() {
		if logger.WriteInfo != nil {
			writer = logger.WriteInfo
		} else {
			writer = logger.WriteError
		}
	} else {
		if logger.WriteError != nil {
			writer = logger.WriteError
		} else {
			writer = logger.WriteInfo
		}
	}
	if writer == nil {
		return
	}
	logger.mutex.Lock()
	for _, line := range lines {
		writer(line)
	}
	logger.mutex.Unlock()
}

func(logger *TextLogger) Close() {
	if logger.CloseStream != nil {
		logger.CloseStream()
		logger.CloseStream = nil
	}
}

func(logger *TextLogger) SubLoggers() []Logger {
	return nil
}

func(logger *TextLogger) Identity() uintptr {
	return logger.ID
}

func TextFileLogger(path string, formatter TextFormatter) (*TextLogger, error) {
	f, err := os.OpenFile(path, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &TextLogger {
		ID: NewLoggerID(),
		WriteInfo: func(line string) {
			f.WriteString(line)
		},
		CloseStream: func() {
			f.Close()
		},
		Formatter: formatter,
	}, nil
}

func WriteLineToStdout(line string) {
	fmt.Println(line)
}

func WriteLineToStderr(line string) {
	fmt.Fprintln(os.Stderr, line)
}

var DumbLogger Logger = &TextLogger {
	ID: NewLoggerID(),
	WriteInfo: WriteLineToStdout,
	WriteError: WriteLineToStderr,
}
