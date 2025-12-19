<script lang="ts">
  import type { PageData } from './$types';
  import PageRenderer from '$lib/renderer/PageRenderer.svelte';
  import ErrorPage from '$lib/renderer/ErrorPage.svelte';

  let { data }: { data: PageData } = $props();
</script>

<svelte:head>
  <title>{data.compiled?.page?.settings?.header?.name || 'Link Bio'}</title>
</svelte:head>

{#if data.error}
  <ErrorPage error={data.error} />
{:else if data.compiled}
  <PageRenderer compiled={data.compiled} />
{/if}
