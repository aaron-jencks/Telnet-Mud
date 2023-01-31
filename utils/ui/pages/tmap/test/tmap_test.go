package test

import (
	"mud/entities"
	"mud/parsing_services/player"
	"mud/services/room"
	"mud/services/tile"
	stmap "mud/services/tmap"
	"mud/utils"
	"mud/utils/io/db"
	mtesting "mud/utils/testing"
	"mud/utils/ui/logger"
	"mud/utils/ui/pages/tmap"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testTiles []entities.Tile = []entities.Tile{
	{
		Name:  "Stone Wall",
		Icon:  "#",
		FG:    90,
		BG:    0,
		Solid: true,
	},
}

func resetTables() {
	db.DeleteTable("tiles")
	db.DeleteTable("players")
	db.DeleteTable("rooms")
	db.DeleteTable("map")
	entities.SetupTileTable()
	entities.SetupRoomTable()
	entities.SetupPlayerTable()
	entities.SetupMapTable()
	for _, tt := range testTiles {
		tile.CRUD.Create(tt.Name, tt.IconType, tt.Icon, tt.BG, tt.FG)
	}
}

func TestMain(m *testing.M) {
	resetTables()

	err := m.Run()

	// Cleanup
	if db.DbDirectoryExists() {
		os.RemoveAll(filepath.Dir(utils.DB_LOCATION))
	}

	os.Exit(err)
}

type testTile struct {
	X    int
	Y    int
	Tile string
}

type mapWindowTestCaseResults struct {
	WindowLX int
	WindowTY int
	WindowRX int
	WindowBY int
}

type mapWindowTestCase struct {
	Name           string
	PlayerX        int
	PlayerY        int
	RoomH          int
	RoomW          int
	FillBackground bool
	FillTile       string
	Tiles          []testTile
	Results        mapWindowTestCaseResults
}

func TestGetMapWindow(t *testing.T) {
	testCases := []mapWindowTestCase{
		{
			Name:           "Outside of map window range",
			PlayerX:        0,
			PlayerY:        10,
			RoomH:          19,
			RoomW:          17,
			FillBackground: false,
			Tiles: []testTile{
				{
					X:    0,
					Y:    0,
					Tile: "Stone Wall",
				},
			},
			Results: mapWindowTestCaseResults{
				WindowLX: 0,
				WindowTY: 6,
				WindowRX: 16,
				WindowBY: 13,
			},
		},
		{
			Name:           "Outside of portrangeY",
			PlayerX:        0,
			PlayerY:        5,
			RoomH:          19,
			RoomW:          17,
			FillBackground: false,
			Tiles: []testTile{
				{
					X:    0,
					Y:    0,
					Tile: "Stone Wall",
				},
			},
			Results: mapWindowTestCaseResults{
				WindowLX: 0,
				WindowTY: 1,
				WindowRX: 16,
				WindowBY: 8,
			},
		},
		{
			Name:           "Outside of portrangeX with smallWidth",
			PlayerX:        10,
			PlayerY:        0,
			RoomH:          19,
			RoomW:          17,
			FillBackground: false,
			Tiles: []testTile{
				{
					X:    0,
					Y:    0,
					Tile: "Stone Wall",
				},
			},
			Results: mapWindowTestCaseResults{
				WindowLX: 0,
				WindowTY: 0,
				WindowRX: 16,
				WindowBY: 7,
			},
		},
		{
			Name:           "Outside of portrangeX",
			PlayerX:        30,
			PlayerY:        0,
			RoomH:          19,
			RoomW:          200,
			FillBackground: false,
			Tiles: []testTile{
				{
					X:    0,
					Y:    0,
					Tile: "Stone Wall",
				},
			},
			Results: mapWindowTestCaseResults{
				WindowLX: 11,
				WindowTY: 0,
				WindowRX: 48,
				WindowBY: 7,
			},
		},
		{
			Name:           "Near right side",
			PlayerX:        195,
			PlayerY:        0,
			RoomH:          19,
			RoomW:          200,
			FillBackground: false,
			Tiles: []testTile{
				{
					X:    0,
					Y:    0,
					Tile: "Stone Wall",
				},
			},
			Results: mapWindowTestCaseResults{
				WindowLX: 180,
				WindowTY: 0,
				WindowRX: 199,
				WindowBY: 7,
			},
		},
	}

	for _, tc := range testCases {
		logger.Info("Testing %s", tc.Name)

		// Sets up the test case
		resetTables()

		// Creates the player
		p := player.CRUD.Create(
			mtesting.GenerateRandomAlnumString(10),
			mtesting.GenerateRandomAlnumString(10),
		).(entities.Player)
		p.RoomX = tc.PlayerX
		p.RoomY = tc.PlayerY
		p = player.CRUD.Update(p, p.Id).(entities.Player)

		// Creates the test room
		r := room.CRUD.Create(
			mtesting.GenerateRandomAlnumString(10),
			mtesting.GenerateRandomAlnumString(10),
			tc.RoomH,
			tc.RoomW,
			testTiles[0].Name,
		).(room.ExpandedRoom)

		// Populates the room's tiles
		if tc.FillBackground {
			for row := 0; row < tc.RoomH; row++ {
				for col := 0; col < tc.RoomW; col++ {
					stmap.CRUD.Create(r.Id, tc.FillTile, col, row)
				}
			}
		}

		for _, rt := range tc.Tiles {
			stmap.CRUD.Create(r.Id, rt.Tile, rt.X, rt.Y)
		}

		// Performs the actual test
		tlx, tly, brx, bry := tmap.GetMapPortCoords(p.Room, p.RoomX, p.RoomY)

		assert.Equal(t, tc.Results.WindowLX, tlx,
			"%s: Map box Left X should be calculated correctly", tc.Name)
		assert.Equal(t, tc.Results.WindowTY, tly,
			"%s: Map box Top Y should be calculated correctly", tc.Name)
		assert.Equal(t, tc.Results.WindowRX, brx,
			"%s: Map box Right X should be calculated correctly", tc.Name)
		assert.Equal(t, tc.Results.WindowBY, bry,
			"%s: Map box Bottom Y should be calculated correctly", tc.Name)
	}

}
