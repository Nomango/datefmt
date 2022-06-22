package datefmt

import (
	"sync"
	"time"
	"unsafe"
)

// FastFormat ...
func FastFormat(t time.Time, generalLayout string) string {
	l := FastLayout(generalLayout)
	return l.Format(t)
}

type layout struct {
	max    int
	flag   formatFlag
	args   []*formatArg
	layout string
}

func (l *layout) Format(t time.Time) string {
	var (
		year       int
		month      time.Month
		day        int
		hour       int
		minute     int
		second     int
		zoneName   string
		zoneOffset int
		fb         fastBuilder
	)

	fb.Grow(l.max)
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
			arg.write(&fb, year)
		case formatFlagMonth:
			arg.write(&fb, int(month))
		case formatFlagDay:
			arg.write(&fb, day)
		case formatFlagHour:
			arg.write(&fb, hour)
		case formatFlagMinute:
			arg.write(&fb, minute)
		case formatFlagSecond:
			arg.write(&fb, second)
		case formatFlagWeekDay:
			arg.write(&fb, int(t.Weekday()))
		case formatFlagNanosecond:
			arg.write(&fb, t.Nanosecond())
		case formatFlagZoneName:
			fb.WriteString(zoneName)
		case formatFlagZoneOffset:
			arg.write(&fb, zoneOffset)
		case formatFlagYearDay:
			arg.write(&fb, t.YearDay())
		case formatFlagNone:
			arg.write(&fb, 0)
		}
	}
	// fmt.Println("len =", fb.Len(), ", cap =", fb.Cap(), ", max =", l.max)
	return fb.String()
}

var fastLayoutCache sync.Map

func FastLayout(generalLayout string) *layout {
	v, ok := fastLayoutCache.Load(generalLayout)
	if ok {
		return v.(*layout)
	}
	v, _ = fastLayoutCache.LoadOrStore(generalLayout, newFastLayout(generalLayout))
	return v.(*layout)
}

func newFastLayout(generalLayout string) *layout {
	var (
		l  = layout{layout: generalLayout}
		gl = []byte(generalLayout)
		n  = len(gl)
		sb = fastBuilder{}
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
			// fmt.Println("text =", sb.String())
			l.args = append(l.args, newBytesFormatArg(getBytesOfBuilder(&sb)))
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
		// fmt.Println("ph =", string(gl[s:e+1]))
		arg := getFormatArg(gl[s : e+1])
		l.args = append(l.args, arg)
		l.max += arg.max
		l.flag.Add(arg.flag)
	}
	// flush buffer
	if sb.Len() > 0 {
		// fmt.Println("text =", sb.String())
		l.args = append(l.args, newBytesFormatArg(getBytesOfBuilder(&sb)))
		l.max += sb.Len()
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

type formatArg struct {
	max  int
	flag formatFlag

	write func(fb *fastBuilder, v int)
}

func newBytesFormatArg(p []byte) *formatArg {
	return &formatArg{
		max:   len(p),
		flag:  formatFlagNone,
		write: func(fb *fastBuilder, _ int) { fb.Write(p) },
	}
}

var (
	formatArgs = map[string]formatArg{
		"YYYY": {max: 4, flag: formatFlagYear, write: func(fb *fastBuilder, v int) { formatMax9999(fb, v) }},
		"yyyy": {max: 4, flag: formatFlagYear, write: func(fb *fastBuilder, v int) { formatMax9999(fb, v) }},
		"YY":   {max: 2, flag: formatFlagYear, write: func(fb *fastBuilder, v int) { formatZero2(fb, v%100) }},
		"yy":   {max: 2, flag: formatFlagYear, write: func(fb *fastBuilder, v int) { formatZero2(fb, v%100) }},
		"MMMM": {max: 9, flag: formatFlagMonth, write: func(fb *fastBuilder, v int) { fb.WriteString(time.Month(v).String()) }},
		"MMM":  {max: 3, flag: formatFlagMonth, write: func(fb *fastBuilder, v int) { fb.WriteString(time.Month(v).String()[:3]) }},
		"MM":   {max: 2, flag: formatFlagMonth, write: func(fb *fastBuilder, v int) { formatZero2(fb, v) }},
		"M":    {max: 2, flag: formatFlagMonth, write: func(fb *fastBuilder, v int) { formatMax99(fb, v) }},
		"DDD":  {max: 3, flag: formatFlagNone, write: func(fb *fastBuilder, v int) { formatZero3(fb, v) }},
		"dd":   {max: 2, flag: formatFlagDay, write: func(fb *fastBuilder, v int) { formatZero2(fb, v) }},
		"d":    {max: 2, flag: formatFlagDay, write: func(fb *fastBuilder, v int) { formatMax99(fb, v) }},
		"EEEE": {max: 9, flag: formatFlagWeekDay, write: func(fb *fastBuilder, v int) { fb.WriteString(time.Weekday(v).String()) }},
		"EEE":  {max: 3, flag: formatFlagWeekDay, write: func(fb *fastBuilder, v int) { fb.WriteString(time.Weekday(v).String()[:3]) }},
		"HH":   {max: 2, flag: formatFlagHour, write: func(fb *fastBuilder, v int) { formatZero2(fb, v) }},
		"hh":   {max: 2, flag: formatFlagHour, write: func(fb *fastBuilder, v int) { formatZero2(fb, halfDayHour(v)) }},
		"H":    {max: 2, flag: formatFlagHour, write: func(fb *fastBuilder, v int) { formatMax99(fb, v) }},
		"h":    {max: 2, flag: formatFlagHour, write: func(fb *fastBuilder, v int) { formatMax99(fb, halfDayHour(v)) }},
		"mm":   {max: 2, flag: formatFlagMinute, write: func(fb *fastBuilder, v int) { formatZero2(fb, v) }},
		"m":    {max: 2, flag: formatFlagMinute, write: func(fb *fastBuilder, v int) { formatMax99(fb, v) }},
		"ss":   {max: 2, flag: formatFlagSecond, write: func(fb *fastBuilder, v int) { formatZero2(fb, v) }},
		"s":    {max: 2, flag: formatFlagSecond, write: func(fb *fastBuilder, v int) { formatMax99(fb, v) }},
		"SSS":  {max: 3, flag: formatFlagNanosecond, write: func(fb *fastBuilder, v int) { formatZero3(fb, v/1000000) }},
		"a":    {max: 2, flag: formatFlagHour, write: formatPM},
		"z":    {max: 3, flag: formatFlagZoneName},
		"Z":    {max: 0, flag: formatFlagZoneOffset, write: func(fb *fastBuilder, v int) {}},
		"X":    {max: 0, flag: formatFlagZoneOffset, write: func(fb *fastBuilder, v int) {}},
		"XX":   {max: 0, flag: formatFlagZoneOffset, write: func(fb *fastBuilder, v int) {}},
		"XXX":  {max: 0, flag: formatFlagZoneOffset, write: func(fb *fastBuilder, v int) {}},
	}
)

func formatPM(fb *fastBuilder, v int) {
	if v >= 12 {
		fb.WriteString("PM")
	} else {
		fb.WriteString("AM")
	}
}

func halfDayHour(hour int) int {
	hour %= 12
	if hour == 0 {
		return 12
	}
	return hour
}

func getFormatArg(ph []byte) *formatArg {
	tmp := ph
	if len(tmp) > 4 {
		tmp = tmp[:4]
	}
	if arg, ok := formatArgs[string(tmp)]; ok {
		return &arg
	}
	return newBytesFormatArg(ph) // Do not modify
}

func formatMax9(b *fastBuilder, v int) {
	b.WriteByte(byte('0' + v))
}

func formatMax99(b *fastBuilder, v int) {
	if v > 9 {
		formatMax9(b, v/10)
		formatMax9(b, v%10)
		return
	}
	formatMax9(b, v)
}

func formatMax999(b *fastBuilder, v int) {
	if v > 99 {
		formatMax9(b, v/100)
		formatMax9(b, (v/10)%10)
		formatMax9(b, v%10)
		return
	}
	formatMax99(b, v)
}

func formatMax9999(b *fastBuilder, v int) {
	if v > 999 {
		formatMax9(b, v/1000)
		formatMax9(b, (v/100)%10)
		formatMax9(b, (v/10)%10)
		formatMax9(b, v%10)
		return
	}
	formatMax999(b, v)
}

func formatZero2(b *fastBuilder, v int) {
	if v > 9 {
		formatMax99(b, v)
		return
	}
	formatMax9(b, 0)
	formatMax9(b, v)
}

func formatZero3(b *fastBuilder, v int) {
	if v > 99 {
		formatMax999(b, v)
		return
	}
	if v > 9 {
		formatMax9(b, 0)
		formatMax9(b, v/10)
		formatMax9(b, v%10)
		return
	}
	formatMax9(b, 0)
	formatMax9(b, 0)
	formatMax9(b, v)
}

func readOnlyBytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func string2ReadOnlyBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// type fastBuilder = strings.Builder

func getBytesOfBuilder(b *fastBuilder) []byte {
	return string2ReadOnlyBytes(b.String())
}

type fastBuilder struct {
	buf []byte
}

func (b *fastBuilder) Grow(n int) {
	if b.buf == nil {
		b.buf = make([]byte, 0, n)
		return
	}
	if cap(b.buf)-len(b.buf) < n {
		buf := make([]byte, len(b.buf), 2*cap(b.buf)+n)
		copy(buf, b.buf)
		b.buf = buf
	}
}

func (b *fastBuilder) Write(p []byte) {
	b.buf = append(b.buf, p...)
}

func (b *fastBuilder) WriteByte(c byte) error {
	b.buf = append(b.buf, c)
	return nil
}

func (b *fastBuilder) WriteString(s string) {
	b.buf = append(b.buf, s...)
}

func (b *fastBuilder) Reset() {
	b.buf = nil
}

func (b *fastBuilder) Cap() int {
	return cap(b.buf)
}

func (b *fastBuilder) Len() int {
	return len(b.buf)
}

func (b *fastBuilder) Bytes() []byte {
	return b.buf
}

func (b *fastBuilder) String() string {
	return readOnlyBytes2String(b.buf)
}
