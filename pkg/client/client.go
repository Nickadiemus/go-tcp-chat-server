package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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
		msg, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			return // connection closed
		}
		msg = strings.Trim(msg, "\r\n ")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])
		switch cmd {
		case "/setname":
			c.Commands <- Command{
				Id:     CMD_SET_NAME,
				Client: c,
				Args:   args,
			}
		case "/join":
			c.Commands <- Command{
				Id:     CMD_JOIN,
				Client: c,
				Args:   args,
			}
		case "/createroom":
			c.Commands <- Command{
				Id:     CMD_CREATE,
				Client: c,
				Args:   args,
			}
		case "/listrooms":
			c.Commands <- Command{
				Id:     CMD_ROOMS,
				Client: c,
				Args:   args,
			}
		case "/msg":
			c.Commands <- Command{
				Id:     CMD_MSG,
				Client: c,
				Args:   args,
			}
		case "/quit":
			c.Commands <- Command{
				Id:     CMD_QUIT,
				Client: c,
				Args:   args,
			}
		default:
			c.Err(fmt.Errorf("invalid command: %s", cmd))
		}
	}
}

func (c *Client) Err(err error) {
	c.Conn.Write([]byte("Error: " + err.Error() + "\n"))
}

// msg is used to broadcast text to specifically the client
func (c *Client) Msg(msg string) {
	if c.Room != nil {
		c.Conn.Write([]byte(fmt.Sprintf("%s:", c.Room.Name) + "> " + msg + "\n"))
	} else {
		c.Conn.Write([]byte("> " + msg + "\n"))
	}
}

func (c *Client) JoinRoom(r Room) {
	c.Room = &r
}
