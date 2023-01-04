package defined

import (
	"mud/actions/definitions"
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	"mud/services/terminal"
	"net"
)

func CreateLoginAction(conn net.Conn, username, password string) definitions.Action {
	return definitions.Action{
		Name:        "Login",
		Duration:    0,
		AlwaysValid: false,
		ValidModes:  []string{"Not Logged In"},
		Handler: func() parsing.CommandResponse {
			var result parsing.CommandResponse = parsing.CommandResponse{
				Person: true,
			}

			if !player.LoginPlayer(username, password, conn) {
				username = player.GetAnonymousUsername(conn)
				player.PushAction(username, CreateInfoAction(conn, "Sorry, either that account doesn't exist or the password is incorrect"))
				result.Info = true
			} else {
				terminal.LoadPlayer(conn, username)
				player.PushAction(username, CreateInfoAction(conn, "Welcome! Please be respectful."))
				result.Chat = true
			}

			return result
		},
	}
}

func CreateLogoutAction(conn net.Conn) definitions.Action {
	return definitions.Action{
		Name:        "Logout",
		Duration:    0,
		AlwaysValid: false,
		ValidModes:  []string{"Logged In"},
		Handler: func() parsing.CommandResponse {
			username := player.GetConnUsername(conn)
			result := parsing.CommandResponse{
				Person: true,
			}

			if !player.LogoutPlayer(username) {
				player.PushAction(username, CreateInfoAction(conn, "Sorry, either that account doesn't exist or isn't currently logged in"))
				result.Info = true
			} else {
				result.Clear = true
			}

			return result
		},
	}
}

func CreateRegisterAction(conn net.Conn, username, password string) definitions.Action {
	return definitions.Action{
		Name:        "Register",
		Duration:    0,
		AlwaysValid: false,
		ValidModes:  []string{"Not Logged In"},
		Handler: func() parsing.CommandResponse {
			anonUsername := player.GetConnUsername(conn)

			if !player.RegisterPlayer(username, password) {
				player.PushAction(anonUsername, CreateInfoAction(conn, "Sorry, that account already exists"))
			} else {
				player.PushAction(anonUsername, CreateInfoAction(conn, "User created successfully, you may now login."))
			}

			return parsing.CommandResponse{
				Info:   true,
				Person: true,
			}
		},
	}
}
