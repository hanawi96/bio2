<script lang="ts">
  let password = $state('');
  let error = $state('');
  let loading = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    loading = true;
    error = '';

    try {
      const response = await fetch('/r/password', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password })
      });

      if (response.ok) {
        window.location.reload();
      } else {
        const data = await response.json();
        error = data.error?.message || 'Invalid password';
      }
    } catch {
      error = 'Failed to verify password';
    } finally {
      loading = false;
    }
  }
</script>

<div class="container">
  <div class="card">
    <h1>🔒 Password Protected</h1>
    <p>This page requires a password to view.</p>

    <form onsubmit={handleSubmit}>
      <input
        type="password"
        bind:value={password}
        placeholder="Enter password"
        disabled={loading}
      />
      <button type="submit" disabled={loading || !password}>
        {loading ? 'Verifying...' : 'Enter'}
      </button>
    </form>

    {#if error}
      <p class="error">{error}</p>
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
    max-width: 400px;
  }

  h1 {
    font-size: 1.5rem;
    margin-bottom: 8px;
  }

  p {
    opacity: 0.7;
    margin-bottom: 24px;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  input {
    padding: 12px 16px;
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 8px;
    background: rgba(255, 255, 255, 0.1);
    color: white;
    font-size: 16px;
  }

  input::placeholder {
    color: rgba(255, 255, 255, 0.5);
  }

  button {
    padding: 12px 24px;
    border: none;
    border-radius: 8px;
    background: #3b82f6;
    color: white;
    font-weight: 500;
    cursor: pointer;
    transition: opacity 0.15s;
  }

  button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .error {
    color: #ef4444;
    margin-top: 12px;
  }
</style>
