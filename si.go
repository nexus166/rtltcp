package rtltcp

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// https://github.com/bemasher/rtltcp/blob/master/si/si.go

var (
	suffixes map[rune]float64 = map[rune]float64{
		'Y': 1e24,
		'Z': 1e21,
		'E': 1e18,
		'P': 1e15,
		'T': 1e12,
		'G': 1e9,
		'M': 1e6,
		'k': 1e3,
		'm': 1e-3,
		'u': 1e-6,
		'n': 1e-9,
		'p': 1e-12,
		'f': 1e-15,
		'a': 1e-18,
		'z': 1e-21,
		'y': 1e-24,
	}
)

type ScientificNotation float64

func (si *ScientificNotation) String() (s string) {
	return strconv.FormatFloat(float64(*si), 'g', -1, 64)
}

func (si *ScientificNotation) Set(value string) error {
	mantissaStr := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || r == '.' || r == '-' {
			return r
		}
		return -1
	}, value)
	suffix := strings.Map(func(r rune) rune {
		if _, ex := suffixes[r]; ex {
			return r
		}
		return -1
	}, value)
	if len(suffix) > 1 {
		return fmt.Errorf("suffix too long: %q", suffix)
	}
	mantissa, err := strconv.ParseFloat(mantissaStr, 64)
	if err != nil {
		return err
	}
	siNew := ScientificNotation(mantissa)
	si = &siNew
	if len(suffix) > 0 {
		if sfx, ex := suffixes[rune(suffix[0])]; ex {
			*si = ScientificNotation(sfx)
		}
	}
	return nil
}
