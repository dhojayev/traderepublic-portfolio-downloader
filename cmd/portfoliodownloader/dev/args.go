package main

import "github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"

type Args struct {
	PhoneNumber auth.PhoneNumber `arg:"positional,required" help:"Phone number in international format: +49xxxxxxxxxxxxx"`

	LocalMode bool `arg:"-l,--local" help:"Read json file from local fs instead of calling trade republic api"`
	TraceMode bool `arg:"--trace" help:"Enable trace mode"`
}
