# Trade Republic Portfolio Downloader v2

This is version 2 of the Trade Republic Portfolio Downloader, which includes a completely rewritten WebSocket client implementation.

## Features

- New WebSocket client with improved connection handling
- Support for pagination in timeline transactions
- Support for timeline detail requests
- Automatic reconnection if the connection is lost
- Automatic unsubscription after receiving data
- Connection reuse for multiple subscriptions

## Directory Structure

```
v2/
├── cmd/
│   └── websocket-downloader/  # Command-line tool for downloading transactions
├── internal/
│   ├── const.go               # Constants used by the WebSocket client
│   └── traderepublic/
│       └── api/
│           └── websocketclient/ # WebSocket client implementation
└── go.mod                     # Module definition for v2
```

## Usage

### WebSocket Client

```go
import "github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/websocketclient"

// Create a new client
client, err := websocketclient.NewClient(
    websocketclient.WithLogger(logger),
    websocketclient.WithSessionToken(sessionToken),
)

// Connect to the WebSocket server
if err := client.Connect(ctx); err != nil {
    log.Fatalf("Failed to connect to WebSocket: %v", err)
}
defer client.Close()

// Subscribe to timeline transactions
transactionsCh, err := client.SubscribeToTimelineTransactions(ctx)
if err != nil {
    log.Fatalf("Failed to subscribe to timeline transactions: %v", err)
}

// Wait for data
select {
case data := <-transactionsCh:
    // Process data
case <-ctx.Done():
    log.Fatalf("Timeout waiting for timeline transactions: %v", ctx.Err())
}
```

### WebSocket Downloader

```bash
# Run with default settings
go run ./v2/cmd/websocket-downloader/main.go

# Run with debug logging
go run ./v2/cmd/websocket-downloader/main.go --debug

# Download only the first 10 transactions
go run ./v2/cmd/websocket-downloader/main.go --max-items=10

# Set a longer timeout (2 minutes)
go run ./v2/cmd/websocket-downloader/main.go --timeout=120
```

## Building

```bash
# Build the WebSocket downloader
go build -o websocket-downloader ./v2/cmd/websocket-downloader
```

## License

Same as the main project.
