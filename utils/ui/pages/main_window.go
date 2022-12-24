package main_window

import (
  "menu"
  "windows"
  "register_window"
  "login_window"
)

func GetMainWindow() windows.Window {
  var window menu.MenuWindow
  var registerFunc func() windows.Window = func() windows.Window {
    register_window.RegisterUser()
    return &window
  }
  var loginFunc func() windows.Window = func() windows.Window {
    login_window.LoginUser()
    return &window
  }

  window = menu.CreateMenuWindow("Mud v1.0.0", "What would you like to do? ",
    []string{
      "Login",
      "Register",
      "Exit",
    }, map[string]func() windows.Window {
      "Register": registerFunc,
      "Login": loginFunc,
    },
    0, 0, 20, 80)
  return &window
}
