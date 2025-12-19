<script lang="ts">
  import { ChevronDown } from 'lucide-svelte';

  interface Option {
    value: string;
    label: string;
  }

  interface Props {
    value?: string;
    options: Option[];
    placeholder?: string;
    label?: string;
    disabled?: boolean;
    onchange?: (value: string) => void;
  }

  let {
    value = $bindable(''),
    options,
    placeholder = 'Select...',
    label = '',
    disabled = false,
    onchange
  }: Props = $props();

  function handleChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    value = target.value;
    onchange?.(value);
  }
</script>

<div class="select-wrapper">
  {#if label}
    <label class="select-label">{label}</label>
  {/if}
  <div class="select-container">
    <select class="select" {value} {disabled} onchange={handleChange}>
      {#if placeholder}
        <option value="" disabled>{placeholder}</option>
      {/if}
      {#each options as opt}
        <option value={opt.value}>{opt.label}</option>
      {/each}
    </select>
    <ChevronDown size={18} class="select-icon" />
  </div>
</div>

<style>
  .select-wrapper {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .select-label {
    font-size: 13px;
    font-weight: 500;
    color: var(--lb-text);
  }

  .select-container {
    position: relative;
  }

  .select {
    width: 100%;
    height: 44px;
    padding: 0 36px 0 14px;
    font-size: 16px;
    color: var(--lb-text);
    background: var(--lb-surface);
    border: 1px solid var(--lb-border);
    border-radius: 10px;
    appearance: none;
    cursor: pointer;
    transition: all 0.2s;
  }

  .select:focus {
    outline: none;
    border-color: var(--lb-primary);
    box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.15);
  }

  .select:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  :global(.select-icon) {
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--lb-muted);
    pointer-events: none;
  }
</style>
