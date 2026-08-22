// Package buildinfo exposes deterministic source/build identity for the Omnexa kernel.
package buildinfo

import "fmt"

// Version and Commit are intentionally variables so release/build tooling can
// inject source identity through -ldflags without introducing runtime config.
var (
	Version = "dev"
	Commit  = "unknown"
)

// Info is the bounded build identity exposed by the P01.01 process skeleton.
type Info struct {
	Version string
	Commit  string
}

// Current returns normalized build identity without local machine metadata.
func Current() Info {
	return Info{
		Version: normalize(Version, "dev"),
		Commit:  normalize(Commit, "unknown"),
	}
}

// String returns a stable human-readable representation suitable for smoke
// verification. It intentionally contains no timestamps, usernames or paths.
func (i Info) String() string {
	return fmt.Sprintf("version=%s commit=%s", i.Version, i.Commit)
}

func normalize(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
