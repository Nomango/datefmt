package bench_test

import (
	"testing"
	"time"

	"github.com/Nomango/datefmt"
	"github.com/vigneshuvi/GoDateFormat"
	"github.com/vjeantet/jodaTime"
)

var (
	zoneCST, _  = time.LoadLocation("Asia/Shanghai")
	testingTime = time.Unix(1655776150, 0).In(zoneCST)

	tts2 = []struct {
		t        time.Time
		layout   string
		goLayout string
		formated string
	}{
		{
			t:        testingTime,
			layout:   `yyyy-MM-dd HH:mm:ss`,
			goLayout: `2006-01-02 15:04:05`,
			formated: `2022-06-21 09:49:10`,
		},
		{
			t:        testingTime,
			layout:   `yyyy-MM-dd'T'HH:mm:ss z`,
			goLayout: `2006-01-02T15:04:05 MST`,
			formated: `2022-06-21T09:49:10 CST`,
		},
		{
			t:        testingTime,
			layout:   `yyMMddHHmmssZ`,
			goLayout: `060102150405-0700`,
			formated: `220621094910+0800`,
		},
		{
			t:        testingTime,
			layout:   `EEE, MMM DDD, ''yy`,
			goLayout: `Mon, Jan 002, '06`,
			formated: `Tue, Jun 172, '22`,
		},
		{
			t:        testingTime,
			layout:   `h:mm a`,
			goLayout: `3:04 PM`,
			formated: `9:49 AM`,
		},
	}
)

func BenchmarkFormat(b *testing.B) {
	for _, tt := range tts2 {
		tt := tt

		l := datefmt.NewLayout(tt.layout)
		b.Run("datefmt     "+tt.layout, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = l.Format(tt.t)
			}
		})

		b.Run("std.Format  "+tt.layout, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = tt.t.Format(tt.goLayout)
			}
		})

		b.Run("jodaTime    "+tt.layout, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = jodaTime.Format(tt.layout, tt.t)
			}
		})

		b.Run("GoDateFormat"+tt.layout, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = tt.t.Format(GoDateFormat.ConvertFormat(tt.layout))
			}
		})
	}
}
