package handlers

import (
	"fmt"
	"mud/actions/defined"
	acrud "mud/actions/defined/crud"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	"mud/services/room"
	"mud/services/terminal"
	"mud/utils/handlers/crud"
	"mud/utils/strings"
	"net"
)

var RoomCrudHandler parsing.CommandHandler = acrud.CreateCrudParser(
	"room",
	"Usage: room create \"name\" \"description\" height width",
	"Usage: room retrieve id",
	"Usage: room update id property:(name|description|height|width) \"newValue\"",
	"Usage: room delete id",
	5, 2, 4, 2,
	func(c net.Conn, s []string) bool {
		usageString := "Usage: room create \"name\" \"description\" height width"

		hParsable, _ := crud.ParseIntegerCheck(c, s[3], usageString, "height")
		wParsable, _ := crud.ParseIntegerCheck(c, s[4], usageString, "width")

		return hParsable && wParsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: room retrieve id", "id")
		return parsable
	},
	func(c net.Conn, s []string) bool {
		usageString := "Usage: room update id property:(name|description) \"newValue\""
		parsable, _ := crud.ParseIntegerCheck(c, s[1], usageString, "id")
		if s[2] == "height" || s[2] == "width" {
			dimParsable, _ := crud.ParseIntegerCheck(c, strings.StripQuotes(s[3]), usageString, "newValue")
			return parsable && dimParsable
		}
		return parsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: room delete id", "id")
		return parsable
	},
	func(c net.Conn, s []string) []interface{} {
		var height, width int
		fmt.Sscanf(s[2], "%d", &height)
		fmt.Sscanf(s[3], "%d", &width)
		return []interface{}{strings.StripQuotes(s[0]), strings.StripQuotes(s[1]), height, width}
	},
	func(c net.Conn, s []string) interface{} {
		var id int
		fmt.Sscanf(s[0], "%d", &id)
		return id
	},
	func(i interface{}) string {
		nv := room.CRUD.Retrieve(i.(int)).(entities.Room)
		return fmt.Sprintf("Room %d(%s) created!", nv.Id, nv.Name)
	},
	func(i interface{}) string {
		if i != nil {
			r := i.(entities.Room)
			return fmt.Sprintf("Room %d:\nName: \"%s\"\nDescription: \"%s\"",
				r.Id, r.Name, r.Description)
		} else {
			return "That room did not exist!"
		}
	},
	func(i interface{}) string {
		if i != nil {
			nv := i.(entities.Room)
			return fmt.Sprintf("Room %d(%s) updated!", nv.Id, nv.Name)
		} else {
			return "That room did not exist!"
		}
	},
	func(i interface{}) string {
		if i != nil {
			nv := i.(entities.Room)
			return fmt.Sprintf("Room %d(%s) deleted!", nv.Id, nv.Name)
		} else {
			return "That room did not exist!"
		}
	},
	func(c net.Conn) {},
	func(c net.Conn) {},
	func(c net.Conn) {
		username := player.GetConnUsername(c)
		player.EnqueueAction(username, defined.CreateGlobalMapRepaint(c))
	},
	func(c net.Conn) {},
	[]string{"name", "description", "height", "width"}, 2,
	func(i interface{}, s1 string, s2 []string) interface{} {
		c := i.(entities.Room)

		var inv int

		nv := strings.StripQuotes(s2[0])
		switch s1 {
		case "name":
			c.Name = nv
		case "description":
			c.Description = nv
		case "height":
			fmt.Sscanf(nv, "%d", &inv)
			c.Height = inv
		case "width":
			fmt.Sscanf(nv, "%d", &inv)
			c.Width = inv
		}

		return c
	},
	acrud.DefaultCrudModes, room.CRUD,
)

func HandleInfo(conn net.Conn, args []string) {
	username := player.GetConnUsername(conn)
	t := terminal.TerminalMap[conn]
	player.EnqueueAction(username, defined.CreateInfoAction(conn, fmt.Sprintf("%s:\n%s", t.Room.Name, t.Room.Description)))
}
