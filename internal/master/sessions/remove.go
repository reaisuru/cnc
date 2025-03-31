package sessions

func (s *Session) Remove() {
	delete(sessions, s.ID)
}
