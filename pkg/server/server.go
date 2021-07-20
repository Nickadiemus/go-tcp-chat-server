package server

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/nickadiemus/go-tcp-chat-server/pkg/client"
)

// Server handles processing command inputs and room deligation
type Server struct {
	Rooms    map[string]*client.Room
	Commands chan client.Command
}

func NewServer() *Server {
	fmt.Println("Creating new instance of chat server...")
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
	// check if user belongs to a room
	if c.Room == nil {
		c.Msg("please join a room first")
		return
	}
	// requirement args > 2
	if len(args) < 2 {
		c.Msg("please use /msg <text>")
		return
	}
	// broadcast messageto current room
	c.Room.Broadcast(c, fmt.Sprintf("%s: %s", c.Name, strings.Join(args[1:], " ")))
}

func (s *Server) joinRoom(c *client.Client, args []string) {
	// TODO: validate arg input
	// check if room exists
	r, ok := s.Rooms[args[1]]
	if !ok {
		c.Msg(fmt.Sprintf("Server doesn't exist. You can create %s new one with the command /create %s", args[1], args[1]))
	}
	// check if current client belongs to a room

	r.Members[c.Conn.RemoteAddr()] = c
	s.leaveCurrentRoom(c)
	c.Room = r
	// else send error to client that no room could be found
}
func (s *Server) listRooms(c *client.Client, args []string) {
	// validate arg input
	if len(args) > 1 {
		c.Msg("please use /createRoom")
		return
	}
	msg := ""
	for room := range s.Rooms {
		msg += room + "\n"
	}
	c.Msg(msg)
	// loop over current rooms and write to client availible rooms
}
func (s *Server) createRoom(c *client.Client, args []string) {
	// validate arg input
	if len(args) < 2 {
		c.Msg("please use /createRoom <text>")
		return
	}
	// check for existing room
	_, ok := s.Rooms[args[1]]
	if !ok {
		s.Rooms[args[1]] = &client.Room{
			Name:    args[1],
			Members: make(map[net.Addr]*client.Client),
		}
	} else {
		c.Msg(fmt.Sprintf("%s already exists. Try /join %s jump in!", args[1], args[1]))
	}

}
func (s *Server) quit(c *client.Client, args []string) {
	// terminate client's connection from server
}

func (s *Server) leaveCurrentRoom(c *client.Client) {
	if c.Room != nil {
		delete(c.Room.Members, c.Conn.RemoteAddr())
		c.Room.Broadcast(c, fmt.Sprintf("%s has left", c.Name))
	}
}
