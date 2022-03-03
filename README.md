# lograte

A small tool (~25 lines of code) to calculate the rate of lines in STDOUT.

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