// Package database implements the P01.04 PostgreSQL connection and migration foundation.
package database

import (
	"errors"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/failure"
)

const (
	codeConfigurationInvalid  failure.Code = "database.configuration.invalid"
	codeConnectionUnavailable failure.Code = "database.connection.unavailable"
	codeTransactionBegin      failure.Code = "database.transaction.begin_failed"
	codeTransactionCommit     failure.Code = "database.transaction.commit_failed"
	codeTransactionRollback   failure.Code = "database.transaction.rollback_failed"
	codeMigrationInvalid      failure.Code = "database.migration.invalid"
	codeMigrationLock         failure.Code = "database.migration.lock_failed"
	codeMigrationLedger       failure.Code = "database.migration.ledger_failed"
	codeMigrationDrift        failure.Code = "database.migration.drift_detected"
	codeMigrationApply        failure.Code = "database.migration.apply_failed"
)

func safeFailure(code failure.Code, category failure.Category, title string, options ...failure.Option) error {
	value, err := failure.New(code, category, title, options...)
	if err != nil {
		return errors.New("database failure could not be classified safely")
	}
	return value
}

func safeWrappedFailure(cause error, code failure.Code, category failure.Category, title string, options ...failure.Option) error {
	if cause == nil {
		return safeFailure(code, category, title, options...)
	}
	value, err := failure.Wrap(cause, code, category, title, options...)
	if err != nil {
		return errors.New("database failure could not be classified safely")
	}
	return value
}
