package main

type Args struct {
	WriteResponseFiles bool `arg:"-w,--write-responses" help:"write API responses to the file system"`
	DebugMode          bool `arg:"--debug" help:"enable debug mode"`
	TraceMode          bool `arg:"--trace" help:"enable trace mode"`
}
