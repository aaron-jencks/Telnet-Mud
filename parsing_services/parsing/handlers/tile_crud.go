package handlers

import (
	"fmt"
	acrud "mud/actions/defined/crud"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/services/tile"
	"mud/utils/handlers/crud"
	"mud/utils/strings"
	"net"
)

var TileCrudHandler parsing.CommandHandler = acrud.CreateCrudParser(
	"tile",
	"Usage: tile create \"name\" \"type\" \"icon\" [bg fg]",
	"Usage: tile retrieve \"name\"",
	"Usage: tile update \"name\" (name|type|icon|bg|fg) \"newValue\"",
	"Usage: tile delete \"name\"",
	4, 2, 4, 2,
	func(c net.Conn, s []string) bool {
		usageString := "Usage: tile create \"name\" \"type\" \"icon\" [bg fg]"
		if len(s) == 6 {
			bgParsed, _ := crud.ParseIntegerCheck(c, s[4], usageString, "bg")
			fgParsed, _ := crud.ParseIntegerCheck(c, s[5], usageString, "fg")

			return bgParsed && fgParsed
		}
		return true
	},
	func(c net.Conn, s []string) bool { return true },
	func(c net.Conn, s []string) bool {
		if s[2] == "bg" || s[2] == "fg" {
			nv := strings.StripQuotes(s[3])
			parsable, _ := crud.ParseIntegerCheck(c, nv, "Usage: tile update \"name\" (name|type|icon|bg|fg) \"newValue\"", "newValue")
			return parsable
		}
		return true
	},
	func(c net.Conn, s []string) bool { return true },
	func(s []string) []interface{} {
		name := strings.StripQuotes(s[0])
		itype := strings.StripQuotes(s[1])
		icon := parsing.ParseIconString(strings.StripQuotes(s[2]))
		if len(s) == 5 {
			var fg, bg int
			fmt.Sscanf(s[3], "%d", &bg)
			fmt.Sscanf(s[4], "%d", &fg)
			return []interface{}{name, itype, icon, bg, fg}
		}
		return []interface{}{name, itype, icon}
	},
	func(s []string) interface{} {
		return strings.StripQuotes(s[0])
	},
	func(i interface{}) string {
		nv := i.(entities.Tile)
		return fmt.Sprintf("Tile %s created!", nv.Name)
	},
	func(i interface{}) string {
		r := i.(entities.Tile)
		return fmt.Sprintf("Tile:\nName: \"%s\"\nType: \"%s\"\nIcon: \"\033[%dm\033%dm%s\033[0m\"",
			r.Name, r.IconType, r.BG, r.FG, r.Icon)
	},
	func(i interface{}) string {
		nv := i.(entities.Tile)
		return fmt.Sprintf("Tile %s updated!", nv.Name)
	},
	func(i interface{}) string {
		nv := i.(entities.Tile)
		return fmt.Sprintf("Tile %s deleted!", nv.Name)
	},
	[]string{"name", "type", "icon", "bg", "fg"}, 2,
	func(i interface{}, s1 string, s2 []string) interface{} {
		r := i.(entities.Tile)

		var newValue int
		fmt.Sscanf(s2[0], "%d", &newValue)

		nv := strings.StripQuotes(s2[0])
		switch s1 {
		case "name":
			r.Name = nv
		case "type":
			r.IconType = nv
		case "icon":
			r.Icon = parsing.ParseIconString(nv)
		case "bg":
			fmt.Sscanf(nv, "%d", &newValue)

			r.BG = newValue
		case "fg":
			fmt.Sscanf(nv, "%d", &newValue)

			r.FG = newValue
		}

		return r
	},
	acrud.DefaultCrudModes, tile.CRUD,
)
