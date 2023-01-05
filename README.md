# Telnet Mud
A Telnet MUD implemented in Gcc-Go

## Connecting From Windows
If you want to connect to the server from a windows machine, my recommendation would be to use PuTTY, you can see a guide below.

## Connecting From Linux
Linux has a built in `telnet` command for connecting to the server, but you do need to enable LINEMODE. To do this, you simply need to hit the escape character `Ctrl+[` by default and then type `mode line` to enable line mode. The server will not work with receiving a single character at a time, at least currently.
