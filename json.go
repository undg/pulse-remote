package main

import (
	"encoding/json"
	"log"
)

type Audio struct {
	volume float32
	mute   bool
}

type Result struct {
	Audio
	schema   string
	response string
	error    string
}

func marshalResult(a Result) []byte {
	jsonData, err := json.Marshal(map[string]interface{}{
		"schema":   a.schema,
		"response": a.response,
		"error":    a.error,
	})

	if err != nil {
		log.Println(err)
	}

	return jsonData
}

func (a Result) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(map[string]interface{}{
		"schema":   a.schema,
		"response": a.response,
		"error":    a.error,
	})

	return jsonData, err
}
