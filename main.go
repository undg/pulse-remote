package main

import (
	"embed"
	"fmt"
	"io"
	"net/http"

	"github.com/undg/go-prapi/api/buildinfo"
	prapiJSON "github.com/undg/go-prapi/api/json"
	"github.com/undg/go-prapi/api/logger"
	"github.com/undg/go-prapi/api/utils"
	"github.com/undg/go-prapi/api/ws"
)

// @TODO (undg) 2024-10-06: different port for dev and production

//go:embed web/dist/*
//go:embed web/dist/assets/*
//go:embed web/dist/fonts/*
//go:embed web/dist/icons/*
var prWebDist embed.FS

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

	// Static files
	fsys := http.FileServer(http.FS(prWebDist))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := "build/pr-web/dist" + r.URL.Path
		_, err := prWebDist.Open(path)
		if err != nil {
			// File not exist, serve index.html
			w.Header().Set("Content-Type", "text/html")
			indexFile, _ := prWebDist.Open("build/pr-web/dist/index.html")
			io.Copy(w, indexFile)
			return
		}
		// File exists, serve it
		r.URL.Path = path
		fsys.ServeHTTP(w, r)
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
