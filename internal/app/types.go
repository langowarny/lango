package app

import (
	"context"

	"github.com/langowarny/lango/internal/agent"
	"github.com/langowarny/lango/internal/config"
	"github.com/langowarny/lango/internal/gateway"
	"github.com/langowarny/lango/internal/security"
	"github.com/langowarny/lango/internal/session"
)

// App is the root application structure
type App struct {
	Config *config.Config

	// Core Components
	Agent          agent.AgentRuntime
	Gateway        *gateway.Server
	CryptoProvider security.CryptoProvider
	Store          session.Store

	// Channels
	Channels []Channel

	// State
	BrowserSessionID string
}

// Channel represents a communication channel (Telegram, Discord, Slack)
type Channel interface {
	Start(ctx context.Context) error
	Stop()
}
