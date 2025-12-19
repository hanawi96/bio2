<script lang="ts">
  import type { Snippet } from 'svelte';
  import { createDialog, melt } from '@melt-ui/svelte';
  import { X } from 'lucide-svelte';

  interface Props {
    open?: boolean;
    title?: string;
    description?: string;
    children: Snippet;
    footer?: Snippet;
  }

  let {
    open = $bindable(false),
    title = '',
    description = '',
    children,
    footer
  }: Props = $props();

  const {
    elements: { overlay, content, title: titleEl, description: descEl, close, portalled },
    states: { open: dialogOpen }
  } = createDialog({ forceVisible: true });

  $effect(() => {
    if (open !== $dialogOpen) $dialogOpen = open;
  });

  $effect(() => {
    if ($dialogOpen !== open) open = $dialogOpen;
  });
</script>

<div use:melt={$portalled}>
  {#if $dialogOpen}
    <div use:melt={$overlay} class="overlay"></div>
    <div use:melt={$content} class="dialog">
      <div class="dialog-header">
        {#if title}
          <h2 use:melt={$titleEl} class="dialog-title">{title}</h2>
        {/if}
        {#if description}
          <p use:melt={$descEl} class="dialog-desc">{description}</p>
        {/if}
        <button use:melt={$close} class="close-btn" aria-label="Close">
          <X size={20} />
        </button>
      </div>

      <div class="dialog-body">
        {@render children()}
      </div>

      {#if footer}
        <div class="dialog-footer">
          {@render footer()}
        </div>
      {/if}
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

  .dialog {
    position: fixed;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    z-index: 51;
    width: 90vw;
    max-width: 420px;
    max-height: 85vh;
    background: var(--lb-surface);
    border-radius: 16px;
    box-shadow: 0 24px 48px rgba(0, 0, 0, 0.2);
    overflow: hidden;
    animation: scaleIn 0.2s ease-out;
  }

  .dialog-header {
    position: relative;
    padding: 20px 20px 0;
  }

  .dialog-title {
    font-size: 18px;
    font-weight: 600;
    color: var(--lb-text);
    margin: 0 0 4px;
    padding-right: 32px;
  }

  .dialog-desc {
    font-size: 14px;
    color: var(--lb-muted);
    margin: 0;
  }

  .close-btn {
    position: absolute;
    top: 16px;
    right: 16px;
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

  .dialog-body {
    padding: 20px;
    overflow-y: auto;
  }

  .dialog-footer {
    padding: 16px 20px;
    border-top: 1px solid var(--lb-border);
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @keyframes scaleIn {
    from { opacity: 0; transform: translate(-50%, -50%) scale(0.95); }
    to { opacity: 1; transform: translate(-50%, -50%) scale(1); }
  }
</style>
