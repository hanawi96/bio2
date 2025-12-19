<script lang="ts">
  interface Props {
    size?: number;
    color?: string;
  }

  let { size = 20, color = 'currentColor' }: Props = $props();
</script>

<svg
  width={size}
  height={size}
  viewBox="0 0 24 24"
  fill="none"
  class="spinner"
  style="color: {color}"
>
  <circle
    cx="12"
    cy="12"
    r="10"
    stroke="currentColor"
    stroke-width="3"
    stroke-linecap="round"
    opacity="0.25"
  />
  <path
    d="M12 2a10 10 0 0 1 10 10"
    stroke="currentColor"
    stroke-width="3"
    stroke-linecap="round"
  />
</svg>

<style>
  .spinner {
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
</style>
