package main

import (
	"fmt"
	"os"
)

const usage = `usage:
  drop            wait for files from a nearby device
  drop <files>    send files to a nearby device`

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "drop:", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println(usage)
	return nil
}
