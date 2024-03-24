# Time Date Parse

[![Go Reference](https://pkg.go.dev/badge/github.com/shadiestgoat/dateparse.svg)](https://pkg.go.dev/github.com/shadiestgoat/dateparse) [![Go Report Card](https://goreportcard.com/badge/github.com/shadiestgoat/dateparse)](https://goreportcard.com/report/github.com/shadiestgoat/dateparse)

This package parses dates and times, in a dynamic manner with multiple types. Note that this is not a relative natural date parsing library - it cannot by default handle very natural language!

## Install

```
go get github.com/shadiestgoat/dateparse
```

## Parsing

The package provides `Parse(str)`, which would parse using the default parser. Default parser supports the following formats:

`<time> <date> (TZ)` & `<date> (at) <time> (TZ)`, along with independent `<time> (TZ)` and `<date> (TZ)`. Independent `<time>` parses as `today, <time>`. Independent `<date>` is parsed as `<date> 00:00`. Finally, this package also supports `<timestamp>`.

`<time>`:
  - HHpm
  - HH:MM(:SS)(pm)

`<date>`:
  - "today"
  - DD/MM
  - (DD/)MM/YY(YY)
  - YYYY
  - YYYY/MM/DD
  - All Formats above, but with `-` instead of `/`
  - (DD(st/nd/d/th)) (of) Month ((of) YY(YY))
  - (DD(st/nd/d/th)) (of) Mon ((of) YY(YY))

`<timestamp>`:
  - 10 digits: unix seconds
  - 13 digits: unix milliseconds
  - 16 digits: unix microseconds
  - 19 digits: unix nanoseconds

Notes:
  - Parts in brackets `()` are optional
  - `<time>` mentions "pm" as a substitute for am/pm. Both work. 
  - While docs use DD, HH, or MM, a preceding 0 is not needed. 1 would work as well as 01. Not applicable to YY.
  - `(TZ)` defaults to UTC. If present, it must be in the IANA zoneinfo format (eg. `Europe/London`, `CET`). [Here is a full list](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)

## Custom Parsing

Parsing is done using a `Parser`. The easiest way to make your own parser is to use the built-in DateTimeParser, and include your own dates & times. These are all built using go's date/time layouts. For more info on those, see the time package.

## Timezones

All parsing provided by the package is made using using the time package. You can use it to convert timezones & set a default tz.
Additionally, supported in the default parser is the `(TZ)` part.
