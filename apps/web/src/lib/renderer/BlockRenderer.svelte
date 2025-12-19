<script lang="ts">
  import type { CompiledBlock, CompiledTheme } from '../../types/page';
  import TextBlock from './blocks/TextBlock.svelte';
  import LinkGroupBlock from './blocks/LinkGroupBlock.svelte';

  let { block, theme }: { block: CompiledBlock; theme: CompiledTheme } = $props();
</script>

{#if block.type === 'text'}
  <TextBlock content={block.content} />
{:else if block.type === 'link_group'}
  <LinkGroupBlock group={block.group} {theme} />
{:else}
  <div class="unsupported">Block type "{block.type}" not supported</div>
{/if}

<style>
  .unsupported {
    padding: 12px;
    background: #f0f0f0;
    border-radius: 8px;
    color: #666;
    font-size: 14px;
  }
</style>
