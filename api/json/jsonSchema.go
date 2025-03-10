package json

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/danielgtaylor/huma/schema"
	"github.com/undg/go-prapi/api/logger"
	"github.com/undg/go-prapi/api/pactl"
)

func serveSchemaJSON(w http.ResponseWriter, t reflect.Type) {
	w.Header().Set("Content-Type", "application/json")

	s, err := schema.Generate(t)
	if err != nil {
		logger.Error().Err(err).Msg("Schema error")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	b, err := json.Marshal(s)
	if err != nil {
		logger.Error().Err(err).Msg("Internal Server Error")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(b))
}

func serveRestJSON(w http.ResponseWriter, restJSON interface{}) {

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(restJSON); err != nil {
		logger.Error().Err(err).Msg("json.NewEncoder(w).Encode(restJSON)")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ServeStatusSchemaJSON(w http.ResponseWriter, r *http.Request) {
	serveSchemaJSON(w, reflect.TypeOf(pactl.Status{}))
}

func ServeMessageSchemaJSON(w http.ResponseWriter, r *http.Request) {
	serveSchemaJSON(w, reflect.TypeOf(Message{}))
}

func ServeResponseSchemaJSON(w http.ResponseWriter, r *http.Request) {
	serveSchemaJSON(w, reflect.TypeOf(Response{}))
}

func ServeStatusRestJSON(w http.ResponseWriter, r *http.Request) {
	serveRestJSON(w, pactl.GetStatus())
}
