package tmap

import (
	"math/rand"
	"mud/entities"
	"mud/utils"
	"mud/utils/io/db"
	"os"
	"path/filepath"
	"testing"
	"time"

	mtesting "mud/utils/testing"

	"github.com/stretchr/testify/assert"
)

func resetTable() {
	db.DeleteTable("map")
	entities.SetupRoomTable()
	entities.SetupTileTable()
	entities.SetupMapTable()
}

func TestMain(m *testing.M) {
	resetTable()
	rand.Seed(time.Now().Unix())

	err := m.Run()

	// Cleanup
	if db.DbDirectoryExists() {
		os.RemoveAll(filepath.Dir(utils.DB_LOCATION))
	}

	os.Exit(err)
}

func TestCreate(t *testing.T) {
	rid := rand.Int()
	tid := mtesting.GenerateRandomAlnumString(rand.Intn(64) + 1)
	x := rand.Int()
	y := rand.Int()
	z := rand.Int()

	args := []interface{}{
		rid,
		tid,
		x,
		y,
		z,
	}

	CRUD.Create(args...)
	is := CRUD.Retrieve(rid, x, y, z).(entities.Map)

	assert.Equal(t, rid, is.Room, "Created map should have the right room id")
	assert.Equal(t, tid, is.Tile, "Created map should have the right tile id")
	assert.Equal(t, x, is.X, "Created map should have the right x")
	assert.Equal(t, y, is.Y, "Created map should have the right y")
	assert.Equal(t, z, is.Z, "Created map should have the right z")
}

func TestCreateNoZ(t *testing.T) {
	rid := rand.Int()
	tid := mtesting.GenerateRandomAlnumString(rand.Intn(64) + 1)
	x := rand.Int()
	y := rand.Int()

	args := []interface{}{
		rid,
		tid,
		x,
		y,
	}

	CRUD.Create(args...)

	tiles := GetCurrentTilesForCoordWithType(rid, x, y, tid)

	assert.Equal(t, 1, len(tiles), "We only inserted one value")

	is := tiles[0]

	assert.Equal(t, rid, is.Room, "Created map should have the right room id")
	assert.Equal(t, tid, is.Tile, "Created map should have the right tile id")
	assert.Equal(t, x, is.X, "Created map should have the right x")
	assert.Equal(t, y, is.Y, "Created map should have the right y")
}

func createRandomTestMap() entities.Map {
	rid := rand.Int()
	tid := mtesting.GenerateRandomAlnumString(rand.Intn(64) + 1)
	x := rand.Int()
	y := rand.Int()
	z := rand.Int()

	args := []interface{}{
		rid,
		tid,
		x,
		y,
		z,
	}

	CRUD.Create(args...)
	return CRUD.Retrieve(rid, x, y, z).(entities.Map)
}

func TestUpdate(t *testing.T) {
	ps := createRandomTestMap()
	newTile := mtesting.GenerateRandomAlnumString(rand.Intn(64) + 1)
	ps.Tile = newTile

	nps := CRUD.Update(ps, ps.Room, ps.X, ps.Y, ps.Z).(entities.Map)

	assert.Equal(t, newTile, nps.Tile, "Tile should've been updated")
}

func TestDelete(t *testing.T) {
	ps := createRandomTestMap()

	CRUD.Retrieve(ps.Room, ps.X, ps.Y, ps.Z)
	CRUD.Delete(ps.Room, ps.X, ps.Y, ps.Z)
	psn := CRUD.Retrieve(ps.Room, ps.X, ps.Y, ps.Z)

	assert.Nil(t, psn, "Entry shouldn't exist after deleting")
}

func TestGetCoordTilesWithType(t *testing.T) {
	rid := rand.Int()
	tile := mtesting.GenerateRandomAsciiString(rand.Intn(64))
	ntile := mtesting.GenerateRandomAsciiString(rand.Intn(64))
	x := rand.Int()
	y := rand.Int()

	cargs := []interface{}{
		rid,
		tile,
		x,
		y,
	}

	CRUD.Create(cargs...)
	CRUD.Create(cargs...)
	CRUD.Create(cargs...)
	CRUD.Create([]interface{}{
		rid,
		ntile,
		x,
		y,
	}...)

	mTiles := GetCurrentTilesForCoordWithType(rid, x, y, tile)

	assert.Equal(t, 3, len(mTiles), "Should return the correct number of tiles")

	expectedz := 2
	for mti, mt := range mTiles {
		assert.Equal(t, rid, mt.Room, "Tiles should have the correct room")
		assert.Equal(t, tile, mt.Tile, "Tiles returned should have correct tile type")
		assert.Equal(t, x, mt.X, "Tiles should have the correct x coord")
		assert.Equal(t, y, mt.Y, "Tiles should have the correct y coord")
		assert.Equal(t, expectedz-mti, mt.Z, "Tiles should have correct z coord")
	}
}

func TestGetTilesForRoom(t *testing.T) {
	rid := rand.Int()

	tile := mtesting.GenerateRandomAsciiString(rand.Intn(64))
	x := rand.Int()
	y := rand.Int()

	CRUD.Create([]interface{}{
		rand.Int(),
		tile,
		x,
		y,
	}...)

	for i := 0; i < 10; i++ {
		tile = mtesting.GenerateRandomAsciiString(rand.Intn(64))
		x = rand.Int()
		y = rand.Int()

		CRUD.Create([]interface{}{
			rid,
			tile,
			x,
			y,
		}...)
	}

	rtiles := GetTilesForRoom(rid)

	assert.Equal(t, 10, len(rtiles), "Response should return all the tiles for a given room")

	for _, rt := range rtiles {
		assert.Equal(t, rid, rt.Room, "Tiles should have correct room id")
	}
}

func TestGetTilesForCoord(t *testing.T) {
	rid := rand.Int()

	x := rand.Int()
	y := rand.Int()

	for i := 0; i < 10; i++ {
		tile := mtesting.GenerateRandomAsciiString(rand.Intn(64))

		CRUD.Create([]interface{}{
			rid,
			tile,
			x,
			y,
		}...)
	}

	rtiles := GetCurrentTilesForCoord(rid, x, y)

	assert.Equal(t, 10, len(rtiles), "Response should return all the tiles for a given room")

	for rti, rt := range rtiles {
		assert.Equal(t, rid, rt.Room, "Tiles should have correct room id")
		assert.Equal(t, x, rt.X, "Tiles should have correct x coord")
		assert.Equal(t, y, rt.Y, "Tiles should have correct y coord")
		assert.Equal(t, 9-rti, rt.Z, "Tiles should have correct z coord")
	}
}

func TestTopmostTile(t *testing.T) {
	rid := rand.Int()
	x := rand.Int()
	y := rand.Int()
	var ltile string

	for i := 0; i < 10; i++ {
		ltile = mtesting.GenerateRandomAsciiString(rand.Intn(64))

		CRUD.Create([]interface{}{
			rid,
			ltile,
			x,
			y,
		}...)
	}

	lt := GetTopMostTile(rid, x, y)

	assert.Equal(t, ltile, lt.Tile, "Top tile should have correct tile")
	assert.Equal(t, rid, lt.Room, "Tiles should have correct room id")
	assert.Equal(t, x, lt.X, "Tiles should have correct x coord")
	assert.Equal(t, y, lt.Y, "Tiles should have correct y coord")
}

func TestRegion(t *testing.T) {
	rid := rand.Int()
	lx := rand.Intn(255)
	ty := rand.Intn(255)
	rx := rand.Intn(255) + lx
	by := rand.Intn(255) + ty
	tile := mtesting.GenerateRandomAsciiString(rand.Intn(64))

	xdiff := rx - lx
	ydiff := by - ty

	for i := 0; i < 10; i++ {
		x := rand.Intn(xdiff) + lx
		y := rand.Intn(ydiff) + ty

		CRUD.Create([]interface{}{
			rid,
			tile,
			x,
			y,
		}...)
	}

	nx := rand.Intn(255) + by
	ny := rand.Intn(255) + rx

	CRUD.Create([]interface{}{
		rid,
		tile,
		nx,
		ny,
	}...)

	nx = rand.Intn(lx)
	ny = rand.Intn(ty)

	CRUD.Create([]interface{}{
		rid,
		tile,
		nx,
		ny,
	}...)

	nx = rand.Intn(lx)
	ny = rand.Intn(ydiff) + ty

	CRUD.Create([]interface{}{
		rid,
		tile,
		nx,
		ny,
	}...)

	nx = rand.Intn(xdiff) + lx
	ny = rand.Intn(ty)

	CRUD.Create([]interface{}{
		rid,
		tile,
		nx,
		ny,
	}...)

	regtiles := GetTilesForRegion(rid, lx, ty, rx, by)

	assert.Equal(t, 10, len(regtiles), "Should return the correct number of tiles")
	for _, rt := range regtiles {
		assert.Equal(t, rid, rt.Room, "Tiles should be in the correct room")
		assert.True(t, rt.X <= rx && rt.X >= lx, "X should be within the region")
		assert.True(t, rt.Y <= by && rt.Y >= ty, "Y should be within the region")
	}
}

func TestSurrounding(t *testing.T) {
	rid := rand.Int()
	x := rand.Intn(255) + 1
	y := rand.Intn(255) + 1
	tile := mtesting.GenerateRandomAsciiString(rand.Intn(64))

	for tx := x - 1; tx <= x+1; tx++ {
		for ty := y - 1; ty <= y+1; ty++ {
			CRUD.Create([]interface{}{
				rid,
				tile,
				tx,
				ty,
			}...)
		}
	}

	regtiles := GetSurroundingTiles(rid, x, y)

	assert.Equal(t, 4, len(regtiles), "Should return the correct number of tiles")
	for _, rt := range regtiles {
		assert.Equal(t, rid, rt.Room, "Tiles should be in the correct room")
		if rt.X == x {
			assert.True(t, rt.Y != y && (rt.Y == y-1 || rt.Y == y+1),
				"The coords should be surrounding but not on")
		} else if rt.Y == y {
			assert.True(t, rt.X != x && (rt.X == x-1 || rt.X == x+1),
				"The coords should be surrounding but not on")
		} else {
			assert.True(t, rt.X == x-1 || rt.X == x+1, "X should be surrounding")
			assert.True(t, rt.Y == y-1 || rt.Y == y+1, "Y should be surrounding")
		}
	}
}
