package actions

import "mud/entities"

type ActionHandler func(entities.Player)

type Action struct {
	Name     string
	Handler  ActionHandler
	Duration int
}
