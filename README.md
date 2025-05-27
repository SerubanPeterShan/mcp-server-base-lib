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
go get github.com/serubanpetershan/mcp-server-base-library
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

## Security

This project uses automated security scanning:
- SonarQube for code quality and security analysis
- Snyk for dependency vulnerability scanning
- Weekly automated security scans
- All pull requests are automatically scanned

## Versioning

This project follows [Semantic Versioning](https://semver.org/). For the versions available, see the [tags on this repository](https://github.com/serubanpetershan/mcp-server-base-library/tags).

To create a new release:
1. Update the VERSION file
2. Create a new tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
3. Push the tag: `git push origin v1.0.0`

The release workflow will automatically:
- Run tests
- Build the project
- Create a GitHub release
- Generate a changelog

## License

Copyright (c) 2024 Seruban Peter Shan. All rights reserved.

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

The software is protected by copyright law and international treaties. Unauthorized reproduction or distribution of this software, or any portion of it, may result in severe civil and criminal penalties, and will be prosecuted to the maximum extent possible under law. 