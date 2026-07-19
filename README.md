# ft_transcendence
ft_transcendence is the final group project of the 42 Common Core curriculum. It is a collaborative full-stack web application designed to simulate a real-world production system.

## Navigation

- [About This Project](#about-this-project)
- [Team Roles](#team-roles)
- [Technologies Used](#technologies-used)
- [Architecture](#architecture)
- [Instructions](#instructions)
    - [Container tooling](#container-tooling)
    - [Trust the local HTTPS certificate](#trust-the-local-https-certificate)
- [Contributing](#contributing)
- [Implementation](#implementation)
    - [Authentification](#authentification)
<br><br>

## About This Project
This project is a clone of the popular Meetup app, implementing its core functionality to help users discover, create, and join events.
<br><br>

## Team Roles
**Product Owner (PO):** yingzhan

**Project Manager (PM) / Scrum Master:** alneumann

**Technical Lead / Architect:** lmiehler

**Developers:** dlemaire, shutan
<br><br>

## Modules

Total achieved points: 4

### List of potential modules

| Module | Status | Points |
|--------|--------|--------|
| Use a framework for both frontend and backend | ✅ | 2 |
| Real-time features (WebSockets or similar) |  | 2 |
| User interaction system (chat, profiles, friends) |  | 2 |
| Public API (auth, rate limiting, 5 endpoints) |  | 2 |
| Organization system (CRUD + user management) |  | 2 |
| User management & authentication |  | 2 |
| Advanced permissions system (roles & access control) |  | 2 |
| ORM usage | ✅ | 1 |
| Notification system |  | 1 |
| Server-Side Rendering (SSR) | ✅ | 1 |
| Custom design system |  | 1 |
| Advanced search (filters, sorting, pagination) |  | 1 |
| OAuth 2.0 authentication |  | 1 |
| Two-Factor Authentication (2FA) |  | 1 |
| User analytics dashboard |  | 1 |
| Monitoring system (Prometheus + Grafana) |  | 2 |
| Health checks & backup system |  | 1 |
| CI/CD deployment pipeline |  | 1 |

## Technologies Used
 
- **Go** — Backend language, chosen for its performance, low resource overhead, and strong support for concurrent request handling, which suits a lightweight API server.
- **Huma** — Go API framework, chosen for generating OpenAPI documentation and request/response validation directly from Go structs and types, reducing boilerplate while keeping the API self-documenting.
- **PostgreSQL** — Primary database, chosen for its reliability, strong support for relational data and constraints, and mature ecosystem for transactional workloads.
- **React** — Frontend framework, chosen for its component-based architecture, large ecosystem, and ease of building interactive, maintainable user interfaces.
- **Caddy** — Reverse proxy and web server, chosen for automatic HTTPS, simple configuration via Caddyfile, and straightforward routing to backend/frontend services.

## Architecture

This project is a containerized full-stack web application. It follows a service-oriented architecture where the frontend, backend, database, and reverse proxy are isolated into separate components and communicate over internal networks.

The diagram below shows the architecture with its most important components.

<p align="center">
  <img src="./docs/assets/prod_architecture.drawio.svg" width="400"/>
</p>

<p align="center">
  <img src="./docs/assets/legend.drawio.svg" width="400"/>
</p>

When a user visits our website, the DNS resolves the domain name to the IP address of the VPS where the ft_transcendence project is hosted. The browser then sends HTTP(S) requests to that IP, which is handled by the server running on the VPS.
In our setup, these requests first reach the system-level reverse proxy, Caddy, which listens on ports 80 (HTTP) and 443 (HTTPS). The purpose of this reverse proxy layer is to manage incoming traffic for multiple services running on the same machine (example: there might be another different website running on the same machine) and route requests to the appropriate application based on the domain name and path.

For the domain ft-transcendence.zaper.io, the system Caddy acts as an entry-level reverse proxy and forwards traffic to the Transcendence-specific Caddy instance. The Transcendence Caddy then routes requests internally, for example directing /api/* requests to the backend API while serving frontend traffic for all other routes.

When a user presses a button that triggers an API call, the request is sent from the browser and routed through Caddy to the backend service, typically via a path such as /api/*.

An interesting aspect of the architecture is the potential communication between the frontend and backend containers which the diagram above shows. In most cases, these services do not communicate directly. Instead, all client-triggered requests pass through the reverse proxy.

However, there is one important exception: when using Server Side Rendering (SSR), the frontend container itself may act as a server. In this case, during Server-Side Rendering (SSR), the Next.js server can make internal HTTP requests to the backend API before sending the fully rendered HTML to the client.

### Architecture during development

During development, the architecture differs slightly from production. Running the entire application inside containers can be a little tedious and more complex. Instead, we use a hybrid approach: **local development with dependencies running in containers**.

While developing entirely inside containers offers better reproducibility and helps avoid “it works on my machine” issues, it requires more setup and can introduce additional complexity. For most development tasks, working directly on your local machine is faster, more convenient, and feels more natural (still comes with its own pitfalls though).

> **Note:** I only know one transcendence repository which does development inside of containers. If you are interested in what that looks like, you can take a look at it here https://github.com/cubernetes/ft-transcendence. The key aspect is the use of Dockers file-watching functionality (in compose it's the use of the watch directive).

#### Example: Working on the frontend

When developing the frontend:

- Ensure that all required backend services are running — typically via containers.
- Start the frontend development server locally using: `npm run dev`


<p align="center">
  <img src="./docs/assets/frontend_dev_architecture.drawio.svg" width="400"/>
</p>

#### Example: Working on the backend

<p align="center">
  <img src="./docs/assets/backend_dev_architecture.drawio.svg" width="400"/>
</p>

## Instructions

### Container tooling

This project uses Compose files to describe the container stack, and those files can be run with either Docker or Podman.

- `Docker` and `Podman` are container runtimes. They build images and run containers.
- `docker-compose` / `docker compose` and `podman compose` are Compose-compatible orchestration commands. They start the full stack from `compose.yml`.
- `compose.yml` is the shared service definition. It is not Docker-specific.
- `compose.override.yml` is used during development to expose frontend, backend, and Postgres ports on the host.

We do not run Docker and Podman together for the same environment. Use one runtime per machine.

The current project convention is:

- Local development defaults to Docker Compose through `.env.example`:
  ```env
  CONTAINER_TOOL=docker
  CONTAINER_ORCHESTRATION_TOOL=docker-compose
  ```
- Production deployment uses Podman Compose. The GitHub Actions deploy workflow runs:
  ```bash
  podman compose down --remove-orphans
  podman compose up -d --build --force-recreate
  ```

If you want to use Podman locally, change your local `.env` to:

```env
CONTAINER_TOOL=podman
CONTAINER_ORCHESTRATION_TOOL="podman compose"
```

The `justfile` reads these variables, so commands like `just up`, `just down`, `just logs backend`, and `just serve postgres` work with whichever runtime you configured.

For this project you need:
  - `just` the command runner
  - A Compose-compatible container setup:
    - Docker + `docker-compose` / `docker compose`, or
    - Podman + `podman compose`
  - `go` if you develop on the backend
  - `nodejs` if you develop on the frontend
Make sure they are installed on your system.

1. Clone the repository and navigate to the project root:
```bash
git clone <repo_url>
cd <repo_name>
```

2. Set up your environment variables:
  - You can create your own .env file, or
  - Use the provided example setup and execute:
  ```bash
  just init-env-prod
  ```

3. (Optional) Reset your database and seed it with example values.
  - Make sure your postgres container is running
  ```bash
  just serve postgres
  just reset-db
  just schema-db
  just seed-db
  ```

4. Run the app with:
  ```bash
  just prod
  ```

### Trust the local HTTPS certificate

The local Caddy container generates its own certificate authority for HTTPS. Because a container cannot add that authority to the host trust store, Chrome may reject the certificate and prevent secure features such as PWA service-worker registration.

Export Caddy's local root certificate to your home directory. Replace `username` with your Linux username:

```bash
docker compose cp \
  caddy:/data/caddy/pki/authorities/local/root.crt \
  /home/username/caddy-local-root.crt
```

Import the certificate into Chrome without administrator privileges:

1. Open `chrome://certificate-manager`.
2. Go to **Custom → Trusted certificates**.
3. Select **Import**.
4. Choose `/home/username/caddy-local-root.crt`.
5. Enable trust for identifying websites if Chrome asks for a trust purpose.
6. Completely close Chrome and open it again.

Then visit:

```text
https://localhost:7443
```

The page and `https://localhost:7443/sw.js` should load without a certificate warning. If a service-worker installation previously failed, open **Developer Tools → Application → Service Workers**, unregister the failed worker, clear the site's storage, and reload the page.

Only trust the root certificate exported from your own local Caddy container. A trusted root certificate can identify HTTPS websites in the browser profile where it is installed.

## Contributing

To contribute, please create a feature or fix branch from the latest `main` branch and open a pull request.

### Steps

1. **Make sure your local `main` is up to date**
```bash
git checkout main
git pull origin main
```

2. **Create a new branch**  
Use a clear and descriptive name:
```bash
git checkout -b docs/update-readme
```

3. **Make your changes and stage them**
```bash
git add .
```

4. **Commit your changes**  
Follow a clear commit message convention:
```bash
git commit -m "docs: update README"
```

5. **Push your branch**
```bash
git push -u origin docs/update-readme
```

6. **Open a pull request**
```bash
gh pr create
```

### Tips

- Keep branches focused on a single change  
- Use meaningful commit messages (`feat:`, `fix:`, `docs:`, etc.)  
- Rebase or pull latest `main` if your branch gets outdated

## Implementation

### Authentification

#### General Principle

- We are using JSON Web Tokens (JWTs) and sliding sessions for session management.
- When a user successfully logs in the server will respond by sending a JWT saved within a cookie
- The JWT encodes information about the user (claims), a creation time and an expiration time
- The token is encrypted with a key that is only known by the server
- When a client requests a protected part or functionality of the website the request has to contain this encrypted token
- The server then ensures the following:
    1. If the token is invalid, i.e. not present, not decryptable with the server's key or expired the request gets rejected.
    2. If the token is valid the request gets forwarded to a request handler and an updated token with extended expiration time gets created.
    3. The request handler checks if the token contained the correct user information for accessing the rquested feature, the request gets accepted or rejected accordingly. The response contains the updated token.

#### JWT Format

A JWT consists of three '.' separated Base64url-encoded JSON strings
1. Header: contains metadate (token type, cryptographic algorithm etc.)
2. Payload: claims
3. Signature: used for validation

#### The Cookie

- The server transmits the JWT by instructing the client (browser) to save it within a cookie.
- For this the HTTP response is send with the "Set-Cookie" header set.
- The following fields are set:
    - auth_token=<the acutual JWT>
    - Path=/api                     -> tells the browser for which requested resources it needs to send the cookie
    - Expires=<time>                -> tells the browser till when the cookie is valid
    - HttpOnly                      -> tells the browser to prevent JavaScript from accessing the cookie
    - Secure                        -> tells the browser to only send the cookie when communicating over HTTPS
    - SameSite=Strict               -> tells the browser to only send the cookie when sending a request from the same website

**example header:**
Set-Cookie: auth_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzk3MDEyOTMsImlhdCI6MTc3OTY5OTQ5Mywic3ViIjoiMSJ9.HEOudYQUwWYyZEy5cQOzFSmd0zioRYn-8LQR37hyiqI; Path=/api; Expires=Mon, 25 May 2026 09:28:13 GMT; HttpOnly; Secure; SameSite=Strict

#### Resources

https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Cookies
https://www.newline.co/@kchan/integrating-jwt-authentication-with-go-and-chi-jwtauth-middleware--ff9a6cec
https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-token-structure
