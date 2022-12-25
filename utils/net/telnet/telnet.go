package telnet

import (
	"mud/services/parsing"
	"mud/utils"
	"mud/utils/strings"
	"mud/utils/ui/logger"
	"net"
	"sync"
	"time"
)

const (
	IAC  byte = 255
	SE   byte = 240
	NOP  byte = 241
	DM   byte = 242
	BRK  byte = 243
	IP   byte = 244
	AO   byte = 245
	AYT  byte = 246
	EC   byte = 247
	EL   byte = 248
	GA   byte = 249
	SB   byte = 250
	WILL byte = 251
	WONT byte = 252
	DO   byte = 253
	DONT byte = 254
)

var Clients []net.Conn
var ClientLock sync.Mutex = sync.Mutex{}

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

func SendGlobal(body []byte) {
	ClientLock.Lock()
	for _, client := range Clients {
		client.Write(body)
	}
	ClientLock.Unlock()
}

func SendOthers(body []byte, avoid net.Conn) {
	ClientLock.Lock()
	for _, client := range Clients {
		if client != avoid {
			client.Write(body)
		}
	}
	ClientLock.Unlock()
}

func SendTarget(body []byte, target net.Conn) {
	target.Write(body)
}

// Handles incoming requests.
func TelnetHandler(conn net.Conn) {
	ClientLock.Lock()
	Clients = append(Clients, conn)
	cid := len(Clients) - 1
	ClientLock.Unlock()
	defer removeClient(cid)
	defer conn.Close()

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	for {
		// Read the incoming connection into the buffer.
		reqLen, err := conn.Read(buf)
		logger.Info("Read in %d bytes with error: %v", reqLen, err)
		logger.Info(buf[:reqLen])
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
				if datum == IAC {
					switch line[li+1] {
					case WILL:
						var can byte = WONT
						if DoesOption(line[li+2]) {
							can = WILL
						}
						headerResponse = append(headerResponse, []byte{IAC, can, line[li+2]}...)
						skipLength = 2
					case DO:
						if DoesOption(line[li+2]) {
							doResp, argCount := DoOption(line[li+2], conn, line[li+3:])
							headerResponse = append(headerResponse, doResp...)
							skipLength = 2 + argCount
						} else {
							headerResponse = append(headerResponse, []byte{IAC, DONT, line[li+2]}...)
							skipLength = 2
						}
					default:
						text = append(text, datum)
					}
				} else {
					text = append(text, datum)
				}

				li += skipLength
			}

			if len(headerResponse) > 0 {
				SendTarget(headerResponse, conn)
			}

			if len(text) > 0 && strings.IsNonEmpty(text) {
				response := parsing.HandlePacket(conn, text)
				if response.Global {

				} else {
					if response.Person {

					}
					if len(response.Specific) > 0 {

					}
				}
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

func removeClient(conn int) {
	ClientLock.Lock()
	Clients = append(Clients[:conn], Clients[conn:]...)
	ClientLock.Unlock()
}
