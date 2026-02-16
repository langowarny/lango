package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/langowarny/lango/internal/approval"
)

// approvalPending holds the response channel and message metadata for a pending approval.
type approvalPending struct {
	ch        chan bool
	chatID    int64
	messageID int
}

// ApprovalProvider implements approval.Provider for Telegram using InlineKeyboard buttons.
type ApprovalProvider struct {
	bot     BotAPI
	pending sync.Map // map[requestID]*approvalPending
	timeout time.Duration
}

var _ approval.Provider = (*ApprovalProvider)(nil)

// NewApprovalProvider creates a Telegram approval provider.
func NewApprovalProvider(bot BotAPI, timeout time.Duration) *ApprovalProvider {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &ApprovalProvider{
		bot:     bot,
		timeout: timeout,
	}
}

// RequestApproval sends an InlineKeyboard message to the chat and waits for a callback.
func (p *ApprovalProvider) RequestApproval(ctx context.Context, req approval.ApprovalRequest) (bool, error) {
	chatID, err := parseTelegramChatID(req.SessionKey)
	if err != nil {
		return false, fmt.Errorf("parse session key: %w", err)
	}

	respChan := make(chan bool, 1)

	// Build inline keyboard
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Approve", "approve:"+req.ID),
			tgbotapi.NewInlineKeyboardButtonData("❌ Deny", "deny:"+req.ID),
		),
	)

	text := fmt.Sprintf("🔐 Tool '%s' requires approval", req.ToolName)
	if req.Summary != "" {
		text += "\n\n" + req.Summary
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard

	sentMsg, err := p.bot.Send(msg)
	if err != nil {
		return false, fmt.Errorf("send approval message: %w", err)
	}

	p.pending.Store(req.ID, &approvalPending{
		ch:        respChan,
		chatID:    chatID,
		messageID: sentMsg.MessageID,
	})
	defer p.pending.Delete(req.ID)

	select {
	case approved := <-respChan:
		return approved, nil
	case <-ctx.Done():
		p.editApprovalMessage(chatID, sentMsg.MessageID, "🔐 Tool approval — ⏱ Expired")
		return false, ctx.Err()
	case <-time.After(p.timeout):
		p.editApprovalMessage(chatID, sentMsg.MessageID, "🔐 Tool approval — ⏱ Expired")
		return false, fmt.Errorf("approval timeout")
	}
}

// HandleCallback processes an InlineKeyboard callback query for approval.
func (p *ApprovalProvider) HandleCallback(query *tgbotapi.CallbackQuery) {
	if query == nil || query.Data == "" {
		return
	}

	var requestID string
	var approved bool

	if strings.HasPrefix(query.Data, "approve:") {
		requestID = strings.TrimPrefix(query.Data, "approve:")
		approved = true
	} else if strings.HasPrefix(query.Data, "deny:") {
		requestID = strings.TrimPrefix(query.Data, "deny:")
		approved = false
	} else {
		return
	}

	// LoadAndDelete first to prevent duplicate clicks (TOCTOU)
	val, ok := p.pending.LoadAndDelete(requestID)
	if !ok {
		// Already processed or expired — answer callback silently
		callback := tgbotapi.NewCallback(query.ID, "")
		if _, err := p.bot.Request(callback); err != nil {
			if !isCallbackExpiredErr(err) {
				logger().Debugw("answer expired callback error", "error", err)
			}
		}
		return
	}

	pending, ok := val.(*approvalPending)
	if !ok {
		logger().Warnw("unexpected pending type", "requestId", requestID)
		return
	}

	// Answer callback to dismiss the loading indicator
	callback := tgbotapi.NewCallback(query.ID, "")
	if _, err := p.bot.Request(callback); err != nil {
		if !isCallbackExpiredErr(err) {
			logger().Debugw("answer callback error", "error", err)
		}
	}

	// Edit original message to remove the keyboard
	status := "❌ Denied"
	if approved {
		status = "✅ Approved"
	}
	p.editApprovalMessage(pending.chatID, pending.messageID, fmt.Sprintf("🔐 Tool approval — %s", status))

	// Send result to waiting goroutine
	select {
	case pending.ch <- approved:
	default:
	}
}

// CanHandle returns true for session keys starting with "telegram:".
func (p *ApprovalProvider) CanHandle(sessionKey string) bool {
	return strings.HasPrefix(sessionKey, "telegram:")
}

// editApprovalMessage edits a message with new text and removes inline keyboard buttons.
func (p *ApprovalProvider) editApprovalMessage(chatID int64, messageID int, newText string) {
	emptyMarkup := tgbotapi.NewInlineKeyboardMarkup()
	edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, newText, emptyMarkup)
	if _, err := p.bot.Send(edit); err != nil {
		if !isMessageNotModifiedErr(err) {
			logger().Warnw("edit approval message error", "error", err)
		}
	}
}

func isCallbackExpiredErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "query is too old")
}

func isMessageNotModifiedErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "message is not modified")
}

// parseTelegramChatID extracts the chatID from a session key like "telegram:<chatID>:<userID>".
func parseTelegramChatID(sessionKey string) (int64, error) {
	parts := strings.SplitN(sessionKey, ":", 3)
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid telegram session key: %s", sessionKey)
	}
	return strconv.ParseInt(parts[1], 10, 64)
}
