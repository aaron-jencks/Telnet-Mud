package telnet_options

import (
  "net"
)

var OptionMap map[byte]func(net.Conn, []byte) ([]byte, int)= make(map[byte]func(net.Conn, []byte) ([]byte, int))

func DoesOption(opt byte) bool {
  _, ok := OptionMap[opt]
  return ok
}

func DoOption(opt byte, conn net.Conn, data []byte) ([]byte, int) {
  if DoesOption(opt) {
    return OptionMap[opt](conn, data)
  }
  return data, 0
}
