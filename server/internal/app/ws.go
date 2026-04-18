package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// inboundMsg is what clients send us.
type inboundMsg struct {
	T string `json:"t"`
	P int    `json:"p"`
	N int    `json:"n"`
	S string `json:"s"`
}

// WSWrite handles the writer end of a notebook. URL path: /ws/write/{id}
// Query: ?name=alice
func (s *Server) WSWrite(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "anon"
	}
	n, ok := s.store.Get(id)
	if !ok {
		http.Error(w, "notebook not found", http.StatusNotFound)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	room := s.hub.room(id)
	sub := NewSub(name, true)
	if ok := room.Join(sub, true); !ok {
		_ = conn.WriteJSON(map[string]any{"t": "err", "msg": "already being written"})
		_ = conn.Close()
		return
	}

	s.runConn(conn, sub, room, n, true)
}

// WSWatch handles viewers. URL path: /ws/watch/{id}
func (s *Server) WSWatch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "anon"
	}
	n, ok := s.store.Get(id)
	if !ok {
		http.Error(w, "notebook not found", http.StatusNotFound)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	room := s.hub.room(id)
	sub := NewSub(name, false)
	room.Join(sub, false)

	s.runConn(conn, sub, room, n, false)
}

// runConn wires the subscriber's send channel to the websocket, handles the
// initial snapshot, and loops on incoming messages.
func (s *Server) runConn(conn *websocket.Conn, sub *Sub, room *Room, n *Notebook, isWriter bool) {
	defer func() {
		room.Leave(sub)
		_ = conn.Close()
	}()

	// Send init snapshot.
	viewers, writerName, live := room.snapshot()
	n.Lock()
	init := map[string]any{
		"t":           "init",
		"id":          n.ID,
		"title":       n.Title,
		"author":      n.Author,
		"text":        n.Text,
		"cursor":      room.CurrentCursor(),
		"viewers":     viewers,
		"live":        live,
		"writer_name": writerName,
		"you_write":   isWriter,
	}
	n.Unlock()
	if err := conn.WriteJSON(init); err != nil {
		return
	}

	// Writer goroutine pushes queued messages to the socket.
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case msg, ok := <-sub.Outbox():
				if !ok {
					_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					return
				}
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					return
				}
			case <-ticker.C:
				_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()

	conn.SetReadLimit(1 << 20)
	_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	})

	opsSinceSave := 0
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			break
		}
		_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		var msg inboundMsg
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}
		switch msg.T {
		case "op":
			if !isWriter {
				continue
			}
			op := Op{P: msg.P, N: msg.N, S: msg.S}
			room.ApplyOp(n, op)
			opsSinceSave++
			if opsSinceSave >= 32 {
				n.Lock()
				_ = s.store.Save(n)
				n.Unlock()
				opsSinceSave = 0
			}
		case "cur":
			if !isWriter {
				continue
			}
			room.SetCursor(msg.P)
		case "end":
			return
		}
	}
	<-done
}
