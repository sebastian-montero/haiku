package app

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// Store persists notebooks as JSON files under a data directory. One file per
// notebook keeps writes cheap and avoids rewriting the world on every flush.
type Store struct {
	dir string
	mu  sync.RWMutex
	nbs map[string]*Notebook
}

func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &Store{dir: dir, nbs: map[string]*Notebook{}}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var n Notebook
		if err := json.Unmarshal(b, &n); err != nil {
			continue
		}
		s.nbs[n.ID] = &n
	}
	return s, nil
}

func (s *Store) Get(id string) (*Notebook, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.nbs[id]
	return n, ok
}

func (s *Store) List() []*Notebook {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Notebook, 0, len(s.nbs))
	for _, n := range s.nbs {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt > out[j].UpdatedAt })
	return out
}

func (s *Store) Create(title, author string) *Notebook {
	now := time.Now().UnixMilli()
	n := &Notebook{
		ID:        newID(),
		Title:     title,
		Author:    author,
		CreatedAt: now,
		UpdatedAt: now,
		Ops:       []Op{},
	}
	s.mu.Lock()
	s.nbs[n.ID] = n
	s.mu.Unlock()
	_ = s.Save(n)
	return n
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	delete(s.nbs, id)
	s.mu.Unlock()
	path := filepath.Join(s.dir, id+".json")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Save writes the notebook to disk. Caller is expected to hold the notebook
// lock (or otherwise know no concurrent writers touch it).
func (s *Store) Save(n *Notebook) error {
	b, err := json.Marshal(n)
	if err != nil {
		return err
	}
	path := filepath.Join(s.dir, n.ID+".json")
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func newID() string {
	var b [6]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
