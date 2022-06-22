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
	args   []func(f *formatter) []byte
	layout string
}

func (l *layout) Format(t time.Time) string {
	var (
		f  = formatter{t: t}
		sb = fastBuilder{}
	)
	f.apply(l.flag)
	sb.Grow(l.max)
	// cap1 := sb.Cap()
	for i := range l.args {
		sb.Write(l.args[i](&f))
	}
	// cap2 := sb.Cap()
	// if sb.Cap() > l.max {
	// 	fmt.Println(cap1, cap2)
	// 	panic("cap greater than max!!")
	// }
	// fmt.Println("len =", sb.Len(), ", cap =", sb.Cap(), ", max =", l.max)
	return sb.String()
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
			s := sb.String()
			l.args = append(l.args, func(*formatter) []byte { return string2ReadOnlyBytes(s) })
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
		l.args = append(l.args, arg.f)
		l.max += arg.max
		l.flag.Add(arg.flag)
	}
	// flush buffer
	if sb.Len() > 0 {
		// fmt.Println("text =", sb.String())
		s := sb.String()
		l.args = append(l.args, func(*formatter) []byte { return string2ReadOnlyBytes(s) })
		l.max += sb.Len()
	}
	return &l
}

type formatFlag int

const (
	formatFlagNeedDate  formatFlag = 1
	formatFlagNeedClock formatFlag = 1 << 1
)

func (f *formatFlag) Add(flag formatFlag) {
	*f |= flag
}

func (f formatFlag) Has(flag formatFlag) bool {
	return f&flag != 0
}

type formatter struct {
	t      time.Time
	year   int
	month  time.Month
	day    int
	hour   int
	minute int
	second int
}

func (f *formatter) apply(flag formatFlag) {
	if flag.Has(formatFlagNeedDate) {
		f.year, f.month, f.day = f.t.Date()
	}
	if flag.Has(formatFlagNeedClock) {
		f.hour, f.minute, f.second = f.t.Clock()
	}
}

type formatArg struct {
	max  int
	flag formatFlag
	f    func(*formatter) []byte
}

var (
	formatArgs = map[string]formatArg{
		"YYYY": {max: 4, flag: formatFlagNeedDate, f: func(f *formatter) []byte { return formatMax9999(f.year) }},
		"yyyy": {max: 4, flag: formatFlagNeedDate, f: func(f *formatter) []byte { return formatMax9999(f.year) }},
		"YY":   {max: 2, flag: formatFlagNeedDate, f: func(f *formatter) []byte { return formatZero2(f.year % 100) }},
		"yy":   {max: 2, flag: formatFlagNeedDate, f: func(f *formatter) []byte { return formatZero2(f.year % 100) }},
		"MMMM": {max: 9, flag: formatFlagNeedDate, f: func(f *formatter) []byte { return string2ReadOnlyBytes(f.month.String()) }},
		"MMM":  {max: 3, flag: formatFlagNeedDate, f: func(f *formatter) []byte { return string2ReadOnlyBytes(f.month.String())[:3] }},
		"MM":   {max: 2, flag: formatFlagNeedDate, f: func(f *formatter) []byte { return formatZero2(int(f.month)) }},
		"M":    {max: 2, flag: formatFlagNeedDate, f: func(f *formatter) []byte { return formatMax99(int(f.month)) }},
		"DDD":  {max: 3, flag: formatFlagNeedDate, f: func(f *formatter) []byte { return (formatZero3(f.t.YearDay())) }},
		"dd":   {max: 2, flag: formatFlagNeedDate, f: func(f *formatter) []byte { return formatZero2(f.day) }},
		"d":    {max: 2, flag: formatFlagNeedDate, f: func(f *formatter) []byte { return formatMax99(f.day) }},
		"EEEE": {max: 9, f: func(f *formatter) []byte { return string2ReadOnlyBytes(f.t.Weekday().String()) }},
		"EEE":  {max: 3, f: func(f *formatter) []byte { return string2ReadOnlyBytes(f.t.Weekday().String())[:3] }},
		"HH":   {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) []byte { return formatZero2(f.hour) }},
		"hh":   {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) []byte { return nil }},
		"H":    {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) []byte { return formatMax99(f.hour) }},
		"h":    {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) []byte { return nil }},
		"mm":   {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) []byte { return formatZero2(f.minute) }},
		"m":    {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) []byte { return formatMax99(f.minute) }},
		"ss":   {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) []byte { return formatZero2(f.second) }},
		"s":    {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) []byte { return formatMax99(f.second) }},
		"SSS":  {max: 0, f: func(f *formatter) []byte { return formatZero3(f.t.Nanosecond() / 1000000) /*TODO*/ }},
		"a":    {max: 2, f: func(f *formatter) []byte { return nil }},
		"z":    {max: 3, f: func(f *formatter) []byte { z, _ := f.t.Zone(); return string2ReadOnlyBytes(z) }},
		"Z":    {max: 0, f: func(f *formatter) []byte { return nil }},
		"X":    {max: 0, f: func(f *formatter) []byte { return nil }},
		"XX":   {max: 0, f: func(f *formatter) []byte { return nil }},
		"XXX":  {max: 0, f: func(f *formatter) []byte { return nil }},
	}
)

var (
	nums = []byte{
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9',
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
	return formatArg{max: len(ph), f: func(*formatter) []byte { return ph }} // Do not modify
}

func formatMax9(v int) byte {
	return nums[v]
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

func formatMax9999(v int) []byte {
	if v > 999 {
		return []byte{formatMax9(v / 1000), formatMax9((v / 100) % 10), formatMax9((v / 10) % 10), formatMax9(v % 10)}
	}
	return formatMax999(v)
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

func string2ReadOnlyBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
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

func (b *fastBuilder) WriteByte(c byte) {
	b.buf = append(b.buf, c)
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
