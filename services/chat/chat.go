package chat

import (
	"mud/utils/ui"
	"mud/utils/ui/pages/chat"
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
		chat.FormatChatEntry(sender,
			ui.BoldText(receiver)+ui.StripIllegalChars(message)))
}

func SendDirectMessage(conn net.Conn, user, message string) {
	MessageLogMap[conn] = append(MessageLogMap[conn],
		chat.FormatChatEntry(user,
			ui.StripIllegalChars(message)))
}

func SendGlobalMessage(user, message string) {
	for userConn := range MessageLogMap {
		SendDirectMessage(userConn, user, message)
	}
}

func FetchMessageLog(conn net.Conn) []string {
	return MessageLogMap[conn]
}
