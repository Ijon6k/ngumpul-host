# Operational Workflows & End-to-End Scenarios

> **Component:** End-to-End Scenarios & Lifecycle Sequences  
> **Target Audience:** Core Contributors, Operators & AI Agents

---

## 1. Member Onboarding & Hosting Request Workflow

This workflow documents a new user joining the community, submitting an application for hosting, and receiving confirmation:

```mermaid
sequenceDiagram
    autonumber
    actor Member as Member (User)
    participant Web as SvelteKit Frontend
    participant API as Go Backend
    participant DB as PostgreSQL Database
    actor Operator as Server Operator (Admin)

    Member->>Web: Navigate to /register or /invite/:token
    Web->>API: GET /api/auth/mode & GET /api/invitations/validate?token=:token
    API-->>Web: Returns registration_mode ('INVITE_ONLY') & token validity
    Member->>Web: Submit registration form with invitation token
    Web->>API: POST /api/auth/register (username, email, password, invitation_token)
    API->>DB: Atomic Tx: Validate & consume invitation (SELECT ... FOR UPDATE)
    API->>DB: INSERT INTO users (role='USER', status='ACTIVE')
    API-->>Web: Sets HTTP-only 'ngumpul_session' cookie & returns user profile
    Web-->>Member: Redirect to /me (Workspace)

    Member->>Web: Fill Hosting Request Form (/me/requests)
    Web->>API: POST /api/hosting-requests (project_name, repo_url, tech_stack, notes)
    API->>DB: INSERT INTO hosting_requests (requester_id, status='PENDING')
    API->>DB: INSERT INTO activities (type='REQUEST_SUBMITTED')
    API-->>Web: HTTP 201 Created
    Web-->>Member: Request displayed with 'PENDING' status badge

    Operator->>Web: Inspect Operator Console (/admin/requests)
    Web->>API: GET /api/admin/hosting-requests
    API->>DB: SELECT * FROM hosting_requests WHERE status='PENDING'
    API-->>Web: Return pending requests
    Web-->>Operator: Display request with repository and tech stack
```

---

## 2. Operator Review, Container Provisioning & Publication

This workflow documents how the operator approves the request, provisions the container on the host node, and publishes the project:

```mermaid
sequenceDiagram
    autonumber
    actor Operator as Server Operator (Admin)
    participant Web as SvelteKit Frontend
    participant API as Go Backend
    participant DB as PostgreSQL Database
    actor Member as Member (User)

    Operator->>Web: Click "Approve Request"
    Web->>API: POST /api/admin/hosting-requests/{id}/approve
    API->>DB: UPDATE hosting_requests SET status='APPROVED', reviewed_at=NOW()
    API->>DB: INSERT INTO projects (owner_id, status='SETUP', visibility='UNPUBLISHED')
    API->>DB: INSERT INTO notifications (user_id=Member, title="Hosting Request Approved")
    API-->>Web: HTTP 200 OK

    Note over Operator: Operator deploys container on host (e.g. docker run -p 3045:3000 myapp)
    Note over Operator: Configures reverse proxy ingress for custom subdomain (myapp.ngumpul.id)

    Operator->>Web: Click "Complete Deployment" (/admin/requests)
    Web->>API: POST /api/admin/hosting-requests/{id}/complete (public_url)
    API->>DB: UPDATE projects SET status='ONLINE', visibility='PUBLIC', public_url='https://myapp.ngumpul.id'
    API->>DB: UPDATE hosting_requests SET status='COMPLETED'
    API->>DB: INSERT INTO activities (type='PROJECT_PUBLISHED')
    API->>DB: INSERT INTO notifications (user_id=Member, title="Project Live!")
    API-->>Web: HTTP 200 OK

    Member->>Web: Visits Public Catalog (/projects)
    Web->>API: GET /api/projects
    API->>DB: SELECT * FROM projects WHERE visibility='PUBLIC'
    API-->>Web: Return published catalog
    Web-->>Member: Project is live with ping probe and edge telemetry
```

---

## 3. Host Reboot & Automated Downtime Gap Reconciliation

This workflow documents how the system detects host reboots without active heartbeat polling:

```mermaid
sequenceDiagram
    autonumber
    participant Kernel as Linux Kernel (/proc)
    participant Worker as Heartbeat Goroutine
    participant DB as PostgreSQL Database
    participant Service as Availability Engine

    Note over Kernel,DB: Server runs continuously. Heartbeat runs every 5 minutes.
    Worker->>DB: UPDATE availability_state SET last_seen_at = 10:00:00 (boot_id = 'UUID-A')

    Note over Kernel: Host reboots at 10:03:00 (Power cycle / Kernel upgrade)
    Note over Kernel: Linux initializes at 10:05:00 (boot_id = 'UUID-B', /proc/uptime = 60s)

    Service->>Kernel: Read boot_id ('UUID-B') and /proc/uptime (60s)
    Service->>DB: SELECT * FROM availability_state WHERE id = 1
    DB-->>Service: prev_boot_id = 'UUID-A', prev_last_seen = 10:00:00

    Service->>Service: Gap = 10:06:00 - 10:00:00 = 6m (>= 120s)
    Service->>Service: prev_boot_id != current_boot_id -> Cause = 'SYSTEM_REBOOT'
    Service->>Service: Incident Interval = [10:00:00, 10:04:00] (duration = 240s)

    Service->>DB: INSERT INTO availability_incidents (cause='SYSTEM_REBOOT', duration=240s)
    Service->>DB: UPDATE availability_state SET boot_id='UUID-B', started_at=10:04:00, last_seen_at=10:06:00

    Note over Service,DB: Next request to /api/status automatically reflects the 4-minute reboot gap in the heatmap.
```

---

## 4. Administrative Member Moderation & Role Safeguards

This workflow documents safety invariants and confirmation controls when managing members and roles (`/admin/users`):

```mermaid
sequenceDiagram
    autonumber
    actor Admin as Server Operator (Admin)
    participant Web as SvelteKit Frontend
    participant API as Go Backend
    participant DB as PostgreSQL Database

    alt Promote Member to Admin
        Admin->>Web: Click "Make Admin" on Member row
        Web-->>Admin: Open high-impact modal requiring dynamic confirmation phrase
        Note over Web: Prompts: 'give admin role to <username>'
        Admin->>Web: Type exact matching phrase
        Web->>API: PATCH /api/admin/users/:id/role { "role": "ADMIN" }
        API->>DB: UPDATE users SET role = 'ADMIN', updated_at = NOW()
        API->>DB: INSERT INTO audit_logs (action='ADMIN_CHANGED_ROLE')
        API-->>Web: HTTP 200 OK
        Web-->>Admin: Success notification & table updated
    else Demote or Suspend User
        Admin->>Web: Click "Demote" or "Suspend"
        Note over Web: If target is self (u.id === $user.id), action is disabled
        Web-->>Admin: Display explicit confirmation dialog with impact warning
        Admin->>Web: Confirm action
        Web->>API: PATCH /api/admin/users/:id/role OR /status
        API->>API: Verify target != current admin
        API->>API: Verify target is not the last remaining active admin
        API->>DB: UPDATE users SET role / status
        API->>DB: If SUSPENDED: DELETE FROM sessions WHERE user_id = :id
        API->>DB: INSERT INTO audit_logs
        API-->>Web: HTTP 200 OK
        Web-->>Admin: Table refreshed
    end
```

### Safety Invariants:
1. **Self-Action Prevention:** An administrator cannot demote or suspend their own account in either frontend UI or backend API.
2. **Last Administrator Protection:** The backend rejects any attempt to demote or suspend the last active administrator on the system, preventing instance lockout.
3. **GitHub-Style Dynamic Typed Phrase:** Promoting any member to administrator requires typing `give admin role to <username>`, preventing accidental privilege elevation.
4. **Interactive Confirmation Modals:** Demoting an administrator, suspending an account, or restoring access requires confirmation with transparent consequence descriptions.

---

## 5. Outbound Project Redirect Engine (`/go/:slug`)

This workflow documents how community visitors navigate from Ngumpul Host showcase pages to external deployed applications while securely and anonymously recording outbound visit metrics:

```mermaid
sequenceDiagram
    autonumber
    actor Visitor as Community Visitor
    participant Nginx as Edge Ingress (Nginx)
    participant API as Go Backend Core
    participant DB as PostgreSQL Database
    actor App as External Project Application

    Visitor->>Nginx: Click "Visit project ↗" (GET /go/:slug)
    Nginx->>API: Proxy pass to backend /go/:slug
    API->>DB: Query project record by slug (SELECT id, public_url, status, visibility)
    alt Project or Public URL Not Found
        DB-->>API: No rows returned / public_url is empty
        API-->>Visitor: HTTP 404 Not Found JSON error
    else Valid Live Project
        DB-->>API: Returns project metadata & target public_url
        Note over API: Non-blocking asynchronous click recording
        API-)DB: INSERT INTO project_visits (project_id, visited_at, date_bucket)
        API-->>Visitor: HTTP 302 Found (Location: public_url)
        Visitor->>App: Browser follows redirect to external service
    end
```

### Architectural & Privacy Guarantees:
1. **Asynchronous Non-Blocking Execution:** Click recording occurs in a lightweight background goroutine (`go func() { ... }()`) with a 5-second context timeout. Visitor redirection is never delayed by database write latency.
2. **Zero Client Fingerprinting:** Outbound visit records strictly store `(project_id, visited_at, CURRENT_DATE)`. The system intentionally does **NOT** store client IP addresses, browser cookies, Canvas fingerprints, or third-party ad beacons.
3. **Dual-Layer Ingress Reliability:**
   - **Primary Layer (Nginx):** Route `location /go/` forwards directly to `backend_upstream` with zero Node.js overhead.
   - **Fallback Layer (SvelteKit):** Universal endpoint in `src/routes/go/[slug]/+server.ts` handles SSR/direct edge fallback, guaranteeing redirection even if ingress routing is bypassed during local development.

---

## 6. Project Traffic Analytics, Impression Deduplication & 30-Day Aggregation Engine

Ngumpul Host provides project owners with calm, actionable visibility into their application's reach without corporate surveillance tooling:

```mermaid
sequenceDiagram
    autonumber
    actor Visitor as Community Visitor
    participant Svelte as SvelteKit Server / Client
    participant Cache as In-Memory sync.Map Cache
    participant API as Go Backend
    participant DB as PostgreSQL Database

    Note over Visitor,DB: 1. Showcase Page Impression Tracking
    Visitor->>Svelte: View /projects/:slug
    Svelte->>API: GET /api/projects/:slug
    API->>Cache: Check clientKey (IP hash / session) last_seen
    alt Seen within past 30 minutes
        Cache-->>API: Cache hit (within 30m window)
        Note over API: View ignored (deduplicated)
    else First view or expired window
        API->>Cache: Update last_seen = NOW()
        API->>DB: UPSERT INTO project_page_views (project_id, date, views_count + 1)
    end

    Note over Visitor,DB: 2. Owner Analytics Retrieval (/api/me/projects/:id/visits)
    actor Owner as Project Owner
    Owner->>API: GET /api/me/projects/:id/visits (Cookie: ngumpul_session)
    API->>API: Verify caller == project.owner_id OR caller.role == 'ADMIN'
    API->>DB: Query 30-day views (project_page_views WHERE date >= CURRENT_DATE - 29d)
    API->>DB: Query 30-day outbound clicks (project_visits WHERE date_bucket >= CURRENT_DATE - 29d)
    API->>API: Generate continuous 30-day timeline map (zero-fill missing dates)
    API-->>Owner: Return JSON { total_views_30d, total_outbound_30d, daily_breakdown }
```

### Technical Invariants:
1. **30-Minute Impression Deduplication:** `RecordPageView` uses a thread-safe `sync.Map` in the Go backend. Repeated refreshes, rapid clicks, or bot crawls from the same IP/session within a 30-minute window increment zero database counters.
2. **Continuous Zero-Filled 30-Day Timeline:** The `GetProjectVisits` handler initializes a 30-day date map (`dailyMap`) from `NOW() - 29 days` to `TODAY`. Missing days are preserved as `0 views` and `0 visits`, ensuring frontend line charts and bento metric widgets render smooth, gapless visualizations.
3. **Strict Authorization Gate:** Non-owners cannot inspect other members' traffic statistics. Access is restricted strictly to the project creator or instance operators.

---

## 7. Smart External Link & Platform Auto-Detection (`linkDetector.ts`)

Instead of forcing users into rigid, fragmented form inputs, Ngumpul Host allows creators to provide up to 5 arbitrary external links while automatically inferring the service, domain, brand iconography, and human labels:

```
User Input URL ───► URL Hostname Parser ───► Domain Matcher ───► Curated Phosphor Icon + Label
```

### Detection Matrix:

| Domain Match | Platform Identified | Phosphor Icon | Default Label Rendered |
| :--- | :--- | :--- | :--- |
| `github.com` | GitHub Repository | `<GithubLogo weight="regular" />` | `GitHub` |
| `gitlab.com` | GitLab Repository | `<GitlabLogo weight="regular" />` | `GitLab` |
| `drive.google.com` | Google Drive Cloud Storage | `<GoogleDriveLogo weight="regular" />` | `Google Drive` |
| `docs.google.com` | Google Docs Specification | `<FileText weight="regular" />` | `Google Docs` |
| `figma.com` | Figma Design Workspace | `<FigmaLogo weight="regular" />` | `Figma` |
| `youtube.com`, `youtu.be` | YouTube Video Demo | `<YoutubeLogo weight="regular" />` | `YouTube Demo` |
| `vimeo.com` | Vimeo Video Demo | `<YoutubeLogo weight="regular" />` | `Video Demo` |
| `twitter.com`, `x.com` | X / Twitter Profile | `<TwitterLogo weight="regular" />` | `X (Twitter)` |
| `discord.gg`, `discord.com` | Discord Community Server | `<DiscordLogo weight="regular" />` | `Discord` |
| `notion.so`, `notion.site` | Notion Documentation | `<FileText weight="regular" />` | `Notion Docs` |
| *Any other valid URL* | Generic External Website | `<LinkSimple weight="regular" />` | Clean hostname or fallback |

### Frontend Consumption:
```svelte
<script lang="ts">
    import { detectLinkInfo } from '$lib/utils/linkDetector';
    
    // Example: Dynamically detect external link
    const info = detectLinkInfo(project.demo_url, 'Live Demo');
    const IconComponent = info.icon;
</script>

<a href={project.demo_url} target="_blank" rel="noopener noreferrer" class="flex items-center gap-2 text-sm text-neutral-600 dark:text-neutral-300 hover:text-sky-500">
    <IconComponent size={16} />
    <span>{info.label}</span>
</a>
```


