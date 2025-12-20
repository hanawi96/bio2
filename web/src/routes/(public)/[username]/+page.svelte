<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	interface CompiledPage {
		page: {
			id: number;
			title?: string;
			locale: string;
			mode: string;
			settings: object;
		};
		theme: object;
		blocks: CompiledBlock[];
	}

	interface CompiledBlock {
		id: number;
		type: string;
		content: object;
		is_visible: boolean;
		group?: {
			id: number;
			title?: string;
			layout_type: string;
			links: { id: number; title: string; url: string }[];
		};
	}

	let data = $state<CompiledPage | null>(null);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		const username = $page.params.username;
		try {
			const res = await fetch(`http://localhost:8080/r?path=/${username}`);
			if (!res.ok) throw new Error('Page not found');
			data = await res.json();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load';
		} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	{#if data?.page.title}
		<title>{data.page.title}</title>
	{/if}
</svelte:head>

<main class="public-page">
	{#if loading}
		<p class="loading">Loading...</p>
	{:else if error}
		<p class="error">{error}</p>
	{:else if data}
		<div class="page-content">
			{#if data.page.title}
				<h1>{data.page.title}</h1>
			{/if}

			{#each data.blocks as block}
				{#if block.type === 'link_group' && block.group}
					<div class="link-group">
						{#if block.group.title}
							<h2>{block.group.title}</h2>
						{/if}
						<div class="links">
							{#each block.group.links as link}
								<a href={link.url} class="link-item" target="_blank" rel="noopener">
									{link.title}
								</a>
							{/each}
						</div>
					</div>
				{:else if block.type === 'text'}
					<div class="text-block">
						<p>{(block.content as { text?: string }).text || ''}</p>
					</div>
				{:else if block.type === 'spacer'}
					<div class="spacer"></div>
				{/if}
			{/each}
		</div>
	{/if}
</main>

<style>
	.public-page {
		min-height: 100vh;
		display: flex;
		justify-content: center;
		padding: 2rem 1rem;
	}

	.loading, .error {
		text-align: center;
		padding: 3rem;
	}

	.error {
		color: #dc2626;
	}

	.page-content {
		width: 100%;
		max-width: 480px;
	}

	h1 {
		text-align: center;
		margin-bottom: 2rem;
	}

	.link-group {
		margin-bottom: 1.5rem;
	}

	.link-group h2 {
		font-size: 0.875rem;
		color: #666;
		margin-bottom: 0.75rem;
	}

	.links {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.link-item {
		display: block;
		padding: 1rem;
		background: white;
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
		text-align: center;
		color: var(--color-text);
		transition: transform 0.1s, box-shadow 0.1s;
	}

	.link-item:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0,0,0,0.1);
	}

	.text-block {
		margin-bottom: 1rem;
		text-align: center;
	}

	.spacer {
		height: 1.5rem;
	}
</style>
