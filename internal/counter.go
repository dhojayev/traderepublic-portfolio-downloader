package internal

import "sync/atomic"

type OperationCounter struct {
	processed *atomic.Uint64
	skipped   *atomic.Uint64
}

func NewOperationCounter() *OperationCounter {
	return &OperationCounter{
		processed: &atomic.Uint64{},
		skipped:   &atomic.Uint64{},
	}
}

func (c *OperationCounter) Processed() *atomic.Uint64 {
	return c.processed
}

func (c *OperationCounter) Skipped() *atomic.Uint64 {
	return c.skipped
}
