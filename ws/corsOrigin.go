package ws

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/undg/go-prapi/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func upgraderCheckOrigin() {
	errPrefix := "ERROR [upgraderCheckOrigin()]: "
	upgrader.CheckOrigin = func(r *http.Request) bool {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("%s splitting host and port: %v\n", errPrefix, err)
			return false
		}

		ip := net.ParseIP(host)
		if ip == nil {
			log.Printf("%s Can't parse IP: %s\n", errPrefix, host)
			return false
		}

		if utils.IsLocalIP(ip) || strings.HasPrefix(r.Host, "localhost") {
			return true
		} else {
			log.Printf("%s IP is not allowed: %s\n", errPrefix, host)
			return false
		}

	}
}
