set dotenv-load

# Default entrypoint (defaults to list)
default: list

# List all available just commands
list:
    @just --list

# ── Setup ────────────────────────────────────────────────────

# Removes your old environment and replaces it with an example environment
init-env:
    just remove-envs
    cp .env.example .env
    just init-frontend-env
    just init-backend-env
    
init-frontend-env:
    echo "LOCAL_API_BASE_URL=http://localhost:{{env_var('BACKEND_HOST_PORT')}}" >> frontend/.env

init-backend-env:
    echo "POSTGRES_USER={{env_var('POSTGRES_USER')}}" >> backend/.env
    echo "POSTGRES_PASSWORD={{env_var('POSTGRES_PASSWORD')}}" >> backend/.env
    echo "POSTGRES_DB={{env_var('POSTGRES_DB')}}" >> backend/.env
    echo "POSTGRES_PORT={{env_var('POSTGRES_HOST_PORT')}}" >> backend/.env
    echo "POSTGRES_HOST=localhost" >> backend/.env
    echo "LOCAL_DEV=true" >> backend/.env

remove-envs:
    rm -f .env
    rm -f frontend/.env
    rm -f backend/.env

# Dump the data of the database into a seed.sql file
dump-seed:
    ${CONTAINER_TOOL} exec -t ft_transcendence_postgres pg_dump -U ${POSTGRES_USER} --data-only ft_transcendence > seed.sql

# Dump the schemas of the database into a schema.sql file
dump-schema:
    ${CONTAINER_TOOL} exec -t ft_transcendence_postgres pg_dump -U ${POSTGRES_USER} --schema-only ft_transcendence > schema.sql

# Seeds the database with example values
schema-db:
    ${CONTAINER_TOOL} exec -i ft_transcendence_postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} < db/schema.sql

# Seeds the database with example values
seed-db:
    ${CONTAINER_TOOL} exec -i ft_transcendence_postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} < db/example_seed.sql

# ── Production ────────────────────────────────────────────────────

# Deploys the app to production
prod:
    ${CONTAINER_ORCHESTRATION_TOOL} down --remove-orphans
    ${CONTAINER_ORCHESTRATION_TOOL} -f compose.yml up -d --build --force-recreate

# ── Development ────────────────────────────────────────────────────

# NOTE: The difference between production and development is that
# development uses the compose.override.yml file which will expose
# ports to the containers which can then be used for local
# development. The production container setup does not expose any
# ports expect for the reverse proxy!

# Standard dev way to start the entire app
up:
    ${CONTAINER_ORCHESTRATION_TOOL} up -d --build

# Stop all services
down:
    ${CONTAINER_ORCHESTRATION_TOOL} down

# Rebuild
re: down up

# Start a specific service (example: just serve postgres)
serve service:
    ${CONTAINER_ORCHESTRATION_TOOL} up -d --build {{service}}

# Show running containers
ps:
    @${CONTAINER_TOOL} ps --format "table {{{{.Names}}\t{{{{.Status}}\t{{{{.Ports}}"

# Show log for a service
logs service:
    ${CONTAINER_ORCHESTRATION_TOOL} logs -f {{service}}

# Enter the postgres database
enter-db:
    ${CONTAINER_TOOL} exec -it ft_transcendence_postgres psql -U $POSTGRES_USER -d $POSTGRES_DB

# Drops the database
drop-db:
    ${CONTAINER_TOOL} exec -t ft_transcendence_postgres psql -U $POSTGRES_USER -d postgres -c "DROP DATABASE IF EXISTS ft_transcendence;"

# Reset the database and gives you a fresh slate
reset-db: drop-db
    ${CONTAINER_TOOL} exec -t ft_transcendence_postgres psql -U $POSTGRES_USER -d postgres -c "CREATE DATABASE ft_transcendence;"

# ── Cleanup ────────────────────────────────────────────────────

# Remove all stopped containers
clean-containers:
    ${CONTAINER_TOOL} rm -f $(${CONTAINER_TOOL} ps -aq) 2>/dev/null || true

# Remove all unused volumes
clean-volumes:
    ${CONTAINER_TOOL} volume prune -f

# Remove dangling images
clean-images:
    ${CONTAINER_TOOL} image prune -f

# You probably dont have to use this
deepclean-images:
    ${CONTAINER_TOOL} rmi $(${CONTAINER_TOOL} images -qa) -f

# Full cleanup (containers, volumes, images)
clean: clean-containers clean-volumes clean-images
    @echo "Cleaned up everything"
