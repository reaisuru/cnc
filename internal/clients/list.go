package clients

import (
	"cnc/pkg/archs"
	"cnc/pkg/logging"
	"fmt"
	"golang.org/x/exp/slices"
)

// Add will add a clients to our client list
func Add(bot *Bot) {
	mutex.Lock()
	defer mutex.Unlock()

	ID++
	bot.ID = ID
	List[bot.ID] = bot

	logging.Global.Info().
		Str("name", bot.Name).
		Str("arch", archs.Arch(bot.Arch).String()).
		Str("version", fmt.Sprintf("%d.%d.%d", bot.Version.Major, bot.Version.Minor, bot.Version.Patch)).
		IPAddr("address", bot.Address).
		Msg("Added client")
}

// Delete will remove a clients from our client list
func Delete(bot *Bot) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(List, bot.ID)
	logging.Global.Info().
		Str("name", bot.Name).
		IPAddr("address", bot.Address).
		Msg("Removed client")
}

// Count returns the current clients count
func Count() int {
	mutex.Lock()
	defer mutex.Unlock()

	return len(List)
}

// Distribution will show the bot distribution by group name(s).
func Distribution() map[string]int {
	var dist = make(map[string]int)

	defer mutex.Unlock()
	mutex.Lock()

	for _, bot := range List {
		dist[bot.Name]++
	}

	return dist
}

// ASN will show the bot distribution by ASN
func ASN() map[string]int {
	var dist = make(map[string]int)

	defer mutex.Unlock()
	mutex.Lock()

	for _, bot := range List {
		dist[bot.ASN]++
	}

	return dist
}

// Architectures will show the bot distribution by architectures.
func Architectures() map[string]int {
	var dist = make(map[string]int)

	defer mutex.Unlock()
	mutex.Lock()

	for _, bot := range List {
		dist[archs.Arch(bot.Arch).String()]++
	}

	return dist
}

// Countries will show the bot distribution by countries.
func Countries() map[string]int {
	var dist = make(map[string]int)

	defer mutex.Unlock()
	mutex.Lock()

	for _, bot := range List {
		dist[bot.Country]++
	}

	return dist
}

// Groups returns all available bot groups.
func Groups() []string {
	var sources []string

	for source, _ := range Distribution() {
		if slices.Contains(sources, source) {
			continue // dont need dupes
		}

		sources = append(sources, source)
	}

	return sources
}
