package main

import (
	"fmt"
	"os"

	"github.com/dlbarduzzi/scopehouse"
)

func main() {
	app := scopehouse.New()

	if err := app.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "[error] %s\n", err)
		os.Exit(1)
	}
}
