package buildinfo

import "testing"

func TestCurrentUsesDefaults(t *testing.T) {
	oldVersion, oldCommit := Version, Commit
	t.Cleanup(func() {
		Version, Commit = oldVersion, oldCommit
	})

	Version = ""
	Commit = ""

	got := Current()
	if got.Version != "dev" {
		t.Fatalf("Version = %q, want dev", got.Version)
	}
	if got.Commit != "unknown" {
		t.Fatalf("Commit = %q, want unknown", got.Commit)
	}
}

func TestInfoString(t *testing.T) {
	got := (Info{Version: "1.2.3", Commit: "abc123"}).String()
	want := "version=1.2.3 commit=abc123"
	if got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
