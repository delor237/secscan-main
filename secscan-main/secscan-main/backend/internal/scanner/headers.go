package scanner

import (
	"net/http"
	"time"
)

type HeaderResult struct {
	Present []string `json:"present"`
	Missing []string `json:"missing"`
	Score   string   `json:"score"`
}

// AnalyzeHeaders hedeflenmiş URL üzerindeki HTTP Response Header'larına bakıp puan verir.
func AnalyzeHeaders(url string) (*HeaderResult, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	securityHeaders := []string{
		"Strict-Transport-Security",
		"Content-Security-Policy",
		"X-Frame-Options",
		"X-Content-Type-Options",
		"X-XSS-Protection",
	}

	res := &HeaderResult{}
	foundCount := 0

	for _, h := range securityHeaders {
		if resp.Header.Get(h) != "" {
			res.Present = append(res.Present, h)
			foundCount++
		} else {
			res.Missing = append(res.Missing, h)
		}
	}

	// Puanlama Motoru
	if foundCount == len(securityHeaders) {
		res.Score = "A+"
	} else if foundCount >= 3 {
		res.Score = "B"
	} else if foundCount >= 1 {
		res.Score = "D"
	} else {
		res.Score = "F"
	}

	return res, nil
}
