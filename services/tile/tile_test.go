package tile

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
	db.DeleteTable("tiles")
	entities.SetupTileTable()
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
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	picontype := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	picon := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	bg := rand.Int()
	fg := rand.Int()
	solid := (rand.Int() % 2) == 1

	args := []interface{}{
		pname,
		picontype,
		picon,
		bg,
		fg,
		solid,
	}

	CRUD.Create(args...)
	is := CRUD.Retrieve(pname).(entities.Tile)

	assert.Equal(t, pname, is.Name, "Created tile should have the right name")
	assert.Equal(t, picontype, is.IconType, "Created tile should have the right icon type")
	assert.Equal(t, picon, is.Icon, "Created tile should have the right icon")
	assert.Equal(t, fg, is.FG, "Created tile should have the right foreground color")
	assert.Equal(t, bg, is.BG, "Created tile should have the right background color")
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

	CRUD.Create(args...)
	return CRUD.Retrieve(pname).(entities.Tile)
}

func TestUpdate(t *testing.T) {
	ps := createRandomTestTile()
	newIconType := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	ps.IconType = newIconType

	nps := CRUD.Update(ps, ps.Name).(entities.Tile)

	assert.Equal(t, newIconType, nps.IconType, "Icon type should've been updated")
}

func TestUpdateSolid(t *testing.T) {
	ps := createRandomTestTile()
	ps.Solid = !ps.Solid

	nps := CRUD.Update(ps, ps.Name).(entities.Tile)

	assert.Equal(t, ps.Solid, nps.Solid, "Solid should've been updated")
}

func TestDelete(t *testing.T) {
	ps := createRandomTestTile()
	pid := ps.Name

	CRUD.Retrieve(pid)
	CRUD.Delete(pid)
	psn := CRUD.Retrieve(pid)

	assert.Nil(t, psn, "Entry shouldn't exist after deleting")
}
