// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package websocket

import (
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/console"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/reader"
	"github.com/dhojayev/traderepublic-portfolio-downloader/internal/writer"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

// Injectors from wire.go:

func ProvideReader(responseWriter writer.Interface, logger *logrus.Logger) (*Reader, error) {
	client := api.NewClient(logger)
	authClient, err := auth.NewClient(client, logger)
	if err != nil {
		return nil, err
	}
	authService := console.NewAuthService(authClient)
	reader, err := NewReader(authService, responseWriter, logger)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

// wire.go:

var DefaultSet = wire.NewSet(
	NewReader, wire.Bind(new(reader.Interface), new(*Reader)),
)
