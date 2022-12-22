package item_controller

import (
  "item_service"
  "server_utils"
  "logger"
  "encoding/json"
  "net/http"
  "fmt"
)

type itemPayload map[string]string

func itemJsonChecker(w http.ResponseWriter, d []byte) bool {
  var dict itemPayload
  err := json.Unmarshal(d, &dict)
  if err != nil {
    http.Error(w, "Unrecognized payload format", http.StatusBadRequest)
    logger.Error(err)
  }
  return err != nil
}

func itemIdChecker(w http.ResponseWriter, d []byte) bool {
  var id int
  _, err := fmt.Sscanf(string(d), "%d", &id)
  if err != nil {
    http.Error(w, "Unrecognized payload format", http.StatusBadRequest)
    logger.Error(err)
  }
  return err != nil
}

func parseItemId(d []byte) interface{} {
  var id int
  fmt.Sscanf(string(d), "%d", &id)
  return id
}

func parseUpdate(d []byte) (interface{}, interface{}) {
  kint, vint := server_utils.DefaultUpdateParser(d)
  return int(kint.(float64)), vint
}

var RouteMap map[string]map[string]func(http.ResponseWriter, *http.Request) = map[string]map[string]func(http.ResponseWriter, *http.Request) {
  "/items": server_utils.CreateCrudRoutes(item_service.CRUD,
    server_utils.CrudParsers {
      func(d []byte) []interface{} {
        var dict itemPayload
        json.Unmarshal(d, &dict)
        return []interface{} { dict["name"], dict["description"] }
      },
      parseItemId,
      parseUpdate,
      parseItemId,
    }, server_utils.CrudErrorHandlers {
      itemJsonChecker,
      itemIdChecker,
      server_utils.DefaultUpdateErrorHandler,
      itemIdChecker,
    }, server_utils.CrudTranslators {
      server_utils.DefaultTranslator,
      server_utils.DefaultTranslator,
      server_utils.DefaultTranslator,
    }),
}
