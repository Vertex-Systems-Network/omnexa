package audit

import (
	"context"
	"sync"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	minMemorySinkCapacity = 1
	maxMemorySinkCapacity = 4096
)

// MemorySink is a bounded deterministic local/test append sink. It is not a
// durable production store and intentionally exposes no mutation of appended records.
type MemorySink struct {
	mu       sync.RWMutex
	capacity int
	records  []Record
}

// NewMemorySink creates a bounded process-local sink for tests and local tooling.
func NewMemorySink(capacity int) (*MemorySink, error) {
	if capacity < minMemorySinkCapacity || capacity > maxMemorySinkCapacity {
		return nil, classifiedFailure(codeSinkInvalid, failure.CategoryValidation, "audit memory sink configuration is invalid", false)
	}
	return &MemorySink{capacity: capacity, records: make([]Record, 0, capacity)}, nil
}

// Append validates record integrity and appends one defensive immutable copy.
// Capacity exhaustion fails closed instead of evicting or overwriting history.
func (sink *MemorySink) Append(ctx context.Context, record Record) error {
	if sink == nil || sink.capacity < minMemorySinkCapacity || sink.capacity > maxMemorySinkCapacity {
		return classifiedFailure(codeSinkInvalid, failure.CategoryInvariant, "audit memory sink is invalid", false)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err := record.Verify(); err != nil {
		return err
	}

	sink.mu.Lock()
	defer sink.mu.Unlock()
	if len(sink.records) >= sink.capacity {
		return classifiedFailure(codeSinkFull, failure.CategoryUnavailable, "audit memory sink capacity is exhausted", false)
	}
	sink.records = append(sink.records, cloneRecord(record))
	return nil
}

// Snapshot returns a defensive local/test projection in append order. Snapshot
// is deliberately absent from the Sink interface so write authority never implies read/export authority.
func (sink *MemorySink) Snapshot() []Record {
	if sink == nil {
		return nil
	}
	sink.mu.RLock()
	defer sink.mu.RUnlock()
	result := make([]Record, len(sink.records))
	for index, record := range sink.records {
		result[index] = cloneRecord(record)
	}
	return result
}

// Len returns only local/test sink occupancy for deterministic assertions.
func (sink *MemorySink) Len() int {
	if sink == nil {
		return 0
	}
	sink.mu.RLock()
	defer sink.mu.RUnlock()
	return len(sink.records)
}
