# Runtime Pressure Test Endpoints

This document describes the local debug endpoints used to demonstrate Go runtime behavior in the Grafana monitoring dashboards.

These endpoints are disabled by default. They are intended for local monitoring demos and should not be enabled in production.

## Enable The Endpoints

Each endpoint has its own environment flag. Pass the flag to Compose when rebuilding the backend:

```bash
ENABLE_GC_PRESSURE_TEST=true docker compose up -d --build backend
ENABLE_SCHEDULER_PRESSURE_TEST=true docker compose up -d --build backend
```

## GC Pressure

Endpoint:

```text
POST /api/debug/gc-pressure
```

Enable flag:

```env
ENABLE_GC_PRESSURE_TEST=true
```

Script:

```bash
./scripts/trigger-gc-pressure.sh
```

Script parameters:

```bash
SIZE_MB=128 HOLD_SECONDS=30 ./scripts/trigger-gc-pressure.sh
```

Direct request:

```bash
curl --insecure -X POST "https://localhost:7443/api/debug/gc-pressure?size_mb=128&hold_seconds=30"
```

Parameters:

| Parameter | Default | Maximum | Description |
| --- | ---: | ---: | --- |
| `size_mb` | `128` | `256` | Temporary heap allocation size in megabytes. |
| `hold_seconds` | `30` | `120` | How long the allocation is kept reachable before it is released. |

Expected response:

```json
{
  "allocation_id": 1,
  "allocated_mb": 128,
  "hold_seconds": 30,
  "active_loads": 1
}
```

Expected Grafana behavior:

| Panel | Expected behavior |
| --- | --- |
| `GC Cycles/sec` | Spikes upward when memory pressure causes GC activity, then falls back toward baseline. |
| `Heap Allocated Bytes` | Rises while the temporary allocation is held, then drops after release and GC. |
| `Allocation Rate` | Spikes while the temporary allocation is created, then falls back toward baseline. |
| `Process Resident Memory Bytes` | Rises from the OS perspective and may stay elevated because Go can keep memory reserved for reuse. |

What this demonstrates:

```text
The backend creates temporary heap pressure. Go allocates memory, GC cycles run, heap usage drops after release, and process memory may remain higher than heap allocation.
```

## Scheduler Pressure

Endpoint:

```text
POST /api/debug/scheduler-pressure
```

Enable flag:

```env
ENABLE_SCHEDULER_PRESSURE_TEST=true
```

Script:

```bash
./scripts/trigger-scheduler-pressure.sh
```

Script parameters:

```bash
MODE=sleep GOROUTINES=1000 HOLD_SECONDS=30 ./scripts/trigger-scheduler-pressure.sh
MODE=cpu GOROUTINES=500 HOLD_SECONDS=30 ./scripts/trigger-scheduler-pressure.sh
```

Direct request:

```bash
curl --insecure -X POST "https://localhost:7443/api/debug/scheduler-pressure?mode=sleep&goroutines=1000&hold_seconds=30"
curl --insecure -X POST "https://localhost:7443/api/debug/scheduler-pressure?mode=cpu&goroutines=500&hold_seconds=30"
```

Parameters:

| Parameter | Default | Maximum | Description |
| --- | ---: | ---: | --- |
| `mode` | `sleep` | n/a | Pressure mode. Valid values are `sleep` and `cpu`. |
| `goroutines` | `1000` | `5000` | Number of temporary goroutines to start. |
| `hold_seconds` | `30` | `120` | How long the goroutines stay active. |

Expected response:

```json
{
  "load_id": 1,
  "mode": "sleep",
  "goroutines": 1000,
  "hold_seconds": 30,
  "active_loads": 1,
  "active_goroutines": 1000
}
```

### Sleep Mode

`sleep` mode creates many goroutines that wait until the load is released. They stay alive but do not continuously perform CPU work.

Expected Grafana behavior:

| Panel | Expected behavior |
| --- | --- |
| `Goroutines` | Clear upward spike for the configured hold duration, then a drop after release. |
| `OS Threads` | Usually stable, but may move slightly depending on runtime behavior. |
| `GOMAXPROCS` | Stays flat because the endpoint does not change Go's parallel execution limit. |
| `Process CPU Usage` | Small spike during goroutine creation, then close to baseline. |

What this demonstrates:

```text
Go can keep many blocked goroutines alive without using much CPU.
```

### CPU Mode

`cpu` mode creates many goroutines that repeatedly execute dummy arithmetic until the load is released.

Expected Grafana behavior:

| Panel | Expected behavior |
| --- | --- |
| `Goroutines` | Clear upward spike for the configured hold duration, then a drop after release. |
| `OS Threads` | May increase slightly if the runtime creates or keeps extra OS threads. |
| `GOMAXPROCS` | Stays flat because it is the limit for parallel Go execution, not the total thread count. |
| `Process CPU Usage` | Clear upward spike while CPU-bound goroutines are active. |

What this demonstrates:

```text
Many runnable CPU-bound goroutines compete for execution. CPU usage rises, while GOMAXPROCS still limits how many OS threads can execute Go code at the same time.
```

## Dashboard Access

Grafana is available through Caddy:

```text
https://localhost:7443/grafana/
```

Prometheus is available through Caddy:

```text
https://localhost:7443/prometheus/
```

Useful Prometheus pages:

| Page | URL |
| --- | --- |
| Alerts | `https://localhost:7443/prometheus/alerts` |
| Targets | `https://localhost:7443/prometheus/targets` |
| Query interface | `https://localhost:7443/prometheus/query` |

The dashboards are provisioned in the `Monitoring` folder:

```text
Go GC Runtime Overview
Go Scheduler Runtime Overview
```

## Test The BackendDown Alert

The `BackendDown` alert fires when Prometheus cannot scrape the backend for one minute. To test it locally, first open the Prometheus alerts page:

```text
https://localhost:7443/prometheus/alerts
```

Stop only the backend container:

```bash
docker compose stop backend
```

Prometheus scrapes the backend every five seconds. The alert initially appears as `Pending` and changes to `Firing` after the backend has remained unavailable for one minute.

You can also confirm that the backend target is down on:

```text
https://localhost:7443/prometheus/targets
```

After confirming the alert, start the backend again:

```bash
docker compose start backend
```

The `BackendDown` alert resolves after Prometheus successfully scrapes the backend again. Starting the backend may also trigger the `BackendRestartedRecently` alert, which is expected because the backend process start time changed.
