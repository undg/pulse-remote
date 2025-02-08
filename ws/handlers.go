package ws

import (
	j "encoding/json"
	"fmt"
	"log"

	"github.com/undg/go-prapi/json"
	"github.com/undg/go-prapi/pactl"
	"github.com/undg/go-prapi/utils"
)

func handleSetSinkVolume(msg *json.Message, res *json.Response) {
	errPrefix := "ERROR [handleSetSinkVolume()]"

	if sinkInfo, ok := msg.Payload.(map[string]interface{}); ok {
		name, ok := sinkInfo["name"].(string)
		if !ok {
			log.Printf("%s sinkInfo['name'].(string) NOT OK\n", errPrefix)
		}

		volume, ok := sinkInfo["volume"].(float64)
		if !ok {
			log.Printf("%s sinkInfo['volume'].(float64) NOT OK\n", errPrefix)
		}

		pactl.SetSinkVolume(name, fmt.Sprintf("%.2f", volume))

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleSetSinkMuted(msg *json.Message, res *json.Response) {
	errPrefix := "ERROR [handleSetSinkMuted()]:"

	if sinkInfo, ok := msg.Payload.(map[string]interface{}); ok {
		name, ok := sinkInfo["name"].(string)
		if !ok {
			log.Printf("%s sinkInfo['name'].(string) NOT OK\n", errPrefix)
		}

		muted, ok := sinkInfo["muted"].(bool)
		if !ok {
			log.Printf("%s sinkInfo['muted'].(bool) NOT OK\n", errPrefix)
		}

		pactl.SetSinkMuted(name, muted)

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleSetSinkInputVolume(msg *json.Message, res *json.Response) {
	errPrefix := "ERROR [handleSetSinkInputVolume()]:"

	if sinkInputInfo, ok := msg.Payload.(map[string]interface{}); ok {
		id, ok := sinkInputInfo["id"].(float64)
		if !ok {
			log.Printf("%s sinkInfo['id'].(float64) NOT OK\n", errPrefix)
		}

		volume, ok := sinkInputInfo["volume"].(float64)
		if !ok {
			log.Printf("%s sinkInfo['volume'].(float64) NOT OK\n", errPrefix)
		}

		pactl.SetSinkInputVolume(fmt.Sprintf("%.0f", id), fmt.Sprintf("%.2f", volume))

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleSetSinkInputMuted(msg *json.Message, res *json.Response) {
	errPrefix := "ERROR [handleSetSinkInputMuted()]:"

	if sinkInputInfo, ok := msg.Payload.(map[string]interface{}); ok {
		id, ok := sinkInputInfo["id"].(float64)
		if !ok {
			log.Printf("%s sinkInfo['id'].(float64) NOT OK\n", errPrefix)
		}

		muted, ok := sinkInputInfo["muted"].(bool)
		if !ok {
			log.Printf("%s sinkInfo['muted'].(bool) NOT OK\n", errPrefix)
		}

		pactl.SetSinkInputMuted(fmt.Sprintf("%.0f", id), muted)

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleMoveSinkInput(msg *json.Message, res *json.Response) {
	errPrefix := "Error [handleMoveSinkInput()]"

	if sinkInputInfo, ok := msg.Payload.(map[string]interface{}); ok {
		sinkInputID, ok := sinkInputInfo["inputId"].(float64)
		if !ok {
			log.Printf("%s sinkInfo['inputId'].(float64) NOT OK\n", errPrefix)
		}

		sinkName, ok := sinkInputInfo["sinkName"].(string)
		if !ok {
			log.Printf("%s sinkInfo['sinkName'].(string) NOT OK\n", errPrefix)
		}

		pactl.MoveSinkInput(fmt.Sprintf("%.0f", sinkInputID), sinkName)

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid sink information format"
		res.Status = json.StatusActionError
	}
}

func handleMoveSourceOutput(msg *json.Message, res *json.Response) {
	errPrefix := "Error [handleMoveSourceOutput()]"

	if sourceOutputInfo, ok := msg.Payload.(map[string]interface{}); ok {
		sourceOutputID, ok := sourceOutputInfo["outputId"].(float64)
		if !ok {
			log.Printf("%s sourceOutputInfo['outputId'].(float64) NOT OK\n", errPrefix)
		}

		sourceName, ok := sourceOutputInfo["sourceName"].(string)
		if !ok {
			log.Printf("%s sourceOutputInfo['sourceName'].(string) NOT OK\n", errPrefix)
		}

		pactl.MoveSourceOutput(fmt.Sprintf("%.0f", sourceOutputID), sourceName)

		res.Payload = pactl.GetStatus()
	} else {
		res.Error = "Invalid source information format"
		res.Status = json.StatusActionError
	}
}

func handleGetSchema(res *json.Response) {
	debugPrefix := "DEBUG [handleGetSchema()]"
	// schema := json.GetSchemaJSON()
	//
	// res.Payload = schema
	if utils.DEBUG {
		log.Printf("%s res.Action: %s\n", debugPrefix, res.Action)
		log.Printf("%s res.Payload: %s\n", debugPrefix, res.Payload)
	}
}

func handleServerLog(msg *json.Message, res *json.Response) {
	errPrefix := "ERROR [handleServerLog()]"

	fmt.Printf("\n")
	log.Printf("\n-->\n")

	if msg != nil {
		msgBytes, err := j.MarshalIndent(msg, "", "	")
		if err != nil {
			fmt.Printf("%s j.MarshalIndent(): %s\n", errPrefix, err)
		}
		fmt.Printf("CLIENT message: %s\n", string(msgBytes))
	}

	if utils.DEBUG {
		resBytes, err := j.MarshalIndent(res, "", "	")
		if err != nil {
			fmt.Printf("%s serverLog res.MarshalJson %s\n", errPrefix, err)
		}

		fmt.Printf("SERVER res: %s\n", string(resBytes))
	} else {
		fmt.Printf("SERVER res.status: %d\n", res.Status)
	}

	fmt.Printf(">--\n\n")
}
