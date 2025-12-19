<script lang="ts">
  import type { PageData } from './$types';
  import { page } from '$app/stores';
  import EditorPanel from '$lib/editor/EditorPanel.svelte';
  import PreviewPanel from '$lib/editor/PreviewPanel.svelte';

  let { data }: { data: PageData } = $props();

  // Editor state - initialize from data
  let draft = $state(structuredClone(data.draft));
  let saving = $state(false);
  let publishing = $state(false);
  let message = $state('');

  const pageId = $derived($page.params.pageId);

  async function save() {
    if (!draft) return;
    saving = true;
    message = '';

    try {
      const res = await fetch(`/api/pages/${pageId}/save`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(draft)
      });

      if (res.ok) {
        message = 'Saved!';
        setTimeout(() => (message = ''), 2000);
      } else {
        const err = await res.json();
        message = `Error: ${err.error?.message || 'Save failed'}`;
      }
    } catch {
      message = 'Error: Network error';
    } finally {
      saving = false;
    }
  }

  async function publish() {
    if (!draft) return;
    publishing = true;
    message = '';

    try {
      // Save first
      await fetch(`/api/pages/${pageId}/save`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(draft)
      });

      // Then publish
      const res = await fetch(`/api/pages/${pageId}/publish`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({})
      });

      if (res.ok) {
        message = 'Published!';
        if (draft) draft.page.status = 'published';
        setTimeout(() => (message = ''), 2000);
      } else {
        const err = await res.json();
        message = `Error: ${err.error?.message || 'Publish failed'}`;
      }
    } catch {
      message = 'Error: Network error';
    } finally {
      publishing = false;
    }
  }
</script>

<svelte:head>
  <title>Edit: {draft?.page?.title || 'Page'}</title>
</svelte:head>

{#if data.error}
  <div class="error-container">
    <h2>Error</h2>
    <p>{data.error.message}</p>
    <a href="/app/pages">Back to pages</a>
  </div>
{:else if draft}
  <div class="editor-layout">
    <div class="editor-header">
      <a href="/app/pages" class="back-link">← Back</a>
      <h2>{draft.page.title || 'Untitled'}</h2>
      <div class="actions">
        {#if message}
          <span class="message">{message}</span>
        {/if}
        <button onclick={save} disabled={saving} class="btn-save">
          {saving ? 'Saving...' : 'Save'}
        </button>
        <button onclick={publish} disabled={publishing} class="btn-publish">
          {publishing ? 'Publishing...' : 'Publish'}
        </button>
      </div>
    </div>

    <div class="editor-body">
      <div class="panel-editor">
        <EditorPanel bind:draft />
      </div>
      <div class="panel-preview">
        <PreviewPanel {draft} />
      </div>
    </div>
  </div>
{/if}

<style>
  .error-container {
    text-align: center;
    padding: 40px;
  }

  .editor-layout {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 60px);
    margin: -24px;
  }

  .editor-header {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 12px 24px;
    background: white;
    border-bottom: 1px solid #e0e0e0;
  }

  .back-link {
    color: #666;
    text-decoration: none;
  }

  .editor-header h2 {
    flex: 1;
    margin: 0;
    font-size: 18px;
  }

  .actions {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .message {
    font-size: 14px;
    color: #16a34a;
  }

  .btn-save,
  .btn-publish {
    padding: 8px 16px;
    border: none;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
  }

  .btn-save {
    background: #e0e0e0;
    color: #333;
  }

  .btn-publish {
    background: #3b82f6;
    color: white;
  }

  .btn-save:disabled,
  .btn-publish:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .editor-body {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  .panel-editor {
    width: 400px;
    background: white;
    border-right: 1px solid #e0e0e0;
    overflow-y: auto;
  }

  .panel-preview {
    flex: 1;
    background: #1a1a2e;
    overflow-y: auto;
  }
</style>
