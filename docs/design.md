# Ngumpul Host — Design System

## Design Direction

Ngumpul Host uses:

**Editorial × Bento × Quiet Technical**

The interface should feel like a carefully designed small corner of the internet built by a real person for a small community.

It should feel:
- calm, personal, technical, editorial, spacious
- understated, intentional, slightly playful, trustworthy

It must NOT feel like:
- a generic SaaS landing page
- an enterprise dashboard
- a commercial hosting company
- an AI startup or Web3 website
- a cybersecurity platform or generic template

---

## 1. Anti-AI-Slop Rules

These rules are mandatory. Do not use visual patterns simply because they are common in AI-generated interfaces:

- **Banned:** Purple-to-blue gradients, neon cyan glows, large background blur blobs.
- **Banned:** Generic floating glass cards on everything, excessive rounded-3xl corners.
- **Banned:** Three identical cards in a row with identical icons.
- **Banned:** Centered hero with generic badge + gradient headline + two equal CTA buttons.
- **Banned:** Fake analytics, fake testimonials, fake uptime counters.
- **Banned:** Corporate buzzwords ("Empower your workflow", "Unlock your potential", "Next-gen platform").

Copy must describe the actual product:
- *Good:* "A small corner of the internet for friends, projects, and experiments."
- *Bad:* "Empowering digital transformation through modern web infrastructure."

---

## 2. Color Palette

The primary identity is a restrained muted sky-blue (`#79AFC4`). It resembles morning sky, cool paper, and quiet infrastructure.

### Light Theme (Paper / Sky / Quiet)
```text
Background       #F4F7F8
Surface          #FFFFFF
Surface Muted    #EAF0F3

Ink              #182126
Ink Secondary    #58656C
Ink Muted        #879299

Primary          #79AFC4
Primary Strong   #477F97
Primary Pale     #DCECF2

Border           #D9E1E5

Success          #668B73
Warning          #A58A58
Danger           #A96868
```

### Dark Theme (Night / Infrastructure / Quiet)
```text
Background       #101618
Surface          #161E21
Surface Muted    #1C272B

Ink              #E7EEF0
Ink Secondary    #A7B5BA
Ink Muted        #718187

Primary          #8FC5D8
Primary Strong   #639AAF
Primary Pale     #20353D

Border           #2B383D

Success          #88A994
Warning          #BBA170
Danger           #C08383
```

Theme switching behaves as **Day ➔ Night in the same material**, not shifting from an editorial paper site into a glowing cyberpunk terminal.

---

## 3. Typography

- **Display & Headings:** `Manrope`, sans-serif (Google Fonts)
  - Hero heading, section titles, project titles, numerical parameters.
- **Body & Interface:** `DM Sans`, sans-serif (Google Fonts)
  - Paragraphs, labels, metadata, buttons, controls.
- **Monospace:** `font-mono` (reserved strictly for URLs, hashes, CLI commands, and technical configs).

Typography Scale:
- Hero: `clamp(2.75rem, 5vw, 4.5rem)` (short, max 2 lines)
- Section headings: `clamp(1.75rem, 3vw, 2.5rem)`
- Project titles: `1.5rem – 2rem`
- Large parameters: `clamp(2.5rem, 4vw, 3.5rem)`
- Body text: `1rem` / `leading-relaxed` (max 65ch)

---

## 4. Layout & Bento System

- Max container width: `1440px` with responsive horizontal padding.
- **Asymmetric Bento Grids:** Bento layouts combine diverse information densities (e.g. large featured project visual + compact numeric parameters + live operational dot) rather than 6 identical cards.
- **Varied Section Rhythms:** Mix editorial text, asymmetric split layouts, timeline lists, and project showcases.
- **Quiet Status Indicators:** Subtle dots (`● Online`, `● Setup`, `● Offline`) without frantic neon pulses.
