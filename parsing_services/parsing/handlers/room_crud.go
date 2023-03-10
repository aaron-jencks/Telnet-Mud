package handlers

import (
	"fmt"
	"mud/actions/defined"
	acrud "mud/actions/defined/crud"
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
	"Usage: room create \"name\" \"description\" height width \"background tile\"",
	"Usage: room retrieve id",
	"Usage: room update id property:(name|description|height|width|background) \"newValue\"",
	"Usage: room delete id",
	6, 2, 4, 2,
	func(c net.Conn, s []string) bool {
		usageString := "Usage: room create \"name\" \"description\" height width \"background tile\""

		hParsable, _ := crud.ParseIntegerCheck(c, s[3], usageString, "height")
		wParsable, _ := crud.ParseIntegerCheck(c, s[4], usageString, "width")

		return hParsable && wParsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: room retrieve id", "id")
		return parsable
	},
	func(c net.Conn, s []string) bool {
		usageString := "Usage: room update id property:(name|description|height|width|background) \"newValue\""
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
		return []interface{}{strings.StripQuotes(s[0]), strings.StripQuotes(s[1]),
			height, width, strings.StripQuotes(s[4])}
	},
	func(c net.Conn, s []string) interface{} {
		var id int
		fmt.Sscanf(s[0], "%d", &id)
		return id
	},
	func(i interface{}) string {
		nv := i.(room.ExpandedRoom)
		return fmt.Sprintf("Room %d(%s) created!", nv.Id, nv.Name)
	},
	func(i interface{}) string {
		if i != nil {
			r := i.(room.ExpandedRoom)
			return fmt.Sprintf("Room %d:\nName: \"%s\"\nDescription: \"%s\"\nSize: (%d, %d)\nBackground: \"\033[%dm\033[%dm%s\033[0m\"",
				r.Id, r.Name, r.Description, r.Width, r.Height,
				r.BackgroundTile.BG, r.BackgroundTile.FG, r.BackgroundTile.Icon)
		} else {
			return "That room did not exist!"
		}
	},
	func(i interface{}) string {
		if i != nil {
			nv := i.(room.ExpandedRoom)
			return fmt.Sprintf("Room %d(%s) updated!", nv.Id, nv.Name)
		} else {
			return "That room did not exist!"
		}
	},
	func(i interface{}) string {
		if i != nil {
			nv := i.(room.ExpandedRoom)
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
	[]string{"name", "description", "height", "width", "background"}, 2,
	func(i interface{}, s1 string, s2 []string) interface{} {
		c := i.(room.ExpandedRoom)

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
		case "background":
			c.BackgroundTile.Name = nv
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
