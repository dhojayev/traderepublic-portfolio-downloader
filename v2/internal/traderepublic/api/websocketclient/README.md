# WebSocket Client for Trade Republic API

This package provides a WebSocket client for the Trade Republic API. It allows you to connect to the WebSocket server and subscribe to various data types.

## Features

- Connect to the Trade Republic WebSocket server
- Subscribe to timeline transactions data
- Subscribe to timeline detail data
- Automatic unsubscription after receiving data
- Automatic reconnection if the connection is lost
- Connection reuse for multiple subscriptions

## Usage

### Creating a Client

```go
import (
    "github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/websocketclient"
)

// Create a new client with options
client, err := websocketclient.NewClient(
    websocketclient.WithLogger(logger),
    websocketclient.WithSessionToken(sessionToken),
)
if err != nil {
    log.Fatalf("Failed to create WebSocket client: %v", err)
}
```

### Connecting to the WebSocket Server

```go
// Create a context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Connect to the WebSocket server
if err := client.Connect(ctx); err != nil {
    log.Fatalf("Failed to connect to WebSocket: %v", err)
}
defer client.Close()
```

### Subscribing to Timeline Transactions

```go
// Subscribe to timeline transactions
transactionsCh, err := client.SubscribeToTimelineTransactions(ctx)
if err != nil {
    log.Fatalf("Failed to subscribe to timeline transactions: %v", err)
}

// Wait for data
select {
case data := <-transactionsCh:
    // Process data
    fmt.Println("Received timeline transactions data")
case <-ctx.Done():
    log.Fatalf("Timeout waiting for timeline transactions: %v", ctx.Err())
}
```

### Subscribing to Timeline Detail

```go
// Subscribe to timeline detail
detailsCh, err := client.SubscribeToTimelineDetail(ctx, itemID)
if err != nil {
    log.Fatalf("Failed to subscribe to timeline detail: %v", err)
}

// Wait for data
select {
case data := <-detailsCh:
    // Process data
    fmt.Println("Received timeline detail data")
case <-ctx.Done():
    log.Fatalf("Timeout waiting for timeline detail: %v", ctx.Err())
}
```

## Integration with Existing Code

The WebSocket client can be integrated with the existing code using the `reader.WebSocketReader` implementation:

```go
// Create WebSocket client
wsClient, err := websocketclient.NewClient(
    websocketclient.WithLogger(logger),
    websocketclient.WithSessionToken(sessionToken),
)
if err != nil {
    log.Fatalf("Failed to create WebSocket client: %v", err)
}

// Create WebSocket reader
wsReader := reader.NewWebSocketReader(wsClient, sessionToken, logger)

// Create clients
transactionsClient := transactions.NewClient(wsReader, logger)
detailsClient := details.NewClient(wsReader, logger)

// Use the clients as usual
var allTransactions []TransactionItem
err = transactionsClient.List(&allTransactions)
```

## Connection Management

- The client maintains a single WebSocket connection for all subscriptions
- The connection is established when you call `Connect()` and is reused for subsequent subscriptions
- The client automatically unsubscribes after receiving data
- The client handles reconnection automatically if the connection is lost

## Error Handling

The client provides detailed error messages for various failure scenarios:

- Connection failures
- Subscription failures
- Authentication errors
- Timeout errors

## Examples

See the `examples` directory for complete examples of how to use the WebSocket client:

- `websocket_timeline_example.go`: Simple example of subscribing to timeline transactions and fetching details
- `websocket_timeline_with_api.go`: Example of using the WebSocket client with the existing API clients
