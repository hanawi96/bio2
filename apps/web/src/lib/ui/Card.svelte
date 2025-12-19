<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'default' | 'elevated' | 'outlined';
    padding?: 'none' | 'sm' | 'md' | 'lg';
    class?: string;
    onclick?: () => void;
    children: Snippet;
  }

  let {
    variant = 'default',
    padding = 'md',
    class: className = '',
    onclick,
    children
  }: Props = $props();
</script>

<div
  class="card card-{variant} pad-{padding} {className}"
  class:clickable={!!onclick}
  onclick={onclick}
  onkeydown={(e) => e.key === 'Enter' && onclick?.()}
  role={onclick ? 'button' : undefined}
  tabindex={onclick ? 0 : undefined}
>
  {@render children()}
</div>

<style>
  .card {
    background: var(--lb-surface);
    border-radius: 12px;
    transition: all 0.2s;
  }

  .card-default {
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  }

  .card-elevated {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  }

  .card-outlined {
    border: 1px solid var(--lb-border);
  }

  .clickable {
    cursor: pointer;
  }

  .clickable:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
  }

  .clickable:active {
    transform: translateY(0);
  }

  .pad-none { padding: 0; }
  .pad-sm { padding: 12px; }
  .pad-md { padding: 16px; }
  .pad-lg { padding: 24px; }
</style>
