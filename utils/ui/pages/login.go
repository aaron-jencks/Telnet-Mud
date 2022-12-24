package login_window

import (
  "windows"
  "fmt"
  "player_service"
  "room_service"
  "terminal"
)

func LoginUser() windows.Window {
  var username string
  for true {
    fmt.Print("Username(enter 'cancel' to exit): ")
    fmt.Scanf("%s", &username)
    if username == "cancel" {
      return nil
    } else if !player_service.PlayerExists(username) {
      fmt.Println("That user doesn't exist, please try again")
    } else {
      break
    }
  }

  player := player_service.Retrieve(username)
  room := room_service.Retrieve(player.Room)

  return terminal.CreateTerminal(player, room)
}
