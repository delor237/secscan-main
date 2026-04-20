package scanner

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func isInternalIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsPrivate() {
		return true
	}
	// Spesifik AWS Metadata IP kontrolü (Capital One Hack vektörü)
	if ip.String() == "169.254.169.254" {
		return true
	}
	return false
}

// ScanPorts hedefteki portları çok kanallı (goroutine) olarak tarar.
// SSRF (Server-Side Request Forgery) koruması içerir.
func ScanPorts(targetURL string) ([]int, error) {
	target := strings.TrimPrefix(targetURL, "http://")
	target = strings.TrimPrefix(target, "https://")
	if idx := strings.Index(target, "/"); idx != -1 {
		target = target[:idx]
	} // Domain yada IP ayıklanması

	// SSRF DNS Çözümleme Koruma Duvarı
	ips, err := net.LookupIP(target)
	if err != nil {
		return nil, fmt.Errorf("SSRF Guard: Hedef IP'ye çözülemedi (%v)", err)
	}

	for _, ip := range ips {
		if isInternalIP(ip) {
			return nil, fmt.Errorf("SSRF KORUMASI DEVREDE: Özel/İç IP adresine (%s) erişim engellendi", ip.String())
		}
	}

	var openPorts []int
	var mu sync.Mutex
	var wg sync.WaitGroup

	// En tipik 15 port ile hızlı concurrency testi
	commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 445, 3306, 3389, 5432, 8080, 8443}

	for _, port := range commonPorts {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", target, p)
			conn, err := net.DialTimeout("tcp", address, 1*time.Second)
			if err == nil {
				mu.Lock()
				openPorts = append(openPorts, p)
				mu.Unlock()
				conn.Close()
			}
		}(port)
	}

	wg.Wait()
	return openPorts, nil
}
