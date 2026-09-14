# Ngumpul Host — Design System & Visual Standards

> **Architectural Core:** Quiet Technical × Editorial Infrastructure × High-Density Utility  
> **Target Aesthetic:** A calibrated balance between human editorial warmth (public landing & showcase) and dense, quiet infrastructure precision (status pages & management dashboards).  
> **Location:** [`docs/frontend/design-system.md`](./design-system.md)

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
| **`radius-sm`** | `2px` – `4px` | `rounded-sm` / `rounded-[4px]` | Interactive controls: Buttons, Form Inputs, Select Dropdowns, Checkboxes. |
| **`radius-md`** | `6px` – `8px` | `rounded-md` / `rounded-[8px]` | Structural containers: Content Cards, Dashboard Panels, Inner Table Containers, Alert Banners. |
| **`radius-lg`** | `12px` – `16px` | `rounded-xl` / `rounded-2xl` | Floating overlays: Modals, Popovers, Slide-over Sheets, Floating Toast Notifications. |
| **`radius-full`** | `9999px` | `rounded-full` | Semantic indicators: Status Badges, System Health Dots, User Avatars, Segmented Pill Switches. |

### Specific Element Rules
1. **Inputs & Micro-Actions (`radius-sm`):** Buttons and form fields must remain crisp (`2px`–`4px`). Firm corners enhance text readability, anchor user focus, and communicate professional utility. *(Exception: Icon-only circular buttons may use `radius-full`).*
2. **Cards & Panels (`radius-md`):** Bento boxes, metric panels, and dashboard widgets use `6px`–`8px`. This avoids cartoonish bubble shapes while providing comfortable separation from the canvas.
3. **Overlays & Dialogs (`radius-lg`):** Modals and toasts floating above the application plane use `12px`–`16px` to naturally decouple from the rigid grid beneath.
4. **Semantic Badges & Indicators:**
   - **Inline & Floating Badges (`radius-full`):** Standard status tags and dots use pill geometry (`rounded-full`) to distinguish metadata from buttons.
   - **Project Card Corner Status (`rounded-bl-md`):** On project cover images, status badges sit flush against the top-right corner (`top-0 right-0`) with internal curvature `rounded-bl-md`. Uses theme-adaptive neutral surface (`bg-(--bg-surface)/90 dark:bg-neutral-900/90 backdrop-blur-xs`) with subtle inner hairline borders (`border-b border-l border-(--border-hairline)`). Renders direct semantic colored text without pill outlines, colored borders, or dot circles.

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

### Hierarchy & Scale Standards
* **Base Typography Scale (Tailwind v4 theme tokens):**
  * `--text-xs: 0.875rem` (14px) — metadata, secondary timestamps, tags, badges. **Strict floor: no font size smaller than 14px is permitted anywhere in the UI.**
  * `--text-sm: 1rem` (16px) — body copy, table cells, form labels, card descriptions, buttons.
  * `--text-base: 1.125rem` (18px) — prominent body, feature intro leads.
* **Page Titles:** `text-2xl sm:text-3xl lg:text-4xl font-normal tracking-[-0.03em]`.
* **Section Titles:** `text-lg sm:text-xl font-medium tracking-tight`.
* **Table Headers:** `text-sm font-semibold text-(--text-secondary) uppercase tracking-wider`.
* **Table Body Cells:** `text-sm text-(--text-main)`.
* **Technical Metadata / Code:** `text-xs font-mono`.
* **Micro-badges (`StatusDot`):** `text-xs px-2.5 py-0.5 font-medium`.
* **Breadcrumb Navigation:** `text-sm` (16px / `1rem`) with `text-(--text-muted)` and active item `text-(--text-main) font-medium`.
* **Project Artwork Aspect Ratio:** Standard `16:9` aspect ratio (`aspect-[16/9]`) across catalog cards, project showcase pages, and personal workspace dashboard views via `ProjectCover.svelte`.
* **Catalog Grid & Card Parity:** Both the `/projects` catalog and the landing page showcase (`/+page.svelte`) share identical card sizing via `max-w-5xl` containers and 2-column grids (`grid-cols-1 md:grid-cols-2 gap-6 sm:gap-8`). The landing page renders a balanced 2x2 grid (max 4 projects).
* **Card Description Line-Clamp:** Card descriptions are clamped to strictly 2 lines (`line-clamp-2 min-h-[2.5rem] sm:min-h-[2.85rem] leading-relaxed`) to prevent 2nd line letter clipping and maintain horizontal alignment across grid rows.

### §4.1 Mobile Text Layout Rules (Android / Small Viewport)
Container wrappers (`container mx-auto px-4 sm:px-6 max-w-*`) are already responsive. The text itself must also be. Mandatory rules:
1. **Responsive heading scale:** Page titles use `text-2xl sm:text-3xl lg:text-4xl`; section titles `text-lg sm:text-xl`. Never ship a mobile-first size above `text-2xl` without an `sm:` step-up.
2. **No cramped `justify-between` rows:** Any flex row containing a long-ish label + value or label + action must use `flex-col sm:flex-row ... gap-*`. At 360px the available content width is ~328px; compute before shipping.
3. **Long unbreakable strings:** URLs, emails, subdomains, and mono prefixed strings get `min-w-0 break-all` (or `truncate` only where truncation is safe) inside flex children, plus `flex-wrap` on the container as a safety valve.
4. **Avoid `tracking-wider uppercase` on unit labels** inside narrow status grids; it inflates width ~40% and breaks the 3-column time bounds row. Prefer plain `text-xs opacity-75`.
5. **No descender clipping:** Never fix `h-5`/`h-6` on text rows; use `min-h-[1.35rem]` (matches the `text-xs` line-height) or line-driven height.
6. **Footer/nav link rows:** apply `flex-wrap gap-x-3 gap-y-1` so link chains degrade gracefully instead of overflowing.
7. **Headlines:** add `text-balance` on project/detail titles to avoid orphaned single words on narrow screens.

### §4.2 Skeleton / Loading Screens
Loading UI must **mirror the real layout** — same containers, spacing, and proportions as the loaded state — to eliminate layout shift.
- Use the `Skeleton` primitive (`$lib/components/ui/Skeleton.svelte`) with `variant="line|block|avatar|badge|button"` and width-override classes; it renders `--bg-muted` + `animate-pulse`.
- Scaffold containers get `role="status" aria-live="polite"` and a `sr-only` message; real layout text is replaced, not appended.
- Composed mirrors live beside their targets: `SkeletonProjectCard.svelte` ↔ `ProjectCard.svelte`. Shared table loaders render skeleton rows inside the `Table` component (`skeletonRows` prop).


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
| `--accent-orange` | `#FF5500` | `#FF5500` | Signature brand punch: logo dot, notification badges, and primary alert indicators. |
| `--cf-blue` | `#0284C7` | `#38BDF8` | Technical links, interactive accents, infrastructure badges. |
| `--cf-pastel-bg` | *(sky pastel)* | *(sky pastel dark)* | Active nav item background (sidebar active state). |
| `--cf-pastel-text` | *(sky pastel text)* | *(sky pastel text dark)* | Active nav item foreground text. |
| `--cf-pastel-border` | *(sky pastel border)* | *(sky pastel border dark)* | Active nav item border. |
| `--color-success` | `#2E8555` | `#45A36A` | Operational status, healthy telemetry, verified indicators. |
| `--color-warning` | `#99732B` | `#C29E4D` | Degraded service, memory threshold warnings, pending states. |
| `--color-danger` | `#B84747` | `#CC5454` | Service down, critical alerts, incident records. |

---

## 6. Dense Technical Dashboard & Management Pages

For administrative views, project management, and server settings:
* **Reusable Table Architecture (`$lib/components/ui/`):**
  * `Table.svelte`: Generic typed wrapper with `columns: TableColumn[]`, automatic loading skeleton, empty state handler, and scoped row slot (`let:item`).
  * `TableRow.svelte`: Consistent bottom border separation and hover transitions (`hover:bg-(--bg-muted)/25`).
  * `TableCell.svelte`: Standardized readable scale (`text-sm` padding `px-5 py-3.5`) with support for alignment, vertical alignment, and monospace numeric modes.
* **Readable Information Density & Font Scale:**
  * Headers: `text-sm font-semibold text-(--text-secondary)`.
  * Standard Cell Copy: `text-sm text-(--text-main)`.
  * Technical Telemetry / Hashes / Tokens: `text-xs font-mono`.
  * Action Buttons & Selects: `text-sm font-medium` or `text-xs font-mono` with `radius-sm`.
* **Grid Alignments:**
  * Alphanumeric text & names: **Left-aligned**.
  * Numerical metrics, timestamps, sizes: **Right-aligned** with `font-mono`.
  * Status badges: **Centered** or **Right-aligned**.
* **Restrained Controls:** Use `radius-sm` for action buttons. Use ghost buttons with subtle borders for secondary options.
* **Inline Feedback Over Modals:** Prefer inline validation messages, subtle toast banners, and clear empty states over blocking modals.

---

## 7. Brand Identity Assets

### Logo
- **Files:** `/static/logo.webp` (square mark, 675×675), `/static/logofull.webp` (horizontal full lockup, 1413×471)
- **Usage Pattern:** All navbar and sidebar logos use `logo.webp` (the square mark) beside the text `ngumpulhost.`:
  ```svelte
  <a href="/" class="flex items-center gap-2 hover:opacity-85 transition-opacity">
    <img src="/logo.webp" alt="" class="h-5 w-5 object-contain" />
    <span class="font-sans font-semibold text-sm tracking-tight text-(--text-main)">
      ngumpul<span class="text-(--text-muted)">host</span><span class="text-(--accent-orange)">.</span>
    </span>
  </a>
  ```
- **Public Navbar (transparent hero state):** Text adapts — white on hero image, `text-(--text-main)` on scroll. The dot always stays `text-(--accent-orange)`.
- **Admin Sidebar Collapsed:** Shows only `logo.webp` icon (24×24), morphs to `SidebarSimple` on hover.

### Favicon
- **File:** `/static/favicon.ico` (32×32 ICO, sourced from `asset/logo.ico`)
- **Declared in:** `src/app.html` as `<link rel="icon" href="/favicon.ico" sizes="32x32" type="image/x-icon" />`

### Active Navigation State (Sidebar)
Both admin and user `/me` sidebars use the same active tab style:
```
bg-(--cf-pastel-bg) text-(--cf-pastel-text) font-semibold border border-(--cf-pastel-border)/50 rounded-lg
```
This applies to all nav items in `routes/admin/+layout.svelte` and `routes/me/+layout.svelte`.

### Notification Badge Color
All unread notification counts use `bg-(--accent-orange) text-white` (orange), matching the admin dashboard badge style. Never use `bg-(--accent-sky)` for notification counts.
