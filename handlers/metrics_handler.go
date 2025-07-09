package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"monitoring-backend/models"
	"monitoring-backend/services"
)

type MetricsHandler struct {
	instanceService *services.InstanceService
}

func NewMetricsHandler(instanceService *services.InstanceService) *MetricsHandler {
	return &MetricsHandler{instanceService: instanceService}
}

type MetricsRequest struct {
	InstanceID     uint     `json:"instance_id" binding:"required"`
	CPUUsage       *float64 `json:"cpu_usage"`
	MemoryUsage    *float64 `json:"memory_usage"`
	DiskUsage      *float64 `json:"disk_usage"`
	NetworkIOIn    *int64   `json:"network_io_in"`
	NetworkIOOut   *int64   `json:"network_io_out"`
	Uptime         *int64   `json:"uptime"`
	ProcessCount   *int     `json:"process_count"`
	LoadAverage1m  *float64 `json:"load_average_1m"`
	LoadAverage5m  *float64 `json:"load_average_5m"`
	LoadAverage15m *float64 `json:"load_average_15m"`
}

func (h *MetricsHandler) ReceiveMetrics(c *gin.Context) {
	var req MetricsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify instance exists
	_, err := h.instanceService.GetInstanceByID(req.InstanceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}

	// Create metrics record
	metrics := &models.InstanceMetric{
		CPUUsage:       req.CPUUsage,
		MemoryUsage:    req.MemoryUsage,
		DiskUsage:      req.DiskUsage,
		NetworkIOIn:    req.NetworkIOIn,
		NetworkIOOut:   req.NetworkIOOut,
		Uptime:         req.Uptime,
		ProcessCount:   req.ProcessCount,
		LoadAverage1m:  req.LoadAverage1m,
		LoadAverage5m:  req.LoadAverage5m,
		LoadAverage15m: req.LoadAverage15m,
	}

	if err := h.instanceService.SaveMetrics(req.InstanceID, metrics); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Metrics received successfully"})
}

