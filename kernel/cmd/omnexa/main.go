// Command omnexa is the Omnexa kernel process and developer CLI entrypoint.
package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/buildinfo"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/developer"
)

func main() {
	os.Exit(runCommand(context.Background(), os.Args[1:], os.Stdout, os.Stderr, os.Environ()))
}

func runCommand(ctx context.Context, arguments []string, stdout io.Writer, stderr io.Writer, environ []string) int {
	if len(arguments) == 0 {
		return run(stdout, stderr, environ)
	}
	return developer.Run(ctx, developer.Options{
		Arguments:   arguments,
		Environment: environ,
		Stdout:      stdout,
		Stderr:      stderr,
	})
}

func run(stdout io.Writer, stderr io.Writer, environ []string) int {
	if _, err := config.LoadApplication(config.OSOptions(environ)); err != nil {
		if _, writeErr := fmt.Fprintf(stderr, "omnexa configuration error: %v\n", err); writeErr != nil {
			return 1
		}
		return 2
	}

	if _, err := fmt.Fprintf(stdout, "omnexa-kernel %s\n", buildinfo.Current()); err != nil {
		return 1
	}
	return 0
}
