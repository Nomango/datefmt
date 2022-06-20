# datefmt

A date format tool based on general layout for Golang.

## Usage

```golang
// Format time with general layout
s := datefmt.Format(time.Now(), "yyyy-MM-dd HH:mm:ss") // variable s may be '2022-06-20 21:49:10'

// Get go-style layout
l := datefmt.Layout("yyyy-MM-dd HH:mm:ss") // l = '2006-01-02T15:04:05'
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
