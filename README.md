# SecScan — OWASP Web Security Scanner
 İsim: NATHANAELLE
 Numera : 24080410150

**Full-stack security tool built for the BMU1208 Final Project.**

## 📋 Project Overview
SecScan is a high-performance web security scanner built with Go and Next.js. It allows users to scan public URLs for common OWASP Top 10 vulnerabilities and provides real-time progress updates via Server-Sent Events (SSE).

## 🛠️ Tech Stack
- **Backend:** Go 1.22, Gin Gonic, SQLite (GORM)
- **Frontend:** Next.js 14, Tailwind CSS, TypeScript
- **DevSecOps:** Docker Compose, Semgrep SAST, Trivy SBOM Scan

## 🛡️ Scanner Features (F01–F10)
- **SQLi Detection:** Identifies potential SQL Injection entry points.
- **XSS Testing:** Scans for reflected and stored XSS vulnerabilities.
- **Security Header Analysis:** Verifies CSP, HSTS, and X-Frame-Options.
- **Parallel Scanning:** Uses Go routines for multi-modular analysis.
- **PDF Export:** Generates downloadable security audit reports.
- **Live Progress:** SSE-based live dashboard with radar charts.

## ⚙️ Setup Instructions

### Backend
`bash
cd backend
go run main.go
`

### Frontend
`bash
cd frontend
npm install
npm run dev
`

## ⚠️ Safety & SSRF Protection
This scanner includes built-in **SSRF (Server-Side Request Forgery) Protection**. It automatically blocks requests to:
- Loopback addresses (127.0.0.1)
- Private IP ranges (10.0.0.0/8, 192.168.0.0/16, etc.)
- Metadata service endpoints.

This ensures the scanner cannot be used to probe your internal network infrastructure.
