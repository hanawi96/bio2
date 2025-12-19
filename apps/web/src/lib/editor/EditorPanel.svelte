<script lang="ts">
  import type { DraftPayload } from '../../types/draft';

  let { draft = $bindable() }: { draft: DraftPayload | null } = $props();

  let newBlockType = $state<'text' | 'link_group'>('text');

  function generateSortKey(): string {
    // Simple sort key generator
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
      draft.blocks = [
        ...draft.blocks,
        {
          id: generateId(),
          type: 'text',
          sort_key: sortKey,
          ref_id: null,
          content: { text: 'New text', variant: 'body' },
          is_visible: true
        }
      ];
    } else if (newBlockType === 'link_group') {
      const groupId = generateId();
      const groupIdStr = String(groupId);

      // Add group
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

      // Add block
      draft.blocks = [
        ...draft.blocks,
        {
          id: generateId(),
          type: 'link_group',
          sort_key: sortKey,
          ref_id: groupId,
          content: {},
          is_visible: true
        }
      ];

      // Initialize empty links array
      draft.links = {
        ...draft.links,
        [groupIdStr]: []
      };
    }
  }

  function removeBlock(index: number) {
    if (!draft) return;
    const block = draft.blocks[index];

    // If link_group, also remove the group and links
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
    // Swap sort_keys
    const tempKey = blocks[index].sort_key;
    blocks[index].sort_key = blocks[newIndex].sort_key;
    blocks[newIndex].sort_key = tempKey;

    // Swap positions
    [blocks[index], blocks[newIndex]] = [blocks[newIndex], blocks[index]];
    draft.blocks = blocks;
  }

  function addLink(groupIdStr: string) {
    if (!draft) return;

    const newLink = {
      id: generateId(),
      title: 'New Link',
      url: 'https://',
      icon_asset_id: null,
      sort_key: generateSortKey(),
      is_active: true
    };

    draft.links = {
      ...draft.links,
      [groupIdStr]: [...(draft.links[groupIdStr] || []), newLink]
    };
  }

  function removeLink(groupIdStr: string, linkIndex: number) {
    if (!draft) return;
    draft.links = {
      ...draft.links,
      [groupIdStr]: draft.links[groupIdStr].filter((_, i) => i !== linkIndex)
    };
  }

  function moveLink(groupIdStr: string, index: number, direction: 'up' | 'down') {
    if (!draft) return;
    const links = [...(draft.links[groupIdStr] || [])];
    const newIndex = direction === 'up' ? index - 1 : index + 1;
    if (newIndex < 0 || newIndex >= links.length) return;

    // Swap sort_keys
    const tempKey = links[index].sort_key;
    links[index].sort_key = links[newIndex].sort_key;
    links[newIndex].sort_key = tempKey;

    [links[index], links[newIndex]] = [links[newIndex], links[index]];
    draft.links = { ...draft.links, [groupIdStr]: links };
  }

  // Sort blocks by sort_key for display
  const sortedBlocks = $derived(
    draft ? [...draft.blocks].sort((a, b) => a.sort_key.localeCompare(b.sort_key)) : []
  );
</script>

{#if draft}
  <div class="editor-panel">
    <section class="section">
      <h3>Page Settings</h3>
      <div class="field">
        <label for="page-title">Title</label>
        <input id="page-title" type="text" bind:value={draft.page.title} />
      </div>
      {#if draft.page.settings?.header}
        <div class="field">
          <label for="header-name">Name</label>
          <input id="header-name" type="text" bind:value={draft.page.settings.header.name} />
        </div>
        <div class="field">
          <label for="header-bio">Bio</label>
          <textarea id="header-bio" bind:value={draft.page.settings.header.bio}></textarea>
        </div>
      {:else}
        <button
          class="btn-init-header"
          onclick={() => {
            if (draft) {
              draft.page.settings = {
                ...draft.page.settings,
                header: { name: '', bio: '', verified: false, social: [], avatar_asset_id: null }
              };
            }
          }}
        >
          + Add Header
        </button>
      {/if}
    </section>

    <section class="section">
      <h3>Blocks</h3>
      <div class="add-block">
        <select bind:value={newBlockType}>
          <option value="text">Text</option>
          <option value="link_group">Link Group</option>
        </select>
        <button onclick={addBlock}>Add Block</button>
      </div>

      <div class="blocks-list">
        {#each sortedBlocks as block, index (block.id)}
          <div class="block-item">
            <div class="block-header">
              <span class="block-type">{block.type}</span>
              <div class="block-actions">
                <button onclick={() => moveBlock(index, 'up')} disabled={index === 0}>↑</button>
                <button onclick={() => moveBlock(index, 'down')} disabled={index === sortedBlocks.length - 1}>↓</button>
                <button onclick={() => removeBlock(index)} class="btn-delete">×</button>
              </div>
            </div>

            {#if block.type === 'text'}
              <div class="block-content">
                <input
                  type="text"
                  value={block.content.text || ''}
                  oninput={(e) => {
                    const target = e.target as HTMLInputElement;
                    block.content = { ...block.content, text: target.value };
                  }}
                  placeholder="Text content"
                />
                <select
                  value={block.content.variant || 'body'}
                  onchange={(e) => {
                    const target = e.target as HTMLSelectElement;
                    block.content = { ...block.content, variant: target.value };
                  }}
                >
                  <option value="body">Body</option>
                  <option value="heading">Heading</option>
                  <option value="caption">Caption</option>
                </select>
              </div>
            {:else if block.type === 'link_group' && block.ref_id}
              {@const groupIdStr = String(block.ref_id)}
              {@const group = draft.link_groups[groupIdStr]}
              {#if group}
                <div class="block-content">
                  <input
                    type="text"
                    value={group.title}
                    oninput={(e) => {
                      const target = e.target as HTMLInputElement;
                      draft.link_groups[groupIdStr] = { ...group, title: target.value };
                    }}
                    placeholder="Group title"
                  />

                  <div class="links-section">
                    <div class="links-header">
                      <span>Links</span>
                      <button onclick={() => addLink(groupIdStr)}>+ Add Link</button>
                    </div>
                    {#each draft.links[groupIdStr] || [] as link, linkIndex (link.id)}
                      <div class="link-item">
                        <input
                          type="text"
                          value={link.title}
                          oninput={(e) => {
                            const target = e.target as HTMLInputElement;
                            draft.links[groupIdStr][linkIndex] = { ...link, title: target.value };
                          }}
                          placeholder="Link title"
                        />
                        <input
                          type="text"
                          value={link.url}
                          oninput={(e) => {
                            const target = e.target as HTMLInputElement;
                            draft.links[groupIdStr][linkIndex] = { ...link, url: target.value };
                          }}
                          placeholder="URL"
                        />
                        <div class="link-actions">
                          <button onclick={() => moveLink(groupIdStr, linkIndex, 'up')} disabled={linkIndex === 0}>↑</button>
                          <button onclick={() => moveLink(groupIdStr, linkIndex, 'down')} disabled={linkIndex === (draft.links[groupIdStr]?.length || 0) - 1}>↓</button>
                          <button onclick={() => removeLink(groupIdStr, linkIndex)} class="btn-delete">×</button>
                        </div>
                      </div>
                    {/each}
                  </div>
                </div>
              {/if}
            {/if}
          </div>
        {/each}
      </div>
    </section>
  </div>
{/if}

<style>
  .editor-panel {
    padding: 16px;
  }

  .section {
    margin-bottom: 24px;
  }

  .section h3 {
    font-size: 14px;
    text-transform: uppercase;
    color: #666;
    margin-bottom: 12px;
  }

  .field {
    margin-bottom: 12px;
  }

  .field label {
    display: block;
    font-size: 12px;
    color: #666;
    margin-bottom: 4px;
  }

  .btn-init-header {
    width: 100%;
    padding: 12px;
    background: #f0f0f0;
    border: 1px dashed #ccc;
    border-radius: 8px;
    cursor: pointer;
    color: #666;
  }

  .btn-init-header:hover {
    background: #e8e8e8;
  }

  .field input,
  .field textarea {
    width: 100%;
    padding: 8px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
  }

  .field textarea {
    min-height: 60px;
    resize: vertical;
  }

  .add-block {
    display: flex;
    gap: 8px;
    margin-bottom: 16px;
  }

  .add-block select {
    flex: 1;
    padding: 8px;
    border: 1px solid #ddd;
    border-radius: 4px;
  }

  .add-block button {
    padding: 8px 16px;
    background: #3b82f6;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .blocks-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .block-item {
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    overflow: hidden;
  }

  .block-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background: #f5f5f5;
  }

  .block-type {
    font-size: 12px;
    font-weight: 500;
    text-transform: uppercase;
    color: #666;
  }

  .block-actions {
    display: flex;
    gap: 4px;
  }

  .block-actions button {
    width: 24px;
    height: 24px;
    border: none;
    background: #e0e0e0;
    border-radius: 4px;
    cursor: pointer;
    font-size: 12px;
  }

  .block-actions button:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  .btn-delete {
    background: #fee2e2 !important;
    color: #dc2626;
  }

  .block-content {
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .block-content input,
  .block-content select {
    padding: 8px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
  }

  .links-section {
    margin-top: 8px;
    padding-top: 8px;
    border-top: 1px solid #e0e0e0;
  }

  .links-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }

  .links-header span {
    font-size: 12px;
    color: #666;
  }

  .links-header button {
    font-size: 12px;
    padding: 4px 8px;
    background: #e0e0e0;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .link-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 8px;
    background: #f9f9f9;
    border-radius: 4px;
    margin-bottom: 8px;
  }

  .link-item input {
    font-size: 13px;
  }

  .link-actions {
    display: flex;
    gap: 4px;
    justify-content: flex-end;
  }

  .link-actions button {
    width: 20px;
    height: 20px;
    border: none;
    background: #e0e0e0;
    border-radius: 4px;
    cursor: pointer;
    font-size: 10px;
  }
</style>
