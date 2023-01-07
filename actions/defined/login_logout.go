package defined

import (
	"mud/actions/definitions"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/parsing_services/parsing/utils"
	"mud/parsing_services/player"
	"mud/services/terminal"
	"net"
)

func CreateLoginAction(conn net.Conn, username, password string) definitions.Action {
	return definitions.Action{
		Name:       "Login",
		ValidModes: []string{"Not Logged In"},
		Handler: func() parsing.CommandResponse {
			var result parsing.CommandResponse = parsing.CommandResponse{
				Conn:   conn,
				Person: true,
			}

			if !player.LoginPlayer(username, password, conn) {
				username = player.GetAnonymousUsername(conn)
				player.EnqueueAction(username, CreateInfoAction(conn, "Sorry, either that account doesn't exist or the password is incorrect"))
				result.Info = true
			} else {
				p := player.CRUD.Retrieve(username).(entities.Player)
				terminal.LoadPlayer(conn, p)
				player.EnqueueAction(username, CreateInfoAction(conn, "Welcome! Please be respectful."))
				result = utils.GetDefaultRepaintCommandResponse(conn)
			}

			return result
		},
	}
}

func CreateLogoutAction(conn net.Conn) definitions.Action {
	return definitions.Action{
		Name:       "Logout",
		ValidModes: []string{"Logged In"},
		Handler: func() parsing.CommandResponse {
			username := player.GetConnUsername(conn)
			result := parsing.CommandResponse{
				Conn:   conn,
				Person: true,
			}

			if !player.LogoutPlayer(username) {
				player.EnqueueAction(username, CreateInfoAction(conn, "Sorry, either that account doesn't exist or isn't currently logged in"))
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
		Name:       "Register",
		ValidModes: []string{"Not Logged In"},
		Handler: func() parsing.CommandResponse {
			anonUsername := player.GetConnUsername(conn)

			if !player.RegisterPlayer(username, password) {
				player.EnqueueAction(anonUsername, CreateInfoAction(conn, "Sorry, that account already exists"))
			} else {
				player.EnqueueAction(anonUsername, CreateInfoAction(conn, "User created successfully, you may now login."))
			}

			return utils.GetDefaultInfoCommandResponse(conn)
		},
	}
}
