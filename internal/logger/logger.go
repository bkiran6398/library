package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	logTimeStamp string = time.RFC3339
)

// New creates a configured zerolog.Logger.
// level: debug|info|warn|error.
func New(level string) (zerolog.Logger, error) {
	lvl, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("failed to parse log level %s: %w", level, err)
	}
	zerolog.SetGlobalLevel(lvl)
	zerolog.TimeFieldFormat = logTimeStamp

	return zerolog.New(os.Stdout).With().Timestamp().Logger(), nil
}
