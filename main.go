package main

import (
	"bufio"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"os"
	"regexp"
	"sort"
	"time"
)

func main() {
	re := regexp.MustCompile(`([\w_-]{1,}\d{1,}[\w_-]{1,})+|([\w_-]{1,}\d{1,})+|(\d{1,}[\w_-]{1,})+|([\d])+`)
	d := xxhash.New()

	scanner := bufio.NewScanner(os.Stdin)
	cnt := 0
	start := time.Now()
	lastReport := start

	counts := map[uint64]uint64{}
	samples := map[uint64]string{0: "Other"}
	entries := make([]entry, 0, buckets)

	for scanner.Scan() {
		cnt++

		line := scanner.Bytes()
		line = head(line, 100)
		filtered := re.ReplaceAll(line, []byte("X"))
		d.Reset()
		_, _ = d.Write(filtered)
		h := d.Sum64()

		if counts[h] == 0 {
			if len(counts) > 100 {
				h = 0
			} else {
				//samples[h] = string(head(scanner.Bytes(), 200))
				samples[h] = string(head(filtered, 200))
			}
		}
		counts[h]++

		if time.Since(lastReport) > time.Second {
			lastReport = time.Now()

			entries = entries[:0]
			for h, c := range counts {
				entries = append(entries, entry{
					hash: h,
					cnt:  c,
				})
			}

			sort.Slice(entries, func(i, j int) bool {
				return entries[i].cnt > entries[j].cnt
			})

			fmt.Println(scanner.Text())
			fmt.Printf("%d lines since %s, %.1f per second\n", cnt, start.Format(time.Stamp), float64(cnt)/time.Since(start).Seconds())

			fmt.Println("Top 10 -------------------------")
			for _, e := range entries[0:10] {
				fmt.Printf("%d lines, %.1f lps: %s\n", e.cnt, float64(e.cnt)/time.Since(start).Seconds(), samples[e.hash])
			}
			fmt.Println("-------------------------")
		}
	}
}

// update_worker -> update_consumer.go -> queryValues method

func head(d []byte, l int) []byte {
	if len(d) < l {
		return d
	}

	return d[0:l]
}

const buckets = 1000

type entry struct {
	hash  uint64
	cnt   uint64
	bytes uint64
}
