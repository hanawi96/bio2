# LinkBio Web (SvelteKit + Svelte 5)

## Structure

```
apps/web/
├─ src/
│  ├─ routes/
│  │  └─ (public)/          # Public renderer routes
│  │     └─ [[...path]]/    # Catch-all for custom domains
│  ├─ lib/
│  │  └─ renderer/          # Shared renderer components
│  │     ├─ PageRenderer.svelte
│  │     ├─ BlockRenderer.svelte
│  │     ├─ Header.svelte
│  │     └─ blocks/
│  │        ├─ TextBlock.svelte
│  │        └─ LinkGroupBlock.svelte
│  └─ types/                # TypeScript types
└─ static/
```

## Run

```bash
npm install
npm run dev
```

## Routes (Slice 1)

- `(public)/[[...path]]` - Public page renderer
  - Fetches compiled_json from API `/r`
  - Renders text and link_group blocks
