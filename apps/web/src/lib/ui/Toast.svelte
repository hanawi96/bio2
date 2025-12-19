<script lang="ts">
  import { Check, X, AlertCircle, Info } from 'lucide-svelte';

  interface ToastData {
    id: string;
    message: string;
    type: 'success' | 'error' | 'info';
  }

  let toasts = $state<ToastData[]>([]);

  export function showToast(message: string, type: 'success' | 'error' | 'info' = 'info') {
    const id = Math.random().toString(36).slice(2);
    toasts = [...toasts, { id, message, type }];

    setTimeout(() => {
      toasts = toasts.filter(t => t.id !== id);
    }, 3000);
  }

  function dismiss(id: string) {
    toasts = toasts.filter(t => t.id !== id);
  }

  const icons = { success: Check, error: AlertCircle, info: Info };
</script>

{#if toasts.length > 0}
  <div class="toast-container">
    {#each toasts as toast (toast.id)}
      <div class="toast toast-{toast.type}">
        <svelte:component this={icons[toast.type]} size={18} />
        <span>{toast.message}</span>
        <button class="dismiss" onclick={() => dismiss(toast.id)} aria-label="Dismiss">
          <X size={14} />
        </button>
      </div>
    {/each}
  </div>
{/if}

<style>
  .toast-container {
    position: fixed;
    bottom: 24px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 100;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .toast {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 16px;
    border-radius: 12px;
    font-size: 14px;
    font-weight: 500;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
    animation: slideUp 0.3s cubic-bezier(0.32, 0.72, 0, 1);
  }

  .toast-success {
    background: var(--lb-success);
    color: white;
  }

  .toast-error {
    background: var(--lb-error);
    color: white;
  }

  .toast-info {
    background: var(--lb-surface);
    color: var(--lb-text);
    border: 1px solid var(--lb-border);
  }

  .dismiss {
    margin-left: 4px;
    padding: 4px;
    background: none;
    border: none;
    color: inherit;
    opacity: 0.7;
    cursor: pointer;
    border-radius: 4px;
  }

  .dismiss:hover {
    opacity: 1;
    background: rgba(255, 255, 255, 0.2);
  }

  @keyframes slideUp {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }
</style>
