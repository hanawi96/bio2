<script lang="ts">
  interface ErrorInfo {
    code: string;
    message: string;
  }

  let { error }: { error: ErrorInfo } = $props();
</script>

<div class="container">
  <div class="card">
    {#if error.code === 'NOT_FOUND'}
      <h1>404</h1>
      <p>Page not found</p>
    {:else}
      <h1>Error</h1>
      <p>{error.message}</p>
    {/if}
  </div>
</div>

<style>
  .container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
    background: #0B0F19;
    color: white;
  }

  .card {
    text-align: center;
  }

  h1 {
    font-size: 4rem;
    margin-bottom: 8px;
    opacity: 0.3;
  }

  p {
    opacity: 0.7;
  }
</style>
