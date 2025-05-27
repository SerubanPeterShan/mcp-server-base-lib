# MCP Server Library

A Go library for implementing Machine Control Protocol (MCP) servers. This is a pure library package that can be imported into your Go applications to add MCP server capabilities.

## Features

- WebSocket-based communication
- JSON message format
- Custom message handlers
- Concurrent client handling
- Structured logging
- Health check endpoint
- Easy to extend and customize

## Installation

```bash
go get github.com/serubanpetershan/mcp-server
```

## Usage

```go
package main

import (
    "log"
    "github.com/serubanpetershan/mcp-server/mcp"
    "github.com/sirupsen/logrus"
    "github.com/gorilla/websocket"
)

func main() {
    // Create server configuration
    config := &mcp.Config{
        Port: "8080",
        Logger: logrus.New(),
        Handlers: map[string]mcp.MessageHandler{
            "echo": func(s *mcp.Server, conn *websocket.Conn, msg mcp.Message) error {
                return s.Send(conn, msg)
            },
        },
    }

    // Create and start the server
    server := mcp.NewServer(config)
    if err := server.Start(config.Port); err != nil {
        log.Fatal(err)
    }
}
```

## Message Format

```json
{
    "type": "message_type",
    "payload": {}
}
```

## API Reference

### Server

```go
// Create a new server
server := mcp.NewServer(config)

// Start the server
server.Start(port string) error

// Register a message handler
server.RegisterHandler(messageType string, handler MessageHandler)

// Broadcast a message to all clients
server.Broadcast(msg Message)

// Send a message to a specific client
server.Send(client *websocket.Conn, msg Message) error

// Get all connected clients
server.GetClients() []*websocket.Conn
```

### Configuration

```go
type Config struct {
    Port     string
    Logger   *logrus.Logger
    Handlers map[string]MessageHandler
}
```

### Message Handler

```go
type MessageHandler func(*Server, *websocket.Conn, Message) error
```

## Examples

See the `examples` directory for complete examples of how to use this library in your applications.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

Copyright (c) 2024 Seruban Peter Shan. All rights reserved.

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

The software is protected by copyright law and international treaties. Unauthorized reproduction or distribution of this software, or any portion of it, may result in severe civil and criminal penalties, and will be prosecuted to the maximum extent possible under law. 