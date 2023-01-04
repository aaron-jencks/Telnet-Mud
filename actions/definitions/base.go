package definitions

import (
	"mud/parsing_services/parsing"
	"time"
)

type ActionHandler func() parsing.CommandResponse

type Action struct {
	Name        string
	Handler     ActionHandler
	Duration    time.Duration
	AlwaysValid bool
	ValidModes  []string
}

var ACTION_NOT_FOUND Action = Action{
	Name:        "Not Found",
	Duration:    100,
	AlwaysValid: false,
}
