
# Contribution — shutan

## Detailed breakdown

I worked mainly on the frontend, with a focus on account self-service, mandatory legal pages, and the overall visual design of the authenticated app.

My work covered:

1. Building the **Settings** page end-to-end on the frontend  
2. Implementing **Privacy Policy** and **Terms of Service** pages (mandatory project requirement)  
3. Leading a **UI redesign** (theme tokens, wallpaper, navigation, page typography, button consistency)  
4. Adding a **user-switchable UI theme** (Aurora / Classic)  
5. Documenting chosen modules and ownership in `docs/MODULES.md`

Backend account APIs (`/api/me`, password update, delete account) were implemented by teammates; I integrated them on the frontend and handled validation, UX states, and post-action auth refresh/logout.

---

## Specific features, modules, and components

**Settings page (main feature)**  
- Route: `/settings` under the protected layout  
- Components:
  - `frontend/app/(protected)/settings/page.tsx`
  - `frontend/components/features/settings/settings.tsx`
  - `EmailSettingsForm.tsx` — update email via `PATCH /api/me`
  - `ChangePasswordForm.tsx` — update password via `PATCH /api/me/password`
  - `DeleteAccountSection.tsx` — delete account via `DELETE /api/me`, then logout
  - `SettingsPanel.tsx` — shared section layout (including danger-zone styling)
- Client validation with **react-hook-form** + **Zod**
- Uses `useAuth()` for loading state, current user data, `refreshUser()`, and `logout()`

**Privacy Policy & Terms of Service**  
- Pages: `/privacy`, `/terms`  
- Accessible from the footer  
- Later restyled to match the global teal / wallpaper theme

**UI / design system work**  
- Updated global tokens and fonts in `globals.css` and root layout  
- Added full-bleed background wallpaper under `frontend/public/images/`  
- Aligned Navigation, Footer, Home, My Events, About, Privacy, and Terms styling  
- Matched action buttons (e.g. Register / View Detail, Load more)

**Theme switcher (Aurora / Classic)**  
- Toggle in the navigation bar (`ThemeToggle`)  
- `data-theme` on `<html>` + CSS variables in `globals.css`  
- Preference persisted in `localStorage` (`camaraderie-theme`)  
- Key files: `lib/theme.ts`, `lib/context/ThemeContext.tsx`, `components/layout/ThemeToggle.tsx`  
- Pages use semantic classes (`text-chrome-title`, `text-chrome-body`, …) so both themes share one layout

**Documentation**  
- Contributed module documentation: chosen modules, point calculation, justification, implementation notes, and ownership (including Settings under user-management work)

Related subject areas: mandatory user management / account self-service UI; supporting work for Web (SSR-capable Next.js app, design consistency). Settings alone does not claim the full Major “Standard user management” module (friends / avatar / public profile remain out of scope for my part).

---

## Challenges and how they were overcome

1. **Integrating Settings with existing auth**  
   Profile updates must refresh global auth state so the navbar and other pages stay in sync.  
   **Solution:** call `refreshUser()` after successful profile updates; on account deletion, call `logout()` and rely on cookie-based JWT session clearing from the backend.

2. **Keeping UI readable on a dark wallpaper**  
   Light default text/cards became hard to read or looked inconsistent.  
   **Solution:** semantic chrome color tokens, opaque white event cards, and theme-aware nav/footer.

3. **Local backend not reachable during frontend development**  
   Missing `ADMIN_PASSWORD` caused the backend container to crash-loop, which looked like a frontend “offline” / login failure.  
   **Solution:** diagnose via container logs, ensure required env vars are passed through Compose, and verify API health before debugging client auth flows.

---

## Resources used

- Next.js App Router docs (layouts, client components, protected routes)  
- React docs (Context / hooks)  
- react-hook-form and Zod documentation  
- Tailwind CSS v4 theme / utility docs (`@theme` in CSS)  
- Project internal docs: `README.md`, architecture notes, OpenAPI/Scalar API docs from the backend  
- MDN: cookies, `fetch`, form validation concepts  
- Subject PDF (`ft_transcendence`) for README / module documentation requirements

---

## AI usage

I used AI assistants as a productivity aid for:

- Exploring the codebase and locating the right files for Settings, auth, and styling  
- Drafting UI token ideas and Tailwind class combinations aligned with an existing design system  
- Helping structure README / module documentation and the theme-switch approach  
- Debugging hydration and offline false-positive issues
