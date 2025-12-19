<script lang="ts">
  import type { Snippet } from 'svelte';
  import { Loader2 } from 'lucide-svelte';

  interface Props {
    variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
    size?: 'sm' | 'md' | 'lg';
    disabled?: boolean;
    loading?: boolean;
    class?: string;
    onclick?: (e: MouseEvent) => void;
    children: Snippet;
  }

  let {
    variant = 'primary',
    size = 'md',
    disabled = false,
    loading = false,
    class: className = '',
    onclick,
    children
  }: Props = $props();
</script>

<button
  class="btn btn-{variant} btn-{size} {className}"
  disabled={disabled || loading}
  {onclick}
>
  {#if loading}
    <Loader2 size={size === 'sm' ? 14 : 16} class="animate-spin" />
  {/if}
  {@render children()}
</button>

<style>
  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    font-weight: 500;
    border: none;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.25, 0.1, 0.25, 1);
    user-select: none;
  }

  .btn:active:not(:disabled) {
    transform: scale(0.97);
  }

  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* Variants */
  .btn-primary {
    background: var(--lb-primary);
    color: white;
  }
  .btn-primary:hover:not(:disabled) {
    background: var(--lb-primary-hover);
  }

  .btn-secondary {
    background: var(--lb-surface);
    color: var(--lb-text);
    box-shadow: 0 0 0 1px var(--lb-border);
  }
  .btn-secondary:hover:not(:disabled) {
    background: var(--lb-bg);
  }

  .btn-ghost {
    background: transparent;
    color: var(--lb-primary);
  }
  .btn-ghost:hover:not(:disabled) {
    background: rgba(0, 122, 255, 0.1);
  }

  .btn-danger {
    background: var(--lb-error);
    color: white;
  }
  .btn-danger:hover:not(:disabled) {
    opacity: 0.9;
  }

  /* Sizes */
  .btn-sm {
    height: 32px;
    padding: 0 12px;
    font-size: 13px;
    border-radius: 8px;
  }

  .btn-md {
    height: 40px;
    padding: 0 16px;
    font-size: 14px;
    border-radius: 10px;
  }

  .btn-lg {
    height: 48px;
    padding: 0 24px;
    font-size: 16px;
    border-radius: 12px;
  }

  .animate-spin {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
</style>
