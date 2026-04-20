package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ScanRequest Frontend'den gelecek payload yapısı
type ScanRequest struct {
	URL     string   `json:"url" binding:"required,url"`
	Modules []string `json:"modules" binding:"required"`
}

// HealthCheck sistem durumunu döner
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// StartScan asenkron taramayı başlatır
func StartScan(c *gin.Context) {
	var req ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Asenkron işlem kuyruğuna atılacak (goroutine/redis)
	scanID := "sc_demo123"
	
	c.JSON(http.StatusAccepted, gin.H{
		"scan_id": scanID,
		"message": "Scan started asynchronously",
	})
}

// GetScanResult senkron rapor çekimi
func GetScanResult(c *gin.Context) {
	id := c.Param("id")
	// TODO: Tam Rapor Dönecek (F02 skeleton version)
	c.JSON(http.StatusOK, gin.H{
		"id": id,
		"status": "completed",
		"score": "A-",
		"results": gin.H{},
	})
}

// StreamScanProgress SSE modülü
func StreamScanProgress(c *gin.Context) {
	// TODO: SSE entegrasyonu (F08)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
}

// DownloadPDF binary dönüştürücü
func DownloadPDF(c *gin.Context) {
	// TODO: PDF dönüştürücü entegre edilecek (F09)
	c.String(http.StatusOK, "PDF Blob Data")
}
