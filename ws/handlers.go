package ws

import (
	j "encoding/json"
	"fmt"
	"log"

	"github.com/undg/go-prapi/audiodeprecated"
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

		pactl.SetSinkVolume(name, fmt.Sprintf("%.2f",volume))

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

		pactl.SetSinkInputVolume(fmt.Sprintf("%.0f", id), fmt.Sprintf("%.2f",volume))

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

func handleGetCards(res *json.Response) {
	errPrefix := "ERROR [handleGetCards()]"
	debugPrefix := "DEBUG [handleGetCards()]"

	cards, err := audiodeprecated.GetCards()
	if err != nil {
		log.Printf("%s audiodeprecated.GetCards(): %s\n", errPrefix, err)
		res.Error = "ERROR can't get cards information from the system"
		res.Status = json.StatusError
	}

	b, err := j.Marshal(cards)
	if err != nil {
		log.Printf("%s j.Marshal(cards): %s\n", errPrefix, err)
		res.Error = "ERROR can't pull cards information"
		res.Status = json.StatusError
	}

	res.Payload = string(b)
	if utils.DEBUG {
		log.Printf("%s handleGetCards(res).Action: %s\n", debugPrefix, res.Action)
		log.Printf("%s handleGetCards(res).Payload: %s\n", debugPrefix, res.Payload)
	}
}

func handleGetOutputs(res *json.Response) {
	errPrefix := "ERROR [handleGetOutputs()]"
	debugPrefix := "DEBUG [handleGetOutputs()]"

	outputs, err := audiodeprecated.GetOutputs()
	if err != nil {
		log.Printf("%s audiodeprecated.GetOutputs(): %s\n", errPrefix, err)

		res.Error = "ERROR can't get outputs information from the system"
		res.Status = json.StatusError
	}

	b, err := j.Marshal(outputs)
	if err != nil {
		log.Printf("%s j.Marshal(): %s\n", errPrefix, err)

		res.Error = "ERROR can't pull outputs information"
		res.Status = json.StatusError
	}

	res.Payload = string(b)

	if utils.DEBUG {
		log.Printf("%s res.Payload: %s\n", debugPrefix, res.Payload)
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
