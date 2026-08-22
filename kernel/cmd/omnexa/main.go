// Command omnexa is the minimal Omnexa kernel process entrypoint.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/buildinfo"
	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/config"
)

func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.Environ()))
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
