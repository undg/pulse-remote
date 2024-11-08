package json

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeStatusTypeJSON(t * testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/some-url", nil)


	ServeStatusTypeJSON(w, r)

	resp := w.Result()

	body,_ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("[Err] StatusCode should be %d but is %d", http.StatusOK, resp.StatusCode)
	}

	if resp.Header.Get("Content-type") != "text/plain" {
		t.Fatalf("[Err] Content-type should be text/plain but is %s", resp.Header.Get("Content-type"))
	}

	if !strings.Contains(string(body), "export") {
		t.Fatalf("[Err] Body should contain TypeScript types with 'export' keyword")
	}
}
