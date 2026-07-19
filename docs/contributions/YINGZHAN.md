### Main modules

- Use a frontend framework
- Frontend of advanced permissions system (roles & access control)
- Server-Side Rendering (SSR)
- Progressive Web App (PWA)

### Main features

- User Authentication
- User Registration
- Event Creation (CRUD)
- Event Discovery
- Event Participation
- Server-Side Rendering
- Progressive Web App (PWA)
- Image Uploading

### Main components

- Authentication flow: login and register pages, protected routing, password confirmation, password-policy alignment, and auth state handling.
- Navigation and layout: navigation layout, mobile dropdown navigation, footer behavior, protected layout adjustments, and global style conflict fixes.
- Event browsing: homepage event feed, event list page, filtering, pagination, and the event detail page.
- Event interaction: create event form, edit event modal integration, delete-event support, and supporting UI state updates.
- Media handling: image upload UI, upload API wiring, preview handling, homepage/event-card/detail-page image display, and cover-image fallback logic.
- Offline and PWA support: PWA setup, service-worker caching rules, offline banner, offline HTML cache, offline logic for home, auth, detail, my events, footer pages, logout, and request handling.
- UI refinement: toast replacement for alerts, form component encapsulation, reusable UI pieces, and visual fixes for larger screens and broken image states.

### Challenges and how they were handled

- Offline caching introduced 401 and network-console noise on login, register, and protected pages. The fix was to tighten cache rules, add offline-aware request handling, and silence expected failures in the affected flows.
- Layout and style conflicts between the root layout, protected layout, and auth pages caused display regressions. These were addressed by isolating layout responsibilities and adjusting the global styling structure.
- Image rendering and preview behavior differed across pages and screen sizes. This was corrected by refining image components, adding fallback covers, and fixing the responsive presentation of large images.

### AI Usage

- Assist visualizing product prototype
- Assist making UI design
- Resolve tailwindCSS classnames
- Ask for best practices
- Ask for debugging frontend codes, especially display errors
- Assist writing documentation

### Resources

- https://nextjs.org/learn/dashboard-app
- https://nextjs.org/learn/react-foundations
- https://react.dev/learn
- https://www.w3schools.com/react/default.asp
- Haverbeke, M. (2024). Eloquent JavaScript (4th ed.). No Starch Press.