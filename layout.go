package datefmt

import (
	"strings"
	"time"
	"unsafe"
)

type Layout struct {
	max    int
	flag   formatFlag
	args   []*formatArg
	layout string
}

func (l *Layout) String() string {
	return l.layout
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
		switch arg.ph.flag {
		case formatFlagNone:
			p = arg.ph.format(p, 0, arg.w)
		case formatFlagYear:
			p = arg.ph.format(p, year, arg.w)
		case formatFlagMonth:
			p = arg.ph.format(p, int(month), arg.w)
		case formatFlagDay:
			p = arg.ph.format(p, day, arg.w)
		case formatFlagHour:
			p = arg.ph.format(p, hour, arg.w)
		case formatFlagMinute:
			p = arg.ph.format(p, minute, arg.w)
		case formatFlagSecond:
			p = arg.ph.format(p, second, arg.w)
		case formatFlagWeekDay:
			p = arg.ph.format(p, int(t.Weekday()), arg.w)
		case formatFlagNanosecond:
			p = arg.ph.format(p, t.Nanosecond(), arg.w)
		case formatFlagZoneName:
			p = formatZoneName(p, zoneName, zoneOffset)
		case formatFlagZoneOffset:
			p = arg.ph.format(p, zoneOffset, arg.w)
		case formatFlagYearDay:
			p = arg.ph.format(p, t.YearDay(), arg.w)
		case formatFlagWeekYear:
			iosYear, _ := t.ISOWeek()
			p = arg.ph.format(p, iosYear, arg.w)
		case formatFlagWeekInYear:
			_, isoWeek := t.ISOWeek()
			p = arg.ph.format(p, isoWeek, arg.w)
		case formatFlagWeekInMonth:
			firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
			firstWeekYear, firstWeek := firstDayOfMonth.ISOWeek()
			isoYear, isoWeek := t.ISOWeek()
			if firstWeekYear == isoYear {
				p = arg.ph.format(p, isoWeek-firstWeek+1, arg.w)
			} else if month == time.January {
				// Jan 01 to Jan 03 of year n might belong to week 52 or 53 of year n-1
				p = arg.ph.format(p, isoWeek+1, arg.w)
			} else {
				// Dec 29 to Dec 31 might belong to week 1 of year n+1
				_, lastWeek := time.Date(year, time.December, 28, 0, 0, 0, 0, t.Location()).ISOWeek()
				p = arg.ph.format(p, lastWeek-firstWeek+2, arg.w)
			}
		}
	}
	// fmt.Println("len =", fb.Len(), ", cap =", fb.Cap(), ", max =", l.max)
	return readOnlyBytes2String(p)
}

func NewLayout(generalLayout string) *Layout {
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
		if _, ok := placeholders[gl[i]]; !ok && gl[i] != '\'' {
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
		arg := newPlaceholderFormatArg(gl[s : e+1])
		l.args = append(l.args, arg)
		l.max += arg.max
		l.flag.Add(arg.ph.flag)
	}
	// flush buffer
	if sb.Len() > 0 {
		// fmt.Println("text =", sb.String())
		flushBuffer()
	}
	return &l
}

type formatArg struct {
	s   string
	w   int
	max int
	ph  placeholder
}

func newBytesFormatArg(p []byte) *formatArg {
	return &formatArg{
		s:   readOnlyBytes2String(p),
		w:   len(p),
		max: len(p),
		ph: placeholder{
			flag:   formatFlagNone,
			format: func(dst []byte, _, _ int) []byte { return append(dst, p...) },
		},
	}
}

func newPlaceholderFormatArg(p []byte) *formatArg {
	ph := placeholders[p[0]]
	return &formatArg{
		s:   readOnlyBytes2String(p),
		w:   len(p),
		max: ph.max(len(p)),
		ph:  *ph,
	}
}

type formatFlag int

const (
	formatFlagNone formatFlag = iota
	formatFlagNanosecond
	formatFlagWeekYear
	formatFlagWeekInYear
	formatFlagYearDay
	formatFlagWeekDay

	formatFlagYear formatFlag = iota + formatFlagNeedDate
	formatFlagMonth
	formatFlagDay
	formatFlagWeekInMonth

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

type placeholder struct {
	max    func(int) int
	flag   formatFlag
	format func(p []byte, v, w int) []byte
}

var (
	placeholders = map[byte]*placeholder{
		'G': {max: fixedMax(2), flag: formatFlagYear, format: formatEra},
		'y': {max: yearMax, flag: formatFlagYear, format: formatYear},
		'Y': {max: yearMax, flag: formatFlagWeekYear, format: formatYear},
		'M': {max: monthMax, flag: formatFlagMonth, format: formatMonth},
		'w': {max: numberMax(2), flag: formatFlagWeekInYear, format: formatNumProbably2Digits},
		'W': {max: numberMax(2), flag: formatFlagWeekInMonth, format: formatNumProbably2Digits},
		'D': {max: numberMax(3), flag: formatFlagYearDay, format: formatNumProbably3Digits},
		'd': {max: numberMax(2), flag: formatFlagDay, format: formatNumProbably2Digits},
		'F': {max: numberMax(1), flag: formatFlagDay, format: func(p []byte, v, w int) []byte { return formatNumProbably2Digits(p, dayOfWeekInMonth(v), w) }},
		'E': {max: textMax(3, 9), flag: formatFlagWeekDay, format: formatWeek},
		'u': {max: numberMax(1), flag: formatFlagWeekDay, format: func(p []byte, v, w int) []byte { return formatNumProbably2Digits(p, dayNumOfWeek(v), w) }},
		'a': {max: fixedMax(2), flag: formatFlagHour, format: formatPM},
		'H': {max: numberMax(2), flag: formatFlagHour, format: formatNumProbably2Digits},
		'k': {max: numberMax(2), flag: formatFlagHour, format: func(p []byte, v, w int) []byte { return formatNumProbably2Digits(p, hour24(v), w) }},
		'K': {max: numberMax(2), flag: formatFlagHour, format: func(p []byte, v, w int) []byte { return formatNumProbably2Digits(p, v%12, w) }},
		'h': {max: numberMax(2), flag: formatFlagHour, format: func(p []byte, v, w int) []byte { return formatNumProbably2Digits(p, hour12(v), w) }},
		'm': {max: numberMax(2), flag: formatFlagMinute, format: formatNumProbably2Digits},
		's': {max: numberMax(2), flag: formatFlagSecond, format: formatNumProbably2Digits},
		'S': {max: nanosecondMax, flag: formatFlagNanosecond, format: formatNanosecond},
		'z': {max: fixedMax(5), flag: formatFlagZoneName},
		'Z': {max: fixedMax(5), flag: formatFlagZoneOffset, format: formatZoneOffsetRFC822},
		'X': {max: fixedMax(6), flag: formatFlagZoneOffset, format: formatZoneOffsetISO8601},
	}
)

// y Y Year

func yearMax(w int) int {
	if w == 2 {
		return 2
	}
	return 4
}

func formatYear(p []byte, year, w int) []byte {
	if w == 2 {
		return formatNumProbably2Digits(p, year%100, 2)
	}
	return formatNumProbably4Digits(p, year, w)
}

// M Month

func monthMax(w int) int {
	if w < 3 {
		return 2
	}
	if w == 3 {
		return 3
	}
	return 9
}

func formatMonth(p []byte, month, w int) []byte {
	if w < 3 {
		return formatNumProbably2Digits(p, month, w)
	}
	if w == 3 {
		return formatString(p, time.Month(month).String()[:3])
	}
	return formatString(p, time.Month(month).String())
}

// E Week

func formatWeek(p []byte, week, w int) []byte {
	if w <= 3 {
		return formatString(p, time.Weekday(week).String()[:3])
	}
	return formatString(p, time.Weekday(week).String())
}

// G Era

func formatEra(p []byte, year, w int) []byte {
	if year < 0 {
		return append(p, "BC"...)
	}
	return append(p, "AD"...)
}

// a PM

func formatPM(p []byte, v, w int) []byte {
	if v >= 12 {
		return append(p, "PM"...)
	}
	return append(p, "AM"...)
}

// S Nanosecond

func nanosecondMax(w int) int {
	return w
}

func formatNanosecond(p []byte, v, w int) []byte {
	// truncate nanosecond number
	div := 1000000000
	i := w
	for i > 0 {
		div /= 10
		i--
	}
	return formatNum(p, v/div, w)
}

// z Zone name

func formatZoneName(p []byte, zoneName string, zoneOffset int) []byte {
	if len(zoneName) > 0 {
		return append(p, zoneName...)
	}
	// No time zone known for this time, but we must print one.
	// Use the RFC822 format.
	return formatZoneOffsetRFC822(p, zoneOffset, 0)
}

// Z Zone RFC822

func formatZoneOffsetRFC822(p []byte, zoneOffset, w int) []byte {
	sign, hour, minute := getZoneOffsetParts(zoneOffset)
	return formatNum(formatNum(append(p, sign), hour, 2), minute, 2)
}

// X Zone ISO8601

func formatZoneOffsetISO8601(p []byte, zoneOffset, w int) []byte {
	if zoneOffset == 0 {
		return append(p, 'Z')
	}
	sign, hour, minute := getZoneOffsetParts(zoneOffset)
	switch w {
	case 1:
		return formatNum(append(p, sign), hour, 2)
	case 2:
		return formatNum(formatNum(append(p, sign), hour, 2), minute, 2)
	default:
		return formatNum(append(formatNum(append(p, sign), hour, 2), ':'), minute, 2)
	}
}

// helper functions

func hour12(hour int) int {
	hour %= 12
	if hour == 0 {
		return 12
	}
	return hour
}

func hour24(hour int) int {
	if hour == 0 {
		return 24
	}
	return hour
}

func dayNumOfWeek(week int) int {
	if time.Weekday(week) == time.Sunday {
		return 7
	}
	return week
}

func dayOfWeekInMonth(day int) int {
	return (day-1)/7 + 1
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

func fixedMax(max int) func(int) int {
	return func(int) int {
		return max
	}
}

func textMax(short, long int) func(int) int {
	return func(w int) int {
		if w >= 4 {
			return long
		}
		return short
	}
}

func numberMax(minWidth int) func(int) int {
	return func(w int) int {
		if w > minWidth {
			return w
		}
		return minWidth
	}
}

func readOnlyBytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func string2ReadOnlyBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
