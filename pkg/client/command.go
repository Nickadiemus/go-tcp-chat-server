package client

type CommandID int

const (
	CMD_SET_NAME CommandID = iota
	CMD_MSG
	CMD_ROOMS
	CMD_JOIN
	CMD_CREATE
	CMD_QUIT
)

// Command type holds information for user sent commands
type Command struct {
	Id     CommandID
	Client *Client
	Args   []string
}
