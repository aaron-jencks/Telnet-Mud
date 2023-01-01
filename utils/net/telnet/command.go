package telnet

import "net"

func IsIAC(data []byte, start int) bool {
	return len(data[start:]) >= 2 && data[start] == IAC
}

func ParseIAC(conn net.Conn, data []byte, start int) ([]byte, int) {
	var headerResponse []byte
	var skipLength int = 0

	if IsIAC(data, start) {
		skipLength = 3

		switch data[start+1] {
		case WILL:
			var can byte = WONT
			if DoesOption(data[start+2]) {
				can = WILL
			}
			headerResponse = append(headerResponse, []byte{IAC, can, data[start+2]}...)
		case DO:
			if DoesOption(data[start+2]) {
				doResp, argCount := DoOption(data[start+2], conn, data[start+3:])
				headerResponse = append(headerResponse, doResp...)
				skipLength += argCount
			} else {
				headerResponse = append(headerResponse, []byte{IAC, DONT, data[start+2]}...)
			}
		}
	}

	return headerResponse, skipLength
}
