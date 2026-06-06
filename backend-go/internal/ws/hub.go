package ws

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/instaagrammeta/crm-real/backend-go/internal/auth"
	"github.com/rs/zerolog/log"
)

// EventType — намуди event-и WebSocket.
type EventType string

const (
	EventChatMessage    EventType = "chat:message"
	EventChatDelete     EventType = "chat:delete"
	EventChatUpdate     EventType = "chat:update"
	EventNotification   EventType = "notification"
	EventLeadCreated    EventType = "lead:created"
	EventLeadMoved      EventType = "lead:moved"
	EventLeadUpdated    EventType = "lead:updated"
	EventTaskAssigned   EventType = "task:assigned"
	EventRequestCreated EventType = "request:created"
	EventRequestMoved   EventType = "request:moved"
	EventTariffExpiring EventType = "tariff:expiring"
	EventUserOnline     EventType = "user:online"
)

// Event — паёми WebSocket.
type Event struct {
	Type   EventType `json:"type"`
	UserID uint      `json:"user_id,omitempty"`
	Data   any       `json:"data"`
	Time   time.Time `json:"time"`
}

// Client — як пайвасткунӣ.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	UserID uint
	Name   string
	Role   string
}

// Hub — мутамаркази hama пайвандҳо.
type Hub struct {
	mu         sync.RWMutex
	clients    map[*Client]struct{}
	byUser     map[uint]map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		byUser:     make(map[uint]map[*Client]struct{}),
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
		broadcast:  make(chan []byte, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = struct{}{}
			if _, ok := h.byUser[c.UserID]; !ok {
				h.byUser[c.UserID] = make(map[*Client]struct{})
			}
			h.byUser[c.UserID][c] = struct{}{}
			h.mu.Unlock()
			log.Info().Uint("uid", c.UserID).Msg("ws connected")

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
				if set, ok := h.byUser[c.UserID]; ok {
					delete(set, c)
					if len(set) == 0 {
						delete(h.byUser, c.UserID)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast — паём ба ҳама.
func (h *Hub) Broadcast(ev Event) {
	if ev.Time.IsZero() {
		ev.Time = time.Now()
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return
	}
	select {
	case h.broadcast <- b:
	default:
		log.Warn().Msg("ws broadcast queue full")
	}
}

// SendToUser — паём ба корбари мушаххас.
func (h *Hub) SendToUser(uid uint, ev Event) {
	if ev.Time.IsZero() {
		ev.Time = time.Now()
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.byUser[uid] {
		select {
		case c.send <- b:
		default:
		}
	}
}

// OnlineUsers возвращает ID-и online корбарон.
func (h *Hub) OnlineUsers() []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]uint, 0, len(h.byUser))
	for uid := range h.byUser {
		out = append(out, uid)
	}
	return out
}

// ====================================================================
// HTTP -> WS upgrade
// ====================================================================

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// HandleWS — Gin-handler барои pуйвастани ws://host/ws?token=xxx.
func (h *Hub) HandleWS(jwtMgr *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			if v, err := c.Cookie("auth_token"); err == nil {
				token = v
			}
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token"})
			return
		}
		claims, err := jwtMgr.Parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Error().Err(err).Msg("ws upgrade")
			return
		}
		client := &Client{
			hub:    h,
			conn:   conn,
			send:   make(chan []byte, 64),
			UserID: claims.UserID,
			Name:   claims.FullName,
			Role:   claims.Role,
		}
		h.register <- client
		go client.writePump()
		go client.readPump()
	}
}

// ====================================================================
// Read / write pumps
// ====================================================================

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1 << 20 // 1 MB
)

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		// Мо клиентҳои гузаронданро (typing/ping)-ро дастгирӣ карда метавонем,
		// аммо барои оғоз танҳо мехонем ва дур меандозем.
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
