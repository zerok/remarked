package commandchain

import (
	"sync"

	"github.com/Sirupsen/logrus"
)

// Hub acts as the main control unit that is used to broadcast commands issued
// by a commander to the receivers.
type Hub struct {
	Log        *logrus.Logger
	lock       sync.RWMutex
	receivers  map[*Receiver]struct{}
	commanders map[*Commander]struct{}
}

// RegisterReceiver registers a command receiver with the hub.
func (h *Hub) RegisterReceiver(r *Receiver) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.receivers == nil {
		h.receivers = make(map[*Receiver]struct{})
	}
	if h.Log != nil {
		h.Log.Infof("Registering receiver %s", r)
	}
	r.Commands = make(chan Command)
	h.receivers[r] = struct{}{}
	return nil
}

// BroadcastCommand is used by a commander to issue a specific command to all
// registered receivers.
func (h *Hub) BroadcastCommand(cmd Command, c *Commander) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	for r := range h.receivers {
		if r != nil {
			r.Commands <- cmd
		}
	}
}

// UnregisterReceiver removes the given receiver from the list of known
// command receivers. This also drains the internal command channel for
// that receiver and closes it.
func (h *Hub) UnregisterReceiver(r *Receiver) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.Log != nil {
		h.Log.Infof("Unregistering receiver %s", r)
	}
	if h.receivers == nil {
		h.receivers = make(map[*Receiver]struct{})
	}
	delete(h.receivers, r)
	if h.Log != nil {
		h.Log.Debug("Draining channel")
	}
loop:
	for {
		select {
		case <-r.Commands:
			continue
		default:
			break loop
		}
	}
	if h.Log != nil {
		h.Log.Debug("Channel drained")
	}
	close(r.Commands)
	return nil
}

// RegisterCommander registers a commander with the hub.
func (h *Hub) RegisterCommander(c *Commander) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.Log != nil {
		h.Log.Infof("Registering commander %s", c)
	}
	if h.commanders == nil {
		h.commanders = make(map[*Commander]struct{})
	}
	h.commanders[c] = struct{}{}
	c.Hub = h
	return nil
}

// UnregisterCommander removes the given commander from the hub.
func (h *Hub) UnregisterCommander(c *Commander) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.Log != nil {
		h.Log.Infof("Unregistering commander %s", c)
	}
	if h.commanders == nil {
		h.commanders = make(map[*Commander]struct{})
	}
	delete(h.commanders, c)
	c.Hub = nil
	return nil
}
