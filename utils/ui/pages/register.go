package register_window

import (
  "windows"
  "fmt"
  "player_service"
)

func RegisterUser() windows.Window {
  var username string
  for true {
    fmt.Print("Please enter your username(enter 'cancel' to exit): ")
    fmt.Scanf("%s", &username)
    if username == "cancel" {
      return nil
    } else if player_service.PlayerExists(username) {
      fmt.Println("Sorry, that user already exists, please try again")
    } else {
      break
    }
  }

  player_service.Create(username)

  return nil
}
