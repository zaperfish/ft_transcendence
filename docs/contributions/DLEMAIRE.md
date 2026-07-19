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

## Feature

The event chat gives event participants a private chat room with persistent history and real-time messaging. The interface displays participant names, avatars, timestamps, connection status, and a message character counter. Access is restricted to authenticated event participants.

## Modules

### Real-time features using WebSockets — Major, 2 points

Each event has an isolated WebSocket room. The Go backend authenticates and authorizes connections, persists messages, broadcasts them to connected participants, handles disconnections, and removes empty rooms. The Next.js frontend combines stored history with live messages and reports the connection state.

## Individual contribution

The event chat was implemented across the Go backend and Next.js frontend:

- Created participant-protected REST and WebSocket endpoints.
- Built concurrent event rooms with message broadcasting and automatic cleanup.
- Added database persistence and retrieval of the latest 50 messages.
- Built the chat interface, including sender details, timestamps, connection feedback, message limits, and controlled auto-scrolling.
- Added validation for empty messages and limits of 2,000 characters and 8,000 bytes.
- Added backend tests for chat persistence, authorization, WebSocket behavior, and validation.

### Challenges and solutions

- **Concurrent users and room cleanup:** A synchronized hub manages event rooms, while channel-based room loops serialize joins, leaves, and broadcasts. Empty rooms remove themselves from the hub.
- **Combining history with live updates:** The frontend loads persisted messages before opening the WebSocket, then merges both sources by message ID to avoid duplicates.
- **Scrolling without interrupting the reader:** The interface auto-scrolls only on initial load or when the user is already viewing recent messages.
- **Consistent message limits:** Character and byte limits are enforced in the frontend, backend, and database so Unicode input is handled consistently and safely.

### AI Usage

- Topic research
- AI-assisted code reviews
- Documentation drafting
- Structural bug detection
- Design discussions and evaluation of implementation alternatives

### Resources

- Complete Intro to React, Brian Holt - Frontend Masters
- Go Programming Language, The (Addison-Wesley Professional Computing Series)
- Learning Go: An Idiomatic Approach to Real-World Go Programming
- Haverbeke, M. (2024). Eloquent JavaScript (4th ed.). No Starch Press.