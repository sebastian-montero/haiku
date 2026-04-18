# haiku·

A writing app where the process *is* the piece. Open a notebook, start typing,
and whoever is tuned in watches every keystroke — the pauses, the reversals,
the crossed-out words — in real time. When you close the notebook, the whole
session can be replayed at any speed.

No accounts. Pick a pen name and write.

## Features

- **Keystroke-level streaming.** Each edit is diffed into an op (`insert N
  chars at position P` / `delete M chars`) and pushed over a WebSocket.
  Viewers apply ops as they arrive — no full-text snapshots, no diff
  reconstruction on the client.
- **Live caret.** Viewers see where the author is typing, in a blinking
  vermilion caret inside the text itself.
- **Presence.** "3 watching · alice is writing" on every live page.
- **Replay.** Every notebook stores its full op log with timestamps. `replay`
  plays the session back at 1× / 2× / 4× / 8×, with the original timing
  between edits.
- **Autosave.** Ops stream to a JSON file on disk; no database.
- **Runs locally.** Two processes (`go run`, `next dev`), no Docker,
  no Postgres, no auth.

## Quick start

Requirements: Go 1.22+ and Node 18+.

```bash
make install   # one-time: fetch Go and npm deps
make dev       # runs the server on :8080 and the client on :3000
```

Then open http://localhost:3000.

(or run the two halves separately: `make server` and `make client` in
different terminals.)

## Project layout

```
server/
  cmd/server/main.go          entry point
  internal/app/
    types.go                  Op, Notebook, Summary
    store.go                  JSON-per-notebook persistence
    hub.go                    in-memory rooms with pub/sub
    http.go                   REST: list/get/create/patch/delete
    ws.go                     /ws/write/{id}, /ws/watch/{id}

client/
  src/
    pages/
      index.js                home: live notebooks + recent notebooks
      write/[id].js           author view: textarea → diff → op
      watch/[id].js           viewer view (live) + replay view
    components/
      TopBar.js               brand + pen name switcher
      Stage.js                read-only text with embedded caret
      NameGate.js             pen-name prompt
    lib/
      api.js                  fetch + WebSocket helpers
      diff.js                 minimal string-diff → op
      name.js                 localStorage pen name
```

## Protocol

All messages are JSON over WebSocket.

**Viewer receives:**

```jsonc
{ "t": "init",   "text": "…", "cursor": 5, "viewers": 3, "live": true,
  "writer_name": "alice" }
{ "t": "op",     "op": { "p": 5, "n": 0, "s": "h", "ts": 1234 } }
{ "t": "cur",    "p": 12 }
{ "t": "presence", "viewers": 4, "live": true, "writer_name": "alice" }
```

**Writer sends:**

```jsonc
{ "t": "op",  "p": 5, "n": 0, "s": "h" }
{ "t": "cur", "p": 12 }
{ "t": "end" }
```

Positions are JS-style UTF-16 code-unit offsets. One notebook has one writer
at a time; a second writer's claim is rejected.

## Data

Each notebook is one JSON file under `server/data/<id>.json` containing its
ops log and materialized text. `make clean` wipes it.
