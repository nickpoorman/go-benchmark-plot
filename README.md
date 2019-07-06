# go-benchmark-plot

Plot Go benchmarks on a graph.

## To Use

Run your benchmarks and then pipe them into benchplot.

```
make && cat bench-2019-07-06-10:38.txt | ./bin/benchplot
```

## Help

```
./bin/benchplot --help
Usage of ./bin/benchplot:
  -format string
        Image extension to be used. Supported extensions are: .eps, .jpg, .jpeg, .pdf, .png, .svg, .tif and .tiff (default "svg")
  -xPredictMultiplier float
        Multiplier used to predict values beyond max benchmarked arg (default 1.7)
```

### Note

Check out [benchstat](https://godoc.org/golang.org/x/perf/cmd/benchstat) if you simply want to compare your benchmarks without a graph.

## License

(c) 2019 Nick Poorman. Licensed under the Apache License, Version 2.0.
