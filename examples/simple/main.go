package main

import (
	"log"
	"os"

	"github.com/serubanpetershan/mcp-server/mcp"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

func main() {
	// Create a custom logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create message handlers
	handlers := map[string]mcp.MessageHandler{
		"echo": func(s *mcp.Server, conn *websocket.Conn, msg mcp.Message) error {
			// Echo the message back to the sender
			return s.Send(conn, msg)
		},
		"broadcast": func(s *mcp.Server, conn *websocket.Conn, msg mcp.Message) error {
			// Broadcast the message to all clients
			s.Broadcast(msg)
			return nil
		},
	}

	// Create server configuration
	config := &mcp.Config{
		Port:     "8080",
		Logger:   logger,
		Handlers: handlers,
	}

	// Create and start the server
	server := mcp.NewServer(config)
	log.Printf("Starting MCP server on port %s", config.Port)
	if err := server.Start(config.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
