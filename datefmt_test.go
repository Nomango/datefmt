package datefmt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Nomango/datefmt"
)

func ExampleFormat() {
	t := time.Date(2022, time.June, 20, 9, 49, 10, 0, time.UTC)
	s := datefmt.Format(t, "yyyy-MM-dd HH:mm:ss")
	fmt.Println(s)
	// Output:
	// 2022-06-20 09:49:10
}

func ExampleParse() {
	t, _ := datefmt.Parse("yyyy-MM-dd HH:mm:ss z Z", "2022-06-20 09:49:10 CST +0800")
	fmt.Println(t)
	// Output:
	// 2022-06-20 09:49:10 +0800 CST
}

func ExampleParseInLocation() {
	t, _ := datefmt.ParseInLocation("yyyy-MM-dd HH:mm:ss", "2022-06-20 09:49:10", time.UTC)
	fmt.Println(t)
	// Output:
	// 2022-06-20 09:49:10 +0000 UTC
}

func ExampleGoLayout() {
	l := datefmt.GoLayout("yyyy-MM-dd HH:mm:ss z")
	fmt.Println(l)
	// Output:
	// 2006-01-02 15:04:05 MST
}

type testCase struct {
	in  time.Time
	out string
}

var (
	formatTestCases = []struct {
		layout    string
		testCases []testCase
	}{
		{
			layout: "G G",
			testCases: []testCase{
				{
					in:  time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					out: "AD AD",
				},
				{
					in:  time.Date(-1, time.January, 1, 0, 0, 0, 0, time.UTC),
					out: "BC BC",
				},
			},
		},
		{
			layout: "y yy yyy yyyy yyyyy Y YY YYY YYYY YYYYY",
			testCases: []testCase{
				{
					in:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
					out: "2022 22 2022 2022 02022 2021 21 2021 2021 02021",
				},
			},
		},
		{
			layout: "M MM MMM MMMM MMMMM",
			testCases: []testCase{
				{
					in:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
					out: "1 01 Jan January January",
				},
				{
					in:  time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC),
					out: "2 02 Feb February February",
				},
				{
					in:  time.Date(2022, time.March, 1, 0, 0, 0, 0, time.UTC),
					out: "3 03 Mar March March",
				},
				{
					in:  time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC),
					out: "4 04 Apr April April",
				},
				{
					in:  time.Date(2022, time.May, 1, 0, 0, 0, 0, time.UTC),
					out: "5 05 May May May",
				},
				{
					in:  time.Date(2022, time.June, 1, 0, 0, 0, 0, time.UTC),
					out: "6 06 Jun June June",
				},
				{
					in:  time.Date(2022, time.July, 1, 0, 0, 0, 0, time.UTC),
					out: "7 07 Jul July July",
				},
				{
					in:  time.Date(2022, time.August, 1, 0, 0, 0, 0, time.UTC),
					out: "8 08 Aug August August",
				},
				{
					in:  time.Date(2022, time.September, 1, 0, 0, 0, 0, time.UTC),
					out: "9 09 Sep September September",
				},
				{
					in:  time.Date(2022, time.October, 1, 0, 0, 0, 0, time.UTC),
					out: "10 10 Oct October October",
				},
				{
					in:  time.Date(2022, time.November, 1, 0, 0, 0, 0, time.UTC),
					out: "11 11 Nov November November",
				},
				{
					in:  time.Date(2022, time.December, 1, 0, 0, 0, 0, time.UTC),
					out: "12 12 Dec December December",
				},
			},
		},
		{
			layout: "w ww www W WW WWW F FF FFF",
			testCases: []testCase{
				{
					in:  time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
					out: "1 01 001 2 02 002 1 01 001",
				},
				{
					in:  time.Date(2022, time.July, 3, 12, 0, 0, 0, time.UTC),
					out: "26 26 026 1 01 001 1 01 001",
				},
				{
					in:  time.Date(2022, time.July, 4, 0, 0, 0, 0, time.UTC),
					out: "27 27 027 2 02 002 1 01 001",
				},
				{
					in:  time.Date(2022, time.July, 8, 0, 0, 0, 0, time.UTC),
					out: "27 27 027 2 02 002 2 02 002",
				},
				{
					in:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
					out: "52 52 052 1 01 001 1 01 001",
				},
				{
					in:  time.Date(2024, time.December, 30, 0, 0, 0, 0, time.UTC),
					out: "1 01 001 6 06 006 5 05 005",
				},
			},
		},
		{
			layout: "E EE EEE EEEE EEEEE u uu",
			testCases: []testCase{
				{
					in:  time.Date(2022, time.June, 20, 0, 0, 0, 0, time.UTC),
					out: "Mon Mon Mon Monday Monday 1 01",
				},
				{
					in:  time.Date(2022, time.June, 21, 0, 0, 0, 0, time.UTC),
					out: "Tue Tue Tue Tuesday Tuesday 2 02",
				},
				{
					in:  time.Date(2022, time.June, 22, 0, 0, 0, 0, time.UTC),
					out: "Wed Wed Wed Wednesday Wednesday 3 03",
				},
				{
					in:  time.Date(2022, time.June, 23, 0, 0, 0, 0, time.UTC),
					out: "Thu Thu Thu Thursday Thursday 4 04",
				},
				{
					in:  time.Date(2022, time.June, 24, 0, 0, 0, 0, time.UTC),
					out: "Fri Fri Fri Friday Friday 5 05",
				},
				{
					in:  time.Date(2022, time.June, 25, 0, 0, 0, 0, time.UTC),
					out: "Sat Sat Sat Saturday Saturday 6 06",
				},
				{
					in:  time.Date(2022, time.June, 26, 0, 0, 0, 0, time.UTC),
					out: "Sun Sun Sun Sunday Sunday 7 07",
				},
			},
		},
		{
			layout: "d dd ddd D DD DDD DDDD",
			testCases: []testCase{
				{
					in:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
					out: "1 01 001 1 01 001 0001",
				},
				{
					in:  time.Date(2022, time.January, 31, 0, 0, 0, 0, time.UTC),
					out: "31 31 031 31 31 031 0031",
				},
				{
					in:  time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC),
					out: "1 01 001 32 32 032 0032",
				},
			},
		},
		{
			layout: "H HH h hh K KK k kk a",
			testCases: []testCase{
				{
					in:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
					out: "0 00 12 12 0 00 24 24 AM",
				},
				{
					in:  time.Date(2022, time.January, 1, 1, 0, 0, 0, time.UTC),
					out: "1 01 1 01 1 01 1 01 AM",
				},
				{
					in:  time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC),
					out: "12 12 12 12 0 00 12 12 PM",
				},
				{
					in:  time.Date(2022, time.January, 1, 15, 0, 0, 0, time.UTC),
					out: "15 15 3 03 3 03 15 15 PM",
				},
			},
		},
		{
			layout: "m mm s ss S SS SSS",
			testCases: []testCase{
				{
					in:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
					out: "0 00 0 00 0 00 000",
				},
				{
					in:  time.Date(2022, time.January, 1, 0, 1, 1, 1999999, time.UTC),
					out: "1 01 1 01 0 00 001",
				},
				{
					in:  time.Date(2022, time.January, 1, 0, 59, 59, 181999999, time.UTC),
					out: "59 59 59 59 1 18 181",
				},
			},
		},
		{
			layout: "z zz Z ZZ X XX XXX XXXX",
			testCases: []testCase{
				{
					in:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
					out: "UTC UTC +0000 +0000 Z Z Z Z",
				},
				{
					in:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.FixedZone("CST", int((8*time.Hour).Seconds()))),
					out: "CST CST +0800 +0800 +08 +0800 +08:00 +08:00",
				},
				{
					in:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.FixedZone("EST", -int((5*time.Hour).Seconds()))),
					out: "EST EST -0500 -0500 -05 -0500 -05:00 -05:00",
				},
				{
					in:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.FixedZone("", int((8*time.Hour+30*time.Minute).Seconds()))),
					out: "+0830 +0830 +0830 +0830 +08 +0830 +08:30 +08:30",
				},
			},
		},
		{
			layout: "yyyy-MM-dd'T'HH:mm:ss.SSS z",
			testCases: []testCase{
				{
					in:  time.Date(2022, time.June, 20, 21, 49, 10, 181999999, time.UTC),
					out: "2022-06-20T21:49:10.181 UTC",
				},
			},
		},
		{
			layout: "'o''clock' ''''",
			testCases: []testCase{
				{
					in:  time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
					out: "o'clock ''",
				},
			},
		},
	}

	goLayoutTestCases = []struct {
		in  string
		out string
	}{
		{
			in:  `yyyy-MM-dd HH:mm:ss z`,
			out: `2006-01-02 15:04:05 MST`,
		},
		{
			in:  `yyyy.MM.dd 'at' HH:mm:ss z`,
			out: `2006.01.02 at 15:04:05 MST`,
		},
		{
			in:  `EEE, MMM DDD, ''yy`,
			out: `Mon, Jan 002, '06`,
		},
		{
			in:  `h:mm a`,
			out: `3:04 PM`,
		},
		{
			in:  `hh 'o''clock' a, z`,
			out: `03 o'clock PM, MST`,
		},
		{
			in:  `yyyy.MMMMM.dd hh:mm aaa`,
			out: `2006.January.02 03:04 PM`,
		},
		{
			in:  `EEE, d MMM yyyy HH:mm:ss Z`,
			out: `Mon, 2 Jan 2006 15:04:05 -0700`,
		},
		{
			in:  `yyMMddHHmmssZ`,
			out: `060102150405-0700`,
		},
		{
			in:  `yyyy-MM-dd'T'HH:mm:ss.SSSZ`,
			out: `2006-01-02T15:04:05.000-0700`,
		},
		{
			in:  `YY-M-dd'T'HH:m:s.SSSXXX EEEE`,
			out: `06-1-02T15:4:5.000Z07:00 Monday`,
		},
		{
			in:  `DDD hh a z X XX XXX .SSS S`,
			out: `002 03 PM MST Z07 Z0700 Z07:00 .000 S`,
		},
		{
			in:  `''''`,
			out: `''`,
		},
	}
)

func TestFormat(t *testing.T) {
	for _, tt := range formatTestCases {
		for _, c := range tt.testCases {
			r := datefmt.Format(c.in, tt.layout)
			if r != c.out {
				t.Errorf("Format(%d, %s) = %s; want %s", c.in.Unix(), tt.layout, r, c.out)
			}
		}
	}
}

func TestFormatStable(t *testing.T) {
	for _, tt := range formatTestCases {
		t.Run(tt.layout, func(t *testing.T) {
			t.Parallel()
			for _, c := range tt.testCases {
				r := datefmt.Format(c.in, tt.layout)
				if r != c.out {
					t.Errorf("Format(%d, %s) = %s; want %s", c.in.Unix(), tt.layout, r, c.out)
				}
			}
		})
	}
}

func TestGoLayout(t *testing.T) {
	for _, tt := range goLayoutTestCases {
		l := datefmt.GoLayout(tt.in)
		if l != tt.out {
			t.Errorf("GoLayout(%s) = %s; want %s", tt.in, l, tt.out)
		}
	}
}
