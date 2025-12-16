package converter

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type UnitStandard int

const (
	Binary UnitStandard = iota
	Decimal
)

type RoundMode int

const (
	RoundNearest RoundMode = iota
	RoundUp
	RoundDown
)

var (
	binaryUnits  = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	decimalUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	parseRegex   = regexp.MustCompile(`^([0-9.]+)\s*([A-Za-z]+)$`)
)

type Options struct {
	Standard  UnitStandard
	Precision int
	RoundMode RoundMode
	ForceUnit string
}

func BytesToHuman(bytes uint64, opts Options) (string, error) {
	if opts.Precision < 0 || opts.Precision > 6 {
		return "", fmt.Errorf("precision must be between 0 and 6")
	}

	base := 1024.0
	units := binaryUnits
	if opts.Standard == Decimal {
		base = 1000.0
		units = decimalUnits
	}

	if bytes == 0 {
		return fmt.Sprintf("0 %s", units[0]), nil
	}

	var unitIndex int
	if opts.ForceUnit != "" {
		unitIndex = findUnitIndex(opts.ForceUnit, units)
		if unitIndex == -1 {
			return "", fmt.Errorf("invalid unit: %s", opts.ForceUnit)
		}
	} else {
		unitIndex = int(math.Floor(math.Log(float64(bytes)) / math.Log(base)))
		if unitIndex >= len(units) {
			unitIndex = len(units) - 1
		}
	}

	value := float64(bytes) / math.Pow(base, float64(unitIndex))
	value = applyRounding(value, opts.Precision, opts.RoundMode)

	fmtStr := fmt.Sprintf("%%.%df %%s", opts.Precision)
	return fmt.Sprintf(fmtStr, value, units[unitIndex]), nil
}

func HumanToBytes(human string, opts Options) (uint64, error) {
	human = strings.TrimSpace(human)
	matches := parseRegex.FindStringSubmatch(human)
	if matches == nil {
		return 0, fmt.Errorf("invalid format: %s", human)
	}

	valueStr := matches[1]
	unit := strings.ToUpper(matches[2])

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil || value < 0 {
		return 0, fmt.Errorf("invalid number: %s", valueStr)
	}

	base := 1024.0
	units := binaryUnits
	if opts.Standard == Decimal || strings.HasSuffix(unit, "B") && !strings.HasSuffix(unit, "IB") {
		base = 1000.0
		units = decimalUnits
	}

	unitIndex := findUnitIndex(unit, units)
	if unitIndex == -1 {
		return 0, fmt.Errorf("unknown unit: %s", unit)
	}

	bytes := value * math.Pow(base, float64(unitIndex))
	return uint64(bytes), nil
}

func findUnitIndex(unit string, units []string) int {
	unit = strings.ToUpper(unit)
	for i, u := range units {
		if strings.ToUpper(u) == unit {
			return i
		}
	}
	return -1
}

func applyRounding(value float64, precision int, mode RoundMode) float64 {
	multiplier := math.Pow(10, float64(precision))
	switch mode {
	case RoundUp:
		return math.Ceil(value*multiplier) / multiplier
	case RoundDown:
		return math.Floor(value*multiplier) / multiplier
	default:
		return math.Round(value*multiplier) / multiplier
	}
}
