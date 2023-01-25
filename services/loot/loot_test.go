package loot

import (
	"math/rand"
	"mud/entities"
	"mud/services/item"
	"mud/services/room"
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
	db.DeleteTable("loot")
	entities.SetupRoomTable()
	entities.SetupItemTable()
	entities.SetupLootTable()
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
	iid := rand.Int()
	qty := rand.Int()
	x := rand.Int()
	y := rand.Int()
	z := rand.Int()

	args := []interface{}{
		rid,
		iid,
		qty,
		x,
		y,
		z,
	}

	nid := CRUD.Create(args...)
	is := nid.(entities.Loot)

	assert.Equal(t, rid, is.Room, "Created loot should have the right room id")
	assert.Equal(t, iid, is.Item, "Created loot should have the right item id")
	assert.Equal(t, qty, is.Quantity, "Created loot should have the right quantity")
	assert.Equal(t, x, is.X, "Created loot should have the right x")
	assert.Equal(t, y, is.Y, "Created loot should have the right y")
	assert.Equal(t, z, is.Z, "Created loot should have the right z")
}

func createRandomTestLoot() entities.Loot {
	rid := rand.Int()
	iid := rand.Int()
	qty := rand.Int()
	x := rand.Int()
	y := rand.Int()
	z := rand.Int()

	args := []interface{}{
		rid,
		iid,
		qty,
		x,
		y,
		z,
	}

	nid := CRUD.Create(args...)
	return nid.(entities.Loot)
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

	return room.CRUD.Create(args...).(entities.Room)
}

func createRandomTestItem() entities.Item {
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	pdescription := mtesting.GenerateRandomAsciiString(rand.Intn(64) + 1)

	args := []interface{}{
		pname,
		pdescription,
	}

	return item.CRUD.Create(args...).(entities.Item)
}

func TestUpdate(t *testing.T) {
	ps := createRandomTestLoot()
	ps.Quantity += 5

	nps := CRUD.Update(ps, ps.Id).(entities.Loot)

	assert.Equal(t, ps.Quantity, nps.Quantity, "Quantity should've been updated")
}

func TestDelete(t *testing.T) {
	ps := createRandomTestLoot()
	pid := ps.Id

	CRUD.Retrieve(pid)
	CRUD.Delete(pid)
	psn := CRUD.Retrieve(pid)

	assert.Nil(t, psn, "Entry shouldn't exist after deleting")
}

func TestGetLootForRoom(t *testing.T) {
	// rooms and items must exist first
	r1 := createRandomTestRoom()
	r2 := createRandomTestRoom()
	item := createRandomTestItem()

	rid := r1.Id
	iid := item.Id

	for i := 0; i < 10; i++ {
		qty := rand.Int()
		x := rand.Int()
		y := rand.Int()
		z := rand.Int()

		args := []interface{}{
			rid,
			iid,
			qty,
			x,
			y,
			z,
		}

		CRUD.Create(args...)
	}

	nrid := r2.Id
	qty := rand.Int()
	x := rand.Int()
	y := rand.Int()
	z := rand.Int()

	args := []interface{}{
		nrid,
		iid,
		qty,
		x,
		y,
		z,
	}

	CRUD.Create(args...)

	rloot := GetLootForRoom(room.CRUD.Retrieve(rid).(entities.Room))

	assert.Equal(t, 10, len(rloot), "Should have the right number of loot entities")
	for _, rl := range rloot {
		assert.Equal(t, rid, rl.Room.Id, "Loot should have correct room id")
		assert.Equal(t, iid, rl.Item.Id, "Loot should have correct item id")
	}
}

func TestLootForPosition(t *testing.T) {
	// rooms and items must exist first
	r1 := createRandomTestRoom()
	r2 := createRandomTestRoom()
	item := createRandomTestItem()

	rid := r1.Id
	iid := item.Id
	x := rand.Int()
	y := rand.Int()

	for i := 0; i < 10; i++ {
		qty := rand.Int()
		z := rand.Int()

		args := []interface{}{
			rid,
			iid,
			qty,
			x,
			y,
			z,
		}

		CRUD.Create(args...)
	}

	nrid := r2.Id
	qty := rand.Int()
	z := rand.Int()

	args := []interface{}{
		nrid,
		iid,
		qty,
		x,
		y,
		z,
	}

	CRUD.Create(args...)

	rloot := GetLootForPosition(room.CRUD.Retrieve(rid).(entities.Room), x, y)

	assert.Equal(t, 10, len(rloot), "Should have the right number of loot entities")
	for _, rl := range rloot {
		assert.Equal(t, rid, rl.Room.Id, "Loot should have correct room id")
		assert.Equal(t, iid, rl.Item.Id, "Loot should have correct item id")
		assert.Equal(t, x, rl.X, "Loot should have correct X coord")
		assert.Equal(t, y, rl.Y, "Loot should have correct Y coord")
	}
}
