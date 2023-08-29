package golog

type DispatchingLogger struct {
	ID uintptr
	Rules []*DispatchRule
}

type DispatchRule struct {
	Condition Predicate[*Packet]
	Logger Logger
	Continue Predicate[*Packet]
}

func(logger *DispatchingLogger) Log(packet *Packet) {
	for _, rule := range logger.Rules {
		if rule == nil {
			continue
		}
		if rule.Condition != nil && !rule.Condition.Match(packet) {
			continue
		}
		if rule.Logger != nil {
			rule.Logger.Log(packet)
		}
		if rule.Continue == nil || !rule.Continue.Match(packet) {
			break
		}
	}
}

func(logger *DispatchingLogger) Close() {
	for _, rule := range logger.Rules {
		if rule != nil && rule.Logger != nil {
			rule.Logger.Close()
		}
	}
}

func(logger *DispatchingLogger) SubLoggers() []Logger {
	var children []Logger
	for _, rule := range logger.Rules {
		if rule != nil && rule.Logger != nil {
			children = append(children, rule.Logger)
		}
	}
	return children
}

func(logger *DispatchingLogger) Identity() uintptr {
	return logger.ID
}

var _ Logger = &DispatchingLogger{}
