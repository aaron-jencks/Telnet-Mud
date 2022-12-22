package telnet

import (
	"logger"
	"net"
	"sync"
	"telnet_options"
	"time"
)

const (
	CONN_HOST string = "192.168.0.222"
	CONN_PORT string = "23"
	CONN_TYPE string = "tcp"
	IAC       byte   = 255
	SE        byte   = 240
	NOP       byte   = 241
	DM        byte   = 242
	BRK       byte   = 243
	IP        byte   = 244
	AO        byte   = 245
	AYT       byte   = 246
	EC        byte   = 247
	EL        byte   = 248
	GA        byte   = 249
	SB        byte   = 250
	WILL      byte   = 251
	WONT      byte   = 252
	DO        byte   = 253
	DONT      byte   = 254
)

var Clients []net.Conn
var ClientLock sync.Mutex = sync.Mutex{}

func ListenAndServe(handler func(net.Conn)) {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		logger.Error("Error listening: %v", err.Error())
		panic(err)
	}
	// Close the listener when the application closes.
	defer l.Close()
	logger.Info("Telnet listening on " + CONN_HOST + ":" + CONN_PORT)
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
	ClientLock.Lock()
	Clients = append(Clients, conn)
	cid := len(Clients) - 1
	ClientLock.Unlock()
	defer removeClient(cid)
	defer conn.Close()

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	conn.Write([]byte("Hello Again!\n\r\033[3J\033[;mHello World!\n\r"))

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
			var response []byte
			var transmit []byte
			line := buf[:reqLen]
			for li := 0; li < len(line); li++ {
				datum := line[li]
				var skipLength int = 0
				if datum == IAC {
					switch line[li+1] {
					case WILL:
						var can byte = WONT
						if telnet_options.DoesOption(line[li+2]) {
							can = WILL
						}
						response = append(response, []byte{IAC, can, line[li+2]}...)
						skipLength = 2
					case DO:
						if telnet_options.DoesOption(line[li+2]) {
							doResp, argCount := telnet_options.DoOption(line[li+2], conn, line[li+3:])
							response = append(response, doResp...)
							skipLength = 2 + argCount
						} else {
							response = append(response, []byte{IAC, DONT, line[li+2]}...)
							skipLength = 2
						}
					default:
						response = append(response, datum)
					}
				} else {
					transmit = append(transmit, datum)
				}

				li += skipLength
			}

			logger.Info("Sending %v", response)
			conn.Write(response)
			ClientLock.Lock()
			for _, client := range Clients {
				if client != conn {
					client.Write(transmit)
				}
			}
			ClientLock.Unlock()
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
