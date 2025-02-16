package ws

import (
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/undg/go-prapi/logger"
	"github.com/undg/go-prapi/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func upgraderCheckOrigin() {
	msgPrefix := "ERROR [upgraderCheckOrigin()]: "
	upgrader.CheckOrigin = func(r *http.Request) bool {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logger.Error().Err(err).Msg("SplitHostPort")
			return false
		}

		ip := net.ParseIP(host)
		if ip == nil {
			logger.Error().Err(err).Msg("ParseIP")
			return false
		}

		if utils.IsLocalIP(ip) || strings.HasPrefix(r.Host, "localhost") {
			logger.Info().Str("IP", string(ip)).Str("Host", r.Host).Msgf("ParseIP %s", msgPrefix)
			return true
		} else {
			logger.Error().Str("IP", string(ip)).Str("Host", r.Host).Msgf("ParseIP %s", msgPrefix)
			return false
		}

	}
}
