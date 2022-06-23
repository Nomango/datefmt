package datefmt_test

import (
	"testing"

	"github.com/Nomango/datefmt"
)

func BenchmarkFastFormat(b *testing.B) {
	for _, tt := range tts {
		tt := tt
		if tt.onlyStd {
			continue
		}

		l := datefmt.FastLayout(tt.layout)
		b.Run(tt.layout, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = l.Format(tt.t)
			}
		})
	}
}
