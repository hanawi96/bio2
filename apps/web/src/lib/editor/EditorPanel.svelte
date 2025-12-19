<script lang="ts">
  import type { DraftPayload } from '../../types/draft';
  import { Input, Button, Select, Card } from '$lib/ui';
  import { Plus, Trash2, ChevronUp, ChevronDown, Link, Type, User } from 'lucide-svelte';

  let { draft = $bindable() }: { draft: DraftPayload | null } = $props();

  let newBlockType = $state<'text' | 'link_group'>('text');

  function generateSortKey(): string {
    const chars = 'abcdefghijklmnopqrstuvwxyz';
    let key = '';
    for (let i = 0; i < 4; i++) {
      key += chars[Math.floor(Math.random() * chars.length)];
    }
    return key;
  }

  function generateId(): number {
    return Date.now() + Math.floor(Math.random() * 1000);
  }

  function addBlock() {
    if (!draft) return;

    const sortKey = generateSortKey();

    if (newBlockType === 'text') {
      draft.blocks = [...draft.blocks, {
        id: generateId(),
        type: 'text',
        sort_key: sortKey,
        ref_id: null,
        content: { text: 'New text', variant: 'body' },
        is_visible: true
      }];
    } else if (newBlockType === 'link_group') {
      const groupId = generateId();
      const groupIdStr = String(groupId);

      draft.link_groups = {
        ...draft.link_groups,
        [groupIdStr]: {
          id: groupId,
          title: 'New Group',
          layout_type: 'list',
          layout_config: {},
          style_override: {}
        }
      };

      draft.blocks = [...draft.blocks, {
        id: generateId(),
        type: 'link_group',
        sort_key: sortKey,
        ref_id: groupId,
        content: {},
        is_visible: true
      }];

      draft.links = { ...draft.links, [groupIdStr]: [] };
    }
  }

  function removeBlock(index: number) {
    if (!draft) return;
    const block = draft.blocks[index];

    if (block.type === 'link_group' && block.ref_id) {
      const groupIdStr = String(block.ref_id);
      const { [groupIdStr]: _, ...restGroups } = draft.link_groups;
      const { [groupIdStr]: __, ...restLinks } = draft.links;
      draft.link_groups = restGroups;
      draft.links = restLinks;
    }

    draft.blocks = draft.blocks.filter((_, i) => i !== index);
  }

  function moveBlock(index: number, direction: 'up' | 'down') {
    if (!draft) return;
    const newIndex = direction === 'up' ? index - 1 : index + 1;
    if (newIndex < 0 || newIndex >= draft.blocks.length) return;

    const blocks = [...draft.blocks];
    const tempKey = blocks[index].sort_key;
    blocks[index].sort_key = blocks[newIndex].sort_key;
    blocks[newIndex].sort_key = tempKey;
    [blocks[index], blocks[newIndex]] = [blocks[newIndex], blocks[index]];
    draft.blocks = blocks;
  }

  function addLink(groupIdStr: string) {
    if (!draft) return;
    draft.links = {
      ...draft.links,
      [groupIdStr]: [...(draft.links[groupIdStr] || []), {
        id: generateId(),
        title: 'New Link',
        url: 'https://',
        icon_asset_id: null,
        sort_key: generateSortKey(),
        is_active: true
      }]
    };
  }

  function removeLink(groupIdStr: string, linkIndex: number) {
    if (!draft) return;
    draft.links = {
      ...draft.links,
      [groupIdStr]: draft.links[groupIdStr].filter((_, i) => i !== linkIndex)
    };
  }

  function moveLink(groupIdStr: string, linkIndex: number, direction: 'up' | 'down') {
    if (!draft) return;
    const links = draft.links[groupIdStr];
    if (!links) return;

    const newIndex = direction === 'up' ? linkIndex - 1 : linkIndex + 1;
    if (newIndex < 0 || newIndex >= links.length) return;

    const newLinks = [...links];
    // Swap sort_key
    const tempKey = newLinks[linkIndex].sort_key;
    newLinks[linkIndex].sort_key = newLinks[newIndex].sort_key;
    newLinks[newIndex].sort_key = tempKey;
    // Swap positions
    [newLinks[linkIndex], newLinks[newIndex]] = [newLinks[newIndex], newLinks[linkIndex]];

    draft.links = {
      ...draft.links,
      [groupIdStr]: newLinks
    };
  }

  const sortedBlocks = $derived(
    draft ? [...draft.blocks].sort((a, b) => a.sort_key.localeCompare(b.sort_key)) : []
  );

  const blockTypeOptions = [
    { value: 'text', label: 'Text Block' },
    { value: 'link_group', label: 'Link Group' }
  ];

  const textVariantOptions = [
    { value: 'body', label: 'Body' },
    { value: 'heading', label: 'Heading' },
    { value: 'caption', label: 'Caption' }
  ];
</script>

{#if draft}
  <div class="editor-panel">
    <!-- Page Settings -->
    <section class="section">
      <h3 class="section-title">
        <User size={16} />
        Page Settings
      </h3>
      <div class="section-content">
        <Input label="Title" bind:value={draft.page.title} placeholder="Page title" />

        {#if draft.page.settings?.header}
          <Input label="Name" bind:value={draft.page.settings.header.name} placeholder="Your name" />
          <div class="field">
            <label class="field-label">Bio</label>
            <textarea
              class="textarea"
              bind:value={draft.page.settings.header.bio}
              placeholder="Short bio..."
              rows="3"
            ></textarea>
          </div>
        {:else}
          <button
            class="add-header-btn"
            onclick={() => {
              if (draft) {
                draft.page.settings = {
                  ...draft.page.settings,
                  header: { name: '', bio: '', verified: false, social: [], avatar_asset_id: null }
                };
              }
            }}
          >
            <Plus size={18} />
            Add Header
          </button>
        {/if}
      </div>
    </section>

    <!-- Blocks -->
    <section class="section">
      <h3 class="section-title">
        <Type size={16} />
        Blocks
      </h3>

      <div class="add-block-row">
        <Select
          bind:value={newBlockType}
          options={blockTypeOptions}
        />
        <Button variant="primary" size="md" onclick={addBlock}>
          <Plus size={16} />
          Add
        </Button>
      </div>

      <div class="blocks-list">
        {#each sortedBlocks as block, index (block.id)}
          <Card variant="outlined" padding="none" class="block-card">
            <div class="block-header">
              <span class="block-type">
                {#if block.type === 'text'}
                  <Type size={14} />
                {:else}
                  <Link size={14} />
                {/if}
                {block.type === 'link_group' ? 'Link Group' : 'Text'}
              </span>
              <div class="block-actions">
                <button class="icon-btn" onclick={() => moveBlock(index, 'up')} disabled={index === 0}>
                  <ChevronUp size={16} />
                </button>
                <button class="icon-btn" onclick={() => moveBlock(index, 'down')} disabled={index === sortedBlocks.length - 1}>
                  <ChevronDown size={16} />
                </button>
                <button class="icon-btn danger" onclick={() => removeBlock(index)}>
                  <Trash2 size={16} />
                </button>
              </div>
            </div>

            <div class="block-body">
              {#if block.type === 'text'}
                <Input
                  value={String(block.content.text || '')}
                  placeholder="Text content"
                  oninput={(e) => {
                    const target = e.target as HTMLInputElement;
                    block.content = { ...block.content, text: target.value };
                  }}
                />
                <Select
                  value={String(block.content.variant || 'body')}
                  options={textVariantOptions}
                  onchange={(v) => { block.content = { ...block.content, variant: v }; }}
                />
              {:else if block.type === 'link_group' && block.ref_id}
                {@const groupIdStr = String(block.ref_id)}
                {@const group = draft.link_groups[groupIdStr]}
                {#if group}
                  <Input
                    value={group.title}
                    placeholder="Group title"
                    oninput={(e) => {
                      const target = e.target as HTMLInputElement;
                      draft.link_groups[groupIdStr] = { ...group, title: target.value };
                    }}
                  />

                  <div class="links-section">
                    <div class="links-header">
                      <span>Links ({draft.links[groupIdStr]?.length || 0})</span>
                      <Button variant="ghost" size="sm" onclick={() => addLink(groupIdStr)}>
                        <Plus size={14} />
                        Add
                      </Button>
                    </div>

                    {#each draft.links[groupIdStr] || [] as link, linkIndex (link.id)}
                      <div class="link-item">
                        <div class="link-actions-row">
                          <button
                            class="link-move-btn"
                            onclick={() => moveLink(groupIdStr, linkIndex, 'up')}
                            disabled={linkIndex === 0}
                            title="Move up"
                          >
                            <ChevronUp size={14} />
                          </button>
                          <button
                            class="link-move-btn"
                            onclick={() => moveLink(groupIdStr, linkIndex, 'down')}
                            disabled={linkIndex === (draft.links[groupIdStr]?.length || 0) - 1}
                            title="Move down"
                          >
                            <ChevronDown size={14} />
                          </button>
                          <button class="link-move-btn danger" onclick={() => removeLink(groupIdStr, linkIndex)} title="Remove">
                            <Trash2 size={14} />
                          </button>
                        </div>
                        <Input
                          value={link.title}
                          placeholder="Link title"
                          oninput={(e) => {
                            const target = e.target as HTMLInputElement;
                            draft.links[groupIdStr][linkIndex] = { ...link, title: target.value };
                          }}
                        />
                        <Input
                          value={link.url}
                          placeholder="https://..."
                          oninput={(e) => {
                            const target = e.target as HTMLInputElement;
                            draft.links[groupIdStr][linkIndex] = { ...link, url: target.value };
                          }}
                        />
                      </div>
                    {/each}
                  </div>
                {/if}
              {/if}
            </div>
          </Card>
        {/each}
      </div>
    </section>
  </div>
{/if}

<style>
  .editor-panel {
    padding: 20px;
  }

  .section {
    margin-bottom: 32px;
  }

  .section-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--lb-muted);
    margin: 0 0 16px;
  }

  .section-content {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .field-label {
    display: block;
    font-size: 13px;
    font-weight: 500;
    color: var(--lb-text);
    margin-bottom: 6px;
  }

  .textarea {
    width: 100%;
    padding: 12px 14px;
    font-size: 15px;
    color: var(--lb-text);
    background: var(--lb-surface);
    border: 1px solid var(--lb-border);
    border-radius: 10px;
    resize: vertical;
    font-family: inherit;
  }

  .textarea:focus {
    outline: none;
    border-color: var(--lb-primary);
    box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.15);
  }

  .add-header-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    width: 100%;
    padding: 16px;
    background: var(--lb-bg);
    border: 2px dashed var(--lb-border);
    border-radius: 12px;
    color: var(--lb-muted);
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s;
  }

  .add-header-btn:hover {
    border-color: var(--lb-primary);
    color: var(--lb-primary);
  }

  .add-block-row {
    display: flex;
    gap: 10px;
    margin-bottom: 16px;
  }

  .add-block-row :global(.select-wrapper) {
    flex: 1;
  }

  .blocks-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  :global(.block-card) {
    overflow: hidden;
  }

  .block-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 14px;
    background: var(--lb-bg);
    border-bottom: 1px solid var(--lb-border);
  }

  .block-type {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
    color: var(--lb-muted);
  }

  .block-actions {
    display: flex;
    gap: 4px;
  }

  .icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    background: var(--lb-surface);
    border: 1px solid var(--lb-border);
    border-radius: 6px;
    color: var(--lb-muted);
    cursor: pointer;
    transition: all 0.15s;
  }

  .icon-btn:hover:not(:disabled) {
    background: var(--lb-border);
    color: var(--lb-text);
  }

  .icon-btn:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  .icon-btn.danger:hover:not(:disabled) {
    background: rgba(255, 59, 48, 0.1);
    border-color: var(--lb-error);
    color: var(--lb-error);
  }

  .block-body {
    padding: 14px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .links-section {
    margin-top: 8px;
    padding-top: 12px;
    border-top: 1px solid var(--lb-border);
  }

  .links-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 10px;
  }

  .links-header span {
    font-size: 12px;
    font-weight: 500;
    color: var(--lb-muted);
  }

  .link-item {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 12px;
    background: var(--lb-bg);
    border-radius: 10px;
    margin-bottom: 8px;
  }

  .link-actions-row {
    display: flex;
    gap: 4px;
    justify-content: flex-end;
    margin-bottom: 4px;
  }

  .link-move-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    background: var(--lb-surface);
    border: 1px solid var(--lb-border);
    border-radius: 4px;
    color: var(--lb-muted);
    cursor: pointer;
    transition: all 0.15s;
  }

  .link-move-btn:hover:not(:disabled) {
    background: var(--lb-border);
    color: var(--lb-text);
  }

  .link-move-btn:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  .link-move-btn.danger:hover:not(:disabled) {
    background: rgba(255, 59, 48, 0.1);
    border-color: var(--lb-error);
    color: var(--lb-error);
  }
</style>
