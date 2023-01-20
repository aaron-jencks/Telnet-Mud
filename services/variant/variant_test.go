package variant

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
	db.DeleteTable("variants")
	entities.SetupVariantTable()
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
	vid := rand.Int()
	pname := mtesting.GenerateRandomAlnumString(r1)
	picon := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)

	args := []interface{}{
		vid,
		pname,
		picon,
	}

	CRUD.Create(args...)
	is := CRUD.Retrieve(vid, pname).(entities.TileVariant)

	assert.Equal(t, vid, is.Id, "Created variant should have the right id")
	assert.Equal(t, pname, is.Name, "Created variant should have the right name")
	assert.Equal(t, picon, is.Icon, "Created variant should have the right icon")
}

func createRandomTestVariant() entities.TileVariant {
	r1 := rand.Intn(10) + 1
	vid := rand.Int()
	pname := mtesting.GenerateRandomAlnumString(r1)
	picon := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)

	args := []interface{}{
		vid,
		pname,
		picon,
	}

	CRUD.Create(args...)
	return CRUD.Retrieve(vid, pname).(entities.TileVariant)
}

func TestUpdate(t *testing.T) {
	ps := createRandomTestVariant()
	newIcon := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)
	ps.Icon = newIcon

	nps := CRUD.Update(ps, ps.Id, ps.Name).(entities.TileVariant)

	assert.Equal(t, newIcon, nps.Icon, "Icon should've been updated")
}

func TestDelete(t *testing.T) {
	ps := createRandomTestVariant()

	CRUD.Retrieve(ps.Id, ps.Name)
	CRUD.Delete(ps.Id, ps.Name)
	psn := CRUD.Retrieve(ps.Id, ps.Name)

	assert.Nil(t, psn, "Entry shouldn't exist after deleting")
}

func TestGetAllVariants(t *testing.T) {
	ps := createRandomTestVariant()
	ps2 := createRandomTestVariant()
	oid := ps.Id
	ps.Id = ps2.Id
	CRUD.Update(ps, oid, ps.Name)

	for _, variant := range GetAllVariants(ps.Id) {
		if variant.Name == ps.Name {
			assert.Equal(t, ps.Icon, variant.Icon, "Getting variants should return correct information")
		} else if variant.Name == ps2.Name {
			assert.Equal(t, ps2.Icon, variant.Icon, "Getting variants should return correct information")
		} else {
			assert.Fail(t, "GetAllVariants should return all variants for the given id and nothing else")
		}
	}
}
