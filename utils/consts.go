package utils

import "mud/utils/net"

var CONN_HOST string = net.GetOutboundIP().String()

const (
	SYSTEM_NAME       = "SYSTEM"
	WINDOW_W          = 80
	WINDOW_H          = 20
	TERMINAL_H        = 19
	TERMINAL_W        = 40
	CHAT_H            = 19
	CHAT_W            = 40
	CONN_PORT         = "23"
	CONN_TYPE         = "tcp"
	CHECK_DIE         = 20
	LOOK_FAIL_MESSAGE = "You don't notice anything in particular"
)
