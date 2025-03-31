package swash

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
)

// TokenType allows us to distinguish between different types of lexer tokens
type TokenType int

// typeChecker is a low level type for type validation
type typeChecker func(r rune, i int) bool

// UNDEFINED is a token type that functions without value
var UNDEFINED = &Token{
	TokenType:    STRING,
	TokenArgs:    make([]*Token, 0),
	TokenLiteral: "undefined",
}

const (
	STRING TokenType = iota // "hi"
	NUMBER
	BOOLEAN
	INDENT
	VARIABLE
	FUNCTION
	TAGINDENT
	INDENTWORK

	/* Types below are here for the evaluator specifically */
	VARIADIC_STRING
	VARIADIC_INT
	VARIADIC_BOOL
	VARIADIC_ANY
	ANY

	// Our recognized patterns (added as i worked on this)
	ADD
	DIVIDE
	SUBTRACT
	MULTIPLY
	EQUAL
	DOLLAR
	SPACE
	DOUBLE_EQUAL
	SEMICOLON
	PARANTHESES_OPEN
	PARANTHESES_CLOSE
	BRACE_OPEN
	BRACE_CLOSE
	FULLSTOP
	COMMA
	COLON
	SUBTRACT_ARROW
	GREATERTHAN
	GREATEREQUAL
	LESSTHAN
	LESSEQUAL
	NOTEQUAL
	MODULUS
	SWASH_BODY_OPEN
	SWASH_BODY_CLOSE
	TEXT
	TEXT_BLANKLINE
	COMMENT
	TOKEN
	PLUSPLUS
	MINUSMINUS
	UNDERSCORE

	WORKSPACE_OPEN     string = "<?swash"
	WORKSPACETAG_OPEN  string = "<<"
	WORKSPACETAG_CLOSE string = ">>"
)

// Token will be collected in the mass and returned to the main thread.
type Token struct {
	TokenLine         int
	TokenType         TokenType
	TokenArgs         []*Token /* mainly for functions */
	TokenState        bool
	TokenLiteral      string
	TokenLiteralValue any
	TokenLuggage      reflect.Type
}

// newToken will attempt to create a new token from the information provided
func (tokenizer *Tokenizer) newToken(tokenType TokenType, tokenLiteral string, args ...*Token) *Token {
	return &Token{
		TokenLine:         tokenizer.tokenizerLine,
		TokenType:         tokenType,
		TokenArgs:         args,
		TokenState:        true,
		TokenLiteral:      tokenLiteral,
		TokenLiteralValue: tokenType.Go(tokenLiteral),
	}
}

func ReflectValueToToken(v reflect.Value, t TokenType) *Token {
	return &Token{
		TokenType:         t,
		TokenState:        true,
		TokenLiteral:      fmt.Sprint(v.Interface()),
		TokenLiteralValue: t.Go(fmt.Sprint(v.Interface())),
	}
}

// ChangeValue changes the value of the token
func (token *Token) ChangeValue(tokenLiteral string, tokenType TokenType) {
	token.TokenType = tokenType
	token.TokenLiteral = tokenLiteral
	token.TokenLiteralValue = token.TokenType.Go(tokenLiteral)
}

// NoState means it's ignored by the PeekNext func
func (token *Token) NoState() *Token {
	token.TokenState = false
	return token
}

// SameLine will return whether the token is on the same line as the previous token
func (p *Token) SameLine(s *Token) bool {
	return p.TokenLine == s.TokenLine
}

// Go will actively convert the string into it's original form
func (tokenType TokenType) Go(ltr string) any {
	switch tokenType {

	case STRING, INDENT:
		return ltr

	case NUMBER:
		number, err := strconv.Atoi(ltr)
		if err != nil {
			return 0
		}

		return number

	case BOOLEAN:
		binary, err := strconv.ParseBool(ltr)
		if err != nil {
			return false
		}

		return binary

	case ANY:
		return ltr

	default:
		return nil
	}
}

func (tokenType TokenType) String() string {
	switch tokenType {

	case STRING:
		return "STRING"

	case NUMBER:
		return "INTEGER"

	case BOOLEAN:
		return "BOOLEAN"

	case INDENT:
		return "INDENT"

	case VARIABLE:
		return "VAR"

	case FUNCTION:
		return "FUNC"

	case ANY:
		return "ANY"

	case VARIADIC_STRING:
		return "...STRING"

	case VARIADIC_INT:
		return "...INT"

	case VARIADIC_BOOL:
		return "...BOOL"

	case VARIADIC_ANY:
		return "...ANY"

	case TEXT:
		return "TEXT"

	case TEXT_BLANKLINE:
		return "BLANKLINE"

	case SWASH_BODY_OPEN:
		return "BODY_OPEN"

	case SWASH_BODY_CLOSE:
		return "BODY_CLOSE"

	case SUBTRACT:
		return "SUBTRACT"

	case ADD:
		return "ADDITION"

	case MULTIPLY:
		return "MULTIPLY"

	case DIVIDE:
		return "DIVIDE"

	default:
		return "UNKNOWN"
	}
}

// singleCharTyper will absorb the x rune and compare to the r rune
func singleCharTyper(x rune) typeChecker {
	return func(r rune, i int) bool {
		return x != r
	}
}

func indentCharType(r rune, i int) bool {
	return unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_'
}

// mergeFromPackage allows for easier implementation of different rule sets.
func mergeFromPackage(f func(rune) bool) typeChecker {
	return func(r rune, i int) bool {
		return f(r)
	}
}

// typeCheckers allows you to easily consume multiple different checkers
func typeCheckers(checkers []typeChecker, x rune, p int) bool {
	res := make([]bool, 0)

	for _, checker := range checkers {
		res = append(res, checker(x, p))
	}

	return slices.Contains(res, true)
}

// isVariadic checks whether the tokenType is of variadic type
func (t TokenType) IsVariadic() bool {
	return t == VARIADIC_STRING || t == VARIADIC_INT || t == VARIADIC_BOOL || t == VARIADIC_ANY
}

// MatchVariadic checks whether the tokenType is of variadic type
func (t TokenType) MatchVariadic(parent TokenType) bool {
	return parent == VARIADIC_STRING && t == STRING || parent == VARIADIC_INT && t == NUMBER || parent == VARIADIC_BOOL && t == BOOLEAN || parent == VARIADIC_ANY && t >= 0 && t <= 12
}

// Match will match varadic and non-varadic types
func (t TokenType) Match(parent TokenType) bool {
	if t.MatchVariadic(parent) {
		return true
	}

	return parent == ANY || t == parent
}

var Escapes = map[string]string{
	"\\x1b": "\x1b", "\\u001b": "\u001b", "\\033": "\033",
	"\\r": "\r", "\\n": "\n", "\\a": "\a",
	"\\b": "\b", "\\t": "\t", "\\v": "\v",
	"\\f": "\f", "\\007": "\007",
}

// strip will recombine the internal and then strip all ANSI codes from it.
func (tokenizer *Tokenizer) Strip() *Tokenizer {
	target := strings.Join(tokenizer.tokenizerTarget, "\n")

	for escape, renewed := range Escapes {
		target = strings.ReplaceAll(target, escape, renewed)
	}

	tokenizer.tokenizerTarget = strings.Split(target, "\n")
	return tokenizer
}

// PrefixEscape executes whenever the prefix is an ANSI code
func PrefixEscape(prefix string) bool {
	for _, c := range Escapes {
		if strings.HasPrefix(prefix, c) {
			return true
		}
	}

	return false
}

// ReflectToTokens will transform the array of values into tokens
func ReflectToTokens(values []reflect.Value, dest []*Token) []*Token {
	for _, value := range values {
		switch value.Kind() {

		case reflect.String:
			dest = append(dest, ReflectValueToToken(value, STRING))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			dest = append(dest, ReflectValueToToken(value, NUMBER))

		case reflect.Bool:
			dest = append(dest, ReflectValueToToken(value, BOOLEAN))
		}
	}

	return dest
}
