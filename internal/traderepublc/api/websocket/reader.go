package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/traderepublc/api/header"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"
)

var ErrMsgErrorStateReceived = errors.New("error state received")

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
	u := url.URL{Scheme: "wss", Host: internal.WebsocketBaseHost, Path: "/"}

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

//nolint:cyclop
func (r *Reader) Read(dataType string, req reader.Request) (reader.JSONResponse, error) {
	r.subID++

	resp := reader.NewJSONResponse(nil)

	dataBytes, err := r.createWritableDataBytes(dataType, req)
	if err != nil {
		return resp, err
	}

	err = r.conn.WriteMessage(websocket.TextMessage, dataBytes)
	if err != nil {
		return resp, fmt.Errorf("could not send message: %w", err)
	}

	r.logger.WithField("message", string(dataBytes)).Trace("sent message")

	for {
		_, msg, err := r.conn.ReadMessage()
		if err != nil {
			return resp, fmt.Errorf("could not read message: %w", err)
		}

		r.logger.WithField("message", string(msg)).Trace("received msg")

		message, err := NewMessage(msg)
		if err != nil {
			return resp, fmt.Errorf("could not create message struct: %w", err)
		}

		switch {
		case message.HasContinueState():
			continue
		case message.HasErrorState():
			if message.HasAuthErrMsg() {
				if loginErr := r.authService.Login(); loginErr != nil {
					return resp, fmt.Errorf("could not re-login: %w", loginErr)
				}

				if err = r.reconnect(); err != nil {
					return resp, err
				}

				return r.Read(dataType, req)
			}

			return resp, fmt.Errorf("%w: %s", ErrMsgErrorStateReceived, msg)
		}

		if err := r.writer.Bytes(dataType, message.Data()); err != nil {
			return resp, fmt.Errorf("could not write message: %w", err)
		}

		return reader.NewJSONResponse(message.Data()), nil
	}
}

func (r *Reader) createWritableDataBytes(dataType string, dataMap map[string]any) ([]byte, error) {
	data := dataMap

	if data == nil {
		data = make(map[string]any)
	}

	data["type"] = dataType
	data["token"] = r.authService.SessionToken().Value()

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("could not marshal data into json: %w", err)
	}

	return []byte(fmt.Sprintf("sub %d %s", r.subID, dataBytes)), nil
}
