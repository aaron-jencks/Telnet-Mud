package item

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
	db.DeleteTable("items")
	entities.SetupItemTable()
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

	args := []interface{}{
		pname,
		pdescription,
	}

	inid := CRUD.Create(args...)

	is := CRUD.Retrieve(int(inid)).(entities.Item)

	assert.Equal(t, pname, is.Name, "Created item should have the right name")
	assert.Equal(t, pdescription, is.Description, "Created item should have the right description")
}

func createRandomTestItem() entities.Item {
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	pdescription := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)

	args := []interface{}{
		pname,
		pdescription,
	}

	return CRUD.Retrieve(int(CRUD.Create(args...))).(entities.Item)
}

func TestUpdate(t *testing.T) {
	ps := createRandomTestItem()
	newDescription := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	ps.Description = newDescription

	nps := CRUD.Update(ps, ps.Id).(entities.Item)

	assert.Equal(t, newDescription, nps.Description, "Description should've been updated")
}

func TestDelete(t *testing.T) {
	ps := createRandomTestItem()
	pid := ps.Id

	CRUD.Retrieve(int(pid))
	CRUD.Delete(int(pid))
	psn := CRUD.Retrieve(int(pid))

	assert.Nil(t, psn, "Entry shouldn't exist after deleting")
}
