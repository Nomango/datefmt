package datefmt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Nomango/datefmt"
)

func ExampleFormat() {
	t := time.Unix(1655689750, 0).In(time.UTC)
	s := datefmt.Format(t, "yyyy-MM-dd HH:mm:ss Z z")
	fmt.Println(s)
	// Output:
	// 2022-06-20 01:49:10 +0000 UTC
}

func ExampleParse() {
	t, _ := datefmt.Parse("yyyy-MM-dd HH:mm:ss Z z", "2022-06-20 09:49:10 +0800 CST")
	t = t.In(time.UTC)
	fmt.Println(t)
	// Output:
	// 2022-06-20 01:49:10 +0000 UTC
}

func ExampleParseInLocation() {
	t, _ := datefmt.ParseInLocation("yyyy-MM-dd HH:mm:ss z", "2022-06-20 09:49:10 CST", time.UTC)
	fmt.Println(t)
	// Output:
	// 2022-06-20 09:49:10 +0000 CST
}

func ExampleGoLayout() {
	l := datefmt.GoLayout("yyyy-MM-dd HH:mm:ss z")
	fmt.Println(l)
	// Output:
	// 2006-01-02 15:04:05 MST
}

var (
	zCST, _ = time.LoadLocation("Asia/Shanghai")
	ttt     = time.Unix(1655689750, 0).In(zCST)
	tts     = []struct {
		layout   string
		goLayout string
		formated string
	}{
		{
			layout:   `yyyy-MM-dd HH:mm:ss z`,
			goLayout: `2006-01-02 15:04:05 MST`,
			formated: `2022-06-20 09:49:10 CST`,
		},
		{
			layout:   `yyyy.MM.dd 'at' HH:mm:ss z`,
			goLayout: `2006.01.02 at 15:04:05 MST`,
			formated: `2022.06.20 at 09:49:10 CST`,
		},
		{
			layout:   `EEE, MMM DDD, ''yy`,
			goLayout: `Mon, Jan 002, '06`,
			formated: `Mon, Jun 171, '22`,
		},
		{
			layout:   `h:mm a`,
			goLayout: `3:04 PM`,
			formated: `9:49 AM`,
		},
		{
			layout:   `hh 'o''clock' a, z`,
			goLayout: `03 o'clock PM, MST`,
			formated: `09 o'clock AM, CST`,
		},
		{
			layout:   `yyyyy.MMMMM.dd hh:mm aaa`,
			goLayout: `2006.January.02 03:04 aaa`,
			formated: `2022.June.20 09:49 aaa`,
		},
		{
			layout:   `EEE, d MMM yyyy HH:mm:ss Z`,
			goLayout: `Mon, 2 Jan 2006 15:04:05 -0700`,
			formated: `Mon, 20 Jun 2022 09:49:10 +0800`,
		},
		{
			layout:   `yyMMddHHmmssZ`,
			goLayout: `060102150405-0700`,
			formated: `220620094910+0800`,
		},
		{
			layout:   `yyyy-MM-dd'T'HH:mm:ss.SSSZ`,
			goLayout: `2006-01-02T15:04:05.000-0700`,
			formated: `2022-06-20T09:49:10.000+0800`,
		},
		{
			layout:   `yyyy-MM-dd'T'HH:mm:ss.SSSXXX`,
			goLayout: `2006-01-02T15:04:05.000Z07:00`,
			formated: `2022-06-20T09:49:10.000+08:00`,
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

func TestParse(t *testing.T) {
	for _, tt := range tts {
		r1, err := datefmt.Parse(tt.layout, tt.formated)
		if err != nil {
			t.Errorf("datefmt.Parse(%s, %s) got an error: %v", tt.layout, tt.formated, err)
		}
		r2, err := time.Parse(tt.goLayout, tt.formated)
		if err != nil {
			t.Errorf("time.Parse(%s, %s) got an error: %v", tt.goLayout, tt.formated, err)
		}
		if !r1.Equal(r2) {
			t.Errorf("Parse(%s, %s) = %v; want %v", tt.layout, tt.formated, r1, r2)
		}
	}
}
