package game_controller

import (
  "server_utils"
  "net/http"
  "fmt"
)

var RouteMap map[string]map[string]func(http.ResponseWriter, *http.Request) = map[string]map[string]func(http.ResponseWriter, *http.Request) {
  "/game": map[string]func(http.ResponseWriter, *http.Request) {
    "POST": func(w http.ResponseWriter, r *http.Request) {
      body := server_utils.ReadFullResponse(r.Body)
      fmt.Fprintf(w, string(body))
    },
  },
}
