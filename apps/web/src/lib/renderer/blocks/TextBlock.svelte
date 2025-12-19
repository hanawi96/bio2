<script lang="ts">
  interface TextContent {
    text: string;
    variant: 'body' | 'heading' | 'caption';
  }

  let { content }: { content: TextContent } = $props();
</script>

{#if content.variant === 'heading'}
  <h2 class="text heading">{content.text}</h2>
{:else if content.variant === 'caption'}
  <p class="text caption">{content.text}</p>
{:else}
  <p class="text body">{content.text}</p>
{/if}

<style>
  .text {
    margin: 0;
  }

  .heading {
    font-size: 1.5rem;
    font-weight: 600;
    color: inherit;
  }

  .body {
    font-size: 1rem;
    color: inherit;
  }

  .caption {
    font-size: 0.875rem;
    color: inherit;
    opacity: 0.7;
  }
</style>
