package actions

import (
	"mud/entities"
	"mud/services/parsing"
	"time"
)

type ActionHandler func(entities.Player) parsing.CommandResponse

type Action struct {
	Name        string
	Handler     ActionHandler
	Duration    time.Duration
	AlwaysValid bool
	ValidModes  []string
}
