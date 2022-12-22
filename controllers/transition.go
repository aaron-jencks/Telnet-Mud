package transition_controller

import (
  "transition_service"
  "server_utils"
  "encoding/json"
  "net/http"
  "logger"
  "fmt"
)

type TransitionCreatePayload struct {
  Source int
  Target int
  Command string
  CommandArgs []string
}

func idParser(d []byte) interface{} {
  var id int
  fmt.Sscanf(string(d), "%d", &id)
  return id
}

func idHandler(w http.ResponseWriter, d []byte) bool {
  var id int
  _, err := fmt.Sscanf(string(d), "%d", &id)
  if err != nil {
    http.Error(w, "Unrecognized payload format", http.StatusBadRequest)
    logger.Error(err)
  }
  return err != nil
}

func updateParser(d []byte) (interface{}, interface{}) {
  kint, nv := server_utils.DefaultUpdateParser(d)
  return int(kint.(float64)), nv
}

var RouteMap map[string]map[string]func(http.ResponseWriter, *http.Request) = map[string]map[string]func(http.ResponseWriter, *http.Request) {
  "/transitions": server_utils.CreateCrudRoutes(transition_service.CRUD,
    server_utils.CrudParsers {
      func(d []byte) []interface{} {
        var dict TransitionCreatePayload
        json.Unmarshal(d, &dict)
        return []interface{} {dict.Source, dict.Target, dict.Command, dict.CommandArgs}
      },
      idParser,
      updateParser,
      idParser,
    }, server_utils.CrudErrorHandlers {
      func(w http.ResponseWriter, d []byte) bool {
        var dict TransitionCreatePayload
        err := json.Unmarshal(d, &dict)
        if err != nil {
          http.Error(w, "Payload is an unrecognized format", http.StatusBadRequest)
          logger.Error(err)
        }
        return err != nil
      },
      idHandler,
      server_utils.DefaultUpdateErrorHandler,
      idHandler,
    }, server_utils.CrudTranslators {
      server_utils.DefaultTranslator,
      server_utils.DefaultTranslator,
      server_utils.DefaultTranslator,
    }),
}
