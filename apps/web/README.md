# LinkBio Web (SvelteKit + Svelte 5)

## Structure

```
apps/web/
├─ src/
│  ├─ routes/
│  │  ├─ (public)/          # Public renderer routes (lightweight)
│  │  │  └─ [[...path]]/    # Catch-all for custom domains
│  │  └─ (app)/             # Editor/dashboard routes
│  │     └─ app/pages/[pageId]/edit/  # Page editor
│  ├─ lib/
│  │  ├─ renderer/          # Shared renderer (public + preview)
│  │  ├─ editor/            # Editor components
│  │  ├─ stores/            # Svelte 5 stores
│  │  │  └─ editorSync.svelte.ts  # Optimistic save/publish
│  │  └─ ui/                # UI primitives (editor only)
│  └─ types/                # TypeScript types
└─ static/
```

## Run

```bash
npm install
npm run dev
```

## UI Stack

- **Tailwind CSS** - Utility-first CSS
- **Melt UI** - Headless UI primitives
- **bits-ui** - Svelte UI components
- **lucide-svelte** - Icons
- **motion** - Subtle animations (Dialog/Sheet only)

## Design Tokens

CSS variables defined in `src/app.css`:

- `--lb-bg`, `--lb-surface`, `--lb-text`, `--lb-muted`, `--lb-border`
- `--lb-primary`, `--lb-primary-text`, `--lb-error`, `--lb-success`
- `--lb-radius-sm` (8px), `--lb-radius-md` (14px), `--lb-radius-lg` (18px)
- `--lb-shadow-sm`, `--lb-shadow-md`, `--lb-shadow-soft`

## Theme Modes

Toggle modes by adding classes to `<html>` or `<body>`:

```html
<!-- Dark mode -->
<body class="dark">

<!-- Compact mode -->
<body class="compact">

<!-- Both -->
<body class="dark compact">
```

## Optimistic Editor Status

The editor shows real-time status in the top bar:

- **Ready** - No changes
- **Unsaved** - Changes pending
- **Saving…** - Save in progress
- **Saved** - Successfully saved
- **Error** - Save failed (with retry option)
- **Publishing…** - Publish in progress

Status transitions:
1. Edit → Unsaved (dirty)
2. Click Save → Saving… → Saved (or Error)
3. Click Publish → Publishing… → Saved (or Error)
4. Error → Click Retry → Saving…

## Hard Constraints

**DO NOT** import `bits-ui`, `@melt-ui/svelte`, or `motion` in:
- `src/routes/(public)/**`
- `src/lib/renderer/**`

Public routes must stay lightweight for fast page loads.
