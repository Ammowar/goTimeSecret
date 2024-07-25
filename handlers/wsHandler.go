package handlers

import (
	"net/http"

	"github.com/Ammowar/goTimeSecret/notifications"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader      *websocket.Upgrader
	wsConnManager *notifications.WsConnectionManager
}

func NewWebSocketHandler(wsConnManager *notifications.WsConnectionManager) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		wsConnManager: wsConnManager}
}

func (ws *WebSocketHandler) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to set websocket upgrade", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	id := r.URL.Query().Get("id")
	if id == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Missing secret ID"))
		return
	}

	ws.wsConnManager.RegisterWebSocket(id, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			ws.wsConnManager.UnregisterWebSocket(id)
			break
		}
	}
}
