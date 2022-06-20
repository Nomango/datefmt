package datefmt_test

import (
	"testing"
	"time"

	"github.com/Nomango/datefmt"
)

func TestLayout(t *testing.T) {
	tts := []struct {
		in   string
		want string
	}{
		{
			in:   `yyyy.MM.dd 'at' HH:mm:ss z`,
			want: `2006.01.02 at 15:04:05 MST`,
		},
		{
			in:   `EEE, MMM DDD, ''yy`,
			want: `Mon, Jan 002, '06`,
		},
		{
			in:   `h:mm a`,
			want: `3:04 PM`,
		},
		{
			in:   `hh 'o''clock' a, z`,
			want: `03 o'clock PM, MST`,
		},
		{
			in:   `yyyyy.MMMMM.dd hh:mm aaa`,
			want: `2006.January.02 03:04 aaa`,
		},
		{
			in:   `EEE, d MMM yyyy HH:mm:ss Z`,
			want: `Mon, 2 Jan 2006 15:04:05 -0700`,
		},
		{
			in:   `yyMMddHHmmssZ`,
			want: `060102150405-0700`,
		},
		{
			in:   `yyyy-MM-dd'T'HH:mm:ss.SSSZ`,
			want: `2006-01-02T15:04:05.000-0700`,
		},
		{
			in:   `yyyy-MM-dd'T'HH:mm:ss.SSSXXX`,
			want: `2006-01-02T15:04:05.000Z07:00`,
		},
		{
			in:   `YYYY-'M'MM-u`,
			want: `2006-M01-u`,
		},
		{
			in:   `''''`,
			want: `''`,
		},
	}
	for _, tt := range tts {
		l := datefmt.Layout(tt.in)
		if l != tt.want {
			t.Errorf("Layout(%s) = %s; want %s", tt.in, l, tt.want)
		}
	}
}

func TestFormat(t *testing.T) {
	tts := []struct {
		in     int64
		layout string
		want   string
	}{
		{
			in:     1655732950,
			layout: `yyyy-MM-dd HH:mm:ss z`,
			want:   `2022-06-20 21:49:10 CST`,
		},
	}
	for _, tt := range tts {
		r := datefmt.Format(time.Unix(tt.in, 0), tt.layout)
		if r != tt.want {
			t.Errorf("Format(%d, %s) = %s; want %s", tt.in, tt.layout, r, tt.want)
		}
	}
}
