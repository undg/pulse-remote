package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/undg/go-prapi/buildinfo"
	prapiJSON "github.com/undg/go-prapi/json"
	"github.com/undg/go-prapi/logger"
	"github.com/undg/go-prapi/utils"
	"github.com/undg/go-prapi/ws"
)

// @TODO (undg) 2024-10-06: different port for dev and production

func startServer(mux *http.ServeMux) {
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/schema/status":
			prapiJSON.ServeStatusSchemaJSON(w, r)
		case "/api/v1/schema/message":
			prapiJSON.ServeMessageSchemaJSON(w, r)
		case "/api/v1/schema/response":
			prapiJSON.ServeResponseSchemaJSON(w, r)
		case "/api/v1/status":
			prapiJSON.ServeStatusRestJSON(w, r)
		case "/api/v1/ws":
			ws.HandleWebSocket(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	fs := http.FileServer(http.Dir("/tmp/bin/pr-web/dist"))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat("/tmp/bin/pr-web/dist" + r.URL.Path); os.IsNotExist(err) {
			http.ServeFile(w, r, "/tmp/bin/pr-web/dist/index.html")
		} else {
			fs.ServeHTTP(w, r)
		}
	}))
}

func main() {
	ip, err := utils.GetLocalIP()
	if err != nil {
		logger.Error().Err(err).Msg("can't GetLocalIP()")
	}
	b := buildinfo.Get()

	fmt.Print(`
┌───────────────────────────────────────────────────┐
│                     GO-PRAPI                      │
├───────────────────────────────────────────────────┤
  GitVersion: `, b.GitVersion, `
  GitCommit:  `, b.GitCommit, `
  BuildDate:  `, b.BuildDate, `
  Compiler:   `, b.Compiler, `
  Platform:   `, b.Platform, `
  GoVersion:  `, b.GoVersion, `
  LogLevel:   `, logger.GetLevel(), `
  DEBUG:      `, logger.DebugEnv, `
└───────────────────────────────────────────────────┘
`)
	fmt.Println("\n🔥 Igniting server on ws://" + ip + utils.PORT)
	fmt.Println("🔥 WebApp http://" + ip + utils.PORT + "\n")

	fmt.Print(`──────────────────────────────────────────────────────────────
`)
	logger.Trace().Str("Trace", "ON").Msg("Log LEVEL")
	logger.Debug().Str("Debug", "ON").Msg("Log LEVEL")
	logger.Info().Str("Info", "ON").Msg("Log LEVEL")
	logger.Warn().Str("Warn", "ON").Msg("Log LEVEL")
	logger.Error().Str("Error", "ON").Msg("Log LEVEL")

	fmt.Print(`──────────────────────────────────────────────────────────────

`)

	mux := http.NewServeMux()

	startServer(mux)

	go ws.BroadcastUpdates()

	errListenAndServe := http.ListenAndServe(utils.PORT, mux)
	if errListenAndServe != nil {
		logger.Fatal().Err(errListenAndServe).Msg("server failed to start")
	}
}
