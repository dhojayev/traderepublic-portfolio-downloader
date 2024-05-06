package main

type Args struct {
	LocalMode bool `arg:"-l,--local" help:"Read json file from local fs instead of calling trade republic api"`
	TraceMode bool `arg:"--trace" help:"Enable trace mode"`
}
