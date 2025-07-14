package ws

import (
	"reflect"
	"time"

	"github.com/undg/go-prapi/api/json"
	"github.com/undg/go-prapi/api/logger"
	"github.com/undg/go-prapi/api/pactl"
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
			logger.Debug().Msg("No clients connected. Skip VOLUME update.")
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

		loggerMsg := "broadcasting volume status"

		for conn := range clients {
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := safeWriteJSON(conn, res)
			if err != nil {
				logger.Error().Err(err).Msg(loggerMsg)
				conn.Close()
				delete(clients, conn)
			} else {
				updatedClients++
			}
		}
		clientsMutex.Unlock()

		if res.Error != "" {
			logger.Error().Str("Action", res.Action).Int("Status", int(res.Status)).Str("Error", string(res.Error)).Int("updated_clients", updatedClients).Msg(loggerMsg)
		}

		logger.Info().Str("Action", res.Action).Int("Status", int(res.Status)).Int("updated_clients", updatedClients).Msg(loggerMsg)
		logger.Debug().Str("res.Payload", "DEBUG=TRACE to see Payload").Msg(loggerMsg)
		logger.Trace().Interface("full_res", res).Msg(loggerMsg)
	}
}
