package app

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Server bundles the store + hub for HTTP handlers.
type Server struct {
	store *Store
	hub   *Hub
}

func NewServer(store *Store, hub *Hub) *Server {
	return &Server{store: store, hub: hub}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /notebooks", s.listNotebooks)
	mux.HandleFunc("POST /notebooks", s.createNotebook)
	mux.HandleFunc("GET /notebooks/{id}", s.getNotebook)
	mux.HandleFunc("PATCH /notebooks/{id}", s.patchNotebook)
	mux.HandleFunc("DELETE /notebooks/{id}", s.deleteNotebook)
	mux.HandleFunc("GET /ws/write/{id}", s.WSWrite)
	mux.HandleFunc("GET /ws/watch/{id}", s.WSWatch)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})
	return withCORS(mux)
}

func (s *Server) listNotebooks(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter") // "live" | "" (all)
	live := s.hub.LiveIDs()
	all := s.store.List()
	out := make([]Summary, 0, len(all))
	for _, n := range all {
		isLive := live[n.ID]
		if filter == "live" && !isLive {
			continue
		}
		viewers, writerName, _ := s.hub.Snapshot(n.ID)
		out = append(out, Summary{
			ID:         n.ID,
			Title:      n.Title,
			Author:     n.Author,
			CreatedAt:  n.CreatedAt,
			UpdatedAt:  n.UpdatedAt,
			Preview:    preview(n.Text, 140),
			OpsCount:   len(n.Ops),
			Live:       isLive,
			Viewers:    viewers,
			WriterName: writerName,
		})
	}
	writeJSON(w, 200, out)
}

func (s *Server) createNotebook(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	body.Title = strings.TrimSpace(body.Title)
	body.Author = strings.TrimSpace(body.Author)
	if body.Title == "" {
		body.Title = "untitled"
	}
	if body.Author == "" {
		body.Author = "anon"
	}
	n := s.store.Create(body.Title, body.Author)
	writeJSON(w, 201, n)
}

func (s *Server) getNotebook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, ok := s.store.Get(id)
	if !ok {
		http.Error(w, "not found", 404)
		return
	}
	viewers, writerName, live := s.hub.Snapshot(id)
	n.Lock()
	resp := map[string]any{
		"id":          n.ID,
		"title":       n.Title,
		"author":      n.Author,
		"created_at":  n.CreatedAt,
		"updated_at":  n.UpdatedAt,
		"text":        n.Text,
		"ops":         n.Ops,
		"live":        live,
		"viewers":     viewers,
		"writer_name": writerName,
	}
	n.Unlock()
	writeJSON(w, 200, resp)
}

func (s *Server) patchNotebook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, ok := s.store.Get(id)
	if !ok {
		http.Error(w, "not found", 404)
		return
	}
	var body struct {
		Title *string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	n.Lock()
	if body.Title != nil {
		t := strings.TrimSpace(*body.Title)
		if t != "" {
			n.Title = t
		}
	}
	_ = s.store.Save(n)
	n.Unlock()
	writeJSON(w, 200, n)
}

func (s *Server) deleteNotebook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.store.Delete(id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func preview(text string, max int) string {
	text = strings.TrimSpace(text)
	if len(text) <= max {
		return text
	}
	return text[:max] + "…"
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
