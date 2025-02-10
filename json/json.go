package json

import (
	"encoding/json"
)

type Action string

const (
	// Get composed informations about all sinks, sources, inputs and build
	ActionGetStatus Action = "GetStatus"

	// Metadata about build
	ActionGetBuildInfo Action = "GetBuildInfo"

	// SINKS, e.g. Speakers
	ActionSetSinkVolume Action = "SetSinkVolume"
	ActionSetSinkMuted  Action = "SetSinkMuted"

	// Apps playing audio
	ActionSetSinkInputVolume Action = "SetSinkInputVolume"
	ActionSetSinkInputMuted  Action = "SetSinkInputMuted"
	// Move App to different SINK
	ActionMoveSinkInput Action = "MoveSinkInput"

	// SOURCES, e.g. Microphones
	ActionSetSourceVolume Action = "SetSourceVolume"
	ActionSetSourceMuted  Action = "SetSourceMuted"

	// Apps active access to microphones
	ActionSetSourceInputVolume Action = "SetSourceInputVolume"
	ActionSetSourceInputMuted  Action = "SetSourceInputMuted"
	// Move App to different SOURCE
	ActionMoveSourceOutput Action = "MoveSourceOutput"
)

var AvailableCommands = []Action{
	// Get composed informations about all sinks, sources, inputs and build
	ActionGetStatus,

	// Metadata about build
	ActionGetBuildInfo,

	// SINKS, e.g. Speakers
	ActionSetSinkVolume,
	ActionSetSinkMuted,

	// Apps playing audio
	ActionSetSinkInputVolume,
	ActionSetSinkInputMuted,
	// Move App to different SINK
	ActionMoveSinkInput,

	// SOURCES, e.g. Microphones
	ActionSetSourceVolume,
	ActionSetSourceMuted,

	// Apps active access to microphones
	ActionSetSourceInputVolume,
	ActionSetSourceInputMuted,
	// Move App to different SOURCE
	ActionMoveSourceOutput,
}

// @TODO (undg) 2025-02-10: generate enum's. Check https://github.com/danielgtaylor/huma

// Message is an request from the client
type Message struct {
	// Actions listed in availableCommands slice
	Action Action `json:"action" doc:"Action to perform fe. GetVolume, SetVolume, SetMute..."enum:"GetStatus,GetBuildInfo,SetSinkVolume,SetSinkMuted,SetSinkInputVolume,SetSinkInputMuted,MoveSinkInput,SetSourceVolume,SetSourceMuted,SetSourceInputVolume,SetSourceInputMuted,MoveSourceOutput"`
	// Paylod send with Set* actions if necessary
	Payload interface{} `json:"payload,omitempty" doc:"Paylod send with Set* actions if necessary"`
}

type Response struct {
	// Action performed by API
	Action string `json:"action" doc:"Action performed by API"`
	// Status code
	Status int16 `json:"status" doc:"Status code"`
	// Response payload
	Payload interface{} `json:"payload" doc:"Response payload"`
	// Error description if any
	Error string `json:"error,omitempty" doc:"Error description if any"`
}

const (
	StatusSuccess          int16 = 4000
	StatusError            int16 = 4001
	StatusActionError      int16 = 4002
	StatusPayloadError     int16 = 4003
	StatusErrorInvalidJSON int16 = 4004
)

func (r Response) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"action": r.Action,
		"status": r.Status,
	}

	if r.Payload != nil {
		data["payload"] = r.Payload
	}

	if r.Error != "" {
		data["error"] = r.Error
	}

	return json.Marshal(data)
}
