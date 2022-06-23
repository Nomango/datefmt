package datefmt_test

import (
	"testing"

	"github.com/Nomango/datefmt"
)

func BenchmarkFormat(b *testing.B) {
	for _, tt := range tts {
		tt := tt
		_ = datefmt.GoLayout(tt.layout)
		b.Run(tt.layout, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = datefmt.Format(tt.t, tt.layout)
			}
		})
	}
}

func BenchmarkFormatConcurrent(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = datefmt.Format(testingTime, "yyyy-MM-dd HH:mm:ss z")
		}
	})
}
