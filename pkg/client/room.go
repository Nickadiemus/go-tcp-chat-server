package client

import (
	"net"
)

type Room struct {
	Name    string // name of room
	Members map[net.Addr]*Client
}

func (r *Room) Broadcast(sender *Client, msg string) {
	for addr, m := range r.Members {
		if addr != sender.Conn.RemoteAddr() {
			m.Msg(msg)
		}
	}
}
