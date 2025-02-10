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
		log.Printf("Error upgrading to WebSocket: %v\n", err)
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

	if err := safeWriteJson(conn, initialResponse); err != nil {
		log.Printf("Error sending initial sinks data: %v\n", err)
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
				log.Printf("Error reading JSON: %v\n", err)
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
			print("set source")
		case json.ActionSetSourceMuted:
			print("mute source")

		// App's under SOURCES
		case json.ActionSetSourceInputVolume:
			print("set App under source volume")
		case json.ActionSetSourceInputMuted:
			print("mute App under source")

		case json.ActionMoveSourceOutput:
			handleMoveSourceOutput(&msg, &res)

		default:
			res.Error = "Command not found. Available actions: " + strings.Join(utils.ActionsToStrings(json.AvailableCommands), " ")
			res.Status = json.StatusActionError
		}

		handleServerLog(&msg, &res)

		if err := safeWriteJson(conn, res); err != nil {
			log.Printf("Error writing JSON: %v\n", err)
			break
		}
	}
}
