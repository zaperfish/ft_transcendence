set dotenv-load

# Default entrypoint (defaults to list)
default: list

# List all available just commands
list:
    @just --list

# ── Production ────────────────────────────────────────────────────

# Deploys the app to production
prod:
    podman-compose down --remove-orphans
    podman-compose -f compose.yml up -d --build --force-recreate

# ── Development ────────────────────────────────────────────────────

# NOTE: The difference between production and development is that
# development uses the compose.override.yml file which will expose
# ports to the containers which can then be used for local
# development. The production container setup does not expose any
# ports expect for the reverse proxy!

# Standard dev way to start the entire app
up:
    podman-compose up -d --build

# Stop all services
down:
    podman-compose down

# Rebuild
re: down up

# Start a specific service (example: just up postgres)
serve service:
    podman-compose up -d --build {{service}}

# Show running containers
ps:
    @podman ps --format "table {{{{.Names}}\t{{{{.Status}}\t{{{{.Ports}}"

# Show log for a service
logs service:
    podman-compose logs -f {{service}}

db:
    podman exec -it ft_transcendence_postgres psql -U $POSTGRES_USER -d $POSTGRES_DB

# ── Cleanup ────────────────────────────────────────────────────

# Remove all stopped containers
clean-containers:
    podman rm -f $(podman ps -aq) 2>/dev/null || true

# Remove all unused volumes
clean-volumes:
    podman volume prune -f

# Remove dangling images
clean-images:
    podman image prune -f

# You probably dont have to use this
deepclean-images:
    podman rmi $(podman images -qa) -f

# Full cleanup (containers, volumes, images)
clean: clean-containers clean-volumes clean-images
    @echo "Cleaned up everything"
