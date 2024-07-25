package notifications

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WsConnectionManager struct {
	wsConnections map[string]*websocket.Conn
	mu            sync.RWMutex
}

func NewWsConnectionManager() *WsConnectionManager {
	return &WsConnectionManager{
		wsConnections: make(map[string]*websocket.Conn),
	}
}
func (wsConn *WsConnectionManager) RegisterWebSocket(id string, conn *websocket.Conn) {
	wsConn.mu.Lock()
	defer wsConn.mu.Unlock()
	wsConn.wsConnections[id] = conn
}

func (wsConn *WsConnectionManager) UnregisterWebSocket(id string) {
	wsConn.mu.Lock()
	defer wsConn.mu.Unlock()
	delete(wsConn.wsConnections, id)
}

func (wsConn *WsConnectionManager) SendMessage(id string, message string) {
	wsConn.mu.RLock()
	conn, ok := wsConn.wsConnections[id]
	if !ok {
		wsConn.mu.RUnlock()
		return
	}
	wsConn.mu.RUnlock()
	conn.WriteMessage(websocket.TextMessage, []byte(message+" at "+time.Now().Format(time.RFC3339)))
	conn.Close()

	wsConn.mu.Lock()
	delete(wsConn.wsConnections, id)
	wsConn.mu.Unlock()
}
