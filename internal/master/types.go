package master

import (
	"golang.org/x/crypto/ssh"
)

type Listener struct {
	Port       int
	PrivateKey string

	server *ssh.ServerConfig
}
