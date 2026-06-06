package telegram

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/instaagrammeta/crm-real/backend-go/internal/config"
	"github.com/instaagrammeta/crm-real/backend-go/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Bot — Telegram bot wrapper.
type Bot struct {
	cfg     *config.Config
	db      *gorm.DB
	api     *tgbotapi.BotAPI
	enabled bool

	stopCh chan struct{}
	once   sync.Once
}

func NewBot(cfg *config.Config, db *gorm.DB) *Bot {
	return &Bot{cfg: cfg, db: db, stopCh: make(chan struct{})}
}

// Start initializes the bot if a token is configured.
// Falls back to no-op if no token.
func (b *Bot) Start() error {
	if b.cfg.TelegramBotToken == "" {
		log.Info().Msg("telegram bot token not set – integration disabled")
		return nil
	}
	api, err := tgbotapi.NewBotAPI(b.cfg.TelegramBotToken)
	if err != nil {
		return fmt.Errorf("telegram bot init: %w", err)
	}
	api.Debug = b.cfg.TelegramDebug
	b.api = api
	b.enabled = true
	log.Info().Str("username", api.Self.UserName).Msg("telegram bot connected")

	go b.poll()
	return nil
}

func (b *Bot) Stop() {
	b.once.Do(func() { close(b.stopCh) })
}

// Enabled reports whether the bot is configured.
func (b *Bot) Enabled() bool { return b.enabled }

// SendMessage отправляет text message в чат.
func (b *Bot) SendMessage(chatID int64, text string) error {
	if !b.enabled {
		return nil
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true
	_, err := b.api.Send(msg)
	return err
}

// GenerateLinkToken creates a one-time token to bind a Telegram chat to a user.
func (b *Bot) GenerateLinkToken(db *gorm.DB, userID uint) (string, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	token := hex.EncodeToString(buf)
	sub := models.TelegramSubscriber{
		UserID:    userID,
		LinkToken: token,
		IsActive:  false,
	}
	if err := db.Create(&sub).Error; err != nil {
		return "", err
	}
	return token, nil
}

// LinkURL — URL барои pуйвастан /start <token>.
func (b *Bot) LinkURL(token string) string {
	if !b.enabled {
		return ""
	}
	return fmt.Sprintf("https://t.me/%s?start=%s", b.api.Self.UserName, token)
}

// SubscribersForUser возвращает активные telegram chat-ҳои корбар.
func (b *Bot) SubscribersForUser(userID uint) ([]models.TelegramSubscriber, error) {
	if !b.enabled {
		return nil, nil
	}
	var subs []models.TelegramSubscriber
	err := b.db.Where("user_id = ? AND is_active = TRUE", userID).Find(&subs).Error
	return subs, err
}

func (b *Bot) poll() {
	cfg := tgbotapi.NewUpdate(0)
	cfg.Timeout = 30
	updates := b.api.GetUpdatesChan(cfg)

	for {
		select {
		case <-b.stopCh:
			b.api.StopReceivingUpdates()
			return
		case upd, ok := <-updates:
			if !ok {
				return
			}
			b.handleUpdate(upd)
		}
	}
}

func (b *Bot) handleUpdate(upd tgbotapi.Update) {
	if upd.Message == nil {
		return
	}
	msg := upd.Message
	text := strings.TrimSpace(msg.Text)

	switch {
	case strings.HasPrefix(text, "/start"):
		b.handleStart(msg, text)
	case text == "/me":
		b.handleMe(msg)
	case text == "/stop":
		b.handleStop(msg)
	case text == "/help":
		_ = b.SendMessage(msg.Chat.ID,
			"Дастурҳо:\n"+
				"<b>/start &lt;token&gt;</b> — пайваст кардани ҳисоб\n"+
				"<b>/me</b> — ҳолати ҳисоби ҷорӣ\n"+
				"<b>/stop</b> — қатъ кардани оғоҳномаҳо\n"+
				"<b>/help</b> — ин рӯйхат")
	default:
		_ = b.SendMessage(msg.Chat.ID,
			"Барои оғоз /start <code>token</code> -ро аз кабинети CRM нусха гиред.")
	}
}

func (b *Bot) handleStart(msg *tgbotapi.Message, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		_ = b.SendMessage(msg.Chat.ID,
			"Барои pуйвастан: <code>/start TOKEN</code>\n"+
				"TOKEN-ро аз CRM (бахши Профил → Telegram) нусха гиред.")
		return
	}
	token := parts[1]
	var sub models.TelegramSubscriber
	if err := b.db.Where("link_token = ?", token).First(&sub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = b.SendMessage(msg.Chat.ID, "❌ Токен топ нашуд ё аллакай истифода шуд.")
			return
		}
		_ = b.SendMessage(msg.Chat.ID, "Хатои дохилӣ. Кӯшиш кунед боз як бор.")
		return
	}

	updates := map[string]any{
		"chat_id":    msg.Chat.ID,
		"username":   msg.From.UserName,
		"first_name": msg.From.FirstName,
		"last_name":  msg.From.LastName,
		"is_active":  true,
		"link_token": "",
	}
	if err := b.db.Model(&sub).Updates(updates).Error; err != nil {
		_ = b.SendMessage(msg.Chat.ID, "Хатои сабт. Бо мудир тамос гиред.")
		return
	}
	// синхронизатсияи user.telegram_chat_id
	_ = b.db.Model(&models.User{}).Where("id = ?", sub.UserID).
		Update("telegram_chat_id", msg.Chat.ID).Error

	_ = b.SendMessage(msg.Chat.ID,
		"✅ Шумо муваффақона пайваст шудед!\nҲоло шумо оғоҳномаҳои CRM-ро мегиред.")
}

func (b *Bot) handleMe(msg *tgbotapi.Message) {
	var sub models.TelegramSubscriber
	if err := b.db.Where("chat_id = ? AND is_active = TRUE", msg.Chat.ID).First(&sub).Error; err != nil {
		_ = b.SendMessage(msg.Chat.ID, "Шумо ҳоло пайваст нашудаед. /start TOKEN -ро истифода кунед.")
		return
	}
	var u models.User
	if err := b.db.First(&u, sub.UserID).Error; err != nil {
		_ = b.SendMessage(msg.Chat.ID, "Корбар топ нашуд.")
		return
	}
	_ = b.SendMessage(msg.Chat.ID, fmt.Sprintf(
		"<b>%s</b>\nЛогин: %s\nНақш: %s\nПайваст шуд: %s",
		u.FullName, u.Login, u.Role, sub.CreatedAt.Format(time.RFC822),
	))
}

func (b *Bot) handleStop(msg *tgbotapi.Message) {
	res := b.db.Model(&models.TelegramSubscriber{}).
		Where("chat_id = ?", msg.Chat.ID).
		Update("is_active", false)
	if res.RowsAffected > 0 {
		_ = b.SendMessage(msg.Chat.ID, "🔕 Оғоҳномаҳо қатъ карда шуданд. /start TOKEN барои аз нав фаъол сохтан.")
	} else {
		_ = b.SendMessage(msg.Chat.ID, "Шумо нашуда будед.")
	}
}
