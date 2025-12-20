<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { getEditor, loadDraft, save, publish, addBlock, addLinkGroup, addLink, deleteBlock } from '$lib/stores/editor.svelte';

	const editor = getEditor();
	let loading = $state(true);
	let publishing = $state(false);

	$effect(() => {
		const id = Number($page.params.id);
		if (id) {
			loadDraft(id).finally(() => loading = false);
		}
	});

	async function handleSave() {
		await save();
	}

	async function handlePublish() {
		publishing = true;
		try {
			await publish();
			alert('Published!');
		} finally {
			publishing = false;
		}
	}

	function handleAddLinkGroup() {
		const group = addLinkGroup('Links');
		if (group) {
			addBlock('link_group');
		}
	}
</script>

<div class="editor">
	{#if loading}
		<p>Loading...</p>
	{:else if editor.draft}
		<div class="editor-header">
			<h1>{editor.draft.page.title || 'Untitled'}</h1>
			<div class="actions">
				<button class="btn-secondary" onclick={handleSave} disabled={editor.saving || !editor.dirty}>
					{editor.saving ? 'Saving...' : 'Save'}
				</button>
				<button class="btn-primary" onclick={handlePublish} disabled={publishing}>
					{publishing ? 'Publishing...' : 'Publish'}
				</button>
			</div>
		</div>

		<div class="editor-body">
			<div class="sidebar">
				<h3>Add Block</h3>
				<div class="block-buttons">
					<button onclick={() => addBlock('text')}>Text</button>
					<button onclick={handleAddLinkGroup}>Link Group</button>
					<button onclick={() => addBlock('image')}>Image</button>
					<button onclick={() => addBlock('spacer')}>Spacer</button>
				</div>

				<h3>Settings</h3>
				<a href="/pages/{editor.draft.page.id}/appearance">Appearance</a>
				<a href="/pages/{editor.draft.page.id}/settings">Settings</a>
			</div>

			<div class="preview">
				<div class="preview-frame">
					{#each editor.blocks as block}
						<div class="block" class:hidden={!block.is_visible}>
							<div class="block-header">
								<span>{block.type}</span>
								<button onclick={() => deleteBlock(block.id)}>×</button>
							</div>
							
							{#if block.type === 'text'}
								<div class="block-content">
									<textarea placeholder="Enter text..."></textarea>
								</div>
							{:else if block.type === 'link_group' && block.ref_id}
								<div class="block-content">
									<div class="link-group">
										{#each editor.getLinks(block.ref_id) as link}
											<div class="link-item">
												<span>{link.title}</span>
												<span class="url">{link.url}</span>
											</div>
										{/each}
										<button class="add-link" onclick={() => {
											const title = prompt('Link title');
											const url = prompt('Link URL');
											if (title && url && block.ref_id) {
												addLink(block.ref_id, title, url);
											}
										}}>+ Add Link</button>
									</div>
								</div>
							{:else if block.type === 'spacer'}
								<div class="block-content spacer"></div>
							{:else}
								<div class="block-content">
									<p class="placeholder">{block.type} block</p>
								</div>
							{/if}
						</div>
					{/each}

					{#if editor.blocks.length === 0}
						<p class="empty">Add blocks to get started</p>
					{/if}
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	.editor {
		max-width: 1200px;
		margin: 0 auto;
	}

	.editor-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
	}

	.actions {
		display: flex;
		gap: 0.5rem;
	}

	.editor-body {
		display: grid;
		grid-template-columns: 250px 1fr;
		gap: 2rem;
	}

	.sidebar h3 {
		font-size: 0.875rem;
		color: #666;
		margin-bottom: 0.5rem;
		margin-top: 1.5rem;
	}

	.sidebar h3:first-child {
		margin-top: 0;
	}

	.block-buttons {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.block-buttons button {
		text-align: left;
		background: var(--color-border);
	}

	.sidebar a {
		display: block;
		padding: 0.5rem 0;
		color: var(--color-text);
	}

	.preview {
		background: #f5f5f5;
		border-radius: var(--radius);
		padding: 2rem;
		min-height: 600px;
	}

	.preview-frame {
		max-width: 400px;
		margin: 0 auto;
		background: white;
		border-radius: var(--radius);
		padding: 1rem;
		min-height: 500px;
	}

	.block {
		margin-bottom: 1rem;
		border: 1px solid var(--color-border);
		border-radius: var(--radius);
	}

	.block.hidden {
		opacity: 0.5;
	}

	.block-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.5rem;
		background: #f9fafb;
		border-bottom: 1px solid var(--color-border);
		font-size: 0.75rem;
		color: #666;
	}

	.block-header button {
		background: none;
		padding: 0;
		width: 20px;
		height: 20px;
		color: #999;
	}

	.block-content {
		padding: 0.75rem;
	}

	.block-content textarea {
		min-height: 60px;
		resize: vertical;
	}

	.block-content.spacer {
		height: 32px;
		background: repeating-linear-gradient(
			45deg,
			#f0f0f0,
			#f0f0f0 10px,
			#e5e5e5 10px,
			#e5e5e5 20px
		);
	}

	.placeholder {
		color: #999;
		text-align: center;
		padding: 1rem;
	}

	.link-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.link-item {
		display: flex;
		justify-content: space-between;
		padding: 0.5rem;
		background: #f9fafb;
		border-radius: 4px;
	}

	.link-item .url {
		font-size: 0.75rem;
		color: #666;
	}

	.add-link {
		background: none;
		color: var(--color-primary);
		text-align: left;
		padding: 0.5rem;
	}

	.empty {
		text-align: center;
		color: #999;
		padding: 3rem;
	}
</style>
