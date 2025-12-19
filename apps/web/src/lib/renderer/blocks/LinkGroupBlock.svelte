<script lang="ts">
  import type { CompiledLinkGroup, CompiledTheme } from '../../../types/page';

  let { group, theme }: { group: CompiledLinkGroup; theme: CompiledTheme } = $props();

  const style = $derived(() => {
    const s = group.final_style || {};
    let css = '';
    if (s.textAlign) css += `text-align: ${s.textAlign};`;
    if (s.fontSize) {
      const sizes: Record<string, string> = { S: '14px', M: '16px', L: '18px', XL: '22px' };
      css += `font-size: ${sizes[s.fontSize] || '16px'};`;
    }
    return css;
  });

  const linkRadius = $derived(group.final_style?.radius || 16);
  const layoutGap = $derived(group.layout_config?.gap || 12);
</script>

<div class="link-group" style={style()}>
  {#if group.title}
    <h3 class="group-title">{group.title}</h3>
  {/if}

  <div class="links" style="gap: {layoutGap}px;">
    {#each group.links as link}
      <a
        href={link.url}
        target="_blank"
        rel="noopener noreferrer"
        class="link-item"
        style="border-radius: {linkRadius}px;"
      >
        {#if link.icon_url}
          <img src={link.icon_url} alt="" class="link-icon" />
        {/if}
        <span class="link-title">{link.title}</span>
      </a>
    {/each}
  </div>
</div>

<style>
  .link-group {
    width: 100%;
  }

  .group-title {
    font-size: 1rem;
    font-weight: 500;
    margin-bottom: 8px;
    opacity: 0.8;
  }

  .links {
    display: flex;
    flex-direction: column;
  }

  .link-item {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 14px 20px;
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    color: inherit;
    text-decoration: none;
    transition: transform 0.15s ease, background 0.15s ease;
  }

  .link-item:hover {
    transform: scale(1.02);
    background: rgba(255, 255, 255, 0.15);
  }

  .link-icon {
    width: 24px;
    height: 24px;
    border-radius: 4px;
  }

  .link-title {
    font-weight: 500;
  }
</style>
