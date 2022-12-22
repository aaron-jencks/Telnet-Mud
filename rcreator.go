package main

import (
  "fmt"
  "room"
  "loot"
)

func main() {
  var name string
  var description string
  var roomLoot []loot.Loot
  var connections map[string]string

  fmt.Print("Enter a name for the room: ")
  fmt.Scanf("%s", &name)
  fmt.Print("Enter a description for the room: ")
  fmt.Scanf("%s", &description)
}
