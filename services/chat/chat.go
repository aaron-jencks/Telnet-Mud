package chat

import (
	"mud/utils"
	"mud/utils/ui"
	"mud/utils/ui/gui"
	"net"
)

var MessageLogMap map[net.Conn][]string = make(map[net.Conn][]string)

func RegisterConnection(conn net.Conn) {
	MessageLogMap[conn] = nil
}

func UnregisterConnection(conn net.Conn) {
	delete(MessageLogMap, conn)
}

func SendMentionMessage(conn net.Conn, sender, receiver, message string) {
	MessageLogMap[conn] = append(MessageLogMap[conn],
		gui.FormatChatEntry(sender,
			ui.BoldText(receiver)+" "+ui.StripIllegalChars(message)))
}

func SendDirectMessage(conn net.Conn, sender, message string) {
	MessageLogMap[conn] = append(MessageLogMap[conn],
		gui.FormatChatEntry(sender,
			ui.StripIllegalChars(message)))
}

func SendSystemMessage(conn net.Conn, message string) {
	SendDirectMessage(conn, utils.SYSTEM_NAME, message)
}

func SendGlobalMessage(sender, message string) {
	for userConn := range MessageLogMap {
		SendDirectMessage(userConn, sender, message)
	}
}

func FetchMessageLog(conn net.Conn) []string {
	return MessageLogMap[conn]
}
