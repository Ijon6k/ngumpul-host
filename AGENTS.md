# AI Engineering Guidelines & Anti-Slop Directives (AGENTS.md)

> **Audience:** All AI coding agents (Antigravity, Claude Code, Cursor, Copilot, Gemini) and human contributors.  
> **Authority:** Absolute. These directives supersede generic LLM defaults, creative biases, and speculative habits.  
> **Scope:** Entire repository (`/backend`, `/frontend`, `/docs`, configuration, and migrations).

---

## 1. Repository Identity & Architectural Mental Model

**Ngumpul Host** is a quiet, self-hosted community server and project catalog for friends, experiments, and independent software.

### The 4 Non-Negotiable Mental Invariants:
1. **Single-Instance Community Server:** One deployment represents one specific homelab or community node. Never attempt to turn this into a multi-tenant SaaS. Do not add `tenant_id` or organizational partitioning.
2. **Calm Hardware Editorial:** The visual and product feel is grounded, editorial, and physical (inspired by real homelab hardware, Apple/Samsung calm typography, and Cloudflare-style console density). It is NOT a cyberpunk neon arcade or generic AI SaaS dashboard.
3. **Manual Operator Hosting:** Ngumpul Host is a **management and discovery layer**, NOT an automated Kubernetes PaaS or auto-deploy engine. The operator provisions containers; Ngumpul Host tracks, showcases, and monitors them.
4. **Real Linux Kernel Telemetry:** Every metric and status must bind directly to Linux kernel telemetry (`/proc/uptime`, `/proc/cpuinfo`, `/proc/meminfo`, `statvfs`). Never manufacture synthetic telemetry or decorative progress meters.

---

## 2. Mandatory Rules of Engagement

### 2.1. Pre-Flight Rule: Read `docs/` Before Editing
Before generating code, modifying markup, or changing schemas, you **MUST** view and inspect the relevant documentation in `docs/`:

| If your task involves... | You MUST read this first: |
| :--- | :--- |
| UI, components, colors, radii, typography | [`docs/frontend/design-system.md`](./docs/frontend/design-system.md), [`docs/frontend/components.md`](./docs/frontend/components.md), [`docs/frontend/architecture.md`](./docs/frontend/architecture.md) |
| Backend handlers, routing, middleware | [`docs/backend/architecture.md`](./docs/backend/architecture.md) & [`docs/backend/api.md`](./docs/backend/api.md) |
| Availability, uptime, gap detection, reboots | [`docs/backend/availability.md`](./docs/backend/availability.md) |
| System telemetry, CPU/RAM/NVMe probes | [`docs/backend/telemetry.md`](./docs/backend/telemetry.md) |
| Database migrations, tables, constraints | [`docs/backend/database.md`](./docs/backend/database.md) |
| Architecture, tenancy, ingress routing | [`docs/architecture.md`](./docs/architecture.md) |
| User/admin scenarios and operational flows | [`docs/workflows.md`](./docs/workflows.md) |

### 2.2. Post-Flight Rule: Synchronize Documentation
Code and documentation must never drift apart. Any task that alters an endpoint, table, token, or workflow **MUST** update the matching document in `docs/` in the exact same task:
- Modified API route or payload $\rightarrow$ Update [`docs/backend/api.md`](./docs/backend/api.md).
- Modified table or migration $\rightarrow$ Update [`docs/backend/database.md`](./docs/backend/database.md).
- Modified UI primitive, table, or radius $\rightarrow$ Update [`docs/frontend/design-system.md`](./docs/frontend/design-system.md) or [`docs/frontend/components.md`](./docs/frontend/components.md).
- Added/altered host probe $\rightarrow$ Update [`docs/backend/telemetry.md`](./docs/backend/telemetry.md) or [`docs/backend/availability.md`](./docs/backend/availability.md).

---

## 3. Tooling & Execution Standards

### Frontend Tooling
- **Package Manager:** Strictly use `bun` (never `npm`, `pnpm`, or `yarn`).
- **Type Checking:** `bun run check` (MUST pass with **0 errors and 0 warnings**).
- **Production Bundle:** `bun run build` (MUST build cleanly with `@sveltejs/adapter-node`).
- **Icons:** Exclusively use `phosphor-svelte`. Never use Lucide, FontAwesome, Heroicons, or raw inline `<svg>` markup.

### Backend Tooling
- **Compiler:** Go 1.23+ (`go build -v ./cmd/server`).
- **Tests:** `go test -v ./...`.
- **Database Driver:** `jackc/pgx/v5` with connection pooling (`*pgxpool.Pool`). Do not introduce ORMs (GORM, Ent).
- **Routing:** `go-chi/chi/v5` router with explicit dependency injection.

### Off-Limit Boundaries
- **NEVER** edit auto-generated files: `frontend/.svelte-kit/`, `frontend/build/`, `node_modules/`.
- **NEVER** modify `.env` without explicit instruction from the user.
- **NEVER** add new npm or Go dependencies without confirming they cannot be solved with standard libraries or existing packages (`phosphor-svelte`, `clsx`, `axios`, `chi`, `pgx`).

---

## 4. The Anti-AI-Slop Manifesto

### What is "AI Slop"?
**AI Slop** is mechanically generated code, design, or copy that compiles and passes basic checks, but lacks domain awareness, architectural restraint, and human editorial taste. It creates a false illusion of productivity while degrading maintainability, inflating codebase size, and rotting software over time.

You must actively prevent and eliminate AI Slop across all five engineering vectors:

---

### Vector A: UI & Visual Design Slop

| Slop Pattern | Symptom / Defect | Required Standard in Ngumpul Host |
| :--- | :--- | :--- |
| **Pill / Chip Spam** | Slapping `rounded-full` colored chips on every label, tag, button, and metric. | Rely on **typographic hierarchy and spatial proximity** first. Pills are reserved **strictly for status badges** (`Badge.svelte`, `StatusDot.svelte`). |
| **Monospace Caps Obsession** | Forcing `uppercase tracking-widest font-mono text-[10px]` on section headers and labels. | Use clean, calm sans-serif typography (`font-sans font-normal text-sm sm:text-base`) with normal sentence casing. Monospace is reserved for code, SHA hashes, IPs, and telemetry units. |
| **Cyberpunk / Neon Fantasy** | Random purple/cyan laser gradients, neon glow drop-shadows, radial background lights, and techno buzzwords. | **Calm Hardware Editorial**: Natural photography of real homelab nodes, deep neutral zinc/charcoal surfaces, and brand sky blue highlights (`#79AFC4`). |
| **Incoherent Border Radii** | Arbitrarily scattering `rounded-2xl`, `rounded-full`, and `rounded-none` haphazardly across controls and cards. | **Strict Token Hierarchy:**<br>• `radius-sm` (2px - 4px) for interactive controls (`Button`, `Input`).<br>• `radius-md` (6px - 8px) for structural containers (`Card`, `Table`).<br>• `radius-full` (9999px) for status dots & badges only. |
| **Manufactured Telemetry** | Fake pulsating dots, arbitrary animated progress meters, and decorative SVGs that represent zero real data. | Every telemetry widget must bind to real Linux kernel metrics (`/proc/cpuinfo`, `/proc/meminfo`, `statvfs`, `/proc/uptime`). If no data exists, render a clean `no_data` state. |
| **Arbitrary Cover Ratios** | Using arbitrary, inconsistent aspect ratios for project artwork. | **16:9 Standard Artwork Ratio**: All project covers (cards, overview pages, and dashboard previews) must use `aspect-[16/9]` (`ProjectCover aspectRatio="16/9"`). |
| **Low-Contrast Micro-Text** | Tiny gray-on-dark text (`text-neutral-600` on dark canvas) that sacrifices legibility for "mood". | High contrast readability adhering to WCAG AA standards. Interactive elements must be instantly readable without zooming. |
| **Icon Inconsistency** | Mixing heroicons, lucide, font-awesome, and raw inline `<svg>` blobs. | Use **Phosphor Icons** (`phosphor-svelte`) exclusively with `weight="regular"` (or `"bold"` for directional arrows and checks). |

---

### Vector B: Code & Engineering Slop

| Slop Pattern | Symptom / Defect | Required Standard in Ngumpul Host |
| :--- | :--- | :--- |
| **Near-Duplicate Utilities** | Generating a new helper function because the model didn't search the existing codebase. | Search first! Reuse established utilities (`formatDate`, `extractError`, `api`, `response.JSON`, `response.Error`). Never re-implement date formatters or HTTP helpers. |
| **Speculative Indirection** | Building generic factories, multi-tier interfaces, and pass-through wrappers for single-use operations. | **Write direct, simple, single-purpose code.** Do not build abstract abstractions until at least 3 concrete use cases exist. |
| **Narrative Echo Comments** | Comments that state the obvious syntax: `// check if user is nil`, `// loop through projects`, `// return result`. | **Delete narrative comments.** Comments must only explain *why* non-obvious business rules, race conditions, or hardware constraints exist. |
| **Defensive Null-Check Litter** | Adding `?.` and fallback operators on fields guaranteed non-null by TypeScript or PostgreSQL `NOT NULL` constraints. | Trust types and database invariants. If a field is `NOT NULL`, treat it as non-null. Avoid unnecessary defensive boilerplate. |
| **Swallowing Errors & Silent Failures** | `catch (e) { /* do nothing */ }` or `if err != nil { return nil, nil }` hiding bugs until production incidents. | Always handle or return errors with context: `fmt.Errorf("failed to load project %s: %w", slug, err)`. In HTTP handlers, emit appropriate status codes (`400`, `404`, `500`). |
| **Shallow "Green" Tests** | Writing tests that execute code but assert nothing (`expect(res).toBeDefined()`) to inflate test coverage. | Test real business invariants: state transitions, permission boundaries, calculation edge cases, and error branches. |

---

### Vector C: Folder Structure & Organization Slop

| Slop Pattern | Symptom / Defect | Required Standard in Ngumpul Host |
| :--- | :--- | :--- |
| **Ad-Hoc Dumping Grounds** | Creating arbitrary directories like `helpers/`, `misc/`, `common/`, or single-file folders. | Follow the established package and directory structure: `frontend/src/lib/components/ui/`, `frontend/src/lib/components/layout/`, `backend/internal/<domain>/`. |
| **Layer Confusion** | Putting database queries directly inside HTTP handlers or embedding backend business logic in Svelte views. | Maintain clear architectural layers: Handlers parse HTTP requests and delegate to services; services interact with the database via `pgxpool`; Svelte views render UI and query APIs via TanStack Query. |
| **Naming Inconsistency** | Mixing `camelCase`, `kebab-case`, and `snake_case` randomly within the same layer. | • Go: PascalCase for exported, camelCase for internal.<br>• Svelte: PascalCase for components (`Button.svelte`), camelCase for utility modules (`format.ts`).<br>• Database: snake_case for tables and columns (`created_at`, `hosting_requests`). |

---

### Vector D: Methodology & Approach Slop

| Slop Pattern | Symptom / Defect | Required Standard in Ngumpul Host |
| :--- | :--- | :--- |
| **Assumption-First Coding** | Writing code based on assumptions without viewing or searching active files. | Always view files and read relevant documentation before modifying code. Ground every change in reality. |
| **Brute-Force Rewriting** | Rewriting a 500-line file when a surgical 5-line edit was required. | Use targeted edits (`replace_file_content` or `multi_replace_file_content`). Preserve existing code, comments, and structure. |
| **Dependency Hallucination** | Running `npm install` or `go get` for third-party libraries for problems solved by native APIs or installed packages. | Check existing dependencies first. The repository already includes `phosphor-svelte`, `clsx`, `axios`, `pgx/v5`, and `chi/v5`. Avoid introducing new dependencies. |
| **Scope Creep** | Refactoring unrelated files or changing unrequested visual styles under the guise of "cleaning up". | Stay strictly within the scope requested by the user. Do not touch unrelated files without explicit instruction. |

---

### Vector E: Copywriting & Tone Slop

| Slop Pattern | Forbidden Examples | Required Editorial Standard |
| :--- | :--- | :--- |
| **Silicon Valley AI Buzzwords** | "Revolutionize", "Blazing fast", "Seamlessly", "Elevate", "Next-generation", "Unleash power" | Calm, direct, honest prose: *"A quiet corner of the internet for side projects, bots, and friends."* |
| **Manufactured Marketing Claims** | "Trusted by thousands of leading teams worldwide" | Grounded homelab truth: *"Running on real hardware in Jakarta. Kept online without cloud bills."* |
| **Corporate Jargon Overkill** | "Leverage synergistic edge topologies to maximize availability throughput" | Human engineering clarity: *"Monitored via kernel /proc telemetry with continuous reboot gap detection."* |

---

## 5. Concrete Code Smells vs. Good Implementations

### Example 1: Button Styling (UI Radius & Slop Rule)
```svelte
<!-- ❌ BAD (AI Slop: Incoherent radius, pill spam, loud gradient) -->
<button class="rounded-full bg-gradient-to-r from-purple-500 to-cyan-500 text-white font-mono uppercase tracking-widest text-[10px] py-1 px-4 shadow-lg shadow-cyan-500/50">
    Deploy Application
</button>

<!-- ✅ GOOD (Ngumpul Host Standard: radius-sm, clean typography, calm contrast) -->
<button class="px-4 py-2 rounded-md bg-neutral-900 text-white dark:bg-white dark:text-neutral-950 text-xs sm:text-sm font-medium hover:opacity-90 transition-opacity">
    Deploy application
</button>
```

### Example 2: Error Handling (Backend Slop Rule)
```go
// ❌ BAD (AI Slop: Swallowing error, narrative echo comment)
// get the project
proj, err := s.GetProject(ctx, slug)
if err != nil {
    return nil, nil // ignore error
}

// ✅ GOOD (Ngumpul Host Standard: Contextual wrapping, clear semantics)
proj, err := s.projectRepo.GetBySlug(ctx, slug)
if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrProjectNotFound
    }
    return nil, fmt.Errorf("failed to fetch project %q: %w", slug, err)
}
```

### Example 3: Icon Usage (Phosphor Rule)
```svelte
<!-- ❌ BAD (AI Slop: Raw SVG blob pasted into component) -->
<svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor">
    <path d="M12 2v20M2 12h20" />
</svg>

<!-- ✅ GOOD (Ngumpul Host Standard: Phosphor Svelte Component) -->
<script>
    import { Plus } from 'phosphor-svelte';
</script>
<Plus size={16} weight="bold" />
```

---

## 6. Pre-Completion Quality Gate Checklist

Before declaring any task complete or presenting diffs to the user, run through this mental quality gate:

1. [ ] **Did I read `docs/` first?** Did I inspect existing implementations before creating new code?
2. [ ] **Did I reuse existing design primitives?** (`$lib/components/ui`, `$lib/components/layout`, Phosphor Icons).
3. [ ] **Are all border radii strictly compliant?**
   - Interactive controls: `radius-sm` (2px - 4px).
   - Structural containers: `radius-md` (6px - 8px).
   - Status badges only: `radius-full` (9999px).
4. [ ] **Is there any AI slop present?**
   - [ ] No pill/chip spam.
   - [ ] Strict 14px typography floor maintained (no text-[11px], text-[10px], or sub-14px microtext).
   - [ ] No uppercase monospace headers.
   - [ ] No cyberpunk/neon laser glows.
   - [ ] No narrative syntax echo comments.
   - [ ] No duplicate helper utilities.
   - [ ] No swallowed errors.
   - [ ] No marketing buzzword copy.
5. [ ] **Did I verify the changes?**
   - Frontend: `bun run check` (**0 errors, 0 warnings**) and `bun run build`.
   - Backend: `go build ./cmd/server` and `go test ./...`.
6. [ ] **Did I update `docs/`?** If an API, schema, token, or workflow changed, is the documentation synchronized?
