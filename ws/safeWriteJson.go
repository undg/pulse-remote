package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	pactlMutex sync.Mutex
	writeMutex sync.Mutex
)

func safeWriteJson(conn *websocket.Conn, v interface{}) error {
	writeMutex.Lock()
	defer writeMutex.Unlock()
	return conn.WriteJSON(v)
}
