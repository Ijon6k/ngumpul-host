# Frontend Architecture & Technical Standards

> **Framework:** SvelteKit 2 running in Node.js mode (`@sveltejs/adapter-node`)  
> **Reactivity & Syntax:** Svelte 5  
> **Styling:** Tailwind CSS v4 + Semantic CSS Variables  
> **Data Fetching:** TanStack Svelte Query v5 (`@tanstack/svelte-query`)  
> **Form Handling:** Superforms (`sveltekit-superforms`) + Valibot (`valibot`) via server actions  
> **Icons:** Phosphor Icons (`phosphor-svelte`)  
> **Location:** [`docs/frontend/architecture.md`](./architecture.md)

---

## 1. Directory & Codebase Organization

The frontend codebase is strictly organized into clean architectural layers:

```
frontend/src/
├── app.css                       # Global design tokens, theme variables & resets
├── lib/
│   ├── api/                      # Modular domain HTTP clients (projects, comments, reports, notifications, hostingRequests, activity, system, auth, admin/*)
│   ├── components/
│   │   ├── layout/               # Shell architecture (Header, Footer, ThemeToggle)
│   │   ├── ui/                   # Reusable atomic design system primitives
│   │   ├── bento/                # Host hardware telemetry & silicon visual cards
│   │   ├── status/               # Availability heatmap, indicators & probes
│   │   └── admin/                # Operator console shell & management views
│   ├── stores/                   # Auth store & reactive session state
│   ├── types/                    # Shared TypeScript domain contracts (project, user, comment, report, hosting, notification, activity)
│   ├── schemas/                  # Shared Valibot form schemas (auth.ts: login/register/setup)
│   ├── server/                   # Server-only helpers (session cookie forwarding)
│   └── utils/                    # Date & duration formatting utilities
└── routes/
    ├── +layout.svelte            # Root layout, TanStack Query provider & auth boot
    ├── +page.svelte              # Editorial landing page & hardware bento showcase
    ├── projects/                 # Project directory & [slug] deep-dive view
    ├── people/                   # Member directory & [username] profile view
    ├── activity/                 # Chronological public event feed
    ├── status/                   # Infrastructure availability & telemetry heatmap
    ├── me/                       # Authenticated member workspace & hosting intake
    ├── admin/                    # Operator console shell & management views
    ├── login/                    # Member authentication (superform + server action)
    ├── register/                 # Account registration (superform + server action)
    ├── setup/                    # Initial node provisioning (superform + server action)
    └── invite/[token]/           # Invitation acceptance flow
```

---

## 2. Layout & Spacing Mechanics

### 2.1. Fixed Header Collision Clearance
The global header is fixed at the top with a height of `60px`. To prevent content overlap on non-home pages without adding ad-hoc padding to individual components:

In `frontend/src/routes/+layout.svelte`:
```svelte
<main class="flex-1 shrink-0 {isHome ? '' : 'pt-[60px]'}">
    <slot />
</main>
```

### 2.2. Harmonized Page Padding
All subpages (`/projects`, `/people`, `/activity`, `/status`, `/me`, `/projects/[slug]`) must use the standardized page container padding token:
```html
<div class="container mx-auto px-4 sm:px-6 max-w-5xl py-8 sm:py-12 pb-24 flex flex-col gap-8">
    <!-- Page Content -->
</div>
```
Alternatively, import the design system primitive:
```svelte
<script>
    import { PageContainer, PageHeader } from '$lib/components/ui';
</script>

<PageContainer maxWidth="5xl">
    <PageHeader title="Projects" description="Independent applications hosted here." />
    <!-- Page Content -->
</PageContainer>
```

---

## 3. Data Fetching & State Architecture

### 3.1. Domain-Oriented Modular API Layer (`src/lib/api/`)
All backend interaction is centralized into typed domain modules rather than scattered raw `fetch()` or ad-hoc HTTP calls in components:
- **`projectsApi`** (`$lib/api/projects`): Public project catalog, slug detail, member projects, visit analytics, cover upload.
- **`commentsApi`** (`$lib/api/comments`): Public comments, comment posting, comment moderation and deletion.
- **`reportsApi`** (`$lib/api/reports`): Moderation reporting against projects or comments.
- **`notificationsApi`** (`$lib/api/notifications`): In-app member notifications, individual and bulk read acknowledgments.
- **`hostingRequestsApi`** (`$lib/api/hostingRequests`): Member application submission and personal request history.
- **`activityApi`** (`$lib/api/activity`): Public system activity feed and member personal activity history.
- **`systemApi`** (`$lib/api/system`): Public status availability, host specs, and service health.
- **`authApi`** (`$lib/api/auth`): Login, registration, session checks, profile updates, avatar uploads, invite validation.
- **`usersApi`** (`$lib/api/users`): Public member directory roster and user profile pages.
- **`adminApi`** (`$lib/api/admin`): Operator console subsystems (`stats`, `projects`, `comments`, `reports`, `users`, `requests`, `settings`, `system`, `audit`).

### 3.2. Shared Domain Contracts (`src/lib/types/`)
Domain models are strictly typed and centralized:
- `project.ts`: `Project`, `ProjectAvailability`, `ProjectVisits`, `UpdateProjectInput`.
- `user.ts`: `User`, `PublicUser`, `UpdateProfileInput`.
- `comment.ts`: `Comment`, `CreateCommentInput`.
- `report.ts`: `Report`, `CreateReportInput`.
- `hosting.ts`: `HostingRequest`, `CreateHostingRequestInput`.
- `notification.ts`: `AppNotification`.
- `activity.ts`: `ActivityEvent`.

### 3.3. Universal SvelteKit Load & SSR Proxy Hook (`src/hooks.server.ts`)
To allow SvelteKit universal page data preloading (`+page.ts`) while running behind Docker internal networks:
- `src/hooks.server.ts` implements `handleFetch` to automatically rewrite internal server-side requests targeting `/api/...` to the Go backend container (`BACKEND_URL` / `http://backend:8080`).
- Prevents container network isolation failures during server-side rendering while keeping public browser requests transparently routed via Nginx.
- Server actions that call `fetch('/api/...')` (e.g. auth form submissions) are routed through the same hook, so they transparently reach the Go backend.

### 3.4. Form Handling — Superforms + Server Actions (Auth Routes)
Login, registration, and initial node setup use `sveltekit-superforms` with shared Valibot schemas:
- **Shared schemas** live in `src/lib/schemas/auth.ts` (`loginSchema`, `registerSchema`, `setupSchema`) and are consumed both server-side (validation) and client-side (`valibotClient` for inline constraints/errors).
- **Server actions** (`+page.server.ts` on `/login`, `/register`, `/setup`) validate the request with `superValidate(request, valibot(schema))`, proxy the payload to the Go backend via `event.fetch('/api/...')`, and return `{ form, user }` on success or `fail(status, { form, message })` / `setError(...)` on failure.
- **Session cookie forwarding:** the Go backend issues the `ngumpul_session` cookie on its response. `src/lib/server/session.ts` exposes `forwardSessionCookie(cookies, response)` which re-sets that cookie on the browser via `cookies.set()`, honoring `Secure`/`SameSite`/`HttpOnly`/`expires`/`max-age` attributes.
- **Cross-field validation** (`passwords match`) is intentionally performed in the server action with `setError(form, 'confirmPassword', ...)` rather than a schema-level `v.check`, keeping pathless issues out of the `_errors` bucket.
- On successful submissions the client's `onResult` handler updates the `$user` store and redirects by role.

### 3.5. Reactive Stores & Real-Time Counter Synchronization
- **Session Authentication (`$lib/stores/auth.ts`)**: Initializes via `initAuth()` in `+layout.svelte`, storing `$user`.
- **Member Notifications (`$lib/stores/notifications.ts`)**: `unreadNotificationsCount` store keeps the `/me` sidebar badge synchronized in real-time when notifications are opened or marked as read.
- **Operator Metrics (`$lib/stores/adminStats.ts`)**: `adminPendingCount` and `adminOpenReportsCount` keep operator badges updated when requests or reports are reviewed.

---

## 4. Phosphor Icons Standard

All icons across the application must be imported from `phosphor-svelte`:
```svelte
<script lang="ts">
    import { Cpu, ShieldCheck, ArrowRight, ArrowUpRight } from 'phosphor-svelte';
</script>

<!-- Regular weight for neutral icons -->
<ShieldCheck size={16} weight="regular" class="text-sky-500" />

<!-- Bold weight for directional arrows and confirmed checks -->
<ArrowRight size={14} weight="bold" />
```

### Rules:
- **No Raw SVGs:** Do not paste raw inline `<svg>` elements or raw SVG strings in navigation arrays.
- **No Fill / Heavy Icons:** Use `weight="regular"` (or `"bold"` for arrows and checks).
- **Interactive Buttons:** The navbar `ThemeToggle` button must remain borderless with a transparent hover background.

---

## 5. Project Directory & Showcase Standards

### 5.1. Public Project Catalog (`/projects`)
- **Grid Structure:** Strict 2-column desktop (`grid-cols-1 md:grid-cols-2 gap-6 sm:gap-8`) and 1-column mobile layout. Eliminates crowded SaaS card rows.
- **Project Cover Primitives:** 4:3 aspect ratio (`aspect-[4/3]`) utilizing `<ProjectCover>` with lazy loading and deterministic color/letter fallback when an image has not been uploaded.
- **Pagination:** Controlled pagination (`<Pagination>`) bound to 12 projects per page with query filtering (`type`, `q`).
- **Typographic Hierarchy:** Creator, status, and concise technology stack separated by middle dots (`·`) without pill/chip spam.

### 5.2. Project Detail Showcase (`/projects/[slug]`)
- **Editorial Showcase Hierarchy:** Breadcrumb $\rightarrow$ Title $\rightarrow$ Short Description $\rightarrow$ Status $\rightarrow$ 16:9 Showcase Artwork $\rightarrow$ Asymmetric 2-column body (About, Technology, Project Activity timeline on left; Access Links, Creator Card, Node Privacy notice on right).
- **Zero Telemetry Leaks:** No container IDs, internal ports, cgroups, memory limits, TLS implementation internals, or fake ping animations.
- **Real Project Activity:** Honest, project-specific activity stream queried directly from recorded events (`activities` table with `project_id`).
