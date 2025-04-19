package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/websocketclient"
	"github.com/joho/godotenv"
)

const (
	// Filepaths for saving responses.
	filepathTransactions = "./debug/transactions"

	// Filepaths for saving details.
	filepathDetails = "./debug/details"

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

func main() {
	// Use this variable to track exit code
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	_ = godotenv.Load(".env")

	// Parse command line flags and set up logging
	config := parseFlags()
	logger := setupLogger(config.debug)

	// Create directories for saving responses
	if err := createDirectories(); err != nil {
		logger.Error("Failed to create directories", "error", err)
		
		exitCode = 1

		return
	}

	// Authenticate and get session token
	sessionToken, err := authenticate(logger)
	if err != nil {
		exitCode = 1

		return
	}

	// Create and connect WebSocket client
	wsClient, ctx, cancel, err := setupWebSocketClient(logger, sessionToken, config.timeoutSecs)
	if err != nil {
		exitCode = 1

		return
	}

	defer cancel()

	defer wsClient.Close()

	logger.Info("Will save both raw and formatted responses to the filesystem")

	// Fetch all transactions
	allTransactions, err := fetchAllTransactions(ctx, logger, wsClient)
	if err != nil {
		exitCode = 1

		return
	}

	// Process transactions
	processTransactions(ctx, logger, wsClient, allTransactions, config.maxItems)

	logger.Info("All transactions processed successfully")
}

// Config holds the command line configuration.
type Config struct {
	debug       bool
	maxItems    int
	timeoutSecs int
}

// parseFlags parses command line flags and returns a Config.
func parseFlags() Config {
	var config Config

	flag.BoolVar(&config.debug, "debug", false, "Enable debug logging")
	flag.IntVar(&config.maxItems, "max-items", 0, "Maximum number of items to process (0 = all)")
	flag.IntVar(&config.timeoutSecs, "timeout", defaultTimeoutSeconds, "Timeout in seconds for the entire operation")
	flag.Parse()

	return config
}

// setupLogger creates and configures a logger.
func setupLogger(debug bool) *slog.Logger {
	logOpts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if debug {
		logOpts.Level = slog.LevelDebug
	}

	return slog.New(slog.NewTextHandler(os.Stdout, logOpts))
}

// authenticate performs authentication and returns a session token.
//
//nolint:cyclop,funlen
func authenticate(logger *slog.Logger) (string, error) {
	// Create credentials service for storing tokens
	credentials := auth.NewFileCredentialsService(internal.AuthTokenFilename)

	// Try to load existing tokens first
	if err := credentials.Load(); err == nil {
		sessionToken := credentials.GetSessionToken()
		if sessionToken != "" {
			logger.Info("Using existing session token")

			// We'll verify the token by making a test API call later
			// If it fails, we'll need to re-authenticate
			return sessionToken, nil
		}
	}

	// Create input handler for user interaction
	inputHandler := console.NewInputHandler()

	// If no valid tokens found, proceed with authentication
	phoneNumber := os.Getenv("TR_PHONE_NUMBER")
	pin := os.Getenv("TR_PIN")

	// If environment variables are not set, prompt the user
	if phoneNumber == "" {
		var err error
		
		phoneNumber, err = inputHandler.GetPhoneNumber()
		if err != nil {
			logger.Error("Failed to get phone number", "error", err)

			return "", fmt.Errorf("failed to get phone number: %w", err)
		}
	}

	if pin == "" {
		var err error
		
		pin, err = inputHandler.GetPIN()
		if err != nil {
			logger.Error("Failed to get PIN", "error", err)

			return "", fmt.Errorf("failed to get PIN: %w", err)
		}
	}

	// Create API client and authenticate
	apiClient, err := api.NewClient()
	if err != nil {
		logger.Error("Failed to create API client", "error", err)

		return "", fmt.Errorf("failed to create API client: %w", err)
	}

	// Create auth client
	authClient, err := auth.NewClient(apiClient)
	if err != nil {
		logger.Error("Failed to create auth client", "error", err)

		return "", fmt.Errorf("failed to create auth client: %w", err)
	}

	// Login
	processID, err := authClient.Login(auth.PhoneNumber(phoneNumber), auth.Pin(pin))
	if err != nil {
		logger.Error("Failed to login", "error", err)

		return "", fmt.Errorf("failed to login: %w", err)
	}

	if processID != "" {
		// Get OTP from user
		otp, err := inputHandler.GetOTP()
		if err != nil {
			logger.Error("Failed to get OTP", "error", err)
			
			return "", fmt.Errorf("failed to get OTP: %w", err)
		}

		// Get tokens from OTP verification
		token, err := authClient.ProvideOTP(processID, auth.OTP(otp))
		if err != nil {
			logger.Error("Failed to validate OTP", "error", err)

			return "", fmt.Errorf("failed to validate OTP: %w", err)
		}

		// Store tokens using credentials service
		if err := credentials.Store(token.SessionToken(), token.RefreshToken()); err != nil {
			logger.Error("Failed to store tokens", "error", err)
		}

		return token.SessionToken(), nil
	}

	logger.Info("Successfully authenticated")

	return credentials.GetSessionToken(), nil
}

// setupWebSocketClient creates and connects a WebSocket client.
func setupWebSocketClient(
	logger *slog.Logger,
	sessionToken string,
	timeoutSecs int,
) (*websocketclient.Client, context.Context, context.CancelFunc, error) {
	// Create WebSocket client
	wsClient, err := websocketclient.NewClient(
		websocketclient.WithLogger(logger),
		websocketclient.WithSessionToken(sessionToken),
	)
	if err != nil {
		logger.Error("Failed to create WebSocket client", "error", err)

		return nil, nil, nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)

	// Connect to WebSocket
	if err := wsClient.Connect(ctx); err != nil {
		cancel()
		logger.Error("Failed to connect to WebSocket", "error", err)

		return nil, nil, nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	logger.Info("Connected to WebSocket, subscribing to timeline transactions...")

	return wsClient, ctx, cancel, nil
}

// fetchAllTransactions fetches all transactions with pagination.
func fetchAllTransactions(
	ctx context.Context,
	logger *slog.Logger,
	wsClient *websocketclient.Client,
) ([]TransactionItem, error) {
	var allTransactions []TransactionItem

	var cursor string

	var page int

	for {
		page++
		logger.Info("Fetching transactions page", "page", page)

		// Fetch page data
		transactionsData, err := fetchTransactionPage(ctx, logger, wsClient, cursor, page)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch transaction page: %w", err)
		}

		// Parse and process page data
		transactions, nextCursor, err := processTransactionPage(logger, transactionsData, page, allTransactions)
		if err != nil {
			return nil, fmt.Errorf("failed to process transaction page: %w", err)
		}

		// Add items to our collection
		allTransactions = append(allTransactions, transactions.Items...)
		logger.Info("Added transactions",
			"count", len(transactions.Items),
			"page", page,
			"total", len(allTransactions))

		// Check if we have more pages
		if nextCursor == "" {
			logger.Info("No more pages available")

			break
		}

		// Update cursor for next page
		cursor = nextCursor
		logger.Info("Next page cursor", "cursor", cursor)
	}

	if len(allTransactions) == 0 {
		logger.Error("No transactions found")

		return nil, errors.New("no transactions found")
	}

	// Save all transactions to a single file
	return saveAllTransactions(logger, allTransactions)
}

// fetchTransactionPage fetches a single page of transactions.
func fetchTransactionPage(
	ctx context.Context,
	logger *slog.Logger,
	wsClient *websocketclient.Client,
	cursor string,
	page int,
) ([]byte, error) {
	// Subscribe to timeline transactions with cursor if available
	var transactionsCh <-chan []byte

	var err error

	if cursor == "" {
		transactionsCh, err = wsClient.SubscribeToTimelineTransactions(ctx)
	} else {
		transactionsCh, err = wsClient.SubscribeToTimelineTransactionsWithCursor(ctx, cursor)
	}

	if err != nil {
		logger.Error("Failed to subscribe to timeline transactions", "error", err)

		return nil, fmt.Errorf("failed to subscribe to timeline transactions: %w", err)
	}

	// Wait for transactions data
	logger.Info("Waiting for timeline transactions data", "page", page)

	var transactionsData []byte

	select {
	case data := <-transactionsCh:
		transactionsData = data

		logger.Info("Received timeline transactions data", "page", page)
	case <-ctx.Done():
		logger.Error("Timeout waiting for timeline transactions", "error", ctx.Err())

		return nil, fmt.Errorf("timeout waiting for timeline transactions: %w", ctx.Err())
	}

	// Save transactions data for this page
	pageFilename := fmt.Sprintf("page_%d", page)
	if err := saveTransactionsData(pageFilename, transactionsData); err != nil {
		logger.Error("Failed to save transactions data", "error", err)

		return nil, fmt.Errorf("failed to save transactions data: %w", err)
	}

	return transactionsData, nil
}

// processTransactionPage processes a page of transaction data.
func processTransactionPage(
	logger *slog.Logger,
	transactionsData []byte,
	_ int, // page number, not used but kept for clarity
	_ []TransactionItem, // allTransactions, not used but kept for clarity
) (TransactionsResponse, string, error) {
	// Parse transactions data
	var transactions TransactionsResponse

	if err := json.Unmarshal(transactionsData, &transactions); err != nil {
		logger.Error("Failed to parse transactions data", "error", err)

		return TransactionsResponse{}, "", fmt.Errorf("failed to parse transactions data: %w", err)
	}

	return transactions, transactions.Cursors.After, nil
}

// saveAllTransactions saves all transactions to a single file.
func saveAllTransactions(logger *slog.Logger, allTransactions []TransactionItem) ([]TransactionItem, error) {
	allTransactionsData := map[string]interface{}{
		"items": allTransactions,
	}

	allTransactionsBytes, err := json.Marshal(allTransactionsData)
	if err != nil {
		logger.Error("Failed to marshal all transactions", "error", err)

		return nil, fmt.Errorf("failed to marshal all transactions: %w", err)
	}

	if err := saveTransactionsData("all_transactions", allTransactionsBytes); err != nil {
		logger.Error("Failed to save all transactions data", "error", err)

		return nil, fmt.Errorf("failed to save all transactions data: %w", err)
	}

	logger.Info("Found transactions", "count", len(allTransactions))

	return allTransactions, nil
}

// processTransactions processes the transactions and fetches details for each.
func processTransactions(
	ctx context.Context,
	logger *slog.Logger,
	wsClient *websocketclient.Client,
	allTransactions []TransactionItem,
	maxItems int,
) {
	// Determine how many items to process
	itemsToProcess := len(allTransactions)
	if maxItems > 0 && maxItems < itemsToProcess {
		itemsToProcess = maxItems
	}

	logger.Info("Processing transactions", "count", itemsToProcess)

	for i := 0; i < itemsToProcess; i++ {
		transaction := allTransactions[i]
		logger.Info("Processing transaction",
			"current", i+1,
			"total", itemsToProcess,
			"id", transaction.ID,
			"type", transaction.Type)

		// Get details for the transaction
		logger.Info("Fetching details for transaction", "payload", transaction.Action.Payload)

		detailsCh, err := wsClient.SubscribeToTimelineDetail(ctx, transaction.Action.Payload)
		if err != nil {
			logger.Error("Failed to subscribe to timeline detail", "error", err)

			continue
		}

		// Wait for details data
		logger.Debug("Waiting for timeline detail data")

		var detailsData []byte
		select {
		case data := <-detailsCh:
			detailsData = data

			logger.Debug("Received timeline detail data")
		case <-ctx.Done():
			logger.Error("Timeout waiting for timeline detail", "error", ctx.Err())

			continue
		}

		// Save details data
		if err := saveDetailsData(transaction.ID, detailsData); err != nil {
			logger.Error("Failed to save details data", "error", err)

			continue
		}

		logger.Info("Successfully processed transaction", "current", i+1, "total", itemsToProcess)
	}
}

// createDirectories creates the necessary directories for saving responses.
func createDirectories() error {
	dirs := []string{filepathTransactions, filepathDetails}
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
	rawFilename := filepath.Join(filepathTransactions, filename+".raw.json")
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
	formattedFilename := filepath.Join(filepathTransactions, filename+".json")
	if err := os.WriteFile(formattedFilename, formattedData, filePermissions); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Also save individual transaction files.
	var transactions TransactionsResponse
	if err := json.Unmarshal(data, &transactions); err != nil {
		return fmt.Errorf("failed to parse transactions: %w", err)
	}

	// Process each transaction
	for _, transaction := range transactions.Items {
		if err := saveIndividualTransaction(transaction); err != nil {
			return fmt.Errorf("failed to save individual transaction: %w", err)
		}
	}

	return nil
}

// saveIndividualTransaction saves a single transaction to a file.
func saveIndividualTransaction(transaction TransactionItem) error {
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
	transFilename := filepath.Join(filepathTransactions, transaction.ID+".json")
	if err := os.WriteFile(transFilename, transactionData, filePermissions); err != nil {
		return fmt.Errorf("failed to write transaction file: %w", err)
	}

	// We don't save individual raw transaction data since they're part of the full response
	return nil
}

// saveDetailsData saves the details data to a file.
func saveDetailsData(id string, data []byte) error {
	// Save raw response data
	rawFilename := filepath.Join(filepathDetails, id+".raw.json")
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
	formattedFilename := filepath.Join(filepathDetails, id+".json")
	if err := os.WriteFile(formattedFilename, formattedData, filePermissions); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
