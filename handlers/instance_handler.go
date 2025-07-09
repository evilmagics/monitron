package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"monitoring-backend/models"
	"monitoring-backend/services"
)

type InstanceHandler struct {
	service *services.InstanceService
}

func NewInstanceHandler(service *services.InstanceService) *InstanceHandler {
	return &InstanceHandler{service: service}
}

func (h *InstanceHandler) GetInstances(c *gin.Context) {
	instances, err := h.service.GetInstances()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, instances)
}

func (h *InstanceHandler) GetInstance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instance ID"})
		return
	}

	instance, err := h.service.GetInstanceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}

	c.JSON(http.StatusOK, instance)
}

func (h *InstanceHandler) CreateInstance(c *gin.Context) {
	var instance models.Instance
	if err := c.ShouldBindJSON(&instance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateInstance(&instance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, instance)
}

func (h *InstanceHandler) UpdateInstance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instance ID"})
		return
	}

	var instance models.Instance
	if err := c.ShouldBindJSON(&instance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateInstance(uint(id), &instance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instance updated successfully"})
}

func (h *InstanceHandler) DeleteInstance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instance ID"})
		return
	}

	if err := h.service.DeleteInstance(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instance deleted successfully"})
}

func (h *InstanceHandler) GetInstanceMetrics(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instance ID"})
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

	metrics, err := h.service.GetInstanceMetrics(uint(id), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

