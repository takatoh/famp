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
	progVersion = "v0.5.0"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`%s - Amplitude and Phase angle of seismic wave.

Usage:
  %s <wavefile.csv>

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
	ndata := wave.NData()
	dt := wave.DT()
	t2 := wave.Length() / 2.0
	x, n := makeData(wave.Data, ndata)

	c := fft.FFT(x, n)
	nfold := n / 2
	a, b := discreteFourierCoeff(c, nfold)
	amplitude, phase := amplitudeAndPhase(a, b, nfold)
	for k := 0; k <= nfold; k++ {
		amplitude[k] = amplitude[k] * t2
	}
	f, t := frequencies(ndata, dt)

	if *opt_csv_output {
		printResultAsCSV(t, f, amplitude, phase)
	} else {
		printResult(t, f, a, b, amplitude, phase)
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

func discreteFourierCoeff(c []complex128, nfold int) ([]float64, []float64) {
	var a []float64
	var b []float64
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
		xk := math.Sqrt(a[k]*a[k] + b[k]*b[k])
		amplitude = append(amplitude, xk)
		phase = append(phase, math.Atan2(-b[k], a[k]))
	}
	return amplitude, phase
}

func frequencies(n int, dt float64) ([]float64, []float64) {
	var f []float64
	var t []float64
	f = append(f, 0.0)
	t = append(t, 0.0)
	ndt := float64(n) * dt
	for k := 1; k <= n/2; k++ {
		fk := float64(k) / ndt
		f = append(f, fk)
		t = append(t, 1.0/fk)
	}
	return f, t
}

func printResult(t, f, a, b, x, phi []float64) {
	fmt.Println("    k        T       f       A       B    AMP.   PHASE")
	fmt.Println("")
	for k := 0; k < len(t); k++ {
		fmt.Printf("%5d %8.3f%8.3f%8.3f%8.3f%8.3f%8.3f\n", k, t[k], f[k], a[k], b[k], x[k], phi[k])
	}
}

func printResultAsCSV(t, f, x, phi []float64) {
	fmt.Println("k,T,f,AMP.,PHASE")
	for k := 0; k < len(t); k++ {
		fmt.Printf("%d,%f,%f,%f,%f\n", k, t[k], f[k], x[k], phi[k])
	}
}
