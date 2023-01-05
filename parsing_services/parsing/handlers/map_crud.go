package handlers

import (
	"fmt"
	acrud "mud/actions/defined/crud"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/services/tmap"
	"mud/utils/handlers/crud"
	"mud/utils/strings"
	"net"
)

var MapCrudHandler parsing.CommandHandler = acrud.CreateCrudParserMultiRetrieve(
	"map",
	"Usage: map create room \"icon\" x y [z]",
	"Usage: map retrieve room x y [z]",
	"Usage: map delete room x y z",
	5, 4, 5,
	func(c net.Conn, s []string) bool {
		usageString := "Usage: map create room \"icon\" x y [z]"
		rparsable, _ := crud.ParseIntegerCheck(c, s[1], usageString, "room")
		xparsable, _ := crud.ParseIntegerCheck(c, s[3], usageString, "x")
		yparsable, _ := crud.ParseIntegerCheck(c, s[4], usageString, "y")
		zparsable := len(s) < 6
		if !zparsable {
			zparsable, _ = crud.ParseIntegerCheck(c, s[5], usageString, "z")
		}
		return rparsable && xparsable && yparsable && zparsable
	},
	func(c net.Conn, s []string) bool {
		usageString := "Usage: map retrieve room x y [z]"
		rparsable, _ := crud.ParseIntegerCheck(c, s[1], usageString, "room")
		xparsable, _ := crud.ParseIntegerCheck(c, s[2], usageString, "x")
		yparsable, _ := crud.ParseIntegerCheck(c, s[3], usageString, "y")
		zparsable := len(s) < 5
		if !zparsable {
			zparsable, _ = crud.ParseIntegerCheck(c, s[4], usageString, "z")
		}
		return rparsable && xparsable && yparsable && zparsable
	},
	func(c net.Conn, s []string) bool {
		usageString := "Usage: map delete room x y z"
		rparsable, _ := crud.ParseIntegerCheck(c, s[1], usageString, "room")
		xparsable, _ := crud.ParseIntegerCheck(c, s[2], usageString, "x")
		yparsable, _ := crud.ParseIntegerCheck(c, s[3], usageString, "y")
		zparsable, _ := crud.ParseIntegerCheck(c, s[4], usageString, "z")
		return rparsable && xparsable && yparsable && zparsable
	},
	func(s []string) []interface{} {
		var rid, x, y, z int
		fmt.Sscanf(s[0], "%d", &rid)
		fmt.Sscanf(s[3], "%d", &x)
		fmt.Sscanf(s[4], "%d", &y)
		if len(s) == 6 {
			fmt.Sscanf(s[5], "%d", &z)
			return []interface{}{rid, strings.StripQuotes(s[2]), x, y, z}
		}
		return []interface{}{rid, strings.StripQuotes(s[2]), x, y}
	},
	func(s []string) []interface{} {
		var rid, x, y, z int
		fmt.Sscanf(s[0], "%d", &rid)
		fmt.Sscanf(s[1], "%d", &x)
		fmt.Sscanf(s[2], "%d", &y)
		fmt.Sscanf(s[3], "%d", &z)
		return []interface{}{rid, "Room", x, "X", y, "Y", z, "Z"}
	},
	func(conn net.Conn, args []interface{}) interface{} {
		var rid, x, y, z int
		rid = args[0].(int)
		x = args[1].(int)
		y = args[2].(int)
		if len(args) == 4 {
			z = args[3].(int)
			return tmap.GetTileForCoord(rid, x, y, z)
		}
		return tmap.GetTopMostTile(rid, x, y)
	},
	func(i interface{}) string {
		nr := i.(entities.Map)
		return fmt.Sprintf("Tile %s placed at (Room: %d, X: %d, Y: %d, Z: %d) created!",
			nr.Tile, nr.Room, nr.X, nr.Y, nr.Z)
	},
	func(i interface{}) string {
		r := i.(entities.Map)
		return fmt.Sprintf("Map:\nRoom: %d\nCoord: (%d, %d, %d)\nTile: \"%s\"",
			r.Room, r.X, r.Y, r.Z, r.Tile)
	},
	func(i interface{}) string {
		nv := i.(entities.Map)
		return fmt.Sprintf("Map (%d, %d, %d, %d) deleted!", nv.Room, nv.X, nv.Y, nv.Z)
	},
	acrud.DefaultCrudModes, tmap.CRUD,
)
