package datefmt

import (
	"strings"
	"sync"
	"time"
)

var (
	defaultFormatter = FastFormat
)

// Use replaces the default formatter
func Use(formatter func(t time.Time, generalLayout string) string) {
	defaultFormatter = formatter
}

// Format is a general layout based version of time.Format
func Format(t time.Time, generalLayout string) string {
	return defaultFormatter(t, generalLayout)
}

// Parse is a general layout based version of time.Parse
func Parse(generalLayout, value string) (time.Time, error) {
	return time.Parse(GoLayout(generalLayout), value)
}

// ParseInLocation is a general layout based version of time.ParseInLocation
func ParseInLocation(generalLayout, value string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(GoLayout(generalLayout), value, loc)
}

// GoLayout returns a go-style layout according to the general layout defined by the argument.
func GoLayout(generalLayout string) string {
	v, ok := layoutCache.Load(generalLayout)
	if ok {
		return v.(string)
	}
	v, _ = layoutCache.LoadOrStore(generalLayout, goLayout(generalLayout))
	return v.(string)
}

var layoutCache sync.Map

func goLayout(generalLayout string) string {
	var (
		l   = []byte(generalLayout)
		n   = len(l)
		sb  = strings.Builder{}
		max = n + 20
	)
	sb.Grow(max)
	for i := 0; i < n; i++ {
		if !tokens[l[i]] {
			sb.WriteByte(l[i])
			continue
		}
		// quote
		if l[i] == '\'' {
			for i++; i < n; i++ {
				if l[i] == '\'' {
					if l[i-1] == '\'' {
						// real quote
						sb.WriteByte('\'')
						break
					} else if i < n-1 && l[i+1] == '\'' {
						// real quote
						sb.WriteByte('\'')
						i++
						continue
					} else {
						// end of text
						break
					}
				}
				// text delimiter
				sb.WriteByte(l[i])
			}
			continue
		}
		// find consecutive tokens
		s, e, token := i, i, l[i]
		for i++; i < n; i++ {
			if l[i] != token {
				i--
				break
			}
			e = i
		}
		sb.WriteString(getPlaceholder(l[s : e+1]))
	}
	return sb.String()
}

var (
	tokens = map[byte]bool{
		'G': true,
		'y': true,
		'Y': true,
		'M': true,
		// 'w': true,
		// 'W': true,
		'D': true,
		'd': true,
		// 'F': true,
		'E': true,
		// 'u': true,
		'a': true,
		'H': true,
		// 'k': true,
		// 'K': true,
		'h': true,
		'm': true,
		's': true,
		'S': true,
		'z': true,
		'Z': true,
		'X': true,

		'\'': true,
	}

	placeholders = map[string]string{
		"YYYY": "2006",
		"yyyy": "2006",
		"YY":   "06",
		"yy":   "06",
		"MMMM": "January",
		"MMM":  "Jan",
		"MM":   "01",
		"M":    "1",
		"DDD":  "002",
		"dd":   "02",
		"d":    "2",

		"EEEE": "Monday",
		"EEE":  "Mon",

		"HH":  "15",
		"hh":  "03",
		"h":   "3",
		"mm":  "04",
		"m":   "4",
		"ss":  "05",
		"s":   "5",
		"SSS": "000",

		"a": "PM",

		"z":   "MST",
		"Z":   "-0700",
		"X":   "Z07",
		"XX":  "Z0700",
		"XXX": "Z07:00",
	}
)

func getPlaceholder(ph []byte) string {
	tmp := ph
	if len(tmp) > 4 {
		tmp = tmp[:4]
	}
	if goPh, ok := placeholders[string(tmp)]; ok {
		return goPh
	}
	return string(ph) // Do not modify
}
