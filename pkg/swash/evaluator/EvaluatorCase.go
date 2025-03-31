package evaluator

import (
	"cnc/pkg/swash"
)

// decision will execute the if statements directly.
func (evaluator *Evaluator) decision(index *swash.If, leaves *swash.TokenType) (bool, error) {
	context, err := evaluator.Memory.args(index.Decision)
	if err != nil || context == nil || context.TokenType != swash.BOOLEAN || !context.TokenLiteralValue.(bool) {
		return false, err
	}

	object, ok, err := evaluator.self(index.Body.Nodes(), leaves, false)
	if err != nil || object == nil {
		return ok, err
	}

	evaluator.Closed = object
	return ok, nil
}
