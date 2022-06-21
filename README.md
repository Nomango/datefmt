# datefmt

[![Go Reference](https://pkg.go.dev/badge/github.com/Nomango/datefmt.svg)](https://pkg.go.dev/github.com/Nomango/datefmt)
[![Github status](https://github.com/Nomango/datefmt/actions/workflows/UnitTest.yml/badge.svg?branch=main)](https://github.com/Nomango/datefmt/actions)
[![GitHub release](https://img.shields.io/github/release/nomango/datefmt)](https://github.com/Nomango/datefmt/releases/latest)
[![GitHub license](https://img.shields.io/github/license/nomango/datefmt)](https://github.com/Nomango/datefmt/blob/main/LICENSE)

A date format tool based on general layout for Golang.

It is designed to solve this problem [Golang/go issue: magic number datetime formatting](https://github.com/golang/go/issues/38871).

## Usage

```golang
// Format time with general layout
s := datefmt.Format(time.Now(), "yyyy-MM-dd HH:mm:ss")

// Parse time with general layout
t, err := datefmt.Parse("yyyy-MM-dd HH:mm:ss", "2022-06-20 09:49:10")

// Parse time in location with general layout
t, err := datefmt.ParseInLocation("yyyy-MM-dd HH:mm:ss", "2022-06-20 09:49:10", time.UTC)

// Get go-style layout
l := datefmt.GoLayout("yyyy-MM-dd HH:mm:ss") // l = '2006-01-02 15:04:05'
```

The format of the layout is similar to the time pattern definition in Java.

Support for the standard is as follows:

```golang
placeholders = map[string]string{
    "YYYY": "2006",
    "yyyy": "2006",
    "YY":   "06",
    "yy":   "06",
    "MMMM": "January",
    "MMM":  "Jan",
    "MM":   "01",
    "M":    "1",
    "DDD":  "002",
    "dddd": "Monday",
    "ddd":  "Mon",
    "dd":   "02",
    "d":    "2",

    "EEEE": "Monday",
    "EEE":  "Mon",

    "HH":  "15",
    "hh":  "03",
    "h":   "3",
    "mm":  "04",
    "m":   "4",
    "ss":  "05",
    "s":   "5",
    "SSS": "000",

    "a": "PM",

    "z":   "MST",
    "Z":   "-0700",
    "X":   "Z07",
    "XX":  "Z0700",
    "XXX": "Z07:00",
}
```
