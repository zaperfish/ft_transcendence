# ft_transcendence
ft_transcendence is the final group project of the 42 Common Core curriculum. It is a collaborative full-stack web application designed to simulate a real-world production system.

## Navigation

- [About This Project](#about-this-project)
- [Team Roles](#team-roles)
- [Architecture](#architecture)
- [Setup] (#setup)
- [Contributing](#contributing)
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

## Setup

1. Clone the repository and navigate to the project root:
```bash
git clone <repo_url>
cd <repo_name>
```

2. Set up your environment variables:
  - You can create your own .env file, or
  - Use the provided example setup and execute:
  ```bash
  just init-env
  ```

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
