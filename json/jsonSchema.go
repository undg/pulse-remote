package json

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/danielgtaylor/huma/schema"
	"github.com/undg/go-prapi/pactl"
)

func RenderSchemaJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s, err := schema.Generate(reflect.TypeOf(pactl.Status{}))
	if err != nil {
		log.Println("ERROR ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	b, err := json.Marshal(s)
	if err != nil {
		log.Println("ERROR ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(b))
}
