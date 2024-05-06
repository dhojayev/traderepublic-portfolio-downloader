package main

type Args struct {
	WriteResponseFiles bool `arg:"-w,--write-responses" help:"Write api responses to file system"`
	DebugMode          bool `arg:"--debug" help:"Enable debug mode"`
	TraceMode          bool `arg:"--trace" help:"Enable trace mode"`
}
