package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"time"

	"github.com/bool64/dev/version"
	"github.com/cespare/xxhash/v2"
)

func main() {
	var (
		buckets    int
		top        int
		length     int
		interval   time.Duration
		ver        bool
		bySize     bool
		lineBuf    int
		cpuProfile string
	)

	flag.IntVar(&buckets, "buckets", 500, "max number of buckets to track filtered messages")
	flag.IntVar(&top, "top", 0, "show top filtered messages ordered by rate")
	flag.IntVar(&length, "len", 120, "limit message length")
	flag.DurationVar(&interval, "t", time.Second, "reporting interval")
	flag.BoolVar(&bySize, "by-size", false, "order messages by size instead of count")
	flag.BoolVar(&ver, "version", false, "print version and exit")
	flag.IntVar(&lineBuf, "line-buf", 1e7, "line token buffer size")
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
		fmt.Println(version.Info().Version)

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

	counts := map[uint64]int{}
	byteCounts := map[uint64]int{}
	samples := map[uint64]string{0: "Other"}
	entries := make([]entry, 0, buckets)

	for scanner.Scan() {
		line := scanner.Bytes()
		cnt++

		byteCnt += len(line)

		if time.Since(lastReport) > interval {
			lastReport = time.Now()
			ela := time.Since(start).Seconds()
			lps := float64(cnt) / ela
			MBps := float64(byteCnt) / (ela * 1024 * 1024)

			fmt.Println(scanner.Text())
			fmt.Printf("%d lines since %s, %.1f per second, %.1f MB/s, %d B/avg\n",
				cnt, start.Format(time.Stamp), lps, MBps, byteCnt/cnt)

			if top <= 0 {
				continue
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
				cntPercent := 100 * float64(e.cnt) / float64(cnt)
				bytesPercent := 100 * float64(e.bytes) / float64(byteCnt)
				lps = float64(e.cnt) / ela
				MBps = float64(byteCounts[e.hash]) / (ela * 1024 * 1024)
				fmt.Printf("%d lines, %.1f lps (%.1f%%), %.1f MB/s (%.1f%%), %d B/avg: %s\n",
					e.cnt, lps, cntPercent, MBps, bytesPercent, e.bytes/e.cnt, samples[e.hash])
			}

			fmt.Printf("---------------------------------------\n\n")
		}

		if top > 0 {
			filtered := filterDynamic(line, length)

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

	if err := scanner.Err(); err != nil {
		fmt.Println("Scan error:", err.Error())
	}
}

type entry struct {
	hash  uint64
	cnt   int
	bytes int
}
