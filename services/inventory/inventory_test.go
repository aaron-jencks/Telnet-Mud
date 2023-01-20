package inventory

import (
	"math/rand"
	"mud/entities"
	"mud/parsing_services/player"
	"mud/services/item"
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
	db.DeleteTable("inventory")
	entities.SetupPlayerTable()
	entities.SetupItemTable()
	entities.SetupInventoryTable()
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
	iid := rand.Int()
	qty := rand.Int()

	args := []interface{}{
		pid,
		iid,
		qty,
	}

	inid := CRUD.Create(args...)

	is := CRUD.Retrieve(int(inid)).(entities.Inventory)

	assert.Equal(t, pid, is.Player, "Created inventory item should have correct player id")
	assert.Equal(t, iid, is.Item, "Created inventory item should have correct item id")
	assert.Equal(t, qty, is.Quantity, "Created inventory item should have correct quantity")
}

func createRandomTestUser() entities.Player {
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	pass := mtesting.GenerateRandomString(rand.Intn(64)+1, mtesting.VisibleAscii)

	args := []interface{}{
		pname,
		pass,
	}

	return player.CRUD.Retrieve(int(player.CRUD.Create(args...))).(entities.Player)
}

func createRandomTestItem() entities.Item {
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	pdescription := mtesting.GenerateRandomString(rand.Intn(64)+1, mtesting.VisibleAscii)

	args := []interface{}{
		pname,
		pdescription,
	}

	return item.CRUD.Retrieve(int(item.CRUD.Create(args...))).(entities.Item)
}

func createRandomTestInventory(p entities.Player, i entities.Item) entities.Inventory {
	args := []interface{}{
		p.Id,
		i.Id,
		rand.Int(),
	}

	return CRUD.Retrieve(int(CRUD.Create(args...))).(entities.Inventory)
}

func TestUpdate(t *testing.T) {
	player := createRandomTestUser()
	item := createRandomTestItem()
	ps := createRandomTestInventory(player, item)
	ps.Quantity += 5

	nps := CRUD.Update(ps, ps.Id).(entities.Inventory)

	assert.Equal(t, ps.Quantity, nps.Quantity, "Quantity should've been updated")
}

func TestDelete(t *testing.T) {
	player := createRandomTestUser()
	item := createRandomTestItem()
	ps := createRandomTestInventory(player, item)
	pid := ps.Id

	CRUD.Retrieve(int(pid))
	CRUD.Delete(int(pid))
	psn := CRUD.Retrieve(int(pid))

	assert.Nil(t, psn, "Entry shouldn't exist after deleting")
}

func TestPlayerInventory(t *testing.T) {
	player := createRandomTestUser()
	item := createRandomTestItem()
	item2 := createRandomTestItem()
	ps := createRandomTestInventory(player, item)
	ps2 := createRandomTestInventory(player, item2)

	playerInventory := GetPlayerInventory(player)

	assert.Equal(t, 2, len(playerInventory), "Both items should be returned")

	for _, ei := range playerInventory {
		if ei.Item.Id == item.Id {
			assert.Equal(t, item.Name, ei.Item.Name, "Should return correct information for each item")
			assert.Equal(t, item.Description, ei.Item.Description, "Should return correct information for each item")
			assert.Equal(t, ps.Quantity, ei.Quantity, "Should return correct information for each item")
		} else if ei.Item.Id == item2.Id {
			assert.Equal(t, item2.Name, ei.Item.Name, "Should return correct information for each item")
			assert.Equal(t, item2.Description, ei.Item.Description, "Should return correct information for each item")
			assert.Equal(t, ps2.Quantity, ei.Quantity, "Should return correct information for each item")
		} else {
			assert.Fail(t, "Items in inventory should be one of the test items and nothing else")
		}
	}
}

func TestAddItemToInventory(t *testing.T) {
	player := createRandomTestUser()
	item := createRandomTestItem()
	qty := rand.Int()

	nqty := AddItemToInventory(player, item, qty)

	assert.Equal(t, qty, nqty, "Quantity should reflect how many items were added")

	nnqty := AddItemToInventory(player, item, 5)

	assert.Equal(t, nqty+5, nnqty, "Calling with the same item multiple times should increase quantity")
}
