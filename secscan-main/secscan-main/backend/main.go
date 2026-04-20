package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SSRF PROTECTION logic
func validateURL(target string) error {
	u, err := url.Parse(target)
	if err != nil {
		return err
	}

	ips, err := net.LookupIP(u.Hostname())
	if err != nil {
		return err
	}

	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() {
			return fmt.Errorf("SSRF Protection: Private/Internal IP ranges are not allowed")
		}
	}
	return nil
}

type ScanRequest struct {
	URL     string   `json:"url"`
	Modules []string `json:"modules"`
}

type ScanStatus struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Result   string `json:"result,omitempty"`
}

var scans = make(map[string]*ScanStatus)
var channels = make(map[string]chan string)

func main() {
	r := gin.Default()

	// CORS Setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.POST("/scan", startScan)
		api.GET("/scan/:id", getScan)
		api.GET("/scan/:id/stream", streamScanProgress)
	}

	r.Run(":8080")
}

func startScan(c *gin.Context) {
	var req ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// SSRF Check
	if err := validateURL(req.URL); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	scanID := fmt.Sprintf("scan_%d", time.Now().UnixNano())
	
	scans[scanID] = &ScanStatus{
		ID:       scanID,
		Status:   "running",
		Progress: 0,
	}
	channels[scanID] = make(chan string, 10)

	// Async Scanning Routine (Simulated or Real Logic)
	go processScan(scanID, req)

	c.JSON(http.StatusAccepted, gin.H{"scan_id": scanID})
}

func processScan(scanID string, req ScanRequest) {
	status := scans[scanID]
	ch := channels[scanID]

	// Simulate Nmap / ZAP / Header modules
	modulesToRun := len(req.Modules)
	if modulesToRun == 0 {
		modulesToRun = 3
	}

	for i := 0; i < modulesToRun; i++ {
		time.Sleep(2 * time.Second) // Simulate scanning time
		status.Progress = int((float64(i+1) / float64(modulesToRun)) * 100)
		ch <- fmt.Sprintf(`{"progress": %d, "message": "Modul %d tamamlandı"}`, status.Progress, i+1)
	}

	status.Status = "completed"
	status.Result = `{"high_risk": 0, "medium_risk": 2, "low_risk": 5}`
	ch <- `{"progress": 100, "message": "Tarama Başarıyla Tamamlandı", "status": "completed"}`
	
	// Small delay to allow frontend to catch the last message
	time.Sleep(1 * time.Second)
	close(ch)
}

func getScan(c *gin.Context) {
	id := c.Param("id")
	status, exists := scans[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scan not found"})
		return
	}
	c.JSON(http.StatusOK, status)
}

func streamScanProgress(c *gin.Context) {
	id := c.Param("id")
	ch, exists := channels[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scan stream not found"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	c.Stream(func(w int) bool {
		if msg, ok := <-ch; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}
