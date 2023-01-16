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
	X        int
	Y        int
	Z        int
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

func SetupPlayerTable() {
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
		"Id integer primary key autoincrement",
		"Name text unique",
		"Password text",
		"Dex integer not null",
		"Str integer not null",
		"Int integer not null",
		"Wis integer not null",
		"Con integer not null",
		"Chr integer not null",
		"Room integer references rooms (Id) on delete no action on update no action",
		"RoomX integer not null",
		"RoomY integer not null",
		"ActionCapacity integer not null",
		"CurrentMode text not null",
	}, true)
}

func SetupTables() {
	// Map
	logger.Info("Creating map table")
	db.CreateTableIfNotExist("map", []string{
		"Room",
		"Tile",
		"X",
		"Y",
		"Z",
	}, []string{
		"Room integer references rooms (Id) on delete cascade on update cascade",
		"Tile text not null",
		"X integer not null",
		"Y integer not null",
		"Z integer not null",
	}, false)
	db.RunExec("alter table map add constraint PK_MAPLOC primary key (Room, X, Y, Z);")

	// Variant
	logger.Info("Creating tile variants table")
	db.CreateTableIfNotExist("variants", []string{
		"Id",
		"Name",
		"Icon",
	}, []string{
		"Id integer primary key",
		"Name text not null",
		"Icon text not null",
	}, false)
	db.RunExec("alter table variants add constraint PK_VARID primary key (Id, Name);")

	// Tile
	logger.Info("Creating tiles table")
	db.CreateTableIfNotExist("tiles", []string{
		"Name",
		"IconType",
		"Icon",
		"BG",
		"FG",
	}, []string{
		"Name text primary key",
		"IconType text not null",
		"Icon text not null",
		"BG integer not null",
		"FG integer not null",
	}, false)

	// Loot
	logger.Info("Creating loot table")
	db.CreateTableIfNotExist("loot", []string{
		"Id",
		"Room",
		"Item",
		"Quantity",
		"X",
		"Y",
		"Z",
	}, []string{
		"Id integer primary key autoincrement",
		"Room integer references rooms (Id) on delete cascade on update cascade",
		"Item integer references items (Id) on delete cascade on update cascade",
		"Quantity integer not null",
		"X integer not null",
		"Y integer not null",
		"Z integer not null",
	}, true)

	// Note
	logger.Info("Creating notes table")
	db.CreateTableIfNotExist("notes", []string{
		"Id",
		"Player",
		"Title",
		"Contents",
	}, []string{
		"Id integer primary key autoincrement",
		"Player integer references players (Id) on delete cascade on update cascade",
		"Title text not null",
		"Contents text not null",
	}, true)

	// Command
	logger.Info("Creating commands table")
	db.CreateTableIfNotExist("commands", []string{
		"Name",
		"ArgCount",
		"ArgRegex",
	}, []string{
		"Name text primary key",
		"ArgCount integer not null",
		"ArgRegex text not null",
	}, false)

	// Item
	logger.Info("Creating items table")
	db.CreateTableIfNotExist("items", []string{
		"Id",
		"Name",
		"Description",
	}, []string{
		"Id integer primary key autoincrement",
		"Name text unique",
		"Description text not null",
	}, true)

	SetupPlayerTable()

	// Inventoriy
	logger.Info("Creating inventory table")
	db.CreateTableIfNotExist("inventory", []string{
		"Id",
		"Player",
		"Item",
		"Quantity",
	}, []string{
		"Id integer primary key autoincrement",
		"Player integer references players (Id) on delete cascade on update cascade",
		"Item integer references items (Id) on delete cascade on update cascade",
		"Quantity integer not null",
	}, true)

	// Room
	logger.Info("Creating rooms table")
	db.CreateTableIfNotExist("rooms", []string{
		"Id",
		"Name",
		"Description",
		"Height",
		"Width",
	}, []string{
		"Id integer primary key autoincrement",
		"Name text not null",
		"Description text not null",
		"Height integer not null",
		"Width integer not null",
	}, true)
}
