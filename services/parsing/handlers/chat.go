package handlers

import (
	chatService "mud/services/chat"
	"mud/services/parsing"
	"mud/utils/ui/pages/chat"
	"strings"
)

func HandleChat(body []string) parsing.CommandResponse {
	if len(body) == 0 {
		return parsing.CommandResponse{}
	}

	result := parsing.CommandResponse{}

	if body[0][0] == '@' {
		// Direct message
		result.Specific = make(parsing.DirectMessageMap)

		// TODO add way to get connections from usernames in the player service

		result.Specific[body[0][1:]] = chat.FormatChatEntry("Anonymous", strings.Join(body[1:], " "))
	} else {
		chatService.SendGlobalMessage("Anonymous", strings.Join(body, " "))
	}

	return result
}

func HandleGlobal(body []string) parsing.CommandResponse {
	if len(body) == 0 {
		return parsing.CommandResponse{}
	}

	return parsing.CommandResponse{
		Global: chat.FormatChatEntry("Anonymous", strings.Join(body, " ")),
	}
}
