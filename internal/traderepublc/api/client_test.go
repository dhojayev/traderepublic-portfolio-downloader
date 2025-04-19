package api_test

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/restclient"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api"
)

const (
	testProcessID = "test-process-id"
)

func TestClient_Login_Success(t *testing.T) {
	t.Parallel()
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	mockClient := api.NewMockClientInterface(ctrl)

	// Set up expectations
	loginRequest := restclient.APILoginRequest{
		PhoneNumber: "+1234567890",
		Pin:         "1234",
	}
	refreshToken := api.NewToken(api.TokenNameRefresh, "")
	processID := testProcessID
		expectedResponse := restclient.APILoginResponse{
		ProcessId: &processID,
	}
	expectedSessionToken := api.NewToken(api.TokenNameSession, "test-session-token")

	mockClient.EXPECT().
		Login(gomock.Eq(loginRequest), gomock.Eq(refreshToken)).
		Return(expectedResponse, expectedSessionToken, nil)

	// Call Login
	response, sessionToken, err := mockClient.Login(loginRequest, refreshToken)

	// Verify results
	assert.NoError(t, err)
	assert.Equal(t, testProcessID, *response.ProcessId)
	assert.Equal(t, "test-session-token", sessionToken.Value())
}

func TestClient_Login_Error(t *testing.T) {
	t.Parallel()
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	mockClient := api.NewMockClientInterface(ctrl)

	// Set up expectations
	loginRequest := restclient.APILoginRequest{
		PhoneNumber: "+1234567890",
		Pin:         "wrong-pin",
	}
	refreshToken := api.NewToken(api.TokenNameRefresh, "")
	expectedError := errors.New("invalid credentials")

	mockClient.EXPECT().
		Login(gomock.Eq(loginRequest), gomock.Eq(refreshToken)).
		Return(restclient.APILoginResponse{}, api.NewToken(api.TokenNameSession, ""), expectedError)

	// Call Login
	_, _, err := mockClient.Login(loginRequest, refreshToken)

	// Verify error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestClient_PostOTP_Success(t *testing.T) {
	t.Parallel()
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	mockClient := api.NewMockClientInterface(ctrl)

	// Set up expectations
	processID := "test-process-id"
	otp := "123456"
	expectedSessionToken := api.NewToken(api.TokenNameSession, "test-session-token")
	expectedRefreshToken := api.NewToken(api.TokenNameRefresh, "test-refresh-token")

	mockClient.EXPECT().
		PostOTP(gomock.Eq(processID), gomock.Eq(otp)).
		Return(expectedSessionToken, expectedRefreshToken, nil)

	// Call PostOTP
	sessionToken, refreshToken, err := mockClient.PostOTP(processID, otp)

	// Verify results
	assert.NoError(t, err)
	assert.Equal(t, "test-session-token", sessionToken.Value())
	assert.Equal(t, "test-refresh-token", refreshToken.Value())
}

func TestClient_PostOTP_EmptyProcessID(t *testing.T) {
	t.Parallel()
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	mockClient := api.NewMockClientInterface(ctrl)

	// Set up expectations
	processID := ""
	otp := "123456"
	expectedError := errors.New("processID cannot be empty")

	mockClient.EXPECT().
		PostOTP(gomock.Eq(processID), gomock.Eq(otp)).
		Return(api.NewToken(api.TokenNameSession, ""), api.NewToken(api.TokenNameRefresh, ""), expectedError)

	// Call PostOTP with empty processID
	_, _, err := mockClient.PostOTP(processID, otp)

	// Verify error
	assert.Error(t, err)
	assert.Equal(t, "processID cannot be empty", err.Error())
}

func TestClient_PostOTP_Error(t *testing.T) {
	t.Parallel()
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	mockClient := api.NewMockClientInterface(ctrl)

	// Set up expectations
	processID := testProcessID
	otp := "wrong-otp"
	expectedError := errors.New("invalid OTP")

	mockClient.EXPECT().
		PostOTP(gomock.Eq(processID), gomock.Eq(otp)).
		Return(api.NewToken(api.TokenNameSession, ""), api.NewToken(api.TokenNameRefresh, ""), expectedError)

	// Call PostOTP
	_, _, err := mockClient.PostOTP(processID, otp)

	// Verify error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid OTP")
}

func TestClient_Session_Success(t *testing.T) {
	t.Parallel()
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	mockClient := api.NewMockClientInterface(ctrl)

	// Set up expectations
	refreshToken := api.NewToken(api.TokenNameRefresh, "test-refresh-token")
	expectedSessionToken := api.NewToken(api.TokenNameSession, "new-session-token")

	mockClient.EXPECT().
		Session(gomock.Eq(refreshToken)).
		Return(expectedSessionToken, nil)

	// Call Session
	sessionToken, err := mockClient.Session(refreshToken)

	// Verify results
	assert.NoError(t, err)
	assert.Equal(t, "new-session-token", sessionToken.Value())
}

func TestClient_Session_Error(t *testing.T) {
	t.Parallel()
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	mockClient := api.NewMockClientInterface(ctrl)

	// Set up expectations
	refreshToken := api.NewToken(api.TokenNameRefresh, "invalid-refresh-token")
	expectedError := errors.New("invalid refresh token")

	mockClient.EXPECT().
		Session(gomock.Eq(refreshToken)).
		Return(api.NewToken(api.TokenNameSession, ""), expectedError)

	// Call Session with invalid refresh token
	_, err := mockClient.Session(refreshToken)

	// Verify error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh token")
}

func TestClient_RequestError(t *testing.T) {
	t.Parallel()
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock client
	mockClient := api.NewMockClientInterface(ctrl)

	// Set up expectations for Login
	loginRequest := restclient.APILoginRequest{
		PhoneNumber: "+1234567890",
		Pin:         "1234",
	}
	refreshToken := api.NewToken(api.TokenNameRefresh, "")
	loginError := errors.New("could not login")

	mockClient.EXPECT().
		Login(gomock.Eq(loginRequest), gomock.Eq(refreshToken)).
		Return(restclient.APILoginResponse{}, api.NewToken(api.TokenNameSession, ""), loginError)

	// Set up expectations for PostOTP
	processID := testProcessID
	otp := "123456"
	otpError := errors.New("could not validate otp")

	mockClient.EXPECT().
		PostOTP(gomock.Eq(processID), gomock.Eq(otp)).
		Return(api.NewToken(api.TokenNameSession, ""), api.NewToken(api.TokenNameRefresh, ""), otpError)

	// Set up expectations for Session
	sessionError := errors.New("could not refresh session")

	mockClient.EXPECT().
		Session(gomock.Eq(refreshToken)).
		Return(api.NewToken(api.TokenNameSession, ""), sessionError)

	// Test Login
	_, _, err := mockClient.Login(loginRequest, refreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not login")

	// Test PostOTP
	_, _, err = mockClient.PostOTP(processID, otp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not validate otp")

	// Test Session
	_, err = mockClient.Session(refreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not refresh session")
}
