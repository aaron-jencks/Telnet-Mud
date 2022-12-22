package room_controller

import (
  "room_service"
  "server_utils"
  "logger"
  "encoding/json"
  "net/http"
  "fmt"
)

type roomPayload map[string]string

func roomJsonChecker(w http.ResponseWriter, d []byte) bool {
  var dict roomPayload
  err := json.Unmarshal(d, &dict)
  if err != nil {
    http.Error(w, "Unrecognized payload format", http.StatusBadRequest)
    logger.Error(err)
  }
  return err != nil
}

func roomIdChecker(w http.ResponseWriter, d []byte) bool {
  var id int
  _, err := fmt.Sscanf(string(d), "%d", &id)
  if err != nil {
    http.Error(w, "Unrecognized payload format", http.StatusBadRequest)
    logger.Error(err)
  }
  return err != nil
}

func parseRoomId(d []byte) interface{} {
  var id int
  fmt.Sscanf(string(d), "%d", &id)
  return id
}

func parseUpdate(d []byte) (interface{}, interface{}) {
  kint, vint := server_utils.DefaultUpdateParser(d)
  return int(kint.(float64)), vint
}

var RouteMap map[string]map[string]func(http.ResponseWriter, *http.Request) = map[string]map[string]func(http.ResponseWriter, *http.Request) {
  "/rooms": server_utils.CreateCrudRoutes(room_service.CRUD,
    server_utils.CrudParsers {
      func(d []byte) []interface{} {
        var dict roomPayload
        json.Unmarshal(d, &dict)
        return []interface{} { dict["name"], dict["description"] }
      },
      parseRoomId,
      parseUpdate,
      parseRoomId,
    }, server_utils.CrudErrorHandlers {
      roomJsonChecker,
      roomIdChecker,
      server_utils.DefaultUpdateErrorHandler,
      roomIdChecker,
    }, server_utils.CrudTranslators {
      server_utils.DefaultTranslator,
      server_utils.DefaultTranslator,
      server_utils.DefaultTranslator,
    }),
}
