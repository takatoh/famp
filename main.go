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
	progVersion = "v0.5.5"
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
	opt_phase := flag.Bool("phase", false, "Output phase angle spectrum.")
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
	x, n := fft.MakeComplexData(wave.Data)

	c := fft.FFT(x, n)
	nfold := n / 2
	a, b := discreteFourierCoeff(c, nfold)
	amplitude, phase := amplitudeAndPhase(a, b, nfold)
	for k := 0; k <= nfold; k++ {
		amplitude[k] *= t2
	}
	f, t := frequencies(n, dt)

	if *opt_phase {
		if *opt_csv_output {
			printPhaseSpectrumAsCSV(f, phase)
		} else {
			printPhaseSpectrum(f, phase)
		}
	} else if *opt_csv_output {
		printResultAsCSV(t, f, a, b, amplitude, phase)
	} else {
		printResult(t, f, a, b, amplitude, phase)
	}
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
	f := fft.RFFTFreq(n, dt)
	nfold := len(f)
	t := make([]float64, nfold)
	t[0] = 0.0
	for k := 1; k < nfold; k++ {
		t[k] = 1.0 / f[k]
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
	fmt.Println("k,T,f,A,B,AMP.,PHASE")
	for k := 0; k < len(t); k++ {
		fmt.Printf("%d,%f,%f,%f,%f,%f,%f\n", k, t[k], f[k], a[k], b[k], amp[k], phase[k])
	}
}

func printPhaseSpectrum(f, phase []float64) {
	fmt.Println("   OMEGA   PHASE")
	fmt.Println("")
	for k := 0; k < len(f); k++ {
		fmt.Printf("%8.3f%8.3f\n", 2*math.Pi*f[k], phase[k])
	}
}

func printPhaseSpectrumAsCSV(f, phase []float64) {
	fmt.Println("OMEGA,PHASE")
	for k := 0; k < len(f); k++ {
		fmt.Printf("%f,%f\n", 2*math.Pi*f[k], phase[k])
	}
}
