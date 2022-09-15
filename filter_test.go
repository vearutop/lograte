package main

import (
	"regexp"
	"testing"
)

const logLine = `foo-bar-16 i2 2022/09/02 09:48:29.199655 baz handler failed to get cox: failed to get cox type: unknown cox type string: {quux} [R method:GET path:/abcd5s8 ra:2022-09-02T09:48:29 form:'map[bar:[11_2022-09-02] lox_id:[] cucumber:[RedCat1509_2022-09-02] cucumber_id:[132072] faux_id:[afExxSDFKHgBJcwxDgIxDETR1vEAWVVHqXo6PcBjfoaDF29f_I8jYTZZVyKeiXzPlP9O9k3SrZtY3IeqA] cox_alarm:[{payout}}] cox_carrot:[OCD] cox_type:[{marks}] creative:[62_203206_123ebd32047fe640] foo_boo_99diks:[https://peebee.jeeass-foo.site/pushforw?lockid=sdd32432dUR1vEAWVVHqXo6PcBjfoaDF29f_Ik3SrZtY3FzXvq0fP1IeqA] goal:[99diks] gps_pos:[1231230D-0BE9-41EF-B146-123123123B9BC7] baz:[123123124-0BE9-41EF-B146-123123123123] ip_address:[1.2.333.4] labelle:[72_206706_125329e2047fe640] poob_id:[62_206706_550eb9e2047fe640] baz_limit:[1234]]' header:'map[User-Agent:[Mozilla/5.0 (iPhone; CPU iPhone OS 14_4_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148] X-Boo-Id:[12312380d-0be9-41ef-b146-3452352323] X-Forwarded-For:[111.332.555.333] X-Forwarded-Proto:[https] X-Forwarded-For:[123.321.123.321]]' foocksorized gees_valeed deeedre_lee:0]`

// BenchmarkFilterAlphanumeric checks bespoke implementation.
// BenchmarkFilterAlphanumeric-12    	 1912300	       614.7 ns/op	       0 B/op	       0 allocs/op.
func BenchmarkFilterAlphanumeric(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		s := []byte(logLine)
		a := filterAlphanumeric(s, 150)
		_ = a
	}
}

func TestAlphanum(t *testing.T) {
	filtered := string(filterAlphanumeric([]byte(logLine), 150))
	expected := `X X X/X/X X:X:X.X baz handler failed to get cox: failed to get cox type: unknown cox type string: {quux} [R method:GET path:/X ra:X:X:X form:'map[bar:`

	if expected != filtered {
		t.Fatalf("unexpected filtered: %s", filtered)
	}
}

func TestShortLine(t *testing.T) {
	line := "foo-bar-12 i3 2022/09/15 11:24:10.689412 baz 0-275 foo bar"
	expected := "X X X/X/X X:X:X.X baz X foo bar"
	filtered := string(filterAlphanumeric([]byte(line), 120))

	if expected != filtered {
		t.Fatalf("unexpected filtered: %s", filtered)
	}
}

// BenchmarkRegex checks regex-based filtering implementation, for reference.
// BenchmarkRegex-12                 	   76851	     15707 ns/op	     200 B/op	       5 allocs/op.
func BenchmarkRegex(b *testing.B) {
	// re := regexp.MustCompile(`([\w_-]{1,}\d{1,}[\w_-]{1,})+|([\w_-]{1,}\d{1,})+|(\d{1,}[\w_-]{1,})+|([\d])+`)
	re := regexp.MustCompile(`([\w_-]{1,}\d{1,}[\w_-]{1,})+|([\w_-]{1,}\d{1,})+|(\d{1,}[\w_-]{1,})+|([\d])+`)
	// re := regexp.MustCompile(`([\S]{1,}\d{1,}[\S]{1,})+|([\S]{1,}\d{1,})+|(\d{1,}[\S]{1,})+|([\d])+`)
	// re := regexp.MustCompile(`([\w_-]\d)+|(\d[\w_-]+)+|([\d])+`)

	s := []byte(logLine[0:150])

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = re.ReplaceAll(s, []byte("X"))
	}
}
