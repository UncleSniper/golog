package golog

import (
	"time"
	"strings"
)

type TextFormatter interface {
	PacketToText(*Packet) []string
}

type LineFormatter interface {
	PacketToLine(*Packet, *strings.Builder)
}

type StructFormatter interface {
	StructToText(Structure) string
}

func PacketToLine(packet *Packet, formatters []LineFormatter) string {
	var builder strings.Builder
	for _, formatter := range formatters {
		if formatter != nil {
			formatter.PacketToLine(packet, &builder)
		}
	}
	return builder.String()
}

type MessageTextFormatter struct {}

func(form MessageTextFormatter) PacketToText(packet *Packet) []string {
	if packet.Message == nil {
		return nil
	} else {
		return packet.Message.Lines()
	}
}

type ConcatTextFormatter struct {
	Formatters []TextFormatter
}

func(form ConcatTextFormatter) PacketToText(packet *Packet) []string {
	var all []string
	for _, formatter := range form.Formatters {
		if formatter == nil {
			continue
		}
		some := formatter.PacketToText(packet)
		if len(some) == 0 {
			continue
		}
		if all == nil {
			all = some
		} else {
			all = append(all, some...)
		}
	}
	return all
}

type LineTextFormatter struct {
	Formatters [][]LineFormatter
}

func(form LineTextFormatter) PacketToText(packet *Packet) []string {
	var all []string
	for _, outer := range form.Formatters {
		if outer == nil {
			continue
		}
		line := PacketToLine(packet, outer)
		all = append(all, line)
	}
	return all
}

type PrefixMode uint

const (
	PFX_ALL_SAME PrefixMode = iota
	PFX_THEN_SPACES
	PFX_ONLY_TOP
)

type PrefixedTextFormatter struct {
	Prefix []LineFormatter
	PrefixMode PrefixMode
	Lines TextFormatter
}

func(form *PrefixedTextFormatter) PacketToText(packet *Packet) []string {
	var lines []string
	if form.Lines != nil {
		lines = form.Lines.PacketToText(packet)
	}
	if len(lines) == 0 {
		return nil
	}
	topPrefix := PacketToLine(packet, form.Prefix)
	if len(topPrefix) == 0 {
		return lines
	}
	var restPrefix string
	switch form.PrefixMode {
		case PFX_ALL_SAME:
			restPrefix = topPrefix
		case PFX_THEN_SPACES:
			restPrefix = string(RepeatRune(' ', len(topPrefix)))
	}
	for index := 0; index < len(lines); index++ {
		if index == 0 {
			lines[0] = topPrefix + lines[0]
		} else {
			lines[index] = restPrefix + lines[index]
		}
	}
	return lines
}

type ConcatLineFormatter struct {
	Formatters []LineFormatter
}

func(form ConcatLineFormatter) PacketToLine(packet *Packet, builder *strings.Builder) {
	for _, formatter := range form.Formatters {
		if formatter != nil {
			formatter.PacketToLine(packet, builder)
		}
	}
}

type StringLineFormatter struct {
	Value string
}

func(form StringLineFormatter) PacketToLine(packet *Packet, builder *strings.Builder) {
	builder.WriteString(form.Value)
}

type AffixFlags uint

const (
	AFF_PREFIX_IF_MISSING AffixFlags = 1 << iota
	AFF_PREFIX_IF_EMPTY
	AFF_SUFFIX_IF_MISSING
	AFF_SUFFIX_IF_EMPTY
)

type PieceLineFormatterBase struct {
	Prefix LineFormatter
	Suffix LineFormatter
	Flags AffixFlags
	ReplacementIfMissing LineFormatter
	ReplacementIfEmpty LineFormatter
}

func(form *PieceLineFormatterBase) Missing(packet *Packet, builder *strings.Builder) {
	if form.Prefix != nil && form.Flags & AFF_PREFIX_IF_MISSING != 0 {
		form.Prefix.PacketToLine(packet, builder)
	}
	if form.ReplacementIfMissing != nil {
		form.ReplacementIfMissing.PacketToLine(packet, builder)
	}
	if form.Suffix != nil && form.Flags & AFF_SUFFIX_IF_MISSING != 0 {
		form.Suffix.PacketToLine(packet, builder)
	}
}

func(form *PieceLineFormatterBase) Empty(packet *Packet, builder *strings.Builder) {
	if form.Prefix != nil && form.Flags & AFF_PREFIX_IF_EMPTY != 0 {
		form.Prefix.PacketToLine(packet, builder)
	}
	if form.ReplacementIfEmpty != nil {
		form.ReplacementIfEmpty.PacketToLine(packet, builder)
	}
	if form.Suffix != nil && form.Flags & AFF_SUFFIX_IF_EMPTY != 0 {
		form.Suffix.PacketToLine(packet, builder)
	}
}

func(form *PieceLineFormatterBase) WithString(rendition string, packet *Packet, builder *strings.Builder) {
	if len(rendition) == 0 {
		form.Empty(packet, builder)
		return
	}
	if form.Prefix != nil {
		form.Prefix.PacketToLine(packet, builder)
	}
	builder.WriteString(rendition)
	if form.Suffix != nil {
		form.Suffix.PacketToLine(packet, builder)
	}
}

type GenericLevelLineFormatter struct {
	PieceLineFormatterBase
	Adjustment Adjustment
}

func(form *GenericLevelLineFormatter) PacketToLine(packet *Packet, builder *strings.Builder) {
	if packet.Level == nil {
		form.Missing(packet, builder)
	} else {
		form.WithString(packet.Level.HumanReadable(form.Adjustment), packet, builder)
	}
}

type GenericSourceLineFormatter struct {
	PieceLineFormatterBase
}

func(form *GenericSourceLineFormatter) PacketToLine(packet *Packet, builder *strings.Builder) {
	if packet.Source == nil {
		form.Missing(packet, builder)
	} else {
		form.WithString(packet.Source.StringSource(), packet, builder)
	}
}

type GenericTimestampLineFormatter struct {
	PieceLineFormatterBase
	Format string
}

func(form *GenericTimestampLineFormatter) PacketToLine(packet *Packet, builder *strings.Builder) {
	if packet.Timestamp.IsZero() {
		form.Missing(packet, builder)
	} else {
		format := form.Format
		if len(format) == 0 {
			format = time.DateTime
		}
		form.WithString(packet.Timestamp.Format(format), packet, builder)
	}
}

type GenericStructLineFormatter struct {
	PieceLineFormatterBase
	Formatter StructFormatter
}

func(form *GenericStructLineFormatter) PacketToLine(packet *Packet, builder *strings.Builder) {
	if packet.Message == nil {
		form.Missing(packet, builder)
	} else {
		formatter := form.Formatter
		if formatter == nil {
			if DumbStructFormatter == nil {
				formatter = TextStructFormatter{}
			} else {
				formatter = DumbStructFormatter
			}
		}
		form.WithString(formatter.StructToText(packet.Message), packet, builder)
	}
}

type TextStructFormatter struct {
	KeepOutermostParens bool
}

func(form TextStructFormatter) StructToText(structure Structure) string {
	if structure == nil {
		return ""
	}
	sink := &TextStructSink {
		KeepOutermostParens: form.KeepOutermostParens,
	}
	structure.PutStruct(sink)
	return sink.ToString()
}

var DumbStructFormatter StructFormatter = TextStructFormatter{}

var _ TextFormatter = MessageTextFormatter{}
var _ TextFormatter = ConcatTextFormatter{}
var _ TextFormatter = LineTextFormatter{}
var _ TextFormatter = &PrefixedTextFormatter{}

var _ LineFormatter = ConcatLineFormatter{}
var _ LineFormatter = StringLineFormatter{}
var _ LineFormatter = &GenericLevelLineFormatter{}
var _ LineFormatter = &GenericSourceLineFormatter{}
var _ LineFormatter = &GenericTimestampLineFormatter{}
var _ LineFormatter = &GenericStructLineFormatter{}

var _ StructFormatter = TextStructFormatter{}
