package format

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ModifiedParseDuration(duration string) (time.Duration, error) {
	regexMap := map[string]func(int) time.Duration{
		`(\d+)\s*s\b`:   func(s int) time.Duration { return time.Duration(s) * time.Second },
		`(\d+)\s*m\b`:   func(m int) time.Duration { return time.Duration(m) * time.Minute },
		`(\d+)\s*h\b`:   func(h int) time.Duration { return time.Duration(h) * time.Hour },
		`(\d+)\s*d\b`:   func(d int) time.Duration { return time.Duration(d) * 24 * time.Hour },
		`(\d+)\s*w\b`:   func(w int) time.Duration { return time.Duration(w) * 7 * 24 * time.Hour },
		`(\d+)\s*mo\b`:  func(mo int) time.Duration { return time.Duration(mo) * 30 * 24 * time.Hour },
		`(\d+)\s*y\b`:   func(y int) time.Duration { return time.Duration(y) * 365 * 24 * time.Hour },
		`(\d+)\s*dec\b`: func(dec int) time.Duration { return time.Duration(dec) * 10 * 365 * 24 * time.Hour },
	}

	for regex, durationFn := range regexMap {
		re := regexp.MustCompile(regex)
		matches := re.FindAllStringSubmatch(duration, -1)
		for _, match := range matches {
			value, err := strconv.Atoi(match[1])
			if err != nil {
				return 0, fmt.Errorf("invalid duration: %s", duration)
			}
			replaceStr := durationFn(value).String()
			duration = strings.Replace(duration, match[0], replaceStr, 1)
		}
	}

	parsedDuration, err := time.ParseDuration(duration)
	if err != nil {
		return 0, err
	}

	return parsedDuration, nil
}
