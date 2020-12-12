package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log *zerolog.Logger

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	log = &logger
}

func Log() *zerolog.Logger {
	return log
}

func SetLevel(level zerolog.Level) {
	l := log.Level(level)
	log = &l
}
