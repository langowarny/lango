package gateway

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/langowarny/lango/internal/agent"
	"github.com/langowarny/lango/internal/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock runtime
type mockRuntime struct {
	agent.AgentRuntime
}

func (m *mockRuntime) Run(ctx context.Context, sessionKey string, input string, events chan<- agent.StreamEvent) error {
	return nil
}

func TestCompanionIntegration(t *testing.T) {
	// Setup server
	cfg := Config{
		Host:             "localhost",
		Port:             0, // Random port
		HTTPEnabled:      true,
		WebSocketEnabled: true,
	}

	gateway := New(cfg, &mockRuntime{}, nil, &session.EntStore{}, nil)

	// Start test server
	server := httptest.NewServer(gateway.router)
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/companion"

	// Connect companion
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// 1. Send companion.hello
	helloReq := RPCRequest{
		ID:     "1",
		Method: "companion.hello",
		Params: json.RawMessage(`{"deviceId": "test-device", "publicKey": "dummy-key"}`),
	}
	err = conn.WriteJSON(helloReq)
	require.NoError(t, err)

	// Read response
	var resp RPCResponse
	err = conn.ReadJSON(&resp)
	require.NoError(t, err)
	assert.Equal(t, "1", resp.ID)
	assert.Nil(t, resp.Error)

	// 2. Test RequestApproval
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	approvalResult := make(chan bool)
	go func() {
		approved, err := gateway.RequestApproval(ctx, "Test Approval")
		if err != nil {
			t.Errorf("RequestApproval failed: %v", err)
			approvalResult <- false
		} else {
			approvalResult <- approved
		}
	}()

	// Read approval request from WS
	var msg map[string]interface{}
	err = conn.ReadJSON(&msg)
	require.NoError(t, err)
	assert.Equal(t, "event", msg["type"])
	assert.Equal(t, "approval.request", msg["event"])

	payload, ok := msg["payload"].(map[string]interface{})
	require.True(t, ok)
	reqID := payload["id"].(string)

	// Send approval.response
	approveReq := RPCRequest{
		ID:     "2",
		Method: "approval.response",
		Params: json.RawMessage(`{"requestId": "` + reqID + `", "approved": true}`),
	}
	err = conn.WriteJSON(approveReq)
	require.NoError(t, err)

	// Verify result
	approved := <-approvalResult
	assert.True(t, approved)
}
