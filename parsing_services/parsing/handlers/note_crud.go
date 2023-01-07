package handlers

import (
	"fmt"
	acrud "mud/actions/defined/crud"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/services/note"
	"mud/utils/handlers/crud"
	"mud/utils/strings"
	"net"
)

var NoteCrudHandler parsing.CommandHandler = acrud.CreateCrudParser(
	"note",
	"Usage: note create \"name\" \"description\"",
	"Usage: note retrieve id",
	"Usage: note update id property:(name|description) \"newValue\"",
	"Usage: note delete id",
	3, 2, 4, 2,
	func(c net.Conn, s []string) bool { return true },
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: note retrieve id", "id")
		return parsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: note update id property:(name|description) \"newValue\"", "id")
		return parsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: note delete id", "id")
		return parsable
	},
	func(c net.Conn, s []string) []interface{} {
		return []interface{}{strings.StripQuotes(s[0]), strings.StripQuotes(s[1])}
	},
	func(c net.Conn, s []string) interface{} {
		var id int
		fmt.Sscanf(s[0], "%d", &id)
		return id
	},
	func(i interface{}) string {
		nv := i.(entities.Note)
		return fmt.Sprintf("Note %d(%s) created!", nv.Id, nv.Title)
	},
	func(i interface{}) string {
		r := i.(entities.Note)
		return fmt.Sprintf("Note %d:\nName: \"%s\"\nDescription: \"%s\"",
			r.Id, r.Title, r.Contents)
	},
	func(i interface{}) string {
		nv := i.(entities.Note)
		return fmt.Sprintf("Note %d(%s) updated!", nv.Id, nv.Title)
	},
	func(i interface{}) string {
		nv := i.(entities.Note)
		return fmt.Sprintf("Note %d(%s) deleted!", nv.Id, nv.Title)
	},
	func(c net.Conn) {},
	func(c net.Conn) {},
	func(c net.Conn) {},
	func(c net.Conn) {},
	[]string{"title", "contents"}, 2,
	func(i interface{}, s1 string, s2 []string) interface{} {
		c := i.(entities.Note)

		nv := strings.StripQuotes(s2[0])
		switch s1 {
		case "title":
			c.Title = nv
		case "contents":
			c.Contents = nv
		}

		return c
	},
	acrud.DefaultCrudModes, note.CRUD,
)
