package evaluator

import (
	"bytes"
	"cnc/pkg/swash"
	"errors"
	"io"
)

/*
	Evaluator.go is the main file which helps within the guide rail interface, this will range
	over all the nodes within the internal code and attempt to generate responses to them all,
	but it doesn't have any sort of order of operations as we don't want this to be implemented.
*/

// Evaluator is an object which implements the package parameters
type Evaluator struct {
	tokenizer *swash.Tokenizer
	Memory    *Memory
	Closed    *swash.Token
	standard  *standard
}

// NewEvaluator creates a new Evaluator object
func NewEvaluator(tokenizer *swash.Tokenizer, wr io.Writer, rd io.Reader) *Evaluator {
	eval := &Evaluator{
		tokenizer: tokenizer,
		standard:  new(standard),
		Memory:    NewMemory(wr, rd, make(map[string]any)),
	}

	/* allocates memory for the registers */
	eval.standard.evaluator = eval
	eval.standard.register()
	return eval
}

// Execute is the main keyword which will execute the internal code
func (evaluator *Evaluator) Execute() error {
	_, _, err := evaluator.self(evaluator.tokenizer.Nodes(), nil, false)
	return err
}

// self is a keyword powered runtime evaluator which helps implement the main process evaluation
func (evaluator *Evaluator) self(objects []swash.ParserNode, leaves *swash.TokenType, forced bool) (*swash.Token, bool, error) {
	for _, object := range objects {
		if evaluator.Closed != nil {
			return evaluator.Closed, true, nil
		}

		/* once we check for a closure, we continue looping */
		switch index := object.(type) {

		case *swash.If:
			ok, err := evaluator.decision(index, leaves)
			if err != nil || evaluator.Closed != nil || ok {
				return evaluator.Closed, ok, err
			}

		case *swash.Var:
			if err := evaluator.Memory.allocateVar(index); err != nil {
				return nil, false, err
			}

		case *swash.Function:
			if err := evaluator.Memory.allocateFunc(index); err != nil {
				return nil, false, err
			}

		case *swash.Token:
			switch index.TokenType {

			case swash.TEXT: // evaluate the expression
				payload := bytes.NewBuffer(make([]byte, 0))
				payload.WriteString(index.TokenLiteral)
				if len(index.TokenArgs) > 0 && index.TokenArgs[0].TokenType == swash.TEXT_BLANKLINE {
					payload.Write([]byte("\r\n"))
				}

				index.TokenLiteral = payload.String()
				if err := evaluator.evaluate(index); err != nil {
					return nil, false, err
				}

			case swash.VARIABLE, swash.FUNCTION: // evaluate the variable or function
				object, err := evaluator.Memory.evalToken(index)
				if err != nil {
					return nil, false, err
				}

				if object == nil {
					continue
				}

				evaluator.Memory.wr.Write([]byte(object.TokenLiteral))
			}

		case *swash.Return:
			if len(index.Args) == 0 {
				return nil, true, nil
			}

			context, err := evaluator.Memory.args(index.Args)
			if err != nil || context == nil {
				return nil, false, err
			}

			if *leaves != context.TokenType {
				return nil, false, errors.New("mismatched types")
			}

			return context, true, nil
		}
	}

	return nil, false, nil
}
