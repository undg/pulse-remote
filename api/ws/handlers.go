package ws

import (
	"fmt"

	"github.com/undg/go-prapi/api/json"
	"github.com/undg/go-prapi/api/logger"
	"github.com/undg/go-prapi/api/pactl"
)

func handleSetSinkVolume(msg *json.Message, res *json.Response) {
	if sinkInfo, ok := msg.Payload.(map[string]interface{}); ok {
		name, ok := sinkInfo["name"].(string)
		if !ok {
			logger.Error().Msg("sinkInfo['name'].(string) NOT OK\n")
		}

		volume, ok := sinkInfo["volume"].(float64)
		if !ok {
			logger.Error().Msg("sinkInfo['volume'].(float64) NOT OK")
		}

		pactl.SetSinkVolume(name, fmt.Sprintf("%.2f", volume))

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleSetSinkMuted(msg *json.Message, res *json.Response) {
	if sinkInfo, ok := msg.Payload.(map[string]interface{}); ok {
		name, ok := sinkInfo["name"].(string)
		if !ok {
			logger.Error().Msg("sinkInfo['name'].(string) NOT OK")
		}

		muted, ok := sinkInfo["muted"].(bool)
		if !ok {
			logger.Error().Msg("sinkInfo['muted'].(bool) NOT OK")
		}

		pactl.SetSinkMuted(name, muted)

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleSetSinkInputVolume(msg *json.Message, res *json.Response) {
	if sinkInputInfo, ok := msg.Payload.(map[string]interface{}); ok {
		id, ok := sinkInputInfo["id"].(float64)
		if !ok {
			logger.Error().Msg("sinkInfo['id'].(float64) NOT OK")
		}

		volume, ok := sinkInputInfo["volume"].(float64)
		if !ok {
			logger.Error().Msg("sinkInfo['volume'].(float64) NOT OK")
		}

		pactl.SetSinkInputVolume(fmt.Sprintf("%.0f", id), fmt.Sprintf("%.2f", volume))

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleSetSinkInputMuted(msg *json.Message, res *json.Response) {
	if sinkInputInfo, ok := msg.Payload.(map[string]interface{}); ok {
		id, ok := sinkInputInfo["id"].(float64)
		if !ok {
			logger.Error().Msg("sinkInfo['id'].(float64) NOT OK")
		}

		muted, ok := sinkInputInfo["muted"].(bool)
		if !ok {
			logger.Error().Msg("sinkInfo['muted'].(bool) NOT OK")
		}

		pactl.SetSinkInputMuted(fmt.Sprintf("%.0f", id), muted)

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleMoveSinkInput(msg *json.Message, res *json.Response) {
	if sinkInputInfo, ok := msg.Payload.(map[string]interface{}); ok {
		sinkInputID, ok := sinkInputInfo["id"].(float64)
		if !ok {
			logger.Error().Msg("sinkInfo['id'].(float64) NOT OK")
		}

		sinkName, ok := sinkInputInfo["name"].(string)
		if !ok {
			logger.Error().Msg("sinkInfo['name'].(string) NOT OK")
		}

		pactl.MoveSinkInput(fmt.Sprintf("%.0f", sinkInputID), sinkName)

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

// SOURCES, Microphones
func handleSetSourceVolume(msg *json.Message, res *json.Response) {
	if sourceInfo, ok := msg.Payload.(map[string]interface{}); ok {
		name, ok := sourceInfo["name"].(string)
		if !ok {
			logger.Error().Msg("sourceInfo['name'].(string) NOT OK")
		}

		volume, ok := sourceInfo["volume"].(float64)
		if !ok {
			logger.Error().Msg("sourceInfo['volume'].(float64) NOT OK")
		}

		pactl.SetSourceVolume(name, fmt.Sprintf("%.2f", volume))

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid source information format"
		res.Status = json.StatusActionError
	}
}

func handleSetSourceMuted(msg *json.Message, res *json.Response) {
	if sourceInfo, ok := msg.Payload.(map[string]interface{}); ok {
		name, ok := sourceInfo["name"].(string)
		if !ok {
			logger.Error().Msg("sourceInfo['name'].(string) NOT OK")
		}

		muted, ok := sourceInfo["muted"].(bool)
		if !ok {
			logger.Error().Msg("sourceInfo['muted'].(bool) NOT OK")
		}

		pactl.SetSourceMuted(name, muted)

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid source information format"
		res.Status = json.StatusActionError
	}
}

func handleSetSourceInputVolume(msg *json.Message, res *json.Response) {
	if sourceInputInfo, ok := msg.Payload.(map[string]interface{}); ok {
		id, ok := sourceInputInfo["id"].(float64)
		if !ok {
			logger.Error().Msg("sourceInfo['id'].(float64) NOT OK")
		}

		volume, ok := sourceInputInfo["volume"].(float64)
		if !ok {
			logger.Error().Msg("sourceInfo['volume'].(float64) NOT OK")
		}

		pactl.SetSourceInputVolume(fmt.Sprintf("%.0f", id), fmt.Sprintf("%.2f", volume))

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid source information format"
		res.Status = json.StatusActionError
	}
}

func handleSetSourceInputMuted(msg *json.Message, res *json.Response) {
	if sourceInputInfo, ok := msg.Payload.(map[string]interface{}); ok {
		id, ok := sourceInputInfo["id"].(float64)
		if !ok {
			logger.Error().Msg("sourceInfo['id'].(float64) NOT OK")
		}

		muted, ok := sourceInputInfo["muted"].(bool)
		if !ok {
			logger.Error().Msg("sourceInfo['muted'].(bool) NOT OK")
		}

		pactl.SetSourceInputMuted(fmt.Sprintf("%.0f", id), muted)

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid source information format"
		res.Status = json.StatusActionError
	}
}

func handleMoveSourceOutput(msg *json.Message, res *json.Response) {
	if sourceOutputInfo, ok := msg.Payload.(map[string]interface{}); ok {
		sourceOutputID, ok := sourceOutputInfo["outputId"].(float64)
		if !ok {
			logger.Error().Msg("sourceOutputInfo['outputId'].(float64) NOT OK")
		}

		sourceName, ok := sourceOutputInfo["sourceName"].(string)
		if !ok {
			logger.Error().Msg("sourceOutputInfo['sourceName'].(string) NOT OK")
		}

		pactl.MoveSourceOutput(fmt.Sprintf("%.0f", sourceOutputID), sourceName)

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid source information format"
		res.Status = json.StatusActionError
	}
}

func handleServerLog(msg *json.Message, res *json.Response) {
	if msg != nil {
		logger.Trace().Str("Action", string(msg.Action)).Interface("Payload", msg.Payload).Msg("Incoming msg")
	}

	logger.Trace().Interface("res.Payload", res.Payload).Msg("Response to client")
	logger.Info().Int("res_code", int(res.Status)).Msg("server status")

}
