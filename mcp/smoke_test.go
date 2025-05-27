package mcp

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSmokeBasicFunctionality(t *testing.T) {
	// Basic smoke test to verify core functionality
	server := NewServer(&Config{Logger: logrus.New()})
	defer server.Shutdown()

	// Start server
	go func() {
		err := server.Start("8083")
		require.NoError(t, err)
	}()
	// Wait for server to be ready
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	var resp *http.Response
	var err error
waitLoop:
	for {
		select {
		case <-timeout:
			t.Fatalf("Server did not become ready in time")
		case <-ticker.C:
			resp, err = http.Get("http://localhost:8083/health")
			if err == nil && resp.StatusCode == http.StatusOK {
				break waitLoop
			}
		}
	}
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8083/ws", nil)
	require.NoError(t, err)
	defer conn.Close()

	// Test message sending
	msg := Message{
		Type:    "test",
		Payload: json.RawMessage(`{"data": "smoke test"}`),
	}
	err = server.Send(conn, msg)
	require.NoError(t, err)

	// Verify message received
	var receivedMsg Message
	err = conn.ReadJSON(&receivedMsg)
	require.NoError(t, err)
	assert.Equal(t, msg.Type, receivedMsg.Type)
	assert.Equal(t, msg.Payload, receivedMsg.Payload)
}

func TestSmokeServerStartup(t *testing.T) {
	// Test server startup with different configurations
	testCases := []struct {
		name   string
		config *Config
		port   string
	}{
		{
			name: "Default Config",
			config: &Config{
				Logger: logrus.New(),
			},
			port: "8084",
		},
		{
			name: "With Handlers",
			config: &Config{
				Logger: logrus.New(),
				Handlers: map[string]MessageHandler{
					"echo": func(s *Server, conn *websocket.Conn, msg Message) error {
						return s.Send(conn, msg)
					},
				},
			},
			port: "8085",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := NewServer(tc.config)
			go func() {
				err := server.Start(tc.port)
				require.NoError(t, err)
			}()
			// Wait for server to be ready
			timeout := time.After(5 * time.Second)
			tick := time.Tick(100 * time.Millisecond)
			var resp *http.Response
			var err error
		waitLoop:
			for {
				select {
				case <-timeout:
					t.Fatalf("Server did not become ready in time")
				case <-tick:
					resp, err = http.Get("http://localhost:" + tc.port + "/health")
					if err == nil && resp.StatusCode == http.StatusOK {
						break waitLoop
					}
				}
			}
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestSmokeMessageHandling(t *testing.T) {
	// Test basic message handling functionality
	server := NewServer(&Config{
		Logger: logrus.New(),
		Handlers: map[string]MessageHandler{
			"echo": func(s *Server, conn *websocket.Conn, msg Message) error {
				return s.Send(conn, msg)
			},
		},
	})

	go server.Start("8086")
	// Wait for server to be ready
	timeout := time.After(5 * time.Second)
	tick := time.Tick(100 * time.Millisecond)
	var resp *http.Response
	var err error
waitLoop:
	for {
		select {
		case <-timeout:
			t.Fatalf("Server did not become ready in time")
		case <-tick:
			resp, err = http.Get("http://localhost:8086/health")
			if err == nil && resp.StatusCode == http.StatusOK {
				break waitLoop
			}
		}
	}
	require.NoError(t, err)

	// Connect client
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8086/ws", nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send message
	msg := Message{
		Type:    "echo",
		Payload: json.RawMessage(`{"data": "echo test"}`),
	}
	err = conn.WriteJSON(msg)
	require.NoError(t, err)

	// Verify echo response
	var response Message
	err = conn.ReadJSON(&response)
	require.NoError(t, err)
	assert.Equal(t, msg.Type, response.Type)
	assert.Equal(t, msg.Payload, response.Payload)
}
