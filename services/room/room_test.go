package room

import (
	"math/rand"
	"mud/entities"
	"mud/services/tile"
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
	db.DeleteTable("rooms")
	db.DeleteTable("tiles")
	entities.SetupTileTable()
	entities.SetupRoomTable()
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

func createRandomTestTile() entities.Tile {
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	picontype := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	picon := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)

	args := []interface{}{
		pname,
		picontype,
		picon,
	}

	tile.CRUD.Create(args...)
	return tile.CRUD.Retrieve(pname).(entities.Tile)
}

func TestCreate(t *testing.T) {
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	pdescription := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	width := rand.Int()
	height := rand.Int()
	tt := createRandomTestTile()

	args := []interface{}{
		pname,
		pdescription,
		height,
		width,
		tt.Name,
	}

	inid := CRUD.Create(args...)

	is := inid.(ExpandedRoom)

	assert.Equal(t, pname, is.Name, "Created room should have the right name")
	assert.Equal(t, pdescription, is.Description, "Created room should have the right description")
	assert.Equal(t, width, is.Width, "Created room should have the right width")
	assert.Equal(t, height, is.Height, "Created room should have the right height")
	assert.Equal(t, tt, is.BackgroundTile, "Created room should have the right tile")
}

func createRandomTestRoom() ExpandedRoom {
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	pdescription := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	width := rand.Int()
	height := rand.Int()
	tt := createRandomTestTile()

	args := []interface{}{
		pname,
		pdescription,
		height,
		width,
		tt.Name,
	}

	return CRUD.Create(args...).(ExpandedRoom)
}

func TestUpdate(t *testing.T) {
	ps := createRandomTestRoom()
	newDescription := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	ps.Description = newDescription

	nps := CRUD.Update(ps, ps.Id).(ExpandedRoom)

	assert.Equal(t, newDescription, nps.Description, "Description should've been updated")
}

func TestDelete(t *testing.T) {
	ps := createRandomTestRoom()
	pid := ps.Id

	CRUD.Retrieve(int(pid))
	CRUD.Delete(int(pid))
	psn := CRUD.Retrieve(int(pid))

	assert.Nil(t, psn, "Entry shouldn't exist after deleting")
}
