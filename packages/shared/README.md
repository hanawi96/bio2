# @linkbio/shared

Shared types, contracts, and schemas for Link-in-Bio Builder.

## Structure

```
packages/shared/
├─ contracts/
│  ├─ api.ts      # Request/response types
│  ├─ page.ts     # Draft + compiled payload types
│  ├─ blocks.ts   # Block type enums
│  └─ theme.ts    # Theme contract types
└─ schemas/       # Zod validation schemas (future)
```

## Usage

```typescript
import type { CompiledPayload, CompiledBlock } from '@linkbio/shared/contracts/page';
```
