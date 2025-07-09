package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"monitoring-backend/models"
	"monitoring-backend/services"
)

type DNSHandler struct {
	service *services.DNSService
}

func NewDNSHandler(service *services.DNSService) *DNSHandler {
	return &DNSHandler{service: service}
}

func (h *DNSHandler) GetDNSRecords(c *gin.Context) {
	records, err := h.service.GetDNSRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (h *DNSHandler) GetDNSRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DNS record ID"})
		return
	}

	record, err := h.service.GetDNSRecordByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DNS record not found"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (h *DNSHandler) CreateDNSRecord(c *gin.Context) {
	var record models.DNSRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateDNSRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

func (h *DNSHandler) UpdateDNSRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DNS record ID"})
		return
	}

	var record models.DNSRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateDNSRecord(uint(id), &record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DNS record updated successfully"})
}

func (h *DNSHandler) DeleteDNSRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DNS record ID"})
		return
	}

	if err := h.service.DeleteDNSRecord(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DNS record deleted successfully"})
}

func (h *DNSHandler) GetDNSChecks(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DNS record ID"})
		return
	}

	// Parse time range parameters
	var from, to time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from time format"})
			return
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to time format"})
			return
		}
	}

	checks, err := h.service.GetDNSChecks(uint(id), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, checks)
}

