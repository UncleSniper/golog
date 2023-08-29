package golog

type NullLogger struct {
	ID uintptr
}

func(logger *NullLogger) Log(*Packet) {}

func(logger *NullLogger) Close() {}

func(logger *NullLogger) SubLoggers() []Logger {
	return nil
}

func(logger *NullLogger) Identity() uintptr {
	return logger.ID
}

var _ Logger = &NullLogger{}
