# Ngumpul Host Documentation Hub

Welcome to the technical, architectural, and operational documentation for **Ngumpul Host**—a quiet, self-hosted community platform and project portal built on Linux bare-metal or virtualized infrastructure.

---

## Master Documentation Directory

### 1. General & System Architecture
| Document | Description |
| :--- | :--- |
| [**System Overview & Product Vision**](./overview.md) | High-level product thesis, user mental model, and core capabilities. |
| [**System Architecture & Topology**](./architecture.md) | Network ingress, reverse proxying, single-instance tenancy, security boundaries, and container topology. |
| [**Security & Hardening Reference**](./security.md) | Password hashing, session management, rate limiting, upload defenses, CORS, and the hardening checklist before public launch. |
| [**End-to-End Operational Workflows**](./workflows.md) | Visual sequence diagrams for member onboarding, hosting requests, container provisioning, and reboot reconciliation. |

### 2. Frontend Engineering (`docs/frontend/`)
| Document | Description |
| :--- | :--- |
| [**Frontend Architecture & Standards**](./frontend/architecture.md) | SvelteKit 2 + Svelte 5, TanStack Query v5 state management, layout spacing rules, and Phosphor Icons integration. |
| [**Design System & Aesthetics**](./frontend/design-system.md) | Anti-slop design philosophy, typography hierarchy, border-radius tokens, and semantic color palette. |
| [**Kitab Komponen (Master Component Bible)**](./frontend/components.md) | **Canonical AI & Contributor Reference:** Exhaustive guide with Table of Contents, props, slots, design rules, and code examples for all UI primitives, status cards, bento widgets, and the reusable table system. |

### 3. Backend Engineering (`docs/backend/`)
| Document | Description |
| :--- | :--- |
| [**Backend Modular Monolith Architecture**](./backend/architecture.md) | Go 1.23+ modular monolith, package design, dependency injection, middleware pipeline, and graceful shutdown. |
| [**REST API Specification & Data Contracts**](./backend/api.md) | Complete REST API endpoint reference, authentication cookies, request/response JSON envelopes, and error codes. |
| [**Database Schema & Data Dictionary**](./backend/database.md) | PostgreSQL 16 entity-relationship model (ERD), table definitions, constraints, indexes, and session lifecycle. |
| [**Availability & Mathematical Uptime Engine**](./backend/availability.md) | Kernel-driven gap detection, `/proc/sys/kernel/random/boot_id`, in-place state tables, and bounded uptime formulas. |
| [**Native Host Telemetry Engine**](./backend/telemetry.md) | Direct Linux kernel telemetry probing (`/proc/cpuinfo`, `/proc/meminfo`, `statvfs`, NVMe sysfs probing, and edge ping). |

---

## Guidelines for Contributors & AI Agents
All contributors and AI coding assistants MUST review [`AGENTS.md`](../AGENTS.md) before making changes to code, interfaces, or database schemas. Documentation must be kept in sync with code at all times.
