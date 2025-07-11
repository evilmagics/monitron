# Monitron System Specification

This document outlines the detailed specifications for Monitron, a web application designed for monitoring Instances, Services, Domains, and SSL certificates.

## 1. Instances

- **Primary Key:** UUID
- **Required Fields:** Name, Host, Check Interval (default: 10s), Check Timeout (default: 30s), Agent Port (default: 7773)
- **Optional Fields:** Agent Authentication (default: JWT), Description, Label, Group
- **Metrics Collection:**
    - CPU Usage, GPU Usage, Memory Usage, Storage Usage, Network (bytes and percent).
    - Provided as HTTP API response from agents.
- **Device Details Collection:**
    - OS, CPU, GPU, Memory, Storages, Networks.
    - Provided as HTTP API response from agents.
- **Configuration:** Check interval configurable via HTTP API.
- **Agent Connection:** Connects to Server via Agent using HTTP or WebSocket (configurable).
- **Remote Control:** Instances can be restarted or shut down via Agent using HTTP API (graceful restart for Agent).
- **Agent Restart:** Agent can be restarted via HTTP API (graceful restart).
- **Agent-Server Authentication:** JWT, Basic Auth, API Key, Secret Key (HMAC SHA-256).
- **UI Grouping:** Instances separated by group in UI.
- **Stats:** Response Time, Uptime, Last Checked, Average Response Time, Incident Total (additional relevant parameters can be added).

## 2. Services

- **Primary Key:** UUID
- **Required Fields:** Name, API Type, Check Interval (default: 10s), Timeout (default: 30s)
- **Optional Fields:** Description, Label, Group
- **Service Types:** HTTP API, gRPC, MQTT, TCP, DNS, Ping.
    - **HTTP API:** Method, Health URL (default: /health), Expected Status Code (default: 200).
    - **gRPC:** Host, Port, Authentication (optional), Default proto (client must use).
    - **MQTT:** Host, Port, QoS, Topic, Authentication (optional), other recommended parameters.
    - **Ping:** Host.
    - **DNS:** Domain Name.
    - **TCP:** Hostname/IP, Port.
- **Stats:** Response Time, Uptime, Last Checked, Average Response Time, Incident Total (additional relevant parameters can be added).
- **UI Grouping:** Services separated by group in UI.

## 3. Domain and SSL

- **Primary Key:** UUID
- **Required Fields:** Domain, Warning Threshold, Expiry Threshold, Check Interval (default: 1 day)
- **Optional Fields:** Label
- **Backend Parsing:** Backend parses domain/SSL for details:
    - Certificate Detail, Issuer, Valid from, Resolved IP, Expiry, Day Left (UI only), etc.
- **Stats:** Response Time, Uptime, Last Checked, Average Response Time, Incident Total (additional relevant parameters can be added).

## 4. Operational Page

- **Functionality:** Publish selected Services and Domain/SSL as an operational page.
- **Themes:** Dark and Light themes.
- **Primary Key:** Slug and UUID (for modification/deletion).
- **Components:** A page can have multiple service, domain, and SSL components. Components are reusable across pages.
- **Visibility:** Public or Private (requires authentication).
- **Component Display (on page):** Component Name (user input, different from Service Name), Display order, Description, Uptime History (30 days), Last Checked, Status, Average Response Time.
- **Page Performance Analysis:**
    - Overall Uptime, Incidents Total, Average Response Time, Uptime History.
    - Data stored and displayed for the last 30 days.
    - Refresh button to retrieve data from server.
    - Report button (10 requests/day per IP, 1000 reporting requests/day per page) with predefined reasons.

## 5. Reports

- Stores various reports from Services, Instances, Domain/SSL, and Operational Pages.
- Displayed in the Report menu.
- Stored in the database without additional time during request processing.
- All backend activities must have logs and reports.
- Implements GraphQL for efficient querying.
- **Scheduled Reports:** Allow users to schedule regular reports (daily, weekly, monthly) to be generated and sent via email or other notification channels.
- **Data Export Formats:** Provide options to export report data in various formats (CSV, PDF, Excel) for further analysis outside Monitron.

## 6. Response Structure

```json
{
    "code": XXXXX,
    "message": "",
    "data": {},
    "error": {}
}
```

- **Code:** Integer representing error code (00000 for success).
- **Message:** Primary reason for error or success category.
- **Data:** Object or array of data returned on success.
- **Error:** Error object with input field as key and error reason as value.

## 7. Monitron (Server/Backend)

- **Tech Stack:** Golang, PostgreSQL, RabbitMQ, Loki, Grafana, Prometheus, GraphQL.
- **Authentication:** JWT (default), Basic Auth, API Key, Secret Key (HMAC SHA-256).
- **HTTP Handler:** Fiber (Go Fiber) for HTTP requests/responses (JSON payload).
    - Response payload: Code (error code), Message (error/success message), Data (response data).
    - All incoming requests include `RequestID` in the header.
- **WebSocket:** For UI and Agent connection.
- **Logging:** Stores all application logs for Loki Agent collection.
- **Metrics:** Provides metrics for Prometheus.
- **HTTP API Endpoints:** All necessary endpoints for Agent and UI.
- **Authentication Configuration:** Via Environment variables or UI (Admin), stored in database and in-memory cache.
- **Access Control:** Role-Based Access Control with Admin as the highest user.
- **Health/Usage Checks:** Performs health check, usage check, and device info on Instances via Agent at configured intervals.
    - Agent registration required when adding services.
    - Usage check via HTTP API or WebSocket (user configurable).
- **Service Health Check:** Performs health checks on Services at configured intervals.
- **Domain/SSL Health Check:** Performs health checks and domain detail checks at configured intervals (min 4 hours, default 24 hours).
- **Background Processes:** Health check, usage check, device info, domain detail performed in background with RabbitMQ as Message Broker.
- **Secure Authentication Storage:** Stores authentication securely for service/instance checks.
- **Documentation:** All features documented using Stoplight Element or Swagger (with Tailwind or ShadcnUI for better UI).
- **Notifications:** Email, Telegram, Discord notifications configurable via UI.
    - Email sending in background via RabbitMQ (max 10 retries, min 1s timeout).
- **Startup Check:** On program start, re-checks all Instances, Services, and Domain/SSL that have passed their last checked interval.
- **Reporting & Logging:** Provides reporting and logging accessible by UI via GraphQL for filtering.
- **Alerting Enhancements:**
    - **Customizable Alert Rules:** Allow users to define complex alert rules based on various metrics.
    - **Escalation Policies:** Implement escalation policies for alerts, notifying different teams or individuals if an alert remains unacknowledged or unresolved.
    - **Maintenance Windows:** Allow scheduling of maintenance windows to suppress alerts during planned downtime.
- **Role-Based Access Control (RBAC) Granularity:** Further refine RBAC to allow more granular permissions, e.g., read-only access to specific groups of instances or services.
- **Distributed Tracing:** Integrate with a distributed tracing system (e.g., OpenTelemetry) to gain deeper insights into request flows across microservices.
- **TimescaleDB Extension:** For time-series data (metrics, uptime history), consider using TimescaleDB, an open-source extension for PostgreSQL that optimizes for time-series workloads, offering better performance and compression.
- **Alertmanager:** Explicitly use Alertmanager for alert routing, deduplication, and silencing for a production-grade alerting system.
- **Message Persistence and Durability:** Ensure that critical messages (e.g., notifications, health check tasks) are configured for persistence and durability in RabbitMQ to prevent data loss during broker restarts.
- **Dead Letter Queues (DLQ):** Implement DLQs to handle messages that cannot be processed successfully, allowing for debugging and reprocessing.
- **API Documentation:** Emphasize the importance of maintaining up-to-date and interactive API documentation (Stoplight Element or Swagger UI) throughout the development lifecycle.
- **CI/CD Pipeline:** Establish a robust CI/CD pipeline for automated testing, building, and deployment of all Monitron components.
- **Security Best Practices:** Adhere to security best practices throughout development, including secure coding guidelines, regular security audits, and vulnerability scanning.

## 8. Monitron UI

- **Tech Stack:** NextJS, ShadcnUI, Bun, GraphQL.
- **Themes:** Dark and Light themes (configurable via UI, stored per user in Backend).
- **Color Themes:** Light blue, Orange, Light green (configurable via UI, stored per user in Backend).
- **Pages:**
    - **Profile Page:** Update user profile.
    - **Settings Page:** Configure theme, user, notifications, alerts, and templates.
    - **Dashboard Page:** Summary of Instances, Services, Domain/SSL, and Reports.
    - **Instances Page:** List of registered instances and summary. Detail page, edit modal, delete feature. Add and refresh features.
    - **Services Page:** List of registered services and summary. Detail page, edit modal, delete feature. Add and refresh features.
    - **Domain/SSL Page:** List of registered Domain/SSL and summary. Detail page, edit modal, delete feature. Add and refresh features.
    - **Operational Page:** Card list of registered pages and summary. Detail page, edit modal, delete feature. Add and refresh features.
    - **Report Page:** Reports and logs from Instances, Services, Domain/SSL, and application.
        - Separate pages for Instances, Services, Domain/SSL, and Application.
        - Application reports include HTTP Handler metrics.
        - Log detail pages vary by need.
        - All reports filterable by time and other suggested filters (GraphQL recommended).
- **WebSockets for Real-time Updates:** Leverage WebSockets extensively for real-time updates on dashboards and operational pages to provide immediate feedback on instance/service status changes.
- **Charting Libraries:** For rich data visualizations, explore charting libraries like Recharts, Nivo, or ECharts that integrate well with React and can handle dynamic data.

## 9. Monitron Agent

- **Tech Stack:** Golang.
- **Authentication:** JWT (default), Basic Auth, API Key, Secret Key (HMAC SHA-256).
- **HTTP Handler:** Fiber (Go Fiber).
- **Data Transmission:** HTTP or WebSocket for sending Usage metrics to server.
- **Authentication Details:**
    - **JWT:** Login endpoint for Access token and Refresh token. Refresh token stored in Go-Cache for 24 hours.
    - **Basic Auth:** Uses same login endpoint, no JWT implementation.
    - **API Key:** Taken from environment, used in header.
    - **Secret Key:** Taken from environment, verifies payload with HMAC SHA-256.
- **Background Service:** Collects usage metrics every 1s (default, configurable via env) and stores in memory cache as JSON for 5s (default, configurable via env). Expiry time > interval.
- **Usage Metrics Sending:** Metrics sent to server from cache. If not available, collect and store in cache before sending.
- **HTTP API Endpoint:** For retrieving device details (OS, CPU, GPU, Memory, Storages, Network Interfaces).
- **Agent Health Monitoring:** Implement monitoring of the agent's own health and connectivity to ensure data collection is uninterrupted.
- **Open Source Contributions:** Consider open-sourcing parts of the project (e.g., the agent) to foster community contributions and wider adoption.




- **GraphQL Usage:** GraphQL will be used exclusively for **query operations** (read-only access) to the database. All data creation, update, and deletion (CUD) operations will be handled via ORM (Object-Relational Mapping) within the backend application logic.
- **ORM for CUD Operations:** The backend will utilize an ORM to interact with the database for all CUD operations, ensuring data integrity, type safety, and simplifying database interactions.
- **Backend Database Management:** The backend application will be responsible for database schema management, including:
    - **Schema Creation:** Automatically create the database schema if it does not exist on application startup.
    - **Migrations:** Handle database migrations to evolve the schema over time, ensuring compatibility with new application versions.
    - **Database Initialization:** Initialize necessary seed data or configurations if the database is newly created.




## 10. User Management

Monitron will include comprehensive user management features, primarily accessible to administrators, to ensure secure and controlled access to the system.

### 10.1. User List (Admin Access Only)

-   **Functionality:** Administrators will have access to a dedicated page or section in the UI to view a list of all registered users.
-   **Details Displayed:** User ID, Username, Email, Role, Account Status (active/inactive), Last Login, Created At, Updated At.
-   **Filtering & Sorting:** Ability to filter users by role, status, or search by username/email. Sorting by various fields.
-   **Access Control:** This feature will be strictly restricted to users with an 'admin' role.

### 10.2. User Management Access (Admin Only)

-   **Functionality:** Administrators will be able to perform management actions on individual user accounts from the user list or a dedicated user detail page.
-   **Actions:**
    -   **Change User Role:** Modify a user's role (e.g., from 'user' to 'admin' or vice-versa).
    -   **Activate/Deactivate User:** Enable or disable a user's account, preventing them from logging in.
    -   **Reset Password (Admin Initiated):** Force a password reset for a user, typically generating a temporary password or sending a reset link to their email.
    -   **View User Details:** Access comprehensive details of a specific user.
-   **Access Control:** These management actions will be strictly restricted to users with an 'admin' role.

### 10.3. Forgot Password Flow

-   **Initiation:** User navigates to a 


dedicated "Forgot Password" page.
-   **Step 1 (Request Reset):** User enters their registered email address. The system verifies the email and sends a unique, time-limited password reset token to that email address.
-   **Step 2 (Token Verification):** User receives the email, clicks on the provided link, which directs them to a password reset page. The system verifies the token for validity and expiry.
-   **Step 3 (Set New Password):** If the token is valid, the user is prompted to enter and confirm a new password. The new password must meet defined complexity requirements.
-   **Security:** Tokens will be single-use and expire after a short period (e.g., 15-30 minutes). Rate limiting will be applied to password reset requests to prevent abuse.

### 10.4. Change Password Flow (Authenticated User)

-   **Initiation:** Authenticated user navigates to their profile or settings page.
-   **Step 1 (Current Password Verification):** User is prompted to enter their current password for verification, along with the new password and its confirmation.
-   **Step 2 (Set New Password):** If the current password is correct and the new password meets complexity requirements, the system updates the user's password.
-   **Security:** Implement strong password policies (minimum length, complexity, history). Log password change events.

### 10.5. Delete User (Admin Access Only)

-   **Functionality:** Administrators will have the ability to permanently delete user accounts.
-   **Confirmation:** A confirmation step will be required (e.g., re-entering admin password or a specific confirmation phrase) to prevent accidental deletion.
-   **Impact:** Deleting a user will remove all associated user data. Consider soft deletes or archiving for audit purposes if required in the future.
-   **Access Control:** This feature will be strictly restricted to users with an 'admin' role.

These features will enhance the administrative control and user experience within the Monitron platform, ensuring secure and efficient user management.

