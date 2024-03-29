package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/LUH-VSS/project-ds-j0hax/mapper"
	"github.com/LUH-VSS/project-ds-j0hax/reducer"
)

func subCmdError() {
	outfile := flag.CommandLine.Output()
	fmt.Fprintln(outfile, "Error: must use either `map` or `reduce` subcommand")
	os.Exit(1)
}

func main() {

	readerCmd := flag.NewFlagSet("map", flag.ExitOnError)
	readerHost := readerCmd.String("host", "127.0.0.1:1831", "The host(s) to send trend data to. These can be comma-seperated to send to multiple hosts.")
	readerCmd.Usage = func() {
		outfile := flag.CommandLine.Output()
		fmt.Fprintf(outfile, "Usage: %s [OPTION] [FILE]...\n", os.Args[0])
		fmt.Fprintln(outfile, "Extract words from FILE(s) and send these to HOST.")
		readerCmd.PrintDefaults()
	}

	writerCmd := flag.NewFlagSet("reduce", flag.ExitOnError)
	writerAddr := writerCmd.String("addr", ":1831", "The address to bind to")

	if len(os.Args) <= 1 {
		subCmdError()
	}

	switch os.Args[1] {
	case "map":
		readerCmd.Parse(os.Args[2:])
		hosts := strings.Split(*readerHost, ",")
		r := mapper.NewReader(hosts, readerCmd.Args())
		r.Run()
	case "reduce":
		writerCmd.Parse(os.Args[2:])
		w := reducer.NewWriter(*writerAddr, readerCmd.Arg(0))
		w.Run()
	default:
		subCmdError()
	}
}
