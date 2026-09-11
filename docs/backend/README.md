# Backend Engineering Documentation

> **Runtime:** Go 1.23+  
> **Router:** `go-chi/chi/v5`  
> **Database Driver:** `jackc/pgx/v5`

This folder contains all backend-specific technical documentation.

---

| Document | Description |
| :--- | :--- |
| [**Architecture & Monolith Design**](./architecture.md) | Package structure, dependency injection, middleware pipeline, and graceful shutdown. |
| [**REST API Specification**](./api.md) | Complete endpoint reference with request/response envelopes and error codes. |
| [**Database Schema**](./database.md) | PostgreSQL 16 ERD, table specs, constraints, and migration history. |
| [**Availability & Uptime Engine**](./availability.md) | Kernel gap detection, in-place state tables, and mathematical uptime formulas. |
| [**Native Host Telemetry**](./telemetry.md) | Direct `/proc` and `/sys` Linux kernel probing for CPU, RAM, NVMe, and latency. |

---

> Back to [**Documentation Hub**](../README.md)
