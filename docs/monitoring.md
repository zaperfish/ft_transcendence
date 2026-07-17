# Monitoring

**Contributor:** dlemaire

## Feature

The monitoring system uses Prometheus and Grafana to observe the health and runtime behavior of the Go backend. It collects application metrics, provides dashboards for garbage collection and scheduler activity, and defines alerts for downtime, restarts, and high goroutine usage. Grafana is available through the application's HTTPS reverse proxy and requires authentication.

## Modules

### Monitoring system with Prometheus and Grafana — Major, 2 points

The backend exposes Prometheus metrics that are collected by a dedicated Prometheus service. Grafana uses Prometheus as its provisioned data source and automatically loads custom dashboards. Prometheus evaluates alerting rules for backend availability, restarts, and goroutine usage. Access to Grafana is secured with an administrator login and routed through Caddy over HTTPS.

## Individual contribution

The monitoring system was implemented across the backend and container infrastructure:

- Added the Prometheus metrics endpoint to the Go backend.
- Added Prometheus and Grafana services with persistent storage to the Compose stack.
- Provisioned the Prometheus data source and Grafana dashboards automatically.
- Created dashboards for Go garbage collection, memory allocation, goroutines, OS threads, `GOMAXPROCS`, and CPU usage.
- Added alerting rules for backend downtime, recent or repeated restarts, and sustained high goroutine counts.
- Routed Grafana through Caddy with HTTPS redirection and authenticated access.
- Added optional GC and scheduler pressure endpoints and scripts for demonstrating dashboard behavior locally.
- Documented how to run the pressure tests and interpret the resulting metrics.

### Challenges and solutions

- **Reproducible monitoring setup:** Prometheus configuration, alert rules, the Grafana data source, and dashboards are stored in the repository and provisioned automatically when the containers start.
- **Understanding Go runtime metrics:** Separate dashboards group related garbage-collection and scheduler metrics, making memory, goroutine, thread, and CPU behavior easier to interpret.
- **Generating observable runtime activity:** Disabled-by-default development endpoints create controlled memory or scheduler pressure, with scripts and safety limits for local demonstrations.
- **Serving Grafana below `/grafana/`:** Grafana is configured for a subpath and routed through Caddy, which redirects plain HTTP access to HTTPS while preserving the requested path.
