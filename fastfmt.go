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
	args   []func(f *formatter)
	layout string
}

func (l *layout) Format(t time.Time) string {
	f := &formatter{}
	f.init(l, t)
	for i := range l.args {
		l.args[i](f)
	}
	// fmt.Println("len =", sb.Len(), ", cap =", sb.Cap(), ", max =", l.max)
	return f.fb.String()
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
			l.args = append(l.args, func(f *formatter) { f.fb.WriteString(s) })
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
		l.args = append(l.args, func(f *formatter) { f.fb.WriteString(s) })
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
	fb     fastBuilder
}

func (f *formatter) init(l *layout, t time.Time) {
	if l.flag.Has(formatFlagNeedDate) {
		f.year, f.month, f.day = f.t.Date()
	}
	if l.flag.Has(formatFlagNeedClock) {
		f.hour, f.minute, f.second = f.t.Clock()
	}
	f.fb.Grow(l.max)
}

type formatArg struct {
	max  int
	flag formatFlag
	f    func(*formatter)
}

var (
	formatArgs = map[string]formatArg{
		"YYYY": {max: 4, flag: formatFlagNeedDate, f: func(f *formatter) { formatMax9999(&f.fb, f.year) }},
		"yyyy": {max: 4, flag: formatFlagNeedDate, f: func(f *formatter) { formatMax9999(&f.fb, f.year) }},
		"YY":   {max: 2, flag: formatFlagNeedDate, f: func(f *formatter) { formatZero2(&f.fb, f.year%100) }},
		"yy":   {max: 2, flag: formatFlagNeedDate, f: func(f *formatter) { formatZero2(&f.fb, f.year%100) }},
		"MMMM": {max: 9, flag: formatFlagNeedDate, f: func(f *formatter) { f.fb.WriteString(f.month.String()) }},
		"MMM":  {max: 3, flag: formatFlagNeedDate, f: func(f *formatter) { f.fb.WriteString(f.month.String()[:3]) }},
		"MM":   {max: 2, flag: formatFlagNeedDate, f: func(f *formatter) { formatZero2(&f.fb, int(f.month)) }},
		"M":    {max: 2, flag: formatFlagNeedDate, f: func(f *formatter) { formatMax99(&f.fb, int(f.month)) }},
		"DDD":  {max: 3, flag: formatFlagNeedDate, f: func(f *formatter) { formatZero3(&f.fb, f.t.YearDay()) }},
		"dd":   {max: 2, flag: formatFlagNeedDate, f: func(f *formatter) { formatZero2(&f.fb, f.day) }},
		"d":    {max: 2, flag: formatFlagNeedDate, f: func(f *formatter) { formatMax99(&f.fb, f.day) }},
		"EEEE": {max: 9, f: func(f *formatter) { f.fb.WriteString(f.t.Weekday().String()) }},
		"EEE":  {max: 3, f: func(f *formatter) { f.fb.WriteString(f.t.Weekday().String()[:3]) }},
		"HH":   {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) { formatZero2(&f.fb, f.hour) }},
		"hh":   {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) {}},
		"H":    {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) { formatMax99(&f.fb, f.hour) }},
		"h":    {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) {}},
		"mm":   {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) { formatZero2(&f.fb, f.minute) }},
		"m":    {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) { formatMax99(&f.fb, f.minute) }},
		"ss":   {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) { formatZero2(&f.fb, f.second) }},
		"s":    {max: 2, flag: formatFlagNeedClock, f: func(f *formatter) { formatMax99(&f.fb, f.second) }},
		"SSS":  {max: 0, f: func(f *formatter) { formatZero3(&f.fb, f.t.Nanosecond()/1000000) /*TODO*/ }},
		"a":    {max: 2, f: func(f *formatter) {}},
		"z":    {max: 3, f: func(f *formatter) { z, _ := f.t.Zone(); f.fb.WriteString(z) }},
		"Z":    {max: 0, f: func(f *formatter) {}},
		"X":    {max: 0, f: func(f *formatter) {}},
		"XX":   {max: 0, f: func(f *formatter) {}},
		"XXX":  {max: 0, f: func(f *formatter) {}},
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
	return formatArg{max: len(ph), f: func(f *formatter) { f.fb.Write(ph) }} // Do not modify
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
