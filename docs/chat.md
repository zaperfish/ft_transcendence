# Event Chat

**Contributor:** dlemaire

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
