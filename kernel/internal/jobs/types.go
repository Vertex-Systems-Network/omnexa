package jobs

import (
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	maxJobTypeRunes          = 64
	maxIdempotencyRunes      = 256
	maxRetryAttempts         = 8
	minRetryBackoff          = time.Millisecond
	maxRetryBackoff          = 30 * time.Second
	minRecurringInterval     = time.Second
	maxRecurringInterval     = 30 * 24 * time.Hour
)

var jobTypePattern = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)

// Type is a stable kernel-local job type. It carries no tenant, actor, or
// authorization meaning.
type Type string

func (jobType Type) valid() bool {
	return utf8.RuneCountInString(string(jobType)) <= maxJobTypeRunes && jobTypePattern.MatchString(string(jobType))
}

// Idempotency is an explicit duplicate-safety contract. The key is opaque and
// the fingerprint is supplied by the caller; neither value is an entity ID or
// authorization token.
type Idempotency struct {
	Key         string
	Fingerprint string
}

func (value Idempotency) valid() bool {
	return validOpaque(value.Key, maxIdempotencyRunes) && validOpaque(value.Fingerprint, maxIdempotencyRunes)
}

// RetryPolicy is a bounded deterministic retry policy. Zero value means one
// attempt with no retry.
type RetryPolicy struct {
	MaxAttempts    int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

func (policy RetryPolicy) normalize() (RetryPolicy, bool) {
	if policy == (RetryPolicy{}) {
		return RetryPolicy{MaxAttempts: 1}, true
	}
	if policy.MaxAttempts < 1 || policy.MaxAttempts > maxRetryAttempts {
		return RetryPolicy{}, false
	}
	if policy.MaxAttempts == 1 {
		if policy.InitialBackoff != 0 || policy.MaxBackoff != 0 {
			return RetryPolicy{}, false
		}
		return policy, true
	}
	if policy.InitialBackoff < minRetryBackoff || policy.InitialBackoff > maxRetryBackoff {
		return RetryPolicy{}, false
	}
	if policy.MaxBackoff < policy.InitialBackoff || policy.MaxBackoff > maxRetryBackoff {
		return RetryPolicy{}, false
	}
	return policy, true
}

func (policy RetryPolicy) backoffAfter(attempt int) time.Duration {
	if attempt < 1 || attempt >= policy.MaxAttempts || policy.InitialBackoff <= 0 {
		return 0
	}
	delay := policy.InitialBackoff
	for current := 1; current < attempt; current++ {
		if delay >= policy.MaxBackoff || delay > policy.MaxBackoff/2 {
			return policy.MaxBackoff
		}
		delay *= 2
	}
	if delay > policy.MaxBackoff {
		return policy.MaxBackoff
	}
	return delay
}

// Request is one in-memory kernel-local job submission. Payload is caller-owned
// process memory; P01.09 does not serialize or persist it.
type Request struct {
	Type        Type
	Payload     any
	Retry       RetryPolicy
	Idempotency *Idempotency
}

// Invocation is the bounded handler view of an execution attempt.
type Invocation struct {
	ExecutionID   string
	Type          Type
	Attempt       int
	Payload       any
	IdempotencyKey string
}

// State is the stable terminal execution-state vocabulary.
type State string

const (
	StateSucceeded State = "succeeded"
	StateFailed    State = "failed"
	StateCanceled  State = "canceled"
	StateDeadline  State = "deadline_exceeded"
)

// Reason is the safe machine-readable terminal reason. Raw handler errors are
// never retained in Result.
type Reason string

const (
	ReasonOK            Reason = "ok"
	ReasonHandlerFailed Reason = "handler_failed"
	ReasonCanceled      Reason = "canceled"
	ReasonDeadline      Reason = "deadline_exceeded"
)

// Result is a safe terminal job execution projection.
type Result struct {
	ExecutionID string
	Type        Type
	State       State
	Reason      Reason
	Attempts    int
	Duplicate   bool
}

// ScheduleKind is the deliberately small P01.09 schedule vocabulary.
type ScheduleKind string

const (
	ScheduleOneShot  ScheduleKind = "one_shot"
	ScheduleRecurring ScheduleKind = "recurring"
)

// Schedule defines simple in-memory kernel maintenance timing. It is not a
// durable workflow timer or cron/event fabric.
type Schedule struct {
	Kind     ScheduleKind
	StartAt  time.Time
	Interval time.Duration
}

// NewOneShot returns a UTC-normalized one-shot schedule.
func NewOneShot(at time.Time) (Schedule, error) {
	if at.IsZero() {
		return Schedule{}, classifiedFailure(codeScheduleInvalid, "validation", "job schedule is invalid", false)
	}
	return Schedule{Kind: ScheduleOneShot, StartAt: at.UTC()}, nil
}

// NewRecurring returns a bounded UTC-normalized interval schedule.
func NewRecurring(startAt time.Time, interval time.Duration) (Schedule, error) {
	if startAt.IsZero() || interval < minRecurringInterval || interval > maxRecurringInterval {
		return Schedule{}, classifiedFailure(codeScheduleInvalid, "validation", "job schedule is invalid", false)
	}
	return Schedule{Kind: ScheduleRecurring, StartAt: startAt.UTC(), Interval: interval}, nil
}

// Next returns the first scheduled instant strictly after the supplied instant.
func (schedule Schedule) Next(after time.Time) (time.Time, bool) {
	after = after.UTC()
	switch schedule.Kind {
	case ScheduleOneShot:
		if schedule.StartAt.IsZero() || !schedule.StartAt.After(after) {
			return time.Time{}, false
		}
		return schedule.StartAt.UTC(), true
	case ScheduleRecurring:
		if schedule.StartAt.IsZero() || schedule.Interval < minRecurringInterval || schedule.Interval > maxRecurringInterval {
			return time.Time{}, false
		}
		start := schedule.StartAt.UTC()
		if after.Before(start) {
			return start, true
		}
		steps := after.Sub(start)/schedule.Interval + 1
		next := start.Add(steps * schedule.Interval)
		if !next.After(after) {
			return time.Time{}, false
		}
		return next, true
	default:
		return time.Time{}, false
	}
}

func validOpaque(value string, maxRunes int) bool {
	if strings.TrimSpace(value) == "" || utf8.RuneCountInString(value) > maxRunes {
		return false
	}
	for _, char := range value {
		if unicode.IsControl(char) {
			return false
		}
	}
	return true
}
