# Monitron: Todo

### To do
- [ ] Change logging using Zerolog and write logs to file. Prepare Loki and Promtail config on config directory.
- [ ] Implement API specification on swagger.
- [ ] Add feature services health check. Each services have different interval, use RabbitMQ as message broker. After checked, send next check time to message broker to schedule next health check. Implementation auth refer to saved auth type.
- [ ] Add unit tests

### Done
- [x] Add feature health check to services, instances and domain/SSL
- [x] Change Net/http to use Fiber for HTTP Handler and Resty for HTTP Client
- [x] Change SQLx and native query to use ORM
- [x] Validate user input
- [x] Change Net/http to use Fiber for HTTP Handler and Resty for HTTP Client
- [x] Change SQLx and native query to use ORM
- [x] Validate user input
