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

