package scanner

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CVEResult Uygulamanın kullandığı teknoloji stack'ini tespit edip CVE (Trivy/NVD) simülasyonu sunar (F07)
type CVEResult struct {
	TechStack string   `json:"tech_stack"`
	Threats   []string `json:"threats"`
}

// AnalyzeCVEs HTTP response header'larından stack tespiti
func AnalyzeCVEs(url string) *CVEResult {
	res := &CVEResult{TechStack: "Bilinmiyor, WAF arkasında olabilir"}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return res
	}
	defer resp.Body.Close()

	poweredBy := resp.Header.Get("X-Powered-By")
	server := resp.Header.Get("Server")

	if strings.Contains(poweredBy, "Express") {
		res.TechStack = "Express.js / Node.js"
		res.Threats = append(res.Threats, "[Express] Trivy Alert: CVE-2022-24999 (High) - Request Smuggling")
	} else if strings.Contains(server, "nginx") {
		res.TechStack = fmt.Sprintf("Nginx Proxy (%s)", server)
		res.Threats = append(res.Threats, "[Nginx] NVD Uyarıları düzenli kontrol edilmeli. Eski sürüm.")
	} else if strings.Contains(poweredBy, "PHP") {
		res.TechStack = "PHP " + poweredBy
		res.Threats = append(res.Threats, "[PHP] Kritik Zafiyet Cihazı Olabilir. Mutlaka sürüm gizlenmelidir.")
	} else {
		res.TechStack = "Bilinmeyen Sunucu Tipi"
		res.Threats = append(res.Threats, "Header'lardan tech stack analiz edilemedi, bilgi isabetli değil.")
	}

	return res
}
