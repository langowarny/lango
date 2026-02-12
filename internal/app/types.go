package app

import (
	"context"
	"sync"

	"github.com/langowarny/lango/internal/adk"
	"github.com/langowarny/lango/internal/config"
	"github.com/langowarny/lango/internal/gateway"
	"github.com/langowarny/lango/internal/session"
)

// App is the root application structure
type App struct {
	Config *config.Config

	// Core Components
	Agent   *adk.Agent
	Gateway *gateway.Server
	Store   session.Store

	// Channels
	Channels []Channel

	// wg tracks background goroutines for graceful shutdown
	wg sync.WaitGroup
}

// Channel represents a communication channel (Telegram, Discord, Slack)
type Channel interface {
	Start(ctx context.Context) error
	Stop()
}
