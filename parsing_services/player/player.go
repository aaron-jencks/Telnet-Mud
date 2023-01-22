package player

import (
	"database/sql"
	"fmt"
	"mud/actions/definitions"
	"mud/entities"
	"mud/utils"
	"mud/utils/crud"
	"mud/utils/io/db"
	"net"
	"sync"
)

func playerToArr(ps interface{}) []interface{} {
	pec := ps.(entities.Player)
	return []interface{}{
		pec.Id,
		pec.Name,
		pec.Password,
		pec.Dex,
		pec.Str,
		pec.Int,
		pec.Wis,
		pec.Con,
		pec.Chr,
		pec.Room,
		pec.RoomX,
		pec.RoomY,
		pec.ActionCapacity,
		pec.CurrentMode,
	}
}

func playerFromArr(data []interface{}) interface{} {
	return entities.Player{
		Id:             data[0].(int),
		Name:           data[1].(string),
		Password:       data[2].(string),
		Dex:            data[3].(int),
		Str:            data[4].(int),
		Int:            data[5].(int),
		Wis:            data[6].(int),
		Con:            data[7].(int),
		Chr:            data[8].(int),
		Room:           data[9].(int),
		RoomX:          data[10].(int),
		RoomY:          data[11].(int),
		ActionCapacity: data[12].(int),
		CurrentMode:    data[13].(string),
	}
}

func playerCreateFunc(table db.TableDefinition, args ...interface{}) []interface{} {
	if len(args) > 0 {
		result := make([]interface{}, 13)
		result[0] = args[0]
		result[1] = args[1]

		if len(args) >= 7 {
			for i := 2; i < 8; i++ {
				result[i] = args[i]
			}
			if len(args) > 7 {
				result[8] = args[8]
			} else {
				result[8] = 1
			}
		} else {
			for i := 2; i < 8; i++ {
				result[i] = 5
			}
			if len(args) >= 3 {
				result[8] = args[2]
			} else {
				result[8] = 1
			}
		}

		// room coords
		result[9] = 0
		result[10] = 0

		// action queue limit
		result[11] = utils.DEFAULT_PLAYER_ACTION_LIMIT

		// The game mode of the player
		result[12] = utils.DEFAULT_PLAYER_MODE

		return result
	}

	return []interface{}{}
}

func playerSelectorFormatter(args []interface{}) string {
	if len(args) == 1 {
		return fmt.Sprintf("Id=%d", args[0].(int))
	}
	return ""
}

func playerRowScanner(row *sql.Rows) (interface{}, error) {
	var result entities.Player = entities.Player{}
	err := row.Scan(&result.Id,
		&result.Name, &result.Password,
		&result.Dex, &result.Str, &result.Int, &result.Wis, &result.Con, &result.Chr,
		&result.Room, &result.RoomX, &result.RoomY,
		&result.ActionCapacity, &result.CurrentMode)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func playerUpdateFunc(oldValue, newValue interface{}) []crud.RowModStruct {
	ops := oldValue.(entities.Player)
	nps := newValue.(entities.Player)

	var result []crud.RowModStruct

	if ops.Name != nps.Name {
		result = append(result, crud.RowModStruct{
			Column:   "Name",
			NewValue: fmt.Sprintf("\"%s\"", nps.Name),
		})
	}
	if ops.Password != nps.Password {
		result = append(result, crud.RowModStruct{
			Column:   "Password",
			NewValue: fmt.Sprintf("\"%s\"", nps.Password),
		})
	}
	if ops.Dex != nps.Dex {
		result = append(result, crud.RowModStruct{
			Column:   "Dex",
			NewValue: nps.Dex,
		})
	}
	if ops.Str != nps.Str {
		result = append(result, crud.RowModStruct{
			Column:   "Str",
			NewValue: nps.Str,
		})
	}
	if ops.Int != nps.Int {
		result = append(result, crud.RowModStruct{
			Column:   "Int",
			NewValue: nps.Int,
		})
	}
	if ops.Wis != nps.Wis {
		result = append(result, crud.RowModStruct{
			Column:   "Wis",
			NewValue: nps.Wis,
		})
	}
	if ops.Con != nps.Con {
		result = append(result, crud.RowModStruct{
			Column:   "Con",
			NewValue: nps.Con,
		})
	}
	if ops.Chr != nps.Chr {
		result = append(result, crud.RowModStruct{
			Column:   "Chr",
			NewValue: nps.Chr,
		})
	}
	if ops.Room != nps.Room {
		result = append(result, crud.RowModStruct{
			Column:   "Room",
			NewValue: nps.Room,
		})
	}
	if ops.RoomX != nps.RoomX {
		result = append(result, crud.RowModStruct{
			Column:   "RoomX",
			NewValue: nps.RoomX,
		})
	}
	if ops.RoomY != nps.RoomY {
		result = append(result, crud.RowModStruct{
			Column:   "RoomY",
			NewValue: nps.RoomY,
		})
	}
	if ops.ActionCapacity != nps.ActionCapacity {
		result = append(result, crud.RowModStruct{
			Column:   "ActionCapacity",
			NewValue: nps.ActionCapacity,
		})
	}
	if ops.CurrentMode != nps.CurrentMode {
		result = append(result, crud.RowModStruct{
			Column:   "CurrentMode",
			NewValue: fmt.Sprintf("\"%s\"", nps.CurrentMode),
		})
	}

	return result
}

var CRUD crud.Crud = crud.CreateCrud("players",
	playerSelectorFormatter, playerToArr, playerRowScanner, playerFromArr, playerCreateFunc, playerUpdateFunc)

func FetchPlayerByName(name string) entities.Player {
	if PlayerExists(name) {
		table := CRUD.FetchTable()
		result := table.QueryData(fmt.Sprintf("Name=\"%s\"", name), playerRowScanner)
		return result[0].(entities.Player)
	}
	return entities.Player{}
}

func PlayerExists(name string) bool {
	table := CRUD.FetchTable()
	result := table.QueryData(fmt.Sprintf("Name=\"%s\"", name), playerRowScanner)
	return len(result) > 0
}

var LoggedInPlayerMap map[string]net.Conn = make(map[string]net.Conn)
var PlayerConnectionMap map[net.Conn]string = make(map[net.Conn]string)
var PlayerRegistrationLock sync.Mutex = sync.Mutex{}
var PlayerQueueMap map[string]*definitions.ActionQueue = make(map[string]*definitions.ActionQueue)
var PlayerQueueMapLock sync.Mutex = sync.Mutex{}

func CreateAnonymousHandler(username string) {
	PlayerQueueMapLock.Lock()
	defer PlayerQueueMapLock.Unlock()

	_, ok := PlayerQueueMap[username]
	if !ok {
		PlayerQueueMap[username] = definitions.CreateActionQueue(utils.DEFAULT_GLOBAL_ACTION_LIMIT)

		go ActionHandler(username)
	}
}

func LoginPlayer(name string, password string, conn net.Conn) bool {
	if PlayerExists(name) && FetchPlayerByName(name).Password == password {
		PlayerRegistrationLock.Lock()
		defer PlayerRegistrationLock.Unlock()
		_, ok := LoggedInPlayerMap[name]
		if !ok {
			LoggedInPlayerMap[name] = conn
			PlayerQueueMapLock.Lock()
			PlayerQueueMap[name] = definitions.CreateActionQueue(FetchPlayerByName(name).ActionCapacity)
			PlayerQueueMapLock.Unlock()
			PlayerConnectionMap[conn] = name

			// Launch the action processing loop for the player
			go ActionHandler(name)

			return true
		}
	}
	return false
}

func LogoutPlayer(name string) bool {
	PlayerRegistrationLock.Lock()
	defer PlayerRegistrationLock.Unlock()
	conn, ok := LoggedInPlayerMap[name]
	if ok {
		delete(LoggedInPlayerMap, name)
		delete(PlayerConnectionMap, conn)

		// Signals to the action handler to quit
		UnregisterHandler(name)

		return true
	}
	return false
}

func UnregisterHandler(name string) {
	EnqueueAction(name, definitions.Action{
		Name: "STOP",
	})
}

func PlayerLoggedIn(name string) bool {
	PlayerRegistrationLock.Lock()
	defer PlayerRegistrationLock.Unlock()
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
	PlayerRegistrationLock.Lock()
	defer PlayerRegistrationLock.Unlock()
	_, ok := PlayerConnectionMap[conn]
	return ok
}

func EnqueueAction(p string, a definitions.Action) {
	PlayerQueueMapLock.Lock()
	q := PlayerQueueMap[p]
	PlayerQueueMapLock.Unlock()
	q.Enqueue(a)
}

func EnqueueActions(player string, actions []definitions.Action) {
	for _, action := range actions {
		EnqueueAction(player, action)
	}
}

func GetNextAction(player string) definitions.Action {
	PlayerQueueMapLock.Lock()
	q := PlayerQueueMap[player]
	PlayerQueueMapLock.Unlock()
	return q.Dequeue()
}

func GetConnUsername(conn net.Conn) string {
	if ConnLoggedIn(conn) {
		PlayerRegistrationLock.Lock()
		defer PlayerRegistrationLock.Unlock()
		return PlayerConnectionMap[conn]
	} else {
		return GetAnonymousUsername(conn)
	}
}
