package evaluator

import (
	"cnc/pkg/swash"
	"errors"
	"fmt"
	"strconv"
)

/*
	EvaluatorArgs.go is probably the hardest & most complex thing about this language as it's
	something I'm yet to entirely understand, this allows the evaluator to build the programs
	args which can be found everywhere.

	the we separate operators and operations is via the position, even = operation &
	odd = operator, we then use this it generally iterate across it
*/

// args will compile the args into a single token, this will even execute functions contained within the runtime.
func (m *Memory) args(args []*swash.Token) (*swash.Token, error) {
	if len(args) == 0 {
		return swash.UNDEFINED, nil
	}

	constant, err := m.evalToken(args[0])
	if err != nil {
		return nil, err
	}

	for p := 2; p < len(args); p += 2 {
		arg, err := m.evalToken(args[p])
		if err != nil || arg == nil {
			return nil, errors.New("mismatched types")
		}

		// differentiates between different operators
		switch args[p-1].TokenType {

		case swash.DOUBLE_EQUAL, swash.GREATEREQUAL, swash.LESSEQUAL, swash.GREATERTHAN, swash.LESSTHAN, swash.NOTEQUAL:
			second, err := m.args(args[p:])
			if err != nil {
				return nil, err
			}

			if second.TokenType != constant.TokenType {
				return nil, errors.New("mismatched types")
			}

			switch args[p-1].TokenType {

			case swash.DOUBLE_EQUAL:
				constant.ChangeValue(fmt.Sprint(second.TokenLiteral == constant.TokenLiteral), swash.BOOLEAN)

			case swash.NOTEQUAL:
				constant.ChangeValue(fmt.Sprint(second.TokenLiteral != constant.TokenLiteral), swash.BOOLEAN)

			// these requires conversion to go types
			case swash.GREATEREQUAL, swash.LESSEQUAL, swash.GREATERTHAN, swash.LESSTHAN:
				integerPrimary, err := strconv.Atoi(constant.TokenLiteral)
				if err != nil {
					return nil, err
				}

				integerSecondly, err := strconv.Atoi(second.TokenLiteral)
				if err != nil {
					return nil, err
				}

				switch args[p-1].TokenType {

				case swash.GREATEREQUAL:
					constant.ChangeValue(fmt.Sprint(integerPrimary >= integerSecondly), swash.BOOLEAN)

				case swash.LESSEQUAL:
					constant.ChangeValue(fmt.Sprint(integerPrimary <= integerSecondly), swash.BOOLEAN)

				case swash.LESSTHAN:
					constant.ChangeValue(fmt.Sprint(integerPrimary < integerSecondly), swash.BOOLEAN)

				case swash.GREATERTHAN:
					constant.ChangeValue(fmt.Sprint(integerPrimary > integerSecondly), swash.BOOLEAN)
				}

			}

			return constant, nil

		case swash.MULTIPLY, swash.DIVIDE, swash.SUBTRACT, swash.MODULUS: // multiplication, division, subtraction, modulus
			if arg.TokenType != constant.TokenType || arg.TokenType != swash.NUMBER {
				return nil, errors.New("mismatched types")
			}

			index, err := strconv.Atoi(constant.TokenLiteral)
			if err != nil {
				return nil, err
			}

			indent, err := strconv.Atoi(arg.TokenLiteral)
			if err != nil {
				return nil, err
			}

			switch args[p-1].TokenType {

			case swash.MULTIPLY:
				constant.ChangeValue(fmt.Sprint(index*indent), swash.NUMBER)

			case swash.DIVIDE:
				constant.ChangeValue(fmt.Sprint(index/indent), swash.NUMBER)

			case swash.SUBTRACT:
				constant.ChangeValue(fmt.Sprint(index-indent), swash.NUMBER)

			case swash.MODULUS:
				constant.ChangeValue(fmt.Sprint(index%indent), swash.NUMBER)
			}

		case swash.ADD: // addition
			if arg.TokenType != constant.TokenType {
				return nil, errors.New("mismatched types")
			}

			switch arg.TokenType {

			case swash.STRING: // appends both strings together
				constant.TokenLiteral += arg.TokenLiteral

			case swash.NUMBER: // adds both numbers together
				index, err := strconv.Atoi(constant.TokenLiteral)
				if err != nil {
					return nil, err
				}

				indent, err := strconv.Atoi(arg.TokenLiteral)
				if err != nil {
					return nil, err
				}

				constant.ChangeValue(fmt.Sprint(index+indent), swash.NUMBER)
			}
		}
	}

	return constant, nil
}

// eval will compile each individual expression into a single token
func (m *Memory) evalToken(token *swash.Token) (*swash.Token, error) {
	switch token.TokenType {

	case swash.INDENTWORK:
		return m.args(token.TokenArgs)

	case swash.NUMBER, swash.STRING, swash.BOOLEAN:
		return token, nil

	case swash.VARIABLE, swash.FUNCTION:
		value, err := m.index(token)
		if err != nil || value == nil {
			return nil, err
		}

		return &swash.Token{TokenLine: value.TokenLine, TokenType: value.TokenType, TokenArgs: value.TokenArgs, TokenState: value.TokenState, TokenLiteral: value.TokenLiteral, TokenLiteralValue: value.TokenType.Go(value.TokenLiteral)}, nil

	default:
		return nil, errors.New("invalid type")
	}
}
