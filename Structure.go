package golog

import (
	"fmt"
	"strings"
)

type StructSink interface {
	Map() StructMap
	List() StructList
}

type StructMap interface {
	BoolProperty(string, bool)
	StringProperty(string, string)
	IntProperty(string, int64)
	FloatProperty(string, float64)
	MapProperty(string) StructMap
	ListProperty(string) StructList
	EndMap()
}

type StructList interface {
	StructSink
	Bool(bool)
	String(string)
	Int(int64)
	Float(float64)
	EndList()
}

type Structure interface {
	PutStruct(StructSink)
}

type TextStructSink struct {
	builder strings.Builder
	stack BoolStack
	KeepOutermostParens bool
}

func(sink *TextStructSink) Map() StructMap {
	if !sink.enterElement() || sink.KeepOutermostParens {
		sink.builder.WriteRune('{')
	}
	sink.stack.Push(false)
	return sink
}

func(sink *TextStructSink) List() StructList {
	if !sink.enterElement() || sink.KeepOutermostParens {
		sink.builder.WriteRune('[')
	}
	sink.stack.Push(false)
	return sink
}

func(sink *TextStructSink) enterElement() bool {
	if sink.stack.IsEmpty() {
		return true
	}
	if sink.stack.Top() {
		sink.builder.WriteString(", ")
	} else {
		sink.stack.Replace(true)
	}
	return false
}

func(sink *TextStructSink) BoolProperty(name string, value bool) {
	sink.enterElement()
	sink.builder.WriteString(fmt.Sprintf("%s: %v", name, value))
}

func(sink *TextStructSink) StringProperty(name string, value string) {
	sink.enterElement()
	sink.builder.WriteString(fmt.Sprintf("%s: %q", name, value))
}

func(sink *TextStructSink) IntProperty(name string, value int64) {
	sink.enterElement()
	sink.builder.WriteString(fmt.Sprintf("%s: %d", name, value))
}

func(sink *TextStructSink) FloatProperty(name string, value float64) {
	sink.enterElement()
	sink.builder.WriteString(fmt.Sprintf("%s: %G", name, value))
}

func(sink *TextStructSink) MapProperty(name string) StructMap {
	sink.enterElement()
	sink.builder.WriteString(fmt.Sprintf("%s: {", name))
	sink.stack.Push(false)
	return sink
}

func(sink *TextStructSink) ListProperty(name string) StructList {
	sink.enterElement()
	sink.builder.WriteString(fmt.Sprintf("%s: [", name))
	sink.stack.Push(false)
	return sink
}

func(sink *TextStructSink) EndMap() {
	sink.stack.Pop()
	if sink.KeepOutermostParens || !sink.stack.IsEmpty() {
		sink.builder.WriteRune('}')
	}
}

func(sink *TextStructSink) Bool(value bool) {
	sink.enterElement()
	if value {
		sink.builder.WriteString("true")
	} else {
		sink.builder.WriteString("false")
	}
}

func(sink *TextStructSink) String(value string) {
	sink.enterElement()
	sink.builder.WriteString(fmt.Sprintf("%q", value))
}

func(sink *TextStructSink) Int(value int64) {
	sink.enterElement()
	sink.builder.WriteString(fmt.Sprintf("%d", value))
}

func(sink *TextStructSink) Float(value float64) {
	sink.enterElement()
	sink.builder.WriteString(fmt.Sprintf("%G", value))
}

func(sink *TextStructSink) EndList() {
	sink.stack.Pop()
	if sink.KeepOutermostParens || !sink.stack.IsEmpty() {
		sink.builder.WriteRune(']')
	}
}

func(sink *TextStructSink) ToString() string {
	return sink.builder.String()
}

var _ StructSink = &TextStructSink{}
