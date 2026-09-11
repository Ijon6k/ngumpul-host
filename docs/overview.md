# Ngumpul Host

> A small corner of the internet for friends, projects, and experiments.

Ngumpul Host is an open-source, self-hostable community and project portal built around a shared hosting environment.

It allows a small group of people to:

- create a simple public profile
- showcase their projects
- request a project to be hosted
- see the status of their hosted projects
- discover projects and people in the community
- receive notifications about important project events

Ngumpul Host does **not** perform arbitrary hosting or deployment automatically. Project hosting is intentionally controlled by an administrator.

The administrator decides where and how a project is actually hosted, then registers the resulting project and public domain in Ngumpul Host.

The application is therefore a **community/project management layer**, not a hosting platform or deployment engine.

---

## 1. Product Vision

Ngumpul Host is designed around a simple idea:

> A small, shared corner of the internet where friends can build, host, and showcase things together.

The application feels like a personal community website rather than:
- a commercial hosting platform
- a DevOps dashboard
- a social media clone
- a SaaS administration panel
- a server monitoring application

The interface communicates that real infrastructure exists behind the website, but infrastructure details remain mostly invisible to normal users. Technical information appears through carefully designed metadata, status indicators, project parameters, and activity events rather than through raw server dashboards.

---

## 2. Core Product Model

Ngumpul Host has three major concepts:
1. **People**
2. **Projects**
3. **Hosting**

People own projects. Projects can be hosted by the administrator. Hosting is an operational process controlled outside Ngumpul Host.

```text
User
  │
  │ submits hosting request
  ▼
Ngumpul Host
  │
  │ administrator reviews request
  ▼
Administrator
  │
  │ manually sets up hosting
  ▼
External infrastructure
  │
  │ public domain assigned
  ▼
Ngumpul Host
  │
  │ project becomes published
  ▼
Public project page
```

The actual hosting infrastructure can be anything:
- Docker / Docker Compose
- Dokploy / Coolify
- Kubernetes
- VPS / Bare Metal / Home Server

Ngumpul Host does not depend on any particular hosting provider.

---

## 3. Important Product Boundary

**Ngumpul Host manages:**
- people & profiles
- projects & project metadata
- hosting requests & administrator review
- public project pages
- project ownership
- safe application activity events
- notifications & public status

**Ngumpul Host does NOT manage:**
- arbitrary server provisioning or VPS creation
- DNS provisioning
- arbitrary Docker execution or shell commands
- arbitrary deployment from user submissions
- server credentials & infrastructure secrets
- raw container management & server logs
- SSH access

---

## 4. Visual Design Direction

Primary design direction: **Editorial × Bento × Quiet Technical**

The visual language is:
- calm, minimal, structured, editorial
- slightly technical, personal, spacious
- information-oriented

Visual characteristics:
- warm off-white / light gray backgrounds (`#F4F7F8` paper light, `#101618` night dark)
- charcoal or near-black text
- restrained sky-blue accent color (`#79AFC4` / `#8FC5D8`)
- thin borders & subtle rounded corners (8px–12px)
- large typography (`Manrope` display, `DM Sans` interface)
- large numerical parameters
- generous whitespace
- asymmetric bento layouts & project screenshots
- minimal status indicators (`● Online`)

---

## 5. Technology Stack

- **Frontend:** SvelteKit + TypeScript (Bun tooling)
- **Backend:** Go (Modular Monolith)
- **Database:** PostgreSQL
- **Reverse Proxy:** Nginx
- **Deployment:** Docker Compose
- **Authentication:** Server-side sessions + Argon2id
- **File Storage:** Local filesystem abstraction (`/data/uploads`)
- **API:** REST `/api/v1`

---

## 6. Technical Product Specification

### Tenancy Model
Single-instance, multi-user application. One deployment represents one community with many users, projects, and hosting requests. Data isolation is enforced via user ownership and role-based access control (`USER`, `ADMIN`).

### Roles
- **USER:** Normal community member. Can manage own profile, view public projects/members/activity, submit hosting requests, view and manage allowed metadata of own projects, and receive notifications.
- **ADMIN:** Instance operator. Reviews hosting requests, approves/rejects, creates/edits projects, publishes/archives projects, sets public status, manages members, and monitors application health.

### Project Lifecycle
`UNPUBLISHED` ➔ `SETUP` ➔ `ONLINE` (or `OFFLINE` if health check fails) ➔ `ARCHIVED`.

### Hosting Request Lifecycle
`PENDING` ➔ `REVIEWING` ➔ `APPROVED` ➔ `SETUP` ➔ `COMPLETED` (or `REJECTED` / `CANCELLED`).

### Reverse Proxy & Routing Principles
Nginx acts as the single public entry point:
- Routes explicitly configured upstreams (SvelteKit frontend, Go backend API `/api/v1/`, uploaded media `/uploads/`).
- Any unmatched/default request falls back to: `http://1111:80`.
- Nginx is never configured as an open proxy.

---

## 7. Deep-Dive Documentation Links

- **Technical Architecture, Database Schema, APIs & Uptime Formulas:** [`docs/technical.md`](./technical.md)
- **Design System, Anti-Slop Guidelines & UI Tokens:** [`docs/design.md`](./design.md)

