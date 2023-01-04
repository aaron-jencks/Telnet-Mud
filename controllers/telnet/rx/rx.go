package rx

import (
	"mud/actions/definitions"
	"mud/controllers"
	"mud/controllers/telnet/tx"
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	"mud/utils"
	"mud/utils/net/telnet"
	"mud/utils/strings"
	"mud/utils/ui/logger"
	"net"
	"time"
)

func ListenAndServe(handler func(net.Conn)) {
	// Listen for incoming connections.
	l, err := net.Listen(utils.CONN_TYPE, utils.CONN_HOST+":"+utils.CONN_PORT)
	if err != nil {
		logger.Error("Error listening: %v", err.Error())
		panic(err)
	}
	// Close the listener when the application closes.
	defer l.Close()
	logger.Info("Telnet listening on " + utils.CONN_HOST + ":" + utils.CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Error accepting: %v", err.Error())
			panic(err)
		}
		// Handle connections in a new goroutine.
		go handler(conn)
	}
}

// Handles incoming requests.
func TelnetHandler(conn net.Conn) {
	cid := tx.RegisterConnection(conn)
	defer tx.RemoveClient(cid)
	defer conn.Close()

	anonUsername := player.GenerateRandomUsername(conn)
	defer player.UnregisterAnonymousName(anonUsername)

	logger.Info("New Connection from %s. Assigned to %s", conn.RemoteAddr().String(), anonUsername)

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	player.PushAction(anonUsername, definitions.Action{
		Name:        "Initial Display",
		Duration:    0,
		AlwaysValid: true,
		Handler: func() parsing.CommandResponse {
			return parsing.CommandResponse{
				Clear:  true,
				Chat:   true,
				Map:    true,
				Info:   true,
				Conn:   conn,
				Person: true,
			}
		},
	})

	for {
		// Read the incoming connection into the buffer.
		reqLen, err := conn.Read(buf)
		logger.Info("Read in %d bytes with error: %v", reqLen, err)
		if err != nil {
			logger.Error("Error reading: %v", err.Error())
			return
		}

		if reqLen > 0 {
			var text []byte
			var headerResponse []byte
			line := buf[:reqLen]
			for li := 0; li < len(line); li++ {
				datum := line[li]
				var skipLength int = 0
				if telnet.IsIAC(line, li) {
					head, length := telnet.ParseIAC(conn, line, li)
					headerResponse = append(headerResponse, head...)
					li += length - 1
					continue
				} else {
					text = append(text, datum)
				}

				li += skipLength
			}

			if len(headerResponse) > 0 {
				tx.SendTarget(headerResponse, conn)
			}

			if len(text) > 0 && strings.IsNonEmpty(text) {
				controllers.HandlePacket(conn, text)
			}
		} else {
			logger.Info("Waiting for input...")
			time.Sleep(2 * time.Second)
		}
	}
}

// Serves and handles telnet protocol
func TelnetListenAndServe() {
	ListenAndServe(TelnetHandler)
}
