# Features Implemented by alneuma

## User CRUD & API Foundation

### User Reading and Creation
Implemented user CRUD endpoints utilizing the a layered repository/service/handler pattern. Including validation of username and password format rules.

### Basic User Authentication
Full JWT-based authentication system with login, register, logout endpoints, using the sliding sessions principle. Implemented password hashing using argon2id (128MB memory, 4 iterations, parallelism=NumCPU), cookie-based session management (HttpOnly, Secure, SameSite=Strict), and input validation using ozzo-validation. Created middleware for token verification and automatic cookie refresh.

### /api/me Endpoints
Added authenticated user self-service endpoints: `GET /api/me`, `PATCH /api/me`, `PATCH /api/me/password`, `DELETE /api/me`. Handlers extract user ID from JWT claims rather than URL path parameters.

## Event User Roles & Participant Management

### Join Event Endpoint
Added `POST /api/me/join/{id}` endpoint allowing authenticated users to join events as participants. Integrated event service dependency and JWT claim extraction for user identification.

### User's Events List
Added `GET /api/me/events` endpoint returning paginated list of events the authenticated user is registered for. Implemented SQL JOIN on event_participants table with limit/offset pagination.

### Event User Roles
Major refactor replacing implicit many-to-many GORM association with explicit `event_users` join table including a Role field (admin/member). Implemented cascade-delete hooks, role-based participant queries using LEFT JOIN with COALESCE, and converted all IDs from string to uint throughout the event package.

### Leave Event Endpoint
Added `POST /api/me/leave/{id}` endpoint allowing users to remove themselves from events. Implemented JWT-based user identification and participant removal.

### Auto-Add Creator as Admin
When creating an event, the authenticated user is automatically added as an admin participant using database transactions. Added role support to add-participant endpoint with validation for "admin" or "member" roles.

### By endpoint access control
Implemented admin only access for endpoints that modify events in any way.

## Event Features

### Event Image Endpoints
Full CRUD for event images (create, get, update, delete) with file storage on disk and image path tracking in database. Implemented MIME type detection, validation (max 100MB, JPEG/PNG only), and atomic file operations. Integrated mimetype library for content-type verification.

## Error Handling & Auth Infrastructure

### Custom Error Type System
Introduced `CamaError` type carrying both error kind (for programmatic matching) and human-readable message. Implemented `Is()` interface for `errors.Is()` pattern matching. Created error kind enum and constructor for consistent error creation across the application.

# Challenges

The most persistent and recurring challenges had to do with misunderstandings about HUMA, the backend framework, and GORM, the ORM. Those were solved by a mix of reading documentation and consulting AI.

# Resources

https://go.dev/tour/list
https://go-chi.io/#/README
https://huma.rocks/
https://gorm.io/

Go Programming Language, The (Addison-Wesley Professional Computing Series) (ISBN-10: 9780134190440)
Learning Go: An Idiomatic Approach to Real-World Go Programming (ISBN-10: 1492077216)

https://pkg.go.dev/github.com/alexedwards/argon2id
https://pkg.go.dev/github.com/danielgtaylor/huma/v2
https://pkg.go.dev/github.com/go-chi/chi/v5
https://pkg.go.dev/github.com/go-chi/jwtauth
https://pkg.go.dev/github.com/go-ozzo/ozzo-validation

# AI usage

- research best practices
- discuss architecture
- ask for common patterns to solve specific problems
- research language features or libraries to solve specific problems
- explain unfamiliar concepts
- assist with documentation
