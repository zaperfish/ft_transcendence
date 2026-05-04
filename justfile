# Default entrypoint (defaults to list)
default: list

# List all available just commands
list:
    @just --list

# Deploys the app to production
prod:
    podman-compose down --remove-orphans
    podman-compose -f compose.yml up -d --build --force-recreate

# Standard dev way to start the entire app
up:
    podman-compose up -d
    sleep 2
    @podman ps --format "table {{{{.Names}}\t{{{{.Status}}\t{{{{.Ports}}"

# Stop all services
down:
    podman-compose down

# Start a specific service (example: just up postgres)
up-service service:
    podman-compose up -d {{service}}

# Show running containers
ps:
    @podman ps --format "table {{{{.Names}}\t{{{{.Status}}\t{{{{.Ports}}"

# Show log for a service
logs service:
    podman-compose logs -f {{service}}

# ── Cleanup ────────────────────────────────────────────────────
# Remove all stopped containers
clean-containers:
    podman rm -f $(podman ps -aq) 2>/dev/null || true

# Remove all unused volumes
clean-volumes:
    podman volume prune -f

# Remove all unused images
clean-images:
    podman image prune -f

# Full cleanup (containers, volumes, images)
clean: clean-containers clean-volumes clean-images
    @echo "Cleaned up everything"
