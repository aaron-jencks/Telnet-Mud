package handlers

import (
	"fmt"
	acrud "mud/actions/defined/crud"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/services/item"
	"mud/utils/handlers/crud"
	"mud/utils/strings"
	"net"
)

var ItemCrudHandler parsing.CommandHandler = acrud.CreateCrudParser(
	"item",
	"Usage: item create \"name\" \"description\"",
	"Usage: item retrieve id",
	"Usage: item update id property:(name|description) \"newValue\"",
	"Usage: item delete id",
	3, 2, 4, 2,
	func(c net.Conn, s []string) bool { return true },
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: item retrieve id", "id")
		return parsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: item update id property:(name|description) \"newValue\"", "id")
		return parsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: item delete id", "id")
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
		nv := i.(entities.Item)
		return fmt.Sprintf("Item %d(%s) created!", nv.Id, nv.Name)
	},
	func(i interface{}) string {
		r := i.(entities.Item)
		return fmt.Sprintf("Item %d:\nName: \"%s\"\nDescription: \"%s\"",
			r.Id, r.Name, r.Description)
	},
	func(i interface{}) string {
		nv := i.(entities.Item)
		return fmt.Sprintf("Item %d(%s) updated!", nv.Id, nv.Name)
	},
	func(i interface{}) string {
		nv := i.(entities.Item)
		return fmt.Sprintf("Item %d(%s) deleted!", nv.Id, nv.Name)
	},
	[]string{"name", "description"}, 2,
	func(i interface{}, s1 string, s2 []string) interface{} {
		c := i.(entities.Item)

		nv := strings.StripQuotes(s2[0])
		switch s1 {
		case "name":
			c.Name = nv
		case "description":
			c.Description = nv
		}

		return c
	},
	acrud.DefaultCrudModes, item.CRUD,
)
