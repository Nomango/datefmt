package datefmt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Nomango/datefmt"
)

func ExampleFormat() {
	t := time.Unix(1655732950, 0)
	s := datefmt.Format(t, "yyyy-MM-dd HH:mm:ss z")
	fmt.Println(s)
	// Output:
	// 2022-06-20 21:49:10 CST
}

func ExampleGoLayout() {
	l := datefmt.GoLayout("yyyy-MM-dd HH:mm:ss z")
	fmt.Println(l)
	// Output:
	// 2006-01-02 15:04:05 MST
}

var (
	zCST, _ = time.LoadLocation("Asia/Shanghai")
	ttt     = time.Unix(1655732950, 0).In(zCST)
	tts     = []struct {
		layout   string
		goLayout string
		formated string
	}{
		{
			layout:   `yyyy-MM-dd HH:mm:ss z`,
			goLayout: `2006-01-02 15:04:05 MST`,
			formated: `2022-06-20 21:49:10 CST`,
		},
		{
			layout:   `yyyy.MM.dd 'at' HH:mm:ss z`,
			goLayout: `2006.01.02 at 15:04:05 MST`,
			formated: `2022.06.20 at 21:49:10 CST`,
		},
		{
			layout:   `EEE, MMM DDD, ''yy`,
			goLayout: `Mon, Jan 002, '06`,
			formated: `Mon, Jun 171, '22`,
		},
		{
			layout:   `h:mm a`,
			goLayout: `3:04 PM`,
			formated: `9:49 PM`,
		},
		{
			layout:   `hh 'o''clock' a, z`,
			goLayout: `03 o'clock PM, MST`,
			formated: `09 o'clock PM, CST`,
		},
		{
			layout:   `yyyyy.MMMMM.dd hh:mm aaa`,
			goLayout: `2006.January.02 03:04 aaa`,
			formated: `2022.June.20 09:49 aaa`,
		},
		{
			layout:   `EEE, d MMM yyyy HH:mm:ss Z`,
			goLayout: `Mon, 2 Jan 2006 15:04:05 -0700`,
			formated: `Mon, 20 Jun 2022 21:49:10 +0800`,
		},
		{
			layout:   `yyMMddHHmmssZ`,
			goLayout: `060102150405-0700`,
			formated: `220620214910+0800`,
		},
		{
			layout:   `yyyy-MM-dd'T'HH:mm:ss.SSSZ`,
			goLayout: `2006-01-02T15:04:05.000-0700`,
			formated: `2022-06-20T21:49:10.000+0800`,
		},
		{
			layout:   `yyyy-MM-dd'T'HH:mm:ss.SSSXXX`,
			goLayout: `2006-01-02T15:04:05.000Z07:00`,
			formated: `2022-06-20T21:49:10.000+08:00`,
		},
		{
			layout:   `YYYY-'M'MM-u`,
			goLayout: `2006-M01-u`,
			formated: `2022-M06-u`,
		},
		{
			layout:   `''''`,
			goLayout: `''`,
			formated: `''`,
		},
	}
)

func TestGoLayout(t *testing.T) {
	for _, tt := range tts {
		l := datefmt.GoLayout(tt.layout)
		if l != tt.goLayout {
			t.Errorf("GoLayout(%s) = %s; want %s", tt.layout, l, tt.goLayout)
		}
	}
}

func TestFormat(t *testing.T) {
	for _, tt := range tts {
		r := datefmt.Format(ttt, tt.layout)
		if r != tt.formated {
			t.Errorf("Format(%d, %s) = %s; want %s", ttt.Unix(), tt.layout, r, tt.formated)
		}
	}
}

func TestFormatStable(t *testing.T) {
	for _, tt := range tts {
		tt := tt
		t.Run(tt.layout, func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 5; i++ {
				r := datefmt.Format(ttt, tt.layout)
				if r != tt.formated {
					t.Errorf("Format(%d, %s) = %s; want %s", ttt.Unix(), tt.layout, r, tt.formated)
				}
			}
		})
	}
}

func BenchmarkFormat(b *testing.B) {
	for _, tt := range tts {
		tt := tt
		b.Run(tt.layout, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = datefmt.Format(ttt, tt.layout)
			}
		})
	}
}

func BenchmarkFormatConcurrent(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = datefmt.Format(ttt, "yyyy-MM-dd HH:mm:ss z")
		}
	})
}
