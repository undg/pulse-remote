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

	ActionSetSinkVolume Action = "SetSinkVolume"
	ActionSetSinkMuted  Action = "SetSinkMuted"

	ActionSetSinkInputVolume Action = "SetSinkInputVolume"
	ActionSetSinkInputMuted  Action = "SetSinkInputMuted"

	ActionMoveSinkInput    Action = "MoveSinkInput"
	ActionMoveSourceOutput Action = "MoveSourceOutput"

	// Not used for now
	ActionGetSinks   Action = "GetSinks"
	ActionGetCards   Action = "GetCards"
	ActionGetOutputs Action = "GetOutputs"
	ActionGetSources Action = "GetSources"
	ActionGetSchema  Action = "GetSchema"
)

var AvailableCommands = []Action{
	ActionGetStatus,
	ActionGetBuildInfo,

	ActionSetSinkVolume,
	ActionSetSinkMuted,

	ActionSetSinkInputVolume,
	ActionSetSinkInputMuted,

	ActionMoveSinkInput,

	// Not used for now
	ActionGetSinks,
	ActionGetCards,
	ActionGetOutputs,
	ActionGetSources,
	ActionGetSchema,
}

// Message is an request from the client
type Message struct {
	// Actions listed in availableCommands slice
	Action Action `json:"action" doc:"Action to perform fe. GetVolume, SetVolume, SetMute..."`
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
