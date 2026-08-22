package database

import (
	"testing"
	"time"
)

func TestMigrationChecksumIsDeterministicAndContentSensitive(t *testing.T) {
	first := Migration{Version: 1, Name: "create_fixture", SQL: "CREATE TABLE example(id bigint);"}
	second := Migration{Version: 1, Name: "create_fixture", SQL: "CREATE TABLE example(id bigint);"}
	changed := Migration{Version: 1, Name: "create_fixture", SQL: "CREATE TABLE example(id text);"}

	if first.checksum() != second.checksum() {
		t.Fatal("identical migration SQL produced different checksums")
	}
	if first.checksum() == changed.checksum() {
		t.Fatal("changed migration SQL produced the same checksum")
	}
	if len(first.checksum()) != 64 {
		t.Fatalf("checksum length = %d, want 64", len(first.checksum()))
	}
}

func TestNewMigratorRejectsInvalidOwnerSequenceAndTimeout(t *testing.T) {
	tests := []struct {
		name       string
		owner      string
		migrations []Migration
		timeout    time.Duration
	}{
		{name: "bad owner", owner: "Kernel.Data", timeout: time.Second},
		{name: "zero timeout", owner: "kernel.data", timeout: 0},
		{name: "zero version", owner: "kernel.data", timeout: time.Second, migrations: []Migration{{Version: 0, Name: "bad", SQL: "SELECT 1"}}},
		{name: "duplicate version", owner: "kernel.data", timeout: time.Second, migrations: []Migration{{Version: 1, Name: "one", SQL: "SELECT 1"}, {Version: 1, Name: "two", SQL: "SELECT 2"}}},
		{name: "descending version", owner: "kernel.data", timeout: time.Second, migrations: []Migration{{Version: 2, Name: "two", SQL: "SELECT 2"}, {Version: 1, Name: "one", SQL: "SELECT 1"}}},
		{name: "bad name", owner: "kernel.data", timeout: time.Second, migrations: []Migration{{Version: 1, Name: "Bad-Name", SQL: "SELECT 1"}}},
		{name: "empty sql", owner: "kernel.data", timeout: time.Second, migrations: []Migration{{Version: 1, Name: "empty", SQL: "   "}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewMigrator(nil, test.owner, test.migrations, test.timeout)
			if err == nil {
				t.Fatal("NewMigrator() error = nil, want validation failure")
			}
			assertFailureCode(t, err, codeMigrationInvalid)
		})
	}
}

func TestAdvisoryLockKeyIsStableAndOwnerScoped(t *testing.T) {
	first := advisoryLockKey("kernel.data")
	if first != advisoryLockKey("kernel.data") {
		t.Fatal("advisory lock key is not deterministic")
	}
	if first == advisoryLockKey("kernel.other") {
		t.Fatal("distinct migration owners share an advisory lock key")
	}
}
