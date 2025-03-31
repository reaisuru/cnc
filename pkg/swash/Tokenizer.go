package swash

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Tokenizer is the interface in which the tokenizer requires
type Tokenizer struct {
	tokenizerText        bool
	tokenizerLine        int
	tokenizerPrevious    []*Token
	TokenizerStream      []*Token
	tokenizerTarget      []string
	tokenizerParserNodes []ParserNode
}

// NewTokenizer creates a new Tokenizer with the object tokenizerTarget as the focus
func NewTokenizer(tokenizerTarget string, tokenizerText bool) *Tokenizer {
	return &Tokenizer{
		tokenizerText:        tokenizerText,
		tokenizerLine:        0,
		TokenizerStream:      make([]*Token, 0),
		tokenizerTarget:      strings.Split(tokenizerTarget, "\n"),
		tokenizerPrevious:    make([]*Token, 0),
		tokenizerParserNodes: make([]ParserNode, 0),
	}
}

// NewTokenizerSourcedFromFile will attempt to create the tokenizer imported from a file
func NewTokenizerSourcedFromFile(file string) (*Tokenizer, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	tokenizer := NewTokenizer(string(content), true)
	return tokenizer.Strip(), nil
}

// PeekNext attempts to peek one token into the future, once peeked it removes it from the tokenizerTarget
func (tokenizer *Tokenizer) PeekNext() (*Token, bool) {
	if tokenizer.tokenizerTarget == nil {
		if len(tokenizer.TokenizerStream) == 0 {
			return nil, false
		}

		current := tokenizer.TokenizerStream[0]
		tokenizer.TokenizerStream = tokenizer.TokenizerStream[1:]
		return current, true
	}

	if len(tokenizer.tokenizerTarget[tokenizer.tokenizerLine]) <= 0 {
		if len(tokenizer.tokenizerTarget) <= tokenizer.tokenizerLine+1 {
			return nil, false
		}

		tokenizer.tokenizerLine++
		return tokenizer.PeekNext()
	}

	token := tokenizer.tokenize(tokenizer.tokenizerTarget[tokenizer.tokenizerLine][0])
	if token == nil {
		return nil, true
	}

	tokenizer.shift(len(token.TokenLiteral))
	if !token.TokenState {
		return tokenizer.PeekNext()
	}

	/* skips to the next line */
	if token.TokenType == COMMENT {
		tokenizer.tokenizerTarget[tokenizer.tokenizerLine] = ""
		return tokenizer.PeekNext()
	}

	tokenizer.tokenizerPrevious = append(tokenizer.tokenizerPrevious, token)
	return token, true
}

// shift will move the current scope across one
func (tokenizer *Tokenizer) shift(i int) bool {
	if len(tokenizer.tokenizerTarget[tokenizer.tokenizerLine]) < i {
		tokenizer.tokenizerTarget[tokenizer.tokenizerLine] = ""
		return true
	}

	tokenizer.tokenizerTarget[tokenizer.tokenizerLine] = tokenizer.tokenizerTarget[tokenizer.tokenizerLine][i:]
	return true
}

// nextInline will return the next inline character
func (tokenizer *Tokenizer) nextInline(i int) rune {
	if len(tokenizer.tokenizerTarget[tokenizer.tokenizerLine]) <= i {
		return 0
	}

	return rune(tokenizer.tokenizerTarget[tokenizer.tokenizerLine][i])
}

// tokenize will take the current x value and attempt to form it into a single token
func (tokenizer *Tokenizer) tokenize(x byte) *Token {
	if tokenizer.tokenizerText {
		return tokenizer.newWorkspace()
	}

	switch x {

	case '+': // Math/arg operator
		if tokenizer.nextInline(1) == '+' {
			return tokenizer.newToken(PLUSPLUS, "++")
		}

		return tokenizer.newToken(ADD, string(x))

	case '-': // Math operator
		if tokenizer.nextInline(1) == '>' {
			return tokenizer.newToken(SUBTRACT_ARROW, "->")
		} else if tokenizer.nextInline(1) == '-' {
			return tokenizer.newToken(MINUSMINUS, "--")
		}

		return tokenizer.newToken(SUBTRACT, string(x))

	case '/': // Math operator
		if tokenizer.nextInline(1) == '/' {
			return tokenizer.newToken(COMMENT, "//")
		}

		return tokenizer.newToken(DIVIDE, string(x))

	case '*': // Math operator
		return tokenizer.newToken(MULTIPLY, string(x))

	case '%':
		return tokenizer.newToken(MODULUS, string(x))

	case ';': // Semicolon
		return tokenizer.newToken(SEMICOLON, string(x))

	case '(': /* implementing support for arg references */
		token := tokenizer.tokenizeArgs(tokenizer.runUntilVoid(0, singleCharTyper('(')))
		token.TokenType = INDENTWORK
		return token

	case ')':
		return tokenizer.newToken(PARANTHESES_CLOSE, string(x))

	case '{':
		return tokenizer.newToken(BRACE_OPEN, string(x))

	case '}':
		return tokenizer.newToken(BRACE_CLOSE, string(x))

	case '.':
		return tokenizer.newToken(FULLSTOP, string(x))

	case ',':
		return tokenizer.newToken(COMMA, string(x))

	case ':':
		return tokenizer.newToken(COLON, string(x))

	case '_':
		return tokenizer.newToken(VARIABLE, string(x))

	case '?':
		if tokenizer.nextInline(1) == '>' {
			tokenizer.tokenizerText = true
			return tokenizer.newToken(SWASH_BODY_CLOSE, "?>")
		}

	case '>':
		if tokenizer.nextInline(1) == '=' {
			return tokenizer.newToken(GREATEREQUAL, ">=")
		}

		return tokenizer.newToken(GREATERTHAN, ">")

	case '!':
		if tokenizer.nextInline(1) == '=' {
			return tokenizer.newToken(NOTEQUAL, "!=")
		}

	case '<':
		if tokenizer.nextInline(1) == '=' {
			return tokenizer.newToken(LESSEQUAL, "<=")
		} else if strings.HasPrefix(tokenizer.tokenizerTarget[tokenizer.tokenizerLine], WORKSPACE_OPEN) {
			return tokenizer.newToken(SWASH_BODY_OPEN, WORKSPACE_OPEN)
		}

		return tokenizer.newToken(LESSTHAN, "<")

	case '=': // Equal operator
		if tokenizer.nextInline(1) == '=' {
			return tokenizer.newToken(DOUBLE_EQUAL, "==")
		}

		return tokenizer.newToken(EQUAL, string(x))

	case '$': // Dollar operator
		if unicode.IsLetter(tokenizer.nextInline(1)) {
			literal := string(x) + tokenizer.runUntilVoid(1, indentCharType)
			if tokenizer.nextInline(len(literal)) == '[' {
				return tokenizer.tokenizeArgs(literal)
			}

			return tokenizer.newToken(VARIABLE, literal)
		}

		return tokenizer.newToken(DOLLAR, string(x))

	case '"', '\'': // String operator
		tokenizer.shift(1)
		inspector := tokenizer.runUntilVoid(0, singleCharTyper(rune(x)))
		tokenizer.shift(1)

		return tokenizer.newToken(STRING, inspector)

	case '@': // tag reference
		return tokenizer.newToken(TAGINDENT, fmt.Sprintf("@%s", tokenizer.runUntilVoid(1, indentCharType)))

	default: // Indent, int & space operators
		if unicode.IsLetter(rune(x)) {
			inspector := tokenizer.runUntilVoid(0, indentCharType)
			if inspector == "true" || inspector == "false" {
				return tokenizer.newToken(BOOLEAN, inspector)
			}

			switch tokenizer.nextInline(len(inspector)) {

			case '.', '(': // func/deep var
				return tokenizer.tokenizeArgs(inspector)

			default: // indent
				return tokenizer.newToken(INDENT, inspector)
			}
		} else if unicode.IsDigit(rune(x)) {
			return tokenizer.newToken(NUMBER, tokenizer.runUntilVoid(0, mergeFromPackage(unicode.IsDigit)))
		} else if unicode.IsSpace(rune(x)) {
			return tokenizer.newToken(SPACE, string(x)).NoState()
		}
	}

	return nil
}

// runUntilVoid will continue until the checker func returns false
func (tokenizer *Tokenizer) runUntilVoid(i int, checker ...typeChecker) string {
	content := bytes.NewBuffer(make([]byte, 0))
	for i, char := range tokenizer.tokenizerTarget[tokenizer.tokenizerLine][i:] {
		if !typeCheckers(checker, char, i) {
			break
		}

		content.WriteRune(char)
	}

	return content.String()
}

// Nodes will return all the nodes within the tokenizer guidelines
func (tokenizer *Tokenizer) Nodes() []ParserNode {
	return tokenizer.tokenizerParserNodes
}
