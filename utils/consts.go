package utils

import "mud/utils/net"

var CONN_HOST string = net.GetOutboundIP().String()

const (
	SYSTEM_NAME                 = "SYSTEM"
	WINDOW_W                    = 80
	WINDOW_H                    = 20
	TERMINAL_H                  = 9
	TERMINAL_W                  = 40
	MAP_H                       = 10
	MAP_W                       = 40
	DEFAULT_MAP_BACKGROUND      = "\033[42m\u00b0\033[0m"
	CHAT_H                      = 19
	CHAT_W                      = 40
	CONN_PORT                   = "23"
	CONN_TYPE                   = "tcp"
	CHECK_DIE                   = 20
	LOOK_FAIL_MESSAGE           = "You don't notice anything in particular"
	CACHE_SIZE_LIMIT            = 1000
	DEFAULT_PLAYER_ACTION_LIMIT = 10
	DEFAULT_PLAYER_MODE         = "normal"
)
