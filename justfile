set dotenv-load

# Default entrypoint (defaults to list)
default: list

# List all available just commands
list:
    @just --list

# ── Setup ────────────────────────────────────────────────────

# Removes your old environment and replaces it with an example environment
init-env:
    rm .env
    rm frontend/.env
    rm backend/.env
    cp .env.example .env
    echo "LOCAL_API_BASE_URL=http://localhost:{{env_var('BACKEND_HOST_PORT')}}" >> frontend/.env
    echo "POSTGRES_USER={{env_var('POSTGRES_USER')}}" >> backend/.env
    echo "POSTGRES_PASSWORD={{env_var('POSTGRES_PASSWORD')}}" >> backend/.env
    echo "POSTGRES_DB={{env_var('POSTGRES_DB')}}" >> backend/.env
    echo "POSTGRES_PORT={{env_var('POSTGRES_HOST_PORT')}}" >> backend/.env
    echo "POSTGRES_HOST=localhost" >> backend/.env
    echo "LOCAL_DEV=true" >> backend/.env

# Dump the data of the database into a seed.sql file
dump-seed:
    podman exec -t ft_transcendence_postgres pg_dump -U ${POSTGRES_USER} --data-only ft_transcendence > seed.sql

# Dump the schemas of the database into a schema.sql file
dump-schema:
    podman exec -t ft_transcendence_postgres pg_dump -U ${POSTGRES_USER} --schema-only ft_transcendence > schema.sql

# Seeds the database with example values
schema-db:
    podman exec -i ft_transcendence_postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} < db/schema.sql

# Seeds the database with example values
seed-db:
    podman exec -i ft_transcendence_postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} < db/example_seed.sql

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

# Start a specific service (example: just serve postgres)
serve service:
    podman-compose up -d --build {{service}}

# Show running containers
ps:
    @podman ps --format "table {{{{.Names}}\t{{{{.Status}}\t{{{{.Ports}}"

# Show log for a service
logs service:
    podman-compose logs -f {{service}}

# Enter the postgres database
enter-db:
    podman exec -it ft_transcendence_postgres psql -U $POSTGRES_USER -d $POSTGRES_DB

# Drops the database
drop-db:
    podman exec -t ft_transcendence_postgres psql -U $POSTGRES_USER -d postgres -c "DROP DATABASE IF EXISTS ft_transcendence;"

# Reset the database and gives you a fresh slate
reset-db: drop-db
    podman exec -t ft_transcendence_postgres psql -U $POSTGRES_USER -d postgres -c "CREATE DATABASE ft_transcendence;"

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
