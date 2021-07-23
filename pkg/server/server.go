package server

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/nickadiemus/go-tcp-chat-server/pkg/client"
)

// Server defines model for Server
type Server struct {
	rooms    map[string]*client.Room
	commands chan client.Command
}

// NewServer creates a new server instance
func NewServer() *Server {
	log.Println("Creating new instance of chat server...")
	return &Server{
		rooms:    make(map[string]*client.Room),
		commands: make(chan client.Command),
	}
}

// NewClient creates a new cient instance
func (s *Server) NewClient(conn net.Conn) {
	log.Print("Acquired new connection | addr: %s", conn.RemoteAddr().String())
	c := &client.Client{
		Conn:     conn,
		Name:     "anonymous",
		Commands: s.commands,
	}
	c.ReadInput()
}

// Run starts the process of command input
func (s *Server) Run() {
	for cmd := range s.commands {
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

// setName changes the client's name provided arguments
func (s *Server) setName(c *client.Client, args []string) {
	// validate arg input
	c.Name = args[1]
	c.Msg(fmt.Sprintf("name changed to %s", c.Name))
}

// msg broadcastes a client's message to the associated room
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

// joinRoom handles client room joining logic
func (s *Server) joinRoom(c *client.Client, args []string) {
	// TODO: validate arg input
	// check if room exists
	r, ok := s.rooms[args[1]]
	if !ok {
		c.Msg(fmt.Sprintf("Server doesn't exist. You can create %s new one with the command /create %s", args[1], args[1]))
	}
	// check if current client belongs to a room
	r.Members[c.Conn.RemoteAddr()] = c
	s.leaveCurrentRoom(c)
	c.Room = r
	// else send error to client that no room could be found
}

// listRooms sends a list of available rooms to the client
func (s *Server) listRooms(c *client.Client, args []string) {
	// validate arg input
	if len(args) > 1 {
		c.Msg("please use /createRoom")
		return
	}
	msg := ""
	for room := range s.rooms {
		msg += room + "\n"
	}
	c.Msg(msg)
	// loop over current rooms and write to client availible rooms
}

// createRoom creates a new room
func (s *Server) createRoom(c *client.Client, args []string) {
	// validate arg input
	if len(args) < 2 {
		c.Msg("please use /createRoom <text>")
		return
	}
	// check for existing room
	_, ok := s.rooms[args[1]]
	if !ok {
		s.rooms[args[1]] = &client.Room{
			Name:    args[1],
			Members: make(map[net.Addr]*client.Client),
		}
	} else {
		c.Msg(fmt.Sprintf("%s already exists. Try /join %s to jump in!", args[1], args[1]))
	}

}

// quit removes the client from its current room and closes the connection
func (s *Server) quit(c *client.Client, args []string) {
	log.Printf("user %s has disconnected: %s", c.Name, c.Conn.RemoteAddr())
	s.leaveCurrentRoom(c)
	c.Msg("Goodbye!")
	// terminate client's connection from server
	c.Conn.Close()
}

// leaveCurrentRoom leaves the client's current room
func (s *Server) leaveCurrentRoom(c *client.Client) {
	if c.Room != nil {
		delete(c.Room.Members, c.Conn.RemoteAddr())
		c.Room.Broadcast(c, fmt.Sprintf("%s has left", c.Name))
	}
}
