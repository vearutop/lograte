package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cnt := 0
	start := time.Now()
	lastReport := start

	for scanner.Scan() {
		cnt++

		if time.Since(lastReport) > time.Second {
			lastReport = time.Now()
			fmt.Println(scanner.Text())
			fmt.Printf("%d lines since %s, %.2f per second\n", cnt, start.Format(time.Stamp), float64(cnt)/time.Since(start).Seconds())
		}
	}
}
