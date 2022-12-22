package roomMap

import (
  "room"
)

type RoomMap struct {
  Rooms []room.Room
  RoomMapping map[string]int
  Edges map[string][]int
}
