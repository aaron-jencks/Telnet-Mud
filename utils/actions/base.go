package actions

import (
	"mud/entities"
	"time"
)

type ActionHandler func(entities.Player)

type Action struct {
	Name        string
	Handler     ActionHandler
	Duration    time.Duration
	AlwaysValid bool
	ValidModes  []string
}
