package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/takatoh/fft"
	"github.com/takatoh/seismicwave"
)

const (
	progVersion = "v0.1.0"
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
	opt_csv_output := flag.Bool("csv-output", false, "Output as CSV.")
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
	x, n := makeData(wave.Data, ndata)

	c := fft.FFT(x, n)
	a, b := discreteFourierCoeff(c, n)
	nfold := n / 2
	amp, phi := amplitudeAndPhase(a, b, nfold)
	f, t := frequencies(nfold, wave.DT())

	if *opt_csv_output {
		fmt.Println("k,T,f,X,PHI")
		for k := 0; k <= nfold; k++ {
			fmt.Fprintf(os.Stdout, "%d,%.3f,%.1f,%.3f,%.3f\n", k, t[k], f[k], amp[k], phi[k])
		}
	} else {
		fmt.Println("    k        T       f       A       B       X     PHI")
		for k := 0; k <= nfold; k++ {
			fmt.Fprintf(os.Stdout, "%5d %8.3f%8.1f%8.3f%8.3f%8.3f%8.3f\n", k, t[k], f[k], a[k], b[k], amp[k], phi[k])
		}
	}
}

func makeData(data []float64, ndata int) ([]complex128, int) {
	var n int = 2
	for {
		if n >= ndata {
			break
		} else {
			n *= 2
		}
	}
	var x []complex128
	for k := 0; k < ndata; k++ {
		x = append(x, complex(data[k], 0.0))
	}
	for k := ndata; k < n; k++ {
		x = append(x, complex(0.0, 0.0))
	}
	return x, n
}

func discreteFourierCoeff(c []complex128, n int) ([]float64, []float64) {
	var a []float64
	var b []float64
	nfold := n / 2
	for k := 0; k <= nfold; k++ {
		a = append(a, 2.0*real(c[k]))
		b = append(b, -2.0*imag(c[k]))
	}
	b[0] = 0.0
	b[nfold] = 0.0
	return a, b
}

func amplitudeAndPhase(a []float64, b []float64, nfold int) ([]float64, []float64) {
	var amplitude []float64
	var phase []float64
	for k := 0; k <= nfold; k++ {
		amplitude = append(amplitude, math.Sqrt(math.Pow(a[k], 2.0)+math.Pow(b[k], 2.0)))
		phase = append(phase, math.Atan(-b[k]/a[k]))
	}
	return amplitude, phase
}

func frequencies(nfold int, dt float64) ([]float64, []float64) {
	var f []float64
	var t []float64
	for k := 0; k <= nfold; k++ {
		fk := float64(k) / 2.0 / dt
		f = append(f, fk)
		t = append(t, 1.0/fk)
	}
	return f, t
}
