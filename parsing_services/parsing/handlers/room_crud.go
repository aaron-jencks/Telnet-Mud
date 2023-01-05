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
	"Usage: room create \"name\" \"description\"",
	"Usage: room retrieve id",
	"Usage: room update id property:(name|description) \"newValue\"",
	"Usage: room delete id",
	3, 2, 4, 2,
	func(c net.Conn, s []string) bool { return true },
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: room retrieve id", "id")
		return parsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: room update id property:(name|description) \"newValue\"", "id")
		return parsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: room delete id", "id")
		return parsable
	},
	func(s []string) []interface{} {
		return []interface{}{strings.StripQuotes(s[0]), strings.StripQuotes(s[1])}
	},
	func(s []string) interface{} {
		var id int
		fmt.Sscanf(s[0], "%d", &id)
		return id
	},
	func(i interface{}) string {
		nv := i.(entities.Room)
		return fmt.Sprintf("Room %d(%s) created!", nv.Id, nv.Name)
	},
	func(i interface{}) string {
		r := i.(entities.Room)
		return fmt.Sprintf("Room %d:\nName: \"%s\"\nDescription: \"%s\"",
			r.Id, r.Name, r.Description)
	},
	func(i interface{}) string {
		nv := i.(entities.Room)
		return fmt.Sprintf("Room %d(%s) updated!", nv.Id, nv.Name)
	},
	func(i interface{}) string {
		nv := i.(entities.Room)
		return fmt.Sprintf("Room %d(%s) deleted!", nv.Id, nv.Name)
	},
	[]string{"name", "description"}, 2,
	func(i interface{}, s1 string, s2 []string) interface{} {
		c := i.(entities.Room)

		nv := strings.StripQuotes(s2[0])
		switch s1 {
		case "name":
			c.Name = nv
		case "description":
			c.Description = nv
		}

		return c
	},
	acrud.DefaultCrudModes, room.CRUD,
)

func HandleInfo(conn net.Conn, args []string) {
	username := player.GetConnUsername(conn)
	t := terminal.TerminalMap[conn]
	player.PushAction(username, defined.CreateInfoAction(conn, fmt.Sprintf("%s:\n%s", t.Room.Name, t.Room.Description)))
}
