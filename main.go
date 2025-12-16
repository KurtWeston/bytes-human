package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/bytes-human/converter"
)

type config struct {
	toHuman   bool
	toBytes   bool
	decimal   bool
	precision int
	roundUp   bool
	roundDown bool
	forceUnit string
	jsonOut   bool
	batch     bool
}

type result struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

func main() {
	cfg := parseFlags()

	if cfg.batch {
		processBatch(cfg)
		return
	}

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Error: no input provided")
		flag.Usage()
		os.Exit(1)
	}

	input := flag.Arg(0)
	output, err := processInput(input, cfg)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if cfg.jsonOut {
		r := result{Input: input, Output: output}
		data, _ := json.MarshalIndent(r, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println(output)
	}
}

func parseFlags() config {
	cfg := config{}
	flag.BoolVar(&cfg.toHuman, "h", false, "Convert bytes to human-readable format")
	flag.BoolVar(&cfg.toBytes, "b", false, "Convert human-readable to bytes")
	flag.BoolVar(&cfg.decimal, "d", false, "Use decimal units (KB, MB) instead of binary (KiB, MiB)")
	flag.IntVar(&cfg.precision, "p", 1, "Decimal precision (0-6)")
	flag.BoolVar(&cfg.roundUp, "u", false, "Round up")
	flag.BoolVar(&cfg.roundDown, "D", false, "Round down")
	flag.StringVar(&cfg.forceUnit, "f", "", "Force specific unit (B, KiB, MiB, etc.)")
	flag.BoolVar(&cfg.jsonOut, "j", false, "Output in JSON format")
	flag.BoolVar(&cfg.batch, "B", false, "Batch mode: read from stdin")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: bytes-human [options] <value>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  bytes-human -h 1048576        # Convert bytes to human\n")
		fmt.Fprintf(os.Stderr, "  bytes-human -b \"1.5 GB\"      # Convert human to bytes\n")
		fmt.Fprintf(os.Stderr, "  bytes-human -h -d 1000000     # Use decimal units\n")
		fmt.Fprintf(os.Stderr, "  cat file.txt | bytes-human -B -h  # Batch mode\n")
	}

	flag.Parse()

	if !cfg.toHuman && !cfg.toBytes {
		cfg.toHuman = true
	}

	return cfg
}

func processInput(input string, cfg config) (string, error) {
	opts := converter.Options{
		Standard:  converter.Binary,
		Precision: cfg.precision,
		RoundMode: converter.RoundNearest,
		ForceUnit: cfg.forceUnit,
	}

	if cfg.decimal {
		opts.Standard = converter.Decimal
	}
	if cfg.roundUp {
		opts.RoundMode = converter.RoundUp
	} else if cfg.roundDown {
		opts.RoundMode = converter.RoundDown
	}

	if cfg.toBytes {
		bytes, err := converter.HumanToBytes(input, opts)
		if err != nil {
			return "", err
		}
		return strconv.FormatUint(bytes, 10), nil
	}

	bytes, err := strconv.ParseUint(strings.TrimSpace(input), 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid byte count: %s", input)
	}

	return converter.BytesToHuman(bytes, opts)
}

func processBatch(cfg config) {
	scanner := bufio.NewScanner(os.Stdin)
	results := []result{}

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		output, err := processInput(input, cfg)
		r := result{Input: input, Output: output}
		if err != nil {
			r.Error = err.Error()
		}
		results = append(results, r)

		if !cfg.jsonOut {
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: error: %v\n", input, err)
			} else {
				fmt.Println(output)
			}
		}
	}

	if cfg.jsonOut {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
	}
}
