package cli

import (
	"fmt"
	"os"
)

func Version() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("v0.1")
		os.Exit(2)
	}
}
