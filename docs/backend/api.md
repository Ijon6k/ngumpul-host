# REST API Specification & Data Contracts

> **Base URL:** `/api` (canonical) and `/go/{slug}` (outbound tracking redirect)  
> **Content-Type:** `application/json`  
> **Authentication:** Stateful cookie `ngumpul_session` (`HttpOnly; SameSite=Lax`)  
> **Location:** [`docs/backend/api.md`](./api.md)

---

## 1. Response Envelope Protocol

### Standard Success Response
```json
{
  "data": { ... },
  "message": "Optional human-readable confirmation message"
}
```

### Standard Error Response
```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "A detailed explanation of the error",
    "details": null
  }
}
```

#### Common Error Codes
- `BAD_REQUEST`: Malformed JSON or invalid parameter validation.
- `UNAUTHORIZED`: Request missing a valid `ngumpul_session` cookie.
- `FORBIDDEN`: User does not possess the necessary role (e.g. non-admin accessing `/api/admin/*`).
- `NOT_FOUND`: Target resource (slug, username, or UUID) does not exist.
- `INTERNAL_SERVER_ERROR`: Unhandled backend failure.

---

## 2. Public Endpoints (Unauthenticated)

### 2.1. System Health
- **`GET /health`** or **`GET /api/health`**
  - **Response (200 OK):**
    ```json
    {
      "service": "ngumpul-host-backend",
      "status": "healthy",
      "time": "2026-09-11T20:00:00Z"
    }
    ```

### 2.2. Authentication & Access
- **`GET /api/auth/mode`**
  - **Response (200 OK):**
    ```json
    {
      "registration_mode": "INVITE_ONLY"
    }
    ```
- **`GET /api/invitations/validate?token={raw_token}`**
  - **Response (200 OK):**
    ```json
    {
      "valid": true,
      "invited_email": "friend@example.com",
      "expires_at": "2026-09-18T20:00:00Z"
    }
    ```
- **`POST /api/auth/register`**
  - **Request Body:**
    ```json
    {
      "username": "erik",
      "email": "erik@example.com",
      "password": "securepassword123",
      "display_name": "Erik Maulana",
      "invitation_token": "a1b2c3d4e5f6..."
    }
    ```
    *Note: `invitation_token` is required when `registration_mode` is `INVITE_ONLY`.*
  - **Response (201 Created):** Sets `ngumpul_session` cookie and returns user profile.

- **`POST /api/auth/login`**
  - **Request Body:**
    ```json
    {
      "email": "erik@example.com",
      "password": "securepassword123"
    }
    ```
  - **Response (200 OK):** Sets `ngumpul_session` cookie and returns user profile.

- **`POST /api/auth/logout`**
  - **Response (200 OK):** Invalidates session in database and clears `ngumpul_session` cookie.

### 2.3. Public Discovery & Catalog
- **`GET /api/projects`**
  - Optional Query Parameters: `page` (default `1`), `limit` (default `12`, max `50`), `q` (search term), `type` (`HOSTED_HERE` | `EXTERNAL`), `tech` (technology filter), `status` (`ONLINE` | `OFFLINE`).
  - **Response (200 OK):**
    ```json
    {
      "projects": [ ... ],
      "total": 24,
      "page": 1,
      "limit": 12,
      "total_pages": 2
    }
    ```

- **`GET /api/projects/{slug}`**
  - **Response (200 OK):**
    ```json
    {
      "project": {
        "id": "...",
        "name": "Atlas",
        "slug": "atlas",
        "owner": { ... },
        "status": "ONLINE"
      },
      "activities": [
        {
          "id": "...",
          "type": "PROJECT_PUBLISHED",
          "metadata": { "title": "Atlas published" },
          "created_at": "2026-09-10T12:00:00Z",
          "actor_name": "Erik Maulana"
        }
      ]
    }
    ```

- **`GET /api/projects/{slug}/comments`**
  - **Response (200 OK):**
    ```json
    {
      "comments": [
        {
          "id": "...",
          "content": "Impressive performance on bare metal.",
          "created_at": "2026-09-11T14:30:00Z",
          "author": {
            "id": "...",
            "username": "erik",
            "display_name": "Erik Maulana",
            "avatar_url": null
          },
          "is_deleted": false
        }
      ]
    }
    ```

- **`GET /api/users`**
  - **Response (200 OK):** `{ "members": [ ... ] }`

- **`GET /api/users/{username}`**
  - **Response (200 OK):** Returns member profile, bio, joined date, and owned public projects.

- **`GET /api/activity`**
  - **Response (200 OK):** Chronological activity feed for the public timeline.

### 2.4. Infrastructure Telemetry
- **`GET /api/status?range={1d|7d|30d}`**
  - **Response (200 OK):**
    ```json
    {
      "availability_percent": 100.0,
      "status": "operational",
      "uptime_formatted": "3 days, 4 hours",
      "observed_period": "Past 30 Days",
      "daily_blocks": [
        {
          "index": 0,
          "label": "Aug 12",
          "date": "2026-08-12",
          "status": "operational",
          "uptime_percent": 100.0,
          "details": "100.0% operational"
        }
      ],
      "incidents": []
    }
    ```

- **`GET /go/{slug}`**: Public outbound redirect tracking endpoint. Increments `project_visits` and redirects (HTTP 302) to the project's external application URL.

---

## 3. Authenticated Member Endpoints (`RequireAuth`)

Requests must include a valid `ngumpul_session` cookie.

- **`GET /api/me`**: Returns the currently authenticated user identity and role.
- **`PATCH /api/me`**: Updates `display_name`, `bio`, or `avatar_url`.
- **`GET /api/me/projects`**: Lists projects owned by the authenticated member.
- **`GET /api/me/projects/{id}`**: Detailed project view for project owner including availability probes and activity history.
- **`GET /api/me/projects/{id}/visits`**: Returns 30-day traffic analytics (`total_page_views`, `total_visits`, and daily breakdown `[{ date, page_views, visits }]`).
- **`PATCH /api/me/projects/{id}`**: Updates allowed fields on an owned project (description, cover image, documentation URL, repository, demo, tech stack).
- **`POST /api/hosting-requests`**: Submits a new hosting request.
- **`GET /api/me/hosting-requests`**: Returns the user's submitted requests and operator notes.
- **`GET /api/me/activity`**: Lists recent personal audit and project events.
- **`GET /api/me/notifications`**: Lists in-app member notifications.
- **`PATCH /api/me/notifications/{id}/read`**: Marks a notification as read.
- **`POST /api/me/notifications/read-all`**: Marks all member notifications as read.
- **`POST /api/projects/{slug}/comments`**: Posts a comment on a project showcase (rate limit: 5/10m per user).
- **`DELETE /api/comments/{id}`**: Soft-deletes a comment. Permitted for comment author, project owner, or operator.
- **`POST /api/projects/{slug}/report`**: Submits a moderation report against a project (rate limit: 5/10m per user/IP).
  - **Payload:** `{ "reason": "SPAM" | "ABUSE_HARASSMENT" | "INAPPROPRIATE" | "MALICIOUS_SUSPICIOUS" | "OTHER", "details": "..." }`
- **`POST /api/comments/{id}/report`**: Submits a moderation report against a comment (rate limit: 5/10m per user/IP).
  - **Payload:** `{ "reason": "SPAM" | "ABUSE_HARASSMENT" | "INAPPROPRIATE" | "MALICIOUS_SUSPICIOUS" | "OTHER", "details": "..." }`
- **`POST /api/upload`**: Multipart file upload (`multipart/form-data`, file key: `file`, optional form field: `purpose` (`avatar` | `cover`)). Hard size limit: 10MB. Allowed MIME types: `image/jpeg`, `image/png`, `image/webp` (SVGs rejected). Enforces dimension sanity checks ($\le 4096\text{px}$). Returns `{ "url": "/uploads/covers/uuid.webp", "object_key": "covers/uuid.webp", "byte_size": 123456, "content_type": "image/webp", "width": 1200, "height": 900 }`.

---

## 4. Administrative Endpoints (`RequireAdmin`)

Requests must possess `role == 'ADMIN'`.

- **`GET /api/admin/settings`**: Retrieves instance settings (e.g. `{ "registration_mode": "INVITE_ONLY" }`).
- **`PATCH /api/admin/settings`**: Updates instance settings (`{ "registration_mode": "OPEN" | "INVITE_ONLY" | "CLOSED" }`).
- **`GET /api/admin/invitations`**: Lists all generated invitations with usage counts, email restrictions, and expiration timestamps.
- **`POST /api/admin/invitations`**: Generates a new invitation token (`{ "invited_email": "friend@example.com", "max_uses": 1, "expires_in_days": 7 }`). Returns one-time view of `raw_token` and `invite_url`.
- **`POST /api/admin/invitations/{id}/revoke`**: Revokes an active invitation token immediately.
- **`GET /api/admin/stats`**: Aggregate counts for total members, active projects, pending requests, and `open_reports`.
- **`GET /api/admin/reports`**: Lists community moderation reports (`?status=OPEN|REVIEWED|RESOLVED|DISMISSED`).
- **`POST /api/admin/reports/{id}/resolve`**: Resolves a moderation report (`{ "resolution_notes": "..." }`).
- **`POST /api/admin/reports/{id}/dismiss`**: Dismisses a moderation report (`{ "resolution_notes": "..." }`).
- **`GET /api/admin/comments`**: Lists all project comments with report count, author, and soft-delete status. Supports `?search=...` and `?project_id=...`.
- **`DELETE /api/admin/comments/{id}`**: Soft-deletes a comment as operator.
- **`GET /api/admin/users`**: Lists all members with email, status, and role.
- **`PATCH /api/admin/users/{id}/role`**: Updates member role (`{ "role": "ADMIN" | "USER" }`). Protected against self-demotion and demoting the last active administrator on the server.
- **`PATCH /api/admin/users/{id}/status`**: Updates member account status (`{ "status": "ACTIVE" | "SUSPENDED" }`). Protected against self-suspension and suspending the last active administrator.
- **`GET /api/admin/hosting-requests`**: Lists incoming hosting requests.
- **`POST /api/admin/hosting-requests/{id}/approve`**: Approves request and generates project record.
- **`POST /api/admin/hosting-requests/{id}/reject`**: Rejects request (`{ "reason": "Repository inaccessible" }`).
- **`POST /api/admin/hosting-requests/{id}/complete`**: Finalizes container deployment, sets public URL, and marks project `ONLINE`.
- **`GET /api/admin/system`**: Detailed Linux host telemetry, mounts, and probe diagnostics.
- **`GET /api/admin/audit`**: Immutable operator audit log.

