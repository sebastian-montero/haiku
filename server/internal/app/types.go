package app

import (
	"sync"
	"unicode/utf16"
)

// Op is a single edit at UTF-16 position P: delete N code units, then insert S.
// Insert-only → N=0. Delete-only → S="". Replace → both.
// Ts is milliseconds since the notebook was created.
type Op struct {
	P  int    `json:"p"`
	N  int    `json:"n"`
	S  string `json:"s"`
	Ts int64  `json:"ts"`
}

// Apply returns text with op applied. Positions are JS-string (UTF-16) units.
func Apply(text string, op Op) string {
	units := utf16.Encode([]rune(text))
	p := op.P
	if p < 0 {
		p = 0
	}
	if p > len(units) {
		p = len(units)
	}
	end := p + op.N
	if end > len(units) {
		end = len(units)
	}
	insert := utf16.Encode([]rune(op.S))
	out := make([]uint16, 0, len(units)-(end-p)+len(insert))
	out = append(out, units[:p]...)
	out = append(out, insert...)
	out = append(out, units[end:]...)
	return string(utf16.Decode(out))
}

// Notebook is one piece of writing, live or finished.
type Notebook struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Text      string `json:"text"`
	Ops       []Op   `json:"ops"`

	mu sync.Mutex
}

func (n *Notebook) Lock()   { n.mu.Lock() }
func (n *Notebook) Unlock() { n.mu.Unlock() }

// Summary is the compact view used for listings.
type Summary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
	Preview    string `json:"preview"`
	OpsCount   int    `json:"ops_count"`
	Live       bool   `json:"live"`
	Viewers    int    `json:"viewers"`
	WriterName string `json:"writer_name"`
}
