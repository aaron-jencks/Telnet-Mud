package player

import (
	"mud/actions/definitions"
	"mud/controllers/ui"
	"time"
)

func ActionHandler(player string) {
	for {
		nextAction := GetNextAction(player)
		if nextAction.Name != definitions.ACTION_NOT_FOUND.Name {
			if nextAction.Name == "STOP" {
				break
			} else {
				if nextAction.Duration > 0 {
					time.Sleep(nextAction.Duration * time.Millisecond)
				}
				ui.HandleCommandResponse(nextAction.Handler())
			}
		}
	}
}
