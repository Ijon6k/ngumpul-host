# Ngumpul Host Component Directory — Master Component Bible

> **Audience:** AI Coding Agents (Antigravity, Claude Code, Cursor, Copilot) & Human Engineers  
> **Authority:** Canonical Reference. Supersedes all speculative UI generation.  
> **Source Directory:** `frontend/src/lib/components/`  
> **Location:** [`docs/frontend/components.md`](./components.md)

---

## Table of Contents

1. [Core Rules & Design Invariants](#1-core-rules--design-invariants)
   - [Border Radius Hierarchy](#border-radius-hierarchy-mandatory)
   - [Typography Scale Standards](#typography-scale-standards-mandatory)
   - [Icon Standards (Phosphor Svelte Exclusive)](#icon-standards-phosphor-svelte-exclusive)
   - [Anti-AI-Slop Rules for Components](#anti-ai-slop-rules-for-components)
2. [UI Primitives (`$lib/components/ui/`)](#2-ui-primitives-libcomponentsui)
   - [Table, TableRow, TableCell](#table-tablerow-tablecell-centralized-reusable-table-architecture)
   - [Button](#button)
   - [Badge](#badge)
   - [Card](#card)
   - [Input](#input)
   - [Textarea](#textarea)
   - [Select](#select)
   - [Alert](#alert)
   - [ImageUpload](#imageupload)
   - [PageContainer & PageHeader](#pagecontainer--pageheader)
   - [ProjectCover](#projectcover)
   - [Pagination](#pagination)
   - [EmptyState](#emptystate)
   - [Timeline & TimelineItem](#timeline--timelineitem)
   - [Tabs](#tabs)
   - [Breadcrumb, BreadcrumbItem, & BreadcrumbDropdown](#breadcrumb-breadcrumbitem--breadcrumbdropdown)
   - [MarkdownView](#markdownview)
3. [Status & Telemetry Components (`$lib/components/status/`)](#3-status--telemetry-components-libcomponentsstatus)
   - [StatusDot](#statusdot)
   - [UptimeHistory](#uptimehistory)
   - [AvailabilityGrid](#availabilitygrid)
   - [StatusOverviewCard](#statusoverviewcard)
   - [ServiceHealthCard & ServiceStatusRow](#servicehealthcard--servicestatusrow)
   - [IncidentHistoryCard](#incidenthistorycard)
4. [Bento & Hardware Components (`$lib/components/bento/`)](#4-bento--hardware-components-libcomponentsbento)
   - [HardwareCard](#hardwarecard)
   - [MemoryCard](#memorycard)
   - [StorageCard](#storagecard)
   - [NetworkCard](#networkcard)
   - [ServerPhotoCard](#serverphotocard)
   - [CpuPackageVisual](#cpupackagevisual)
5. [Layout & Shell Components (`$lib/components/layout/`)](#5-layout--shell-components-libcomponentslayout)
   - [Header](#header)
   - [Footer](#footer)
   - [ThemeToggle](#themetoggle)
6. [Showcase & Public Feed Components (`$lib/components/`)](#6-showcase--public-feed-components-libcomponents)
   - [ProjectCard](#projectcard)
   - [linkDetector](#linkdetector-libutilslinkdetectorts)
   - [Comments](#comments)
   - [ReportModal](#reportmodal)
   - [ProjectAvailability](#projectavailability)
   - [ActivityTimeline](#activitytimeline)
7. [Admin Console Components (`$lib/components/admin/`)](#7-admin-console-components-libcomponentsadmin)
   - [AdminPanel](#adminpanel)
   - [AdminStatCard](#adminstatcard)
8. [AI Decision Matrix (Anti-Slop Rules)](#8-ai-decision-matrix)

---

## 1. Core Rules & Design Invariants

Before creating or modifying components, adhere strictly to these 4 absolute rules:

### Border Radius Hierarchy (Mandatory)
| Element Category | Token | Tailwind Class | Example Components |
| :--- | :--- | :--- | :--- |
| **Interactive Controls** | `radius-sm` | `rounded-[4px]` / `rounded-sm` | `Button`, `Input`, `Select`, Pagination Buttons |
| **Structural Containers** | `radius-md` | `rounded-[8px]` / `rounded-md` | `Card`, `Table`, `ProjectCover`, `EmptyState` |
| **Floating Overlays** | `radius-lg` | `rounded-xl` / `rounded-2xl` | Modals, Confirmation Dialogs, Popovers |
| **Semantic Badges** | `radius-full` | `rounded-full` | `Badge`, `StatusDot`, Avatars |

> [!CAUTION]
> NEVER mix or randomize radii (e.g., creating `rounded-full` action buttons or `rounded-none` cards).

### Typography Scale Standards (Mandatory)
- **Base Typography Tokens (Tailwind v4 theme tokens):**
  - `--text-xs: 0.875rem` (14px) — metadata, secondary labels, status badges, timestamps.
  - `--text-sm: 1rem` (16px) — body copy, table cells, form inputs, sidebar navigation, card descriptions, breadcrumbs.
  - `--text-base: 1.125rem` (18px) — lead paragraphs, feature summaries.
- **Page Titles:** `text-2xl sm:text-3xl lg:text-4xl font-normal tracking-[-0.03em]`.
- **Section Titles:** `text-lg sm:text-xl font-medium tracking-tight`.
- **Table Column Headers:** `text-sm font-semibold text-(--text-secondary) uppercase tracking-wider`.
- **Table Body Cells:** `text-sm text-(--text-main)`.
- **Technical Metadata / Identifiers:** `text-xs font-mono` (Ports, IPs, Token Hashes, Timestamps).
- **Badges:** `text-xs font-medium px-2.5 py-0.5`.

### Icon Standards (Phosphor Svelte Exclusive)
- **Source:** Strictly `phosphor-svelte`.
- **Weight:** Default `weight="regular"`. Use `weight="bold"` ONLY for directional indicators (`ArrowRight`, `CaretLeft`, `CaretDown`) and confirmation checks (`Check`).
- **Forbidden:** Never embed raw `<svg>` markup or other icon libraries (Lucide, Heroicons, FontAwesome).

### Anti-AI-Slop Rules for Components
1. **No Pill/Chip Spam:** Do not insert floating pill badges above section titles or random text labels. Pills are strictly reserved for semantic statuses (`StatusDot`, `Badge`).
2. **No All-Caps Monospace on Headers:** Never format titles like `SYSTEM_ADMIN_USERS`. Use clean, proportional sans-serif in standard Title Case or sentence case.
3. **Semantic CSS Tokens Only:** Always reference CSS custom properties (`--bg-surface`, `--text-main`, `--border-hairline`, `--accent-sky`). Never insert arbitrary hex values inside component markup.

---

## 2. UI Primitives (`$lib/components/ui/`)

Import all primitives directly from the barrel index:
```svelte
import {
  Table,
  TableRow,
  TableCell,
  Button,
  Badge,
  Card,
  Input,
  PageContainer,
  PageHeader,
  ProjectCover,
  Pagination,
  EmptyState,
  Timeline,
  TimelineItem,
  Tabs,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbDropdown
} from '$lib/components/ui';
```

---

### `Table`, `TableRow`, `TableCell` (Centralized Reusable Table Architecture)

This reusable table architecture serves as the foundation for tabular presentations across the application. All margins, padding, borders, responsive wrapping, loading skeletons, and empty states are centralized here.

#### Column Schema (`TableColumn`)
```typescript
export interface TableColumn {
    key: string;                         // Unique column key / data field
    label: string;                       // Column header title rendered in <thead>
    align?: 'left' | 'center' | 'right'; // Horizontal alignment (default: 'left')
    width?: string;                      // Optional width (e.g., '120px', '25%')
    class?: string;                      // Additional class for both header & cell
    headerClass?: string;                // Additional class specifically for <th>
    cellClass?: string;                  // Additional class specifically for <td>
}
```

#### Props `Table.svelte`
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `columns` | `TableColumn[]` | `[]` | Column definitions array |
| `items` | `T[]` | `undefined` | Data rows rendered in table |
| `loading` | `boolean` | `false` | Shows loading spinner state |
| `loadingMessage` | `string` | `'Loading records...'` | Message displayed while loading |
| `emptyMessage` | `string` | `'No records found.'` | Message displayed when items array is empty |
| `keyField` | `string` | `'id'` | Unique identity field for `{#each (key)}` |
| `alignTop` | `boolean` | `false` | Aligns all cells to the top (`align-top`) |
| `class` | `string` | `''` | Additional class for outer wrapper |

#### Props `TableRow.svelte`
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `class` | `string` | `''` | Additional class for `<tr>` (hover state is built-in) |

#### Props `TableCell.svelte`
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `align` | `'left' \| 'center' \| 'right'` | `'left'` | Horizontal text alignment |
| `alignTop` | `boolean` | `false` | Vertically aligns content to the top (`align-top`) |
| `mono` | `boolean` | `false` | Sets typography to `font-mono text-xs sm:text-sm` |
| `colSpan` | `number` | `undefined` | Populates HTML `colspan` attribute |
| `class` | `string` | `''` | Additional class overrides |

#### Usage Example (Custom Row Slot)
```svelte
<script lang="ts">
  import { Table, TableRow, TableCell, type TableColumn, Button } from '$lib/components/ui';
  import StatusDot from '$lib/components/StatusDot.svelte';

  const columns: TableColumn[] = [
    { label: 'Project', key: 'name' },
    { label: 'Port', key: 'port', width: '120px' },
    { label: 'Status', key: 'status', align: 'center', width: '140px' },
    { label: 'Actions', key: 'actions', align: 'right', width: '160px' }
  ];

  let projects = [
    { id: '1', name: 'Atlas Bot', port: 8081, status: 'ONLINE' },
    { id: '2', name: 'Story Web', port: 3000, status: 'SETUP' }
  ];
</script>

<Table {columns} items={projects} loading={false} emptyMessage="No projects registered.">
  <svelte:fragment let:item={proj}>
    <TableRow>
      <TableCell class="font-medium">{proj.name}</TableCell>
      <TableCell mono>:{proj.port}</TableCell>
      <TableCell align="center">
        <StatusDot status={proj.status} />
      </TableCell>
      <TableCell align="right">
        <Button size="sm">Manage</Button>
      </TableCell>
    </TableRow>
  </svelte:fragment>
</Table>
```

---

### `Button`

Interactive action element conforming strictly to the `radius-sm` (`rounded-[4px]`) token.
- **`variant="primary"`**: Uses calm sky-blue accent (`#0284C7`, hover `#0369A1`) with solid white text for high contrast (WCAG AA 4.54:1) that remains legible across both light and dark themes.
- **`variant="secondary"`**: Uses standardized neutral surfaces (`bg-(--bg-muted)` with border `--border-hairline`).

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `variant` | `'primary' \| 'secondary' \| 'ghost' \| 'danger'` | `'secondary'` | Visual style (`primary` = sky-blue, `secondary` = neutral border) |
| `size` | `'sm' \| 'md' \| 'lg'` | `'sm'` | Button sizing (height: 32px, 36px, 40px) |
| `href` | `string \| undefined` | `undefined` | When provided, automatically renders an `<a>` tag |
| `disabled` | `boolean` | `false` | Disables interaction and dims opacity |
| `loading` | `boolean` | `false` | Displays an inline loading spinner |
| `type` | `'button' \| 'submit' \| 'reset'` | `'button'` | HTML button type |

#### Usage Example
```svelte
<Button variant="primary" on:click={handleSave}>Save changes</Button>
<Button variant="secondary" href="/projects">Browse catalog</Button>
<Button variant="danger" size="sm" loading={isDeleting}>Delete</Button>
```

---

### `Badge`

Semantic metadata indicator with `radius-full` token. Used for status, type tags, or categories.

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `variant` | `'success' \| 'warning' \| 'danger' \| 'sky' \| 'neutral'` | `'neutral'` | Semantic color palette |
| `size` | `'sm' \| 'md'` | `'sm'` | Text and padding size |
| `dot` | `boolean` | `false` | Renders a small circular dot on the left side |

#### Usage Example
```svelte
<Badge variant="success" dot>Operational</Badge>
<Badge variant="sky">Docker</Badge>
```

---

### `Card`

Structural content container with `radius-md` (`rounded-[8px]`) token.

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `variant` | `'surface' \| 'muted' \| 'interactive'` | `'surface'` | Background surface & hover behavior |
| `padding` | `'none' \| 'sm' \| 'md' \| 'lg'` | `'md'` | Standardized padding scale |
| `href` | `string \| undefined` | `undefined` | When provided, renders as an interactive link card (`<a>`) |

#### Usage Example
```svelte
<Card variant="surface" padding="md">
  <h3 class="text-base font-medium">Server Specification</h3>
  <p class="text-sm text-(--text-secondary)">Bare-metal hardware details.</p>
</Card>
```

---

### `Input`

Form input element with `radius-sm` (`rounded-[4px]`) token, automated label handling, helper text, and error states. Passes through arbitrary attributes via `{...$$restProps}`, so Superforms `$constraints` (e.g. `minlength`, `maxlength`, `pattern`, `required`) and `aria-invalid` spread onto the underlying `<input>` work directly with zero component changes.

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `value` | `string \| number` | `''` | Input value (bindable) |
| `type` | `string` | `'text'` | HTML input type |
| `label` | `string` | `''` | Header label above input |
| `placeholder` | `string` | `''` | Placeholder text |
| `helperText` | `string` | `''` | Descriptive guidance below input |
| `error` | `string` | `''` | Validation error message (rendered in danger color) |
| `disabled` | `boolean` | `false` | Disables interaction |
| `readonly` | `boolean` | `false` | Read-only mode |

#### Usage Example
```svelte
<Input
  label="Repository URL"
  placeholder="https://github.com/user/repo"
  bind:value={repoUrl}
  error={formErrors.repoUrl}
  helperText="Public repository only."
/>
```

---

### `Textarea`

Multi-line text entry primitive conforming to `radius-sm` (`rounded-[4px]`), supporting optional automated labels, character counter or extra metadata via `label-extra` slot, error states, and helper text.

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `value` | `string` | `''` | Textarea value (bindable) |
| `label` | `string` | `''` | Header label above textarea |
| `placeholder` | `string` | `''` | Placeholder guidance text |
| `rows` | `number` | `3` | Visible text lines count |
| `helperText` | `string` | `''` | Subtle text rendered beneath the field |
| `error` | `string` | `''` | Validation error message |
| `disabled` | `boolean` | `false` | Disables interaction |
| `readonly` | `boolean` | `false` | Disallows text input |
| *Slots* | `label-extra` | | Slot for right-aligned label extras (e.g. character counter) |

#### Usage Example
```svelte
<Textarea
  label="Project Description"
  rows={3}
  maxlength={280}
  bind:value={desc}
  placeholder="Short summary..."
>
  <svelte:fragment slot="label-extra">
    <span class="text-xs font-mono text-(--text-muted)">{desc.length}/280</span>
  </svelte:fragment>
</Textarea>
```

---

### `Select`

Drop-down selection primitive styled with `radius-sm` (`rounded-[4px]`), custom Phosphor `CaretDown` adornment, and full dark-theme token adherence.

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `value` | `any` | `''` | Selected option value (bindable) |
| `label` | `string` | `''` | Header label above dropdown |
| `helperText` | `string` | `''` | Descriptive text below dropdown |
| `error` | `string` | `''` | Validation error message |
| `disabled` | `boolean` | `false` | Disables interaction |
| *Slots* | `default` | | `<option>` elements |

#### Usage Example
```svelte
<Select label="Hosting Runtime" bind:value={runtime}>
  <option value="DOCKER">Docker Container</option>
  <option value="STATIC">Static HTML / SPA</option>
</Select>
```

---

### `Alert`

Contextual inline alert and notification banner with `radius-md` (`rounded-[8px]`), semantic border accents, and Phosphor icons (`CheckCircle`, `WarningCircle`, `XCircle`, `Info`).

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `variant` | `'danger' \| 'warning' \| 'success' \| 'info'` | `'info'` | Alert visual severity and coloring |
| `title` | `string` | `''` | Optional bold header text |
| `showIcon` | `boolean` | `true` | Controls whether leading Phosphor icon appears |

#### Usage Example
```svelte
<Alert variant="danger">
  Invalid credentials provided. Please try again.
</Alert>

<Alert variant="warning" title="Notice">
  Registrations are currently closed on this host.
</Alert>
```

---

### `ImageUpload`

Centralized project artwork cover uploader conforming to the **16:9 standard ratio**. Features client-side WebP compression (1600x1200 max, 0.82 quality), dashed dropzone, hover actions, and direct integration with `projectsApi.uploadCover`.

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `value` | `string` | `''` | Cover image URL (bindable) |
| `label` | `string` | `'Cover Artwork'` | Header label |
| `helperText` | `string` | `'(Optional · 16:9 standard)'` | Guidance text |
| `projectName` | `string` | `''` | Fallback initials title for preview |
| `aspectRatio` | `string` | `'16/9'` | Strict cover aspect ratio |
| `disabled` | `boolean` | `false` | Disables file selection and upload |

#### Usage Example
```svelte
<ImageUpload
  bind:value={coverImageUrl}
  projectName={projectName}
  label="Project Cover"
/>
```

---

### `Skeleton`

Dependency-free loading primitive that mirrors the **exact visual structure** of the real UI. Built with the `--bg-muted` token and `animate-pulse`; never use it in isolation, always compose it inside the same containers the loaded state uses so there is **zero layout shift**.

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `variant` | `'block' \| 'line' \| 'avatar' \| 'badge' \| 'button'` | `'block'` | Shape: `line` = single text row (`h-3.5`), `avatar` = circle, `badge` = pill, `button` = action stub |
| `animate` | `boolean` | `true` | Disable the pulse animation (e.g. inside already-animated parents) |
| `ariaLabel` | `string` | `''` | When set, switched to `role="status"` for screen readers |
| `class` | `string` | `''` | Size/width overrides via `cn()` (tailwind-merge resolves conflicts) |

#### Usage Example — Card content skeleton
```svelte
<script>
  import { Skeleton } from '$lib/components/ui';
</script>

<div class="flex items-center gap-3" role="status" aria-live="polite">
  <Skeleton variant="avatar" />
  <div class="flex flex-col gap-2 flex-1">
    <Skeleton variant="line" class="w-1/2" />
    <Skeleton variant="line" class="w-2/3" />
  </div>
</div>
<span class="sr-only">Loading...</span>
```

#### Rules
1. **Mirror the real layout:** use the same card/grid/padding classes as the loaded state; only swap data elements for `<Skeleton>`.
2. **Never plain text loaders:** `Loading...` text alone is forbidden; keep a `sr-only` message for a11y, always set `role="status"` on the scaffold.
3. **Wide variety via `line` widths:** vary width (`w-1/2`, `w-2/3`, `w-full`) to mimic natural text lengths; never equal-width stripes everywhere.
4. **Composed helpers:** `SkeletonProjectCard` (`$lib/components/SkeletonProjectCard.svelte`) mirrors `ProjectCard` and is the canonical catalog/landing loader.

---

### `PageContainer` & `PageHeader`

Standardized layout containers for sub-views to eliminate layout shifts and inconsistent padding.

#### Props `PageContainer`
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `maxWidth` | `'4xl' \| '5xl' \| '6xl' \| 'full'` | `'6xl'` | Maximum content container width |
| `paddingY` | `string` | `'py-8 sm:py-12 pb-24'` | Vertical pacing classes |

#### Props `PageHeader`
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `title` | `string` | *(required)* | Main page title |
| `description` | `string` | `''` | Subtitle / brief editorial overview |
| *Slots* | `breadcrumb`, `actions` | | Slots for navigation breadcrumbs and action buttons |

#### Usage Example
```svelte
<PageContainer maxWidth="5xl">
  <PageHeader title="Members" description="The independent developers hosting here.">
    <svelte:fragment slot="actions">
      <Button variant="primary" href="/invite">Invite Member</Button>
    </svelte:fragment>
  </PageHeader>
  
  <!-- Content -->
</PageContainer>
```

---

### `ProjectCover`

Renders visual project artwork **full-bleed** with 4:3, 16:9, or flexible container ratios (`aspectRatio="full"`).
- **Full-Bleed Fallback:** When a project does not have an uploaded cover image, this component renders a solid, calm pastel blue surface (`bg-sky-100/80 dark:bg-sky-950/40` with `border-sky-300/60 dark:border-sky-800/60`) filling the entire card area without letterboxing or shrinking.
- **Initial Monogram:** Displays 1–2 initial letters of the project title centered proportionally.
- **Aspect Ratio Normalization:** Supports `'4/3'`, `'16/9'`, `'aspect-[4/3]'`, or `'full'` (`w-full h-full min-h-full`).

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `src` | `string \| null` | `undefined` | Cover image URL |
| `name` | `string` | `'Project'` | Project name for monogram generation |
| `aspectRatio` | `string` | `'full'` | Aspect ratio string (`'full'`, `'4/3'`, `'16/9'`, or Tailwind class) |
| `class` | `string` | `''` | Additional classes for outer wrapper |

---

### `Pagination`

Tabular page navigation control with automated page window calculation, Previous/Next labeled buttons, and Phosphor icon controls.

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `currentPage` | `number` | `1` | Currently active page index (1-based) |
| `totalItems` | `number` | `0` | Total record count |
| `pageSize` | `number` | `12` | Number of records per page |
| `onPageChange` | `(page: number) => void` | *(required)* | Callback triggered upon page change |
| `variant` | `'full' \| 'compact' \| 'simple'` | `'full'` | Full numbered buttons or compact `Page X of Y` indicator |
| `showLabels` | `boolean` | `true` | Show 'Previous' and 'Next' text labels alongside carets |
| `showInfo` | `boolean` | `true` | Show 'Showing X–Y of Z' counter description |
| `hideOnSinglePage` | `boolean` | `true` | Automatically hide navigation bar when total pages is $\le 1$ |

---

### `EmptyState`

Standardized presentation when tables, lists, or feeds contain zero records.

#### Props & Slots
| Prop / Slot | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `title` | `string` | `'Nothing here yet'` | Main empty state heading |
| `description` | `string` | `''` | Clarifying guidance explaining the empty state |
| `slot="icon"` | Slot | Default dot | Custom icon above title |
| `slot="action"` | Slot | | Call-to-action button to initiate new records |

---

### `Timeline` & `TimelineItem`

Standardized vertical timeline primitive for project history, workspace activity, audit logs, and public feeds. Replaces repetitive table structures or nested card patterns with crisp connector lines and semantic status dots.

#### Props `Timeline.svelte`
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `density` | `'comfortable' \| 'compact'` | `'comfortable'` | Vertical spacing density between items |

#### Props `TimelineItem.svelte`
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `title` | `string` | `''` | Event title (or via slot `title`) |
| `timestamp` | `string` | `''` | Relative timestamp / formatted date (or via slot `timestamp`) |
| `description` | `string` | `''` | Summary detail text |
| `dotColor` | `string` | `''` | Tailwind dot color class (e.g., `bg-emerald-500`) |
| `status` | `string` | `''` | Automated dot color resolution via status (`ONLINE`, `OFFLINE`) |
| `density` | `'comfortable' \| 'compact'` | `'comfortable'` | Vertical padding density |
| `href` | `string` | `''` | Optional target URL (e.g., `/projects/slug`) |
| `linkText` | `string` | `''` | Clickable anchor label |

---

### `Tabs`

Minimalist tab primitive based on clean typography and subtle underlines (never thick capsules or pill tabs). Used in the Project Workspace (`Overview | Visits | Activity | Settings`) and Public Project Detail (`Overview | Activity`).

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `tabs` | `Array<{ id: string; label: string; icon?: any; badge?: number }>` | `[]` | List of navigation tab objects |
| `active` | `string` | *(first tab)* | Active tab key (bindable: `bind:active={...}`) |
| `size` | `'sm' \| 'default'` | `'default'` | Typography and padding size scale |

---

### `Breadcrumb`, `BreadcrumbItem`, & `BreadcrumbDropdown`

Standardized hierarchical navigation primitive calibrated to **`text-sm` (16px / `1rem`)** across the application. Replaces ad-hoc `<nav>` tags to preserve typography consistency, spacing, contrast, and ARIA accessibility in the admin console, user workspace, and public showcase.

#### Usage Patterns:
1. **Quick Declarative (`Breadcrumb` with `items` prop):** Ideal for standard 1–3 level navigation paths.
2. **Slot Composition (`Breadcrumb` + `BreadcrumbItem`):** Ideal for dynamic layouts or embedding interactive switchers such as `BreadcrumbDropdown`.

#### Props `Breadcrumb.svelte`
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `items` | `BreadcrumbItemData[]` | `[]` | Array of breadcrumb items `{ label, href?, icon?, current? }` |
| `ariaLabel` | `string` | `'Breadcrumb'` | Accessible ARIA label for `<nav>` |
| `class` | `string` | `''` | Additional Tailwind classes for container (e.g., `mb-4`, `mb-6`) |

#### Props `BreadcrumbItem.svelte`
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `href` | `string \| undefined` | `undefined` | Target link URL. If omitted or `current={true}`, renders as active `<span>` |
| `current` | `boolean` | `false` | Marks item as active page (`aria-current="page"`, `font-medium`) |
| `icon` | `any` | `undefined` | Optional Phosphor icon component (e.g., `CaretLeft`) |
| `separator` | `boolean` | `true` | Renders preceding slash `/` separator. Set to `false` for first child |
| `label` | `string` | `''` | Alternative text label if not using default slot |

#### Props `BreadcrumbDropdown.svelte`
Switcher component embedded inside breadcrumb as an active terminal item.
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `label` | `string` | *(required)* | Label text displayed on the trigger button |
| `activeId` | `string` | `''` | ID of currently active item (marked with `Check` icon) |
| `options` | `Array<{ id: string; label: string }>` | `[]` | List of selectable options |
| `listboxLabel` | `string` | `'Switch item'` | Accessible label for flyout listbox |
| `dropdownHeader` | `string` | same as `listboxLabel` | Small header text at top of listbox |
| *Events* | `select` | `{ id: string }` | Emitted when an option is selected |

#### Usage Example

```svelte
<script>
  import { CaretLeft } from 'phosphor-svelte';
  import { Breadcrumb, BreadcrumbItem, BreadcrumbDropdown } from '$lib/components/ui';
</script>

<!-- Pattern 1: Quick Declarative -->
<Breadcrumb
  items={[
    { label: 'Projects', href: '/projects' },
    { label: project.name }
  ]}
  class="mb-6"
/>

<!-- Pattern 2: With Back Icon & BreadcrumbDropdown Switcher -->
<Breadcrumb class="mb-4">
  <BreadcrumbItem href="/me/projects" icon={CaretLeft} separator={false}>
    Projects
  </BreadcrumbItem>
  <BreadcrumbItem current>
    <BreadcrumbDropdown
      label={project.name}
      activeId={project.id}
      options={allProjects}
      on:select={(e) => goto(`/me/projects/${e.detail.id}`)}
    />
  </BreadcrumbItem>
</Breadcrumb>
```

---

### `MarkdownView`

A GitHub-flavored markdown renderer component powered by `marked`. Supports an optional outer marker container with a top header strip (`title="README.md"`) and thin outer border to clearly differentiate user-authored markdown content from platform UI chrome.

#### Props
| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `content` | `string` | `''` | Raw markdown string to parse and render |
| `title` | `string` | `''` | Optional header title (e.g. `'README.md'`). Triggers top header strip with `FileText` icon |
| `bordered` | `boolean` | `false` | When true (or when `title` is set), wraps rendered content in a thin marker border container (`radius-md`) |
| `class` | `string` | `''` | Additional classes applied to inner markdown container |

#### Usage Example
```svelte
<MarkdownView
  content={project.readme}
  title="README.md"
  bordered={true}
/>
```

---

## 3. Status & Telemetry Components (`$lib/components/status/`)

Dedicated components for system availability and telemetry visualizations:

- **`StatusDot.svelte`**:
  - Calm operational indicator (`ONLINE`, `SETUP`, `PENDING`, `OFFLINE`, `ARCHIVED`).
  - Supports `variant="inline"` (default: dot + text without pill wrapper), `variant="dot"` (circle dot only), and `variant="badge"` (semantic pill).
  - Supports `showLabel={true | false}`. When rendering inside tables or lists alongside custom colored text spans (`resolved.colorClass`), always set `showLabel={false}` to prevent duplicate uncolored status text.
- **`UptimeHistory.svelte`**:
  - Unified multi-range discrete vertical uptime bar visualization (`1d` / 24 or 48 bars, `7d` / 56 bars at 3-hour intervals, `30d` / 30 bars).
  - Binds strictly to real probe telemetry without synthetic data fabrication; dates prior to project deployment are non-penalized (`no_data` neutral blocks).
  - Anti-slop empty state: When no telemetry records exist yet (`total_checks == 0` or all bars `no_data`), displays an editorial calm pending notice instead of empty spam bars and dashes.
  - Interactive tooltip upon hover displaying exact time window, percentage, and probe check count ($X/Y$).
  - Supports range toggle buttons (`24h`, `7d`, `30d`) with instant client-side switching across pre-calculated range buckets.
  - Props: `availability` (multi-range payload), `blocks` (direct block array), `status`, `range`, `variant` (`'compact' | 'detailed'`), `showRangeSelector`.
- **`AvailabilityGrid.svelte`**:
  - High-resolution host availability vertical bar sequence replacing legacy square grids.
  - Slices observation windows into 24 bars (`1d`), 56 bars (`7d`), or 30 bars (`30d`).
  - Displays detailed percentage and incident duration tooltip upon hover.
- **`StatusOverviewCard.svelte`**:
  - Global status banner showing trailing availability window percentage and formatted host uptime.
- **`ServiceHealthCard.svelte` & `ServiceStatusRow.svelte`**:
  - Renders list of internal infrastructure components (Web Ingress, PostgreSQL, Telemetry Probes) and edge latency.
- **`IncidentHistoryCard.svelte`**:
  - Recapitulation of real downtime incident records ($\ge 120$ seconds) with root causes (`SYSTEM_REBOOT` or `SERVICE_SUSPENDED`).

---

## 4. Bento & Hardware Components (`$lib/components/bento/`)

Dedicated to physical host hardware specification showcase on landing page (`+page.svelte`):

- **`HardwareCard.svelte`**: Displays Intel/AMD processor details, physical core count, logical threads, and x86_64 architecture.
- **`MemoryCard.svelte`**: Displays physical RAM (`MemTotal` vs `MemAvailable`) and utilization bar.
- **`StorageCard.svelte`**: Displays physical NVMe drive models (`/sys/block/nvme*`) and Linux `/` partition capacity.
- **`NetworkCard.svelte`**: Displays uplink bandwidth and edge round-trip latency to Cloudflare 1.1.1.1.
- **`ServerPhotoCard.svelte`**: Tactile editorial visual card showcasing the real homelab server node in Jakarta.
- **`CpuPackageVisual.svelte`**: Visual die package diagram representing the processor silicon.

---

## 5. Layout & Shell Components (`$lib/components/layout/`)

- **`Header.svelte`**:
  - Sticky top header with standardized `60px` height.
  - Navigation links: `/projects`, `/people`, `/activity`, `/status`.
  - User session state (Login action or User profile dropdown).
  - Responsive mobile drawer.
- **`Footer.svelte`**:
  - Editorial footer with homelab node identity, physical server location, and documentation links.
- **`ThemeToggle.svelte`**:
  - Clean light/dark mode switch without excessive borders, using Phosphor `Sun` and `Moon` icons.

---

## 6. Showcase & Public Feed Components (`$lib/components/`)

- **`ProjectCard.svelte`**:
  - Standard public catalog card for `/projects` and the landing page showcase (`/+page.svelte`).
  - **16:9 Artwork Ratio:** Uses a strict 16:9 aspect container (`aspect-[16/9]`) with `ProjectCover` (`aspectRatio="full"`).
  - **Corner-Flush Minimalist Status Badge:** Positioned flush in the top-right corner of the image (`top-0 right-0`) with subtle `rounded-bl-md` curvature. Uses theme-adaptive neutral surface backdrop (`bg-(--bg-surface)/90 dark:bg-neutral-900/90 backdrop-blur-xs`) and hairline borders (`border-b border-l border-(--border-hairline)`). Renders direct semantic colored text (`{projectStatus.label}`) **without** pill shapes, colored borders, or dot circles.
  - **Strict 2-Line Description Clamp:** Implements `line-clamp-2 min-h-[2.5rem] sm:min-h-[2.85rem] leading-relaxed font-normal`. Guarantees exactly 2 lines maximum without clipping font descenders ('y', 'p', 'g', 'q') and preserves vertical alignment across cards.
  - **Grid & Dimension Parity:** Standardized in a 2-column grid (`grid-cols-1 md:grid-cols-2 gap-6 sm:gap-8`) inside a `max-w-5xl` container. The landing page showcases up to 4 items in a clean 2x2 grid, perfectly matching the catalog card dimensions 1:1.
  - Displays title, concise summary, creator/owner, and technology stack (separated by `·`).
  - Supports direct external visit link connected to `/go/{slug}` tracking route.
- **`linkDetector` (`$lib/utils/linkDetector.ts`)**:
  - Automated URL domain detector for project resource links:
    - `github.com` $\rightarrow$ `GithubLogo` icon, label "GitHub"
    - `gitlab.com` $\rightarrow$ `GitlabLogo` icon, label "GitLab"
    - `drive.google.com` $\rightarrow$ `GoogleDriveLogo` icon, label "Google Drive"
    - `docs.google.com` $\rightarrow$ `FileText` icon, label "Google Docs"
    - `figma.com` $\rightarrow$ `FigmaLogo` icon, label "Figma"
    - `youtube.com` / `youtu.be` / `vimeo.com` $\rightarrow$ `YoutubeLogo` icon, label "Video Demo"
    - `twitter.com` / `x.com` $\rightarrow$ `TwitterLogo` icon, label "X (Twitter)"
    - `discord.gg` / `discord.com` $\rightarrow$ `DiscordLogo` icon, label "Discord"
    - Other arbitrary domains $\rightarrow$ `LinkSimple` icon with intelligent fallback label.
- **`Comments.svelte`**:
  - Chronological discussion thread on public project showcase (`/projects/[slug]`).
  - Integrated with `commentsApi` domain module (`getProjectComments`, `createComment`, `deleteComment`).
  - Comment composer with 5/10 min rate limiting, input sanitization, author avatar, and display name.
  - Soft-delete handling (*"This comment was removed"*).
  - Contextual action menus: Report modal trigger and Delete with `ConfirmModal` confirmation (restricted to author, project owner, or admin).
- **`ReportModal.svelte`**:
  - Content moderation and violation reporting dialog for Projects and Comments.
  - Integrated with `reportsApi.submitReport`.
  - Standardized reason dropdown, optional detail text, and calm blurred backdrop.
- **`ProjectAvailability.svelte`**:
  - Project health telemetry widget driven by 5-minute HTTP probes.
  - Displays 30-day uptime in %, average round-trip latency, and real-time operational status.
- **`ActivityTimeline.svelte`**:
  - Chronological public activity feed (project publications, updates, verifications) backed by `activityApi`.

---

## 7. Admin Console Components (`$lib/components/admin/`)

- **`AdminPanel.svelte`**: Operator sub-navigation bar (`Overview`, `Projects`, `Requests`, `Users`, `Audit`, `System`, `Settings`).
- **`AdminStatCard.svelte`**: Metric card summarizing key operational statistics.

---

## 8. AI Decision Matrix

When AI Coding Agents create or modify user interfaces, follow this decision matrix:

| Scenario | MANDATORY Component | FORBIDDEN Pattern |
| :--- | :--- | :--- |
| **Displaying Tabular Data** | `<Table>`, `<TableRow>`, `<TableCell>` from `$lib/components/ui` | ❌ Raw `<table>`, `<th>`, `<td>` tags with ad-hoc styles |
| **Rendering Buttons** | `<Button variant="..." size="...">` | ❌ `<button class="rounded-full bg-gradient-to-r ...">` |
| **Project/Node Operational Status** | `<StatusDot status={...} />` | ❌ Handcrafted chip badges with random colors |
| **Form Inputs** | `<Input label="..." bind:value={...} />` | ❌ `<input class="border p-2 ...">` without design tokens |
| **Content Container / Bento Panel** | `<Card variant="surface" padding="md">` | ❌ `div` wrappers with `rounded-2xl` or arbitrary colors |
| **Project Imagery** | `<ProjectCover src={...} name={...} />` | ❌ Bare `<img>` tags without fallbacks or monogram handling |
| **Interface Icons** | Imported from `phosphor-svelte` | ❌ Raw `<svg>` markup or other libraries (Lucide, Heroicons) |
| **Breadcrumbs** | `<Breadcrumb>` / `<BreadcrumbItem>` | ❌ Inline ad-hoc `<nav aria-label="Breadcrumb">` tags |
