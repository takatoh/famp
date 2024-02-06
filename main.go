package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/takatoh/fft"
	"github.com/takatoh/seismicwave"
)

const (
	progName    = "famp"
	progVersion = "v0.5.3"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`%s - Fourier AMplitude and Phase angle of seismic wave.

Usage:
  %s [options] <wavefile.csv>

Options:
`, progName, progName)
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
	dt := wave.DT()
	t2 := wave.Length() / 2.0
	x, n := makeData(wave.Data)

	c := fft.FFT(x, n)
	nfold := n / 2
	a, b := discreteFourierCoeff(c, nfold)
	amplitude, phase := amplitudeAndPhase(a, b, nfold)
	for k := 0; k <= nfold; k++ {
		amplitude[k] *= t2
	}
	f, t := frequencies(n, dt)

	if *opt_csv_output {
		printResultAsCSV(t, f, a, b, amplitude, phase)
	} else {
		printResult(t, f, a, b, amplitude, phase)
	}
}

func makeData(data []float64) ([]complex128, int) {
	ndata := len(data)
	n := 2
	for n < ndata {
		n *= 2
	}
	x := make([]complex128, n)
	for k := 0; k < ndata; k++ {
		x[k] = complex(data[k], 0.0)
	}
	for k := ndata; k < n; k++ {
		x[k] = complex(0.0, 0.0)
	}
	return x, n
}

func discreteFourierCoeff(c []complex128, nfold int) ([]float64, []float64) {
	a := make([]float64, nfold+1)
	b := make([]float64, nfold+1)
	for k := 0; k <= nfold; k++ {
		a[k] = 2.0 * real(c[k])
		b[k] = -2.0 * imag(c[k])
	}
	b[0] = 0.0
	b[nfold] = 0.0
	return a, b
}

func amplitudeAndPhase(a []float64, b []float64, nfold int) ([]float64, []float64) {
	amplitude := make([]float64, nfold+1)
	phase := make([]float64, nfold+1)
	for k := 0; k <= nfold; k++ {
		amplitude[k] = math.Sqrt(a[k]*a[k] + b[k]*b[k])
		phase[k] = math.Atan2(-b[k], a[k])
	}
	return amplitude, phase
}

func frequencies(n int, dt float64) ([]float64, []float64) {
	nfold := n / 2
	f := make([]float64, nfold+1)
	t := make([]float64, nfold+1)
	f[0] = 0.0
	t[0] = 0.0
	ndt := float64(n) * dt
	for k := 1; k <= n/2; k++ {
		fk := float64(k) / ndt
		f[k] = fk
		t[k] = 1.0 / fk
	}
	return f, t
}

func printResult(t, f, a, b, amp, phase []float64) {
	fmt.Println("    k        T       f       A       B    AMP.   PHASE")
	fmt.Println("")
	for k := 0; k < len(t); k++ {
		fmt.Printf("%5d %8.3f%8.3f%8.3f%8.3f%8.3f%8.3f\n", k, t[k], f[k], a[k], b[k], amp[k], phase[k])
	}
}

func printResultAsCSV(t, f, a, b, amp, phase []float64) {
	fmt.Println("k,T,f,AMP.,PHASE")
	for k := 0; k < len(t); k++ {
		fmt.Printf("%d,%f,%f,%f,%f,%f,%f\n", k, t[k], f[k], a[k], b[k], amp[k], phase[k])
	}
}
