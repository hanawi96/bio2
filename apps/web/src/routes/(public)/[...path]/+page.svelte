<script lang="ts">
  import type { PageData } from './$types';
  import PageRenderer from '$lib/renderer/PageRenderer.svelte';
  import PasswordPrompt from '$lib/renderer/PasswordPrompt.svelte';
  import ErrorPage from '$lib/renderer/ErrorPage.svelte';

  let { data }: { data: PageData } = $props();
</script>

<svelte:head>
  <title>{data.compiled?.page?.settings?.header?.name || 'Link Bio'}</title>
</svelte:head>

{#if data.error}
  {#if data.error.code === 'PASSWORD_REQUIRED'}
    <PasswordPrompt />
  {:else}
    <ErrorPage error={data.error} />
  {/if}
{:else if data.compiled}
  <PageRenderer compiled={data.compiled} />
{/if}
