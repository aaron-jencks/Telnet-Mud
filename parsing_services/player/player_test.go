package player

import (
	"math/rand"
	"mud/entities"
	"mud/utils"
	"mud/utils/io/db"
	"os"
	"path/filepath"
	"testing"

	mtesting "mud/utils/testing"
)

func resetTable() {
	db.DeleteTable("players")
	entities.SetupPlayerTable()
}

func TestMain(m *testing.M) {
	resetTable()

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
		result[i] = rand.Intn(20)
	}
	return result
}

func TestFullCreate(t *testing.T) {
	pname := mtesting.GenerateRandomAlnumString(rand.Intn(10))
	pass := mtesting.GenerateRandomString(rand.Intn(64), mtesting.VisibleAscii)
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

	CRUD.Create(args...)
}
