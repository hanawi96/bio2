<script lang="ts">
  import type { DraftPayload } from '../../types/draft';
  import type { CompiledTheme, CompiledLinkGroup } from '../../types/page';
  import Header from '$lib/renderer/Header.svelte';
  import TextBlock from '$lib/renderer/blocks/TextBlock.svelte';
  import LinkGroupBlock from '$lib/renderer/blocks/LinkGroupBlock.svelte';

  let { draft }: { draft: DraftPayload | null } = $props();

  interface PreviewTextBlock {
    type: 'text';
    content: { text: string; variant: 'body' | 'heading' | 'caption' };
  }

  interface PreviewLinkGroupBlock {
    type: 'link_group';
    group: CompiledLinkGroup;
  }

  type PreviewBlock = PreviewTextBlock | PreviewLinkGroupBlock;

  // Convert draft to preview format
  const previewBlocks = $derived((): PreviewBlock[] => {
    if (!draft) return [];

    const blocks: PreviewBlock[] = [];

    [...draft.blocks]
      .filter((b) => b.is_visible)
      .sort((a, b) => a.sort_key.localeCompare(b.sort_key))
      .forEach((block) => {
        if (block.type === 'text') {
          blocks.push({
            type: 'text',
            content: {
              text: String(block.content.text || ''),
              variant: (block.content.variant as 'body' | 'heading' | 'caption') || 'body'
            }
          });
        } else if (block.type === 'link_group' && block.ref_id) {
          const groupIdStr = String(block.ref_id);
          const group = draft.link_groups[groupIdStr];
          if (!group) return;

          const links = (draft.links[groupIdStr] || [])
            .filter((l) => l.is_active)
            .sort((a, b) => a.sort_key.localeCompare(b.sort_key))
            .map((l) => ({
              title: l.title,
              url: l.url,
              icon_url: null
            }));

          blocks.push({
            type: 'link_group',
            group: {
              id: group.id,
              title: group.title,
              layout_type: (group.layout_type as 'list' | 'cards' | 'grid') || 'list',
              layout_config: group.layout_config,
              final_style: group.style_override as { textAlign?: string; fontSize?: string; radius?: number },
              links
            }
          });
        }
      });

    return blocks;
  });

  const previewTheme = $derived((): CompiledTheme => ({
    preset_key: draft?.page?.theme?.preset_key || 'theme_a',
    mode: (draft?.page?.theme?.mode as 'light' | 'dark' | 'compact') || 'light',
    compiled: {
      page: {
        layout: { textAlign: 'center', baseFontSize: 'M', pagePadding: 16, blockGap: 12 },
        defaults: { linkGroup: { textAlign: 'center', fontSize: 'M', radius: 16 } }
      },
      background: { kind: 'color', color: '#0B0F19' }
    }
  }));
</script>

<div class="preview-panel">
  <div class="preview-frame">
    {#if draft?.page?.settings?.header}
      <Header header={draft.page.settings.header} />
    {/if}

    <div class="blocks" style="gap: 12px;">
      {#each previewBlocks() as block}
        {#if block.type === 'text'}
          <TextBlock content={block.content} />
        {:else if block.type === 'link_group'}
          <LinkGroupBlock group={block.group} theme={previewTheme()} />
        {/if}
      {/each}
    </div>
  </div>
</div>

<style>
  .preview-panel {
    display: flex;
    justify-content: center;
    padding: 24px;
    min-height: 100%;
  }

  .preview-frame {
    width: 100%;
    max-width: 420px;
    background: #0b0f19;
    border-radius: 24px;
    padding: 16px;
    color: white;
    min-height: 600px;
  }

  .blocks {
    display: flex;
    flex-direction: column;
  }
</style>
