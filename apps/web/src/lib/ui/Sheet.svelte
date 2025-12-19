<script lang="ts">
  import type { Snippet } from 'svelte';
  import { createDialog, melt } from '@melt-ui/svelte';
  import { X } from 'lucide-svelte';

  interface Props {
    open?: boolean;
    title?: string;
    side?: 'bottom' | 'right';
    children: Snippet;
  }

  let {
    open = $bindable(false),
    title = '',
    side = 'bottom',
    children
  }: Props = $props();

  const {
    elements: { overlay, content, title: titleEl, close, portalled },
    states: { open: sheetOpen }
  } = createDialog({ forceVisible: true });

  $effect(() => {
    if (open !== $sheetOpen) $sheetOpen = open;
  });

  $effect(() => {
    if ($sheetOpen !== open) open = $sheetOpen;
  });
</script>

<div use:melt={$portalled}>
  {#if $sheetOpen}
    <div use:melt={$overlay} class="overlay"></div>
    <div use:melt={$content} class="sheet sheet-{side}">
      {#if side === 'bottom'}
        <div class="handle-bar"></div>
      {/if}

      <div class="sheet-header">
        {#if title}
          <h2 use:melt={$titleEl} class="sheet-title">{title}</h2>
        {/if}
        <button use:melt={$close} class="close-btn" aria-label="Close">
          <X size={20} />
        </button>
      </div>

      <div class="sheet-body">
        {@render children()}
      </div>
    </div>
  {/if}
</div>

<style>
  .overlay {
    position: fixed;
    inset: 0;
    z-index: 50;
    background: var(--lb-overlay);
    animation: fadeIn 0.2s ease-out;
  }

  .sheet {
    position: fixed;
    z-index: 51;
    background: var(--lb-surface);
    overflow: hidden;
  }

  .sheet-bottom {
    left: 0;
    right: 0;
    bottom: 0;
    max-height: 90vh;
    border-radius: 16px 16px 0 0;
    animation: slideUp 0.3s cubic-bezier(0.32, 0.72, 0, 1);
  }

  .sheet-right {
    top: 0;
    right: 0;
    bottom: 0;
    width: 100%;
    max-width: 400px;
    animation: slideRight 0.3s cubic-bezier(0.32, 0.72, 0, 1);
  }

  .handle-bar {
    width: 36px;
    height: 5px;
    background: var(--lb-border);
    border-radius: 3px;
    margin: 8px auto;
  }

  .sheet-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border-bottom: 1px solid var(--lb-border);
  }

  .sheet-title {
    font-size: 17px;
    font-weight: 600;
    color: var(--lb-text);
    margin: 0;
  }

  .close-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--lb-bg);
    border: none;
    border-radius: 50%;
    color: var(--lb-muted);
    cursor: pointer;
    transition: all 0.15s;
  }

  .close-btn:hover {
    background: var(--lb-border);
    color: var(--lb-text);
  }

  .sheet-body {
    padding: 16px;
    overflow-y: auto;
    max-height: calc(90vh - 80px);
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @keyframes slideUp {
    from { transform: translateY(100%); }
    to { transform: translateY(0); }
  }

  @keyframes slideRight {
    from { transform: translateX(100%); }
    to { transform: translateX(0); }
  }
</style>
