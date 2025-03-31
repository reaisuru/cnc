package simpleconfig

import (
	"path/filepath"
)

func (s *SimpleConfig) Unmarshal(v interface{}, file string) error {
	return s.decode(filepath.Join(s.directory, file+s.fileExtension), v)
}

func (s *SimpleConfig) Marshal(v interface{}, file string) error {
	return s.encode(filepath.Join(s.directory, file+s.fileExtension), true, v)
}
