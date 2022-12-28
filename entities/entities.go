package entities

import (
	"mud/utils/io/db"
	"mud/utils/ui/logger"
)

type Tile struct {
	Name     string
	IconType string
	Icon     string
}

type Detail struct {
	Id         int
	Room       int
	Direction  string
	Detail     string
	Perception int
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
	Id       int
	Name     string
	Password string
	Dex      int
	Str      int
	Int      int
	Wis      int
	Con      int
	Chr      int
	Room     int
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
}

type Transition struct {
	Id          int
	Source      int
	Target      int
	Command     string
	CommandArgs string
}

func SetupTables() {
	var table db.TableDefinition
	var index map[string][]int64

	// Tile
	logger.Info("Creating tiles table")
	table = db.CreateTableIfNotExist("tiles", []string{
		"Name",
		"IconType",
		"Icon",
	}, []string{
		"string",
		"string",
		"string",
	}, 0, true)

	// Detail
	logger.Info("Creating details table")
	table = db.CreateTableIfNotExist("details", []string{
		"Id",
		"Room",
		"Direction",
		"Detail",
		"Perception",
	}, []string{
		"integer",
		"integer",
		"string",
		"string",
		"integer",
	}, 0, true)
	index = db.CreateIndex(table.CSV, "Room")
	table.Info.Indices["Room"] = index

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
	}, []string{
		"integer",
		"string",
		"string",
	}, 0, true)

	// Transition
	table = db.CreateTableIfNotExist("transitions", []string{
		"Id",
		"Source",
		"Target",
		"Command",
		"CommandArgs",
	}, []string{
		"integer",
		"integer",
		"integer",
		"string",
		"string",
	}, 0, true)
	index = db.CreateIndex(table.CSV, "Target")
	table.Info.Indices["Target"] = index
	index = db.CreateIndex(table.CSV, "Source")
	table.Info.Indices["Source"] = index
}
