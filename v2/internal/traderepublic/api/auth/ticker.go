package auth

import (
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal"
)

//nolint:gochecknoglobals
var SessionRefreshTicker = time.NewTicker(internal.SessionRefreshInterval * time.Second)
