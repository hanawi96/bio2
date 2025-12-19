<script lang="ts">
  interface Props {
    checked?: boolean;
    disabled?: boolean;
    label?: string;
    onchange?: (checked: boolean) => void;
  }

  let {
    checked = $bindable(false),
    disabled = false,
    label = '',
    onchange
  }: Props = $props();

  function toggle() {
    if (disabled) return;
    checked = !checked;
    onchange?.(checked);
  }
</script>

<label class="toggle-wrapper" class:disabled>
  <button
    type="button"
    role="switch"
    aria-checked={checked}
    class="toggle"
    class:checked
    onclick={toggle}
    {disabled}
  >
    <span class="thumb"></span>
  </button>
  {#if label}
    <span class="toggle-label">{label}</span>
  {/if}
</label>

<style>
  .toggle-wrapper {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
  }

  .toggle-wrapper.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .toggle {
    position: relative;
    width: 51px;
    height: 31px;
    background: var(--lb-border);
    border: none;
    border-radius: 16px;
    cursor: pointer;
    transition: background 0.2s;
  }

  .toggle.checked {
    background: var(--lb-success);
  }

  .thumb {
    position: absolute;
    top: 2px;
    left: 2px;
    width: 27px;
    height: 27px;
    background: white;
    border-radius: 50%;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    transition: transform 0.2s cubic-bezier(0.25, 0.1, 0.25, 1);
  }

  .toggle.checked .thumb {
    transform: translateX(20px);
  }

  .toggle-label {
    font-size: 15px;
    color: var(--lb-text);
  }
</style>
