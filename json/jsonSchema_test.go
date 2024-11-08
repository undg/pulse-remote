package json

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeStatusSchemaJSON(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/any-endpoint", nil)

	ServeStatusSchemaJSON(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("[Err] Expected status code %d but bot %d", http.StatusOK, resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("[Err] Expected Content-Type json, got %v", resp.Header.Get("Content-Type"))
	}

	var js map[string]interface{}
	if err := json.Unmarshal(body, &js); err != nil {
		t.Errorf("[Err] Response not valid JSON: %v", err)
	}
}
