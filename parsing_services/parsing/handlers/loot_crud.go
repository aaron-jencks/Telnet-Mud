package handlers

import (
	"fmt"
	acrud "mud/actions/defined/crud"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/services/loot"
	"mud/utils/handlers/crud"
	"net"
)

var LootCrudHandler parsing.CommandHandler = acrud.CreateCrudParser(
	"loot",
	"Usage: loot create room item quantity",
	"Usage: loot retrieve id",
	"Usage: loot update id property:(room|item|quantity) newValue",
	"Usage: loot delete id",
	4, 2, 4, 2,
	func(c net.Conn, s []string) bool {
		usageString := "Usage: loot create room item quantity"
		rparsable, _ := crud.ParseIntegerCheck(c, s[1], usageString, "room")
		iparsable, _ := crud.ParseIntegerCheck(c, s[2], usageString, "item")
		qparsable, _ := crud.ParseIntegerCheck(c, s[3], usageString, "quantity")
		return rparsable && iparsable && qparsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: loot retrieve id", "id")
		return parsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: loot update id property:(room|item|quantity) newValue", "id")
		return parsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: loot delete id", "id")
		return parsable
	},
	func(s []string) []interface{} {
		var rid, iid, qty int
		fmt.Sscanf(s[0], "%d", &rid)
		fmt.Sscanf(s[1], "%d", &iid)
		fmt.Sscanf(s[2], "%d", &qty)
		return []interface{}{rid, iid, qty}
	},
	func(s []string) interface{} {
		var id int
		fmt.Sscanf(s[0], "%d", &id)
		return id
	},
	func(i interface{}) string {
		nv := i.(entities.Loot)
		return fmt.Sprintf("Loot %d created!", nv.Id)
	},
	func(i interface{}) string {
		r := i.(entities.Loot)
		return fmt.Sprintf("Loot %d:\nRoom: %d\nItem: %d\nQuantity: %d",
			r.Id, r.Room, r.Item, r.Quantity)
	},
	func(i interface{}) string {
		nv := i.(entities.Loot)
		return fmt.Sprintf("Loot %d updated!", nv.Id)
	},
	func(i interface{}) string {
		nv := i.(entities.Loot)
		return fmt.Sprintf("Loot %d updated!", nv.Id)
	},
	[]string{"room", "item", "quantity"}, 2,
	func(i interface{}, s1 string, s2 []string) interface{} {
		c := i.(entities.Loot)

		var newValue int
		fmt.Sscanf(s2[0], "%d", &newValue)

		switch s1 {
		case "room":
			c.Room = newValue
		case "item":
			c.Item = newValue
		case "quantity":
			c.Quantity = newValue
		}

		return c
	},
	acrud.DefaultCrudModes, loot.CRUD,
)
