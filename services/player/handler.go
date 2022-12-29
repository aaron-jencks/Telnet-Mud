package player

import (
	"mud/entities"
	"time"
)

func ActionHandler(player string) {
	for {
		nextAction := GetNextAction(player)
		if nextAction.Name == "STOP" {
			break
		} else {
			time.Sleep(nextAction.Duration * time.Millisecond)
			nextAction.Handler(CRUD.Retrieve(player).(entities.Player))
		}
	}
}
