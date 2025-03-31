package evaluator

import (
	"cnc/pkg/swash"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
	EvaluatorMemory.go represents all the memory handling procedures within the program, the
	Pointer structure within this file is what we use to transfer information between functions
	and methods, it means that it can represent Functions, Variables & scopes entirely.

*/

type Memory struct {
	packages map[string]any
	memory   []any
	wr       io.Writer
	rd       io.Reader
}

type Object struct {
	Descriptor *swash.Token
	Values     *Memory
	TrueType   any
	Exporter   *swash.Token
}

type Variable struct {
	Descriptor *swash.Token
	Value      *swash.Token
	Var        *swash.Var
}

type Function struct {
	Descriptor *swash.Token
	Tokenizer  *swash.Tokenizer
	Exporter   *swash.Token
	Function   *swash.Function
	Return     *swash.TokenType
	Args       []swash.FunctionArg
}

// NewMemory inits a brand new memory object to keep track of information.
func NewMemory(wr io.Writer, rd io.Reader, pkgs map[string]any) *Memory {
	memory := &Memory{
		packages: pkgs,
		memory:   make([]any, 0),
		wr:       wr,
		rd:       rd,
	}

	if memory.rd == nil {
		memory.rd = os.Stdin
	}

	memory.WritePackage("", nil)
	return memory
}

// WritePackage will write the package into memory
func (m *Memory) WritePackage(pkg string, v any) error {
	if _, ok := m.packages[pkg]; ok {
		return errors.New("package already imported")
	}

	m.packages[pkg] = v
	return nil
}

// allocateFunc will allocate the infrastructure to register the swash function registry
func (memory *Memory) allocateFunc(index *swash.Function) error {
	if strings.Count(index.Descriptor.TokenLiteral, ".") > 0 {
		return errors.New("can't allocate function to object function")
	}

	/* indexes within memory to check if the function already exists */
	if object, err := memory.index(&index.Descriptor); err == nil && object != nil {
		return errors.New("function object already registered")
	}

	function := &Function{
		Descriptor: &index.Descriptor,
		Tokenizer:  index.Nodes,
		Function:   index,
		Return:     index.Returns,
		Args:       index.Args,
	}

	memory.memory = append(memory.memory, function)
	return nil
}

// allocateVar will look through the current scope & then allocate if it can.
func (memory *Memory) allocateVar(index *swash.Var) error {
	if len(index.Args) == 1 {
		switch index.Args[0].TokenType {

		case swash.FUNCTION: /* registers the function */
			constant, err := memory.indexGoFunc(index.Args[0])
			if err != nil {
				return err
			}

			switch constant.Descriptor {

			case "require":
				token, err := memory.args(index.Args[0].TokenArgs)
				if err != nil {
					return err
				}

				return memory.require(index.Descriptor, token, index)

			default:
				tags, err := constant.execute(index.Args[0], memory)
				if err != nil || len(tags) != 1 {
					return err
				}

				return memory.Go2Swash(index.Descriptor.TokenLiteral, tags[0].Interface())
			}
		}
	}

	value, err := memory.args(index.Args)
	if err != nil {
		return err
	}

	/* We check if the memory has already been allocated*/
	if object, err := memory.index(index.Descriptor); err == nil && object != nil {
		if index.Keyword.TokenLiteral == "const" {
			return nil
		}

		return memory.change(index, value)
	}

	/* appends the allocated memory */
	memory.memory = append(memory.memory, &Variable{
		Descriptor: index.Descriptor,
		Value:      value,
		Var:        index,
	})

	return nil
}

// index will only return the values of the token representing a func or variable
func (memory *Memory) index(token *swash.Token) (*swash.Token, error) {
	allocatedMemory := NewMemory(memory.wr, memory.rd, memory.packages)
	allocatedMemory.memory = append(allocatedMemory.memory, memory.memory...)

	// implements the scope based interface indexing
	for i, arg := range strings.Split(strings.Split(token.TokenLiteral, "(")[0], ".") {
		if i+1 >= len(strings.Split(strings.Split(token.TokenLiteral, "(")[0], ".")) {
			reference, err := allocatedMemory.search(arg)
			if err != nil {
				return nil, err
			}
			// switches depending on the reference
			switch object := reference.(type) {

			case *Variable:
				return object.Value, nil

			case *Function:
				return memory.execute(object, token)

			case *GoFunc:
				values, err := object.execute(token, memory)
				if err != nil || len(values) == 0 || len(values) >= 2 || object.GoFuncReturns == nil {
					return nil, err
				}

				return swash.ReflectValueToToken(values[0], object.GoFuncReturns.TokenType), nil

			default:
				break
			}
		}

		/* looks for object components and nothing else. */
		body, err := allocatedMemory.search(arg)
		if err != nil {
			return nil, err
		}

		switch object := body.(type) {

		case *Object:
			allocatedMemory = object.Values

		default:
			return nil, errors.New("non object type found in object reference")
		}
	}

	return nil, fmt.Errorf("undefined value: %s", token.TokenLiteral)
}

func (memory *Memory) searchIndex(indent string) (any, error) {
	allocatedMemory := NewMemory(memory.wr, memory.rd, memory.packages)
	allocatedMemory.memory = append(allocatedMemory.memory, memory.memory...)

	// implements the scope based interface indexing
	for i, arg := range strings.Split(strings.Split(indent, "(")[0], ".") {
		if i+1 >= len(strings.Split(strings.Split(indent, "(")[0], ".")) {
			reference, err := allocatedMemory.search(arg)
			if err != nil {
				return nil, err
			}

			return reference, nil
		}

		/* looks for object components and nothing else. */
		body, err := allocatedMemory.search(arg)
		if err != nil {
			return nil, err
		}

		switch object := body.(type) {

		case *Object:
			allocatedMemory = object.Values

		default:
			return nil, errors.New("non object type found in object reference")
		}
	}

	return nil, nil
}

// object will attempt to look within the memory scope for the object
func (memory *Memory) search(indent string) (any, error) {
	for _, pointer := range memory.memory {
		switch object := pointer.(type) {

		case *Object: /* object inspection */
			if object.Descriptor.TokenLiteral != indent {
				continue
			}

			return object, nil

		case *Function: /* function inspection */
			if object.Descriptor.TokenLiteral != indent {
				continue
			}

			return object, nil

		case *GoFunc: /* GoFunction inspection */
			if object.Descriptor != indent {
				continue
			}

			return object, nil

		case *Variable: /* variable inspection */
			if object.Descriptor.TokenLiteral != indent {
				continue
			}

			return object, nil
		}
	}

	return nil, fmt.Errorf("undefined reference: %s", indent)
}

// indexGoFunc will specifically index go functions
func (memory *Memory) indexGoFunc(token *swash.Token) (*GoFunc, error) {
	allocatedMemory := NewMemory(memory.wr, memory.rd, memory.packages)
	allocatedMemory.memory = append(allocatedMemory.memory, memory.memory...)

	// implements the scope based interface indexing
	for i, arg := range strings.Split(strings.Split(token.TokenLiteral, "(")[0], ".") {
		if i+1 >= len(strings.Split(strings.Split(token.TokenLiteral, "(")[0], ".")) {
			reference, err := allocatedMemory.search(arg)
			if err != nil {
				return nil, err
			}
			// switches depending on the reference
			switch object := reference.(type) {

			case *GoFunc:
				return object, nil
			}
		}

		/* looks for object components and nothing else. */
		body, err := allocatedMemory.search(arg)
		if err != nil {
			return nil, err
		}

		switch object := body.(type) {

		case *Object:
			allocatedMemory = object.Values

		default:
			return nil, errors.New("non object type found in object reference")
		}
	}

	return nil, fmt.Errorf("undefined value: %s", token.TokenLiteral)
}

// change will attempt to change the value of the statement in memory
func (memory *Memory) change(statement *swash.Var, value *swash.Token) error {
	for _, pointer := range memory.memory {
		switch object := pointer.(type) {

		case *Variable:
			if object.Descriptor.TokenLiteral != statement.Descriptor.TokenLiteral {
				continue
			}

			object.Value = value
			return nil
		}
	}

	return errors.New("undefined reference")
}

// execute will compile the current function into a forced executor context
func (memory *Memory) execute(f *Function, t *swash.Token) (*swash.Token, error) {
	evaluator := NewEvaluator(f.Tokenizer, memory.wr, memory.rd)
	evaluator.Memory.memory = append(memory.memory, evaluator.Memory.memory...)

	if len(f.Args) > 0 {
		arguments := make([][]*swash.Token, 0)
		arguments = append(arguments, make([]*swash.Token, 0))
		for _, arg := range t.TokenArgs {
			switch arg.TokenType {

			case swash.COMMA:
				arguments = append(arguments, make([]*swash.Token, 0))

			default:
				arguments[len(arguments)-1] = append(arguments[len(arguments)-1], arg)
			}
		}

		if len(t.TokenArgs) <= 0 || len(f.Args) != len(arguments) {
			return nil, fmt.Errorf(errorFormatFunctionMissingArgs, f.Descriptor.TokenLiteral, len(f.Args), len(arguments))
		}

		for p, data := range f.Args {
			switch data.Type {

			default:
				value, err := memory.args(arguments[p])
				if err != nil {
					return nil, err
				}

				if data.Type != value.TokenType {
					return nil, fmt.Errorf(errorFormatFunctionMismatchType, f.Descriptor.TokenLiteral, data.Type.String(), value.TokenType.String())
				}

				if err := evaluator.Memory.Go2Swash(strings.TrimPrefix(data.Descriptor.TokenLiteral, "$"), value.TokenLiteralValue); err != nil {
					return nil, err
				}

				continue
			}
		}
	}

	context, _, err := evaluator.self(evaluator.tokenizer.Nodes(), f.Return, f.Return != nil)
	if err != nil {
		return nil, err
	}

	/* replaces all the different variables within the internal */
	for _, memoryArg := range memory.memory {
		switch object := memoryArg.(type) {

		case *Variable:
			if object.Var.Keyword.TokenLiteral == "const" {
				return context, nil
			}

			val, err := evaluator.Memory.index(object.Descriptor)
			if err != nil {
				return nil, err
			}

			if err := memory.change(object.Var, val); err != nil {
				return nil, err
			}
		}
	}

	return context, nil
}
