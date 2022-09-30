package meta

import "time"

var (
	Version   = "dev"
	BuildDate = time.Now()
	date      = ""
)

func init() {
	if buildDate, err := time.Parse(time.RFC3339, date); err == nil {
		BuildDate = buildDate
	}
}
