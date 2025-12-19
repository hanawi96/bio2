<script lang="ts">
  import type { PageData } from './$types';
  import { Button, Input, Card, Badge } from '$lib/ui';
  import { Plus, FileText, Globe, ChevronRight } from 'lucide-svelte';

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

<div class="page-container">
  <header class="page-header">
    <div>
      <h1>My Pages</h1>
      <p>Create and manage your link-in-bio pages</p>
    </div>
  </header>

  <Card variant="outlined" padding="lg" class="create-card">
    <div class="create-form">
      <Input
        bind:value={title}
        placeholder="Enter page title..."
        disabled={creating}
        onkeydown={(e) => e.key === 'Enter' && createPage()}
      />
      <Button
        variant="primary"
        size="lg"
        onclick={createPage}
        disabled={creating || !title.trim()}
        loading={creating}
      >
        <Plus size={18} />
        Create Page
      </Button>
    </div>
  </Card>

  <section class="pages-section">
    <h2>Your Pages</h2>

    {#if data.pages && data.pages.length > 0}
      <div class="pages-grid">
        {#each data.pages as page}
          <a href="/app/pages/{page.id}/edit" class="page-card">
            <div class="page-icon">
              <FileText size={24} />
            </div>
            <div class="page-info">
              <h3>{page.title || 'Untitled'}</h3>
              <div class="page-meta">
                <Badge variant={page.status === 'published' ? 'success' : 'default'}>
                  {#if page.status === 'published'}
                    <Globe size={10} />
                  {/if}
                  {page.status}
                </Badge>
                <span class="locale">{page.locale.toUpperCase()}</span>
              </div>
            </div>
            <ChevronRight size={20} class="chevron" />
          </a>
        {/each}
      </div>
    {:else}
      <div class="empty-state">
        <div class="empty-icon">
          <FileText size={48} />
        </div>
        <h3>No pages yet</h3>
        <p>Create your first page to get started</p>
      </div>
    {/if}
  </section>
</div>

<style>
  .page-container {
    max-width: 800px;
    margin: 0 auto;
    padding: 32px 24px;
  }

  .page-header {
    margin-bottom: 32px;
  }

  .page-header h1 {
    font-size: 32px;
    font-weight: 700;
    color: var(--lb-text);
    margin: 0;
    letter-spacing: -0.5px;
  }

  .page-header p {
    font-size: 16px;
    color: var(--lb-muted);
    margin: 8px 0 0;
  }

  :global(.create-card) {
    margin-bottom: 40px;
  }

  .create-form {
    display: flex;
    gap: 12px;
  }

  .create-form :global(.input-wrapper) {
    flex: 1;
  }

  .pages-section h2 {
    font-size: 13px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--lb-muted);
    margin: 0 0 16px;
  }

  .pages-grid {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .page-card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px 20px;
    background: var(--lb-surface);
    border: 1px solid var(--lb-border);
    border-radius: 14px;
    text-decoration: none;
    color: inherit;
    transition: all 0.2s;
  }

  .page-card:hover {
    border-color: var(--lb-primary);
    box-shadow: 0 4px 12px rgba(0, 122, 255, 0.1);
  }

  .page-card:active {
    transform: scale(0.99);
  }

  .page-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 48px;
    height: 48px;
    background: var(--lb-bg);
    border-radius: 12px;
    color: var(--lb-primary);
  }

  .page-info {
    flex: 1;
    min-width: 0;
  }

  .page-info h3 {
    font-size: 16px;
    font-weight: 600;
    color: var(--lb-text);
    margin: 0 0 6px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .page-meta {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .locale {
    font-size: 11px;
    font-weight: 500;
    color: var(--lb-muted);
  }

  :global(.chevron) {
    color: var(--lb-muted);
    flex-shrink: 0;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 60px 20px;
    text-align: center;
  }

  .empty-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 80px;
    height: 80px;
    background: var(--lb-bg);
    border-radius: 20px;
    color: var(--lb-muted);
    margin-bottom: 20px;
  }

  .empty-state h3 {
    font-size: 18px;
    font-weight: 600;
    color: var(--lb-text);
    margin: 0 0 8px;
  }

  .empty-state p {
    font-size: 14px;
    color: var(--lb-muted);
    margin: 0;
  }
</style>
