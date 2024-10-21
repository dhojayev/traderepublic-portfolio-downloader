package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	minDataParts  = 2
	stateContinue = "C"
	stateError    = "E"

	errorCodeAuth = "AUTHENTICATION_ERROR"
)

type Message struct {
	id    uint
	state string
	data  []byte
}

func NewMessage(data []byte) (Message, error) {
	msg := Message{}
	parts := strings.Split(string(data), " ")

	if len(parts) < minDataParts {
		return msg, errors.New("could not parse the contents")
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return msg, fmt.Errorf("could not convert id string to int: %w", err)
	}

	//nolint:gosec // disable G115 until its fixed
	msg.id = uint(id)
	msg.state = parts[1]
	msg.data = []byte(strings.Join(parts[2:], " "))

	return msg, nil
}

func (m Message) HasErrorState() bool {
	return m.state == stateError
}

func (m Message) HasContinueState() bool {
	return m.state == stateContinue
}

func (m Message) HasAuthErrMsg() bool {
	var respErr ResponseErrors

	if err := json.Unmarshal(m.data, &respErr); err != nil {
		return false
	}

	for _, v := range respErr.Errors {
		if v.ErrorCode == errorCodeAuth {
			return true
		}
	}

	return false
}

func (m Message) ID() uint {
	return m.id
}

func (m Message) State() string {
	return m.state
}

func (m Message) Data() []byte {
	return m.data
}
