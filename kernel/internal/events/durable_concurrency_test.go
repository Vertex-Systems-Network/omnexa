package events

import (
	"context"
	"testing"
)

func TestDurableConsumerConcurrentSameScopeCheckpointRaceIsFailClosed(t *testing.T) {
	ready := make(chan struct{}, 2)
	release := make(chan struct{})
	registry := durableTestRegistry(t, false, func(context.Context, Envelope) error {
		ready <- struct{}{}
		<-release
		return nil
	})
	store := NewMemoryCheckpointStore()
	binding := durableTestBinding("")
	consumer := durableTestConsumer(t, registry, store, binding)
	first := testEnvelope(t)
	second := testEnvelope(t)

	errorsCh := make(chan error, 2)
	go func() {
		errorsCh <- consumer.Process(context.Background(), Delivery{Position: 1, Envelope: first})
	}()
	go func() {
		errorsCh <- consumer.Process(context.Background(), Delivery{Position: 1, Envelope: second})
	}()

	<-ready
	<-ready
	close(release)

	firstErr := <-errorsCh
	secondErr := <-errorsCh
	if firstErr == nil && secondErr == nil {
		t.Fatal("concurrent same-position deliveries both advanced the checkpoint")
	}
	if firstErr != nil && secondErr != nil {
		t.Fatalf("concurrent same-position deliveries both failed: first=%v second=%v", firstErr, secondErr)
	}
	if firstErr != nil {
		assertFailureCode(t, firstErr, codeCheckpointStale)
	} else {
		assertFailureCode(t, secondErr, codeCheckpointStale)
	}

	checkpoint, exists, err := consumer.LastCheckpoint(context.Background())
	if err != nil {
		t.Fatalf("LastCheckpoint() error = %v", err)
	}
	if !exists || checkpoint.Position != 1 {
		t.Fatalf("checkpoint after concurrent race = %#v exists=%v", checkpoint, exists)
	}
	if checkpoint.EventID != first.ID && checkpoint.EventID != second.ID {
		t.Fatalf("checkpoint event ID %q does not match either raced delivery", checkpoint.EventID)
	}

	position, err := consumer.ResumePosition(context.Background())
	if err != nil {
		t.Fatalf("ResumePosition() error = %v", err)
	}
	if position != 2 {
		t.Fatalf("resume position after concurrent race = %d, want 2", position)
	}
}
