package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const writeWait = 10 * time.Second

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}

		msg := string(p)
		msgOut := string(p)

		if msg == "get" {
			msgOut = string(marshalAudio(getVol()))
		}

		if msg == "set" {
			msgOut = string(marshalAudio(setVol(0.5)))
		}

		if msg == "muteToggle" {
			msgOut = string(marshalAudio(toggleMute()))
		}

		if msg == "mute" {
			msgOut = string(marshalAudio(mute()))
		}

		if msg == "unMute" {
			msgOut = string(marshalAudio(unMute()))
		}

		if err := conn.WriteMessage(messageType, []byte(msgOut)); err != nil {
			log.Println(err)
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("wsEndpoint visited by:", w.Header())

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	reader(ws)
}
