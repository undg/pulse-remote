package utils

import (
	"net"

	"github.com/undg/go-prapi/json"
)

const PORT = ":8448"

func ActionsToStrings(actions []json.Action) []string {
	strs := make([]string, len(actions))
	for i, action := range actions {
		strs[i] = string(action)
	}
	return strs
}

func IsLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	if ip4 := ip.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return true
		case ip4[0] == 127 && ip4[1] == 0 && ip4[2] == 0 && ip4[3] == 1:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		}
	}

	return false
}

func GetLocalIP() (string, error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}
