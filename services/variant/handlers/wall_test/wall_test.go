package wall_test

import (
	"math/rand"
	"mud/entities"
	"mud/services/room"
	"mud/services/tile"
	"mud/services/tmap"
	"mud/services/variant"
	"mud/services/variant/handlers"
	"mud/utils"
	"mud/utils/io/db"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func resetTable() {
	db.DeleteTable("map")
	db.DeleteTable("rooms")
	db.DeleteTable("tiles")
	db.DeleteTable("variants")
	entities.SetupTileTable()
	entities.SetupRoomTable()
	entities.SetupMapTable()
	entities.SetupVariantTable()
}

var wallTile entities.Tile

func TestMain(m *testing.M) {
	resetTable()
	rand.Seed(time.Now().Unix())

	// Setup Test Data
	variant.CRUD.Create(1, "H", "\u2500")
	variant.CRUD.Create(1, "V", "\u2502")
	variant.CRUD.Create(1, "TLC", "\u250c")
	variant.CRUD.Create(1, "TRC", "\u2510")
	variant.CRUD.Create(1, "BLC", "\u2514")
	variant.CRUD.Create(1, "BRC", "\u2518")
	variant.CRUD.Create(1, "RT", "\u251c")
	variant.CRUD.Create(1, "LT", "\u2524")
	variant.CRUD.Create(1, "DT", "\u252c")
	variant.CRUD.Create(1, "UT", "\u2534")
	variant.CRUD.Create(1, "Cross", "\u253c")

	wallTile = tile.CRUD.Create("Wall", "wall", "1", true).(entities.Tile)
	tile.CRUD.Create("RBG", "normal", " ")

	room.CRUD.Create("Test Room", "It's a test room", 3, 3, "RBG")

	err := m.Run()

	// Cleanup
	if db.DbDirectoryExists() {
		os.RemoveAll(filepath.Dir(utils.DB_LOCATION))
	}

	os.Exit(err)
}

func TestWallHandlerTLC(t *testing.T) {
	tt := tmap.CRUD.Create(1, "Wall", 1, 1).(entities.Map)
	tmap.CRUD.Create(1, "Wall", 1, 2)
	tmap.CRUD.Create(1, "Wall", 2, 1)

	icon := handlers.ParseWallVariant(wallTile, 1, 1, 1, tt.Z)
	assert.Equal(t, "\u250c", icon, "TLC should parse correctly")
}
