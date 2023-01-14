package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/takatoh/seismicwave"
)

const (
	progVersion = "v0.0.0"
)

func main() {
	progName := filepath.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`Usage:
  %s <wavefile.csv>

Options:
`, progName)
		flag.PrintDefaults()
	}
	opt_version := flag.Bool("version", false, "Show version.")
	flag.Parse()

	if *opt_version {
		fmt.Println(progVersion)
		os.Exit(0)
	}

	filename := flag.Arg(0)

	var waves []*seismicwave.Wave
	var err error
	waves, err = seismicwave.LoadCSV(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	var wave *seismicwave.Wave = waves[0]
	fmt.Fprintf(os.Stdout, "NAME:   %s\n", wave.Name)
	fmt.Fprintf(os.Stdout, "NDATA:  %d\n", wave.NData())
	fmt.Fprintf(os.Stdout, "DT:     %f\n", wave.DT())
	fmt.Fprintf(os.Stdout, "LENGTH: %f\n", wave.Length())
}
