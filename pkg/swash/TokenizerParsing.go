package swash

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

/*
	TokenizerParsing.go is a method which controls the parsing of tokens into it's own segmanted abstract syntax tree. this
	is a format which the interpreter can directly understand and execute the code as we handle it, this brings many advantages
	to the interpreter as it means all the parsing & lexing is done before it ever reaches the execution route and any errors
	can be caught and safely discarded off before the end user reaches them.
*/

type ParserNode any

type Var struct {
	Keyword    *Token
	Exporter   *Token
	Descriptor *Token
	Args       []*Token
}

type Function struct {
	Descriptor Token
	Args       []FunctionArg
	Returns    *TokenType
	Exporter   *Token
	Nodes      *Tokenizer
}

// Object is a reference towards GoToSwash
type FunctionArg struct {
	Descriptor *Token
	Object     reflect.Type
	Type       TokenType
}

type Return struct {
	Descriptor *Token
	Args       []*Token
}

type If struct {
	Descriptor *Token
	Decision   []*Token
	Body       *Tokenizer
}

// Parse is the function which will force the runtime around the internal code being consumed
func (tokenizer *Tokenizer) Parse() error {
	token, ok := tokenizer.PeekNext()
	if !ok || token == nil {
		return nil
	}

	return tokenizer.handleBuf(token)
}

// handleBuf handles the current state of the tokenizer position and in return this will tokenize -> parse form
func (tokenizer *Tokenizer) handleBuf(token *Token) error {
	switch token.TokenLiteral {

	case "var", "const":
		return tokenizer.handleVar(token)

	case "func":
		return tokenizer.handleFunc(token)

	case "return":
		return tokenizer.handleReturn(token)

	case "if":
		return tokenizer.handleIf(token)
	}

	// references to type-based tokenization, whereas above is keyword-based tokenization
	switch token.TokenType {

	case FUNCTION, VARIABLE, TEXT, TEXT_BLANKLINE:
		operation, err := tokenizer.handleOperation(token)
		if err != nil {
			return err
		}

		tokenizer.tokenizerParserNodes = append(tokenizer.tokenizerParserNodes, operation)
		peek, ok := tokenizer.PeekNext()
		if !ok || peek == nil {
			return nil
		}

		return tokenizer.handleBuf(peek)

	case SWASH_BODY_OPEN, SWASH_BODY_CLOSE, TAGINDENT, SEMICOLON:
		peek, ok := tokenizer.PeekNext()
		if !ok || peek == nil {
			return nil
		}

		return tokenizer.handleBuf(peek)
	}

	return nil
}

// handleOperation will implement the fixing and bindings for operating within statements
func (tokenizer *Tokenizer) handleOperation(token *Token) (any, error) {
	if token.TokenType == TEXT_BLANKLINE || token.TokenType == TEXT || token.TokenType == FUNCTION {
		return token, nil
	}

	operator, ok := tokenizer.PeekNext()
	if !ok || operator == nil {
		return token, nil
	}

	indent := &Var{
		Args:       make([]*Token, 0),
		Keyword:    &Token{TokenLiteral: "var", TokenType: INDENT},
		Exporter:   new(Token),
		Descriptor: token,
	}

	switch operator.TokenType {

	case PLUSPLUS:
		indent.Args = append(indent.Args, token)
		indent.Args = append(indent.Args, &Token{TokenType: ADD, TokenLiteral: "+"})
		indent.Args = append(indent.Args, &Token{TokenType: NUMBER, TokenLiteral: "1"})
		return indent, nil

	case MINUSMINUS:
		indent.Args = append(indent.Args, token)
		indent.Args = append(indent.Args, &Token{TokenType: SUBTRACT, TokenLiteral: "-"})
		indent.Args = append(indent.Args, &Token{TokenType: NUMBER, TokenLiteral: "1"})
		return indent, nil
	}

	return indent, nil
}

// handleVar summons whenever the variable declaration is encountered
func (tokenizer *Tokenizer) handleVar(token *Token) error {
	exporter := new(Token)
	if len(tokenizer.tokenizerPrevious)-2 > 0 && tokenizer.tokenizerPrevious[len(tokenizer.tokenizerPrevious)-2].TokenType == TAGINDENT {
		exporter = tokenizer.tokenizerPrevious[len(tokenizer.tokenizerPrevious)-2]
	}

	descriptor, ok := tokenizer.PeekNext()
	if !ok || descriptor == nil || descriptor.TokenType != VARIABLE && descriptor.TokenType != INDENT {
		return errors.New("error occurred while parsing variable declaration descriptor due to conflicts")
	}

	sequel, ok := tokenizer.PeekNext()
	if !ok || sequel == nil || sequel.TokenType != EQUAL && sequel.TokenType != SEMICOLON {
		return errors.New("error occurred while parsing variable declaration separator due to conflicts")
	}

	args := make([]*Token, 0)

	for {
		arg, ok := tokenizer.PeekNext()
		if !ok || arg == nil {
			tokenizer.tokenizerParserNodes = append(tokenizer.tokenizerParserNodes, &Var{
				Descriptor: descriptor,
				Exporter:   exporter,
				Keyword:    token,
				Args:       args,
			})

			return nil
		}

		args = append(args, arg)

		/* Whenever a semi colon is found which enforces the closement */
		if ok := token.SameLine(arg); !ok || arg.TokenType == SEMICOLON {
			tokenizer.tokenizerParserNodes = append(tokenizer.tokenizerParserNodes, &Var{
				Descriptor: descriptor,
				Exporter:   exporter,
				Keyword:    token,
				Args:       args[:len(args)-1],
			})

			return tokenizer.handleBuf(arg)
		}
	}
}

// handleFunc will directly handle the function embedded within the statement
func (tokenizer *Tokenizer) handleFunc(token *Token) error {
	exporter := new(Token)
	if len(tokenizer.tokenizerPrevious)-2 >= 1 && tokenizer.tokenizerPrevious[len(tokenizer.tokenizerPrevious)-2].TokenType == TAGINDENT {
		exporter = tokenizer.tokenizerPrevious[len(tokenizer.tokenizerPrevious)-2]
	}

	descriptor, ok := tokenizer.PeekNext()
	if !ok || descriptor == nil || descriptor.TokenType != FUNCTION {
		return errors.New("error occurred while parsing function declaration descriptor due to conflicts")
	}

	Func := Function{
		Descriptor: *descriptor,
		Exporter:   exporter,
		Args:       make([]FunctionArg, 0),
		Returns:    nil,
		Nodes:      new(Tokenizer),
	}

	/* strips anything else included within it*/
	Func.Descriptor.TokenLiteral = strings.Split(Func.Descriptor.TokenLiteral, "(")[0]
	if len(strings.Split(strings.Join(strings.Split(descriptor.TokenLiteral, "(")[1:], "("), ")")[:strings.Count(descriptor.TokenLiteral, ")")][0]) > 0 {
		for _, arg := range strings.Split(strings.Split(strings.Join(strings.Split(descriptor.TokenLiteral, "(")[1:], "("), ")")[:strings.Count(descriptor.TokenLiteral, ")")][0], ",") {
			index := strings.Split(strings.ReplaceAll(arg, " ", ""), ":")
			switch index[1] {

			case "string":
				Func.Args = append(Func.Args, FunctionArg{
					Descriptor: &Token{TokenType: VARIABLE, TokenLiteral: index[0]},
					Type:       STRING,
				})

			case "int", "number":
				Func.Args = append(Func.Args, FunctionArg{
					Descriptor: &Token{TokenType: VARIABLE, TokenLiteral: index[0]},
					Type:       NUMBER,
				})

			case "bool":
				Func.Args = append(Func.Args, FunctionArg{
					Descriptor: &Token{TokenType: VARIABLE, TokenLiteral: index[0]},
					Type:       BOOLEAN,
				})

			default:
				return fmt.Errorf("unsupported type declaration inside %s: %s", Func.Descriptor.TokenLiteral, index[1])
			}
		}
	}

	question, ok := tokenizer.PeekNext()
	if !ok || question == nil || question.TokenType != SUBTRACT_ARROW && question.TokenType != BRACE_OPEN {
		return errors.New("error occurred while parsing function declaration due to conflicts")
	}

	/* handles the return state */
	if question.TokenType == SUBTRACT_ARROW {
		returnable, ok := tokenizer.PeekNext()
		if !ok || returnable == nil || returnable.TokenType != INDENT {
			return errors.New("error occurred while parsing function declaration due to return value conflicts")
		}

		switch returnable.TokenLiteral {

		case "string":
			returnable.TokenType = STRING

		case "int", "number":
			returnable.TokenType = NUMBER

		case "bool":
			returnable.TokenType = BOOLEAN

		}

		Func.Returns = &returnable.TokenType
		question, ok = tokenizer.PeekNext()
		if !ok || question == nil || question.TokenType != BRACE_OPEN {
			return errors.New("error occurred while parsing function declaration due to conflicts")
		}
	}

	tokens := make([]*Token, 0)
	if err := tokenizer.handleBody(BRACE_OPEN, BRACE_CLOSE, &tokens); err != nil {
		return err
	}

	mimic := &Tokenizer{
		tokenizerText:        false,
		TokenizerStream:      tokens,
		tokenizerParserNodes: make([]ParserNode, 0),
	}

	if err := mimic.Parse(); err != nil {
		return err
	}

	Func.Nodes = mimic
	tokenizer.tokenizerParserNodes = append(tokenizer.tokenizerParserNodes, &Func)
	indentation, ok := tokenizer.PeekNext()
	if !ok || indentation == nil {
		return nil
	}

	return tokenizer.handleBuf(indentation)
}

// handleReturn will handle and parse the return values
func (tokenizer *Tokenizer) handleReturn(token *Token) error {
	args := make([]*Token, 0)

	for {
		arg, ok := tokenizer.PeekNext()
		if !ok || arg == nil {
			tokenizer.tokenizerParserNodes = append(tokenizer.tokenizerParserNodes, &Return{
				Descriptor: token,
				Args:       args,
			})

			return nil
		}

		args = append(args, arg)

		/* Whenever a semi colon is found which enforces the closement */
		if ok := token.SameLine(arg); !ok || arg.TokenType == SEMICOLON {
			tokenizer.tokenizerParserNodes = append(tokenizer.tokenizerParserNodes, &Return{
				Descriptor: token,
				Args:       args[:len(args)-1],
			})

			return tokenizer.handleBuf(arg)
		}
	}
}

// handleIf will parse the context and return the information depending on it
func (tokenizer *Tokenizer) handleIf(token *Token) error {
	statement := &If{
		Descriptor: token,
		Decision:   make([]*Token, 0),
		Body:       new(Tokenizer),
	}

	indent, ok := tokenizer.PeekNext()
	if !ok || indent == nil || indent.TokenType != INDENTWORK {
		return errors.New("missing statement decision")
	}

	statement.Decision = indent.TokenArgs
	next, ok := tokenizer.PeekNext()
	if !ok || next == nil || next.TokenType != BRACE_OPEN {
		return errors.New("missing statement body")
	}

	cache := make([]*Token, 0)
	err := tokenizer.handleBody(BRACE_OPEN, BRACE_CLOSE, &cache)
	if err != nil {
		return err
	}

	mimic := &Tokenizer{
		tokenizerText:        false,
		TokenizerStream:      cache,
		tokenizerParserNodes: make([]ParserNode, 0),
	}

	if err := mimic.Parse(); err != nil {
		return err
	}

	statement.Body = mimic
	tokenizer.tokenizerParserNodes = append(tokenizer.tokenizerParserNodes, statement)
	indentation, ok := tokenizer.PeekNext()
	if !ok || indentation == nil {
		fmt.Println(indentation)
		return nil
	}

	return tokenizer.handleBuf(indentation)
}

// handleBody does a polled search within the context and returns the information
func (tokenizer *Tokenizer) handleBody(open, close TokenType, tokens *[]*Token) error {
	depth := 0

	for {
		token, ok := tokenizer.PeekNext()
		if !ok || token == nil {
			return errors.New("body not closed")
		}

		switch token.TokenType {

		case close:
			if depth == 0 {
				return nil
			}

			depth--
			*tokens = append(*tokens, token)
			continue

		case open:
			depth++
			*tokens = append(*tokens, token)
			continue

		default:
			*tokens = append(*tokens, token)
		}
	}
}
