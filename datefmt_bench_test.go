package datefmt_test

import (
	"testing"
	"time"

	"github.com/Nomango/datefmt"
)

func BenchmarkFormat(b *testing.B) {
	for _, tt := range formatTestCases {
		tt := tt
		for _, c := range tt.testCases {
			c := c

			l := datefmt.NewLayout(tt.layout)
			b.Run(tt.layout, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = l.Format(c.in)
				}
			})
		}
	}
}

func BenchmarkFormatConcurrent(b *testing.B) {
	t := time.Date(2022, time.June, 20, 9, 49, 10, 0, time.UTC)
	b.ResetTimer()

	b.ReportAllocs()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = datefmt.Format(t, "yyyy-MM-dd HH:mm:ss z")
		}
	})
}
