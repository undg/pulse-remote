package json

import (
	"encoding/json"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	response := Response{
		Action:  string(ActionGetStatus),
		Status:  StatusSuccess,
		Payload: "test payload",
	}

	expected := `{"action":"GetStatus","payload":"test payload","status":4000}`

	assertJSON(t, response, expected)
}

func TestMarshalJSONWithError(t *testing.T) {
	response := Response{
		Action: string(ActionGetStatus),
		Status: StatusError,
		Error:  "test error",
	}

	expected := `{"action":"GetStatus","error":"test error","status":4001}`

	assertJSON(t, response, expected)
}

func assertJSON(t *testing.T, response Response, expected string) {
	result, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("ERROR marshaling JSON: %v", err)
	}

	if string(result) != expected {
		t.Errorf("\nExpected %s\nGot      %s", expected, result)
	}
}
