package scanner

import (
	"net/http"
	"time"
)

// TestSQLi Error-based veya Boolean-based minik SQL Injection sinyali arar (F06-3)
func TestSQLi(baseURL string) bool {
	payload := "' OR 1=1 --"
	target := baseURL + "?id=" + payload

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(target)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Laboratuvar iskeleti gereği çok basit bir Warning/500 analizi:
	// Eğer tek tırnak ve payload eklenince sistem Internal Error veriyorsa muhtemel zafiyettir.
	if resp.StatusCode == http.StatusInternalServerError {
		return true // SQL Syntax veya DB Hatası düşmüş olabilir
	}
	return false
}
