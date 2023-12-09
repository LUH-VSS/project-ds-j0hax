package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/LUH-VSS/project-ds-j0hax/mapper"
	"github.com/LUH-VSS/project-ds-j0hax/reducer"
)

func main() {

	readerCmd := flag.NewFlagSet("map", flag.ExitOnError)
	readerHost := readerCmd.String("host", "http://127.0.0.1:1831/words", "The host to send trend data to")

	writerCmd := flag.NewFlagSet("reduce", flag.ExitOnError)
	writerAddr := readerCmd.String("addr", ":1831", "The address to bind to")
	writerPattern := readerCmd.String("pattern", "/words", "The pattern to listen to")

	if len(os.Args) < 2 {
		fmt.Println("expected 'reader' or 'writer' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "map":
		readerCmd.Parse(os.Args[2:])
		r := mapper.NewReader(*readerHost, readerCmd.Args())
		r.Run()
	case "reduce":
		writerCmd.Parse(os.Args[2:])
		w := reducer.NewWriter(*writerAddr, *writerPattern)
		w.Run()
	}
}
