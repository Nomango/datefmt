package datefmt

import (
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// FastFormat ...
func FastFormat(t time.Time, generalLayout string) string {
	l := FastLayout(generalLayout)
	return l.Format(t)
}

type formatArg = func(t *time.Time) string

type layout struct {
	args []interface{}
	max  int
}

func (l *layout) Format(t time.Time) string {
	var sb strings.Builder
	sb.Grow(l.max)
	for _, arg := range l.args {
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
			l.max += sb.Len()
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
		f, max := getFormatArg(gl[s : e+1])
		l.args = append(l.args, f)
		l.max += max
	}
	if sb.Len() > 0 {
		l.args = append(l.args, sb.String())
		l.max += sb.Len()
	}
	return &l
}

type formatArgWithMax struct {
	max int
	f   formatArg
}

var (
	formatArgs = map[string]formatArgWithMax{
		"YYYY": {max: 4, f: func(t *time.Time) string { return strconv.Itoa(t.Year()) }},
		"yyyy": {max: 4, f: func(t *time.Time) string { return strconv.Itoa(t.Year()) }},
		"YY":   {max: 2, f: func(t *time.Time) string { return strconv.Itoa(t.Year())[2:4] }},
		"yy":   {max: 2, f: func(t *time.Time) string { return strconv.Itoa(t.Year())[2:4] }},
		"MMMM": {max: 9, f: func(t *time.Time) string { return t.Month().String() }},
		"MMM":  {max: 3, f: func(t *time.Time) string { return t.Month().String()[:3] }},
		"MM":   {max: 2, f: func(t *time.Time) string { return readOnlyBytes2String(formatZero2(int(t.Month()))) }},
		"M":    {max: 2, f: func(t *time.Time) string { return readOnlyBytes2String(formatMax99(int(t.Month()))) }},
		"DDD":  {max: 3, f: func(t *time.Time) string { return readOnlyBytes2String(formatZero3(t.YearDay())) }},
		"dd":   {max: 2, f: func(t *time.Time) string { return readOnlyBytes2String(formatZero2(t.Day())) }},
		"d":    {max: 2, f: func(t *time.Time) string { return readOnlyBytes2String(formatMax99(t.Day())) }},
		"EEEE": {max: 9, f: func(t *time.Time) string { return t.Weekday().String() }},
		"EEE":  {max: 3, f: func(t *time.Time) string { return t.Weekday().String()[:3] }},
		"HH":   {max: 2, f: func(t *time.Time) string { return readOnlyBytes2String(formatZero2(t.Hour())) }},
		"hh":   {max: 2, f: func(t *time.Time) string { return "TODO" }},
		"H":    {max: 2, f: func(t *time.Time) string { return readOnlyBytes2String(formatMax99(t.Hour())) }},
		"h":    {max: 2, f: func(t *time.Time) string { return "TODO" }},
		"mm":   {max: 2, f: func(t *time.Time) string { return readOnlyBytes2String(formatZero2(t.Minute())) }},
		"m":    {max: 2, f: func(t *time.Time) string { return readOnlyBytes2String(formatMax99(t.Minute())) }},
		"ss":   {max: 2, f: func(t *time.Time) string { return readOnlyBytes2String(formatZero2(t.Second())) }},
		"s":    {max: 2, f: func(t *time.Time) string { return readOnlyBytes2String(formatMax99(t.Second())) }},
		"SSS":  {max: 0, f: func(t *time.Time) string { return readOnlyBytes2String(formatZero3(t.Nanosecond() / 1000000)) /*TODO*/ }},
		"a":    {max: 2, f: func(t *time.Time) string { return "" }},
		"z":    {max: 0, f: func(t *time.Time) string { z, _ := t.Zone(); return z }},
		"Z":    {max: 0, f: func(t *time.Time) string { return "" }},
		"X":    {max: 0, f: func(t *time.Time) string { return "" }},
		"XX":   {max: 0, f: func(t *time.Time) string { return "" }},
		"XXX":  {max: 0, f: func(t *time.Time) string { return "" }},
	}
)

func getFormatArg(ph []byte) (formatArg, int) {
	tmp := ph
	if len(tmp) > 4 {
		tmp = tmp[:4]
	}
	if arg, ok := formatArgs[string(tmp)]; ok {
		return arg.f, arg.max
	}
	return func(*time.Time) string { return string(ph) }, len(ph) // Do not modify
}

func formatMax9(v int) byte {
	return byte('0' + v)
}

func formatMax99(v int) []byte {
	if v > 9 {
		return []byte{formatMax9(v / 10), formatMax9(v % 10)}
	}
	return []byte{formatMax9(v)}
}

func formatMax999(v int) []byte {
	if v > 99 {
		return []byte{formatMax9(v / 100), formatMax9((v / 10) % 10), formatMax9(v % 10)}
	}
	return formatMax99(v)
}

func formatZero2(v int) []byte {
	if v > 9 {
		return formatMax99(v)
	}
	return []byte{'0', formatMax9(v)}
}

func formatZero3(v int) []byte {
	if v > 99 {
		return formatMax999(v)
	}
	if v > 9 {
		return []byte{'0', formatMax9(v / 10), formatMax9(v % 10)}
	}
	return []byte{'0', '0', formatMax9(v)}
}

func readOnlyBytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
