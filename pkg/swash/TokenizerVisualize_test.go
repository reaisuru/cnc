package swash_test

import (
	"cnc/pkg/swash"
	"strings"
	"testing"
)

/*
	TokenizerVisualize_test.go allows for the tokens which have been tokenized concurrently to be visualized in a format
	which humans can understand easily
*/

func TestTokenizerVisualize(t *testing.T) {
	tokenizer, err := swash.NewTokenizerSourcedFromFile("examples/example.tfx")
	if err != nil {
		t.Fatal(err)
	}

	tokens := make([]*swash.Token, 0)

	for {
		token, ok := tokenizer.PeekNext()
		if !ok || token == nil {
			break
		}

		tokens = append(tokens, token)
	}

	iterateTokensPrint(tokens, 0, t)
}

// iterateTokensPrint will print all the tokens within the array
func iterateTokensPrint(tokens []*swash.Token, depth int, t *testing.T) {
	for _, token := range tokens {
		switch token.TokenType {

		default:
			t.Logf("%-12s"+strings.Repeat("\t", depth)+token.TokenLiteral, token.TokenType.String())

		case swash.FUNCTION:
			t.Logf("%-12s"+strings.Repeat("\t", depth)+strings.Split(token.TokenLiteral, "(")[0]+"", token.TokenType.String())
		}

		if len(token.TokenArgs) > 0 {
			iterateTokensPrint(token.TokenArgs, depth+1, t)
		}
	}
}
