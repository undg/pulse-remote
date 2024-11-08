package json

import (
	"fmt"
	"net/http"
)

func ServeStatusTypeJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")


	fmt.Fprint(w, "")
}
