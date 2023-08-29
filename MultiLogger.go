package golog

type MultiLogger struct {
	ID uintptr
	Children []Logger
}

func(logger *MultiLogger) Log(packet *Packet) {
	for _, child := range logger.Children {
		if child != nil {
			child.Log(packet)
		}
	}
}

func(logger *MultiLogger) Close() {
	for _, child := range logger.Children {
		if child != nil {
			child.Close()
		}
	}
}

func(logger *MultiLogger) SubLoggers() []Logger {
	var children []Logger
	for _, child := range logger.Children {
		if child != nil {
			children = append(children, child)
		}
	}
	return children
}

func(logger *MultiLogger) Identity() uintptr {
	return logger.ID
}

var _ Logger = &MultiLogger{}
