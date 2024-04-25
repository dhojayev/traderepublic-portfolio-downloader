package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/header"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/portfolio"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"
)

const (
	baseHost = "api.traderepublic.com"
)

type Reader struct {
	authService console.AuthServiceInterface
	writer      writer.Interface
	logger      *log.Logger
	conn        *websocket.Conn
	subID       uint
}

func NewReader(authService console.AuthServiceInterface, writer writer.Interface, logger *log.Logger) (*Reader, error) {
	client := &Reader{
		authService: authService,
		writer:      writer,
		logger:      logger,
	}

	return client, client.connect()
}

func (r *Reader) connect() error {
	u := url.URL{Scheme: "wss", Host: baseHost, Path: "/"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header.NewHeaders().AsHTTPHeader())
	if err != nil {
		return fmt.Errorf("could not connect to websocket: %w", err)
	}

	r.conn = conn

	if err = r.conn.WriteMessage(websocket.TextMessage, []byte(`connect 31 {"locale": "de"}`)); err != nil {
		return fmt.Errorf("could not send connect msg: %w", err)
	}

	_, msg, err := r.conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("could not read connect msg: %w", err)
	}

	r.logger.WithField("message", string(msg)).Trace("received msg")

	return nil
}

func (r *Reader) reconnect() error {
	_ = r.Close()

	return r.connect()
}

func (r *Reader) Close() error {
	if r.conn == nil {
		return errors.New("cannot close websocket: connection not established")
	}

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("could not close websocket connection: %w", err)
	}

	return nil
}

//nolint:cyclop,ireturn
func (r *Reader) Read(dataType string, dataMap map[string]any) (portfolio.OutputDataInterface, error) {
	r.subID++

	dataBytes, err := r.createWritableDataBytes(dataType, dataMap)
	if err != nil {
		return Message{}, err
	}

	err = r.conn.WriteMessage(websocket.TextMessage, dataBytes)
	if err != nil {
		return Message{}, fmt.Errorf("could not send message: %w", err)
	}

	r.logger.WithField("message", string(dataBytes)).Trace("sent message")

	for {
		_, msg, err := r.conn.ReadMessage()
		if err != nil {
			return Message{}, fmt.Errorf("could not read message: %w", err)
		}

		r.logger.WithField("message", string(msg)).Trace("received msg")

		message, err := NewMessage(msg)
		if err != nil {
			return message, fmt.Errorf("could not create message struct: %w", err)
		}

		switch {
		case message.HasContinueState():
			continue
		case message.HasErrorState():
			if message.HasAuthErrMsg() {
				if loginErr := r.authService.Login(); loginErr != nil {
					return message, fmt.Errorf("could not re-login: %w", loginErr)
				}

				if err = r.reconnect(); err != nil {
					return message, err
				}

				return r.Read(dataType, dataMap)
			}

			return message, fmt.Errorf("error state received: %s", msg)
		}

		if err := r.writer.Bytes(dataType, message.Data()); err != nil {
			return message, fmt.Errorf("could not write message: %w", err)
		}

		return message, nil
	}
}

func (r *Reader) createWritableDataBytes(dataType string, dataMap map[string]any) ([]byte, error) {
	data := dataMap
	data["type"] = dataType
	data["token"] = r.authService.SessionToken().Value()

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("could not marshal data into json: %w", err)
	}

	return []byte(fmt.Sprintf("sub %d %s", r.subID, dataBytes)), nil
}
