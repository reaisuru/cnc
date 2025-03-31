package clients

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	nameSanitizeRegex = regexp.MustCompile(`[\r\n\t\f\v]+`)
	validNameRegex    = regexp.MustCompile(`^[a-zA-Z.0-9]+$`)

	nameMap = map[string]string{
		"weed.mnt":     "ktcctv",
		"multidvr.mnt": "multidvr",
	}
)

func (bot *Bot) sanitizeName() error {
	if bot.Name == "'" {
		bot.Name = "unnamed"
		return nil
	}

	// bot doesn't have a name, just call it unnamed
	if len(bot.Name) <= 2 {
		bot.Name = "unnamed"
		return nil
	}

	// don't need uppercase bot names
	bot.Name = strings.ToLower(bot.Name)
	bot.Name = nameSanitizeRegex.ReplaceAllString(bot.Name, "")

	// check if it's still not valid
	if !validNameRegex.MatchString(bot.Name) {
		return fmt.Errorf("%w: '%s'", ErrInvalidSource, bot.Name)
	}

	// remap it
	for oldName, newName := range nameMap {
		if strings.HasPrefix(bot.Name, oldName) {
			bot.Name = newName
		}
	}

	return nil
}
