package main

import "github.com/dhojayev/traderepublic-portfolio-downloader/internal/api/auth"

type Args struct {
	PhoneNumber auth.PhoneNumber `arg:"positional,required" help:"Phone number in international format: +49xxxxxxxxxxxxx"`

	WriteResponseFiles bool `arg:"-w,--write-responses" help:"Write api responses to file system"`
	DebugMode          bool `arg:"--debug" help:"Enable debug mode"`
	TraceMode          bool `arg:"--trace" help:"Enable trace mode"`
}
