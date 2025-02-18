package ws

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/undg/go-prapi/api/json"
	"github.com/undg/go-prapi/api/logger"
	"github.com/undg/go-prapi/api/pactl"
	"github.com/undg/go-prapi/api/utils"
)

var clients = make(map[*websocket.Conn]bool)
var clientsMutex = &sync.Mutex{}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	logger.Info().Str("server_ip", r.Host).Str("client_ip", r.RemoteAddr).Msg("New client attempting to connect")

	upgraderCheckOrigin()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Upgrading Websocket unsuccessful")
		return
	}

	clientsMutex.Lock()
	clients[conn] = true
	clientCount := len(clients)
	clientsMutex.Unlock()

	logger.Info().Int("clients_connected", clientCount).Msg("Client connection established")

	// Execute ActionGetStatus when a new client connects
	status := pactl.GetStatus()

	initialResponse := json.Response{
		Action:  string(json.ActionGetStatus),
		Status:  json.StatusSuccess,
		Payload: status,
	}

	if err := safeWriteJSON(conn, initialResponse); err != nil {
		logger.Error().Err(err).Msg("Initial sinks data FAIL")
	}

	// Cleanup after client is disconnected
	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		clientCounts := len(clients)
		clientsMutex.Unlock()
		conn.Close()
		logger.Info().Int("clients_count", clientCounts).Msg("Client disconnected")
	}()

	// Messaging system with client
	for {
		var msg json.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error().Err(err).Msg("Can't read JSON")
			}
			break
		}

		// Same Action and StatusSuccess if everyting is OK
		res := json.Response{
			Action: string(msg.Action),
			Status: json.StatusSuccess,
		}

		switch msg.Action {

		case json.ActionGetStatus:
			status := pactl.GetStatus()
			res.Payload = status

		// SINKS, Speakers
		case json.ActionSetSinkVolume:
			handleSetSinkVolume(&msg, &res)
		case json.ActionSetSinkMuted:
			handleSetSinkMuted(&msg, &res)

		// App's under SiNKS
		case json.ActionSetSinkInputVolume:
			handleSetSinkInputVolume(&msg, &res)
		case json.ActionSetSinkInputMuted:
			handleSetSinkInputMuted(&msg, &res)
		case json.ActionMoveSinkInput:
			handleMoveSinkInput(&msg, &res)

		// SOURCES, Microphones
		case json.ActionSetSourceVolume:
			handleSetSourceVolume(&msg, &res)
		case json.ActionSetSourceMuted:
			handleSetSourceMuted(&msg, &res)

		// App's under SOURCES
		case json.ActionSetSourceInputVolume:
			handleSetSourceInputVolume(&msg, &res)
		case json.ActionSetSourceInputMuted:
			handleSetSourceInputMuted(&msg, &res)

		case json.ActionMoveSourceOutput:
			handleMoveSourceOutput(&msg, &res)

		default:
			res.Error = "Command not found. Available actions: " + strings.Join(utils.ActionsToStrings(json.AvailableCommands), " ")
			res.Status = json.StatusActionError
		}

		handleServerLog(&msg, &res)

		if err := safeWriteJSON(conn, res); err != nil {
			logger.Error().Err(err).Msg("Can't write JSON")
			break
		}
	}
}
