package handlers

import (
	"fmt"
	acrud "mud/actions/defined/crud"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/services/variant"
	"mud/utils/handlers/crud"
	"mud/utils/strings"
	"net"
)

var VariantCrudHandler parsing.CommandHandler = acrud.CreateCrudParserMultiRetrieve(
	"variant",
	"Usage: variant create [id] \"name\" \"icon\"",
	"Usage: variant retrieve id \"name\"",
	"Usage: variant delete id \"name\"",
	3, 3, 3,
	func(c net.Conn, s []string) bool {
		usageString := "Usage: variant create [id] \"name\" \"icon\""
		if len(s) == 3 {
			parsable, _ := crud.ParseIntegerCheck(c, s[1], usageString, "id")

			return parsable
		}
		return true
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: variant retrieve id \"name\"", "id")

		return parsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: variant retrieve id \"name\"", "id")

		return parsable
	},
	func(c net.Conn, s []string) []interface{} {
		if len(s) == 3 {
			var id int
			fmt.Sscanf(s[0], "%d", &id)
			return []interface{}{id, strings.StripQuotes(s[1]), parsing.ParseIconString(strings.StripQuotes(s[2]))}
		}
		return []interface{}{strings.StripQuotes(s[0]), parsing.ParseIconString(strings.StripQuotes(s[1]))}
	},
	func(c net.Conn, s []string) []interface{} {
		var id int
		fmt.Sscanf(s[0], "%d", &id)
		return []interface{}{id, strings.StripQuotes(s[1])}
	},
	func(c net.Conn, i []interface{}) interface{} {
		return variant.CRUD.Retrieve(i[0].(int), i[1].(string))
	},
	func(i interface{}) string {
		if i == nil {
			return "That Variant does not exist"
		}
		nr := i.(entities.TileVariant)
		return fmt.Sprintf("Variant %d(%s) created!", nr.Id, nr.Name)
	},
	func(i interface{}) string {
		if i == nil {
			return "That Variant does not exist"
		}
		r := i.(entities.TileVariant)
		return fmt.Sprintf("Variant:\nId: %d\nName: \"%s\"\nIcon: \"%s\"",
			r.Id, r.Name, r.Icon)
	},
	func(i interface{}) string {
		if i == nil {
			return "That Variant does not exist"
		}
		nr := i.(entities.TileVariant)
		return fmt.Sprintf("Variant %d(%s) deleted!", nr.Id, nr.Name)
	},
	func(c net.Conn) {},
	func(c net.Conn) {},
	func(c net.Conn) {},
	acrud.DefaultCrudModes, variant.CRUD,
)
