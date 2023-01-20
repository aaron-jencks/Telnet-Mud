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
