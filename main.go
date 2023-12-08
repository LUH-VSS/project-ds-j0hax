package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/LUH-VSS/project-ds-j0hax/reader"
)

func main() {

	readerCmd := flag.NewFlagSet("reader", flag.ExitOnError)
	readerHost := readerCmd.String("host", "http://127.0.0.1:1831/words", "The host to send trend data to")

	if len(os.Args) < 2 {
		fmt.Println("expected 'reader' or 'writer' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "reader":
		readerCmd.Parse(os.Args[2:])
		r := reader.NewReader()
	case "writer":
	}
}
