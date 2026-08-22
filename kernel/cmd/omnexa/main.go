// Command omnexa is the minimal P01.01 kernel process skeleton.
package main

import (
	"fmt"

	"github.com/Vertex-Systems-Network/omnexa/kernel/internal/buildinfo"
)

func main() {
	fmt.Printf("omnexa-kernel %s\n", buildinfo.Current())
}
