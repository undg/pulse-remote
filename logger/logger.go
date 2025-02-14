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
	Debug = logger.Debug
	Info  = logger.Info
	Warn  = logger.Warn
	Error = logger.Error
	Fatal = logger.Fatal
	Panic = logger.Panic
)

func init() {
	bi := buildinfo.Get()

	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Str("commit", bi.GitCommit).
		Logger()
}
