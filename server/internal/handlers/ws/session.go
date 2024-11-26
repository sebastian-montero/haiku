package ws_handler

import (
	"encoding/json"
	"fmt"
	"haiku/internal/models"
	"haiku/internal/services"
	"haiku/internal/utils/logger"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	Service        *services.SessionService
	Clients        map[int]map[*websocket.Conn]bool
	SessionContent map[int]string
	Mutex          sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebSocketHandler) WebSocketEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notebookIDStr := vars["notebook_id"]
	connType := vars["notebook_id"]
	ownerIDStr := r.URL.Query().Get("owner_id")

	notebookID, _ := strconv.Atoi(notebookIDStr)
	ownerID, _ := strconv.Atoi(ownerIDStr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("WebSocket upgrade error: %v", err))
		return
	}
	defer conn.Close()

	h.Mutex.Lock()
	if h.Clients[notebookID] == nil {
		h.Clients[notebookID] = make(map[*websocket.Conn]bool)
	}
	h.Clients[notebookID][conn] = true
	logger.Info(fmt.Sprintf("New connection created for notebook %d (Owner ID: %d). Total connections: %d", notebookID, ownerID, len(h.Clients[notebookID])))
	h.Mutex.Unlock()

	if connType == "read" {
		h.handleRead(conn, notebookID)
	} else if connType == "write" {
		h.handleWrite(conn, notebookID, ownerID)
	}
}

func (h *WebSocketHandler) handleRead(conn *websocket.Conn, notebookID int) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			logger.Error(fmt.Sprintf("WebSocket read error: %v", err))
			h.removeClient(conn, notebookID)
			break
		}
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
	h.Mutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Error(fmt.Sprintf("WebSocket read error: %v", err))
			h.removeClient(conn, notebookID)
			h.removeAllClientsFromNotebook(notebookID)
			break
		}

		var parsedMsg map[string]interface{}
		if err := json.Unmarshal(msg, &parsedMsg); err != nil {
			logger.Error(fmt.Sprintf("JSON parsing error: %v", err))
			continue
		}

		// Handle "send" and "end" message types
		messageType, ok := parsedMsg["type"].(string)
		if !ok {
			logger.Error("Message type not specified or invalid")
			continue
		}

		if messageType == "end" {
			logger.Info(fmt.Sprintf("Ending session for notebook %d", notebookID))
			h.removeAllClientsFromNotebook(notebookID)
			err := h.Service.EndSessionByID(strconv.Itoa(session.ID), ownerID)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to end session: %v", err))
			}
			err = h.Service.CreateContent(notebookID, h.SessionContent[notebookID])
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to create notebook content: %v", err))
			}
			err = h.Service.UpdateNotebookContent(strconv.Itoa(notebookID), h.SessionContent[notebookID])
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to update notebook content: %v", err))
			}

			break
		}

		if messageType == "send" {
			content, ok := parsedMsg["content"].(string)
			if !ok {
				logger.Error("Content field missing or invalid in 'send' message")
				continue
			}

			logger.Info(fmt.Sprintf("Broadcasting 'send' message on notebook %d: %s", notebookID, content))

			responseMsg := map[string]interface{}{
				"status":      "success",
				"type":        "send",
				"content":     content,
				"owner_id":    ownerID,
				"notebook_id": notebookID,
				"session_id":  session.ID,
			}

			err = conn.WriteJSON(responseMsg)
			if err != nil {
				logger.Error(fmt.Sprintf("WebSocket write error: %v", err))
				h.removeClient(conn, notebookID)
				break
			}

			h.SessionContent[notebookID] = content
			h.broadcastMessage(notebookID, websocket.TextMessage, msg)
		}
	}
}

func (h *WebSocketHandler) removeClient(conn *websocket.Conn, notebookID int) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	if _, exists := h.Clients[notebookID][conn]; exists {
		delete(h.Clients[notebookID], conn)
		logger.Info(fmt.Sprintf("Connection removed for notebook %d. Remaining connections: %d", notebookID, len(h.Clients[notebookID])))
	}

	if len(h.Clients[notebookID]) == 0 {
		delete(h.Clients, notebookID)
		logger.Info(fmt.Sprintf("All connections closed for notebook %d. Notebook entry removed.", notebookID))
	}
}

func (h *WebSocketHandler) broadcastMessage(notebookID int, messageType int, message []byte) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	for client := range h.Clients[notebookID] {
		if err := client.WriteMessage(messageType, message); err != nil {
			logger.Error(fmt.Sprintf("WebSocket broadcast error to client: %v", err))
			client.Close()
			delete(h.Clients[notebookID], client)
			logger.Info(fmt.Sprintf("Faulty connection removed for notebook %d. Remaining connections: %d", notebookID, len(h.Clients[notebookID])))
		}
	}
}

func (h *WebSocketHandler) removeAllClientsFromNotebook(notebookID int) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	connections, exists := h.Clients[notebookID]
	if !exists {
		logger.Info(fmt.Sprintf("No connections found for notebook %d", notebookID))
		return
	}

	for client := range connections {
		// Close the WebSocket connection
		err := client.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("Error closing connection for notebook %d: %v", notebookID, err))
		} else {
			logger.Info(fmt.Sprintf("Connection closed for notebook %d", notebookID))
		}
		// Remove the client from the notebook's map
		delete(connections, client)
	}

	// Remove the notebook entry if there are no more connections
	if len(connections) == 0 {
		delete(h.Clients, notebookID)
		logger.Info(fmt.Sprintf("All connections removed for notebook %d", notebookID))
	}

	logger.Info(fmt.Sprintf("All clients have been removed and connections closed for notebook %d", notebookID))
}
