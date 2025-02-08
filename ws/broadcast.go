package ws

import (
	"log"
	"reflect"
	"time"

	"github.com/undg/go-prapi/json"
	"github.com/undg/go-prapi/pactl"
	"github.com/undg/go-prapi/utils"
)

var prevRes json.Response

const writeWait = 10 * time.Second

func BroadcastUpdates() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		clientsMutex.Lock()
		clientsCount := len(clients)
		clientsMutex.Unlock()

		if clientsCount == 0 {
			if utils.DEBUG {
				log.Println("No clients connected. Skipping VOLUME update.")
			}

			continue
		}

		// Same Action and StatusSuccess if everything is OK
		res := json.Response{
			Action: string(json.ActionGetStatus),
			Status: json.StatusSuccess,
		}

		res.Payload = pactl.GetStatus()

		equal := reflect.DeepEqual(res, prevRes)
		if equal {
			continue
		}

		prevRes = res

		clientsMutex.Lock()
		updatedClients := 0
		for conn := range clients {
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := safeWriteJson(conn, res)
			if err != nil {
				log.Printf("Error broadcast VOLUME update to client: %v\n", err)
				conn.Close()
				delete(clients, conn)
			} else {
				updatedClients++
			}
		}
		clientsMutex.Unlock()

		log.Printf("Volume broadcast sent to %d/%d clients. Value: %v\n", updatedClients, clientsCount, res.Payload)
	}
}
