package services

import (
	"fmt"
	"time"

	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"github.com/instaagrammeta/crm-real/backend-go/internal/services/telegram"
	"github.com/instaagrammeta/crm-real/backend-go/internal/ws"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Notifier marshalls notifications across DB, WebSocket, Telegram.
type Notifier struct {
	db  *gorm.DB
	hub *ws.Hub
	tg  *telegram.Bot
}

// NewNotifier — конструктор.
func NewNotifier(db *gorm.DB, hub *ws.Hub, tg *telegram.Bot) *Notifier {
	return &Notifier{db: db, hub: hub, tg: tg}
}

// Payload — параметрҳои оғоҳнома.
type Payload struct {
	UserID     uint
	Title      string
	Body       string
	Type       string
	EntityType string
	EntityID   uint
	Link       string
}

// Notify сохраняет в БД + push WS + Telegram.
func (n *Notifier) Notify(p Payload) error {
	if p.UserID == 0 {
		return nil
	}
	rec := models.Notification{
		UserID:     p.UserID,
		Title:      p.Title,
		Body:       p.Body,
		Type:       p.Type,
		EntityType: p.EntityType,
		EntityID:   p.EntityID,
		Link:       p.Link,
	}
	if err := n.db.Create(&rec).Error; err != nil {
		return err
	}

	// WebSocket push
	if n.hub != nil {
		n.hub.SendToUser(p.UserID, ws.Event{
			Type:   ws.EventNotification,
			UserID: p.UserID,
			Data:   rec,
			Time:   time.Now(),
		})
	}

	// Telegram push
	if n.tg != nil && n.tg.Enabled() {
		go n.sendTelegram(p)
	}
	return nil
}

// NotifyMany — оғоҳномаи оммавӣ.
func (n *Notifier) NotifyMany(userIDs []uint, p Payload) {
	for _, uid := range userIDs {
		p.UserID = uid
		if err := n.Notify(p); err != nil {
			log.Error().Err(err).Uint("uid", uid).Msg("notify")
		}
	}
}

// NotifyAdmins — ба тамоми админҳо.
func (n *Notifier) NotifyAdmins(p Payload) {
	var admins []models.User
	if err := n.db.Where("role = ?", "admin").Find(&admins).Error; err != nil {
		return
	}
	ids := make([]uint, 0, len(admins))
	for _, a := range admins {
		ids = append(ids, a.ID)
	}
	n.NotifyMany(ids, p)
}

// Broadcast — оғоҳномаи real-time бе сабт дар БД.
func (n *Notifier) Broadcast(eventType ws.EventType, data any) {
	if n.hub == nil {
		return
	}
	n.hub.Broadcast(ws.Event{Type: eventType, Data: data, Time: time.Now()})
}

func (n *Notifier) sendTelegram(p Payload) {
	subs, err := n.tg.SubscribersForUser(p.UserID)
	if err != nil {
		return
	}
	text := fmt.Sprintf("<b>%s</b>\n%s", p.Title, p.Body)
	if p.Link != "" {
		text += fmt.Sprintf("\n\n🔗 %s", p.Link)
	}
	for _, s := range subs {
		if err := n.tg.SendMessage(s.ChatID, text); err != nil {
			log.Warn().Err(err).Int64("chat", s.ChatID).Msg("telegram send")
		}
	}
}
