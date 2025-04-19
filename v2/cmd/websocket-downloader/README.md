# WebSocket Downloader

This tool downloads transaction data and details from the Trade Republic API using WebSockets and saves them to the filesystem.

## Features

- Authenticates with the Trade Republic API
- Downloads all timeline transactions
- Downloads details for each transaction
- Saves responses to the filesystem in JSON format
- Supports pagination for large transaction lists
- Configurable timeout and maximum number of items to process

## Usage

### Environment Variables

Set the following environment variables before running the tool:

```bash
export TR_PHONE_NUMBER="+49123456789"
export TR_PIN="123456"
```

### Running the Tool

You can run the tool using the Makefile targets:

```bash
# Run with default settings
make run-websocket-downloader

# Run with debug logging
make run-websocket-downloader-debug

# Build the tool
make build-websocket-downloader
```

Or you can run it directly:

```bash
go run ./cmd/websocket-downloader/main.go [options]
```

### Command Line Options

- `--debug`: Enable debug logging
- `--max-items=N`: Maximum number of items to process (0 = all)
- `--timeout=N`: Timeout in seconds for the entire operation (default: 60)

### Examples

```bash
# Download all transactions with debug logging
go run ./cmd/websocket-downloader/main.go --debug

# Download only the first 10 transactions
go run ./cmd/websocket-downloader/main.go --max-items=10

# Set a longer timeout (2 minutes)
go run ./cmd/websocket-downloader/main.go --timeout=120
```

## Output

The tool creates the following directories and files:

### Raw Data (Exactly as Received from API)
- `transactions/page_N.raw.json`: Raw response for each page of transactions
- `transactions/all_transactions.raw.json`: Raw response for combined transactions
- `details/{id}.raw.json`: Raw response for transaction details

### Formatted Data (Pretty-Printed JSON)
- `transactions/page_N.json`: Formatted response for each page of transactions
- `transactions/all_transactions.json`: Formatted response for combined transactions
- `transactions/{id}.json`: Formatted data for individual transactions
- `details/{id}.json`: Formatted data for transaction details

## Implementation Details

The tool uses the WebSocket client to connect to the Trade Republic API and subscribe to timeline transactions and details. It then saves the responses to the filesystem in JSON format.

The WebSocket client automatically handles:
- Authentication
- Connection management
- Subscription and unsubscription
- Reconnection if the connection is lost

## Dependencies

- github.com/sirupsen/logrus: For logging
- github.com/antonfisher/nested-logrus-formatter: For formatted logging
- Internal WebSocket client implementation
