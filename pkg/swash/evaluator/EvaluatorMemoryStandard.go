package evaluator

import (
	"cnc/pkg/swash"
	"cnc/pkg/swash/packages"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

/*
	EvaluatorMemoryStandard.go maintains all the standard memory operations for the evaluator, this file maintains everything which is required by the evaluator for
	basic operational purposes. It converts a bunch of functions into Swash GoFunc's which can be called. we save a huge amount of time with this file and means we
	even import a few standard library's too. Happy Swash programming!
*/

// standard is what handles the individual inspection
type standard struct {
	evaluator *Evaluator
}

// register will register all the std methods
func (std *standard) register() {
	value := reflect.ValueOf(std)
	for i := 0; i < value.NumMethod(); i++ {
		context := reflect.TypeOf(std).Method(i)
		std.evaluator.Memory.Go2Swash(strings.ToLower(context.Name), value.Method(i).Interface())
	}

	/* We use this register to init some of our base pkg in Swash. */
	std.evaluator.Memory.WritePackage(packages.JSON, packages.JSONFunctions)
	std.evaluator.Memory.WritePackage(packages.TIME, packages.TIMEFunctions)
	std.evaluator.Memory.WritePackage(packages.HTTP, packages.HTTPFunctions)
}

// Print implements the standard interface from formatting
func (std standard) Print(values ...any) {
	fmt.Fprint(std.evaluator.Memory.wr, append(values, "\r\n")...)
}

// Printn Print implements the standard interface from formatting
func (std standard) Printn(values ...any) {
	fmt.Fprint(std.evaluator.Memory.wr, strings.ReplaceAll(fmt.Sprint(append(values)...), "\n", "\r\n"))
}

func (std standard) Repeat(v string, amount int) string {
	return strings.Repeat(v, amount)
}

// Printf implements the standard formatting package
func (std standard) Printf(format string, values ...any) {
	fmt.Fprintf(std.evaluator.Memory.wr, format+"\r\n", values...)
}

// Sprint implements the sprint function from standard formatting
func (std standard) Sprint(values ...any) string {
	return fmt.Sprint(values...)
}

// Sprintf implements the sprint format function from standard formatting
func (std standard) Sprintf(format string, values ...any) string {
	return fmt.Sprintf(format, values...)
}

// Len implements the standard interface for measuring the length of the value as a string
func (std standard) Len(value any) int {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		return len(value.([]any))
	}

	return len(fmt.Sprint(value))
}

// Itoa will convert a number to a string
func (std standard) Itoa(value int) string {
	return strconv.Itoa(value)
}

// Atoi will convert the string into the number
func (std standard) Atoi(value string) int {
	indent, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return indent
}

// TypeOf will convert into the string type of the object
func (std standard) TypeOf(value any) string {
	return reflect.TypeOf(value).String()
}

// PadRight will pad all the text towards the right
func (std standard) PadRight(text string, pad int) string {
	if pad-utf8.RuneCountInString(text) <= 0 {
		return text
	}

	return text + strings.Repeat(" ", pad-utf8.RuneCountInString(text))
}

// Replace implements the strings.ReplaceAll method
func (std standard) Replace(src, repl, to string) string {
	return strings.ReplaceAll(src, repl, to)
}

// Trim will trim the amount of trim of the content
func (std standard) Trim(src string, trim int) string {
	return src[:len(src)-trim]
}

// Require is the register for the require root note function
func (std standard) Require(path string) {}

// require will implement the function set for requiring routes
func (memory *Memory) require(constructor, arg *swash.Token, index *swash.Var) error {
	if len(filepath.Ext(arg.TokenLiteral)) == 0 {
		context, ok := memory.packages[arg.TokenLiteral]
		if !ok {
			return errors.New("unknown package imported")
		}

		return memory.Go2Swash(constructor.TokenLiteral, context)
	}

	tokenizer, err := swash.NewTokenizerSourcedFromFile(arg.TokenLiteral)
	if err != nil {
		return err
	}

	/* attempts to parse the tokenization route */
	if err := tokenizer.Parse(); err != nil {
		return err
	}

	eval := NewEvaluator(tokenizer, memory.wr, memory.rd)
	eval.Memory.memory = append(eval.Memory.memory, memory.memory...)
	if err := eval.Execute(); err != nil {
		return err
	}

	synchronize := &Object{
		Descriptor: constructor,
		Exporter:   index.Exporter,
		Values:     NewMemory(memory.wr, memory.rd, memory.packages),
	}

	/* once executed, loops through memory */
	for _, pointer := range eval.Memory.memory {
		switch object := pointer.(type) {

		case *Variable:
			if object.Var.Exporter == nil || object.Var.Exporter.TokenType != swash.TAGINDENT {
				continue
			}

			/* allocates into memory the variable */
			if err := synchronize.Values.allocateVar(object.Var); err != nil {
				return err
			}

		case *Function:
			if object.Exporter == nil || object.Exporter.TokenType != swash.TAGINDENT {
				continue
			}

			/* allocates into memory the function */
			if err := synchronize.Values.allocateFunc(object.Function); err != nil {
				return err
			}

		case *Object:
			if object.Exporter == nil || object.Exporter.TokenType != swash.TAGINDENT {
				continue
			}

			synchronize.Values.memory = append(synchronize.Values.memory, object)
		}
	}

	memory.memory = append(memory.memory, synchronize)
	return nil
}

// converts the str to uppercase
func (std standard) Uppercase(str string) string {
	return strings.ToUpper(str)
}

// converts the str to lowercase
func (std standard) Lowercase(str string) string {
	return strings.ToLower(str)
}
