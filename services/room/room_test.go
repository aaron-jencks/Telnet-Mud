package room

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
	db.DeleteTable("rooms")
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

func TestCreate(t *testing.T) {
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	pdescription := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	width := rand.Int()
	height := rand.Int()

	args := []interface{}{
		pname,
		pdescription,
		height,
		width,
	}

	inid := CRUD.Create(args...)

	is := inid.(entities.Room)

	assert.Equal(t, pname, is.Name, "Created room should have the right name")
	assert.Equal(t, pdescription, is.Description, "Created room should have the right description")
}

func createRandomTestRoom() entities.Room {
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	pdescription := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	width := rand.Int()
	height := rand.Int()

	args := []interface{}{
		pname,
		pdescription,
		height,
		width,
	}

	return CRUD.Create(args...).(entities.Room)
}

func TestUpdate(t *testing.T) {
	ps := createRandomTestRoom()
	newDescription := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	ps.Description = newDescription

	nps := CRUD.Update(ps, ps.Id).(entities.Room)

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
