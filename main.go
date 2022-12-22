package main

import (
  "logger"
  "entities"
  "server_utils"
  "player_controller"
  "room_controller"
  "game_controller"
  "command_controller"
  "transition_controller"
  "item_controller"
  "inventory_controller"
  "net/http"
)

func main() {
  entities.SetupTables()

  logger.Info("Starting Server on port 8080")

  server_utils.CreateHandlers(player_controller.RouteMap)
  server_utils.CreateHandlers(room_controller.RouteMap)
  server_utils.CreateHandlers(game_controller.RouteMap)
  server_utils.CreateHandlers(command_controller.RouteMap)
  server_utils.CreateHandlers(transition_controller.RouteMap)
  server_utils.CreateHandlers(item_controller.RouteMap)
  server_utils.CreateHandlers(inventory_controller.RouteMap)

  if err := http.ListenAndServe(":8080", nil); err != nil {
    logger.Error(err)
    panic(err)
  }
}
