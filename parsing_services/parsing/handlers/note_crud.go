package handlers

import (
	"fmt"
	"mud/entities"
	"mud/services/chat"
	"mud/services/note"
	"mud/services/parsing"
	"mud/services/player"
	"mud/utils/strings"
	"net"
)

func HandleNoteCrud(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Chat:   true,
		Person: true,
	}

	if CrudChecks(conn, "note", args) {
		return result
	}

	switch args[0] {
	case "create":
		if CheckMinArgs(conn, args, 3, "Usage: note create \"title\" \"contents\"") {
			return result
		}

		p := player.PlayerConnectionMap[conn]
		pe := player.CRUD.Retrieve(p).(entities.Player)

		nr := note.CRUD.Create(
			pe.Id,
			strings.StripQuotes(args[1]),
			strings.StripQuotes(args[2])).(entities.Note)
		chat.SendSystemMessage(conn, fmt.Sprintf("Note %d(%s) created!", nr.Id, nr.Title))

	case "retrieve":
		if CheckMinArgs(conn, args, 2, "Usage: note retrieve id") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: note retrieve id", "id")
		if !idParsed {
			return result
		}

		r := note.CRUD.Retrieve(id).(entities.Note)
		chat.SendSystemMessage(conn,
			fmt.Sprintf("Note %d:\nTitle: \"%s\"\nContents: \"%s\"",
				r.Id, r.Title, r.Contents))

	case "update":
		if CheckMinArgs(conn, args, 4, "Usage: note update id (title|contents) \"newValue\"") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: note update id (title|contents) \"newValue\"", "id")
		if !idParsed {
			return result
		}

		if CheckStringOptions(conn, args[2], []string{"title", "contents"},
			"Usage: note update id property \"newValue\"", "property") {
			return result
		}

		r := note.CRUD.Retrieve(id).(entities.Note)
		nv := strings.StripQuotes(args[3])
		switch args[2] {
		case "title":
			r.Title = nv
		case "contents":
			r.Contents = nv
		}

		nr := note.CRUD.Update(id, r).(entities.Note)
		chat.SendSystemMessage(conn, fmt.Sprintf("Note %d(%s) updated!", nr.Id, nr.Title))

	case "delete":
		if CheckMinArgs(conn, args, 2, "Usage: note delete id") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: note delete id", "id")
		if !idParsed {
			return result
		}

		note.CRUD.Delete(id)
		chat.SendSystemMessage(conn, fmt.Sprintf("Note %d deleted!", id))
	}

	return result
}
