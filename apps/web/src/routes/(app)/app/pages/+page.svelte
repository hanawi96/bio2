<script lang="ts">
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();

  let title = $state('');
  let creating = $state(false);

  async function createPage() {
    if (!title.trim()) return;
    creating = true;

    try {
      const res = await fetch('/api/pages', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: title.trim(),
          locale: 'vi',
          theme_preset_key: 'theme_a'
        })
      });

      if (res.ok) {
        const { page_id } = await res.json();
        window.location.href = `/app/pages/${page_id}/edit`;
      }
    } finally {
      creating = false;
    }
  }
</script>

<svelte:head>
  <title>My Pages</title>
</svelte:head>

<div class="pages-container">
  <h1>My Pages</h1>

  <div class="create-form">
    <input
      type="text"
      bind:value={title}
      placeholder="Page title..."
      disabled={creating}
    />
    <button onclick={createPage} disabled={creating || !title.trim()}>
      {creating ? 'Creating...' : 'Create Page'}
    </button>
  </div>

  <div class="pages-list">
    {#if data.pages && data.pages.length > 0}
      {#each data.pages as page}
        <a href="/app/pages/{page.id}/edit" class="page-card">
          <h3>{page.title || 'Untitled'}</h3>
          <div class="page-meta">
            <span class="status" class:published={page.status === 'published'}>
              {page.status}
            </span>
            <span class="locale">{page.locale}</span>
          </div>
        </a>
      {/each}
    {:else}
      <p class="empty">No pages yet. Create your first page!</p>
    {/if}
  </div>
</div>

<style>
  .pages-container {
    max-width: 800px;
  }

  h1 {
    margin-bottom: 24px;
  }

  .create-form {
    display: flex;
    gap: 12px;
    margin-bottom: 32px;
  }

  .create-form input {
    flex: 1;
    padding: 12px 16px;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-size: 16px;
  }

  .create-form button {
    padding: 12px 24px;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 8px;
    font-weight: 500;
    cursor: pointer;
  }

  .create-form button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .pages-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .page-card {
    display: block;
    padding: 16px 20px;
    background: white;
    border-radius: 8px;
    text-decoration: none;
    color: inherit;
    border: 1px solid #e0e0e0;
    transition: box-shadow 0.15s;
  }

  .page-card:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .page-card h3 {
    margin: 0 0 8px;
  }

  .page-meta {
    display: flex;
    gap: 12px;
    font-size: 14px;
    color: #666;
  }

  .status {
    padding: 2px 8px;
    background: #f0f0f0;
    border-radius: 4px;
  }

  .status.published {
    background: #dcfce7;
    color: #166534;
  }

  .empty {
    color: #666;
    text-align: center;
    padding: 40px;
  }
</style>
