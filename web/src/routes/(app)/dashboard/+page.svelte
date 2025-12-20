<script lang="ts">
	import { onMount } from 'svelte';
	import { bio, type BioData } from '$lib/api/client';
	import { getAuth } from '$lib/stores/auth.svelte';

	const auth = getAuth();

	let bioData = $state<BioData | null>(null);
	let loading = $state(true);

	onMount(async () => {
		try {
			bioData = await bio.get();
		} catch (err) {
			console.error(err);
		} finally {
			loading = false;
		}
	});

	const totalBlocks = $derived(bioData?.blocks.length || 0);
	const totalLinks = $derived(
		bioData?.blocks
			.filter((b) => b.type === 'link_group' && b.group)
			.reduce((sum, b) => sum + (b.group?.links.length || 0), 0) || 0
	);
	const activeLinks = $derived(
		bioData?.blocks
			.filter((b) => b.type === 'link_group' && b.group)
			.reduce((sum, b) => sum + (b.group?.links.filter((l) => l.is_active).length || 0), 0) || 0
	);
	const recentLinks = $derived(
		bioData?.blocks
			.filter((b) => b.type === 'link_group' && b.group)
			.flatMap((b) => b.group?.links || [])
			.slice(0, 5) || []
	);
</script>

<div class="dashboard">
	<!-- Welcome -->
	<div class="welcome">
		<h1>Xin chào{auth.user?.username ? `, @${auth.user.username}` : ''}! 👋</h1>
		<p>Quản lý trang bio của bạn từ đây</p>
	</div>

	<!-- Stats Grid -->
	<div class="stats-grid">
		<div class="stat-card">
			<div class="stat-icon blue">📦</div>
			<div class="stat-info">
				<span class="stat-value">{totalBlocks}</span>
				<span class="stat-label">Blocks</span>
			</div>
		</div>
		<div class="stat-card">
			<div class="stat-icon purple">🔗</div>
			<div class="stat-info">
				<span class="stat-value">{totalLinks}</span>
				<span class="stat-label">Links</span>
			</div>
		</div>
		<div class="stat-card">
			<div class="stat-icon green">👁️</div>
			<div class="stat-info">
				<span class="stat-value">{activeLinks}</span>
				<span class="stat-label">Hiển thị</span>
			</div>
		</div>
		<div class="stat-card">
			<div class="stat-icon orange">👆</div>
			<div class="stat-info">
				<span class="stat-value">0</span>
				<span class="stat-label">Clicks</span>
			</div>
		</div>
	</div>

	<!-- Quick Actions -->
	<div class="section">
		<div class="section-title">Bắt đầu nhanh</div>
		<div class="actions-grid">
			<a href="/bio" class="action-card">
				<div class="action-icon">📝</div>
				<div class="action-info">
					<span class="action-title">Chỉnh sửa Bio</span>
					<span class="action-desc">Thêm links, text, hình ảnh</span>
				</div>
				<span class="action-arrow">→</span>
			</a>
			<a href="/appearance" class="action-card">
				<div class="action-icon">🎨</div>
				<div class="action-info">
					<span class="action-title">Tùy chỉnh giao diện</span>
					<span class="action-desc">Theme, màu sắc, font</span>
				</div>
				<span class="action-arrow">→</span>
			</a>
			{#if auth.user?.username}
				<a href="/{auth.user.username}" target="_blank" class="action-card">
					<div class="action-icon">👁️</div>
					<div class="action-info">
						<span class="action-title">Xem trang Bio</span>
						<span class="action-desc">linkbio.com/{auth.user.username}</span>
					</div>
					<span class="action-arrow">↗</span>
				</a>
			{/if}
		</div>
	</div>

	<!-- Recent Links -->
	<div class="section">
		<div class="section-header-row">
			<div class="section-title">Links gần đây</div>
			<a href="/bio" class="section-link">Xem tất cả</a>
		</div>

		{#if loading}
			<div class="loading-box">
				<div class="spinner"></div>
			</div>
		{:else if recentLinks.length === 0}
			<div class="empty-box">
				<div class="empty-icon">🔗</div>
				<p>Chưa có link nào</p>
				<a href="/bio" class="btn-primary">Thêm link đầu tiên</a>
			</div>
		{:else}
			<div class="links-card">
				{#each recentLinks as link}
					<div class="link-row" class:inactive={!link.is_active}>
						<div class="link-icon">🔗</div>
						<div class="link-info">
							<span class="link-title">{link.title}</span>
							<span class="link-url">{link.url}</span>
						</div>
						<div class="link-status" class:active={link.is_active}>
							{link.is_active ? 'Hiện' : 'Ẩn'}
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<style>
	.dashboard {
		max-width: 900px;
		margin: 0 auto;
	}

	/* Welcome */
	.welcome {
		margin-bottom: var(--space-6);
	}

	.welcome h1 {
		font-size: var(--text-2xl);
		font-weight: 700;
		color: var(--color-text);
		margin-bottom: var(--space-1);
	}

	.welcome p {
		font-size: var(--text-base);
		color: var(--color-text-secondary);
	}

	/* Stats */
	.stats-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: var(--space-4);
		margin-bottom: var(--space-6);
	}

	.stat-card {
		background: var(--color-bg-secondary);
		border-radius: var(--radius-lg);
		padding: var(--space-4);
		display: flex;
		align-items: center;
		gap: var(--space-3);
		box-shadow: var(--shadow-card);
	}

	.stat-icon {
		width: 44px;
		height: 44px;
		border-radius: var(--radius-md);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.25rem;
	}

	.stat-icon.blue {
		background: rgba(0, 122, 255, 0.1);
	}
	.stat-icon.purple {
		background: rgba(88, 86, 214, 0.1);
	}
	.stat-icon.green {
		background: rgba(52, 199, 89, 0.1);
	}
	.stat-icon.orange {
		background: rgba(255, 149, 0, 0.1);
	}

	.stat-info {
		display: flex;
		flex-direction: column;
	}

	.stat-value {
		font-size: var(--text-xl);
		font-weight: 700;
		color: var(--color-text);
	}

	.stat-label {
		font-size: var(--text-xs);
		color: var(--color-text-secondary);
	}

	/* Sections */
	.section {
		margin-bottom: var(--space-6);
	}

	.section-title {
		font-size: var(--text-xs);
		font-weight: 600;
		color: var(--color-text-secondary);
		text-transform: uppercase;
		letter-spacing: 0.5px;
		margin-bottom: var(--space-3);
	}

	.section-header-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--space-3);
	}

	.section-link {
		font-size: var(--text-sm);
		color: var(--color-primary);
	}

	/* Actions */
	.actions-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: var(--space-3);
	}

	.action-card {
		background: var(--color-bg-secondary);
		border-radius: var(--radius-lg);
		padding: var(--space-4);
		display: flex;
		align-items: center;
		gap: var(--space-3);
		box-shadow: var(--shadow-card);
		transition: all var(--transition-base);
	}

	.action-card:hover {
		transform: translateY(-2px);
		box-shadow: var(--shadow-md);
		opacity: 1;
	}

	.action-card:active {
		transform: translateY(0);
	}

	.action-icon {
		font-size: 1.5rem;
	}

	.action-info {
		flex: 1;
	}

	.action-title {
		display: block;
		font-size: var(--text-sm);
		font-weight: 600;
		color: var(--color-text);
		margin-bottom: 2px;
	}

	.action-desc {
		font-size: var(--text-xs);
		color: var(--color-text-secondary);
	}

	.action-arrow {
		color: var(--color-text-tertiary);
		font-size: var(--text-lg);
	}

	/* Links Card */
	.links-card {
		background: var(--color-bg-secondary);
		border-radius: var(--radius-lg);
		overflow: hidden;
		box-shadow: var(--shadow-card);
	}

	.link-row {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-3) var(--space-4);
		border-bottom: 1px solid var(--color-separator);
		transition: background var(--transition-fast);
	}

	.link-row:last-child {
		border-bottom: none;
	}

	.link-row:hover {
		background: var(--color-bg);
	}

	.link-row.inactive {
		opacity: 0.5;
	}

	.link-icon {
		font-size: 1rem;
	}

	.link-info {
		flex: 1;
		min-width: 0;
	}

	.link-title {
		display: block;
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--color-text);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.link-url {
		font-size: var(--text-xs);
		color: var(--color-text-secondary);
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.link-status {
		font-size: var(--text-xs);
		padding: 2px var(--space-2);
		border-radius: var(--radius-full);
		background: var(--color-bg);
		color: var(--color-text-secondary);
	}

	.link-status.active {
		background: rgba(52, 199, 89, 0.1);
		color: var(--color-success);
	}

	/* Loading & Empty */
	.loading-box {
		display: flex;
		justify-content: center;
		padding: var(--space-8);
	}

	.spinner {
		width: 24px;
		height: 24px;
		border: 2px solid var(--color-border);
		border-top-color: var(--color-primary);
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.empty-box {
		background: var(--color-bg-secondary);
		border-radius: var(--radius-lg);
		padding: var(--space-8);
		text-align: center;
		box-shadow: var(--shadow-card);
	}

	.empty-icon {
		font-size: 2.5rem;
		margin-bottom: var(--space-3);
	}

	.empty-box p {
		color: var(--color-text-secondary);
		margin-bottom: var(--space-4);
	}

	/* Responsive */
	@media (max-width: 768px) {
		.stats-grid {
			grid-template-columns: repeat(2, 1fr);
		}

		.actions-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
