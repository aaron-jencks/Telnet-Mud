package handlers

import (
	"fmt"
	acrud "mud/actions/defined/crud"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	"mud/services/loot"
	"mud/services/tmap"
	"mud/utils/handlers/crud"
	"net"
)

var LootCrudHandler parsing.CommandHandler = acrud.CreateCrudParser(
	"loot",
	"Usage: loot create room item quantity [x y [z]]",
	"Usage: loot retrieve id",
	"Usage: loot update id property:(room|item|quantity|x|y|z) newValue",
	"Usage: loot delete id",
	4, 2, 4, 2,
	func(c net.Conn, s []string) bool {
		usageString := "Usage: loot create room item quantity [x y [z]]"
		rparsable, _ := crud.ParseIntegerCheck(c, s[1], usageString, "room")
		iparsable, _ := crud.ParseIntegerCheck(c, s[2], usageString, "item")
		qparsable, _ := crud.ParseIntegerCheck(c, s[3], usageString, "quantity")
		firstPart := rparsable && iparsable && qparsable
		if firstPart && len(s) > 5 {
			xparsable, _ := crud.ParseIntegerCheck(c, s[4], usageString, "x")
			yparsable, _ := crud.ParseIntegerCheck(c, s[5], usageString, "y")

			secondPart := xparsable && yparsable

			if secondPart && len(s) == 7 {
				zparsable, _ := crud.ParseIntegerCheck(c, s[6], usageString, "z")

				return zparsable
			}

			return secondPart
		}
		return firstPart
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: loot retrieve id", "id")
		return parsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: loot update id property:(room|item|quantity|x|y|z) newValue", "id")
		return parsable
	},
	func(c net.Conn, s []string) bool {
		parsable, _ := crud.ParseIntegerCheck(c, s[1], "Usage: loot delete id", "id")
		return parsable
	},
	func(c net.Conn, s []string) []interface{} {
		var rid, iid, qty, x, y, z int
		fmt.Sscanf(s[0], "%d", &rid)
		fmt.Sscanf(s[1], "%d", &iid)
		fmt.Sscanf(s[2], "%d", &qty)

		if len(s) > 4 {
			fmt.Sscanf(s[3], "%d", &x)
			fmt.Sscanf(s[4], "%d", &y)

			if len(s) == 6 {
				fmt.Sscanf(s[5], "%d", &z)
			} else {
				z = tmap.GetTopMostTile(rid, x, y).Z + 1
			}
		} else {
			username := player.GetConnUsername(c)
			p := player.CRUD.Retrieve(username).(entities.Player)
			x = p.RoomX
			y = p.RoomY
			z = tmap.GetTopMostTile(rid, x, y).Z + 1
		}

		return []interface{}{rid, iid, qty, x, y, z}
	},
	func(c net.Conn, s []string) interface{} {
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
		return fmt.Sprintf("Loot %d deleted!", nv.Id)
	},
	func(c net.Conn) {},
	func(c net.Conn) {},
	func(c net.Conn) {},
	func(c net.Conn) {},
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
		case "x":
			c.X = newValue
		case "y":
			c.Y = newValue
		case "z":
			c.Z = newValue
		}

		return c
	},
	acrud.DefaultCrudModes, loot.CRUD,
)
