package command_controller

import (
  "command_service"
  "server_utils"
  "encoding/json"
  "net/http"
  "logger"
)

var RouteMap map[string]map[string]func(http.ResponseWriter, *http.Request) = map[string]map[string]func(http.ResponseWriter, *http.Request) {
  "/commands": server_utils.CreateCrudRoutes(command_service.CRUD,
    server_utils.CrudParsers {
      func(d []byte) []interface{} {
        var dict command_service.ExpandedCommand
        json.Unmarshal(d, &dict)
        return []interface{} {dict.Name, dict.Args}
      },
      func(d []byte) interface{} { return string(d) },
      server_utils.DefaultUpdateParser,
      func(d []byte) interface{} { return string(d) },
    }, server_utils.CrudErrorHandlers {
      func(w http.ResponseWriter, d []byte) bool {
        var dict command_service.ExpandedCommand
        err := json.Unmarshal(d, &dict)
        if err != nil {
          http.Error(w, "Payload is an unrecognized format", http.StatusBadRequest)
          logger.Error(err)
        }
        return err != nil
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
