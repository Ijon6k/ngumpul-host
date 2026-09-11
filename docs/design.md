# Ngumpul Host — Design System & Interface Guidelines

> **Architectural Core:** Quiet Technical × Editorial Infrastructure × High-Density Utility  
> **Target Aesthetic:** A calibrated balance between human editorial warmth (public landing & showcase) and dense, quiet infrastructure precision (status pages & management dashboards).

---

## 1. Visual Thesis & Philosophy

Ngumpul Host is an open-source, self-hosted platform running on real hardware at home. The interface must communicate honesty, calm engineering, and tactile digital craftsmanship.

### Core Tenets
1. **Calm & Restrained:** Low visual noise, generous whitespace, hair-thin borders (`1px`), and deliberate pacing. The interface feels like an intentionally crafted physical instrument.
2. **Dense Technical Utility:** Dashboards and telemetry tables prioritize fast scannability, tabular alignment, and unambiguous data over decorative flourishes.
3. **Day ➔ Night in the Same Material:** Dark mode is not an oversaturated cyberpunk neon shift. It is the exact same anodized metal and matte paper material viewed under evening light.
4. **Authenticity Over Illusion:** Every statistic, uptime bar, and latency metric reflects real native host OS probes. Never simulate, fake, or exaggerate.

---

## 2. Anti-AI-Slop Directives (Mandatory Quality Gates)

These rules are non-negotiable. Any AI assistant or developer writing UI code must strictly adhere to these exclusions:

### 🚫 Forbidden Patterns (Slop Catalog)
* **No Random Pill/Chip Badges Floating Everywhere:** Do not insert pill chips above section headings (e.g., `[ INFRASTRUCTURE V2 ]` or `[ FEATURES ]`). Badges are strictly reserved for genuine semantic status (`● Operational`, `Active`, `Admin`). Hierarchy must be established through font size, weight, and proximity—not floating badge boxes.
* **No Capslock Monospace Headings:** Headings must **never** be formatted in all-caps monospace (`SYSTEM_TELEMETRY_LOGS`, `CORE_SERVER_SPECS`). All headings must use clean, proportional sans-serif in standard Title Case or sentence case. Monospace is strictly reserved for raw numbers, hashes, code snippets, timestamps, and terminal commands.
* **No Edgy AI Buzzwords & Marketing Cringe:** Never generate hallucinated corporate jargon (*"Empowering next-gen decentralized quantum synergy"* or *"Revolutionizing seamless cloud paradigms"*). Copy must be calm, direct, and human:
  * *Good:* *"A single machine at home, kept online for friends and experiments. No cloud bills, no hidden layers—just a server we manage together."*
  * *Bad:* *"Enterprise-grade hyper-resilient containerized infrastructure solutions."*
* **No Cyberpunk Neon & Blob Gradients:** Absolutely ban purple-to-cyan gradient text, floating blurred color blobs, glowing neon borders, and dark glassmorphic transparency overlays that reduce contrast and legibility.
* **No Fake or Extrapolated Telemetry:** Never claim 100% 30-day uptime if monitoring only started 4 hours ago. Missing history must be represented honestly as neutral unmonitored states (e.g., neutral gray blocks), never as false uptime or artificial incidents.
* **No Repetitive 3-Card AI Templates:** Never generate 3 identical rounded cards in a row with generic Lucide icons (*Rocket, Shield, Zap*). Bento layouts must feature asymmetric weights, varied data densities, and functional differentiation.

---

## 3. Border Radius & Geometric Hierarchy

To preserve structural discipline and eliminate awkward visual rounding, all UI elements must strictly follow this token hierarchy. Arbitrary radius values outside this scale are forbidden.

### Token Scale

| Token | Size | Tailwind Class | Primary Element Usage |
| :--- | :--- | :--- | :--- |
| **`radius-none`** | `0px` | `rounded-none` | Full-width mobile screen edges, continuous dividers, table outer boundaries. |
| **`radius-sm`** | `2px` – `4px` | `rounded-sm` / `rounded` | Interactive controls: Buttons, Form Inputs, Select Dropdowns, Checkboxes. |
| **`radius-md`** | `6px` – `8px` | `rounded-md` / `rounded-lg` | Structural containers: Content Cards, Dashboard Panels, Inner Table Containers, Alert Banners. |
| **`radius-lg`** | `12px` – `16px` | `rounded-xl` / `rounded-2xl` | Floating overlays: Modals, Popovers, Slide-over Sheets, Floating Toast Notifications. |
| **`radius-full`** | `9999px` | `rounded-full` | Semantic indicators: Status Badges, System Health Dots, User Avatars, Segmented Pill Switches. |

### Specific Element Rules
1. **Inputs & Micro-Actions (`radius-sm`):** Buttons and form fields must remain crisp (`2px`–`4px`). Firm corners enhance text readability, anchor user focus, and communicate professional utility. *(Exception: Icon-only circular buttons may use `radius-full`).*
2. **Cards & Panels (`radius-md`):** Bento boxes, metric panels, and dashboard widgets use `6px`–`8px`. This avoids cartoonish bubble shapes while providing comfortable separation from the canvas.
3. **Overlays & Dialogs (`radius-lg`):** Modals and toasts floating above the application plane use `12px`–`16px` to naturally decouple from the rigid grid beneath.
4. **Semantic Badges & Pills (`radius-full`):** Badges and tags use pill geometry to immediately distinguish informative metadata from clickable rectangular buttons.

### The Nested Radius Formula
When nesting a rounded card or element inside a parent container with padding:
$$\text{Radius}_{\text{inner}} = \text{Radius}_{\text{outer}} - \text{Padding}$$
*Never allow inner corners to collide with or pinch against parent borders.*

---

## 4. Typography & Information Architecture

### Font Families
* **Display & Headings:** `Geist` / `Inter` / `Plus Jakarta Sans` (sans-serif)
  * Clean, geometric, neutral grotesque with high legibility and tight letter-spacing (`tracking-[-0.02em]` to `tracking-[-0.035em]`).
* **Body & Labels:** `Inter` / `Geist` (sans-serif)
  * Optimized for sustained readability across light and dark modes. Pacing: `leading-relaxed`, line length capped at `65ch`.
* **Technical Telemetry & Code:** `JetBrains Mono` / `ui-monospace`
  * Exclusively for IP addresses, git commit SHAs, latency numbers (`29 ms`), port numbers, memory sizes (`23.1 GB`), and terminal outputs.

### Hierarchy Rules
* **Proximity over Containers:** Group related elements with tight spacing (`gap-1.5` / `gap-2`) before introducing an enclosing box or border.
* **Scale over Tags:** Differentiate primary metrics using scale (`text-3xl font-normal tracking-tight`) rather than redundant labeling.
* **Sentence Case Dominance:** Headings, button labels, and descriptions use natural sentence or title capitalization. No uppercase yelling.

---

## 5. Color Palette & Design Tokens

Always reference predefined CSS custom properties. Ad-hoc hex values in component templates are strictly prohibited.

### Surface & Canvas Tokens
| Token | Light Value | Dark Value | Purpose |
| :--- | :--- | :--- | :--- |
| `--bg-canvas` | `#FAFBFC` | `#080A0D` | Main page background canvas. |
| `--bg-surface` | `#FFFFFF` | `#0E1217` | Cards, panels, table surfaces, elevated blocks. |
| `--bg-muted` | `#F3F5F7` | `#151A21` | Secondary backgrounds, bento grid foundations, table headers. |
| `--bg-hover` | `#E9EDF1` | `#1C232D` | Hover states, active list rows, subtle dividers. |

### Border & Hairline Tokens
| Token | Light Value | Dark Value | Purpose |
| :--- | :--- | :--- | :--- |
| `--border-hairline` | `#E1E6EB` | `rgba(255, 255, 255, 0.08)` | Crisp 1px structural borders and card dividers. |
| `--border-subtle` | `#ECF0F4` | `rgba(255, 255, 255, 0.04)` | Secondary inner dividers, subtle row separators. |

### Text & Ink Tokens
| Token | Light Value | Dark Value | Purpose |
| :--- | :--- | :--- | :--- |
| `--text-main` | `#0E1114` | `#F3F5F7` | Primary headings, prominent data metrics, active labels. |
| `--text-secondary` | `#47545E` | `#8E9BA4` | Body copy, secondary parameters, table headers. |
| `--text-muted` | `#798791` | `#56646E` | Tertiary metadata, timestamps, unmonitored indicators. |

### Brand & Utility Accents
| Token | Light Value | Dark Value | Purpose |
| :--- | :--- | :--- | :--- |
| `--accent-sky` | `#4A8FA8` | `#7CB1C7` | Core brand identity: calm sky-blue for focus states and graphs. |
| `--accent-orange` | `#FF5500` | `#FF5500` | Signature brand punch: subtle logo dot and primary notifications. |
| `--cf-blue` | `#0284C7` | `#38BDF8` | Technical links, interactive accents, infrastructure badges. |
| `--color-success` | `#2E8555` | `#45A36A` | Operational status, healthy telemetry, verified indicators. |
| `--color-warning` | `#99732B` | `#C29E4D` | Degraded service, memory threshold warnings, pending states. |
| `--color-danger` | `#B84747` | `#CC5454` | Service down, critical alerts, incident records. |

---

## 6. Component Architecture & Reusability Workflow

### Rule 1: Reuse Existing Components First
Before writing a new component, developers and AI agents **must inspect** existing components:
* **Status & Telemetry:** `$lib/components/status/` (`StatusOverviewCard`, `AvailabilityGrid`, `ServiceHealthCard`, `IncidentHistoryCard`, `StatusIndicator`).
* **Bento & Hardware:** `$lib/components/bento/` (`HardwareCard`, `MemoryCard`, `StorageCard`, `NetworkCard`, `AvailabilityCard`, `ServerPhotoCard`, `CpuPackageVisual`).
* **General UI:** `$lib/components/` (`Header`, `ThemeToggle`, `ProjectCard`, `ActivityTimeline`).

### Rule 2: The Atomic Construction Protocol
When a new component is strictly necessary:
1. **Build Atomically:** Construct from the smallest reusable unit up (e.g., `Badge` ➔ `StatusRow` ➔ `HealthTable`).
2. **Parameterize via TypeScript Interfaces:** Never hardcode text or values. Expose clean Svelte props with TypeScript interfaces located in `$lib/types/`.
3. **No Embedded Colors:** All colors must route through semantic CSS tokens (`--bg-surface`, `--border-hairline`, `--text-main`).
4. **Theme Resilience:** Verify both Light mode and Dark mode transitions without layout jumps or contrast loss.

---

## 7. Dense Technical Dashboard & Management Pages

For administrative views, project management, and server settings:
* **High Information Density:** Keep vertical padding compact (`py-2.5` to `py-3` on table rows).
* **Grid Alignments:**
  * Alphanumeric text & names: **Left-aligned**.
  * Numerical metrics, timestamps, sizes: **Right-aligned** with `font-mono`.
  * Status badges: **Centered** or **Right-aligned**.
* **Restrained Controls:** Use `radius-sm` for action buttons. Use ghost buttons with subtle borders for secondary options.
* **Inline Feedback Over Modals:** Prefer inline validation messages, subtle toast banners, and clear empty states over blocking modals.
