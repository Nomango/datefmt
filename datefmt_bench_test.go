package datefmt_test

import (
	"testing"

	"github.com/Nomango/datefmt"
)

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
