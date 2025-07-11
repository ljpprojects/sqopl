package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type InstructionKind int

const (
	Invalid InstructionKind = -2
	Unknown InstructionKind = -1
	Repeat  InstructionKind = iota
	With
	Queueinit
	Reserve
	Grow
	Shrink
	Extend
	Ref
	Run
	Push
	Add
	Sub
	Mul
	Div
	Rem
	Pow
	Subroutine
	Dump
	Drain
	DumpN
	DrainN
	Goto
	GotoPoint
	GotoCondition
	Equal
	Inequal
	GreaterThan
	LessThan
	GreaterThanEqual
	LessThanEqual
	Point
	Zero
	Fill
	Configure
)

type ObjectAddress uint16

type Instruction interface {
	kind() InstructionKind
	operands() any
}

type RepeatInstruction struct {
	repititions uint8
	code        []Instruction
}

func (i RepeatInstruction) kind() InstructionKind {
	return Repeat
}

func (i RepeatInstruction) operands() any {
	return []any{i.repititions, i.code}
}

type WithInstruction struct {
	queueIndex uint8
	code       []Instruction
}

func (i WithInstruction) kind() InstructionKind {
	return With
}

func (i WithInstruction) operands() any {
	return []any{i.queueIndex, i.code}
}

type QueueInitInstruction struct{}

func (i QueueInitInstruction) kind() InstructionKind {
	return Queueinit
}

func (i QueueInitInstruction) operands() any {
	return nil
}

type ReserveInstruction struct {
	bytes uint16
}

func (i ReserveInstruction) kind() InstructionKind {
	return Reserve
}

func (i ReserveInstruction) operands() any {
	return i.bytes
}

type GrowInstruction struct {
	bytes int16
}

func (i GrowInstruction) kind() InstructionKind {
	return Grow
}

func (i GrowInstruction) operands() any {
	return i.bytes
}

type ShrinkInstruction struct {
	bytes int16
}

func (i ShrinkInstruction) kind() InstructionKind {
	return Shrink
}

func (i ShrinkInstruction) operands() any {
	return i.bytes
}

type ExtendInstruction struct {
	arrayReference ObjectAddress
}

func (i ExtendInstruction) kind() InstructionKind {
	return Shrink
}

func (i ExtendInstruction) operands() any {
	return i.arrayReference
}

type PushInstruction struct {
	data uint8
}

func (i PushInstruction) kind() InstructionKind {
	return Push
}

func (i PushInstruction) operands() any {
	return i.data
}

type DumpInstructionFormat uint8

const (
	DumpInstrRawFormat DumpInstructionFormat = iota
	DumpInstrTextFormat
)

type DumpInstruction struct {
	format DumpInstructionFormat

	/// If this is -1, dump everything
	count int32

	concat bool
}

func (i DumpInstruction) kind() InstructionKind {
	if i.count == -1 {
		return Dump
	} else {
		return DumpN
	}
}

func (i DumpInstruction) operands() any {
	return []any{i.format, i.count, i.concat}
}

type DrainInstruction struct {
	/// If this is -1, drain everything
	count int32
}

func (i DrainInstruction) kind() InstructionKind {
	if i.count == -1 {
		return Drain
	} else {
		return DrainN
	}
}

func (i DrainInstruction) operands() any {
	return i.count
}

type ConfigureInstrMode uint8

const (
	ConfigureInstrEnableMode ConfigureInstrMode = iota
	ConfigureInstrDisableMode
)

type ConfigureInstrKey uint8

const (
	ConfigureInstrAutoExpandQueuesKey ConfigureInstrKey = iota
	ConfigureInstrAutoShrinkQueuesKey
)

type ConfigureInstruction struct {
	mode ConfigureInstrMode
	key  ConfigureInstrKey
}

func (i ConfigureInstruction) kind() InstructionKind {
	return Configure
}

// {ConfigureInstrMode, ConfigureInstrKey}
func (i ConfigureInstruction) operands() any {
	return []any{i.mode, i.key}
}

type ArithmeticInstruction struct {
	_kind InstructionKind
}

func (i ArithmeticInstruction) kind() InstructionKind {
	return i._kind
}

func (i ArithmeticInstruction) operands() any {
	return nil
}

/**** FILLER FUNCTIONS ****/

func arr_map[A any, B any](original []A, predicate func(A) B) []B {
	var new []B

	for _, v := range original {
		new = append(new, predicate(v))
	}

	return new
}

func arr_filter[A any](original []A, predicate func(A) bool) []A {
	var new []A

	for _, v := range original {
		if predicate(v) {
			new = append(new, v)
		}
	}

	return new
}

func arr_flat_map[A any, B any](original []A, predicate func(A) []B) []B {
	var new []B

	for _, v := range original {
		for _, w := range predicate(v) {
			new = append(new, w)
		}
	}

	return new
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func joinNumericSlice[N Numeric](numbers []N) string {
	s := make([]string, len(numbers))

	for i, n := range numbers {
		s[i] = fmt.Sprint(n)
	}
	return strings.Join(s, ", ")
}

/**** LEXER + PARSER ****/

var ErrInvalidWithSpecifier = fmt.Errorf("Invalid queue given for 'with' statement, expected specifier to match /q[0-9]+/")
var ErrUnknownInstruction = fmt.Errorf("Unknown instruction")
var ErrInvalidConfigureMode = fmt.Errorf("Invalid configure mode for configure instruction.")
var ErrInvalidConfigureKey = fmt.Errorf("Invalid configuration key for configure instruction.")

var li uint32
var ci uint32
var depth uint8
var lines []string

func init() {
	li = 0
	ci = 0
	depth = 0
	lines = []string{}
}

func ParseInteger() int32 {
	__n := ""

	if lines[li][ci] == '\'' {
		ci++

		c := int32(rune(lines[li][ci]))

		ci++

		return c
	}

	for len([]byte(lines[li])) > int(ci) && unicode.IsDigit(rune(lines[li][ci])) {
		__n = __n + string(lines[li][ci])

		ci++
	}

	_n, _ := strconv.Atoi(__n)

	return int32(_n)
}

func ParseInstruction() (Instruction, error) {
	oneShotInstructions := map[string]Instruction{
		"queueinit": QueueInitInstruction{},
		"add":       ArithmeticInstruction{_kind: Add},
		"sub":       ArithmeticInstruction{_kind: Sub},
		"mul":       ArithmeticInstruction{_kind: Mul},
		"div":       ArithmeticInstruction{_kind: Div},
		"rem":       ArithmeticInstruction{_kind: Rem},
		"pow":       ArithmeticInstruction{_kind: Pow},
	}

	for k, v := range oneShotInstructions {
		if strings.HasPrefix(lines[li][ci:], k) {
			return v, nil
		}
	}

	if strings.HasPrefix(lines[li][ci:], "cfg") {
		ci += 3

		var mode ConfigureInstrMode
		var key ConfigureInstrKey

		if strings.HasPrefix(lines[li][ci:], ".enable") {
			mode = ConfigureInstrEnableMode
			ci += 7
		} else if strings.HasPrefix(lines[li][ci:], ".disable") {
			mode = ConfigureInstrDisableMode
			ci += 8
		} else {
			return nil, ErrInvalidConfigureMode
		}

		ci++

		var validConfigurationKeys = map[string]ConfigureInstrKey{
			"autoExpandQueues": ConfigureInstrAutoExpandQueuesKey,
			"autoShrinkQueues": ConfigureInstrAutoShrinkQueuesKey,
		}

		for s, k := range validConfigurationKeys {
			if strings.HasPrefix(lines[li][ci:], s) {
				key = k

				ci += uint32(len([]byte(s)))

				break
			}
		}

		return ConfigureInstruction{
			key:  key,
			mode: mode,
		}, nil
	} else if strings.HasPrefix(lines[li][ci:], "drain") {
		ci += 5

		count := int32(-1)

		if strings.HasPrefix(lines[li][ci:], ".") {
			ci++

			count = ParseInteger()
		}

		ci++

		return DrainInstruction{
			count: count,
		}, nil
	} else if strings.HasPrefix(lines[li][ci:], "push") {
		ci += 5

		n := uint8(ParseInteger())

		return PushInstruction{
			data: n,
		}, nil
	} else if strings.HasPrefix(lines[li][ci:], "dump") {
		ci += 4

		format := DumpInstrRawFormat
		count := int32(-1)
		concat := false

		if strings.HasPrefix(lines[li][ci:], ".text") {
			format = DumpInstrTextFormat
			ci += 5
		}

		if strings.HasPrefix(lines[li][ci:], ".concat") {
			ci += 7

			concat = true
		}

		if strings.HasPrefix(lines[li][ci:], ".") {
			ci++

			count = ParseInteger()
		}

		ci++

		return DumpInstruction{
			format: format,
			count:  count,
			concat: concat,
		}, nil
	} else if strings.HasPrefix(lines[li][ci:], "reserve") {
		ci += 8

		n := uint16(ParseInteger())

		return ReserveInstruction{
			bytes: n,
		}, nil
	} else if strings.HasPrefix(lines[li][ci:], "with") {
		ci += 5

		if lines[li][ci] != 'q' {
			return nil, ErrInvalidWithSpecifier
		}

		index := uint8(ParseInteger())

		depth++

		li++
		ci = 0

		body, err := ParseBody()

		if err != nil {
			return nil, err
		}

		depth--

		return WithInstruction{
			queueIndex: index,
			code:       body,
		}, nil
	} else if strings.HasPrefix(lines[li][ci:], "repeat") {
		ci += 7

		n := uint8(ParseInteger())

		depth++

		li++
		ci = 0

		body, err := ParseBody()

		if err != nil {
			return nil, err
		}

		depth--

		return RepeatInstruction{
			repititions: n,
			code:        body,
		}, nil
	}

	fmt.Println(lines[li][ci:])

	return nil, ErrUnknownInstruction
}

func ParseBody() ([]Instruction, error) {
	code := []Instruction{}

	for int(li) < len(lines) && (depth == 0 || (string(lines[li][ci:depth+1]) == (strings.Repeat(">", int(depth)) + " "))) {
		ci += uint32(depth) + uint32(boolToInt(depth > 0))

		inst, err := ParseInstruction()

		if err != nil {
			return nil, err
		}

		code = append(code, inst)

		li++
		ci = 0
	}

	li--
	ci = 0

	return code, nil
}

/**** INTERPETER ****/

type Numeric interface {
	~uint8 | ~int8 | ~uint16 | ~int16 | ~uint32 | ~int32 | ~uint64 | ~int64 | ~int | ~uint | ~float32 | ~float64
}

type Range[T Numeric] struct {
	Start T
	End   T
}

const MaxAllowedObjectSpace uint16 = 1024

var ErrQueueOverflow = fmt.Errorf("Queue overflow")
var ErrQueueOveflowsObjectSpace = fmt.Errorf("Queue would overflow object space if expanded")
var ErrCannotWriteOutsideOfCurrentQueue = fmt.Errorf("Cannot modify data outside of the current queue")
var ErrAutoExpandQueueNotCompatibleWithReserve = fmt.Errorf("The reserve instruction cannot be run with autoExpandQueue configured to enabled")

type ConfigurationState uint8

const (
	DisabledConfigurationState ConfigurationState = iota
	EnabledConfigurationState
)

type InterpreterSubroutine struct {
	parentContext *InterpreterStateContext
	instructions  []Instruction
}

type InterpreterStateContext struct {
	currentQueue *Range[ObjectAddress]
	subroutines  map[string]InterpreterSubroutine
}

type InterpreterState struct {
	data           [MaxAllowedObjectSpace]uint8
	queueBounds    []Range[ObjectAddress]
	currentContext InterpreterStateContext
	qi             ObjectAddress
	configuration  map[ConfigureInstrKey]ConfigurationState
}

func InterpretConfigure(instr Instruction, state *InterpreterState) error {
	mode := (instr.operands().([]any))[0].(ConfigureInstrMode)
	key := (instr.operands().([]any))[1].(ConfigureInstrKey)

	switch mode {
	case ConfigureInstrEnableMode:
		state.configuration[key] = EnabledConfigurationState
	case ConfigureInstrDisableMode:
		state.configuration[key] = DisabledConfigurationState
	}

	return nil
}

func InterpretQueueinit(_instr Instruction, state *InterpreterState) error {
	var startAddr ObjectAddress = 0

	if len(state.queueBounds) > 0 {
		startAddr = ObjectAddress(state.queueBounds[len(state.queueBounds)-1].End + 1)
	}

	state.queueBounds = append(state.queueBounds, Range[ObjectAddress]{
		Start: startAddr,
		End:   startAddr,
	})

	return nil
}

func InterpretWith(instr Instruction, state *InterpreterState) error {
	queueIndex := (instr.operands().([]any))[0].(uint8)
	code := (instr.operands().([]any))[1].([]Instruction)

	state.currentContext.currentQueue = &state.queueBounds[queueIndex]

	err := Interpet(code, state)

	if err != nil {
		return err
	}

	return nil
}

func InterpretRepeat(instr Instruction, state *InterpreterState) error {
	repititions := (instr.operands().([]any))[0].(uint8)
	code := (instr.operands().([]any))[1].([]Instruction)

	for i := 0; i < int(repititions); i++ {
		err := Interpet(code, state)

		if err != nil {
			return err
		}
	}

	return nil
}

func InterpretReserve(instr Instruction, state *InterpreterState) error {
	if state.configuration[ConfigureInstrAutoExpandQueuesKey] == EnabledConfigurationState {
		return ErrAutoExpandQueueNotCompatibleWithReserve
	}

	count := instr.operands().(uint16)

	if uint16(state.currentContext.currentQueue.End)+count > MaxAllowedObjectSpace {
		return ErrQueueOveflowsObjectSpace
	}

	state.currentContext.currentQueue.End += ObjectAddress(count)

	return nil
}

func InterpretPush(instr Instruction, state *InterpreterState) error {
	data := instr.operands().(uint8)

	if err := state.QueuePush(data); err != nil {
		return err
	}

	return nil
}

func InterpretDrain(instr Instruction, state *InterpreterState) error {
	count := instr.operands().(int32)

	if count > 0 {
		if state.currentContext.currentQueue.Start+ObjectAddress(count) > state.currentContext.currentQueue.End {
			return ErrCannotWriteOutsideOfCurrentQueue
		}

		for i := uint16(0); i < uint16(count); i++ {
			state.data[uint16(state.currentContext.currentQueue.Start)+i] = 0
		}

		state.qi -= ObjectAddress(count)
	} else {
		for i := uint16(0); i < (uint16(state.currentContext.currentQueue.End) - uint16(state.currentContext.currentQueue.Start)); i++ {
			state.data[uint16(state.currentContext.currentQueue.Start)+i] = 0
		}

		state.qi = 0
	}

	return nil
}

func InterpretDump(instr Instruction, state *InterpreterState) error {
	format := (instr.operands().([]any))[0].(DumpInstructionFormat)
	count := (instr.operands().([]any))[1].(int32)
	concat := (instr.operands().([]any))[2].(bool)

	switch format {
	case DumpInstrRawFormat:
		switch {
		case count > 0 && concat:
			fmt.Println(joinNumericSlice(state.data[uint16(state.currentContext.currentQueue.Start) : uint16(state.currentContext.currentQueue.Start)+uint16(count)]))
		case count == -1 && concat:
			fmt.Println(joinNumericSlice(state.data[state.currentContext.currentQueue.Start:state.currentContext.currentQueue.End]))
		case count > 0:
			for _, b := range state.data[uint16(state.currentContext.currentQueue.Start) : uint16(state.currentContext.currentQueue.Start)+uint16(count)] {
				fmt.Println(b)
			}
		case count == -1:
			for _, b := range state.data[state.currentContext.currentQueue.Start:state.currentContext.currentQueue.End] {
				fmt.Println(b)
			}
		}
	case DumpInstrTextFormat:
		switch {
		case count > 0 && concat:
			fmt.Println(string(state.data[uint16(state.currentContext.currentQueue.Start) : uint16(state.currentContext.currentQueue.Start)+uint16(count)]))
		case count == -1 && concat:
			fmt.Println(string(state.data[state.currentContext.currentQueue.Start:state.currentContext.currentQueue.End]))
		case count > 0:
			for _, b := range state.data[uint16(state.currentContext.currentQueue.Start) : uint16(state.currentContext.currentQueue.Start)+uint16(count)] {
				fmt.Printf("%c\n", b)
			}
		case count == -1:
			for _, b := range state.data[state.currentContext.currentQueue.Start:state.currentContext.currentQueue.End] {
				fmt.Printf("%c\n", b)
			}
		}
	}

	return nil
}

func (state *InterpreterState) QueuePop() uint8 {
	val := state.data[state.currentContext.currentQueue.Start]

	for i, v := range state.data[state.currentContext.currentQueue.Start:state.currentContext.currentQueue.End] {
		if i == 0 {
			continue
		}

		state.data[int(state.currentContext.currentQueue.Start)+i-1] = v
	}

	state.data[state.currentContext.currentQueue.End-1] = 0

	state.qi--

	if state.configuration[ConfigureInstrAutoShrinkQueuesKey] == EnabledConfigurationState {
		state.currentContext.currentQueue.End--
	}

	return val
}

func (state *InterpreterState) QueuePush(v uint8) error {
	if state.configuration[ConfigureInstrAutoExpandQueuesKey] == EnabledConfigurationState && state.qi+1 > state.currentContext.currentQueue.End {
		count := uint16(1)

		if uint16(state.currentContext.currentQueue.End)+count > MaxAllowedObjectSpace {
			return ErrQueueOveflowsObjectSpace
		}

		state.currentContext.currentQueue.End += ObjectAddress(count)
	} else if state.qi+1 > state.currentContext.currentQueue.End {
		return ErrQueueOverflow
	}

	state.data[state.currentContext.currentQueue.Start+state.qi] = v

	state.qi++

	return nil
}

func InterpretArithmetic(instr Instruction, state *InterpreterState) error {
	x := state.QueuePop()
	y := state.QueuePop()

	switch instr.kind() {
	case Add:
		if err := state.QueuePush(x + y); err != nil {
			return err
		}
	case Mul:
		if err := state.QueuePush(x * y); err != nil {
			return err
		}
	case Sub:
		if err := state.QueuePush(x - y); err != nil {
			return err
		}
	case Div:
		if err := state.QueuePush(x / y); err != nil {
			return err
		}
	case Rem:
		if err := state.QueuePush(x % y); err != nil {
			return err
		}
	case Pow:
		if err := state.QueuePush(uint8(math.Pow(float64(x), float64(y)))); err != nil {
			return err
		}
	}

	return nil
}

func Interpet(instrs []Instruction, state *InterpreterState) error {
	for _, instr := range instrs {
		switch instr.kind() {
		case Queueinit:
			err := InterpretQueueinit(instr, state)

			if err != nil {
				return err
			}
		case Repeat:
			err := InterpretRepeat(instr, state)

			if err != nil {
				return err
			}
		case With:
			err := InterpretWith(instr, state)

			if err != nil {
				return err
			}
		case Reserve:
			err := InterpretReserve(instr, state)

			if err != nil {
				return err
			}
		case Push:
			err := InterpretPush(instr, state)

			if err != nil {
				return err
			}
		case Dump, DumpN:
			err := InterpretDump(instr, state)

			if err != nil {
				return err
			}
		case Drain, DrainN:
			err := InterpretDrain(instr, state)

			if err != nil {
				return err
			}
		case Configure:
			err := InterpretConfigure(instr, state)

			if err != nil {
				return err
			}
		case Add, Sub, Mul, Div, Rem, Pow:
			err := InterpretArithmetic(instr, state)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	file, err := os.Open("/Volumes/monster/gibberish/hello.sqopl")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len([]byte(line)) > 0 && ([]byte(line))[0] != '#' {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	body, err := ParseBody()

	if err != nil {
		log.Fatal(err)
	}

	state := new(InterpreterState)

	state.data = [MaxAllowedObjectSpace]uint8{0}
	state.queueBounds = []Range[ObjectAddress]{}
	state.qi = ObjectAddress(0)

	state.currentContext = InterpreterStateContext{
		currentQueue: nil,
		subroutines:  map[string]InterpreterSubroutine{},
	}

	state.configuration = map[ConfigureInstrKey]ConfigurationState{
		ConfigureInstrAutoExpandQueuesKey: DisabledConfigurationState,
	}

	ierr := Interpet(body, state)

	if ierr != nil {
		log.Fatal(ierr)
	}
}
