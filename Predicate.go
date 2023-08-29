package golog

import (
	"time"
)

type Predicate[SubjectT any] interface {
	Match(SubjectT) bool
}

type LevelPredicate struct {
	Predicate Predicate[Level]
	MissingResult bool
}

func(pred *LevelPredicate) Match(packet *Packet) bool {
	if packet == nil || packet.Level == nil || pred.Predicate == nil {
		return pred.MissingResult
	} else {
		return pred.Predicate.Match(packet.Level)
	}
}

type MessagePredicate struct {
	Predicate Predicate[Message]
	MissingResult bool
}

func(pred *MessagePredicate) Match(packet *Packet) bool {
	if packet == nil || packet.Message == nil || pred.Predicate == nil {
		return pred.MissingResult
	} else {
		return pred.Predicate.Match(packet.Message)
	}
}

type SourcePredicate struct {
	Predicate Predicate[Source]
	MissingResult bool
}

func(pred *SourcePredicate) Match(packet *Packet) bool {
	if packet == nil || packet.Source == nil || pred.Predicate == nil {
		return pred.MissingResult
	} else {
		return pred.Predicate.Match(packet.Source)
	}
}

type TimestampPredicate struct {
	Predicate Predicate[time.Time]
	MissingResult bool
}

func(pred *TimestampPredicate) Match(packet *Packet) bool {
	if packet == nil || packet.Timestamp.IsZero() || pred.Predicate == nil {
		return pred.MissingResult
	} else {
		return pred.Predicate.Match(packet.Timestamp)
	}
}

type LevelOrderPredicate struct {
	Threshold int
	Relation OrderRel
	MissingResult bool
}

func(pred *LevelOrderPredicate) Match(level Level) bool {
	if level == nil {
		return pred.MissingResult
	}
	actual := level.Numerical()
	switch pred.Relation {
		case ORDR_GREATER_EQUAL:
			return actual >= pred.Threshold
		case ORDR_GREATER:
			return actual > pred.Threshold
		case ORDR_LESS:
			return actual < pred.Threshold
		case ORDR_LESS_EQUAL:
			return actual <= pred.Threshold
		default:
			return actual == pred.Threshold
	}
}

type TruePredicate[SubjectT any] struct {}

func(pred TruePredicate[SubjectT]) Match(SubjectT) bool {
	return true
}

type FalsePredicate[SubjectT any] struct {}

func(pred FalsePredicate[SubjectT]) Match(SubjectT) bool {
	return false
}

type AllPredicate[SubjectT any] struct {
	Children []Predicate[SubjectT]
}

func(pred AllPredicate[SubjectT]) Match(subject SubjectT) bool {
	for _, child := range pred.Children {
		if child == nil {
			continue
		}
		if !child.Match(subject) {
			return false
		}
	}
	return true
}

type AnyPredicate[SubjectT any] struct {
	Children []Predicate[SubjectT]
}

func(pred AnyPredicate[SubjectT]) Match(subject SubjectT) bool {
	for _, child := range pred.Children {
		if child == nil {
			continue
		}
		if child.Match(subject) {
			return true
		}
	}
	return false
}

type NonePredicate[SubjectT any] struct {
	Children []Predicate[SubjectT]
}

func(pred NonePredicate[SubjectT]) Match(subject SubjectT) bool {
	for _, child := range pred.Children {
		if child == nil {
			continue
		}
		if child.Match(subject) {
			return false
		}
	}
	return true
}

var _ Predicate[*Packet] = &LevelPredicate{}
var _ Predicate[*Packet] = &MessagePredicate{}
var _ Predicate[*Packet] = &SourcePredicate{}
var _ Predicate[*Packet] = &TimestampPredicate{}

var _ Predicate[Level] = &LevelOrderPredicate{}

var _ Predicate[int] = TruePredicate[int]{}
var _ Predicate[int] = FalsePredicate[int]{}
var _ Predicate[int] = AllPredicate[int]{}
var _ Predicate[int] = AnyPredicate[int]{}
var _ Predicate[int] = NonePredicate[int]{}
