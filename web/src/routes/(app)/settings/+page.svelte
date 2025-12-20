<script lang="ts">
	import { getAuth } from '$lib/stores/auth.svelte';
	import { User, Mail, AtSign, Globe, LogOut, Trash2 } from 'lucide-svelte';

	const auth = getAuth();
</script>

<div class="settings-page">
	<header class="settings-header">
		<h1>Settings</h1>
		<p class="text-secondary">Manage your account and preferences</p>
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
				<input type="text" class="settings-input" placeholder="Your name" value={auth.user?.display_name || ''} />
			</div>
			<div class="settings-item column">
				<div class="item-header">
					<span class="item-label">Bio</span>
				</div>
				<textarea class="settings-textarea" placeholder="Tell people about yourself..." maxlength="150"></textarea>
				<span class="char-count">0/150</span>
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

	.settings-header h1 {
		font-size: var(--text-xl);
		font-weight: 600;
		color: var(--color-text);
		margin-bottom: var(--space-1);
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
</style>
