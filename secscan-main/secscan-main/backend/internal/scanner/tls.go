package scanner

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

// TLSResult sertifika ve yapılandırma durumu
type TLSResult struct {
	Version            string   `json:"version"`
	CipherSuite        string   `json:"cipher_suite"`
	ValidCert          bool     `json:"valid_cert"`
	Issuer             string   `json:"issuer"`
	ExpirationDaysLeft int      `json:"expiration_days_left"`
	Vulnerabilities    []string `json:"vulnerabilities"`
}

// AnalyzeTLS TLS 1.0-1.3 destek ve sertifika kontrolünü gerçekleştirir (F05)
func AnalyzeTLS(targetURL string) (*TLSResult, error) {
	url := strings.TrimPrefix(targetURL, "https://")
	url = strings.TrimPrefix(url, "http://")
	if idx := strings.Index(url, "/"); idx != -1 {
		url = url[:idx]
	}
	if !strings.Contains(url, ":") {
		url = url + ":443"
	}

	conf := &tls.Config{
		InsecureSkipVerify: true, // Zafiyet taramak için engeli kaldırıyoruz
	}

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 3 * time.Second}, "tcp", url, conf)
	if err != nil {
		return nil, fmt.Errorf("TLS destegi yok veya erisim kapali")
	}
	defer conn.Close()

	state := conn.ConnectionState()
	res := &TLSResult{}

	switch state.Version {
	case tls.VersionTLS10:
		res.Version = "TLS 1.0"
		res.Vulnerabilities = append(res.Vulnerabilities, "Kritik: Çok eski TLS 1.0 kullanımı")
	case tls.VersionTLS11:
		res.Version = "TLS 1.1"
		res.Vulnerabilities = append(res.Vulnerabilities, "Riskli: Eski TLS 1.1 kullanımı")
	case tls.VersionTLS12:
		res.Version = "TLS 1.2"
	case tls.VersionTLS13:
		res.Version = "TLS 1.3"
	default:
		res.Version = "Bilinmeyen"
	}

	res.CipherSuite = tls.CipherSuiteName(state.CipherSuite)

	if len(state.PeerCertificates) > 0 {
		cert := state.PeerCertificates[0]
		if len(cert.Issuer.Organization) > 0 {
			res.Issuer = cert.Issuer.Organization[0]
		}
		days := int(time.Until(cert.NotAfter).Hours() / 24)
		res.ExpirationDaysLeft = days
		if days < 0 {
			res.Vulnerabilities = append(res.Vulnerabilities, "Sertifika Süresi Dolmuş (Expired)")
			res.ValidCert = false
		} else {
			res.ValidCert = true
		}
	}

	return res, nil
}
