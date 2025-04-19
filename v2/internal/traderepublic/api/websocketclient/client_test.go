package websocketclient_test

import (
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/traderepublic/api/websocketclient"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	// Create a logger
	logger := log.New()

	// Create a client
	client, err := websocketclient.NewClient(
		websocketclient.WithLogger(logger),
		websocketclient.WithSessionToken("test-token"),
	)

	// Verify results
	assert.NoError(t, err)
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
	portfolioCh := make(chan []byte, 1)
	portfolioReadCh := (<-chan []byte)(portfolioCh)
	instrumentCh := make(chan []byte, 1)
	instrumentReadCh := (<-chan []byte)(instrumentCh)

	mockClient.EXPECT().
		Connect(gomock.Any()).
		Return(nil)

	mockClient.EXPECT().
		SubscribeToTimeline(gomock.Any()).
		Return(timelineReadCh, nil)

	mockClient.EXPECT().
		SubscribeToPortfolio(gomock.Any()).
		Return(portfolioReadCh, nil)

	mockClient.EXPECT().
		SubscribeToInstrument(gomock.Any(), gomock.Eq("US0378331005")).
		Return(instrumentReadCh, nil)

	mockClient.EXPECT().
		Close().
		Return(nil)

	// Call methods
	err := mockClient.Connect(ctx)
	assert.NoError(t, err)

	ch1, err := mockClient.SubscribeToTimeline(ctx)
	assert.NoError(t, err)
	assert.Equal(t, timelineReadCh, ch1)

	ch2, err := mockClient.SubscribeToPortfolio(ctx)
	assert.NoError(t, err)
	assert.Equal(t, portfolioReadCh, ch2)

	ch3, err := mockClient.SubscribeToInstrument(ctx, "US0378331005")
	assert.NoError(t, err)
	assert.Equal(t, instrumentReadCh, ch3)

	err = mockClient.Close()
	assert.NoError(t, err)
}
