package websocketclient_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/websocketclient"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a client
	client, err := websocketclient.NewClient(
		websocketclient.WithLogger(logger),
	)

	// Verify results
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestMockClient(t *testing.T) {
	t.Parallel()

	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	mockClient := websocketclient.NewMockClientInterface(ctrl)

	// Set up expectations
	ctx := t.Context()
	timelineCh := make(chan []byte, 1)
	timelineReadCh := (<-chan []byte)(timelineCh)
	timelineWithCursorCh := make(chan []byte, 1)
	timelineWithCursorReadCh := (<-chan []byte)(timelineWithCursorCh)
	timelineDetailCh := make(chan []byte, 1)
	timelineDetailReadCh := (<-chan []byte)(timelineDetailCh)

	mockClient.EXPECT().
		Connect(gomock.Any()).
		Return(nil)

	mockClient.EXPECT().
		SubscribeToTimelineTransactions(gomock.Any()).
		Return(timelineReadCh, nil)

	mockClient.EXPECT().
		SubscribeToTimelineTransactionsWithCursor(gomock.Any(), gomock.Eq("cursor123")).
		Return(timelineWithCursorReadCh, nil)

	mockClient.EXPECT().
		SubscribeToTimelineDetail(gomock.Any(), gomock.Eq("US0378331005")).
		Return(timelineDetailReadCh, nil)

	mockClient.EXPECT().
		Close().
		Return(nil)

	// Call methods
	err := mockClient.Connect(ctx)
	require.NoError(t, err)

	ch1, err := mockClient.SubscribeToTimelineTransactions(ctx)
	require.NoError(t, err)
	assert.Equal(t, timelineReadCh, ch1)

	ch2, err := mockClient.SubscribeToTimelineTransactionsWithCursor(ctx, "cursor123")
	require.NoError(t, err)
	assert.Equal(t, timelineWithCursorReadCh, ch2)

	ch3, err := mockClient.SubscribeToTimelineDetail(ctx, "US0378331005")
	require.NoError(t, err)
	assert.Equal(t, timelineDetailReadCh, ch3)

	err = mockClient.Close()
	require.NoError(t, err)
}
