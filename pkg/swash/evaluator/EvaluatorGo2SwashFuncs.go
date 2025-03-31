package evaluator

import (
	"cnc/pkg/swash"
	"errors"
	"fmt"
	"reflect"
)

// GoFunc represents the function internally for memory allocation processing
type GoFunc struct {
	Descriptor      string
	RequiredArgs    []*swash.Token
	ValueOfFunction reflect.Value
	GoFuncReturns   *swash.Token
}

// registerGoFunc will register the go function into memory
func (memory *Memory) registerGoFunc(funcName string, funcValue any) error {
	functionMemory := &GoFunc{
		Descriptor:      funcName,
		RequiredArgs:    make([]*swash.Token, 0),
		ValueOfFunction: reflect.ValueOf(funcValue),
		GoFuncReturns:   nil,
	}

	// appends the item into memory
	memory.memory = append(memory.memory, functionMemory)
	if err := functionMemory.Args(); err != nil || functionMemory.ValueOfFunction.Type().NumOut() == 0 {
		return err
	}

	// we only handle the first out handler
	shared := functionMemory.ValueOfFunction.Type().Out(0)
	switch shared.Kind() {

	case reflect.String:
		functionMemory.GoFuncReturns = &swash.Token{
			TokenType:    swash.STRING,
			TokenLiteral: shared.String(),
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		functionMemory.GoFuncReturns = &swash.Token{
			TokenType:    swash.NUMBER,
			TokenLiteral: shared.String(),
		}

	case reflect.Bool:
		functionMemory.GoFuncReturns = &swash.Token{
			TokenType:    swash.BOOLEAN,
			TokenLiteral: shared.String(),
		}

	case reflect.Map:
		functionMemory.GoFuncReturns = &swash.Token{
			TokenType:    swash.ANY,
			TokenLiteral: shared.String(),
		}
	}

	return nil
}

// funcArgs will register the go functions arguments
func (function *GoFunc) Args() error {
	for pos := 0; pos < function.ValueOfFunction.Type().NumIn(); pos++ {
		index := function.ValueOfFunction.Type().In(pos)
		switch index.Kind() {

		case reflect.Pointer:
			if index.Elem().Name() != "Token" {
				continue
			}

			function.RequiredArgs = append(function.RequiredArgs, &swash.Token{
				TokenType:    swash.TOKEN,
				TokenLiteral: index.Name(),
			})

		case reflect.String:
			function.RequiredArgs = append(function.RequiredArgs, &swash.Token{
				TokenType:    swash.STRING,
				TokenLiteral: index.Name(),
			})

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			function.RequiredArgs = append(function.RequiredArgs, &swash.Token{
				TokenType:    swash.NUMBER,
				TokenLiteral: index.Name(),
			})

		case reflect.Bool:
			function.RequiredArgs = append(function.RequiredArgs, &swash.Token{
				TokenType:    swash.BOOLEAN,
				TokenLiteral: index.Name(),
			})

		case reflect.Interface:
			function.RequiredArgs = append(function.RequiredArgs, &swash.Token{
				TokenType:    swash.ANY,
				TokenLiteral: index.Name(),
			})

		case reflect.Slice:
			switch index.String() {

			case "[]interface {}":
				function.RequiredArgs = append(function.RequiredArgs, &swash.Token{
					TokenType:    swash.VARIADIC_ANY,
					TokenLuggage: reflect.TypeOf(index.Name()),
					TokenLiteral: index.Name(),
				})

			case "[]int", "[]int8", "[]int16", "[]int32", "[]int64", "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64":
				function.RequiredArgs = append(function.RequiredArgs, &swash.Token{
					TokenType:    swash.VARIADIC_INT,
					TokenLuggage: reflect.TypeOf(index.Name()),
					TokenLiteral: index.Name(),
				})

			case "[]bool":
				function.RequiredArgs = append(function.RequiredArgs, &swash.Token{
					TokenType:    swash.VARIADIC_BOOL,
					TokenLuggage: reflect.TypeOf(index.Name()),
					TokenLiteral: index.Name(),
				})

			case "[]string":
				function.RequiredArgs = append(function.RequiredArgs, &swash.Token{
					TokenType:    swash.VARIADIC_STRING,
					TokenLuggage: reflect.TypeOf(index.Name()),
					TokenLiteral: index.Name(),
				})
			}

			return nil

		case reflect.Func:
			function.RequiredArgs = append(function.RequiredArgs, &swash.Token{
				TokenType:    swash.INDENT,
				TokenLuggage: index,
				TokenLiteral: index.Name(),
			})

		}
	}

	return nil
}

// execute will execute the function which is registered as Go
func (function *GoFunc) execute(token *swash.Token, memory *Memory) ([]reflect.Value, error) {
	args, err := function.classify(token.TokenArgs, token, memory)
	if err != nil {
		return make([]reflect.Value, 0), err
	}

	return function.ValueOfFunction.Call(args), nil
}

// classify will consume the functions arguments and then parse them into an understandable format.
func (function *GoFunc) classify(unclassified []*swash.Token, parent *swash.Token, memory *Memory) ([]reflect.Value, error) {
	classified := make([][]*swash.Token, 0)
	for _, unclassifiedArg := range unclassified {
		if len(classified) == 0 {
			classified = append(classified, make([]*swash.Token, 0))
		}

		switch unclassifiedArg.TokenType {

		default: /* non-seperator */
			classified[len(classified)-1] = append(classified[len(classified)-1], unclassifiedArg)

		case swash.COMMA: /* operation */
			classified = append(classified, make([]*swash.Token, 0))
		}
	}

	if len(function.RequiredArgs) == 0 {
		return make([]reflect.Value, 0), nil
	}

	/* checks the length of the arguments providied */
	if len(function.RequiredArgs) != len(classified) && !function.RequiredArgs[len(function.RequiredArgs)-1].TokenType.IsVariadic() {
		return nil, fmt.Errorf(errorFormatFunctionMissingArgs, function.Descriptor, len(function.RequiredArgs), len(classified))
	}

	functionArgs := make([]reflect.Value, 0)

	/* classified will range throughout the argument arrays */
	for classifiedPos, classifiedArgs := range classified {
		if classifiedPos < len(function.RequiredArgs) && classifiedArgs[0].TokenType == swash.INDENT || classifiedPos < len(function.RequiredArgs) && classifiedArgs[0].TokenType == swash.VARIABLE && len(classifiedArgs) == 1 || classifiedPos < len(function.RequiredArgs) && function.RequiredArgs[classifiedPos].TokenType == swash.TOKEN {
			constant, err := function.classifyHeader(classifiedArgs, memory, classifiedPos)
			if err != nil {
				return nil, err
			}

			functionArgs = append(functionArgs, constant)
			continue
		}

		TrueForm, err := memory.args(classifiedArgs)
		if err != nil || TrueForm == nil {
			return nil, err
		}

		/* checks the type ratio */
		if classifiedPos >= len(function.RequiredArgs) && !function.RequiredArgs[len(function.RequiredArgs)-1].TokenType.IsVariadic() || classifiedPos >= len(function.RequiredArgs) && function.RequiredArgs[len(function.RequiredArgs)-1].TokenType.IsVariadic() && !TrueForm.TokenType.Match(function.RequiredArgs[len(function.RequiredArgs)-1].TokenType) || classifiedPos < len(function.RequiredArgs) && !TrueForm.TokenType.Match(function.RequiredArgs[classifiedPos].TokenType) {
			tokenType := function.RequiredArgs[len(function.RequiredArgs)-1]
			if classifiedPos < len(function.RequiredArgs) {
				tokenType = function.RequiredArgs[classifiedPos]
			}

			return nil, fmt.Errorf(errorFormatFunctionMismatchType, function.Descriptor, tokenType.TokenType.String(), TrueForm.TokenType.String())
		}

		functionArgs = append(functionArgs, reflect.ValueOf(TrueForm.TokenType.Go(TrueForm.TokenLiteral)))
	}

	return functionArgs, nil
}

// classifyHeader will convert the special case header into it's true form for the interpreter.
func (function *GoFunc) classifyHeader(classified []*swash.Token, memory *Memory, pos int) (reflect.Value, error) {
	if function.RequiredArgs[pos].TokenType == swash.TOKEN {
		return reflect.ValueOf(classified[0]), nil
	}

	index, err := memory.searchIndex(classified[0].TokenLiteral)
	if err != nil || index == nil {
		return reflect.Value{}, err
	}

	switch object := index.(type) {

	case *Variable:
		return reflect.ValueOf(object.Value.TokenType.Go(object.Value.TokenLiteral)), nil

	case *Object:
		if object.TrueType == nil {
			return reflect.ValueOf(make(map[string]any)), nil
		}

		return reflect.ValueOf(object.TrueType), nil

	case *Function: /* port function from Swash to Go */
		if function.RequiredArgs[pos].TokenLuggage.NumIn() != len(object.Args) {
			return reflect.Value{}, fmt.Errorf(errorFormatFunctionMissingArgs, function.Descriptor, function.RequiredArgs[pos].TokenLuggage.NumIn(), len(object.Args))
		}

		for i := 0; i < function.RequiredArgs[pos].TokenLuggage.NumIn() && i < len(object.Args); i++ {
			switch function.RequiredArgs[pos].TokenLuggage.In(i).Kind() {

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if object.Args[i].Type == swash.NUMBER {
					continue
				}

				return reflect.Value{}, fmt.Errorf(errorFormatFunctionMismatchType, function.Descriptor, function.RequiredArgs[pos].TokenLuggage.In(i).Kind().String(), object.Args[i].Type.String())

			case reflect.String:
				if object.Args[i].Type == swash.STRING {
					continue
				}

				return reflect.Value{}, fmt.Errorf(errorFormatFunctionMismatchType, function.Descriptor, function.RequiredArgs[pos].TokenLuggage.In(i).Kind().String(), object.Args[i].Type.String())

			case reflect.Bool:
				if object.Args[i].Type == swash.BOOLEAN {
					continue
				}

				return reflect.Value{}, fmt.Errorf(errorFormatFunctionMismatchType, function.Descriptor, function.RequiredArgs[pos].TokenLuggage.In(i).Kind().String(), object.Args[i].Type.String())
			}
		}

		// reflect.MakeFunc will create the call point for the function.
		return reflect.MakeFunc(function.RequiredArgs[pos].TokenLuggage, function.interop(object, classified, memory, pos)), nil
	}

	return reflect.Value{}, errors.New("unknown factor implemented")
}

// interop is the segment between Go & Swash
func (function *GoFunc) interop(index *Function, classified []*swash.Token, memory *Memory, pos int) func([]reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		token, swashArgs := new(swash.Token), swash.ReflectToTokens(args, make([]*swash.Token, 0))
		token.TokenLiteral = index.Descriptor.TokenLiteral + "("
		for i, arg := range swashArgs {
			token.TokenLiteral += arg.TokenLiteral
			token.TokenArgs = append(token.TokenArgs, arg)
			if i+1 < len(swashArgs) {
				token.TokenArgs = append(token.TokenArgs, &swash.Token{
					TokenType:    swash.COMMA,
					TokenLiteral: ",",
				})
			}
		}

		token.TokenLiteral += ")"
		token.TokenType = swash.FUNCTION
		constant, err := memory.execute(index, token)
		if err != nil || constant == nil {
			return make([]reflect.Value, 0)
		}

		args = make([]reflect.Value, 0)
		return append(args, reflect.ValueOf(constant.TokenType.Go(constant.TokenLiteral)))
	}
}
