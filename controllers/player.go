package player_controller

import (
  "player_service"
  "server_utils"
  "logger"
  "encoding/json"
  "net/http"
  "fmt"
)

var RouteMap map[string]map[string]func(http.ResponseWriter, *http.Request) = map[string]map[string]func(http.ResponseWriter, *http.Request) {
  "/players/exists": map[string]func(http.ResponseWriter, *http.Request) {
    "POST": func(w http.ResponseWriter, r *http.Request) {
      body := r.Body

      data := server_utils.ReadFullResponse(body)
      response, err := json.Marshal(player_service.PlayerExists(string(data)))
      if err != nil {
        logger.Error(err)
        panic(err)
      }

      fmt.Fprintf(w, string(response))
    },
  },
  "/players": server_utils.CreateCrudRoutes(player_service.CRUD,
    server_utils.CrudParsers {
      func(a []byte) []interface{} { return []interface{}{string(a)} },
      func(k []byte) interface{} { return string(k) },
      server_utils.DefaultUpdateParser,
      func(k []byte) interface{} { return string(k) },
    }, server_utils.CrudErrorHandlers {
      func(w http.ResponseWriter, k []byte) bool {
        if player_service.PlayerExists(string(k)) {
          http.Error(w, "Sorry, that player already exists", http.StatusNotAcceptable)
          return true
        }
        return false
      },
      server_utils.DefaultErrorHandler,
      server_utils.DefaultUpdateErrorHandler,
      server_utils.DefaultErrorHandler,
    }, server_utils.CrudTranslators {
      server_utils.DefaultTranslator,
      server_utils.DefaultTranslator,
      server_utils.DefaultTranslator,
    }),
}
