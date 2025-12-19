<script lang="ts">
  import type { CompiledPayload } from '../../types/page';
  import BlockRenderer from './BlockRenderer.svelte';
  import Header from './Header.svelte';

  let { compiled }: { compiled: CompiledPayload } = $props();

  const pageStyle = $derived(() => {
    const layout = compiled.theme?.compiled?.page?.layout;
    const bg = compiled.theme?.compiled?.background;
    
    let style = '';
    if (layout?.pagePadding) style += `padding: ${layout.pagePadding}px;`;
    if (layout?.textAlign) style += `text-align: ${layout.textAlign};`;
    if (bg?.kind === 'color' && bg.color) style += `background-color: ${bg.color};`;
    
    return style;
  });

  const blockGap = $derived(compiled.theme?.compiled?.page?.layout?.blockGap || 12);
</script>

<main class="page" style={pageStyle()}>
  {#if compiled.page?.settings?.header}
    <Header header={compiled.page.settings.header} />
  {/if}

  <div class="blocks" style="gap: {blockGap}px;">
    {#each compiled.blocks as block}
      <BlockRenderer {block} theme={compiled.theme} />
    {/each}
  </div>
</main>

<style>
  .page {
    min-height: 100vh;
    max-width: 680px;
    margin: 0 auto;
  }

  .blocks {
    display: flex;
    flex-direction: column;
  }
</style>
