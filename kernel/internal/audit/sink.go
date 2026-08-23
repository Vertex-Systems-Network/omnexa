package audit

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/observability"
)

// Sink is the append-only P01.11 audit transport contract. It grants no read,
// update, delete, export, identity, tenancy or authorization capability.
type Sink interface {
	Append(context.Context, Record) error
}

// Writer validates/seals records, appends them to a sink and exposes only
// payload-free transport health. It does not define retention or legal durability.
type Writer struct {
	sink      Sink
	logger    *observability.Logger
	submitted atomic.Uint64
	recorded  atomic.Uint64
	degraded  atomic.Uint64
	failed    atomic.Uint64
}

// NewWriter returns a protected audit writer over an explicit append sink.
func NewWriter(sink Sink, logger *observability.Logger) (*Writer, error) {
	if sink == nil {
		return nil, classifiedFailure(codeWriterInvalid, failure.CategoryValidation, "audit writer configuration is invalid", false)
	}
	return &Writer{sink: sink, logger: logger}, nil
}

// Write creates one immutable audit record and attempts exactly one sink append.
// Required failure returns an error. Best-effort failure returns an explicit
// degraded receipt; caller cancellation/deadline always propagates as an error.
func (writer *Writer) Write(ctx context.Context, requirement Requirement, input RecordInput) (Receipt, error) {
	if writer == nil || writer.sink == nil {
		return Receipt{}, classifiedFailure(codeWriterInvalid, failure.CategoryInvariant, "audit writer is invalid", false)
	}
	if !requirement.valid() {
		return Receipt{}, classifiedFailure(codeRequirementInvalid, failure.CategoryValidation, "audit delivery requirement is invalid", false)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx.Err() != nil {
		return Receipt{}, contextFailure(ctx.Err())
	}

	record, err := NewRecord(input)
	if err != nil {
		return Receipt{}, err
	}
	writer.submitted.Add(1)

	appendErr := appendSafely(ctx, writer.sink, record)
	if appendErr == nil {
		receipt := Receipt{RecordID: record.ID(), Status: DeliveryRecorded, Reason: DeliveryReasonNone}
		writer.recorded.Add(1)
		writer.logDelivery(ctx, requirement, record, receipt, nil)
		return receipt, nil
	}

	if ctx.Err() != nil || errors.Is(appendErr, context.Canceled) || errors.Is(appendErr, context.DeadlineExceeded) {
		writer.failed.Add(1)
		reason := DeliveryReasonCanceled
		if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(appendErr, context.DeadlineExceeded) {
			reason = DeliveryReasonDeadline
		}
		receipt := Receipt{RecordID: record.ID(), Status: DeliveryFailed, Reason: reason}
		writer.logDelivery(ctx, requirement, record, receipt, appendErr)
		cause := appendErr
		if ctx.Err() != nil {
			cause = ctx.Err()
		}
		return receipt, contextFailure(cause)
	}

	if requirement == RequirementRequired {
		writer.failed.Add(1)
		receipt := Receipt{RecordID: record.ID(), Status: DeliveryFailed, Reason: DeliveryReasonSinkFailure}
		writer.logDelivery(ctx, requirement, record, receipt, appendErr)
		return receipt, wrappedFailure(appendErr, codeSinkRequiredFailed, failure.CategoryDependency, "required audit append failed", false)
	}

	writer.degraded.Add(1)
	receipt := Receipt{RecordID: record.ID(), Status: DeliveryDegraded, Reason: DeliveryReasonSinkFailure}
	writer.logDelivery(ctx, requirement, record, receipt, appendErr)
	return receipt, nil
}

// Health returns a fixed-size, payload-free transport snapshot.
func (writer *Writer) Health() Health {
	if writer == nil {
		return Health{}
	}
	return Health{
		Submitted: writer.submitted.Load(),
		Recorded:  writer.recorded.Load(),
		Degraded:  writer.degraded.Load(),
		Failed:    writer.failed.Load(),
	}
}

func appendSafely(ctx context.Context, sink Sink, record Record) (err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("audit sink failed safely")
		}
	}()
	return sink.Append(ctx, record)
}

func (writer *Writer) logDelivery(ctx context.Context, requirement Requirement, record Record, receipt Receipt, err error) {
	if writer.logger == nil {
		return
	}
	attrs := []slog.Attr{
		slog.Any("audit_record_id", observability.Classified{
			Classification: observability.ClassificationInternal,
			Value:          record.ID(),
		}),
		slog.String("audit_delivery_status", string(receipt.Status)),
		slog.String("audit_delivery_reason", string(receipt.Reason)),
		slog.String("audit_classification", string(record.Classification())),
		slog.Bool("audit_required", requirement == RequirementRequired),
	}
	if err != nil {
		attrs = append(attrs, slog.Any("error", observability.SafeError(err)))
	}

	switch receipt.Status {
	case DeliveryRecorded:
		writer.logger.Info(ctx, "audit transport append completed", attrs...)
	case DeliveryDegraded:
		writer.logger.Warn(ctx, "audit transport append degraded", attrs...)
	case DeliveryFailed:
		writer.logger.Error(ctx, "audit transport append failed", attrs...)
	}
}
