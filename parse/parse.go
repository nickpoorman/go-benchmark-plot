package parse

import (
	"regexp"
	"strconv"

	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"golang.org/x/tools/benchmark/parse"
)

// Regex: Benchmark[Function_name]_[Function_argument](b *testing.B)
var re *regexp.Regexp = regexp.MustCompile(`Benchmark([a-zA-Z0-9]+)_([_a-zA-Z0-9]+)-([0-9]+)$`)

type ArgSet map[float64]float64
type NameSet map[string]ArgSet

// parseName parses name from the benchmark
func parseName(line string) (name string, arg string, err error) {
	arr := re.FindStringSubmatch(line)

	// we expect 4 columns
	if len(arr) != 4 {
		return "", "", errors.New("Problem parsing benchmarks")
	}

	return arr[1], arr[2], nil
}

func ParseBenchmarks(useMillis bool) NameSet {
	nsPerOpMultiplier := 1.0
	if useMillis {
		nsPerOpMultiplier = 0.000001
	}
	benchResults := make(NameSet)
	scan := bufio.NewScanner(os.Stdin)
	green := color.New(color.FgGreen).SprintfFunc()
	red := color.New(color.FgRed).SprintFunc()
	for scan.Scan() {
		line := scan.Text()

		mark := green("√")

		b, err := parse.ParseLine(line)
		if err != nil {
			mark = red("?")
		}

		// read bench name and arguments
		if b != nil {
			name, arg, err := parseName(b.Name)
			if err != nil {
				mark = red("!")
				fmt.Printf("%s %s\n", mark, line)
				continue
			}

			floatArg, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				panic(err)
			}

			if _, ok := benchResults[name]; !ok {
				benchResults[name] = make(ArgSet)
			}

			benchResults[name][floatArg] = b.NsPerOp * nsPerOpMultiplier
		}

		fmt.Printf("%s %s\n", mark, line)
	}

	if err := scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading stdin: %v", err)
		os.Exit(1)
	}

	if len(benchResults) == 0 {
		fmt.Fprintf(os.Stderr, "no data.\n\n")
		os.Exit(1)
	}

	fmt.Println()

	return benchResults
}
