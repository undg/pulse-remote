package ws

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/undg/go-prapi/json"
	"github.com/undg/go-prapi/pactl"
	"github.com/undg/go-prapi/utils"
)

var clients = make(map[*websocket.Conn]bool)
var clientsMutex = &sync.Mutex{}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("wsEndpoint visited by: %s %s\n", r.Host, r.RemoteAddr)

	upgraderCheckOrigin()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR upgrading to WebSocket: %v\n", err)
		return
	}

	clientsMutex.Lock()
	clients[conn] = true
	clientCount := len(clients)
	clientsMutex.Unlock()

	log.Printf("New client connected. Total clients: %d\n", clientCount)

	// Execute ActionGetStatus when a new client connects
	status := pactl.GetStatus()

	initialResponse := json.Response{
		Action:  string(json.ActionGetStatus),
		Status:  json.StatusSuccess,
		Payload: status,
	}

	if err := safeWriteJSON(conn, initialResponse); err != nil {
		log.Printf("ERROR sending initial sinks data: %v\n", err)
	}

	// Cleanup after client is disconnected
	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		clientCount := len(clients)
		clientsMutex.Unlock()
		conn.Close()
		log.Printf("Client disconnected. Total clients: %d\n", clientCount)
	}()

	// Messaging system with client
	for {
		var msg json.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ERROR reading JSON: %v\n", err)
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
			log.Printf("ERROR writing JSON: %v\n", err)
			break
		}
	}
}
