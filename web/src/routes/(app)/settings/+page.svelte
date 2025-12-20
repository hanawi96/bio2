<script lang="ts">
	import { getAuth, refreshUser } from '$lib/stores/auth.svelte';
	import { bio, type SocialLinks, pages } from '$lib/api/client';
	import { User, Mail, AtSign, Globe, LogOut, Trash2, Save, AlertCircle } from 'lucide-svelte';
	import { onMount } from 'svelte';

	const auth = getAuth();
	
	// Form state
	let displayName = $state('');
	let bioText = $state('');
	let social = $state<SocialLinks>({});
	
	// Original values for comparison
	let originalDisplayName = $state('');
	let originalBioText = $state('');
	let originalSocial = $state<SocialLinks>({});
	
	// UI state
	let loading = $state(false);
	let saving = $state(false);
	let dirty = $state(false);
	let message = $state('');
	let messageType = $state<'success' | 'error'>('success');

	onMount(async () => {
		loading = true;
		try {
			const data = await bio.get();
			console.log('[Settings] Loaded bio data:', data);
			console.log('[Settings] Page settings:', data.page.settings);
			
			// Load display name from user
			displayName = auth.user?.display_name || '';
			originalDisplayName = displayName;
			console.log('[Settings] Display name:', displayName);
			
			// Load bio and social from page settings
			if (data.page.settings && typeof data.page.settings === 'object') {
				const settings = data.page.settings as any;
				console.log('[Settings] Parsed settings:', settings);
				if (settings.bio) {
					bioText = settings.bio;
					originalBioText = settings.bio;
					console.log('[Settings] Loaded bio:', bioText);
				}
				if (settings.social) {
					social = settings.social;
					originalSocial = { ...settings.social };
					console.log('[Settings] Loaded social:', social);
				}
			}
		} catch (err) {
			console.error('[Settings] Failed to load settings:', err);
		} finally {
			loading = false;
		}
	});

	// Track changes
	$effect(() => {
		// Check if any field changed
		const hasChanges = 
			displayName !== originalDisplayName ||
			bioText !== originalBioText ||
			JSON.stringify(social) !== JSON.stringify(originalSocial);
		dirty = hasChanges;
	});

	async function handleSaveAll() {
		saving = true;
		message = '';
		
		try {
			console.log('[Settings] Saving profile:', { displayName, bioText });
			// Save profile (display name and bio)
			await bio.updateProfile(displayName, bioText);
			console.log('[Settings] Profile saved');
			
			console.log('[Settings] Saving social:', social);
			// Save social links
			await bio.updateSocial(social);
			console.log('[Settings] Social saved');
			
			// Auto-publish bio page to update preview
			const bioData = await bio.get();
			console.log('[Settings] Bio data after save:', bioData);
			if (bioData?.page?.id) {
				try {
					console.log('[Settings] Publishing page:', bioData.page.id);
					await pages.publish(bioData.page.id);
					console.log('[Settings] Page published successfully');
				} catch (err) {
					console.warn('[Settings] Failed to auto-publish bio page:', err);
				}
			}
			
			// Refresh user info in auth store
			await refreshUser();
			console.log('[Settings] User refreshed');
			
			// Update original values
			originalDisplayName = displayName;
			originalBioText = bioText;
			originalSocial = { ...social };
			
			message = 'All settings saved successfully!';
			messageType = 'success';
			dirty = false;
			
			setTimeout(() => message = '', 3000);
		} catch (err) {
			message = 'Failed to save settings';
			messageType = 'error';
			console.error('[Settings] Save error:', err);
		} finally {
			saving = false;
		}
	}

	function updateSocial(key: keyof SocialLinks, value: string) {
		social = { ...social, [key]: value };
		dirty = true;
	}
</script>

<div class="settings-page">
	<header class="settings-header">
		<div class="header-content">
			<div>
				<h1>Settings</h1>
				<p class="text-secondary">Manage your account and preferences</p>
			</div>
			<button 
				class="btn-save-all" 
				onclick={handleSaveAll} 
				disabled={!dirty || saving}
				class:dirty={dirty}
			>
				{#if saving}
					<span class="spinner"></span>
				{:else}
					<Save size={18} />
				{/if}
				<span>{saving ? 'Saving...' : 'Save All'}</span>
			</button>
		</div>
		
		{#if message}
			<div class="message" class:success={messageType === 'success'} class:error={messageType === 'error'}>
				{#if messageType === 'error'}
					<AlertCircle size={16} />
				{/if}
				<span>{message}</span>
			</div>
		{/if}
		
		{#if dirty && !message}
			<div class="unsaved-notice">
				<AlertCircle size={16} />
				<span>You have unsaved changes</span>
			</div>
		{/if}
	</header>

	<!-- Account Info -->
	<div class="settings-section">
		<div class="section-header">
			<h2>Account Information</h2>
		</div>
		<div class="settings-list">
			<div class="settings-item">
				<div class="item-icon">
					<Mail size={20} />
				</div>
				<div class="item-content">
					<span class="item-label">Email</span>
					<span class="item-value">{auth.user?.email}</span>
				</div>
			</div>
			<div class="settings-item">
				<div class="item-icon">
					<AtSign size={20} />
				</div>
				<div class="item-content">
					<span class="item-label">Username</span>
					<span class="item-value">@{auth.user?.username}</span>
				</div>
			</div>
			<div class="settings-item">
				<div class="item-icon">
					<Globe size={20} />
				</div>
				<div class="item-content">
					<span class="item-label">Bio Page URL</span>
					<span class="item-value url">linkbio.com/{auth.user?.username}</span>
				</div>
			</div>
		</div>
	</div>

	<!-- Profile Settings -->
	<div class="settings-section">
		<div class="section-header">
			<h2>Profile Settings</h2>
		</div>
		<div class="settings-list">
			<div class="settings-item column">
				<div class="item-header">
					<User size={18} />
					<span class="item-label">Display Name</span>
				</div>
				<input 
					type="text" 
					class="settings-input" 
					placeholder="Your name" 
					bind:value={displayName}
					disabled={loading}
				/>
			</div>
			<div class="settings-item column">
				<div class="item-header">
					<span class="item-label">Bio</span>
				</div>
				<textarea 
					class="settings-textarea" 
					placeholder="Tell people about yourself..." 
					maxlength="150"
					bind:value={bioText}
					disabled={loading}
				></textarea>
				<span class="char-count">{bioText.length}/150</span>
			</div>
		</div>
	</div>

	<!-- Social Links -->
	<div class="settings-section">
		<div class="section-header">
			<h2>Social Media Links</h2>
			<p class="section-desc">Add your social media profiles to display on your bio page</p>
		</div>
		<div class="settings-list">
			<div class="social-links-inline">
				{#each [
					{ key: 'instagram', label: 'Instagram', placeholder: 'https://instagram.com/username' },
					{ key: 'facebook', label: 'Facebook', placeholder: 'https://facebook.com/username' },
					{ key: 'twitter', label: 'Twitter/X', placeholder: 'https://twitter.com/username' },
					{ key: 'tiktok', label: 'TikTok', placeholder: 'https://tiktok.com/@username' },
					{ key: 'youtube', label: 'YouTube', placeholder: 'https://youtube.com/@username' },
					{ key: 'linkedin', label: 'LinkedIn', placeholder: 'https://linkedin.com/in/username' },
					{ key: 'github', label: 'GitHub', placeholder: 'https://github.com/username' },
					{ key: 'website', label: 'Website', placeholder: 'https://yourwebsite.com' }
				] as platform}
					<div class="social-input-group">
						<label class="social-label">{platform.label}</label>
						<input
							type="url"
							class="settings-input"
							placeholder={platform.placeholder}
							value={social[platform.key as keyof SocialLinks] || ''}
							oninput={(e) => updateSocial(platform.key as keyof SocialLinks, e.currentTarget.value)}
							disabled={loading}
						/>
					</div>
				{/each}
			</div>
		</div>
	</div>

	<!-- Actions -->
	<div class="settings-section">
		<div class="section-header">
			<h2>Actions</h2>
		</div>
		<div class="settings-list">
			<button class="settings-item action">
				<div class="item-icon">
					<LogOut size={20} />
				</div>
				<div class="item-content">
					<span class="item-label">Sign Out</span>
				</div>
			</button>
		</div>
	</div>

	<!-- Danger Zone -->
	<div class="settings-section danger">
		<div class="section-header">
			<h2>Danger Zone</h2>
		</div>
		<div class="settings-list">
			<button class="settings-item action danger">
				<div class="item-icon">
					<Trash2 size={20} />
				</div>
				<div class="item-content">
					<span class="item-label">Delete Account</span>
					<span class="item-desc">Permanently delete your account and all data</span>
				</div>
			</button>
		</div>
	</div>
</div>

<style>
	.settings-page {
		max-width: 640px;
		margin: 0 auto;
		padding: var(--space-4) var(--space-4) var(--space-6);
	}

	.settings-header {
		margin-bottom: var(--space-5);
	}

	.header-content {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: var(--space-4);
		margin-bottom: var(--space-3);
	}

	.btn-save-all {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: 10px 20px;
		background: var(--color-bg-secondary);
		color: var(--color-text-secondary);
		border: 1px solid var(--color-separator);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-weight: 500;
		font-family: var(--font-sans);
		cursor: pointer;
		transition: all 0.2s ease;
		white-space: nowrap;
	}

	.btn-save-all.dirty {
		background: var(--color-primary);
		color: white;
		border-color: var(--color-primary);
	}

	.btn-save-all:hover:not(:disabled) {
		transform: translateY(-1px);
		box-shadow: var(--shadow-sm);
	}

	.btn-save-all:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		transform: none;
	}

	.spinner {
		display: inline-block;
		width: 18px;
		height: 18px;
		border: 2px solid rgba(255,255,255,0.3);
		border-top-color: white;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.message {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-3);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		margin-bottom: var(--space-3);
	}

	.message.success {
		background: rgba(52, 199, 89, 0.1);
		color: var(--color-success);
	}

	.message.error {
		background: rgba(255, 59, 48, 0.1);
		color: var(--color-danger);
	}

	.unsaved-notice {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-3);
		background: rgba(255, 149, 0, 0.1);
		color: var(--color-warning);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
	}

	.settings-header h1 {
		font-size: var(--text-xl);
		font-weight: 600;
		color: var(--color-text);
		margin-bottom: var(--space-1);
	}

	.social-links-inline {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: var(--space-4);
		padding: var(--space-4);
	}

	.social-input-group {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}

	.social-label {
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--color-text);
	}

	.settings-section {
		margin-bottom: var(--space-5);
	}

	.settings-section.danger .section-header h2 {
		color: var(--color-danger);
	}

	.section-header {
		margin-bottom: var(--space-3);
	}

	.section-header h2 {
		font-size: var(--text-sm);
		font-weight: 600;
		color: var(--color-text-secondary);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.section-desc {
		margin-top: var(--space-1);
		font-size: var(--text-sm);
		color: var(--color-text-secondary);
	}

	.settings-list {
		background: var(--color-card-bg);
		border-radius: var(--radius-lg);
		overflow: hidden;
		box-shadow: var(--shadow-sm);
	}

	.settings-item {
		display: flex;
		align-items: center;
		padding: var(--space-4);
		border-bottom: 1px solid var(--color-separator);
		background: transparent;
		width: 100%;
		text-align: left;
		transition: background 0.2s;
	}

	.settings-item:last-child {
		border-bottom: none;
	}

	.settings-item.action {
		cursor: pointer;
	}

	.settings-item.action:hover {
		background: var(--color-bg);
	}

	.settings-item.action.danger:hover {
		background: rgba(255, 59, 48, 0.05);
	}

	.settings-item.column {
		flex-direction: column;
		align-items: stretch;
		gap: var(--space-3);
	}

	.item-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		border-radius: var(--radius-md);
		background: var(--color-bg);
		color: var(--color-text-secondary);
		flex-shrink: 0;
		margin-right: var(--space-3);
	}

	.settings-item.action.danger .item-icon {
		color: var(--color-danger);
		background: rgba(255, 59, 48, 0.1);
	}

	.item-header {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		color: var(--color-text-secondary);
	}

	.item-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}

	.item-label {
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--color-text);
	}

	.item-value {
		font-size: var(--text-sm);
		color: var(--color-text-secondary);
	}

	.item-value.url {
		font-family: ui-monospace, 'SF Mono', Monaco, 'Cascadia Code', monospace;
		font-size: var(--text-xs);
		background: var(--color-bg);
		padding: var(--space-1) var(--space-2);
		border-radius: var(--radius-sm);
		display: inline-block;
	}

	.item-desc {
		font-size: var(--text-xs);
		color: var(--color-text-secondary);
	}

	.settings-input,
	.settings-textarea {
		width: 100%;
		padding: var(--space-3);
		border: 1px solid var(--color-separator);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		background: var(--color-bg);
		color: var(--color-text);
		transition: all 0.2s;
	}

	.settings-input:focus,
	.settings-textarea:focus {
		outline: none;
		border-color: var(--color-primary);
		background: var(--color-card-bg);
		box-shadow: 0 0 0 3px var(--color-primary-light);
	}

	.settings-textarea {
		min-height: 80px;
		resize: vertical;
		font-family: inherit;
	}

	.char-count {
		font-size: var(--text-xs);
		color: var(--color-text-secondary);
		text-align: right;
	}

	@media (max-width: 640px) {
		.settings-page {
			padding: var(--space-3);
		}

		.settings-item {
			padding: var(--space-3);
		}

		.item-icon {
			width: 28px;
			height: 28px;
		}
	}

	@media (max-width: 768px) {
		.header-content {
			flex-direction: column;
			align-items: stretch;
		}

		.btn-save-all {
			width: 100%;
			justify-content: center;
		}

		.social-links-inline {
			grid-template-columns: 1fr;
			padding: var(--space-3);
		}
	}
</style>
