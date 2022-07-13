package main

import (
	"regexp"
	"testing"
)

func TestRegex(t *testing.T) {
	s := `esh-backend-1 i2 2022/07/12 19:23:04.627593 DEBUG#9433 found tracker ID. url tracker: ocdix8t (labels: ["TikTok"], partnerIds: ["1678"]) form labels: ["unknown" "ArcadeHole_drd_T6_new (1737523290876930)" "T6 (1737523291334706)" "ARHO_010722_TT_19_30sec_9x16_TT.mp4_005 (1737523293074466)"], form ids: ["1737523293074466" "1737523290876930" "1737523291334706"]`

	//re := regexp.MustCompile(`([a-zA-Z_-]{1,}\d{1,})+|(\d{1,}[a-zA-Z_-]{1,})+|([\d])+`)
	//re := regexp.MustCompile(`([\w_-]+\d{1,})+|(\d{1,}[\w_-]{1,})+|([\d])+`)
	//re := regexp.MustCompile(`([\w_-]+\d+)|(\d+[\w_-]+)|([\w_-]+\d+[\w_-]+)|([\d])+`)
	re := regexp.MustCompile(`([\S]{1,}\d{1,}[\S]{1,})+|([\S]{1,}\d{1,})+|(\d{1,}[\S]{1,})+|([\d])+`)

	ss := re.ReplaceAllString(s, "X")
	println(ss)
}

func BenchmarkRegex(b *testing.B) {
	s := []byte(`esh-backend-1 i2 2022/07/12 19:23:04.627593 DEBUG#9433 found tracker ID. url tracker: ocdix8t (labels: ["TikTok"], partnerIds: ["1678"]) form labels: ["unknown" "ArcadeHole_drd_T6_new (1737523290876930)" "T6 (1737523291334706)" "ARHO_010722_TT_19_30sec_9x16_TT.mp4_005 (1737523293074466)"], form ids: ["1737523293074466" "1737523290876930" "1737523291334706"]`)

	//re := regexp.MustCompile(`([\w_-]{1,}\d{1,}[\w_-]{1,})+|([\w_-]{1,}\d{1,})+|(\d{1,}[\w_-]{1,})+|([\d])+`)
	re := regexp.MustCompile(`([\w_-]{1,}\d{1,}[\w_-]{1,})+|([\w_-]{1,}\d{1,})+|(\d{1,}[\w_-]{1,})+|([\d])+`)
	//re := regexp.MustCompile(`([\S]{1,}\d{1,}[\S]{1,})+|([\S]{1,}\d{1,})+|(\d{1,}[\S]{1,})+|([\d])+`)

	//re := regexp.MustCompile(`([\w_-]\d)+|(\d[\w_-]+)+|([\d])+`)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = re.ReplaceAll(s, []byte("X"))
	}
}
