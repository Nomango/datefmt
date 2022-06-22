package datefmt_test

import (
	"testing"

	"github.com/Nomango/datefmt"
)

var (
	tts2 = []struct {
		layout   string
		goLayout string
		formated string
	}{
		{
			layout:   `yyyy-MM-dd HH:mm:ss z`,
			goLayout: `2006-01-02 15:04:05 MST`,
			formated: `2022-06-20 09:49:10 CST`,
		},
	}
)

func BenchmarkStdFormat(b *testing.B) {
	for _, tt := range tts2 {
		tt := tt
		b.Run(tt.layout, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = ttt.Format(tt.goLayout)
			}
		})
	}
}

func BenchmarkFastFormat(b *testing.B) {
	for _, tt := range tts2 {
		tt := tt
		l := datefmt.FastLayout(tt.layout)
		b.Run(tt.layout, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = l.Format(ttt)
			}
		})
	}
}
