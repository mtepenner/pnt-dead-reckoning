package telemetry

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/mtepenner/pnt-dead-reckoning/navigation_core/internal/filter"
)

type Snapshot struct {
	State   filter.StateVector    `json:"state"`
	History []filter.HistoryPoint `json:"history"`
}

type Hub struct {
	mu       sync.RWMutex
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]struct{}
	latest   Snapshot
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]struct{}),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(*http.Request) bool { return true },
		},
	}
}

func (hub *Hub) Publish(snapshot Snapshot) {
	hub.mu.Lock()
	hub.latest = snapshot
	clients := make([]*websocket.Conn, 0, len(hub.clients))
	for client := range hub.clients {
		clients = append(clients, client)
	}
	hub.mu.Unlock()

	payload, err := json.Marshal(snapshot)
	if err != nil {
		return
	}

	for _, client := range clients {
		_ = client.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
		if err := client.WriteMessage(websocket.TextMessage, payload); err != nil {
			hub.removeClient(client)
		}
	}
}

func (hub *Hub) Latest() Snapshot {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	return hub.latest
}

func (hub *Hub) HandleState(writer http.ResponseWriter, _ *http.Request) {
	respondJSON(writer, http.StatusOK, hub.Latest().State)
}

func (hub *Hub) HandleHistory(writer http.ResponseWriter, _ *http.Request) {
	respondJSON(writer, http.StatusOK, hub.Latest().History)
}

func (hub *Hub) HandleWebSocket(writer http.ResponseWriter, request *http.Request) {
	conn, err := hub.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return
	}

	hub.mu.Lock()
	hub.clients[conn] = struct{}{}
	snapshot := hub.latest
	hub.mu.Unlock()

	if err := conn.WriteJSON(snapshot); err != nil {
		hub.removeClient(conn)
		return
	}

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			hub.removeClient(conn)
			return
		}
	}
}

func respondJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

func (hub *Hub) removeClient(conn *websocket.Conn) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	delete(hub.clients, conn)
	_ = conn.Close()
}
