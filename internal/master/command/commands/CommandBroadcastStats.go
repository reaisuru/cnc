package commands

import (
	"cnc/internal/database"
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"cnc/pkg/pattern"
	"github.com/dustin/go-humanize"
)

func init() {
	command.Create(&command.Command{
		Aliases:     []string{"bcstats", "broadcaststats"},
		Description: "Shows the statistics of your last attack.",
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			log, err := database.Logs.LastByUserID(session.UserProfile.ID)
			if err != nil {
				return session.Notification("There is no recent attack of yours.")
			}

			_ = session.Printfln("Statistics of your last attacks:")
			_ = session.Printfln("- Attack ID: %d", log.AttackID)
			_ = session.Printfln("- Targets: %s", pattern.FormatArray(log.Targets))
			_ = session.Printfln("- Duration: %ds", log.Duration)
			_ = session.Printfln("- Clients: %d", log.Clients)
			_ = session.Printfln("- Started: %s", humanize.Time(log.Started))
			_ = session.Printfln("- Ended: %s", humanize.Time(log.Started))

			return nil
		},
	})
}
