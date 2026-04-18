package app

import (
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"
)

// Hub owns the live-session state keyed by notebook ID.
// Each Room holds the set of subscribers (one writer + many viewers) watching
// a particular notebook in real time.
type Hub struct {
	store *Store
	mu    sync.Mutex
	rooms map[string]*Room
}

func NewHub(store *Store) *Hub {
	return &Hub{store: store, rooms: map[string]*Room{}}
}

// room returns the existing room for id, creating an empty one if needed.
func (h *Hub) room(id string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	r, ok := h.rooms[id]
	if !ok {
		r = &Room{hub: h, id: id, subs: map[*Sub]bool{}}
		h.rooms[id] = r
	}
	return r
}

// Snapshot returns live state for a notebook: viewer count, writer name, live
// flag. Safe to call for any notebook (even ones with no live room yet).
func (h *Hub) Snapshot(id string) (viewers int, writer string, live bool) {
	h.mu.Lock()
	r, ok := h.rooms[id]
	h.mu.Unlock()
	if !ok {
		return 0, "", false
	}
	return r.snapshot()
}

// LiveIDs lists notebooks that currently have an active writer.
func (h *Hub) LiveIDs() map[string]bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make(map[string]bool, len(h.rooms))
	for id, r := range h.rooms {
		r.mu.Lock()
		if r.writer != nil {
			out[id] = true
		}
		r.mu.Unlock()
	}
	return out
}

// Room coordinates a live notebook's subscribers and persists ops as they arrive.
type Room struct {
	hub *Hub
	id  string

	mu     sync.Mutex
	subs   map[*Sub]bool
	writer *Sub
	cursor int
}

type Sub struct {
	ID       uint64
	Name     string
	IsWriter bool
	send     chan []byte
}

var subSeq uint64

func NewSub(name string, isWriter bool) *Sub {
	return &Sub{
		ID:       atomic.AddUint64(&subSeq, 1),
		Name:     name,
		IsWriter: isWriter,
		send:     make(chan []byte, 64),
	}
}

func (s *Sub) Outbox() <-chan []byte { return s.send }

func (r *Room) snapshot() (int, string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	viewers := 0
	for s := range r.subs {
		if !s.IsWriter {
			viewers++
		}
	}
	name := ""
	live := false
	if r.writer != nil {
		name = r.writer.Name
		live = true
	}
	return viewers, name, live
}

// Join adds a sub. If claim=true, tries to become the writer; returns false if
// another writer already holds the room.
func (r *Room) Join(s *Sub, claim bool) bool {
	r.mu.Lock()
	if claim {
		if r.writer != nil {
			r.mu.Unlock()
			return false
		}
		r.writer = s
	}
	r.subs[s] = true
	r.mu.Unlock()
	r.broadcastPresence()
	return true
}

// Leave removes a sub. If it was the writer, the room is marked not-live and
// the notebook is flushed to disk.
func (r *Room) Leave(s *Sub) {
	r.mu.Lock()
	delete(r.subs, s)
	wasWriter := r.writer == s
	if wasWriter {
		r.writer = nil
	}
	empty := len(r.subs) == 0 && r.writer == nil
	r.mu.Unlock()
	close(s.send)
	if wasWriter {
		if n, ok := r.hub.store.Get(r.id); ok {
			n.Lock()
			_ = r.hub.store.Save(n)
			n.Unlock()
		}
	}
	r.broadcastPresence()
	if empty {
		r.hub.mu.Lock()
		delete(r.hub.rooms, r.id)
		r.hub.mu.Unlock()
	}
}

// ApplyOp validates, applies, and broadcasts an op. Only the writer should call.
func (r *Room) ApplyOp(n *Notebook, op Op) {
	n.Lock()
	op.Ts = time.Now().UnixMilli() - n.CreatedAt
	n.Text = Apply(n.Text, op)
	n.Ops = append(n.Ops, op)
	n.UpdatedAt = time.Now().UnixMilli()
	n.Unlock()
	msg := mustJSON(map[string]any{"t": "op", "op": op})
	r.broadcastToViewers(msg)
}

// SetCursor updates the writer's caret and broadcasts it to viewers.
func (r *Room) SetCursor(pos int) {
	r.mu.Lock()
	r.cursor = pos
	r.mu.Unlock()
	r.broadcastToViewers(mustJSON(map[string]any{"t": "cur", "p": pos}))
}

// CurrentCursor returns the last-known writer cursor position.
func (r *Room) CurrentCursor() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.cursor
}

func (r *Room) broadcast(msg []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for s := range r.subs {
		select {
		case s.send <- msg:
		default:
			// Slow consumer: drop the message rather than block the room.
		}
	}
}

func (r *Room) broadcastToViewers(msg []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for s := range r.subs {
		if s.IsWriter {
			continue
		}
		select {
		case s.send <- msg:
		default:
		}
	}
}

func (r *Room) broadcastPresence() {
	viewers, writerName, live := r.snapshot()
	r.broadcast(mustJSON(map[string]any{
		"t":           "presence",
		"viewers":     viewers,
		"live":        live,
		"writer_name": writerName,
	}))
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
