package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/websocketclient"
)

const (
	// Default timeout in seconds.
	defaultTimeoutSeconds = 60

	// Directory permissions.
	dirPermissions = 0700

	// File permissions.
	filePermissions = 0600
)

// TransactionItem represents a transaction item from the timeline.
type TransactionItem struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Action struct {
		Payload string `json:"payload"`
	} `json:"action"`
	Timestamp string `json:"timestamp"`
}

// TransactionsResponse represents the response from the timeline transactions API.
type TransactionsResponse struct {
	Items   []TransactionItem `json:"items"`
	Cursors struct {
		After string `json:"after"`
	} `json:"cursors"`
}

//nolint:gocognit,cyclop,funlen
func main() {
	// Parse command line flags
	var (
		debug       bool
		maxItems    int
		timeoutSecs int
	)

	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.IntVar(&maxItems, "max-items", 0, "Maximum number of items to process (0 = all)")
	flag.IntVar(&timeoutSecs, "timeout", defaultTimeoutSeconds, "Timeout in seconds for the entire operation")
	flag.Parse()

	// Set up logging
	logger := log.New()
	logger.SetFormatter(&nested.Formatter{})

	if debug {
		logger.SetLevel(log.DebugLevel)
	}

	// Get credentials from environment variables
	phoneNumber := os.Getenv("TR_PHONE_NUMBER")
	pin := os.Getenv("TR_PIN")

	if phoneNumber == "" || pin == "" {
		logger.Fatal("Please set TR_PHONE_NUMBER and TR_PIN environment variables")
	}

	// Create API client and authenticate
	apiClient := api.NewClient(logger)

	authClient, err := auth.NewClient(apiClient, logger)
	if err != nil {
		logger.Fatalf("Failed to create auth client: %v", err)
	}

	// Login
	resp, err := authClient.Login(phoneNumber, pin)
	if err != nil {
		logger.Fatalf("Failed to login: %v", err)
	}

	// Handle 2FA if needed
	processID := resp.ProcessID
	if processID != "" {
		fmt.Println("2FA required, please check your phone for the code")

		var otp string

		fmt.Print("Enter 2FA code: ")
		_, _ = fmt.Scanln(&otp)

		err = authClient.ProvideOTP(processID, otp)
		if err != nil {
			logger.Fatalf("Failed to validate OTP: %v", err)
		}
	}

	// Get session token
	sessionToken := authClient.SessionToken().Value()

	logger.Info("Successfully authenticated")

	// Create WebSocket client
	wsClient, err := websocketclient.NewClient(
		websocketclient.WithLogger(logger),
		websocketclient.WithSessionToken(sessionToken),
	)
	if err != nil {
		logger.Fatalf("Failed to create WebSocket client: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	// Connect to WebSocket
	if err := wsClient.Connect(ctx); err != nil {
		logger.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer wsClient.Close()

	logger.Info("Connected to WebSocket, subscribing to timeline transactions...")

	// Create directories for saving responses
	if err := createDirectories(); err != nil {
		logger.Fatalf("Failed to create directories: %v", err)
	}

	logger.Info("Will save both raw and formatted responses to the filesystem")

	// Get all transactions with pagination
	var allTransactions []TransactionItem

	var cursor string

	var page int

	for {
		page++
		logger.Infof("Fetching transactions page %d...", page)

		// Subscribe to timeline transactions with cursor if available
		var transactionsCh <-chan []byte

		var err error

		if cursor == "" {
			transactionsCh, err = wsClient.SubscribeToTimelineTransactions(ctx)
		} else {
			transactionsCh, err = wsClient.SubscribeToTimelineTransactionsWithCursor(ctx, cursor)
		}

		if err != nil {
			logger.Fatalf("Failed to subscribe to timeline transactions: %v", err)
		}

		// Wait for transactions data
		logger.Infof("Waiting for timeline transactions data (page %d)...", page)

		var transactionsData []byte

		select {
		case data := <-transactionsCh:
			transactionsData = data

			logger.Infof("Received timeline transactions data for page %d", page)
		case <-ctx.Done():
			logger.Fatalf("Timeout waiting for timeline transactions: %v", ctx.Err())
		}

		// Save transactions data for this page
		pageFilename := fmt.Sprintf("page_%d", page)
		if err := saveTransactionsData(pageFilename, transactionsData); err != nil {
			logger.Fatalf("Failed to save transactions data: %v", err)
		}

		// Parse transactions data
		var transactions TransactionsResponse
		if err := json.Unmarshal(transactionsData, &transactions); err != nil {
			logger.Fatalf("Failed to parse transactions data: %v", err)
		}

		// Add items to our collection
		allTransactions = append(allTransactions, transactions.Items...)
		logger.Infof("Added %d transactions from page %d (total: %d)",
			len(transactions.Items), page, len(allTransactions))

		// Check if we have more pages
		if transactions.Cursors.After == "" {
			logger.Info("No more pages available")

			break
		}

		// Update cursor for next page
		cursor = transactions.Cursors.After
		logger.Infof("Next page cursor: %s", cursor)
	}

	if len(allTransactions) == 0 {
		logger.Fatal("No transactions found")
	}

	// Save all transactions to a single file
	allTransactionsData := map[string]interface{}{
		"items": allTransactions,
	}

	allTransactionsBytes, err := json.Marshal(allTransactionsData)
	if err != nil {
		logger.Fatalf("Failed to marshal all transactions: %v", err)
	}

	if err := saveTransactionsData("all_transactions", allTransactionsBytes); err != nil {
		logger.Fatalf("Failed to save all transactions data: %v", err)
	}

	logger.Infof("Found %d transactions", len(allTransactions))

	// Process transactions
	itemsToProcess := len(allTransactions)
	if maxItems > 0 && maxItems < itemsToProcess {
		itemsToProcess = maxItems
	}

	logger.Infof("Processing %d transactions", itemsToProcess)

	for i := 0; i < itemsToProcess; i++ {
		transaction := allTransactions[i]
		logger.Infof("Processing transaction %d/%d (ID: %s, Type: %s)",
			i+1, itemsToProcess, transaction.ID, transaction.Type)

		// Get details for the transaction
		logger.Infof("Fetching details for transaction with payload: %s", transaction.Action.Payload)

		detailsCh, err := wsClient.SubscribeToTimelineDetail(ctx, transaction.Action.Payload)
		if err != nil {
			logger.Errorf("Failed to subscribe to timeline detail: %v", err)

			continue
		}

		// Wait for details data
		logger.Debugf("Waiting for timeline detail data...")

		var detailsData []byte
		select {
		case data := <-detailsCh:
			detailsData = data

			logger.Debugf("Received timeline detail data")
		case <-ctx.Done():
			logger.Errorf("Timeout waiting for timeline detail: %v", ctx.Err())

			continue
		}

		// Save details data
		if err := saveDetailsData(transaction.ID, detailsData); err != nil {
			logger.Errorf("Failed to save details data: %v", err)

			continue
		}

		logger.Infof("Successfully processed transaction %d/%d", i+1, itemsToProcess)
	}

	logger.Info("All transactions processed successfully")
}

// createDirectories creates the necessary directories for saving responses.
func createDirectories() error {
	dirs := []string{"transactions", "details"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, dirPermissions); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// saveTransactionsData saves the transactions data to a file.
func saveTransactionsData(filename string, data []byte) error {
	// Save raw response data
	rawFilename := filepath.Join("transactions", filename+".raw.json")
	if err := os.WriteFile(rawFilename, data, filePermissions); err != nil {
		return fmt.Errorf("failed to write raw file: %w", err)
	}

	// Format the data for better readability.
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	formattedData, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}

	// Save formatted data to file.
	formattedFilename := filepath.Join("transactions", filename+".json")
	if err := os.WriteFile(formattedFilename, formattedData, filePermissions); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Also save individual transaction files.
	var transactions TransactionsResponse
	if err := json.Unmarshal(data, &transactions); err != nil {
		return fmt.Errorf("failed to parse transactions: %w", err)
	}

	for _, transaction := range transactions.Items {
		// Create a map with just this transaction.
		singleTransaction := map[string]interface{}{
			"id":        transaction.ID,
			"type":      transaction.Type,
			"action":    transaction.Action,
			"timestamp": transaction.Timestamp,
		}

		// Format the data
		transactionData, err := json.MarshalIndent(singleTransaction, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format transaction JSON: %w", err)
		}

		// Save to file
		transFilename := filepath.Join("transactions", transaction.ID+".json")
		if err := os.WriteFile(transFilename, transactionData, filePermissions); err != nil {
			return fmt.Errorf("failed to write transaction file: %w", err)
		}

		// We don't save individual raw transaction data since they're part of the full response
		// and we already saved the raw full response
	}

	return nil
}

// saveDetailsData saves the details data to a file.
func saveDetailsData(id string, data []byte) error {
	// Save raw response data
	rawFilename := filepath.Join("details", id+".raw.json")
	if err := os.WriteFile(rawFilename, data, filePermissions); err != nil {
		return fmt.Errorf("failed to write raw file: %w", err)
	}

	// Format the data for better readability
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	formattedData, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}

	// Save formatted data to file
	formattedFilename := filepath.Join("details", id+".json")
	if err := os.WriteFile(formattedFilename, formattedData, filePermissions); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
