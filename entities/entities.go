package entities

import (
	"mud/utils/io/db"
	"mud/utils/ui/logger"
)

type TileVariant struct {
	Id   int
	Name string
	Icon string
}

// type VariantType struct {
// 	Id          int
// 	HandlerType string
// }

type Map struct {
	Room int
	Tile string
	X    int
	Y    int
	Z    int
}

type Tile struct {
	Name     string
	IconType string
	Icon     string
	BG       int
	FG       int
}

type Loot struct {
	Id       int
	Room     int
	Item     int
	Quantity int
}

type Note struct {
	Id       int
	Player   int
	Title    string
	Contents string
}

type Command struct {
	Name     string
	ArgCount int
	ArgRegex string
}

type Item struct {
	Id          int
	Name        string
	Description string
}

type Player struct {
	Id             int
	Name           string
	Password       string
	Dex            int
	Str            int
	Int            int
	Wis            int
	Con            int
	Chr            int
	Room           int
	RoomX          int
	RoomY          int
	ActionCapacity int
	CurrentMode    string
}

type Inventory struct {
	Id       int
	Player   int
	Item     int
	Quantity int
}

type Room struct {
	Id          int
	Name        string
	Description string
	Height      int
	Width       int
}

func SetupTables() {
	var table db.TableDefinition
	var index map[string][]int64

	// Map
	logger.Info("Creating map table")
	table = db.CreateTableIfNotExist("map", []string{
		"Room",
		"Tile",
		"X",
		"Y",
		"Z",
	}, []string{
		"integer",
		"string",
		"integer",
		"integer",
		"integer",
	}, 0, false)

	// // Variant Type
	// logger.Info("Creating variant types table")
	// table = db.CreateTableIfNotExist("variantTypes", []string{
	// 	"Id",
	// 	"HandlerType",
	// }, []string{
	// 	"integer",
	// 	"string",
	// }, 0, true)

	// Variant
	logger.Info("Creating tile variants table")
	table = db.CreateTableIfNotExist("variants", []string{
		"Id",
		"Name",
		"Icon",
	}, []string{
		"integer",
		"string",
		"string",
	}, 0, false)

	// Tile
	logger.Info("Creating tiles table")
	table = db.CreateTableIfNotExist("tiles", []string{
		"Name",
		"IconType",
		"Icon",
		"BG",
		"FG",
	}, []string{
		"string",
		"string",
		"string",
		"integer",
		"integer",
	}, 0, true)

	// Loot
	logger.Info("Creating loot table")
	table = db.CreateTableIfNotExist("loot", []string{
		"Id",
		"Room",
		"Item",
		"Quantity",
	}, []string{
		"integer",
		"integer",
		"integer",
		"integer",
	}, 0, true)
	index = db.CreateIndex(table.CSV, "Room")
	table.Info.Indices["Room"] = index

	// Note
	logger.Info("Creating notes table")
	table = db.CreateTableIfNotExist("notes", []string{
		"Id",
		"Player",
		"Title",
		"Contents",
	}, []string{
		"integer",
		"integer",
		"string",
		"string",
	}, 0, true)
	index = db.CreateIndex(table.CSV, "Player")
	table.Info.Indices["Player"] = index

	// Command
	logger.Info("Creating commands table")
	db.CreateTableIfNotExist("commands", []string{
		"Name",
		"ArgCount",
		"ArgRegex",
	}, []string{
		"string",
		"integer",
		"string",
	}, 0, true)

	// Item
	logger.Info("Creating items table")
	db.CreateTableIfNotExist("items", []string{
		"Id",
		"Name",
		"Description",
	}, []string{
		"integer",
		"string",
		"string",
	}, 0, true)

	// Player
	logger.Info("Creating players table")
	db.CreateTableIfNotExist("players", []string{
		"Id",
		"Name",
		"Password",
		"Dex",
		"Str",
		"Int",
		"Wis",
		"Con",
		"Chr",
		"Room",
		"RoomX",
		"RoomY",
		"ActionCapacity",
		"CurrentMode",
	}, []string{
		"integer",
		"string",
		"string",
		"integer",
		"integer",
		"integer",
		"integer",
		"integer",
		"integer",
		"integer",
		"integer",
		"integer",
		"integer",
		"string",
	}, 1, true)

	// Inventoriy
	logger.Info("Creating inventory table")
	table = db.CreateTableIfNotExist("inventory", []string{
		"Id",
		"Player",
		"Item",
		"Quantity",
	}, []string{
		"integer",
		"integer",
		"integer",
		"integer",
	}, 0, false)
	index = db.CreateIndex(table.CSV, "Item")
	table.Info.Indices["Item"] = index

	// Room
	logger.Info("Creating rooms table")
	db.CreateTableIfNotExist("rooms", []string{
		"Id",
		"Name",
		"Description",
		"Height",
		"Width",
	}, []string{
		"integer",
		"string",
		"string",
		"integer",
		"integer",
	}, 0, true)
}
