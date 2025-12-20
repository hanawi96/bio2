<script lang="ts">
	import { bio, type SocialLinks } from '$lib/api/client';
	import { Instagram, Facebook, Twitter, Music, Youtube, Linkedin, Github, Globe, Save } from 'lucide-svelte';
	import { onMount } from 'svelte';

	let social = $state<SocialLinks>({});
	let loading = $state(false);
	let saving = $state(false);
	let message = $state('');

	const platforms = [
		{ key: 'instagram', label: 'Instagram', icon: Instagram, placeholder: 'https://instagram.com/username' },
		{ key: 'facebook', label: 'Facebook', icon: Facebook, placeholder: 'https://facebook.com/username' },
		{ key: 'twitter', label: 'Twitter/X', icon: Twitter, placeholder: 'https://twitter.com/username' },
		{ key: 'tiktok', label: 'TikTok', icon: Music, placeholder: 'https://tiktok.com/@username' },
		{ key: 'youtube', label: 'YouTube', icon: Youtube, placeholder: 'https://youtube.com/@username' },
		{ key: 'linkedin', label: 'LinkedIn', icon: Linkedin, placeholder: 'https://linkedin.com/in/username' },
		{ key: 'github', label: 'GitHub', icon: Github, placeholder: 'https://github.com/username' },
		{ key: 'website', label: 'Website', icon: Globe, placeholder: 'https://yourwebsite.com' }
	];

	onMount(async () => {
		loading = true;
		try {
			const data = await bio.get();
			if (data.page.settings && typeof data.page.settings === 'object') {
				const settings = data.page.settings as any;
				if (settings.social) {
					social = settings.social;
				}
			}
		} catch (err) {
			console.error('Failed to load social links:', err);
		} finally {
			loading = false;
		}
	});

	async function handleSave() {
		saving = true;
		message = '';
		try {
			await bio.updateSocial(social);
			message = 'Social links saved successfully!';
			setTimeout(() => message = '', 3000);
		} catch (err) {
			message = 'Failed to save social links';
			console.error(err);
		} finally {
			saving = false;
		}
	}
</script>

<div class="social-links-editor">
	{#if loading}
		<div class="loading">Loading...</div>
	{:else}
		<div class="social-grid">
			{#each platforms as platform}
				<div class="social-input-group">
					<div class="input-header">
						<svelte:component this={platform.icon} size={18} />
						<span>{platform.label}</span>
					</div>
					<input
						type="url"
						class="social-input"
						placeholder={platform.placeholder}
						bind:value={social[platform.key as keyof SocialLinks]}
					/>
				</div>
			{/each}
		</div>

		<div class="actions">
			<button class="btn-save" onclick={handleSave} disabled={saving}>
				{#if saving}
					<span class="spinner"></span>
				{:else}
					<Save size={16} />
				{/if}
				<span>{saving ? 'Saving...' : 'Save Social Links'}</span>
			</button>
		</div>

		{#if message}
			<div class="message" class:success={message.includes('success')}>
				{message}
			</div>
		{/if}
	{/if}
</div>

<style>
	.social-links-editor {
		width: 100%;
		padding: var(--space-4);
	}

	.loading {
		padding: var(--space-4);
		text-align: center;
		color: var(--color-text-secondary);
		font-size: var(--text-sm);
	}

	.social-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: var(--space-4);
		margin-bottom: var(--space-4);
	}

	.social-input-group {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}

	.input-header {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--color-text);
	}

	.social-input {
		width: 100%;
		padding: 10px 12px;
		border: 1px solid var(--color-separator);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-family: var(--font-sans);
		background: var(--color-surface);
		color: var(--color-text);
		transition: all 0.2s ease;
		box-sizing: border-box;
	}

	.social-input:focus {
		outline: none;
		border-color: var(--color-primary);
		box-shadow: 0 0 0 3px var(--color-primary-light);
	}

	.social-input::placeholder {
		color: var(--color-text-secondary);
	}

	.actions {
		display: flex;
		justify-content: flex-end;
		padding-top: var(--space-3);
		border-top: 1px solid var(--color-separator);
	}

	.btn-save {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: 10px 20px;
		background: var(--color-primary);
		color: white;
		border: none;
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-weight: 500;
		font-family: var(--font-sans);
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.btn-save:hover:not(:disabled) {
		opacity: 0.9;
		transform: translateY(-1px);
	}

	.btn-save:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.spinner {
		display: inline-block;
		width: 16px;
		height: 16px;
		border: 2px solid rgba(255,255,255,0.3);
		border-top-color: white;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.message {
		margin-top: var(--space-3);
		padding: var(--space-3);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		text-align: center;
	}

	.message.success {
		background: rgba(52, 199, 89, 0.1);
		color: var(--color-success);
	}

	@media (max-width: 768px) {
		.social-grid {
			grid-template-columns: 1fr;
		}
		
		.social-links-editor {
			padding: var(--space-3);
		}
	}
</style>
