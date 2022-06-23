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
	zCST, _     = time.LoadLocation("Asia/Shanghai")
	testingTime = time.Unix(1655776150, 0).In(zCST)
	tts         = []struct {
		t        time.Time
		layout   string
		goLayout string
		formated string
		onlyStd  bool
		onlyFast bool
	}{
		{
			t:        testingTime,
			layout:   `yyyy-MM-dd HH:mm:ss z`,
			goLayout: `2006-01-02 15:04:05 MST`,
			formated: `2022-06-21 09:49:10 CST`,
		},
		{
			t:        testingTime,
			layout:   `yyyy.MM.dd 'at' HH:mm:ss z`,
			goLayout: `2006.01.02 at 15:04:05 MST`,
			formated: `2022.06.21 at 09:49:10 CST`,
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
		{
			t:        testingTime,
			layout:   `hh 'o''clock' a, z`,
			goLayout: `03 o'clock PM, MST`,
			formated: `09 o'clock AM, CST`,
		},
		{
			t:        testingTime,
			layout:   `yyyyy.MMMMM.dd hh:mm aaa`,
			goLayout: `2006.January.02 03:04 aaa`,
			formated: `2022.June.21 09:49 aaa`,
		},
		{
			t:        testingTime,
			layout:   `EEE, d MMM yyyy HH:mm:ss Z`,
			goLayout: `Mon, 2 Jan 2006 15:04:05 -0700`,
			formated: `Tue, 21 Jun 2022 09:49:10 +0800`,
		},
		{
			t:        testingTime,
			layout:   `yyMMddHHmmssZ`,
			goLayout: `060102150405-0700`,
			formated: `220621094910+0800`,
		},
		{
			t:        time.Unix(1655776150, 181999999).In(zCST),
			layout:   `yyyy-MM-dd'T'HH:mm:ss.SSSZ`,
			goLayout: `2006-01-02T15:04:05.000-0700`,
			formated: `2022-06-21T09:49:10.181+0800`,
		},
		{
			t:        time.Unix(1655776150, 181999999).In(zCST),
			layout:   `YY-M-dd'T'HH:m:s.SSSXXX EEEE`,
			goLayout: `06-1-02T15:4:5.000Z07:00 Monday`,
			formated: `22-6-21T09:49:10.181+08:00 Tuesday`,
			onlyStd:  true,
		},
		{
			t:      time.Unix(1655776150, 181999999).In(zCST),
			layout: `YY-M-dd'T'H:m:s.SSSXXX EEEE`,
			// goLayout: `06-1-02T15:4:5.000Z07:00 Monday`,
			formated: `22-6-21T9:49:10.181+08:00 Tuesday`,
			onlyFast: true,
		},
		{
			t:        testingTime,
			layout:   `YYYY-'M'MM-u`,
			formated: `2022-M06-2`,
			onlyFast: true,
		},
		{
			t:        time.Unix(1642740550, 0).In(zCST),
			layout:   `DDD hh a X`,
			goLayout: `002 03 PM Z07`,
			formated: `021 12 PM +08`,
		},
		{
			t:      time.Unix(1642740550, 0).In(zCST),
			layout: `D`,
			// goLayout: `002 03 PM Z07`,
			formated: `21`,
			onlyFast: true,
		},
		{
			t:        testingTime,
			layout:   `XX`,
			goLayout: `Z0700`,
			formated: `+0800`,
		},
		{
			t:        time.Unix(1642772950, 0).In(time.FixedZone("", -120)),
			layout:   `z`,
			goLayout: `MST`,
			formated: `-0002`,
		},
		{
			t:        time.Unix(1642772950, 0).In(time.FixedZone("", 0)),
			layout:   `X`,
			goLayout: `Z07`,
			formated: `Z`,
		},
		{
			t:        time.Unix(1642772950, 81999999),
			layout:   `.SSS`,
			goLayout: `.000`,
			formated: `.081`,
		},
		{
			t:        testingTime,
			layout:   `''''`,
			goLayout: `''`,
			formated: `''`,
		},
		{
			t:        time.Date(2022, time.July, 3, 12, 0, 0, 0, time.UTC),
			layout:   `G F FF W WW w ww u k kk K KK`,
			formated: `AD 1 01 1 01 26 26 7 12 12 0 00`,
			onlyFast: true,
		},
		{
			t:        time.Date(2022, time.July, 4, 0, 0, 0, 0, time.UTC),
			layout:   `G F FF W WW w ww u k kk K KK`,
			formated: `AD 1 01 2 02 27 27 1 24 24 0 00`,
			onlyFast: true,
		},
		{
			t:        time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			layout:   `ww WW FF`,
			formated: `52 01 01`,
			onlyFast: true,
		},
		{
			t:        time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
			layout:   `ww WW FF`,
			formated: `01 02 01`,
			onlyFast: true,
		},
		{
			t:        time.Date(2024, time.December, 30, 0, 0, 0, 0, time.UTC),
			layout:   `ww WW FF`,
			formated: `01 06 05`,
			onlyFast: true,
		},
		{
			t:        time.Date(-1, time.January, 1, 0, 0, 0, 0, time.UTC),
			layout:   `G`,
			formated: `BC`,
			onlyFast: true,
		},
	}
)

func TestTest(t *testing.T) {
	for i := 1; i < 40; i++ {
		t := time.Date(2024, time.December, i, 0, 1, 0, 0, time.UTC)
		fmt.Println(datefmt.Format(t, "YYYY yyyy-MM-dd kk:mm:ss DDD ww WW FF EEEE z Z"))
	}
	// for i := 1; i < 40; i++ {
	// 	t := time.Date(2022, time.January, i, 0, 1, 0, 0, time.UTC)
	// 	fmt.Println(datefmt.Format(t, "YYYY yyyy-MM-dd kk:mm:ss DDD ww WW FF EEEE z Z"))
	// }
	// for i := 1; i < 40; i++ {
	// 	t := time.Date(2022, time.July, i, 0, 1, 0, 0, time.UTC)
	// 	fmt.Println(datefmt.Format(t, "YYYY yyyy-MM-dd kk:mm:ss DDD ww WW FF EEEE z Z"))
	// }
}

func TestGoLayout(t *testing.T) {
	for _, tt := range tts {
		if tt.onlyFast {
			continue
		}
		l := datefmt.GoLayout(tt.layout)
		if l != tt.goLayout {
			t.Errorf("GoLayout(%s) = %s; want %s", tt.layout, l, tt.goLayout)
		}
	}
}

func TestFormat(t *testing.T) {
	for _, tt := range tts {
		r := datefmt.Format(tt.t, tt.layout)
		if r != tt.formated {
			t.Errorf("Format(%d, %s) = %s; want %s", tt.t.Unix(), tt.layout, r, tt.formated)
		}
	}
}

func TestFormatStable(t *testing.T) {
	for _, tt := range tts {
		tt := tt
		t.Run(tt.layout, func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 5; i++ {
				r := datefmt.Format(tt.t, tt.layout)
				if r != tt.formated {
					t.Errorf("Format(%d, %s) = %s; want %s", tt.t.Unix(), tt.layout, r, tt.formated)
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	for _, tt := range tts {
		if tt.onlyFast {
			continue
		}
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

func TestUse(t *testing.T) {
	datefmt.Use(datefmt.StdFormat)
	t.Cleanup(func() {
		datefmt.Use(datefmt.FastFormat)
	})
	for _, tt := range tts {
		if tt.onlyFast {
			continue
		}
		r := datefmt.Format(tt.t, tt.layout)
		if r != tt.formated {
			t.Errorf("Format(%d, %s) = %s; want %s", tt.t.Unix(), tt.layout, r, tt.formated)
		}
	}
}
