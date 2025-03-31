package simpleconfig

import (
	"cnc/pkg/simpleconfig/encoders"
	"os"
)

func (s *SimpleConfig) encode(path string, overwrite bool, v interface{}) error {
	if !overwrite {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	err = encoders.NewEncoder(int(s.coder), file).Encode(v)
	if err != nil {
		return err
	}

	return nil
}
