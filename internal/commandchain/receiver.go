package commandchain

import (
	"context"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

// Receiver is a simple receipient of commands coming through the Commands
// channel. These are forwarded through the websocket connection to the
// client's browser.
type Receiver struct {
	Conn     *websocket.Conn
	Commands chan Command
	Log      *logrus.Logger
}

// Handle waits for input from the commands channel in order to forward the
// command to the client's browser.
func (r *Receiver) Handle(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	cancelCtx, cancel := context.WithCancel(ctx)
	r.Conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	r.Conn.SetPongHandler(func(data string) error {
		if r.Log != nil {
			r.Log.Infof("Pong received: %s", data)
		}
		return r.Conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	})
	r.Conn.SetCloseHandler(func(code int, text string) error {
		if r.Log != nil {
			r.Log.Infof("Close received: %d %s", code, text)
		}
		cancel()
		return nil
	})
	for {
		select {
		case <-ticker.C:
			if err := r.Conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second*5)); err != nil {
				return fmt.Errorf("failed to ping client: %s", err.Error())
			}
		case <-cancelCtx.Done():
			return ctx.Err()
		case cmd := <-r.Commands:
			if err := r.Conn.WriteJSON(cmd); err != nil {
				return fmt.Errorf("failed to send command: %s", err.Error())
			}
		}
	}
	return nil
}

func (r *Receiver) String() string {
	if r.Conn == nil {
		return "<Receiver [unconnected]>"
	}
	return fmt.Sprintf("<Receiver from %s>", r.Conn.RemoteAddr())
}
