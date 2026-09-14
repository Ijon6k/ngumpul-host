# Security & Hardening Reference

> **Component:** Whole stack (Nginx inbound, SvelteKit frontend, Go backend, PostgreSQL, SeaweedFS)
> **Status:** Active / Living Document
> **Location:** [`docs/security.md`](./security.md)

This document is the canonical reference for every security control enforced in Ngumpul Host. It exists so that operators and contributors know exactly **what is protected, how it is protected, and what is still open**. If you change a heading, a middleware, or a validation rule, update this document in the same task.

---

## 1. Tenancy & Authorization Model

Ngumpul Host is a **single-instance community server**. There is no multi-tenancy, no `tenant_id`, and no organizational partitioning. Authorization is purely **ownership-based + role-based**:

| Role | Can do | Enforced by |
| :--- | :--- | :--- |
| `Public (Anonymous)` | Read public projects, member directory, activity ledger, telemetry | — |
| `USER` | Authenticate, manage profile, submit hosting requests, edit own project metadata | `auth.RequireAuth` |
| `ADMIN` | Everything in `USER` plus `/api/admin/*`, approvals, role promotion, audit log | `auth.RequireAuth` + `auth.RequireRole("ADMIN")` |

Ownership checks inside handlers (not just middleware):
- Comment deletion → comment author, project owner, or operator.
- Project updates → owner or admin only (`403` otherwise).
- Subdomain changes → owner or admin only.
- Private projects → visible only to owner or admin.

Suspended users are rejected at the session layer: any request with a session belonging to a `SUSPENDED` user returns `403 Forbidden` immediately, and the user's sessions are deleted at suspension time.

---

## 2. Authentication & Session Management

### Session Token Generation (`internal/auth/session.go`)
- **32 random bytes** from `crypto/rand`, hex-encoded to a 64-character token.
- `crypto/rand` is a cryptographically secure source; tokens are unguessable.
- Tokens are stored server-side in PostgreSQL (`sessions` table) — not in a self-contained JWT — so they can be revoked instantly.

### Session Lifetime
- Default expiry: **30 days** from creation.
- Expired sessions are deleted on access (lazy cleanup). No periodic cleanup worker currently runs — stale rows accumulate until accessed.

### Cookie Settings (`ngumpul_session`)
| Attribute | Value | Rationale |
| :--- | :--- | :--- |
| `HttpOnly` | `true` | Prevents JavaScript access (XSS can't read the session). |
| `SameSite` | `Lax` | Blocks cross-site POSTs from sending the cookie (CSRF mitigation for state-changing requests from external origins). |
| `Secure` | `SESSION_SECURE` env | **Defaults to `false`.** MUST be set to `true`/`1` in production HTTPS. |
| `Path` | `/` | Valid for whole site. |
| `MaxAge` | derived from `Expires` | Matches the 30-day DB expiry. |

### Password Hashing (`internal/auth/password.go`)
**Algorithm: Argon2id** (winner of the Password Hashing Competition, memory-hard).

| Parameter | Value |
| :--- | :--- |
| Memory | 64 MB (`64 * 1024` KiB) |
| Iterations | 3 |
| Parallelism | 2 |
| Salt | 16 random bytes (`crypto/rand`) |
| Output key | 32 bytes |

Hash format: `$argon2id$v=19$m=65536,t=3,p=2$<base64-salt>$<base64-hash>`

- Comparison uses `crypto/subtle.ConstantTimeCompare` (timing-attack resistant).
- Parameters are parsed from the stored hash, so cost parameters can be upgraded later **without forcing password resets**.

### Login Failure Behavior
- Unsuccessful logins return a generic `"Invalid credentials"` message — they do **not** reveal whether the username or the password was wrong (prevents account enumeration via login).

---

## 3. Rate Limiting

### 3.1 IP-based Sliding-Window Limiter (`internal/ratelimit`)
A custom in-memory sliding-window limiter keys on the client's real IP (first value of `X-Forwarded-For`, falling back to `X-Real-IP`, then `RemoteAddr`). Stale entries are purged in the background every `2 × window` (minimum 1 minute). Stores per-IP timestamps in a `map[string][]time.Time` guarded by `sync.RWMutex`.

**Limits applied in `cmd/server/main.go`:**

| Endpoint | Limit | Window |
| :--- | :--- | :--- |
| `POST /api/auth/login` | 10 | 1 minute |
| `POST /api/auth/register` | 5 | 1 minute |
| `POST /api/setup` | 5 | 1 minute |
| `GET /api/invitations/validate` | 15 | 1 minute |
| `POST /api/upload` | 20 | 1 minute |

On exceed: `HTTP 429` + `Retry-After: 60` header + JSON error envelope.

> ⚠️ **Operational caveat:** This limiter is **in-memory** — limits reset on restart and do not scale across multiple backend instances. It is sufficient for a single-instance homelab deployment.

### 3.2 Per-User Action Limiting (application level)
These are checked inside handlers and keyed on the authenticated user (not the IP):

| Endpoint | Limit | Window |
| :--- | :--- | :--- |
| `POST /api/projects/{slug}/comments` | 5 | 1 minute |
| `POST /api/projects/{slug}/report` | 3 | 1 minute |
| `POST /api/comments/{id}/report` | 3 | 1 minute |

### 3.3 Not Currently Limited
| Endpoint | Why it's acceptable |
| :--- | :--- |
| All `GET` read endpoints (projects, users, activity, status) | Read-only; no state mutation. Pagination caps (`limit <= 50/100`) bound the cost of each request. |
| `POST /api/auth/logout` | Idempotent, cheap. |
| `PATCH /api/me` | Authenticated; throttling burden considered low priority. |
| `POST /api/hosting-requests` | Authenticated; not yet limited. |
| All `/api/admin/*` | Authenticated + role-gated; operator-only surface. |

### 3.4 No Rate Limiting at Nginx
The Nginx layer performs **no** `limit_req`/`limit_conn`. All throttling happens in the Go backend. Keep this in mind for the deployment contract: the backend must not be directly exposed to the public internet, or the app-layer limits can be bypassed.

---

## 4. Input Validation & Sanitization

### Strengths Implemented Today
| Area | Rules enforced |
| :--- | :--- |
| **Setup** | `name` ≤ 128 chars; `username` regex `^[a-zA-Z0-9_-]{3,30}$`; `email` parsed via `mail.ParseAddress`; `password` ≥ 8 chars; race-guarded against double-init in a transaction. |
| **Registration** | Non-empty `username` (trimmed, lowercased) and `email`; `password` ≥ 8 chars; `invitation_token` required when mode is `INVITE_ONLY`. |
| **Login** | Non-empty `email_or_username` and `password`. |
| **Comments** | Non-empty; **max 1000 chars**; HTML-escaped via `html.EscapeString()` before storage (neutralizes stored XSS). |
| **Reports** | `reason` whitelist (`SPAM`, `ABUSE_HARASSMENT`, `INAPPROPRIATE`, `MALICIOUS_SUSPICIOUS`, `OTHER`); `details` truncated to 500 chars; duplicate open reports blocked. |
| **Hosting subdomain** | Regex `^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`, length 2–63, reserved-word blocklist. |
| **Admin role/status** | Must be exactly `USER`/`ADMIN` and `ACTIVE`/`SUSPENDED`; self-demotion and last-admin demotion blocked. |
| **Access settings** | `registration_mode` must be `OPEN`/`INVITE_ONLY`/`CLOSED`; domain protocol prefixes stripped. |
| **Pagination** | `page`/`limit` integer-parsed with defaults; `limit` capped (max 50/100 depending on endpoint). |

### Known Gaps (validation)
- **No email format validation at registration** — registration accepts any non-empty string; only setup validates the format.
- **No max length** on `display_name`, `bio`, `avatar_url`, project `description`, `readme`, `repository_url`.
- Profile/project update payloads have no length bounds in handlers (DB columns are `TEXT`).

---

## 5. File Upload Security (`internal/storage/service.go`)

Uploads are one of the most abused surfaces on any server. Ngumpul Host applies layered defenses:

| Control | Value |
| :--- | :--- |
| Hard size limit | **10 MB** (`MaxUploadSizeBytes = 10 << 20`), enforced by `http.MaxBytesReader`, `ParseMultipartForm`, and Nginx `client_max_body_size 10M`. |
| Allowed MIME types | `image/jpeg`, `image/png`, `image/webp` only — sniffed from real content via `http.DetectContentType`, **not** trusted from the client-supplied `Content-Type`. |
| SVG | **Explicitly rejected** (SVGs can carry script payloads). |
| Max dimension | ≤ 4096 px per side (`MaxDimensionPixels`). |
| Max total pixels | ≤ 16 MP (`MaxTotalPixels = 16 * 1024 * 1024`) — prevents decompression bombs / pixel-flood DoS. |
| Dimension validation | Images actually decoded (`image.DecodeConfig`); WebP parsed via a dedicated VP8/VP8L/VP8X parser; malformed data rejected. |
| Filename | Server-generated UUID (`uuid.New()`) — unguessable; client filename never used. |
| Storage paths | Purpose-prefixed: `covers/`, `avatars/`, `uploads/`. |
| Object key validation on serve | `strings.Contains(key, "..")` rejected (path traversal protection). |
| Serving headers | `X-Content-Type-Options: nosniff`, `Cache-Control: public, max-age=2592000, no-transform`. |
| Authentication | Upload endpoint requires `RequireAuth`; anonymous upload impossible. |

---

## 6. CORS Configuration (`cnoscmd/server/main.go`)
- **Method whitelist:** `GET, POST, PUT, PATCH, DELETE, OPTIONS`.
- **Header whitelist:** `Accept, Authorization, Content-Type, X-CSRF-Token`.
- **Origin policy:** exact-match `AllowOriginFunc` — no wildcards. Allowed origins are:
  1. Hardcoded dev origins (`localhost`/`127.0.0.1` on 5173/3000/8080).
  2. The configured `APP_URL`.
  3. The instance domain read from the database (both `http://` and `https://` variants).
- **Credentials:** `AllowCredentials: true` but only when the origin is verified against the list above.
- **Preflight cache:** `MaxAge: 300` (5 minutes).

> ⚠️ `X-CSRF-Token` is an allowed header but **no CSRF token is currently generated or validated**. SameSite=Lax mitigates external-origin CSRF; there is no same-origin CSRF scenario in a cookie session unless a stored XSS exists. Document known-gap.

---

## 7. SQL Injection Surface

Every SQL query in the codebase uses **parameterized placeholders** (`$1`, `$2`, …) via `pgx`. Dynamic filters build parameter indices with `fmt.Sprintf("... = $%d", argIdx)` and append values to a typed `args` slice — never string-concatenated SQL. No ORM is used (by architectural mandate). Review checklist: `auth`, `comment`, `report`, `project`, `hosting`, `admin`, `setup`, `storage`, `system` handlers all verified parameterized.

---

## 8. Edge Ingress & Security Headers (Nginx)

### Present (`nginx/nginx.conf`)
```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
```

### Proxy rules
- `/api/`, `/uploads/`, `/go/`, `/health` → Go backend (`backend:8080`).
- `/` → SvelteKit SSR (`frontend:3000`), with WebSocket upgrade support.
- `/uploads/` gets a 30-day cache header.
- `client_max_body_size 10M` (matches backend upload cap).
- `/api/` upstream read timeout 90s (backend itself enforces 15s read / 30s write / 60s idle + chi `Timeout(60s)`).

### Missing headers (gap)
- ❌ **`Content-Security-Policy`** — no CSP anywhere. XSS payloads that survive escaping would execute unrestricted.
- ❌ **`Strict-Transport-Security`** — no HSTS; browsers won't force HTTPS for repeat visits.
- ❌ **`Permissions-Policy`** — camera/mic/geolocation not gated.

---

## 9. Server Hardening (Go)

| Control | Value |
| :--- | :--- |
| `ReadTimeout` | 15s |
| `WriteTimeout` | 30s |
| `IdleTimeout` | 60s |
| `middleware.Timeout` | 60s (context cancellation propagates to DB queries) |
| `middleware.Recoverer` | Converts goroutine panics into safe JSON 500 responses |
| `middleware.RequestID` | Unique request ID for traceability |
| `middleware.RealIP` | Trusted proxy header resolution before rate limiting |

---

## 10. Error Disclosure & Audit

### Good behavior
- Login: generic `"Invalid credentials"` (no account enumeration).
- Audit log table (`audit_logs`): role changes, status changes, project operations, comment deletions, report resolutions, hosting request reviews, and initial setup — actor, action, target, metadata, timestamps.
- `PasswordHash` never serialized to JSON responses.
- Registration conflict messages (`"Username already taken"` / `"Email already registered"`) reveal account existence — accepted tradeoff for self-hosted community UX.

### Gap
- Several **admin** handlers embed raw `err.Error()` into error responses (e.g. hosting subdomain updates, comment admin list). Low severity for a homelab, but should be sanitized to generic messages with the detail logged server-side.

---

## 11. Frontend-Layer Notes

- SvelteKit `hooks.server.ts` rewrites `/api/*` calls to the backend; it applies no security transformation itself (the browser holds cookies directly).
- Forms are validated **both** client-side (Superforms + Valibot, per-field inline errors) and server-side (server action `superValidate` + the Go backend's own checks). Schema-driven validation means client and server share the same rules.
- Uploads: no client-side file is ever trusted; the Go backend re-sniffs MIME and re-decodes dimensions.

---

## 12. Hardening Checklist (Recommended Before Public Launch)

Follow this list to close the documented gaps before exposing the instance beyond a trusted LAN:

- [ ] `SESSION_SECURE=true` in production (cookie `Secure` flag).
- [ ] Add `Strict-Transport-Security: max-age=63072000` to Nginx (HTTPS only).
- [ ] Add a `Content-Security-Policy` header; start restrictive (`default-src 'self'`) and loosen per feature.
- [ ] Add `Permissions-Policy` to disable camera/microphone/geolocation.
- [ ] Email format validation on registration (mirror the setup flow's `mail.ParseAddress`).
- [ ] Max-length bounds on `display_name`, `bio`, profile/project update fields.
- [ ] Sanitize admin-handler error messages (log the detail, return a generic message).
- [ ] Periodic expired-session cleanup worker (currently lazy on-access only).
- [ ] Optional: move per-IP rate limiting to Nginx `limit_req_zone` so the backend need not trust proxy headers.
- [ ] If sessions must survive without restart, replace/reinforce the in-memory limiter with a persistent store (Redis/DB) — only relevant for multi-instance deployments, which are out of scope by architectural mandate.