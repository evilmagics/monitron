package graphql

import (
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"

	"monitron-server/models"
)

// InstanceType defines the GraphQL type for an Instance
var InstanceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Instance",
	Fields: graphql.Fields{
		"id":             &graphql.Field{Type: graphql.ID},
		"name":           &graphql.Field{Type: graphql.String},
		"host":           &graphql.Field{Type: graphql.String},
		"check_interval": &graphql.Field{Type: graphql.Int},
		"check_timeout":  &graphql.Field{Type: graphql.Int},
		"agent_port":     &graphql.Field{Type: graphql.Int},
		"agent_auth":     &graphql.Field{Type: graphql.String},
		"description":    &graphql.Field{Type: graphql.String},
		"label":          &graphql.Field{Type: graphql.String},
		"group":          &graphql.Field{Type: graphql.String},
		"created_at":     &graphql.Field{Type: graphql.DateTime},
		"updated_at":     &graphql.Field{Type: graphql.DateTime},
	},
})

// ServiceType defines the GraphQL type for a Service
var ServiceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Service",
	Fields: graphql.Fields{
		"id":                   &graphql.Field{Type: graphql.ID},
		"name":                 &graphql.Field{Type: graphql.String},
		"api_type":             &graphql.Field{Type: graphql.String},
		"check_interval":       &graphql.Field{Type: graphql.Int},
		"timeout":              &graphql.Field{Type: graphql.Int},
		"description":          &graphql.Field{Type: graphql.String},
		"label":                &graphql.Field{Type: graphql.String},
		"group":                &graphql.Field{Type: graphql.String},
		"http_method":          &graphql.Field{Type: graphql.String},
		"http_health_url":      &graphql.Field{Type: graphql.String},
		"http_expected_status": &graphql.Field{Type: graphql.Int},
		"grpc_host":            &graphql.Field{Type: graphql.String},
		"grpc_port":            &graphql.Field{Type: graphql.Int},
		"grpc_auth":            &graphql.Field{Type: graphql.String},
		"grpc_proto":           &graphql.Field{Type: graphql.String},
		"mqtt_host":            &graphql.Field{Type: graphql.String},
		"mqtt_port":            &graphql.Field{Type: graphql.Int},
		"mqtt_qos":             &graphql.Field{Type: graphql.Int},
		"mqtt_topic":           &graphql.Field{Type: graphql.String},
		"mqtt_auth":            &graphql.Field{Type: graphql.String},
		"tcp_host":             &graphql.Field{Type: graphql.String},
		"tcp_port":             &graphql.Field{Type: graphql.Int},
		"dns_domain_name":      &graphql.Field{Type: graphql.String},
		"ping_host":            &graphql.Field{Type: graphql.String},
		"created_at":           &graphql.Field{Type: graphql.DateTime},
		"updated_at":           &graphql.Field{Type: graphql.DateTime},
	},
})

// DomainSSLType defines the GraphQL type for a DomainSSL
var DomainSSLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DomainSSL",
	Fields: graphql.Fields{
		"id":                 &graphql.Field{Type: graphql.ID},
		"domain":             &graphql.Field{Type: graphql.String},
		"warning_threshold":  &graphql.Field{Type: graphql.Int},
		"expiry_threshold":   &graphql.Field{Type: graphql.Int},
		"check_interval":     &graphql.Field{Type: graphql.Int},
		"label":              &graphql.Field{Type: graphql.String},
		"certificate_detail": &graphql.Field{Type: graphql.String},
		"issuer":             &graphql.Field{Type: graphql.String},
		"valid_from":         &graphql.Field{Type: graphql.DateTime},
		"resolved_ip":        &graphql.Field{Type: graphql.String},
		"expiry":             &graphql.Field{Type: graphql.DateTime},
		"created_at":         &graphql.Field{Type: graphql.DateTime},
		"updated_at":         &graphql.Field{Type: graphql.DateTime},
	},
})

// UserType defines the GraphQL type for a User
var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.ID},
		"username":   &graphql.Field{Type: graphql.String},
		"email":      &graphql.Field{Type: graphql.String},
		"role":       &graphql.Field{Type: graphql.String},
		"status":     &graphql.Field{Type: graphql.String},
		"last_login": &graphql.Field{Type: graphql.DateTime},
		"created_at": &graphql.Field{Type: graphql.DateTime},
		"updated_at": &graphql.Field{Type: graphql.DateTime},
	},
})

// ReportType defines the GraphQL type for a Report
var ReportType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Report",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.ID},
		"name":         &graphql.Field{Type: graphql.String},
		"report_type":  &graphql.Field{Type: graphql.String},
		"format":       &graphql.Field{Type: graphql.String},
		"status":       &graphql.Field{Type: graphql.String},
		"generated_at": &graphql.Field{Type: graphql.DateTime},
		"file_path":    &graphql.Field{Type: graphql.String},
		"user_id":      &graphql.Field{Type: graphql.ID},
		"created_at":   &graphql.Field{Type: graphql.DateTime},
		"updated_at":   &graphql.Field{Type: graphql.DateTime},
	},
})

// LogEntryType defines the GraphQL type for a LogEntry
var LogEntryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LogEntry",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.ID},
		"level":      &graphql.Field{Type: graphql.String},
		"message":    &graphql.Field{Type: graphql.String},
		"timestamp":  &graphql.Field{Type: graphql.DateTime},
		"service":    &graphql.Field{Type: graphql.String},
		"request_id": &graphql.Field{Type: graphql.String},
	},
})

// OperationalPageType defines the GraphQL type for an OperationalPage
var OperationalPageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "OperationalPage",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.ID},
		"slug":        &graphql.Field{Type: graphql.String},
		"name":        &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"is_public":   &graphql.Field{Type: graphql.Boolean},
		"created_at":  &graphql.Field{Type: graphql.DateTime},
		"updated_at":  &graphql.Field{Type: graphql.DateTime},
	},
})

// OperationalPageComponentType defines the GraphQL type for an OperationalPageComponent
var OperationalPageComponentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "OperationalPageComponent",
	Fields: graphql.Fields{
		"id":             &graphql.Field{Type: graphql.ID},
		"page_id":        &graphql.Field{Type: graphql.ID},
		"component_type": &graphql.Field{Type: graphql.String},
		"component_id":   &graphql.Field{Type: graphql.ID},
		"component_name": &graphql.Field{Type: graphql.String},
		"display_order":  &graphql.Field{Type: graphql.Int},
		"description":    &graphql.Field{Type: graphql.String},
		"created_at":     &graphql.Field{Type: graphql.DateTime},
		"updated_at":     &graphql.Field{Type: graphql.DateTime},
	},
})

// RootQuery defines the root query for GraphQL
func RootQuery(db *gorm.DB) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"instances": &graphql.Field{
				Type:        graphql.NewList(InstanceType),
				Description: "Get all instances",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var instances []models.Instance
					err := db.Find(&instances).Error
					if err != nil {
						log.Printf("Error fetching instances for GraphQL: %v", err)
						return nil, err
					}
					return instances, nil
				},
			},
			"instance": &graphql.Field{
				Type:        InstanceType,
				Description: "Get a single instance by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if !ok {
						return nil, fmt.Errorf("invalid instance ID")
					}
					var instance models.Instance
					err := db.Find(&instance, id).Error
					if err != nil {
						log.Printf("Error fetching instance for GraphQL: %v", err)
						return nil, err
					}
					return instance, nil
				},
			},
			"services": &graphql.Field{
				Type:        graphql.NewList(ServiceType),
				Description: "Get all services",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var services []models.Service
					err := db.Find(&services).Error
					if err != nil {
						log.Printf("Error fetching services for GraphQL: %v", err)
						return nil, err
					}
					return services, nil
				},
			},
			"service": &graphql.Field{
				Type:        ServiceType,
				Description: "Get a single service by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if !ok {
						return nil, fmt.Errorf("invalid service ID")
					}
					var service models.Service
					err := db.Find(&service, id).Error
					if err != nil {
						log.Printf("Error fetching service for GraphQL: %v", err)
						return nil, err
					}
					return service, nil
				},
			},
			"domainSSLs": &graphql.Field{
				Type:        graphql.NewList(DomainSSLType),
				Description: "Get all domain/SSL entries",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var domainSSLs []models.DomainSSL
					err := db.Find(&domainSSLs).Error
					if err != nil {
						log.Printf("Error fetching domain/SSLs for GraphQL: %v", err)
						return nil, err
					}
					return domainSSLs, nil
				},
			},
			"domainSSL": &graphql.Field{
				Type:        DomainSSLType,
				Description: "Get a single domain/SSL entry by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if !ok {
						return nil, fmt.Errorf("invalid domain/SSL ID")
					}
					var domainSSL models.DomainSSL
					err := db.Find(&domainSSL, id).Error
					if err != nil {
						log.Printf("Error fetching domain/SSL for GraphQL: %v", err)
						return nil, err
					}
					return domainSSL, nil
				},
			},
			"users": &graphql.Field{
				Type:        graphql.NewList(UserType),
				Description: "Get all users",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var users []models.User
					err := db.Select("id", "username", "role", "status", "last_login", "created_at", "updated_at").Find(&users).Error
					if err != nil {
						log.Printf("Error fetching users for GraphQL: %v", err)
						return nil, err
					}
					return users, nil
				},
			},
			"user": &graphql.Field{
				Type:        UserType,
				Description: "Get a single user by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if !ok {
						return nil, fmt.Errorf("invalid user ID")
					}
					var user models.User
					err := db.Select("id", "username", "role", "status", "last_login", "created_at", "updated_at").Find(&user, id).Error
					if err != nil {
						log.Printf("Error fetching user for GraphQL: %v", err)
						return nil, err
					}
					return user, nil
				},
			},
			"reports": &graphql.Field{
				Type:        graphql.NewList(ReportType),
				Description: "Get all reports",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var reports []models.Report
					err := db.Find(&reports).Error
					if err != nil {
						log.Printf("Error fetching reports for GraphQL: %v", err)
						return nil, err
					}
					return reports, nil
				},
			},
			"report": &graphql.Field{
				Type:        ReportType,
				Description: "Get a single report by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if !ok {
						return nil, fmt.Errorf("invalid report ID")
					}
					var report models.Report
					err := db.Find(&report, id).Error
					if err != nil {
						log.Printf("Error fetching report for GraphQL: %v", err)
						return nil, err
					}
					return report, nil
				},
			},
			"logEntries": &graphql.Field{
				Type:        graphql.NewList(LogEntryType),
				Description: "Get all log entries",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var logEntries []models.LogEntry
					err := db.Order("timestamp DESC").Find(&logEntries).Error
					if err != nil {
						log.Printf("Error fetching log entries for GraphQL: %v", err)
						return nil, err
					}
					return logEntries, nil
				},
			},
			"logEntry": &graphql.Field{
				Type:        LogEntryType,
				Description: "Get a single log entry by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if !ok {
						return nil, fmt.Errorf("invalid log entry ID")
					}
					var logEntry models.LogEntry
					err := db.Find(&logEntry, id).Error
					if err != nil {
						log.Printf("Error fetching log entry for GraphQL: %v", err)
						return nil, err
					}
					return logEntry, nil
				},
			},
			"operationalPages": &graphql.Field{
				Type:        graphql.NewList(OperationalPageType),
				Description: "Get all operational pages",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var pages []models.OperationalPage
					err := db.Find(&pages).Error
					if err != nil {
						log.Printf("Error fetching operational pages for GraphQL: %v", err)
						return nil, err
					}
					return pages, nil
				},
			},
			"operationalPage": &graphql.Field{
				Type:        OperationalPageType,
				Description: "Get a single operational page by ID or slug",
				Args: graphql.FieldConfigArgument{
					"idOrSlug": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idOrSlug, ok := p.Args["idOrSlug"].(string)
					if !ok {
						return nil, fmt.Errorf("invalid operational page ID or slug")
					}
					var page models.OperationalPage
					err := db.Where("id = ? OR slug = ?", idOrSlug, idOrSlug).Find(&page).Error
					if err != nil {
						log.Printf("Error fetching operational page for GraphQL: %v", err)
						return nil, err
					}
					return page, nil
				},
			},
			"operationalPageComponents": &graphql.Field{
				Type:        graphql.NewList(OperationalPageComponentType),
				Description: "Get components for a specific operational page",
				Args: graphql.FieldConfigArgument{
					"pageID": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					pageID, ok := p.Args["pageID"].(string)
					if !ok {
						return nil, fmt.Errorf("invalid page ID")
					}
					var components []models.OperationalPageComponent
					err := db.Where("page_id = ?", pageID).Order("display_order ASC").Find(&components).Error
					if err != nil {
						log.Printf("Error fetching operational page components for GraphQL: %v", err)
						return nil, err
					}
					return components, nil
				},
			},
		},
	})
}

// CreateSchema defines the executable GraphQL schema
func CreateSchema(db *gorm.DB) (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: RootQuery(db),
	})
}
