package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/takatoh/fft"
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

	wave := waves[0]
	ndata := wave.NData()
	var n int = 2
	for {
		if n >= ndata {
			break
		} else {
			n *= 2
		}
	}
	var x []complex128
	for i := 0; i < ndata; i++ {
		x = append(x, complex(wave.Data[i], 0.0))
	}
	for i := ndata; i < n; i++ {
		x = append(x, complex(0.0, 0.0))
	}

	c := fft.FFT(x, n)

	for i := 0; i < n; i++ {
		fmt.Fprintf(os.Stdout, "%v\n", c[i])
	}
	fmt.Println(len(c))
}
