package auth

import (
	"time"

	"github.com/dhojayev/traderepublic-portfolio-downloader/internal"
)

var SessionRefreshTicker = time.NewTicker(internal.SessionRefreshInterval * time.Second)
