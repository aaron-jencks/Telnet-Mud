package handlers

import (
	"mud/services/parsing"
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
		result.Specific[body[0][1:]] = strings.Join(body[1:], " ")
	} else {
		result.Others = strings.Join(body, " ")
	}

	return result
}

func HandleGlobal(body []string) parsing.CommandResponse {
	if len(body) == 0 {
		return parsing.CommandResponse{}
	}

	return parsing.CommandResponse{
		Global: strings.Join(body, " "),
	}
}
