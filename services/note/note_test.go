package note

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
	db.DeleteTable("notes")
	entities.SetupNoteTable()
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
	pid := rand.Int()
	r1 := rand.Intn(10) + 1
	ntitle := mtesting.GenerateRandomAlnumString(r1)
	ncontents := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)

	args := []interface{}{
		pid,
		ntitle,
		ncontents,
	}

	nid := CRUD.Create(args...)
	is := nid.(entities.Note)

	assert.Equal(t, pid, is.Player, "Created note should have the right player id")
	assert.Equal(t, ntitle, is.Title, "Created note should have the right title")
	assert.Equal(t, ncontents, is.Contents, "Created note should have the right contents")
}

func createRandomTestNote() entities.Note {
	pid := rand.Int()
	r1 := rand.Intn(10) + 1
	ntitle := mtesting.GenerateRandomAlnumString(r1)
	ncontents := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)

	args := []interface{}{
		pid,
		ntitle,
		ncontents,
	}

	nid := CRUD.Create(args...)
	return nid.(entities.Note)
}

func TestUpdate(t *testing.T) {
	ps := createRandomTestNote()
	newContents := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	ps.Contents = newContents

	nps := CRUD.Update(ps, ps.Id).(entities.Note)

	assert.Equal(t, newContents, nps.Contents, "Contents should've been updated")
}

func TestDelete(t *testing.T) {
	ps := createRandomTestNote()
	pid := ps.Id

	CRUD.Retrieve(pid)
	CRUD.Delete(pid)
	psn := CRUD.Retrieve(pid)

	assert.Nil(t, psn, "Entry shouldn't exist after deleting")
}
