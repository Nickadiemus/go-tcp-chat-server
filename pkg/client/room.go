package client

import (
	"net"
)

type Room struct {
	name string // name of room
	size map[net.Addr]*Client
}
