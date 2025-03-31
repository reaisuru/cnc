package swash

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
	TokenizerArgs.go is an official way of parsing the implementation args of a statement, you can find arguments within
	certain areas of the internal code of a program, these areas include function call & variable declarations (there are
	several other locations where they can also occur), This function will create a customizable interpreter instruction
	which can be decode and executed quicker with a high accuracy level than the standard method.

	Argument definition. An argument is a way for you to provide more information to a function. The function can then use
	that information as it runs, like a variable. Said differently, when you create a function, you can pass in data in
	the form of an argument, also called a parameter.

	Another definition for TokenizerArgs.go would be how it's used to parse the array elements within the internal code. it
	parses all the parameters for an array within it's self contained AST branch.
*/

// tokenizeArgs will help within the function for handling arguments directly
func (tokenizer *Tokenizer) tokenizeArgs(indent string) *Token {
	line := strings.Split(tokenizer.tokenizerTarget[tokenizer.tokenizerLine], "")

	for i := len(indent); i < len(line); i++ {
		index := line[i]
		switch index {

		case "(": // implements body specific parsing
			indexScale := 0

			// whenever indexScale equals to 0, everything is ok
			for i := utf8.RuneCountInString(indent); i < len(line); i++ {
				char := line[i]
				indent += string(char)

				switch char {

				case "(":
					indexScale++
					continue

				case ")":
					if indexScale == 1 {
						break
					}

					indexScale--
					continue

				default:
					continue
				}

				break
			}

			/* processes the arguments for the func */
			args := make([]*Token, 0)
			t := NewTokenizer(indent[i+1:][:len(indent[i+1:])-1], false)

			for {
				token, ok := t.PeekNext()
				if !ok || token == nil {
					break
				}

				args = append(args, token)
			}

			return tokenizer.newToken(FUNCTION, indent, args...)

		case "$": // variable
			indent += string(index)
			indent += tokenizer.runUntilVoid(i+1, mergeFromPackage(unicode.IsLetter))
			return tokenizer.newToken(VARIABLE, indent)

		default:
			indent += string(index)
		}
	}

	return nil
}
