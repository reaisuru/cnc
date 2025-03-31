package pattern

import (
	"path/filepath"
	"strings"
)

func Filter(possibles []string, patterns []string) []string {
	if len(possibles) <= 0 || len(patterns) <= 0 {
		return make([]string, 0)
	}

	var matchedThings []string

	inclusionPatterns, exclusionPatterns := separatePatterns(patterns)
	for _, thing := range possibles {
		if Matches(thing, inclusionPatterns) {
			if !Matches(thing, exclusionPatterns) {
				matchedThings = append(matchedThings, thing)
			}
		}
	}

	return matchedThings
}

func separatePatterns(patterns []string) (inclusion []string, exclusion []string) {
	for _, pattern := range patterns {
		if strings.HasPrefix(pattern, "!") {
			exclusion = append(exclusion, pattern[1:])
		} else {
			inclusion = append(inclusion, pattern)
		}
	}

	return
}

func Matches(v string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, v)
		if err != nil {
			continue
		}

		if matched {
			return true
		}
	}

	return false
}
