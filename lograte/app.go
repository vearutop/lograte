// Package lograte is a CLI app.
package lograte

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"sort"
	"time"

	"github.com/bool64/dev/version"
	"github.com/cespare/xxhash/v2"
	"github.com/vearutop/lograte/filter"
)

// ParseTime is an option to extract line timestamp.
var ParseTime func(line []byte) (time.Time, error)

// Main is the lograte application.
func Main() {
	var (
		buckets         int
		top             int
		length          int
		interval        time.Duration
		ver             bool
		bySize          bool
		noTime          bool
		lineBuf         int
		cpuProfile      string
		parseTimeRegex  string
		parseTimeFormat string
		skipLines       int
		maxLines        int
	)

	flag.IntVar(&buckets, "buckets", 500, "max number of buckets to track filtered messages")
	flag.IntVar(&top, "top", 0, "show top filtered messages ordered by rate")
	flag.IntVar(&length, "len", 120, "limit message length")
	flag.DurationVar(&interval, "t", time.Second, "reporting interval")
	flag.BoolVar(&bySize, "by-size", false, "order messages by size instead of count")
	flag.BoolVar(&ver, "version", false, "print version and exit")
	flag.IntVar(&lineBuf, "line-buf", 1e7, "line token buffer size")
	flag.BoolVar(&noTime, "no-time", false, "do not use time metrics, for non tailing mode")
	flag.IntVar(&skipLines, "skip-lines", 0, "number of lines to skip at the beginning")
	flag.IntVar(&maxLines, "max-lines", 0, "stop after a number of processed lines")
	flag.StringVar(&parseTimeRegex, "parse-time-regex", "", "regex to parse time value from log line")
	flag.StringVar(&parseTimeFormat, "parse-time-format", "", "format to parse time from log line")
	flag.StringVar(&cpuProfile, "dbg-cpu-prof", "", "write first 10 seconds of CPU profile to file")

	flag.Parse()

	if cpuProfile != "" {
		f, err := os.Create(cpuProfile) //nolint:gosec
		if err != nil {
			log.Fatal(err)
		}

		if err = pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}

		go func() {
			time.Sleep(10 * time.Second)
			pprof.StopCPUProfile()
		}()
	}

	if ver {
		fmt.Println(version.Module("github.com/vearutop/lograte").Version)

		return
	}

	d := xxhash.New()

	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, lineBuf)
	scanner.Buffer(buf, len(buf))

	cnt := 0
	byteCnt := 0

	start := time.Now()
	lastReport := start
	lastTS := start

	counts := map[uint64]int{}
	byteCounts := map[uint64]int{}
	samples := map[uint64]string{0: "Other"}
	entries := make([]entry, 0, buckets)

	var timeRegex *regexp.Regexp
	if parseTimeRegex != "" {
		timeRegex = regexp.MustCompile(parseTimeRegex)
		start = time.Time{}
	}

	report := func() {
		lastReport = time.Now()

		if timeRegex == nil && ParseTime == nil {
			lastTS = lastReport
		}

		ela := lastTS.Sub(start).Seconds()
		lps := float64(cnt) / ela
		MB := float64(byteCnt) / (1024 * 1024)
		MBps := MB / ela

		if cnt == 0 {
			if noTime {
				fmt.Printf("0 lines, %.1f MB\n", MB)
			} else {
				fmt.Printf("0 lines,  %.1f MB since %s\n", MB, start.Format(time.Stamp))
			}

			return
		}

		fmt.Println(scanner.Text())

		if noTime {
			fmt.Printf("%d lines, %.1f MB, %d B/avg\n",
				cnt, MB, byteCnt/cnt)
		} else {
			fmt.Printf("%d lines since %s, %.1f per second, %.1f MB/s, %d B/avg\n",
				cnt, start.Format(time.Stamp), lps, MBps, byteCnt/cnt)
		}

		if top <= 0 {
			return
		}

		entries = entries[:0]
		for h, c := range counts {
			entries = append(entries, entry{
				hash:  h,
				cnt:   c,
				bytes: byteCounts[h],
			})
		}

		if bySize {
			sort.Slice(entries, func(i, j int) bool {
				return entries[i].bytes > entries[j].bytes
			})
		} else {
			sort.Slice(entries, func(i, j int) bool {
				return entries[i].cnt > entries[j].cnt
			})
		}

		fmt.Printf("------ Top %d -------------------------\n", top)

		if len(entries) > top {
			entries = entries[0:top]
		}

		for _, e := range entries {
			cntPercent := 100 * float64(e.cnt) / float64(cnt-skipLines)
			bytesPercent := 100 * float64(e.bytes) / float64(byteCnt)
			lps = float64(e.cnt) / ela
			MB = float64(byteCounts[e.hash]) / (1024 * 1024)
			MBps = MB / ela

			if noTime {
				fmt.Printf("%d lines (%.1f%%), %.1f MB (%.1f%%), %d B/avg: %s\n",
					e.cnt, cntPercent, MB, bytesPercent, e.bytes/e.cnt, samples[e.hash])
			} else {
				fmt.Printf("%d lines, %.1f lps (%.1f%%), %.1f MB/s (%.1f%%), %d B/avg: %s\n",
					e.cnt, lps, cntPercent, MBps, bytesPercent, e.bytes/e.cnt, samples[e.hash])
			}
		}

		fmt.Printf("---------------------------------------\n\n")
	}

	for scanner.Scan() {
		line := scanner.Bytes()
		cnt++

		if skipLines > 0 && cnt < skipLines {
			continue
		}

		if maxLines > 0 && cnt >= maxLines {
			break
		}

		byteCnt += len(line)

		if timeRegex != nil {
			matches := timeRegex.FindSubmatch(line)
			if len(matches) == 0 {
				panic("no matches for timestamp regex")
			}

			ts, err := time.Parse(parseTimeFormat, string(matches[1]))
			if err != nil {
				panic(err)
			}

			if start.IsZero() {
				start = ts
			}

			lastTS = ts
		} else if ParseTime != nil {
			ts, err := ParseTime(line)
			if err != nil {
				panic(err)
			}

			if !ts.IsZero() {
				if start.IsZero() {
					start = ts
				}

				lastTS = ts
			}
		}

		if time.Since(lastReport) > interval {
			report()
		}

		if top > 0 {
			filtered := filter.Dynamic(line, length)

			d.Reset()

			_, err := d.Write(filtered)
			if err != nil {
				log.Fatal(err.Error())
			}

			h := d.Sum64()

			if counts[h] == 0 {
				if len(counts) > buckets {
					h = 0
				} else {
					samples[h] = string(filtered)
				}
			}

			counts[h]++

			byteCounts[h] += len(line)
		}
	}

	report()

	if err := scanner.Err(); err != nil {
		fmt.Println("Scan error:", err.Error())
	}
}

type entry struct {
	hash  uint64
	cnt   int
	bytes int
}
