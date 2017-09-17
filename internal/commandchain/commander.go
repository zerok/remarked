package commandchain

import (
	"context"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

// Commander receives command from a websocket connection and broadcasts them
// through the hub.
type Commander struct {
	Hub   *Hub
	Conn  *websocket.Conn
	Log   *logrus.Logger
	Token string
}

// Handle receives commands from the configured websocket connection and
// forwards them through the hub. Note that the first package received from the
// connection has to be the "auth" command with the correct token.
func (c *Commander) Handle(ctx context.Context) error {
	var authenticated bool
	if c.Hub == nil {
		return fmt.Errorf("no hub set")
	}
	if c.Conn == nil {
		return fmt.Errorf("no conn set")
	}
	if c.Token == "" {
		return fmt.Errorf("no token set")
	}
loop:
	for {
		var cmd Command
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := c.Conn.ReadJSON(&cmd); err != nil {
			return fmt.Errorf("failed to read command from websocket: %s", err)
		}
		if cmd.Type == "auth" {
			if cmd.Token != c.Token {
				return fmt.Errorf("incorrect token")
			}
			authenticated = true
			continue loop
		}
		if !authenticated {
			return fmt.Errorf("channel not authenticated")
		}
		c.Hub.BroadcastCommand(cmd, c)
	}
	return nil
}
func (c *Commander) String() string {
	if c.Conn == nil {
		return "<Commander [unconnected]>"
	}
	return fmt.Sprintf("<Commander from %s>", c.Conn.RemoteAddr())
}
