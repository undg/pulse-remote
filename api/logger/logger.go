package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

// logger based on zerolog
func init() {

	debug := strings.ToUpper(os.Getenv("DEBUG"))

	level := zerolog.InfoLevel

	// zerolog allows for logging at the following levels (from highest to lowest):
	//
	// panic (zerolog.PanicLevel, 5)
	// fatal (zerolog.FatalLevel, 4)
	// error (zerolog.ErrorLevel, 3)
	// warn (zerolog.WarnLevel, 2)
	// info (zerolog.InfoLevel, 1)
	// debug (zerolog.DebugLevel, 0)
	// trace (zerolog.TraceLevel, -1)
	switch debug {
	case "TRACE", "-1":
		level = zerolog.TraceLevel
	case "DEBUG", "0":
		level = zerolog.DebugLevel
	case "INFO", "1":
		level = zerolog.InfoLevel
	case "WARN", "2":
		level = zerolog.WarnLevel
	case "ERR", "3":
		level = zerolog.ErrorLevel
	case "FATAL", "4":
		level = zerolog.FatalLevel
	case "PANIC", "5":
		level = zerolog.PanicLevel
	default:
		level = zerolog.InfoLevel
		debug = "INFO"
	}

	logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	Trace = logger.Trace
	Debug = logger.Debug
	Info = logger.Info
	Warn = logger.Warn
	Error = logger.Error
	Fatal = logger.Fatal
	Panic = logger.Panic
	GetLevel = logger.GetLevel
	DebugEnv = debug
}

var (
	Trace    func() *zerolog.Event
	Debug    func() *zerolog.Event
	Info     func() *zerolog.Event
	Warn     func() *zerolog.Event
	Error    func() *zerolog.Event
	Fatal    func() *zerolog.Event
	Panic    func() *zerolog.Event
	GetLevel func() zerolog.Level
	DebugEnv string
)
