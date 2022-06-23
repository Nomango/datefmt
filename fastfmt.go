package datefmt

import (
	"strings"
	"sync"
	"time"
	"unsafe"
)

// FastFormat is a faster and general layout based version of time.Format
func FastFormat(t time.Time, generalLayout string) string {
	l := FastLayout(generalLayout)
	return l.Format(t)
}

// Format is a general layout based version of time.Format
func StdFormat(t time.Time, generalLayout string) string {
	return t.Format(GoLayout(generalLayout))
}

type Layout struct {
	max    int
	flag   formatFlag
	args   []*formatArg
	layout string
}

func (l *Layout) Format(t time.Time) string {
	var (
		year       int
		month      time.Month
		day        int
		hour       int
		minute     int
		second     int
		zoneName   string
		zoneOffset int
		p          []byte
	)

	p = make([]byte, 0, l.max)
	if l.flag.Has(formatFlagNeedDate) {
		year, month, day = t.Date()
	}
	if l.flag.Has(formatFlagNeedClock) {
		hour, minute, second = t.Clock()
	}
	if l.flag.Has(formatFlagNeedZone) {
		zoneName, zoneOffset = t.Zone()
	}
	for _, arg := range l.args {
		switch arg.flag {
		case formatFlagYear:
			p = arg.appendInt(p, year)
		case formatFlagMonth:
			p = arg.appendInt(p, int(month))
		case formatFlagDay:
			p = arg.appendInt(p, day)
		case formatFlagHour:
			p = arg.appendInt(p, hour)
		case formatFlagMinute:
			p = arg.appendInt(p, minute)
		case formatFlagSecond:
			p = arg.appendInt(p, second)
		case formatFlagWeekDay:
			p = arg.appendInt(p, int(t.Weekday()))
		case formatFlagNanosecond:
			p = arg.appendInt(p, t.Nanosecond())
		case formatFlagZoneName:
			p = formatZoneName(p, zoneName, zoneOffset)
		case formatFlagZoneOffset:
			p = arg.appendInt(p, zoneOffset)
		case formatFlagYearDay:
			p = arg.appendInt(p, t.YearDay())
		case formatFlagNone:
			p = arg.appendInt(p, 0)
		}
	}
	// fmt.Println("len =", fb.Len(), ", cap =", fb.Cap(), ", max =", l.max)
	return readOnlyBytes2String(p)
}

var fastLayoutCache sync.Map

func FastLayout(generalLayout string) *Layout {
	v, ok := fastLayoutCache.Load(generalLayout)
	if ok {
		return v.(*Layout)
	}
	v, _ = fastLayoutCache.LoadOrStore(generalLayout, newFastLayout(generalLayout))
	return v.(*Layout)
}

func newFastLayout(generalLayout string) *Layout {
	var (
		l    = Layout{layout: generalLayout}
		gl   = []byte(generalLayout)
		n    = len(gl)
		sb   = strings.Builder{}
		tmax = 8 // text max length
	)

	flushBuffer := func() {
		l.args = append(l.args, newBytesFormatArg(string2ReadOnlyBytes(sb.String())))
		l.max += sb.Len()
		sb.Reset()
		sb.Grow(tmax)
	}
	sb.Grow(tmax)
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
			// fmt.Println("text =", sb.String())
			flushBuffer()
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
		// fmt.Println("ph =", string(gl[s:e+1]))
		arg := getPlaceholderFormatArg(gl[s : e+1])
		l.args = append(l.args, arg)
		l.max += arg.max
		l.flag.Add(arg.flag)
	}
	// flush buffer
	if sb.Len() > 0 {
		// fmt.Println("text =", sb.String())
		flushBuffer()
	}
	return &l
}

type formatFlag int

const (
	formatFlagNone formatFlag = iota
	formatFlagYearDay
	formatFlagWeekDay
	formatFlagNanosecond

	formatFlagYear formatFlag = iota + formatFlagNeedDate
	formatFlagMonth
	formatFlagDay

	formatFlagHour formatFlag = iota + formatFlagNeedClock
	formatFlagMinute
	formatFlagSecond

	formatFlagZoneName formatFlag = iota + formatFlagNeedZone
	formatFlagZoneOffset

	formatFlagNeedDate  formatFlag = 1 << 0 << 7
	formatFlagNeedClock formatFlag = 1 << 1 << 7
	formatFlagNeedZone  formatFlag = 1 << 2 << 7
)

func (f *formatFlag) Add(flag formatFlag) {
	*f |= flag
}

func (f formatFlag) Has(flag formatFlag) bool {
	return f&flag != 0
}

var (
	formatArgs = map[string]formatArg{
		"YYYY": {max: 4, flag: formatFlagYear, appendInt: func(p []byte, v int) []byte { return formatYear(p, v) }},
		"yyyy": {max: 4, flag: formatFlagYear, appendInt: func(p []byte, v int) []byte { return formatYear(p, v) }},
		"YY":   {max: 2, flag: formatFlagYear, appendInt: func(p []byte, v int) []byte { return formatZero2(p, v%100) }},
		"yy":   {max: 2, flag: formatFlagYear, appendInt: func(p []byte, v int) []byte { return formatZero2(p, v%100) }},
		"MMMM": {max: 9, flag: formatFlagMonth, appendInt: func(p []byte, v int) []byte { return formatString(p, time.Month(v).String()) }},
		"MMM":  {max: 3, flag: formatFlagMonth, appendInt: func(p []byte, v int) []byte { return formatString(p, time.Month(v).String()[:3]) }},
		"MM":   {max: 2, flag: formatFlagMonth, appendInt: func(p []byte, v int) []byte { return formatZero2(p, v) }},
		"M":    {max: 2, flag: formatFlagMonth, appendInt: func(p []byte, v int) []byte { return formatMax99(p, v) }},
		"DDD":  {max: 3, flag: formatFlagYearDay, appendInt: func(p []byte, v int) []byte { return formatZero3(p, v) }},
		"D":    {max: 3, flag: formatFlagYearDay, appendInt: func(p []byte, v int) []byte { return formatMax999(p, v) }},
		"dd":   {max: 2, flag: formatFlagDay, appendInt: func(p []byte, v int) []byte { return formatZero2(p, v) }},
		"d":    {max: 2, flag: formatFlagDay, appendInt: func(p []byte, v int) []byte { return formatMax99(p, v) }},
		"EEEE": {max: 9, flag: formatFlagWeekDay, appendInt: func(p []byte, v int) []byte { return formatString(p, time.Weekday(v).String()) }},
		"EEE":  {max: 3, flag: formatFlagWeekDay, appendInt: func(p []byte, v int) []byte { return formatString(p, time.Weekday(v).String()[:3]) }},
		"HH":   {max: 2, flag: formatFlagHour, appendInt: func(p []byte, v int) []byte { return formatZero2(p, v) }},
		"hh":   {max: 2, flag: formatFlagHour, appendInt: func(p []byte, v int) []byte { return formatZero2(p, halfDayHour(v)) }},
		"H":    {max: 2, flag: formatFlagHour, appendInt: func(p []byte, v int) []byte { return formatMax99(p, v) }},
		"h":    {max: 2, flag: formatFlagHour, appendInt: func(p []byte, v int) []byte { return formatMax99(p, halfDayHour(v)) }},
		"mm":   {max: 2, flag: formatFlagMinute, appendInt: func(p []byte, v int) []byte { return formatZero2(p, v) }},
		"m":    {max: 2, flag: formatFlagMinute, appendInt: func(p []byte, v int) []byte { return formatMax99(p, v) }},
		"ss":   {max: 2, flag: formatFlagSecond, appendInt: func(p []byte, v int) []byte { return formatZero2(p, v) }},
		"s":    {max: 2, flag: formatFlagSecond, appendInt: func(p []byte, v int) []byte { return formatMax99(p, v) }},
		"SSS":  {max: 3, flag: formatFlagNanosecond, appendInt: func(p []byte, v int) []byte { return formatZero3(p, v/1000000) }},
		"a":    {max: 2, flag: formatFlagHour, appendInt: formatPM},
		"z":    {max: 5, flag: formatFlagZoneName},
		"Z":    {max: 5, flag: formatFlagZoneOffset, appendInt: formatZoneOffsetRFC822},
		"X":    {max: 3, flag: formatFlagZoneOffset, appendInt: formatZoneOffsetISO8601(1)},
		"XX":   {max: 5, flag: formatFlagZoneOffset, appendInt: formatZoneOffsetISO8601(2)},
		"XXX":  {max: 6, flag: formatFlagZoneOffset, appendInt: formatZoneOffsetISO8601(3)},
	}
)

type formatArg struct {
	max  int
	flag formatFlag

	appendInt func([]byte, int) []byte
}

func newBytesFormatArg(p []byte) *formatArg {
	return &formatArg{
		max:       len(p),
		flag:      formatFlagNone,
		appendInt: func(dst []byte, _ int) []byte { return append(dst, p...) },
	}
}

func getPlaceholderFormatArg(ph []byte) *formatArg {
	tmp := ph
	if len(tmp) > 4 {
		tmp = tmp[:4]
	}
	if arg, ok := formatArgs[string(tmp)]; ok {
		return &arg
	}
	return newBytesFormatArg(ph) // Do not modify
}

func halfDayHour(hour int) int {
	hour %= 12
	if hour == 0 {
		return 12
	}
	return hour
}

func getZoneOffsetParts(zoneOffset int) (sign byte, hour, minute int) {
	zone := zoneOffset / 60 // convert to minutes
	if zone < 0 {
		sign = '-'
		zone = -zone
	} else {
		sign = '+'
	}
	hour = zone / 60
	minute = zone % 60
	return
}

func formatYear(p []byte, v int) []byte {
	return formatMax9(formatMax9(formatMax9(formatMax9(p, v/1000), (v/100)%10), (v/10)%10), v%10)
}

func formatPM(p []byte, v int) []byte {
	if v >= 12 {
		return formatString(p, "PM")
	}
	return formatString(p, "AM")
}

func formatZoneName(p []byte, zoneName string, zoneOffset int) []byte {
	if len(zoneName) > 0 {
		return append(p, zoneName...)
	}
	// No time zone known for this time, but we must print one.
	// Use the RFC822 format.
	return formatZoneOffsetRFC822(p, zoneOffset)
}

func formatZoneOffsetRFC822(p []byte, zoneOffset int) []byte {
	sign, hour, minute := getZoneOffsetParts(zoneOffset)
	return formatZero2(formatZero2(append(p, sign), hour), minute)
}

func formatZoneOffsetISO8601(n int) func([]byte, int) []byte {
	return func(p []byte, zoneOffset int) []byte {
		if zoneOffset == 0 {
			return append(p, 'Z')
		}
		sign, hour, minute := getZoneOffsetParts(zoneOffset)
		switch n {
		case 1:
			return formatZero2(append(p, sign), hour)
		case 2:
			return formatZero2(formatZero2(append(p, sign), hour), minute)
		default:
			return formatZero2(append(formatZero2(append(p, sign), hour), ':'), minute)
		}
	}
}

func formatString(p []byte, s string) []byte {
	return append(p, s...)
}

func formatMax9(p []byte, v int) []byte {
	return append(p, byte('0'+v))
}

func formatMax99(p []byte, v int) []byte {
	if v > 9 {
		return formatMax9(formatMax9(p, v/10), v%10)
	}
	return formatMax9(p, v)
}

func formatMax999(p []byte, v int) []byte {
	if v > 99 {
		return formatMax9(formatMax9(formatMax9(p, v/100), (v/10)%10), v%10)
	}
	return formatMax99(p, v)
}

func formatZero2(p []byte, v int) []byte {
	if v > 9 {
		return formatMax99(p, v)
	}
	return formatMax9(formatMax9(p, 0), v)
}

func formatZero3(p []byte, v int) []byte {
	if v > 99 {
		return formatMax999(p, v)
	}
	return formatZero2(formatMax9(p, 0), v)
}

func readOnlyBytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func string2ReadOnlyBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
