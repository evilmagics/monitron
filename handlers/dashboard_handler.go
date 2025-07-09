package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"monitoring-backend/services"
)

type DashboardHandler struct {
	instanceService *services.InstanceService
	serviceService  *services.ServiceService
	dnsService      *services.DNSService
}

func NewDashboardHandler(instanceService *services.InstanceService, serviceService *services.ServiceService, dnsService *services.DNSService) *DashboardHandler {
	return &DashboardHandler{
		instanceService: instanceService,
		serviceService:  serviceService,
		dnsService:      dnsService,
	}
}

type DashboardResponse struct {
	Summary struct {
		TotalInstances int `json:"total_instances"`
		TotalServices  int `json:"total_services"`
		TotalDNS       int `json:"total_dns"`
		HealthyServices int `json:"healthy_services"`
		UnhealthyServices int `json:"unhealthy_services"`
		HealthyDNS     int `json:"healthy_dns"`
		UnhealthyDNS   int `json:"unhealthy_dns"`
	} `json:"summary"`
	RecentAlerts []interface{} `json:"recent_alerts"`
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	var response DashboardResponse

	// Get instances count
	instances, err := h.instanceService.GetInstances()
	if err == nil {
		response.Summary.TotalInstances = len(instances)
	}

	// Get services count and health status
	services, err := h.serviceService.GetServices()
	if err == nil {
		response.Summary.TotalServices = len(services)
		
		// Check service health
		for _, service := range services {
			latestCheck, err := h.serviceService.GetLatestCheck(service.ID)
			if err == nil && latestCheck != nil {
				if latestCheck.Status == "up" {
					response.Summary.HealthyServices++
				} else {
					response.Summary.UnhealthyServices++
				}
			} else {
				response.Summary.UnhealthyServices++
			}
		}
	}

	// Get DNS records count and health status
	dnsRecords, err := h.dnsService.GetDNSRecords()
	if err == nil {
		response.Summary.TotalDNS = len(dnsRecords)
		
		// Check DNS health
		for _, record := range dnsRecords {
			latestCheck, err := h.dnsService.GetLatestCheck(record.ID)
			if err == nil && latestCheck != nil {
				if latestCheck.Status == "resolved" {
					response.Summary.HealthyDNS++
				} else {
					response.Summary.UnhealthyDNS++
				}
			} else {
				response.Summary.UnhealthyDNS++
			}
		}
	}

	// TODO: Get recent alerts
	response.RecentAlerts = []interface{}{}

	c.JSON(http.StatusOK, response)
}

