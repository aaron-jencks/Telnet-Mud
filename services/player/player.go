package player

import (
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
	"net"
)

func playerToArr(ps map[string]interface{}) []interface{} {
	return []interface{}{
		int(ps["Id"].(float64)),
		ps["Name"],
		ps["Password"],
		int(ps["Dex"].(float64)),
		int(ps["Str"].(float64)),
		int(ps["Int"].(float64)),
		int(ps["Wis"].(float64)),
		int(ps["Con"].(float64)),
		int(ps["Chr"].(float64)),
		int(ps["Room"].(float64)),
	}
}

func playerFromArr(data []interface{}) interface{} {
	return entities.Player{
		Id:       data[1].(int),
		Name:     data[2].(string),
		Password: data[3].(string),
		Dex:      data[4].(int),
		Str:      data[5].(int),
		Int:      data[6].(int),
		Wis:      data[7].(int),
		Con:      data[8].(int),
		Chr:      data[9].(int),
		Room:     data[10].(int),
	}
}

func playerCreateFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	if len(args) > 0 {
		id := 0
		if table.CSV.LineCount > 0 {
			id = table.RetrieveLine(table.CSV.LineCount - 1)[1].(int) + 1
		}

		result := make([]interface{}, 10)
		result[0] = id
		result[1] = args[0]
		result[2] = args[1]

		if len(args) >= 7 {
			for i := 2; i < 8; i++ {
				result[i+1] = args[i]
			}
			if len(args) > 7 {
				result[9] = args[9]
			} else {
				result[9] = 0
			}
		} else {
			for i := 3; i < 9; i++ {
				result[i] = 5
			}
			if len(args) >= 3 {
				result[9] = args[2]
			} else {
				result[9] = 0
			}
		}

		return result
	}

	return []interface{}{}
}

var CRUD crud.Crud = crud.CreateCrud("players", playerToArr, playerFromArr, playerCreateFunc)

func PlayerExists(name string) bool {
	table := CRUD.FetchTable()
	result := table.Query(name, "Name")
	return len(result) > 0
}

var LoggedInPlayerMap map[string]net.Conn = make(map[string]net.Conn)
var PlayerConnectionMap map[net.Conn]string = make(map[net.Conn]string)

func LoginPlayer(name string, password string, conn net.Conn) bool {
	if PlayerExists(name) && CRUD.Retrieve(name).(entities.Player).Password == password {
		_, ok := LoggedInPlayerMap[name]
		if !ok {
			LoggedInPlayerMap[name] = conn
			PlayerConnectionMap[conn] = name
			return true
		}
	}
	return false
}

func LogoutPlayer(name string) bool {
	conn, ok := LoggedInPlayerMap[name]
	if ok {
		delete(LoggedInPlayerMap, name)
		delete(PlayerConnectionMap, conn)
		return true
	}
	return false
}

func PlayerLoggedIn(name string) bool {
	_, ok := LoggedInPlayerMap[name]
	return ok
}

func RegisterPlayer(name string, password string) bool {
	if !PlayerExists(name) {
		CRUD.Create(name, password)
		return true
	}
	return false
}

func ConnLoggedIn(conn net.Conn) bool {
	_, ok := PlayerConnectionMap[conn]
	return ok
}
