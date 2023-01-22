package player

import (
	"math/rand"
	"mud/utils"
	"net"
	"strings"
)

var namesCurrentlyInUse map[string]net.Conn = make(map[string]net.Conn)
var connAnonNameMap map[net.Conn]string = make(map[net.Conn]string)

func GenerateRandomUsername(conn net.Conn) string {
	// TODO check to make sure that we haven't used every
	// possible username

	var result []byte = make([]byte, utils.MIN_USERNAME_LENGTH)

	for {
		for ri := range result {
			result[ri] = byte(rand.Int()%94) + 33
		}

		sResult := "Anon." + strings.ReplaceAll(string(result), "\"", "\\\"")

		_, ok := namesCurrentlyInUse[sResult]
		if !ok && !PlayerExists(sResult) {
			namesCurrentlyInUse[sResult] = conn
			connAnonNameMap[conn] = sResult

			CreateAnonymousHandler(sResult)
			return sResult
		}
	}
}

func UnregisterAnonymousName(name string) {
	conn, ok := namesCurrentlyInUse[name]
	if ok {
		delete(namesCurrentlyInUse, name)
		delete(connAnonNameMap, conn)
		UnregisterHandler(name)
	}
}

func GetAnonymousUsername(conn net.Conn) string {
	name, ok := connAnonNameMap[conn]
	if ok {
		return name
	} else {
		return ""
	}
}
