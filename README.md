# datefmt

[![Go Reference](https://pkg.go.dev/badge/github.com/Nomango/datefmt.svg)](https://pkg.go.dev/github.com/Nomango/datefmt)
[![Github status](https://github.com/Nomango/datefmt/actions/workflows/UnitTest.yml/badge.svg?branch=main)](https://github.com/Nomango/datefmt/actions)
[![GitHub release](https://img.shields.io/github/release/nomango/datefmt)](https://github.com/Nomango/datefmt/releases/latest)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/9dd32aab87364903a57f32543f3bf738)](https://www.codacy.com/gh/Nomango/datefmt/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Nomango/datefmt&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/9dd32aab87364903a57f32543f3bf738)](https://www.codacy.com/gh/Nomango/datefmt/dashboard?utm_source=github.com&utm_medium=referral&utm_content=Nomango/datefmt&utm_campaign=Badge_Coverage)
[![License](https://img.shields.io/github/license/nomango/datefmt)](https://github.com/Nomango/datefmt/blob/main/LICENSE)

A **FASTER** date format tool based on **`yyyyMMdd`** layout for Golang.

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

## Pattern

The format of the layout is similar to the time and date pattern defined in Java.

Support for the standard is as follows:

| pattern | example   | format (fast, default) | format (std) | parse (std) |
| :---    | :---      |:-:|:-:|:-:|
| G       | AD        | ✓ |   |   |
| YYYY    | 2006      | ✓ | [^1] | [^1] |
| yyyy    | 2006      | ✓ | ✓ | ✓ |
| YY      | 06        | ✓ | [^1] | [^1] |
| yy      | 06        | ✓ | ✓ | ✓ |
| MMMM    | February  | ✓ | ✓ | ✓ |
| MMM     | Feb       | ✓ | ✓ | ✓ |
| MM      | 02        | ✓ | ✓ | ✓ |
| M       | 2         | ✓ | ✓ | ✓ |
| FF      | 07        | ✓ |   |   |
| F       | 7         | ✓ |   |   |
| WW      | 07        | ✓ |   |   |
| W       | 7         | ✓ |   |   |
| ww      | 07        | ✓ |   |   |
| w       | 7         | ✓ |   |   |
| DDD     | 032       | ✓ | ✓ | ✓ |
| D       | 32        | ✓ |   |   |
| dd      | 01        | ✓ | ✓ | ✓ |
| d       | 1         | ✓ | ✓ | ✓ |
| EEEE    | Wednesday | ✓ | ✓ | ✓ |
| EEE     | Web       | ✓ | ✓ | ✓ |
| u       | 3         | ✓ |   |   |
| HH      | 15        | ✓ | ✓ | ✓ |
| hh      | 03        | ✓ | ✓ | ✓ |
| H       | 15        | ✓ |   |   |
| h       | 3         | ✓ | ✓ | ✓ |
| KK      | 11        | ✓ |   |   |
| K       | 11        | ✓ |   |   |
| kk      | 24        | ✓ |   |   |
| k       | 24        | ✓ |   |   |
| mm      | 04        | ✓ | ✓ | ✓ |
| m       | 4         | ✓ | ✓ | ✓ |
| ss      | 05        | ✓ | ✓ | ✓ |
| s       | 5         | ✓ | ✓ | ✓ |
| SSS     | 978       | ✓ | ✓ | ✓ |
| a       | PM        | ✓ | ✓ | ✓ |
| z       | MST       | ✓ | ✓ | ✓ |
| Z       | -0700     | ✓ | ✓ | ✓ |
| X       | Z07       | ✓ | ✓ | ✓ |
| XX      | Z0700     | ✓ | ✓ | ✓ |
| XXX     | Z07:00    | ✓ | ✓ | ✓ |
| '       | 'o''clock'| ✓ | [^2] | [^2] |

> **Note**  
> [^1]: Not week year  
> [^2]: Only support text delimiter in layout convertion  
