package buildinfo

import "fmt"

var (
	Version = "dev"
	Date    = "unknown"
)

func String() string {
	return fmt.Sprintf("confy version %s\nbuilt at: %s", Version, Date)
}
