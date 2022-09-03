# lograte

A small tool ([used to be](https://github.com/vearutop/lograte/blob/v0.1.0/main.go) ~25 lines of code) to calculate the 
rate of lines in STDOUT and group by them by count with alphanumeric filter.

## Installation

```
go install github.com/vearutop/lograte@latest
```

or download prebuilt binary from [releases](https://github.com/vearutop/lograte/releases).

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

```
tail -f /var/log/quick.log | lograte -top 5 -t 1s -by-size
```
```
...
foo-bar-18 i2 2022/09/03 08:22:17.206522 <recent log entry>
271382 lines since Sep  3 08:22:14, 90447.4 per second, 68.3 MB/s, 791 B/avg
------ Top 5 -------------------------
29810 lines, 9935.2 lps (11.0%), 14.7 MB/s (21.5%), 1546 B/avg: X X X/X/X X:X:X.X filtered entry 
26559 lines, 8851.7 lps (9.8%), 9.1 MB/s (13.4%), 1083 B/avg: X X X/X/X X:X:X.X another filtered entry
...
---------------------------------------
```