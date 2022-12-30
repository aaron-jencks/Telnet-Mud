package player

import (
	"mud/entities"
	"mud/utils/actions"
	"time"
)

func ActionHandler(player string) {
	for {
		nextAction := GetNextAction(player)
		if nextAction.Name != actions.ACTION_NOT_FOUND.Name {
			if nextAction.Name == "STOP" {
				break
			} else {
				time.Sleep(nextAction.Duration * time.Millisecond)
				nextAction.Handler(CRUD.Retrieve(player).(entities.Player))
			}
		}
	}
}
