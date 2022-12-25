package net

import (
	"mud/utils/ui/logger"
	"net"
)

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "192.168.0.1:80") // Ping the router
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
