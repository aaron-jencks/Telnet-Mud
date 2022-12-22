package player_service

import (
  "entities"
  "db"
  "crud"
  "logger"
)
func playerToArr(ps map[string]interface{}) []interface{} {
  return []interface{}{
    int(ps["Id"].(float64)),
    ps["Name"],
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
  logger.Info(data)
  return entities.Player{
    data[1].(int),
    data[2].(string),
    data[3].(int),
    data[4].(int),
    data[5].(int),
    data[6].(int),
    data[7].(int),
    data[8].(int),
    data[9].(int),
  }
}

func playerCreateFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
  if len(args) > 0 {
    id := 0
    if table.CSV.LineCount > 0 {
      id = table.RetrieveLine(table.CSV.LineCount - 1)[1].(int) + 1
    }

    result := make([]interface{}, 9)
    result[0] = id
    result[1] = args[0]

    if len(args) >= 7 {
      for i := 1; i < 7; i++ {
        result[i + 1] = args[i]
      }
      if len(args) > 7 {
        result[8] = args[8]
      } else {
        result[8] = 0
      }
    } else {
      for i := 2; i < 8; i++ {
        result[i] = 5
      }
      if len(args) >= 2 {
        result[8] = args[1]
      } else {
        result[8] = 0
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

