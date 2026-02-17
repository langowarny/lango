package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/langowarny/lango/internal/channels/discord"
	"github.com/langowarny/lango/internal/channels/slack"
	"github.com/langowarny/lango/internal/channels/telegram"
)

// channelSender implements cron.ChannelSender, background.ChannelNotifier, and
// workflow.ChannelSender by dispatching messages to the configured channel adapters.
type channelSender struct {
	app *App
}

func newChannelSender(a *App) *channelSender {
	return &channelSender{app: a}
}

// SendMessage sends a text message to the specified channel type.
// For telegram: sends to the first allowlisted chat ID.
// For discord/slack: sends to the general/first available channel.
func (s *channelSender) SendMessage(_ context.Context, channel, message string) error {
	ch := strings.ToLower(strings.TrimSpace(channel))

	for _, c := range s.app.Channels {
		switch ch {
		case "telegram":
			if tg, ok := c.(*telegram.Channel); ok {
				// Use the first allowlisted chat ID for delivery.
				chatID := s.firstTelegramChatID()
				if chatID == 0 {
					return fmt.Errorf("telegram delivery requires at least one allowlisted chat ID")
				}
				return tg.Send(chatID, &telegram.OutgoingMessage{Text: message})
			}
		case "discord":
			if dc, ok := c.(*discord.Channel); ok {
				// Empty channelID lets the adapter use a default channel.
				return dc.Send("", &discord.OutgoingMessage{Content: message})
			}
		case "slack":
			if sl, ok := c.(*slack.Channel); ok {
				return sl.Send("", &slack.OutgoingMessage{Text: message})
			}
		}
	}

	return fmt.Errorf("channel %q not available", channel)
}

func (s *channelSender) firstTelegramChatID() int64 {
	if len(s.app.Config.Channels.Telegram.Allowlist) > 0 {
		return s.app.Config.Channels.Telegram.Allowlist[0]
	}
	return 0
}
