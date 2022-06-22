package datefmt

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

type formatArg = func(t *time.Time) string

type layout struct {
	args []interface{}
	max  int
}

func FastFormat(t time.Time, generalLayout string) string {
	fl := FastLayout(generalLayout)
	var sb strings.Builder
	sb.Grow(fl.max)
	for _, arg := range fl.args {
		switch v := arg.(type) {
		case string:
			sb.WriteString(v)
		case formatArg:
			sb.WriteString(v(&t))
		}
	}
	return sb.String()
}

var fastLayoutCache sync.Map

func FastLayout(generalLayout string) *layout {
	v, ok := fastLayoutCache.Load(generalLayout)
	if ok {
		return v.(*layout)
	}
	v, _ = fastLayoutCache.LoadOrStore(generalLayout, fastLayout(generalLayout))
	return v.(*layout)
}

func fastLayout(generalLayout string) *layout {
	var (
		l  = layout{}
		gl = []byte(generalLayout)
		n  = len(gl)
		sb = strings.Builder{}
	)
	for i := 0; i < n; i++ {
		if !tokens[gl[i]] {
			sb.WriteByte(gl[i])
			continue
		}
		// quote
		if gl[i] == '\'' {
			for i++; i < n; i++ {
				if gl[i] == '\'' {
					if gl[i-1] == '\'' {
						// real quote
						sb.WriteByte('\'')
						break
					} else if i < n-1 && gl[i+1] == '\'' {
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
				sb.WriteByte(gl[i])
			}
			continue
		}
		// flush buffer
		if sb.Len() > 0 {
			l.args = append(l.args, sb.String())
			sb.Reset()
		}
		// find consecutive tokens
		s, e, token := i, i, gl[i]
		for i++; i < n; i++ {
			if gl[i] != token {
				i--
				break
			}
			e = i
		}
		l.args = append(l.args, getFormatArg(gl[s:e+1]))
	}
	if sb.Len() > 0 {
		l.args = append(l.args, sb.String())
	}
	return &l
}

var (
	formatArgs = map[string]formatArg{
		"YYYY": func(t *time.Time) string { return strconv.Itoa(t.Year()) },
		"yyyy": func(t *time.Time) string { return strconv.Itoa(t.Year()) },
		"YY":   func(t *time.Time) string { return strconv.Itoa(t.Year())[2:4] },
		"yy":   func(t *time.Time) string { return strconv.Itoa(t.Year())[2:4] },
		"MMMM": func(t *time.Time) string { return t.Month().String() },
		"MMM":  func(t *time.Time) string { return t.Month().String()[:3] },
		"MM":   func(t *time.Time) string { return formatZero2(int(t.Month())) },
		"M":    func(t *time.Time) string { return formatMax99(int(t.Month())) },
		"DDD":  func(t *time.Time) string { return formatZero3(t.YearDay()) },
		"dd":   func(t *time.Time) string { return formatZero2(t.Day()) },
		"d":    func(t *time.Time) string { return formatMax99(t.Day()) },
		"EEEE": func(t *time.Time) string { return t.Weekday().String() },
		"EEE":  func(t *time.Time) string { return t.Weekday().String()[:3] },
		"HH":   func(t *time.Time) string { return formatZero2(t.Hour()) },
		"hh":   func(t *time.Time) string { return "TODO" },
		"H":    func(t *time.Time) string { return formatMax99(t.Hour()) },
		"h":    func(t *time.Time) string { return "TODO" },
		"mm":   func(t *time.Time) string { return formatZero2(t.Minute()) },
		"m":    func(t *time.Time) string { return formatMax99(t.Minute()) },
		"ss":   func(t *time.Time) string { return formatZero2(t.Second()) },
		"s":    func(t *time.Time) string { return formatMax99(t.Second()) },
		"SSS":  func(t *time.Time) string { return formatZero3(t.Nanosecond() / 1000000) }, // TODO
		"a":    func(t *time.Time) string { return "" },
		"z":    func(t *time.Time) string { z, _ := t.Zone(); return z },
		"Z":    func(t *time.Time) string { return "" },
		"X":    func(t *time.Time) string { return "" },
		"XX":   func(t *time.Time) string { return "" },
		"XXX":  func(t *time.Time) string { return "" },
	}
)

func getFormatArg(ph []byte) formatArg {
	tmp := ph
	if len(tmp) > 4 {
		tmp = tmp[:4]
	}
	if arg, ok := formatArgs[string(tmp)]; ok {
		return arg
	}
	return func(*time.Time) string { return string(ph) } // Do not modify
}

func formatMax9(v int) string {
	return string(byte('0' + v))
}

func formatMax99(v int) string {
	if v > 10 {
		return formatMax9(v/10) + formatMax9(v%10)
	}
	return formatMax9(v)
}

func formatMax999(v int) string {
	if v > 100 {
		return formatMax9(v/100) + formatMax9((v/10)%10) + formatMax9(v%10)
	}
	return formatMax99(v)
}

func formatZero2(v int) string {
	if v > 10 {
		return formatMax99(v)
	}
	return "0" + formatMax9(v)
}

func formatZero3(v int) string {
	if v > 100 {
		return formatMax999(v)
	}
	if v > 10 {
		return "0" + formatMax99(v)
	}
	return "00" + formatMax9(v)
}
