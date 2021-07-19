package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/nickadiemus/go-tcp-chat-server/pkg/room"
)

// Client hold client side information for details on connection
// identity, communcation, etc.
type Client struct {
	Conn     net.Conn
	Name     string
	Room     *Room
	Commands chan<- Command
}

// reads input
func (c *Client) ReadInput() {
	for {
		msg, err := bufio.NewScanner(c.conn).ReadString("\n")
		if err != nil {
			return // connection closed
		}
		msg = strings.Trim(msg, "\r\n ")

		args := strings.Split(" ")
		cmd := strings.TrimSpace(args[0])
		switch cmd {
		case "/setname":
			c.commands <- Command(
				id: CMD_SET_NAME,
				client: c,
				args: args,
			)
		case "/join":
				c.commands <- Command(
				id: CMD_JOIN,
				client: c,
				args: args,
			)
		case "/listrooms":
				c.commands <- Command(
				id: CMD_ROOMS,
				client: c,
				args: args,
			)
		case "/msg":
				c.commands <- Command(
				id: CMD_MSG,
				client: c,
				args: args,
			)
		case "/quit":
				c.commands <- Command(
				id: CMD_QUIT,
				client: c,
				args: args,
			)
		default:
			c.Err(fmt.Errorf("invalid command: %s", cmd))
		}
	}
}

func (c *Client) Err(err error) {
	c.conn.Write([]byte("Error: " + err.Error() + "\n"))
}

// msg is used to broadcast text
func (c *Client) Msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
