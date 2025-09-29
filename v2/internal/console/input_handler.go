// Package console provides utilities for console-based user interaction.
package console

import (
	"fmt"
)

// InputHandler provides methods to handle user input for authentication.
type InputHandler struct{}

// NewInputHandler creates a new InputHandler.
func NewInputHandler() InputHandler {
	return InputHandler{}
}

// GetPhoneNumber prompts the user to enter their phone number.
func (h InputHandler) GetPhoneNumber() (string, error) {
	var phoneNumber string

	fmt.Println("Enter phone number in international format (+49xxxxxxxxxxxxx):")

	if _, err := fmt.Scanln(&phoneNumber); err != nil {
		return "", fmt.Errorf("could not read phone number: %w", err)
	}

	return phoneNumber, nil
}

// GetPIN prompts the user to enter their PIN.
func (h InputHandler) GetPIN() (string, error) {
	pin, err := ReadPassword("PIN")
	if err != nil {
		return "", fmt.Errorf("could not read PIN: %w", err)
	}

	return string(pin), nil
}

// GetOTP prompts the user to enter their one-time password (OTP).
func (h InputHandler) GetOTP() (string, error) {
	fmt.Println("2FA required, please check your phone for the code")

	otp, err := ReadPassword("2FA token")
	if err != nil {
		return "", fmt.Errorf("could not read OTP: %w", err)
	}

	return string(otp), nil
}
