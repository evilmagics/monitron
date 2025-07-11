CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS instances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    host VARCHAR(255) NOT NULL,
    check_interval INT NOT NULL DEFAULT 10,
    check_timeout INT NOT NULL DEFAULT 30,
    agent_port INT NOT NULL DEFAULT 7773,
    agent_auth VARCHAR(255) NOT NULL DEFAULT 'JWT',
    description TEXT,
    label VARCHAR(255),
    "group" VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS instance_stats (
    instance_id UUID PRIMARY KEY REFERENCES instances(id) ON DELETE CASCADE,
    response_time DECIMAL(10, 2) NOT NULL,
    uptime DECIMAL(5, 2) NOT NULL,
    last_checked TIMESTAMP WITH TIME ZONE NOT NULL,
    average_response_time DECIMAL(10, 2) NOT NULL,
    incident_total INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS instance_metrics (
    instance_id UUID REFERENCES instances(id) ON DELETE CASCADE,
    metric_type VARCHAR(255) NOT NULL,
    value DECIMAL(10, 2) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (instance_id, metric_type, timestamp)
);

CREATE TABLE IF NOT EXISTS device_info (
    instance_id UUID PRIMARY KEY REFERENCES instances(id) ON DELETE CASCADE,
    os VARCHAR(255),
    cpu VARCHAR(255),
    gpu VARCHAR(255),
    memory VARCHAR(255),
    storages TEXT,
    networks TEXT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);




CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    api_type VARCHAR(255) NOT NULL,
    check_interval INT NOT NULL DEFAULT 10,
    timeout INT NOT NULL DEFAULT 30,
    description TEXT,
    label VARCHAR(255),
    "group" VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- HTTP API fields
    http_method VARCHAR(10),
    http_health_url TEXT,
    http_expected_status INT,

    -- gRPC fields
    grpc_host VARCHAR(255),
    grpc_port INT,
    grpc_auth TEXT,
    grpc_proto TEXT,

    -- MQTT fields
    mqtt_host VARCHAR(255),
    mqtt_port INT,
    mqtt_qos INT,
    mqtt_topic TEXT,
    mqtt_auth TEXT,

    -- TCP fields
    tcp_host VARCHAR(255),
    tcp_port INT,

    -- DNS fields
    dns_domain_name VARCHAR(255),

    -- Ping fields
    ping_host VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS service_stats (
    service_id UUID PRIMARY KEY REFERENCES services(id) ON DELETE CASCADE,
    response_time DECIMAL(10, 2) NOT NULL,
    uptime DECIMAL(5, 2) NOT NULL,
    last_checked TIMESTAMP WITH TIME ZONE NOT NULL,
    average_response_time DECIMAL(10, 2) NOT NULL,
    incident_total INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);




CREATE TABLE IF NOT EXISTS domain_ssl (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    domain VARCHAR(255) NOT NULL,
    warning_threshold INT NOT NULL,
    expiry_threshold INT NOT NULL,
    check_interval INT NOT NULL DEFAULT 1,
    label VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Parsed SSL/Domain details
    certificate_detail TEXT,
    issuer VARCHAR(255),
    valid_from TIMESTAMP WITH TIME ZONE,
    resolved_ip VARCHAR(255),
    expiry TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS domain_ssl_stats (
    domain_ssl_id UUID PRIMARY KEY REFERENCES domain_ssl(id) ON DELETE CASCADE,
    response_time DECIMAL(10, 2) NOT NULL,
    uptime DECIMAL(5, 2) NOT NULL,
    last_checked TIMESTAMP WITH TIME ZONE NOT NULL,
    average_response_time DECIMAL(10, 2) NOT NULL,
    incident_total INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);




CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);




CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    report_type VARCHAR(255) NOT NULL,
    format VARCHAR(50) NOT NULL,
    generated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    file_path TEXT NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS log_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    level VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    service VARCHAR(255),
    request_id VARCHAR(255)
);




CREATE TABLE IF NOT EXISTS operational_pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS operational_page_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id UUID NOT NULL REFERENCES operational_pages(id) ON DELETE CASCADE,
    component_type VARCHAR(50) NOT NULL,
    component_id UUID NOT NULL,
    component_name VARCHAR(255) NOT NULL,
    display_order INT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(page_id, component_id) -- A component can only be added once per page
);

CREATE TABLE IF NOT EXISTS operational_page_stats (
    page_id UUID PRIMARY KEY REFERENCES operational_pages(id) ON DELETE CASCADE,
    overall_uptime DECIMAL(5, 2) NOT NULL,
    incidents_total INT NOT NULL DEFAULT 0,
    average_response_time DECIMAL(10, 2) NOT NULL,
    uptime_history TEXT, -- JSON string of 30-day history
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

