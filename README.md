# lograte

A small tool ([used to be](https://github.com/vearutop/lograte/blob/v0.1.0/main.go) ~25 lines of code) to calculate the 
rate of lines in STDOUT and group them by count with alphanumeric filter.

## Installation

```
go install github.com/vearutop/lograte@latest
```

or download prebuilt binary from [releases](https://github.com/vearutop/lograte/releases).

```
wget https://github.com/vearutop/lograte/releases/latest/download/linux_amd64.tar.gz && tar xf linux_amd64.tar.gz && rm linux_amd64.tar.gz
./lograte -version
```

## Usage

Pipe the verbose output (for example tail of logs) to `lograte`
```
tail -f /var/log/nginx/error.log | lograte
```

Once a second, the rate of lines is printed together with the last line in that second.

```
2022/03/03 23:44:17 [error] 30064#30064: *27785232 open() "/home/ubuntu/tbex.ru/w/0!8816!www.sibmail.com!c.js" failed (2: No such file or directory), client: 212.107.254.94, server: tbex.ru, request: "GET /w/0!8816!www.sibmail.com!c.js?rev=4-1646347458270 HTTP/1.1", host: "c.tbex.ru", referrer: "http://www.sibmail.com/"
18 lines since Mar  3 23:43:43, 0.52 per second
```

You can also check the rate of remote logs filtered with grep.

```
ssh -C log-collector.acme.com 'tail -f /var/log/app-error.log | grep "my feature" | grep -v "well known error"' | lograte
```

```
...
2022/03/03 22:46:42.368373 EVENT failed to pass with ID foo
6218 lines since Mar  3 23:46:32, 604.20 per second
2022/03/03 22:46:43.933489 EVENT failed to pass with ID bar
6862 lines since Mar  3 23:46:32, 605.47 per second
```

Or show top filtered messages by count or total size. 

> **_NOTE:_** Filtered messages have all case insensitive sequences of `[a-z]-_%` with at least one digit or all digits replaced with `X`. 
> This is usually enough to remove dynamic data from message and decrease cardinality.
> Filtered messages are collected in a limited number of buckets, once the limit is met all other messages are collected into the `Other` bucket.

```
tail -f /var/log/quick.log | lograte -top 5 -t 1s -by-size
```
```
...
foo-bar-18 i2 2022/09/03 08:22:17.206522 <recent log entry>
271382 lines since Sep  3 08:22:14, 90447.4 per second, 68.3 MB/s, 791 B/avg
------ Top 5 -------------------------
29810 lines, 9935.2 lps (11.0%), 14.7 MB/s (21.5%), 1546 B/avg: X X X/X/X X:X:X.X filtered entry X.X
26559 lines, 8851.7 lps (9.8%), 9.1 MB/s (13.4%), 1083 B/avg: X X X/X/X X:X:X.X another filtered X entry
...
---------------------------------------
```

If you want to analyze existing log files, you can use `lograte` with `cat` or something similar to read them.
For such case you can use `-parse-time-format` and `-parse-time-regex` to parse time from log line instead of current clock.


For example such command would read all `*.zst` files in current directory, filter them with `zstdgrep` and then analyze
using `time.RFC3339Nano` format for the first value between space and tab as a timestamp.
```
zstdgrep "fancy error" *.zst | ~/lograte -top 100 -buckets 1000 -parse-time-regex " ([\d-T:.Z]+)\t" -parse-time-format "2006-01-02T15:04:05.999999999Z07:00"
```

or if you want to analyze lines/bytes distribution without time, use `-no-time`
```
zstdgrep "fancy error" *.zst | ~/lograte -top 100 -buckets 1000 -no-time
```


### Flags

```
lograte -help
```
```
Usage of lograte:
  -buckets int
        max number of buckets to track filtered messages (default 500)
  -by-size
        order messages by size instead of count
  -dbg-cpu-prof string
        write first 10 seconds of CPU profile to file
  -len int
        limit message length (default 120)
  -line-buf int
        line token buffer size (default 10000000)
  -max-lines int
        stop after a number of processed lines
  -no-time
        do not use time metrics, for non tailing mode
  -parse-time-format string
        format to parse time from log line
  -parse-time-regex string
        regex to parse time value from log line
  -skip-lines int
        number of lines to skip at the beginning
  -t duration
        reporting interval (default 1s)
  -top int
        show top filtered messages ordered by rate
  -version
        print version and exit
```
