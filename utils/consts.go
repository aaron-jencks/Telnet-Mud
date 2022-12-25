package utils

import "mud/utils/net"

var CONN_HOST string = net.GetOutboundIP().String()

const (
	SYSTEM_NAME = "SYSTEM"
	CHAT_H      = 20
	CHAT_W      = 80
	CONN_PORT   = "23"
	CONN_TYPE   = "tcp"
)
