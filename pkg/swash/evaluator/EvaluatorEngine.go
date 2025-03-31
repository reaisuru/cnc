package evaluator

import (
	"cnc/pkg/swash"
	"io"

	"github.com/valyala/fasttemplate"
)

/*
	EvaluatorEngine.go implements the workspace functions for the evaluators
*/

// evaluates the current expression and returns the context required
func (evaluator *Evaluator) evaluate(expression *swash.Token) error {
	expr, err := fasttemplate.ExecuteFuncStringWithErr(expression.TokenLiteral, swash.WORKSPACETAG_OPEN, swash.WORKSPACETAG_CLOSE, func(wr io.Writer, tag string) (int, error) {
		context := swash.NewTokenizer(tag, false)
		if err := context.Parse(); err != nil || len(context.Nodes()) == 0 {
			return 0, err
		}

		switch object := context.Nodes()[0].(type) {

		case *swash.Token:
			content, err := evaluator.Memory.evalToken(object)
			if err != nil || content == nil {
				return 0, err
			}

			return wr.Write([]byte(content.TokenLiteral))
		}

		return 0, nil
	})

	if err != nil {
		return err
	}

	_, err = evaluator.Memory.wr.Write([]byte(expr))
	return err
}
