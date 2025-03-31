package packages

import (
	"fmt"
	"time"
)

// Since will calculate the amount of time since the unix timestamp passed
func Since(unix int) string {
	periodOfTime := time.Since(time.Unix(int64(unix), 0)).Round(1 * time.Second)
	units := map[string]float64{
		"years":   ((periodOfTime.Hours() / 24) / 30.4) / 12,
		"months":  (periodOfTime.Hours() / 24) / 30.4,
		"days":    periodOfTime.Hours() / 24,
		"hours":   periodOfTime.Hours(),
		"minutes": periodOfTime.Minutes(),
		"seconds": periodOfTime.Seconds(),
	}

	/* iterates over each unit inside the map and compares it's time */
	for key, unit := range units {
		if unit >= 1 {
			return fmt.Sprintf("%2.f", unit) + key
		}
	}

	return "0secs"
}
