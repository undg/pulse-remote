package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/undg/go-prapi/buildinfo"
)

var logger zerolog.Logger

var (
	// export these methods directly
	Trace = logger.Trace // (most noise)
	Debug = logger.Debug
	Info  = logger.Info
	Warn  = logger.Warn
	Error = logger.Error
	Fatal = logger.Fatal
	Panic = logger.Panic
)

// logger based on zerolog
func init() {
	bi := buildinfo.Get()

	debug := os.Getenv("DEBUG")

	// @TODO (undg) 2025-02-14: convert number string to number and check if it's greater than TRACE level, so DEBUG=999 could be valid max.
	switch debug {
	case "TRACE", "3":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "DEBUG", "2":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO", "1":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN", "0":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERR", "-1":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Str("commit", bi.GitCommit).
		Logger()
}
