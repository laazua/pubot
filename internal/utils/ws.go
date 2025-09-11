package utils

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type TaskStatusEnum string

const (
	TaskRunning TaskStatusEnum = "running"
	TaskSuccess TaskStatusEnum = "success"
	TaskError   TaskStatusEnum = "error"
	TaskStopped TaskStatusEnum = "stopped" // 可选，和 success 区分
)

type TaskStatus struct {
	ID     uint           `json:"id"`
	Status TaskStatusEnum `json:"status"`
	Count  int            `json:"count"`
}

type Hub struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) AddClient(conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()
}

func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, conn)
	h.mu.Unlock()
}

func (h *Hub) Broadcast(status TaskStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for client := range h.clients {
		if err := client.WriteJSON(status); err != nil {
			slog.Error("ws 发送出错", slog.Any("Err", err.Error()))
			client.Close()
			delete(h.clients, client)
		}
		slog.Info("广播任务状态", slog.Any("status", status), slog.Int("client_count", len(h.clients)))
	}

}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	// 取出客户端请求的 subprotocols（里面就是 token）
	protocols := websocket.Subprotocols(r)
	if len(protocols) == 0 {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	upgrader := websocket.Upgrader{
		CheckOrigin:  func(r *http.Request) bool { return true },
		Subprotocols: []string{protocols[0]}, // 只确认这个 token
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("ws 升级失败", slog.Any("Err", err.Error()))
		return
	}
	slog.Info("WS 连接成功", slog.String("RemoteAddr", r.RemoteAddr))
	h.AddClient(conn)
	defer h.RemoveClient(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
