<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'default' | 'success' | 'warning' | 'error' | 'primary';
    size?: 'sm' | 'md';
    children: Snippet;
  }

  let {
    variant = 'default',
    size = 'sm',
    children
  }: Props = $props();
</script>

<span class="badge badge-{variant} badge-{size}">
  {@render children()}
</span>

<style>
  .badge {
    display: inline-flex;
    align-items: center;
    font-weight: 500;
    border-radius: 6px;
  }

  .badge-sm {
    padding: 2px 8px;
    font-size: 11px;
  }

  .badge-md {
    padding: 4px 10px;
    font-size: 12px;
  }

  .badge-default {
    background: var(--lb-bg);
    color: var(--lb-muted);
  }

  .badge-success {
    background: rgba(52, 199, 89, 0.15);
    color: var(--lb-success);
  }

  .badge-warning {
    background: rgba(255, 149, 0, 0.15);
    color: var(--lb-warning);
  }

  .badge-error {
    background: rgba(255, 59, 48, 0.15);
    color: var(--lb-error);
  }

  .badge-primary {
    background: rgba(0, 122, 255, 0.15);
    color: var(--lb-primary);
  }
</style>
