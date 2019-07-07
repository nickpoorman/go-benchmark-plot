package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/nickpoorman/go-benchmark-plot/parse"
	"go-hep.org/x/hep/fit"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func formula(degree int) string {
	char := 'a'
	asciiA := int(char)

	f := "f(x) = "
	for i := 0; i <= degree+1; i++ {
		if i == 0 {
			f += string(asciiA + i)
		} else if i == 1 {
			f = fmt.Sprintf("%s + %s*x", f, string(asciiA+i))
		} else {
			f = fmt.Sprintf("%s + %s*x^%d", f, string(asciiA+i), i)
		}
	}
	return f
}

func formulaCoef(ps []float64) string {
	f := "f(x) = "
	for i := 0; i < len(ps); i++ {
		if i == 0 {
			f = fmt.Sprintf("%f", ps[i])
		} else if i == 1 {
			f = fmt.Sprintf("%s + %f*x", f, ps[i])
		} else {
			f = fmt.Sprintf("%s + %f*x^%d", f, ps[i], i)
		}
	}
	return f
}

func CreateGrapth(benchName string, xdata, ydata []float64, format string, xPredictMultiplier float64, degree int) {
	// This is our function to minimize.
	// ps is the slice of parameters to optimize during the fit.
	poly := func(x float64, ps []float64) float64 {
		// return ps[0] + ps[1]*x + ps[2]*x*x
		var f float64
		for i := 0; i <= degree+1; i++ {
			f += ps[i] * math.Pow(x, float64(i))
		}
		return f
	}

	ps := make([]float64, degree+2)
	for i := range ps {
		ps[i] = 1.0
	}

	res, err := fit.Curve1D(
		fit.Func1D{
			// F is the function to minimize.
			// ps is the slice of parameters to optimize during the fit.
			F: poly,
			X: xdata,
			Y: ydata,
			// Ps is the initial values for the parameters.
			// If Ps is nil, the set of initial parameters values is a slice of
			// length N filled with zeros.
			Ps: ps,
		},
		nil, &optimize.NelderMead{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		log.Fatal(err)
	}

	{
		xMax := mat.Max(mat.NewVecDense(len(xdata), xdata)) * xPredictMultiplier

		p := hplot.New()
		p.X.Label.Text = fmt.Sprintf("%s\n%s", formula(degree), formulaCoef(res.X))
		p.Y.Label.Text = "ms/op"

		s := hplot.NewS2D(hplot.ZipXY(xdata, ydata))
		s.Color = color.RGBA{0, 0, 255, 255}
		p.Add(s)

		f := plotter.NewFunction(func(x float64) float64 {
			return poly(x, res.X)
		})
		f.Color = color.RGBA{255, 0, 0, 255}
		f.Samples = 1000
		f.XMax = xMax
		p.Add(f)

		p.Add(plotter.NewGrid())

		p.X.Max = xMax
		p.Y.Max = poly(xMax, res.X)

		err := p.Save(20*vg.Centimeter, -1, fmt.Sprintf("bench-%d-%s.%s", time.Now().Unix(), benchName, format))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	var format = flag.String("format", "svg", "Image extension to be used. Supported extensions are: .eps, .jpg, .jpeg, .pdf, .png, .svg, .tif and .tiff")
	var xPredictMultiplier = flag.Float64("xPredictMultiplier", 1.7, "Multiplier used to predict values beyond max benchmarked arg")
	var useMillis = flag.Bool("useMillis", true, "Use milliseconds when writing the ns/op benchmarks")
	var degree = flag.Int("degree", 1, "The degree for the polynomial function to minimize")
	flag.Parse()

	results := parse.ParseBenchmarks(*useMillis)

	for benchName, resultSet := range results {
		xdata := make([]float64, len(resultSet))
		ydata := make([]float64, len(resultSet))
		for arg, opsNs := range resultSet {
			xdata = append(xdata, arg)
			ydata = append(ydata, opsNs)
		}
		CreateGrapth(benchName, xdata, ydata, *format, *xPredictMultiplier, *degree)
	}
}
