<script lang="ts">
  interface Props {
    type?: 'text' | 'email' | 'password' | 'url' | 'number';
    value?: string;
    placeholder?: string;
    label?: string;
    error?: string;
    disabled?: boolean;
    class?: string;
    oninput?: (e: Event) => void;
    onkeydown?: (e: KeyboardEvent) => void;
  }

  let {
    type = 'text',
    value = $bindable(''),
    placeholder = '',
    label = '',
    error = '',
    disabled = false,
    class: className = '',
    oninput,
    onkeydown
  }: Props = $props();

  let inputId = `input-${Math.random().toString(36).slice(2, 9)}`;
</script>

<div class="input-wrapper {className}">
  {#if label}
    <label for={inputId} class="input-label">{label}</label>
  {/if}
  <input
    id={inputId}
    {type}
    bind:value
    {placeholder}
    {disabled}
    class="input"
    class:error
    {oninput}
    {onkeydown}
  />
  {#if error}
    <span class="input-error">{error}</span>
  {/if}
</div>

<style>
  .input-wrapper {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .input-label {
    font-size: 13px;
    font-weight: 500;
    color: var(--lb-text);
  }

  .input {
    height: 44px;
    padding: 0 14px;
    font-size: 16px;
    color: var(--lb-text);
    background: var(--lb-surface);
    border: 1px solid var(--lb-border);
    border-radius: 10px;
    outline: none;
    transition: all 0.2s;
  }

  .input::placeholder {
    color: var(--lb-muted);
  }

  .input:focus {
    border-color: var(--lb-primary);
    box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.15);
  }

  .input:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .input.error {
    border-color: var(--lb-error);
  }

  .input.error:focus {
    box-shadow: 0 0 0 3px rgba(255, 59, 48, 0.15);
  }

  .input-error {
    font-size: 12px;
    color: var(--lb-error);
  }
</style>
