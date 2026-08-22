package database

import (
	"context"
	"errors"
	"time"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
	"github.com/jackc/pgx/v5"
)

const transactionRollbackTimeout = 2 * time.Second

// TxBeginner is the minimum transaction boundary implemented by pgx pools and
// acquired connections. It keeps callers independent from domain repositories.
type TxBeginner interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

// InTransaction executes fn inside one explicit PostgreSQL transaction. Callback
// errors are returned unchanged when rollback succeeds so higher layers retain
// their own governed failure identity. Begin/commit/rollback provider failures
// are mapped through the P01.03 structured failure contract.
func InTransaction(ctx context.Context, beginner TxBeginner, options pgx.TxOptions, fn func(pgx.Tx) error) error {
	if beginner == nil || fn == nil {
		return safeFailure(
			codeTransactionBegin,
			failure.CategoryValidation,
			"database transaction boundary is invalid",
		)
	}

	tx, err := beginner.BeginTx(ctx, options)
	if err != nil {
		return safeWrappedFailure(
			err,
			codeTransactionBegin,
			failure.CategoryDependency,
			"database transaction could not begin",
		)
	}

	finished := false
	defer func() {
		if finished {
			return
		}
		rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), transactionRollbackTimeout)
		defer cancel()
		_ = tx.Rollback(rollbackCtx)
	}()

	if callbackErr := fn(tx); callbackErr != nil {
		rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), transactionRollbackTimeout)
		rollbackErr := tx.Rollback(rollbackCtx)
		cancel()
		finished = true
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			return safeWrappedFailure(
				errors.Join(callbackErr, rollbackErr),
				codeTransactionRollback,
				failure.CategoryInternal,
				"database transaction rollback failed",
			)
		}
		return callbackErr
	}

	if err := tx.Commit(ctx); err != nil {
		finished = true
		return safeWrappedFailure(
			err,
			codeTransactionCommit,
			failure.CategoryDependency,
			"database transaction commit failed",
		)
	}
	finished = true
	return nil
}
