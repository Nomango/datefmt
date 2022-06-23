# datefmt

[![Go Reference](https://pkg.go.dev/badge/github.com/Nomango/datefmt.svg)](https://pkg.go.dev/github.com/Nomango/datefmt)
[![Github status](https://github.com/Nomango/datefmt/actions/workflows/UnitTest.yml/badge.svg?branch=main)](https://github.com/Nomango/datefmt/actions)
[![GitHub release](https://img.shields.io/github/release/nomango/datefmt)](https://github.com/Nomango/datefmt/releases/latest)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/9dd32aab87364903a57f32543f3bf738)](https://www.codacy.com/gh/Nomango/datefmt/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Nomango/datefmt&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/9dd32aab87364903a57f32543f3bf738)](https://www.codacy.com/gh/Nomango/datefmt/dashboard?utm_source=github.com&utm_medium=referral&utm_content=Nomango/datefmt&utm_campaign=Badge_Coverage)
[![License](https://img.shields.io/github/license/nomango/datefmt)](https://github.com/Nomango/datefmt/blob/main/LICENSE)

A **faster** date format tool based on **`yyyyMMdd`** layout for Golang.

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

The format of the layout is similar to the [time and date pattern defined in Java](https://docs.oracle.com/javase/7/docs/api/java/text/SimpleDateFormat.html).

Support for the standard is as follows:

| letter | description              | example            | format (fast, default) | format (std) | parse (std) |
| :---   | :---                     | :---               |:-:|:-:|:-:|
| G      | Era designator           | AD                 | ✓ |   |   |
| y      | Year                     | 1996; 96           | ✓ | ✓[^1] | ✓[^1] |
| Y      | Week year                | 2009; 09           | ✓ | ✓[^2] | ✓[^2] |
| M      | Month in year            | July; Jul; 07      | ✓ | ✓ | ✓ |
| w      | Week in year             | 27                 | ✓ |   |   |
| W      | Week in month            | 2                  | ✓ |   |   |
| D      | Day in year              | 189                | ✓ | ✓ | ✓ |
| d      | Day in month             | 10                 | ✓ | ✓ | ✓ |
| F      | Day of week in month     | 2                  | ✓ |   |   |
| E      | Day name in week         | Tuesday; Tue       | ✓ | ✓ | ✓ |
| u      | Day number of week (1-7) | 1                  | ✓ |   |   |
| a      | Am/pm marker             | PM                 | ✓ | ✓ | ✓ |
| H      | Hour in day (0-23)       | 0                  | ✓ | ✓ | ✓ |
| k      | Hour in day (1-24)       | 24                 | ✓ |   |   |
| K      | Hour in am/pm (0-11)     | 0                  | ✓ |   |   |
| h      | Hour in am/pm (1-12)     | 12                 | ✓ | ✓ | ✓ |
| m      | Minute in hour           | 30                 | ✓ | ✓ | ✓ |
| s      | Second in minute         | 55                 | ✓ | ✓ | ✓ |
| S      | Millisecond              | 978                | ✓ | ✓ | ✓ |
| z      | Time zone                | PST; GMT-08:00     | ✓ | ✓ | ✓ |
| Z      | Time zone                | -800               | ✓ | ✓ | ✓ |
| X      | Time zone                | -08; -0800; -08:00 | ✓ | ✓ | ✓ |
| '      | Text delimiter           | 'o''clock'         | ✓ | ✓[^3] | ✓[^3] |

> **Note**  
> [^1]: Only support common placeholders in std format & parse, eg, `yyyy` and `yy` is valid, but `yyy` is not. Such as the others.  
> [^2]: 'Y' treated as 'y' in std format & parse.  
> [^3]: Only support text delimiter in layout convertion.  

## Performance

`datefmt` performs quite well and in most cases has better performance than the standard library.

Here is the benchmark results for common cases.

```plain
BenchmarkCompareFormat/datefmt______yyyy-MM-dd_HH:mm:ss-8           7447152        161.7 ns/op       24 B/op        1 allocs/op
BenchmarkCompareFormat/std.Format___yyyy-MM-dd_HH:mm:ss-8           6057490        214.9 ns/op       24 B/op        1 allocs/op
BenchmarkCompareFormat/jodaTime_____yyyy-MM-dd_HH:mm:ss-8           1455223        817.3 ns/op      208 B/op       18 allocs/op
BenchmarkCompareFormat/GoDateFormat_yyyy-MM-dd_HH:mm:ss-8           1495375        810.4 ns/op      168 B/op        7 allocs/op

BenchmarkCompareFormat/datefmt______yyyy-MM-dd'T'HH:mm:ss_z-8       5913645        203.9 ns/op       32 B/op        1 allocs/op
BenchmarkCompareFormat/std.Format___yyyy-MM-dd'T'HH:mm:ss_z-8       5850758        217.3 ns/op       24 B/op        1 allocs/op
BenchmarkCompareFormat/jodaTime_____yyyy-MM-dd'T'HH:mm:ss_z-8       1214950        999.5 ns/op      256 B/op       22 allocs/op
BenchmarkCompareFormat/GoDateFormat_yyyy-MM-dd'T'HH:mm:ss_z-8       1296421        903.4 ns/op      208 B/op        8 allocs/op

BenchmarkCompareFormat/datefmt______yyMMddHHmmssZ-8                 6011174        179.6 ns/op       24 B/op        1 allocs/op
BenchmarkCompareFormat/std.Format___yyMMddHHmmssZ-8                 6345073        192.5 ns/op       24 B/op        1 allocs/op
BenchmarkCompareFormat/jodaTime_____yyMMddHHmmssZ-8                 1683231        710.1 ns/op      136 B/op       12 allocs/op
BenchmarkCompareFormat/GoDateFormat_yyMMddHHmmssZ-8                 1394109        869.9 ns/op      128 B/op        8 allocs/op

BenchmarkCompareFormat/datefmt______EEE,_MMM_DDD,_''yy-8            6854894        174.7 ns/op       24 B/op        1 allocs/op
BenchmarkCompareFormat/std.Format___EEE,_MMM_DDD,_''yy-8            8224816        160.3 ns/op       24 B/op        1 allocs/op
BenchmarkCompareFormat/jodaTime_____EEE,_MMM_DDD,_''yy-8            2169121        545.2 ns/op      112 B/op       14 allocs/op
BenchmarkCompareFormat/GoDateFormat_EEE,_MMM_DDD,_''yy-8            2567779        468.8 ns/op       72 B/op        3 allocs/op

BenchmarkCompareFormat/datefmt______h:mm_a-8                       13415570        86.86 ns/op        8 B/op        1 allocs/op
BenchmarkCompareFormat/std.Format___h:mm_a-8                       10746708        112.0 ns/op        8 B/op        1 allocs/op
BenchmarkCompareFormat/jodaTime_____h:mm_a-8                        4197170        289.3 ns/op       32 B/op        6 allocs/op
BenchmarkCompareFormat/GoDateFormat_h:mm_a-8                        3587061        340.9 ns/op       16 B/op        2 allocs/op

PASS
ok   bench 34.733s
```

> **Note**  
> jodaTime: [github.com/vjeantet/jodaTime](github.com/vjeantet/jodaTime)  
> GoDateFormat: [github.com/vigneshuvi/GoDateFormat](github.com/vigneshuvi/GoDateFormat)  
