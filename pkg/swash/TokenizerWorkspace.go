package swash

import (
	"bytes"
	"strings"
)

// newWorkspace will attempt to create a new workspace for the tokenizer8806
func (tokenizer *Tokenizer) newWorkspace() *Token {
	if !strings.Contains(tokenizer.tokenizerTarget[tokenizer.tokenizerLine], "<?") {
		if tokenizer.tokenizerLine+1 >= len(tokenizer.tokenizerTarget) {
			return tokenizer.newToken(TEXT, tokenizer.tokenizerTarget[tokenizer.tokenizerLine])
		}

		return tokenizer.newToken(TEXT, tokenizer.tokenizerTarget[tokenizer.tokenizerLine], &Token{TokenType: TEXT_BLANKLINE})
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	/* iterates over the src */
	for p, line := range tokenizer.tokenizerTarget[tokenizer.tokenizerLine] {
		if line == '<' && strings.HasPrefix(tokenizer.tokenizerTarget[tokenizer.tokenizerLine][p:], "<?swash") {
			tokenizer.tokenizerText = false
			break
		}

		buf.WriteRune(line)
	}

	return tokenizer.newToken(TEXT, buf.String())
}
