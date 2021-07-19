package server

import (
	"fmt"
	"log"
	"net"

	"github.com/nickadiemus/go-tcp-chat-server/pkg/client"
)

// Server handles processing command inputs and room deligation
type Server struct {
	Commands chan client.Command
	Rooms    map[string]*client.Room
}

func NewServer() *Server {
	return &Server{
		Rooms:    make(map[string]*client.Room),
		Commands: make(chan client.Command),
	}
}

func (s *Server) NewClient(conn net.Conn) {
	log.Print("Acquired new connection | addr: %s", conn.RemoteAddr().String())

	c := &client.Client{
		Conn:     conn,
		Name:     "anonymous",
		Commands: s.Commands,
	}
	c.ReadInput()

}

func (s *Server) Run() {
	for cmd := range s.Commands {
		switch cmd.Id {
		case client.CMD_SET_NAME:
			s.setName(cmd.Client, cmd.Args)
		case client.CMD_MSG:
			s.msg(cmd.Client, cmd.Args)
		case client.CMD_JOIN:
			s.joinRoom(cmd.Client, cmd.Args)
		case client.CMD_ROOMS:
			s.listRooms(cmd.Client, cmd.Args)
		case client.CMD_CREATE:
			s.createRoom(cmd.Client, cmd.Args)
		case client.CMD_QUIT:
			s.quit(cmd.Client, cmd.Args)
		}

	}
}

func (s *Server) setName(c *client.Client, args []string) {
	// validate arg input
	c.Name = args[1]
	c.Msg(fmt.Sprintf("name changed to %s", c.Name))
}

func (s *Server) msg(c *client.Client, args []string) {
	// check if msg format is correct
	// requirement args > 2

	// broadcast message to current room

}

func (s *Server) joinRoom(c *client.Client, args []string) {
	// validate arg input

	// check if room exists

	// join if found

	// else send error to client that no room could be found
}
func (s *Server) listRooms(c *client.Client, args []string) {
	// validate arg input

	// loop over current rooms and write to client availible rooms
}
func (s *Server) createRoom(c *client.Client, args []string) {
	// validate arg input

	// check for existing room

	// if room doesn't exist create it

	// else return to user that room already exists
}
func (s *Server) quit(c *client.Client, args []string) {
	// terminate client's connection from server
}
