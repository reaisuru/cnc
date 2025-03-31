package simpleconfig

import "cnc/pkg/simpleconfig/encoders"

type Coder int

const (
	Yaml = encoders.Yaml
	Toml = encoders.Toml
	Json = encoders.Json
)

type SimpleConfig struct {
	directory     string
	coder         Coder
	fileExtension string
}

func New(coder Coder, directory string) *SimpleConfig {
	var fileExtension string
	switch coder {
	case Json:
		fileExtension = ".json"
	case Yaml:
		fileExtension = ".yml"
	case Toml:
		fileExtension = ".toml"
	}

	return &SimpleConfig{directory: directory, coder: coder, fileExtension: fileExtension}
}
