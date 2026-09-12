# Kitab Komponen Ngumpul Host — Master Component Bible

> **Audience:** AI Coding Agents (Antigravity, Claude Code, Cursor, Copilot) & Human Engineers  
> **Authority:** Canonical Reference. Supersedes all speculative UI generation.  
> **Source Directory:** `frontend/src/lib/components/`  
> **Location:** [`docs/frontend/components.md`](./components.md)

---

## Daftar Isi (Table of Contents)

1. [Aturan Inti & Design Invariants](#1-aturan-inti--design-invariants)
   - [Hierarki Border Radius](#hierarki-border-radius-wajib)
   - [Standar Skala Tipografi](#standar-skala-tipografi-wajib)
   - [Standar Ikon (Phosphor Svelte)](#standar-ikon-phosphor-svelte-eksklusif)
   - [Aturan Anti-AI-Slop untuk Komponen](#aturan-anti-ai-slop-untuk-komponen)
2. [UI Primitives (`$lib/components/ui/`)](#2-ui-primitives-libcomponentsui)
   - [Table, TableRow, TableCell](#table-tablerow-tablecell-tabel-reusable-arsitektur-terpusat)
   - [Button](#button)
   - [Badge](#badge)
   - [Card](#card)
   - [Input](#input)
   - [PageContainer](#pagecontainer)
   - [PageHeader](#pageheader)
   - [ProjectCover](#projectcover)
   - [Pagination](#pagination)
   - [EmptyState](#emptystate)
   - [Timeline & TimelineItem](#timeline--timelineitem)
   - [Tabs](#tabs)
3. [Status & Telemetry Components (`$lib/components/status/`)](#3-status--telemetry-components-libcomponentsstatus)
   - [StatusDot](#statusdot)
   - [StatusIndicator](#statusindicator)
   - [StatusOverviewCard](#statusoverviewcard)
   - [AvailabilityGrid](#availabilitygrid)
   - [ServiceHealthCard & ServiceStatusRow](#servicehealthcard--servicestatusrow)
   - [MetricItem](#metricitem)
   - [IncidentHistoryCard](#incidenthistorycard)
   - [StatusHeader](#statusheader)
4. [Bento & Hardware Components (`$lib/components/bento/`)](#4-bento--hardware-components-libcomponentsbento)
   - [HardwareCard](#hardwarecard)
   - [MemoryCard](#memorycard)
   - [StorageCard](#storagecard)
   - [NetworkCard](#networkcard)
   - [AvailabilityCard](#availabilitycard)
   - [ServerPhotoCard](#serverphotocard)
   - [CpuPackageVisual](#cpupackagevisual)
   - [IdeasCard](#ideascard)
5. [Layout & Shell Components (`$lib/components/layout/`)](#5-layout--shell-components-libcomponentslayout)
   - [Header](#header)
   - [Footer](#footer)
   - [ThemeToggle](#themetoggle)
6. [Showcase & Public Feed Components (`$lib/components/`)](#6-showcase--public-feed-components-libcomponents)
   - [ProjectCard](#projectcard)
   - [HeroSection](#herosection)
   - [BentoGrid](#bentogrid)
   - [MetricPillar & ResourceBar](#metricpillar--resourcebar)
   - [ActivityTimeline](#activitytimeline)
7. [Admin Console Components (`$lib/components/admin/`)](#7-admin-console-components-libcomponentsadmin)
   - [AdminPanel](#adminpanel)
   - [AdminStatCard](#adminstatcard)
8. [Panduan Keputusan AI (Decision Matrix & Anti-Slop Rules)](#8-panduan-keputusan-ai-decision-matrix)

---

## 1. Aturan Inti & Design Invariants

Sebelum membuat atau memodifikasi komponen, pahami 4 aturan absolut berikut:

### Hierarki Border Radius (Wajib)
| Kategori Elemen | Token | Class Tailwind | Contoh Komponen |
| :--- | :--- | :--- | :--- |
| **Interactive Controls** | `radius-sm` | `rounded-[4px]` / `rounded-sm` | `Button`, `Input`, `Select`, Pagination Buttons |
| **Structural Containers** | `radius-md` | `rounded-[8px]` / `rounded-md` | `Card`, `Table`, `ProjectCover`, `EmptyState` |
| **Floating Overlays** | `radius-lg` | `rounded-xl` / `rounded-2xl` | Modals, Confirmation Dialogs, Popovers |
| **Semantic Badges** | `radius-full` | `rounded-full` | `Badge`, `StatusDot`, Avatars |

> [!CAUTION]
> DILARANG mengacak radius (misal membuat tombol `rounded-full` atau card `rounded-none`).

### Standar Skala Tipografi (Wajib)
- **Judul Halaman:** `text-2xl sm:text-3xl lg:text-4xl font-normal tracking-[-0.03em]`.
- **Judul Sub-bagian:** `text-lg sm:text-xl font-medium tracking-tight`.
- **Header Kolom Tabel:** `text-sm font-semibold text-(--text-secondary) uppercase tracking-wider`.
- **Isi Sel Tabel:** `text-sm text-(--text-main)`.
- **Metadata Teknis / Identifiers:** `text-xs font-mono` (Ports, IP, Token Hashes, Timestamps).
- **Badges:** `text-xs font-medium px-2.5 py-0.5`.

### Standar Ikon (Phosphor Svelte Eksklusif)
- **Sumber:** Strictly `phosphor-svelte`.
- **Weight:** Default `weight="regular"`. Gunakan `weight="bold"` HANYA untuk directional arrows (`ArrowRight`, `CaretLeft`) dan konfirmasi (`Check`).
- **Dilarang:** Jangan pernah menyematkan raw `<svg>` atau library lain (Lucide, Heroicons, FontAwesome).

### Aturan Anti-AI-Slop untuk Komponen
1. **Tidak Ada Pill/Chip Spam:** Jangan menambahkan chip badge di atas judul section atau di sembarang label. Pill hanya untuk status semantik (`StatusDot`, `Badge`).
2. **Tidak Ada All-Caps Monospace pada Header:** Jangan pernah membuat judul seperti `SYSTEM_ADMIN_USERS`. Gunakan natural sentence/title case sans-serif.
3. **Warna Semantic Token Saja:** Gunakan variabel CSS (`--bg-surface`, `--text-main`, `--border-hairline`, `--accent-sky`). Jangan gunakan arbitrary hex di dalam template komponen.

---

## 2. UI Primitives (`$lib/components/ui/`)

Impor seluruh primitif langsung dari index:
```svelte
import { Table, TableRow, TableCell, Button, Badge, Card, Input, PageContainer, PageHeader, ProjectCover, Pagination, EmptyState } from '$lib/components/ui';
```

---

### `Table`, `TableRow`, `TableCell` (Tabel Reusable Arsitektur Terpusat)

Sistem tabel reusable ini adalah pondasi tampilan tabular aplikasi. Seluruh margin, padding, borders, responsive wrapping, loading skeleton, dan empty states diatur terpusat di sini.

#### Schema Kolom (`TableColumn`)
```typescript
export interface TableColumn {
    key: string;               // Key unik kolom / field data
    label: string;             // Label judul kolom di <thead>
    align?: 'left' | 'center' | 'right'; // Alignment (default: 'left')
    width?: string;            // Contoh: '120px', '25%'
    class?: string;            // Class tambahan untuk header & cell
    headerClass?: string;      // Class tambahan khusus <th>
    cellClass?: string;        // Class tambahan khusus <td>
}
```

#### Props `Table.svelte`
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `columns` | `TableColumn[]` | `[]` | Array definisi kolom |
| `items` | `T[]` | `undefined` | Data baris yang dirender |
| `loading` | `boolean` | `false` | Menampilkan spinner loading state |
| `loadingMessage`| `string` | `'Loading records...'` | Pesan loading |
| `emptyMessage` | `string` | `'No records found.'` | Pesan jika data kosong |
| `keyField` | `string` | `'id'` | Field identitas unik untuk `{#each (key)}` |
| `alignTop` | `boolean` | `false` | Menyelaraskan seluruh sel ke atas (`align-top`) |
| `class` | `string` | `''` | Class tambahan untuk kontainer luar |

#### Props `TableRow.svelte`
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `class` | `string` | `''` | Class tambahan untuk `<tr>` (hover state sudah bawaan) |

#### Props `TableCell.svelte`
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `align` | `'left' \| 'center' \| 'right'` | `'left'` | Perataan horizontal teks |
| `alignTop` | `boolean` | `false` | Menyelaraskan vertikal ke atas (`align-top`) |
| `mono` | `boolean` | `false` | Mengubah font menjadi `font-mono text-xs sm:text-sm` |
| `colSpan` | `number` | `undefined` | Mengisi atribut HTML `colspan` |
| `class` | `string` | `''` | Class override tambahan |

#### Contoh Penggunaan (Custom Row Slot)
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
    { id: '2', name: 'Cerita Web', port: 3000, status: 'SETUP' }
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

Komponen aksi interaktif dengan kepatuhan token `radius-sm` (`rounded-[4px]`).

#### Props
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `variant` | `'primary' \| 'secondary' \| 'ghost' \| 'danger'` | `'secondary'` | Visual style |
| `size` | `'sm' \| 'md' \| 'lg'` | `'sm'` | Ukuran tombol (tinggi: 32px, 36px, 40px) |
| `href` | `string \| undefined` | `undefined` | Jika diisi, otomatis merender tag `<a>` |
| `disabled` | `boolean` | `false` | Mematikan interaksi |
| `loading` | `boolean` | `false` | Menampilkan spinner di dalam tombol |
| `type` | `'button' \| 'submit' \| 'reset'` | `'button'` | Atribut tipe HTML |

#### Contoh Penggunaan
```svelte
<Button variant="primary" on:click={handleSave}>Save changes</Button>
<Button variant="secondary" href="/projects">Browse catalog</Button>
<Button variant="danger" size="sm" loading={isDeleting}>Delete</Button>
```

---

### `Badge`

Indikator semantik metadata dengan token `radius-full`. Digunakan untuk status, tag tipe, atau kategori.

#### Props
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `variant` | `'success' \| 'warning' \| 'danger' \| 'sky' \| 'neutral'` | `'neutral'` | Palet semantik |
| `size` | `'sm' \| 'md'` | `'sm'` | Ukuran teks & padding |
| `dot` | `boolean` | `false` | Menampilkan titik lingkaran warna di sisi kiri |

#### Contoh Penggunaan
```svelte
<Badge variant="success" dot>Operational</Badge>
<Badge variant="sky">Docker</Badge>
```

---

### `Card`

Kontainer struktural dengan token `radius-md` (`rounded-[8px]`).

#### Props
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `variant` | `'surface' \| 'muted' \| 'interactive'` | `'surface'` | Background & hover state |
| `padding` | `'none' \| 'sm' \| 'md' \| 'lg'` | `'md'` | Pilihan padding terstandar |
| `href` | `string \| undefined` | `undefined` | Jika diisi, otomatis menjadi link card (`<a>`) |

#### Contoh Penggunaan
```svelte
<Card variant="surface" padding="md">
  <h3 class="text-base font-medium">Server Specification</h3>
  <p class="text-sm text-(--text-secondary)">Bare-metal hardware details.</p>
</Card>
```

---

### `Input`

Elemen form input dengan token `radius-sm` (`rounded-[4px]`), penanganan label otomatis, helper text, dan error state.

#### Props
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `value` | `string \| number` | `''` | Nilai input (bindable) |
| `type` | `string` | `'text'` | Tipe input HTML |
| `label` | `string` | `''` | Label di atas input |
| `placeholder`| `string` | `''` | Teks placeholder |
| `helperText` | `string` | `''` | Penjelasan bantuan di bawah input |
| `error` | `string` | `''` | Pesan error validasi (warna merah) |
| `disabled` | `boolean` | `false` | Menonaktifkan input |
| `readonly` | `boolean` | `false` | Mode hanya baca |

#### Contoh Penggunaan
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

### `PageContainer` & `PageHeader`

Standarisasi tata letak halaman sub-views untuk mencegah layout jumping dan margin yang tidak konsisten.

#### Props `PageContainer`
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `maxWidth` | `'4xl' \| '5xl' \| '6xl' \| 'full'` | `'6xl'` | Batas lebar maksimum kontainer |
| `paddingY` | `string` | `'py-8 sm:py-12 pb-24'` | Spasi vertikal |

#### Props `PageHeader`
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `title` | `string` | *(wajib)* | Judul halaman utama |
| `description` | `string` | `''` | Subtitle / deskripsi ringkas |
| *Slots* | `breadcrumb`, `actions` | | Tempat navigasi breadcrumb dan tombol aksi |

#### Contoh Penggunaan
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

Render visual artwork cover proyek dengan rasio 4:3 (`aspect-[4/3]`) atau 16:9, lazy loading, dan **deterministic palette generator**. Jika proyek belum memiliki gambar cover, komponen otomatis menghasilkan palet warna editorial yang konsisten berdasarkan nama proyek (tanpa placeholder generik yang berantakan).

#### Props
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `src` | `string \| null` | `undefined` | URL gambar |
| `name` | `string` | `'Project'` | Nama proyek untuk generator inisial & palet |
| `aspectRatio` | `string` | `'aspect-[4/3]'` | Rasio aspek gambar |

---

### `Pagination`

Kontrol navigasi halaman tabular dengan perhitungan window otomatis dan tombol ber-ikon Phosphor.

#### Props
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `currentPage` | `number` | `1` | Halaman aktif |
| `totalItems` | `number` | `0` | Total record |
| `pageSize` | `number` | `12` | Jumlah per halaman |
| `onPageChange`| `(page: number) => void`| *(wajib)* | Callback ketika halaman berpindah |

---

### `EmptyState`

Tampilan standar saat tabel, list, atau feed tidak memiliki data.

#### Props & Slots
| Prop / Slot | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `title` | `string` | `'Nothing here yet'` | Judul empty state |
| `description` | `string` | `''` | Penjelasan mengapa kosong |
| `slot="icon"` | Slot | Default dot | Icon kustom di atas judul |
| `slot="action"`| Slot | | Tombol call-to-action untuk membuat data baru |

---

### `Timeline` & `TimelineItem`

Primitif timeline vertikal dot terstandarisasi untuk histori proyek, aktivitas workspace, audit log, dan feed publik. Menggantikan pola tabel atau card-in-card berulang dengan garis konektor tipis dan dot semantik.

#### Props `Timeline.svelte`
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `density` | `'comfortable' \| 'compact'` | `'comfortable'` | Kerapatan jarak vertikal antar elemen |

#### Props `TimelineItem.svelte`
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `title` | `string` | `''` | Judul event aktivitas (atau via slot `title`) |
| `timestamp` | `string` | `''` | Waktu relatif / format tanggal (atau via slot `timestamp`) |
| `description` | `string` | `''` | Ringkasan teks detail event |
| `dotColor` | `string` | `''` | Class Tailwind warna dot (contoh: `bg-emerald-500`) |
| `status` | `string` | `''` | Resolusi warna dot otomatis via status (`ONLINE`, `OFFLINE`) |
| `density` | `'comfortable' \| 'compact'` | `'comfortable'` | Jarak padding vertikal |
| `href` | `string` | `''` | Link opsional ke halaman target (contoh: `/projects/slug`) |
| `linkText` | `string` | `''` | Label teks link yang dapat diklik |

---

### `Tabs`

Primitif tab minimalis berbasis tipografi dan garis bawah halus (bukan pill atau kotak tebal). Digunakan di Owner Project Workspace (`Overview | Visits | Activity | Settings`) dan Public Project Detail (`Overview | Activity`).

#### Props
| Prop | Tipe | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `tabs` | `Array<{ id: string; label: string; icon?: any; badge?: number }>` | `[]` | Daftar tab navigasi |
| `active` | `string` | *(tab pertama)* | Tab aktif (bindable `bind:active={...}`) |
| `size` | `'sm' \| 'default'` | `'default'` | Ukuran teks dan padding |

---


## 3. Status & Telemetry Components (`$lib/components/status/`)

Didedikasikan untuk halaman status dan visualisasi ketersediaan sistem:

- **`StatusDot.svelte`**:
  - Indikator status tenang (`ONLINE`, `SETUP`, `PENDING`, `OFFLINE`, `ARCHIVED`).
  - Mendukung `variant="inline"` (default: dot + text tanpa pill wrapper), `variant="dot"` (dot lingkaran saja), dan `variant="badge"` (semantic pill).
  - Mendukung `showLabel={true | false}`.
- **`AvailabilityGrid.svelte`**:
  - Heatmap 90-hari ketersediaan Linux host.
  - Setiap balok merepresentasikan bucket waktu matematis (`operational`, `partial`, `incident`, `no_data`).
  - Menampilkan tooltip detail persentase dan tanggal saat di-hover.
- **`StatusOverviewCard.svelte`**:
  - Banner status global dengan persentase availability jendela berjalan dan uptime terformat.
- **`ServiceHealthCard.svelte` & `ServiceStatusRow.svelte`**:
  - Menampilkan daftar layanan internal (Web Ingress, PostgreSQL, Telemetry Probes) dan latensi edge.
- **`IncidentHistoryCard.svelte`**:
  - Rekapitulasi histori insiden downtime nyata ($\ge 120$ detik) dengan penyebab (`SYSTEM_REBOOT` atau `SERVICE_SUSPENDED`).

---

## 4. Bento & Hardware Components (`$lib/components/bento/`)

Didedikasikan untuk showcase spesifikasi hardware fisik pada landing page (`+page.svelte`):

- **`HardwareCard.svelte`**: Menampilkan prosesor Intel/AMD, jumlah physical cores, logical threads, dan arsitektur x86_64.
- **`MemoryCard.svelte`**: Menampilkan RAM fisik (`MemTotal` vs `MemAvailable`) dan bar utilisasi.
- **`StorageCard.svelte`**: Menampilkan model drive NVMe fisik (`/sys/block/nvme*`) dan kapasitas partisi Linux `/`.
- **`NetworkCard.svelte`**: Menampilkan link bandwidth dan latensi edge ke Cloudflare 1.1.1.1.
- **`ServerPhotoCard.svelte`**: Kartu visual foto tactile server homelab fisik di Jakarta.
- **`CpuPackageVisual.svelte`**: Diagram visual die package silikon prosesor.

---

## 5. Layout & Shell Components (`$lib/components/layout/`)

- **`Header.svelte`**:
  - Fixed header dengan tinggi terstandarisasi `60px`.
  - Navigasi link: `/projects`, `/people`, `/activity`, `/status`.
  - Status sesi user (Login button atau User menu).
  - Mobile drawer responsif.
- **`Footer.svelte`**:
  - Footer editorial dengan identitas homelab node, lokasi server, dan tautan dokumentasi.
- **`ThemeToggle.svelte`**:
  - Tombol toggle light/dark mode tanpa border berlebihan, menggunakan ikon Phosphor `Sun` dan `Moon`.

---

## 6. Showcase & Public Feed Components (`$lib/components/`)

- **`ProjectCard.svelte`**:
  - Kartu katalog publik untuk `/projects`.
  - Menggabungkan `ProjectCover`, judul, deskripsi 2 baris, pemilik, stack teknologi (dipisahkan dengan tanda `·`), dan link visit.
  - Bebas dari pill chip spam.
- **`Comments.svelte`**:
  - Thread diskusi kronologis pada showcase proyek publik (`/projects/[slug]`).
  - Composer komentar dengan rate-limiting 5/10 menit, sanitasi teks, author avatar & display name.
  - Penanganan soft-delete (*"This comment was removed"*).
  - Menu aksi kontekstual: Report modal trigger dan Delete dengan konfirmasi `ConfirmModal` (khusus author, project owner, atau admin).
- **`ReportModal.svelte`**:
  - Modal dialog pelaporan konten pelanggaran/spam untuk target Project dan Comment.
  - Dropdown alasan terstandarisasi, rincian opsional, backdrop blur tenang.
- **`ProjectAvailability.svelte`**:
  - Widget telemetri kesehatan proyek berbasis probe HTTP 5-menit.
  - Menampilkan uptime 30 hari dalam %, latensi respons rata-rata, dan status operasional real-time.
- **`ActivityTimeline.svelte`**:
  - Feed kronologis aktivitas publik (publikasi proyek, update, verifikasi).

---

## 7. Admin Console Components (`$lib/components/admin/`)

- **`AdminPanel.svelte`**: Sub-navigasi konsol operator (`Overview`, `Projects`, `Requests`, `Users`, `Audit`, `System`, `Settings`).
- **`AdminStatCard.svelte`**: Kartu ringkasan metrik statistik operasional.

---

## 8. Panduan Keputusan AI (Decision Matrix)

Ketika AI Coding Agent membuat atau memodifikasi tampilan, ikuti matriks keputusan ini:

| Skenario | Komponen yang WAJIB Digunakan | Pola yang DILARANG |
| :--- | :--- | :--- |
| **Menampilkan Data Tabular** | `<Table>`, `<TableRow>`, `<TableCell>` dari `$lib/components/ui` | ❌ Menulis tag `<table>`, `<th>`, `<td>` mentah dengan style ad-hoc |
| **Membuat Tombol** | `<Button variant="..." size="...">` | ❌ `<button class="rounded-full bg-gradient-to-r ...">` |
| **Status Operasional Proyek/Node** | `<StatusDot status={...} />` | ❌ Membuat chip badge buatan sendiri dengan warna acak |
| **Formulir Input** | `<Input label="..." bind:value={...} />` | ❌ `<input class="border p-2 ...">` tanpa token desain |
| **Kontainer Konten / Bento** | `<Card variant="surface" padding="md">` | ❌ Menulis `div` dengan `rounded-2xl` atau background arbitrary |
| **Menampilkan Gambar Proyek** | `<ProjectCover src={...} name={...} />` | ❌ Tag `<img>` telanjang tanpa error handling dan placeholder |
| **Ikon Antarmuka** | Impor dari `phosphor-svelte` | ❌ Raw `<svg>` atau library lain (Lucide, Heroicons) |
