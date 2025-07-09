package services

import (
	"fmt"
	"monitoring-backend/models"
	"net"
	"net/http"
	"time"

	"github.com/miekg/dns"
	"gorm.io/gorm"
)

type MonitoringService struct {
	db                  *gorm.DB
	serviceService      *ServiceService
	dnsService          *DNSService
	notificationService *NotificationService
}

func NewMonitoringService(db *gorm.DB, serviceService *ServiceService, dnsService *DNSService, notificationService *NotificationService) *MonitoringService {
	return &MonitoringService{
		db:                  db,
		serviceService:      serviceService,
		dnsService:          dnsService,
		notificationService: notificationService,
	}
}

func (s *MonitoringService) StartScheduler() {
	// Start service monitoring
	go s.serviceMonitoringLoop()
	
	// Start DNS monitoring
	go s.dnsMonitoringLoop()
}

func (s *MonitoringService) serviceMonitoringLoop() {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAllServices()
		}
	}
}

func (s *MonitoringService) dnsMonitoringLoop() {
	ticker := time.NewTicker(60 * time.Second) // Check every 60 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAllDNSRecords()
		}
	}
}

func (s *MonitoringService) checkAllServices() {
	services, err := s.serviceService.GetServices()
	if err != nil {
		fmt.Printf("Error getting services: %v\n", err)
		return
	}

	for _, service := range services {
		// Check if it's time to monitor this service
		lastCheck, _ := s.serviceService.GetLatestCheck(service.ID)
		if lastCheck != nil && time.Since(lastCheck.Timestamp) < time.Duration(service.IntervalSeconds)*time.Second {
			continue
		}

		go s.checkService(service)
	}
}

func (s *MonitoringService) checkService(service models.Service) {
	var check models.ServiceCheck
	start := time.Now()

	switch service.Type {
	case models.ServiceTypeHTTP:
		check = s.checkHTTPService(service)
	case models.ServiceTypeTCP:
		check = s.checkTCPService(service)
	case models.ServiceTypeGRPC:
		check = s.checkGRPCService(service)
	case models.ServiceTypeMQTT:
		check = s.checkMQTTService(service)
	default:
		check = models.ServiceCheck{
			Status:       models.ServiceCheckStatusUnknown,
			ErrorMessage: "Unknown service type",
		}
	}

	responseTime := int(time.Since(start).Milliseconds())
	check.ResponseTimeMs = &responseTime

	// Save check result
	if err := s.serviceService.SaveCheck(service.ID, &check); err != nil {
		fmt.Printf("Error saving service check: %v\n", err)
	}

	// Send alert if service is down
	if check.Status == models.ServiceCheckStatusDown {
		alert := &models.Alert{
			SourceType:   models.AlertSourceTypeService,
			SourceID:     service.ID,
			MetricName:   "service_status",
			Threshold:    "up",
			CurrentValue: string(check.Status),
			Status:       models.AlertStatusTriggered,
			TriggeredAt:  time.Now(),
			Message:      fmt.Sprintf("Service %s is down: %s", service.Name, check.ErrorMessage),
		}
		s.notificationService.SendAlert(alert)
	}
}

func (s *MonitoringService) checkHTTPService(service models.Service) models.ServiceCheck {
	url := fmt.Sprintf("http://%s:%d", service.Host, service.Port)
	if service.HTTPConfig != nil {
		url += service.HTTPConfig.Path
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return models.ServiceCheck{
			Status:       models.ServiceCheckStatusDown,
			ErrorMessage: err.Error(),
		}
	}
	defer resp.Body.Close()

	expectedStatus := 200
	if service.HTTPConfig != nil && service.HTTPConfig.ExpectedStatus != 0 {
		expectedStatus = service.HTTPConfig.ExpectedStatus
	}

	if resp.StatusCode != expectedStatus {
		return models.ServiceCheck{
			Status:       models.ServiceCheckStatusDown,
			ErrorMessage: fmt.Sprintf("Expected status %d, got %d", expectedStatus, resp.StatusCode),
		}
	}

	return models.ServiceCheck{
		Status: models.ServiceCheckStatusUp,
	}
}

func (s *MonitoringService) checkTCPService(service models.Service) models.ServiceCheck {
	timeout := 5 * time.Second
	if service.TCPConfig != nil && service.TCPConfig.TimeoutMs > 0 {
		timeout = time.Duration(service.TCPConfig.TimeoutMs) * time.Millisecond
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", service.Host, service.Port), timeout)
	if err != nil {
		return models.ServiceCheck{
			Status:       models.ServiceCheckStatusDown,
			ErrorMessage: err.Error(),
		}
	}
	defer conn.Close()

	return models.ServiceCheck{
		Status: models.ServiceCheckStatusUp,
	}
}

func (s *MonitoringService) checkGRPCService(service models.Service) models.ServiceCheck {
	// TODO: Implement gRPC health check
	// For now, just check TCP connection
	return s.checkTCPService(service)
}

func (s *MonitoringService) checkMQTTService(service models.Service) models.ServiceCheck {
	// TODO: Implement MQTT connection check
	// For now, just check TCP connection
	return s.checkTCPService(service)
}

func (s *MonitoringService) checkAllDNSRecords() {
	records, err := s.dnsService.GetDNSRecords()
	if err != nil {
		fmt.Printf("Error getting DNS records: %v\n", err)
		return
	}

	for _, record := range records {
		// Check if it's time to monitor this DNS record
		lastCheck, _ := s.dnsService.GetLatestCheck(record.ID)
		if lastCheck != nil && time.Since(lastCheck.Timestamp) < time.Duration(record.IntervalSeconds)*time.Second {
			continue
		}

		go s.checkDNSRecord(record)
	}
}

func (s *MonitoringService) checkDNSRecord(record models.DNSRecord) {
	start := time.Now()
	
	c := dns.Client{
		Timeout: 5 * time.Second,
	}

	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(record.Target), dns.TypeA)

	server := "8.8.8.8:53"
	if record.DNSServer != "" {
		server = record.DNSServer + ":53"
	}

	r, _, err := c.Exchange(&m, server)
	
	resolutionTime := int(time.Since(start).Milliseconds())
	
	var check models.DNSCheck
	check.ResolutionTimeMs = &resolutionTime

	if err != nil {
		check.Status = models.DNSCheckStatusFailed
		check.ErrorMessage = err.Error()
	} else if len(r.Answer) == 0 {
		check.Status = models.DNSCheckStatusFailed
		check.ErrorMessage = "No DNS answer received"
	} else {
		check.Status = models.DNSCheckStatusResolved
		if aRecord, ok := r.Answer[0].(*dns.A); ok {
			check.ResolvedIP = aRecord.A.String()
		}
	}

	// Save check result
	if err := s.dnsService.SaveCheck(record.ID, &check); err != nil {
		fmt.Printf("Error saving DNS check: %v\n", err)
	}

	// Send alert if DNS resolution failed
	if check.Status == models.DNSCheckStatusFailed {
		alert := &models.Alert{
			SourceType:   models.AlertSourceTypeDNS,
			SourceID:     record.ID,
			MetricName:   "dns_resolution",
			Threshold:    "resolved",
			CurrentValue: string(check.Status),
			Status:       models.AlertStatusTriggered,
			TriggeredAt:  time.Now(),
			Message:      fmt.Sprintf("DNS resolution failed for %s: %s", record.Target, check.ErrorMessage),
		}
		s.notificationService.SendAlert(alert)
	}
}

