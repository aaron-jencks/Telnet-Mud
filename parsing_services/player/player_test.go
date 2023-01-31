package player

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

var attNames []string = []string{
	"Name",
	"Password",
	"Dex",
	"Str",
	"Int",
	"Wis",
	"Con",
	"Chr",
	"Room",
	"RoomX",
	"RoomY",
	"ActionCapacity",
	"CurrentMode",
}

func resetTable() {
	db.DeleteTable("players")
	entities.SetupPlayerTable()
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

func generateRandomAttributes() []int {
	var result []int = make([]int, 6)
	for i := range result {
		result[i] = rand.Intn(20) + 1
	}
	return result
}

func TestFullCreate(t *testing.T) {
	pname := mtesting.GenerateRandomAlnumString(rand.Intn(10) + 1)
	pass := mtesting.GenerateRandomString(rand.Intn(64)+1, mtesting.VisibleAscii)
	attributes := generateRandomAttributes()
	room := rand.Int()
	attributes = append(attributes, room)

	args := []interface{}{
		pname,
		pass,
	}

	for _, att := range attributes {
		args = append(args, att)
	}

	pid := CRUD.Create(args...)

	ps := pid.(entities.Player)

	parr := playerToArr(ps)

	assert.Greater(t, len(parr), len(args), "Player array should include the ID in it")
	assert.Equal(t, ps.Id, parr[0], "Created player should have the correct ID")

	for pi := range args {
		assert.Equal(t, args[pi], parr[pi+1], "Created player should have the correct attributes applied for %s", attNames[pi])
	}

	assert.Equal(t, 0, ps.RoomX, "Player should be initialized to coord 0")
	assert.Equal(t, 0, ps.RoomY, "Player should be initialized to coord 0")
	assert.Equal(t, utils.DEFAULT_PLAYER_ACTION_LIMIT, ps.ActionCapacity,
		"Player should be initialized to the default action capacity")
	assert.Equal(t, utils.DEFAULT_PLAYER_MODE, ps.CurrentMode,
		"Player should be created in the default game mode")
}

func TestPartialCreate(t *testing.T) {
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	pass := mtesting.GenerateRandomString(rand.Intn(64)+1, mtesting.VisibleAscii)

	args := []interface{}{
		pname,
		pass,
	}

	pid := CRUD.Create(args...)

	ps := pid.(entities.Player)

	parr := playerToArr(ps)

	assert.Greater(t, len(parr), len(args), "Player array should include the ID in it")
	assert.Equal(t, ps.Id, parr[0], "Created player should have the correct ID")

	for pi := range args {
		assert.Equal(t, args[pi], parr[pi+1],
			"Created player should have the correct attributes applied for %s",
			attNames[pi])
	}

	for pi := 0; pi < 6; pi++ {
		assert.Equal(t, 5, parr[pi+3],
			"Player should be initialized with the default attributes for %s",
			attNames[pi+2])
	}

	assert.Equal(t, 1, ps.Room, "Player should be initialized to initial room 1")
	assert.Equal(t, 0, ps.RoomX, "Player should be initialized to coord 0")
	assert.Equal(t, 0, ps.RoomY, "Player should be initialized to coord 0")
	assert.Equal(t, utils.DEFAULT_PLAYER_ACTION_LIMIT, ps.ActionCapacity,
		"Player should be initialized to the default action capacity")
	assert.Equal(t, utils.DEFAULT_PLAYER_MODE, ps.CurrentMode,
		"Player should be created in the default game mode")
}

func TestPartialRoomCreate(t *testing.T) {
	pname := mtesting.GenerateRandomAlnumString(rand.Intn(10) + 1)
	pass := mtesting.GenerateRandomString(rand.Intn(64)+1, mtesting.VisibleAscii)
	room := rand.Int()

	args := []interface{}{
		pname,
		pass,
		room,
	}

	pid := CRUD.Create(args...)

	ps := pid.(entities.Player)

	parr := playerToArr(ps)

	assert.Greater(t, len(parr), len(args), "Player array should include the ID in it")
	assert.Equal(t, ps.Id, parr[0], "Created player should have the correct ID")

	assert.Equal(t, pname, ps.Name, "Created player should have correct username")
	assert.Equal(t, pass, ps.Password, "Created player should have correct password")

	for pi := 0; pi < 6; pi++ {
		assert.Equal(t, 5, parr[pi+3],
			"Player should be initialized with the default attributes for %s",
			attNames[pi+2])
	}

	assert.Equal(t, room, ps.Room, "Player should be initialized to supplied room")
	assert.Equal(t, 0, ps.RoomX, "Player should be initialized to coord 0")
	assert.Equal(t, 0, ps.RoomY, "Player should be initialized to coord 0")
	assert.Equal(t, utils.DEFAULT_PLAYER_ACTION_LIMIT, ps.ActionCapacity,
		"Player should be initialized to the default action capacity")
	assert.Equal(t, utils.DEFAULT_PLAYER_MODE, ps.CurrentMode,
		"Player should be created in the default game mode")
}

func createRandomTestUser() entities.Player {
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	pass := mtesting.GenerateRandomString(rand.Intn(64)+1, mtesting.VisibleAscii)

	args := []interface{}{
		pname,
		pass,
	}

	return CRUD.Create(args...).(entities.Player)
}

func TestUpdate(t *testing.T) {
	ps := createRandomTestUser()
	ps.Dex += 5

	nps := CRUD.Update(ps, ps.Id).(entities.Player)

	assert.Equal(t, ps.Dex, nps.Dex, "Dexterity should've been updated")
}

func TestDelete(t *testing.T) {
	ps := createRandomTestUser()
	pid := ps.Id

	CRUD.Retrieve(int(pid))
	CRUD.Delete(int(pid))
	psn := CRUD.Retrieve(int(pid))

	assert.Nil(t, psn, "Player shouldn't exist after deleting")
}

func TestExistence(t *testing.T) {
	r1 := rand.Intn(10) + 1
	pname := mtesting.GenerateRandomAlnumString(r1)
	pass := mtesting.GenerateRandomString(rand.Intn(64)+1, mtesting.VisibleAscii)

	args := []interface{}{
		pname,
		pass,
	}

	assert.False(t, PlayerExists(pname), "Player shouldn't exist before creation")

	CRUD.Create(args...)

	assert.True(t, PlayerExists(pname), "Created Player should exist")
}

func TestLoginLogout(t *testing.T) {
	ps := createRandomTestUser()

	assert.False(t, PlayerLoggedIn(ps.Name), "Player should not be logged in by default")
	assert.False(t, LogoutPlayer(ps.Name), "Should not be able to logout before logging in")
	assert.True(t, LoginPlayer(ps.Name, ps.Password, nil), "Should login successfully")
	assert.True(t, PlayerLoggedIn(ps.Name), "Player should be logged in after logging in")
	assert.True(t, ConnLoggedIn(nil), "Connection should now be logged in")
	assert.Equal(t, ps.Name, GetConnUsername(nil), "The connection should map to the player's name")
	assert.True(t, LogoutPlayer(ps.Name), "Player should be able to logout after logging in")
}

func TestAnonymousCreation(t *testing.T) {
	nuser := GenerateRandomUsername(nil)
	assert.Nil(t, namesCurrentlyInUse[nuser], "Anonymous names should be logged to avoid duplication")
	assert.Equal(t, nuser, GetAnonymousUsername(nil), "username should be fetchable by connection")
	UnregisterAnonymousName(nuser)
	assert.Equal(t, "", GetAnonymousUsername(nil), "Username should be removed when unregistered")
}
