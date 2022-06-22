package datefmt_test

import (
	"strings"
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
			layout:   `yyyy-MM-dd HH:mm:ss.SSS a z`,
			goLayout: `2006-01-02 15:04:05.000 PM MST`,
			formated: `2022-06-20 09:49:10.000 AM CST`,
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

func BenchmarkNormalFormat(b *testing.B) {
	for _, tt := range tts2 {
		tt := tt
		_ = datefmt.GoLayout(tt.layout)
		b.Run(tt.layout, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = datefmt.Format(ttt, tt.layout)
			}
		})
	}
}

func BenchmarkStdFormatParallel(b *testing.B) {
	for _, tt := range tts2 {
		tt := tt
		b.Run(tt.layout, func(b *testing.B) {
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = ttt.Format(tt.goLayout)
				}
			})
		})
	}
}

func BenchmarkFastFormatParallel(b *testing.B) {
	for _, tt := range tts2 {
		tt := tt
		l := datefmt.FastLayout(tt.layout)
		b.Run(tt.layout, func(b *testing.B) {
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = l.Format(ttt)
				}
			})
		})
	}
}

func BenchmarkTest(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var sb strings.Builder
			sb.WriteByte('0')
		}
	})

	b.Run("test2", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := f{}
			test(&f)
		}
	})

	b.Run("test3", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := f{}
			test2(&f.sb)
		}
	})

	fs := []func(sb *strings.Builder){test2, test2}
	b.Run("test4", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := f{}
			for j := range fs {
				fs[j](&f.sb)
			}
		}
	})

	b.Run("test5", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := f{}
			fs[0](&f.sb)
			fs[1](&f.sb)
		}
	})

	b.Run("test6", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := f{}
			fss[0](&f.sb)
			fss[1](&f.sb)
		}
	})

	b.Run("test7", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := f{}
			fs := []func(sb *strings.Builder){f.test3, f.test3}
			for j := range fs {
				fs[j](nil)
			}
		}
	})

	b.Run("test8", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			f := f{}
			fs := []func(sb *strings.Builder){test2, test2}
			for j := range fs {
				fs[j](&f.sb)
			}
		}
	})

	b.Run("test9", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var sb strings.Builder
			ff := ff{}
			fs := []func(sb *strings.Builder){ff.test4, ff.test4}
			for j := range fs {
				fs[j](&sb)
			}
		}
	})

	b.Run("test10", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var sb strings.Builder
			for j := range fs {
				fs[j](&sb)
			}
		}
	})
}

var fss = []func(sb *strings.Builder){test2, test2}

type f struct {
	sb strings.Builder
	b  [6]byte
}

func test(f *f) {
	f.sb.WriteByte('0')
}

func test2(sb *strings.Builder) {
	sb.WriteByte('0')
}

func (f *f) test3(*strings.Builder) {
	f.sb.WriteByte('0')
}

type ff struct {
	b [100]byte
}

func (f *ff) test4(sb *strings.Builder) {
	sb.WriteByte('0')
}
