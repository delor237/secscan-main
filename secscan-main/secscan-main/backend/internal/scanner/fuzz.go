package scanner

import (
	"fmt"
	"net/http"
	"time"
)

// FuzzDirectories kelime listesiyle görünmez dizinlere bruteforce yapar (F06-1)
func FuzzDirectories(baseURL string) []string {
	// SecScan Laboratuvarı gereği optimize edilmiş mini wordlist
	wordlist := []string{".git", ".env", "admin", "config", "backup.zip", "test", "API", ".DS_Store"}
	var found []string
	
	client := &http.Client{Timeout: 2 * time.Second}

	for _, p := range wordlist {
		target := baseURL
		if target[len(target)-1] != '/' {
			target += "/"
		}
		target += p

		resp, err := client.Get(target)
		if err == nil {
			// Sadece bulguları raporlar
			if resp.StatusCode == 200 || resp.StatusCode == 403 {
				found = append(found, fmt.Sprintf("/%s (Status: %d)", p, resp.StatusCode))
			}
			resp.Body.Close()
		}
	}
	return found
}
