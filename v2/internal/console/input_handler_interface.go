// Package console provides utilities for console-based user interaction.
package console

// InputHandlerInterface defines methods for handling user input for authentication.
type InputHandlerInterface interface {
	// GetPhoneNumber prompts the user to enter their phone number.
	GetPhoneNumber() (string, error)

	// GetPIN prompts the user to enter their PIN.
	GetPIN() (string, error)

	// GetOTP prompts the user to enter their one-time password (OTP).
	GetOTP() (string, error)
}
