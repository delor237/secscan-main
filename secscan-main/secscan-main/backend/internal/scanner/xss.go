package scanner

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// TestXSSReflection URL'ye Reflected XSS enjeksiyonu atıp DOM'u inceler (F06-2)
func TestXSSReflection(baseURL string) bool {
	payload := `'"<script>alert(SECScan_XSS)</script>`
	target := baseURL + "?search=" + payload

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(target)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	// Payload html escape edilmeden aynen sunucudan dönmüşse XSS zafiyeti vardır.
	return strings.Contains(string(body), payload)
}
