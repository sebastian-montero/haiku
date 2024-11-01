package ws_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"wout/internal/models"
	"wout/internal/services"
	"wout/internal/utils/logger"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	Service *services.SessionService
	Clients map[int]map[*websocket.Conn]bool
	Mutex   sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust as needed for security
	},
}

func (h *WebSocketHandler) WebSocketEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notebookIDStr := vars["notebook_id"]
	ownerIDStr := r.URL.Query().Get("owner_id")

	notebookID, err := strconv.Atoi(notebookIDStr)
	if err != nil {
		logger.Error("Invalid notebook_id")
		return
	}
	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		logger.Error("Invalid owner_id")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("WebSocket upgrade error: %v", err))
		return
	}
	defer conn.Close()

	if connType == "read" {
		h.handleRead(conn, notebookID)
	} else if connType == "write" {
		h.handleWrite(conn, notebookID, ownerID)
	}
}

func (h *WebSocketHandler) handleRead(conn *websocket.Conn, notebookID int) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Error(fmt.Sprintf("WebSocket read error: %v", err))
			break
		}

		logger.Info(fmt.Sprintf("Received message on notebook %d: %s", notebookID, string(msg)))
	}
}

func (h *WebSocketHandler) handleWrite(conn *websocket.Conn, notebookID, ownerID int) {
	h.Mutex.Lock()
	session, ok := h.Service.SessionExistsByNotebookID(strconv.Itoa(notebookID))
	if !ok {
		session = models.Session{
			OwnerID:    ownerID,
			NotebookID: notebookID,
		}
		if err := h.Service.CreateSession(&session); err != nil {
			logger.Error(fmt.Sprintf("Failed to create session: %v", err))
			h.Mutex.Unlock()
			return
		}
	}

	err := h.Service.ActivateSession(strconv.Itoa(session.ID))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to activate session: %v", err))
		h.Mutex.Unlock()
		return
	}

	if h.Clients[notebookID] == nil {
		h.Clients[notebookID] = make(map[*websocket.Conn]bool)
	}
	h.Clients[notebookID][conn] = true
	h.Mutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Error(fmt.Sprintf("WebSocket read error: %v", err))
			h.removeClient(conn, notebookID)
			break
		}

		var parsedMsg map[string]interface{}
		if err := json.Unmarshal(msg, &parsedMsg); err != nil {
			logger.Error(fmt.Sprintf("JSON parsing error: %v", err))
			continue
		}

		parsedMsg["owner_id"] = ownerID
		parsedMsg["notebook_id"] = notebookID
		parsedMsg["session_id"] = session.ID

		if parsedMsg["type"] == "end" {
			logger.Info(fmt.Sprintf("Ending session for notebook %d", notebookID))
			break
		}
		if parsedMsg["type"] == "write" {
			logger.Info(fmt.Sprintf("Broadcasting message on notebook %d", notebookID))
			h.broadcastMessage(notebookID, parsedMsg)
		}

	}
}

func (h *WebSocketHandler) removeClient(conn *websocket.Conn, notebookID int) {
	h.Mutex.Lock()
	delete(h.Clients[notebookID], conn)
	if len(h.Clients[notebookID]) == 0 {
		delete(h.Clients, notebookID)
	}
	h.Mutex.Unlock()
}

func (h *WebSocketHandler) broadcastMessage(notebookID int, message map[string]interface{}) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	for client := range h.Clients[notebookID] {
		if err := client.WriteJSON(message); err != nil {
			logger.Error(fmt.Sprintf("WebSocket broadcast error: %v", err))
			client.Close()
			delete(h.Clients[notebookID], client)
		}
	}
}
