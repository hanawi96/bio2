<script lang="ts">
  import type { PageData } from './$types';
  import { page } from '$app/stores';
  import { onDestroy } from 'svelte';
  import EditorPanel from '$lib/editor/EditorPanel.svelte';
  import PreviewPanel from '$lib/editor/PreviewPanel.svelte';
  import { Button, Badge, Spinner } from '$lib/ui';
  import { ArrowLeft, Check, AlertCircle, Circle, Save, Globe } from 'lucide-svelte';
  import {
    getEditorState,
    markDirty,
    saveDraft,
    publish,
    retrySave,
    resetState
  } from '$lib/stores/editorSync.svelte';

  let { data }: { data: PageData } = $props();

  let draft = $state(data.draft ? structuredClone(data.draft) : null);
  const editorState = getEditorState();
  const pageId = $derived(Number($page.params.pageId));

  let initialDraftJson = data.draft ? JSON.stringify(data.draft) : '';

  $effect(() => {
    if (draft) {
      const currentJson = JSON.stringify(draft);
      if (currentJson !== initialDraftJson) {
        markDirty();
      }
    }
  });

  onDestroy(() => resetState());

  async function handleSave() {
    if (!draft) return;
    const success = await saveDraft(pageId, draft);
    if (success) initialDraftJson = JSON.stringify(draft);
  }

  async function handlePublish() {
    if (!draft) return;
    const success = await publish(pageId, draft);
    if (success && draft) {
      draft.page.status = 'published';
      initialDraftJson = JSON.stringify(draft);
    }
  }

  const canSave = $derived(editorState.dirty && editorState.status !== 'saving' && editorState.status !== 'publishing');
  const canPublish = $derived(editorState.status !== 'saving' && editorState.status !== 'publishing');

  const statusBadge = $derived.by(() => {
    if (editorState.status === 'error') return { variant: 'error' as const, text: 'Error', icon: AlertCircle };
    if (editorState.status === 'saving') return { variant: 'default' as const, text: 'Saving...', icon: null };
    if (editorState.status === 'publishing') return { variant: 'primary' as const, text: 'Publishing...', icon: null };
    if (editorState.status === 'saved' && !editorState.dirty) return { variant: 'success' as const, text: 'Saved', icon: Check };
    if (editorState.dirty) return { variant: 'warning' as const, text: 'Unsaved', icon: Circle };
    return { variant: 'default' as const, text: 'Ready', icon: Check };
  });
</script>

<svelte:head>
  <title>Edit: {draft?.page?.title || 'Page'}</title>
</svelte:head>

{#if data.error}
  <div class="error-page">
    <AlertCircle size={56} />
    <h2>Something went wrong</h2>
    <p>{data.error.message}</p>
    <Button variant="secondary" onclick={() => window.location.href = '/app/pages'}>
      <ArrowLeft size={16} />
      Back to pages
    </Button>
  </div>
{:else if draft}
  <div class="editor-layout">
    <!-- Header -->
    <header class="editor-header glass">
      <div class="header-left">
        <a href="/app/pages" class="back-link">
          <ArrowLeft size={20} />
        </a>
        <div class="title-section">
          <h1>{draft.page.title || 'Untitled'}</h1>
          <div class="status-row">
            <Badge variant={statusBadge.variant}>
              {#if statusBadge.icon}
                <svelte:component this={statusBadge.icon} size={12} />
              {:else}
                <Spinner size={12} />
              {/if}
              {statusBadge.text}
            </Badge>
            {#if draft.page.status === 'published'}
              <Badge variant="success">
                <Globe size={12} />
                Published
              </Badge>
            {/if}
          </div>
        </div>
      </div>

      <div class="header-actions">
        {#if editorState.status === 'error'}
          <Button variant="ghost" size="sm" onclick={retrySave}>
            Retry
          </Button>
        {/if}
        <Button
          variant="secondary"
          size="md"
          onclick={handleSave}
          disabled={!canSave}
          loading={editorState.status === 'saving'}
        >
          <Save size={16} />
          Save
        </Button>
        <Button
          variant="primary"
          size="md"
          onclick={handlePublish}
          disabled={!canPublish}
          loading={editorState.status === 'publishing'}
        >
          <Globe size={16} />
          Publish
        </Button>
      </div>
    </header>

    <!-- Error Banner -->
    {#if editorState.status === 'error' && editorState.lastError}
      <div class="error-banner">
        <AlertCircle size={18} />
        <span>{editorState.lastError.message}</span>
        <button onclick={retrySave}>Try again</button>
      </div>
    {/if}

    <!-- Main Content -->
    <main class="editor-main">
      <aside class="editor-sidebar">
        <EditorPanel bind:draft />
      </aside>
      <section class="editor-preview">
        <PreviewPanel {draft} />
      </section>
    </main>
  </div>
{/if}

<style>
  .error-page {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 80vh;
    gap: 16px;
    color: var(--lb-error);
    text-align: center;
  }
  .error-page h2 {
    font-size: 24px;
    font-weight: 600;
    color: var(--lb-text);
    margin: 0;
  }
  .error-page p {
    color: var(--lb-muted);
    margin: 0;
  }

  .editor-layout {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background: var(--lb-bg);
  }

  .editor-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 64px;
    padding: 0 20px;
    border-bottom: 1px solid var(--lb-border);
    position: sticky;
    top: 0;
    z-index: 10;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .back-link {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    color: var(--lb-muted);
    border-radius: 10px;
    transition: all 0.15s;
  }
  .back-link:hover {
    background: var(--lb-bg);
    color: var(--lb-text);
  }

  .title-section h1 {
    font-size: 17px;
    font-weight: 600;
    color: var(--lb-text);
    margin: 0;
  }

  .status-row {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 4px;
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .error-banner {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 20px;
    background: rgba(255, 59, 48, 0.1);
    color: var(--lb-error);
    font-size: 14px;
  }
  .error-banner button {
    margin-left: auto;
    background: none;
    border: none;
    color: var(--lb-error);
    font-weight: 600;
    cursor: pointer;
  }
  .error-banner button:hover {
    text-decoration: underline;
  }

  .editor-main {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  .editor-sidebar {
    width: 400px;
    background: var(--lb-surface);
    border-right: 1px solid var(--lb-border);
    overflow-y: auto;
  }

  .editor-preview {
    flex: 1;
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
    overflow-y: auto;
  }

  @media (max-width: 768px) {
    .editor-main {
      flex-direction: column;
    }
    .editor-sidebar {
      width: 100%;
      height: 50%;
      border-right: none;
      border-bottom: 1px solid var(--lb-border);
    }
    .editor-preview {
      height: 50%;
    }
  }
</style>
