package packages

import (
	"time"
)

// TIME is the require path
const TIME = "date/time"

// Swash Time Package (Ported from Golang)
//
// This package provides essential time-related functionality for the Swash programming language,
// incorporating the time package originally written in Golang. It offers a range of functions and
// utilities for handling time and date calculations, parsing and formatting time strings, and
// performing various time-related operations.
//
// By porting the Golang time package to Swash, developers can seamlessly leverage its robust
// features and APIs to efficiently work with time-related data in their Swash applications.

// TimeFunctions is a map[string]any that serves as a registry for all time-related functions.
// The keys in this map represent the names of the available time functions, while the values are the corresponding function objects.
// This map allows easy access to various time operations, such as time parsing, formatting, duration calculations, and more.
// By maintaining all time functions in a single map, it provides a centralized and convenient way to utilize time functionality in Swash.
// Developers can simply look up the desired function using its name as a key and invoke it as needed.
var TIMEFunctions = map[string]any{
	"microsecond": int(time.Microsecond),
	"millisecond": int(time.Millisecond),
	"second":      int(time.Second),
	"minute":      int(time.Minute),
	"hour":        int(time.Hour),

	// sleep will sleep for the duration of time provided
	"sleep": func(duration int) {
		time.Sleep(time.Duration(duration))
	},

	// now returns the current Unix timestamp representing the current time in seconds since January 1, 1970 UTC.
	"now": func() int {
		return int(time.Now().Unix())
	},

	// unix returns a formatted time representation from a unix timestamp.
	"unix": func(unix int, format string) string {
		return time.Unix(int64(unix), 0).Format(format)
	},

	// hour returns the current hour in the day
	"hours": func(unix int) int {
		return time.Unix(int64(unix), 0).Hour()
	},

	// minute returns the current minute in the day
	"minutes": func(unix int) int {
		return time.Unix(int64(unix), 0).Minute()
	},

	// second returns the current second in the day
	"seconds": func(unix int) int {
		return time.Unix(int64(unix), 0).Second()
	},

	// until returns a formatted time representation of how long until a unix timestamp.
	"until": Until,

	// since returns a formatted time representation of how long since  a unix timestamp
	"since": Since,
}
