package ws

import (
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/undg/go-prapi/buildinfo"
	"github.com/undg/go-prapi/json"
	"github.com/undg/go-prapi/pactl"
	"github.com/undg/go-prapi/utils"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)
var clientsMutex = &sync.Mutex{}

func upgraderCheckOrigin() {
	Upgrader.CheckOrigin = func(r *http.Request) bool {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Error splitting host and port: %v\n", err)
			return false
		}

		ip := net.ParseIP(host)
		if ip == nil {
			log.Printf("Invalid IP: %s\n", host)
			return false
		}

		return utils.IsLocalIP(ip) || strings.HasPrefix(r.Host, "localhost")
	}
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("wsEndpoint visited by: %s %s\n", r.Host, r.RemoteAddr)

	upgraderCheckOrigin()

	conn, err := Upgrader.Upgrade(w, r, nil)
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
	if err := conn.WriteJSON(initialResponse); err != nil {
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

		case json.ActionGetBuildInfo:
			b := buildinfo.Get()
			res.Payload = b

		case json.ActionSetSinkVolume:
			handleSetSinkVolume(&msg, &res)

		case json.ActionSetSinkMuted:
			handleSetSinkMuted(&msg, &res)

		case json.ActionSetSinkInputVolume:
			handleSetSinkInputVolume(&msg, &res)

		case json.ActionSetSinkInputMuted:
			handleSetSinkInputMuted(&msg, &res)

		case json.ActionMoveSinkInput:
			handleMoveSinkInput(&msg, &res)
		
		case json.ActionMoveSourceOutput:
			handleMoveSourceOutput(&msg, &res)

		case json.ActionGetCards:
			handleGetCards(&res)
			handleGetOutputs(&res)

		case json.ActionGetSinks:
			status, _ := pactl.GetOutputs()
			res.Payload = status

		case json.ActionGetSources:
			status, _ := pactl.GetSources()
			res.Payload = status

		case json.ActionGetSchema:
			handleGetSchema(&res)

		default:
			res.Error = "Command not found. Available actions: " + strings.Join(utils.ActionsToStrings(json.AvailableCommands), " ")
			res.Status = json.StatusActionError
		}

		handleServerLog(&msg, &res)

		if err := conn.WriteJSON(res); err != nil {
			log.Printf("Error writing JSON: %v\n", err)
			break
		}
	}
}
