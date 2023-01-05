package crud

import "net"

var DefaultCrudModes []string = []string{"Building"}

// This is a function that accepts the arguments being passed to the Crud function
// and formats them into their corresponding representation and returns them as an array
type ArgumentFormatter func([]string) []interface{}

// This is a function that accepts a list of string arguments to be parsed
// and determines if they are valid for the operation
type ArgumentValidator func(net.Conn, []string) bool

// This is a function used for simple retrieval with a Crud function,
// Used with the Crud Retrieve method, it takes an array of string arguments
// and returns the search query argument for Crud.Retrieve
type RetrieveArgumentFormatter func([]string) interface{}

// Takes the given retrieved value from the Crud Retrieve function
// and updates the given property(s) and returns an updated copy
// Syntax is (oldValue, propertyName, newPropertyValue(s))
type UpdateNewValueFormatter func(interface{}, string, []string) interface{}

// This is a function that accepts the response of a Crud call and formats it into a string
// This string is sent to the user via a system message
type ResponseFormatter func(interface{}) string

type CrudValidator func([]string) bool
type CrudExecutor func() interface{}
type CrudResponseHandler func(interface{})
