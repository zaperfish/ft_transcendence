# Modules

Subject rule: **Major = 2 pts**, **Minor = 1 pt**. Target for validation: **14 pts**.

## Point calculation

| # | Module | Type | Points | Status |
|---|--------|------|--------|--------|
| 1 | Use a framework for both frontend and backend | Major | 2 | ✅ |
| 2 | Real-time features (WebSockets) | Major | 2 | ✅ |
| 3 | Public API (API key, rate limiting, ≥5 endpoints, docs) | Major | 2 | ✅ |
| 4 | Monitoring system (Prometheus + Grafana) | Major | 2 | ✅ |
| 5 | Advanced permissions system (roles & access control) | Major | 2 | ✅ |
| 6 | Use an ORM for the database | Minor | 1 | ✅ |
| 7 | Server-Side Rendering (SSR) | Minor | 1 | ✅ |
| 8 | Progressive Web App (PWA) | Minor | 1 | ✅ |
| | **Total (claimable)** | | **15** | |

Additional work in progress / not counted toward the claimable total yet:

| Module | Type | Points | Status | Notes |
|--------|------|--------|--------|-------|
| Standard user management and authentication | Major | 2 | 🟡 Partial | Login/register + **Settings page** done; friends system / avatar upload / dedicated profile page still missing for full Major claim |
| User interaction (chat, profiles, friends) | Major | 2 | 🟡 Partial | Event chat exists; friends system missing |
| File upload and management system | Minor | 1 | ✅ Extra | Event cover image upload/get/update/delete with validation |
| Advanced search (filters, pagination) | Minor | 1 | 🟡 Partial | Pagination + membership filters; no full multi-field search/sort |
| CI/CD deployment pipeline | Minor (Modules of choice) | 1 | ✅ Extra | Tag → SSH → Podman Compose deploy |

## List of chosen modules (Major and Minor)

1. **Major** — Use a framework for both frontend and backend  
2. **Major** — Real-time features using WebSockets  
3. **Major** — Public API with secured API key, rate limiting, documentation, ≥5 endpoints  
4. **Major** — Monitoring system with Prometheus and Grafana  
5. **Major** — Advanced permissions system (roles & access control)  
6. **Minor** — Use an ORM for the database  
7. **Minor** — Server-Side Rendering (SSR)  
8. **Minor** — Progressive Web App (PWA) with offline support and installability  

## Justification, implementation, and ownership

### 1. Major — Framework for frontend and backend (2 pts)

**Why:** A Meetup-style product needs structured routing, forms, and API conventions on both sides. Frameworks keep the team productive and the codebase consistent.

**How:**
- Frontend: **Next.js** (App Router) + React + Tailwind CSS
- Backend: **Go** with **Chi** + **Huma** (OpenAPI-first routing and docs)
- Reverse proxy: **Caddy**; DB: **Postgres** via Compose

**Team members:** lmiehler (architecture / bootstrap), yingzhan (frontend setup & pages), alneumann (backend structure)

---

### 2. Major — Real-time features / WebSockets (2 pts)

**Why:** Event attendees need live discussion while an event is active. WebSockets fit multi-user chat better than polling.

**How:**
- Backend: `coder/websocket` hub in `backend/chat/` (connect/disconnect handling, broadcast, persisted history)
- Frontend: `EventChatRoom` connects to the event chat WebSocket and renders live messages
- Route: event chat under `/events/[id]/chat`

**Team members:** dlemaire

---

### 3. Major — Public API (2 pts)

**Why:** External clients should be able to integrate with the event platform securely, without using browser session cookies.

**How:**
- API keys with dedicated middleware (`backend/apikey/`)
- Versioned public routes under `/api/v1/*` (CRUD-style event endpoints)
- Rate limiting middleware
- OpenAPI docs via Huma / Scalar

**Team members:** lmiehler

---

### 4. Major — Monitoring (Prometheus + Grafana) (2 pts)

**Why:** Production deployment on a VPS needs visibility into runtime health, latency, and resource pressure.

**How:**
- Prometheus scrapes backend `/metrics` (`prometheus/client_golang`)
- Grafana dashboards + alerting rules under `monitoring/`
- Compose services for Prometheus and Grafana; access hardened for prod

**Team members:** dlemaire

---

### 5. Major — Advanced permissions system (2 pts)

**Why:** Event organizers and attendees need different capabilities. Role-based access keeps event management safe (only admins can change sensitive event state or remove participants).

**How:**
- Event membership join table (`event_users`) with roles: `admin` and `member`
- Event creator is automatically assigned `admin` on create
- Join flow assigns `member`; leave/remove participant respects role rules
- Selected event endpoints are admin-only (update/delete event, manage participants, image management, etc.)
- Event listings can be filtered by the caller’s role (`admin` / `member`)

**Team members:** alneumann, lmiehler

---

### 6. Minor — ORM (1 pt)

**Why:** Avoid hand-written SQL for every model while keeping schema migrations and relations maintainable.

**How:**
- **GORM** (`gorm.io/gorm`) with Postgres driver
- Models + repositories for users, events, chat messages, etc.
- `AutoMigrate` on startup

**Team members:** alneumann, lmiehler

---

### 7. Minor — Server-Side Rendering / SSR (1 pt)

**Why:** Faster first paint and better SEO for public/marketing-facing routes; also used where the Next.js server can fetch API data before sending HTML.

**How:**
- Next.js App Router with server components where appropriate (e.g. chat-related server pages)
- During SSR, the Next.js server can call the backend over the internal network (see Architecture)

**Team members:** lmiehler, yingzhan

---

### 8. Minor — Progressive Web App / PWA (1 pt)

**Why:** Users should still browse cached event pages and get a clearer offline experience when connectivity drops.

**How:**
- `@ducanh2912/next-pwa` + service worker caching rules
- Offline banner, localStorage user cache, selected route/API caching
- Installability via web app manifest

**Team members:** yingzhan

---

## Related contribution — Settings page (**shutan**)

> **Owner: shutan**  
> This is **shutan**’s main frontend feature delivery. It supports account self-service under the broader user-management work (mandatory auth + partial Major “Standard user management”).

**What was implemented:**
- Route: `/settings` — `frontend/app/(protected)/settings/page.tsx`
- Feature UI: `frontend/components/features/settings/`
  - **Profile / email** — update email via `PATCH /api/me` (`EmailSettingsForm`)
  - **Password** — change password via `PATCH /api/me/password` (`ChangePasswordForm`)
  - **Danger zone** — permanently delete account via `DELETE /api/me` (`DeleteAccountSection`), then logout
- Navigation entry to Settings from the user menu
- Form validation with `react-hook-form` + Zod; loading / error / dirty states

**Backend counterparts used by Settings** (implemented mainly by alneumann): `/api/me`, `/api/me/password`, `DELETE /api/me`.

**Status vs Major “Standard user management”:** Settings covers profile update + secure password change + account deletion. Full Major still needs friends, avatar upload, and a dedicated public profile page before claiming the 2 pts.
