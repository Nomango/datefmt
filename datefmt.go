package datefmt

import (
	"strings"
	"sync"
	"time"
)

// Format returns a textual representation of the time value formatted according to the general layout.
func Format(t time.Time, generalLayout string) string {
	return t.Format(GoLayout(generalLayout))
}

// GoLayout returns a go-style layout based on general layout.
func GoLayout(generalLayout string) string {
	v, ok := layoutCache.Load(generalLayout)
	if ok {
		return v.(string)
	}
	v, _ = layoutCache.LoadOrStore(generalLayout, goLayout(generalLayout))
	return v.(string)
}

var layoutCache sync.Map

func goLayout(generalLayout string) (goLayout string) {
	var (
		s  state
		sb strings.Builder
		q  bool // quote
		lt rune // last token
	)
	for _, token := range generalLayout {
		if s.token != 0 && token != s.token {
			sb.WriteString(getToken(s))
			s = state{}
		}
		if token == '\'' {
			if lt == '\'' {
				// write a single quotation
				sb.WriteRune('\'')
				lt = 0
				q = !q
				continue
			} else {
				// begin or end of quote
				q = !q
			}
		} else {
			if q {
				// in quote
				sb.WriteRune(token)
			} else if tokens[token] {
				s.token = token
				s.placeholder += string(token)
			} else {
				sb.WriteRune(token)
			}
		}
		lt = token
	}
	if s.token != 0 {
		sb.WriteString(getToken(s))
	}
	return sb.String()
}

type state struct {
	token       rune
	placeholder string
}

var (
	tokens = map[rune]bool{
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

func getToken(s state) string {
	ph := s.placeholder
	if len(ph) > 4 {
		ph = ph[:4]
	}
	if goPh, ok := placeholders[ph]; ok {
		return goPh
	}
	return s.placeholder // Do not modify
}
