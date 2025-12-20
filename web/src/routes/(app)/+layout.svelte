<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getAuth, logout } from '$lib/stores/auth.svelte';

	let { children } = $props();
	const auth = getAuth();

	$effect(() => {
		if (auth.loading) return;

		const currentPath = $page.url.pathname;

		if (!auth.isLoggedIn) {
			goto('/login');
			return;
		}

		if (auth.needsSetup && currentPath !== '/setup') {
			goto('/setup');
			return;
		}

		if (!auth.needsSetup && currentPath === '/setup') {
			goto('/dashboard');
		}
	});

	async function handleLogout() {
		await logout();
		goto('/');
	}

	const isSetupPage = $derived($page.url.pathname === '/setup');
	const currentPath = $derived($page.url.pathname);

	const navItems = [
		{ href: '/dashboard', icon: 'house.fill', label: 'Dashboard' },
		{ href: '/bio', icon: 'link', label: 'My Bio' },
		{ href: '/appearance', icon: 'paintbrush.fill', label: 'Giao diện' },
		{ href: '/settings', icon: 'gearshape.fill', label: 'Cài đặt' }
	];
</script>

{#if auth.loading}
	<div class="loading-screen">
		<div class="loading-spinner"></div>
	</div>
{:else if auth.isLoggedIn}
	{#if isSetupPage}
		{@render children()}
	{:else}
		<div class="app-shell">
			<!-- Sidebar -->
			<aside class="sidebar">
				<div class="sidebar-brand">
					<div class="brand-icon">⚡</div>
					<span class="brand-text">LinkBio</span>
				</div>

				<nav class="sidebar-nav">
					{#each navItems as item}
						<a
							href={item.href}
							class="nav-link"
							class:active={currentPath === item.href || currentPath.startsWith(item.href + '/')}
						>
							<span class="nav-icon">
								{#if item.icon === 'house.fill'}🏠
								{:else if item.icon === 'link'}🔗
								{:else if item.icon === 'paintbrush.fill'}🎨
								{:else if item.icon === 'gearshape.fill'}⚙️
								{/if}
							</span>
							<span class="nav-label">{item.label}</span>
						</a>
					{/each}
				</nav>

				<div class="sidebar-footer">
					{#if auth.user?.username}
						<a href="/{auth.user.username}" target="_blank" class="preview-btn">
							<span>👁️</span>
							<span>Xem trang bio</span>
						</a>
					{/if}
				</div>
			</aside>

			<!-- Main -->
			<div class="main-area">
				<header class="topbar">
					<div class="topbar-left">
						{#if auth.user?.username}
							<div class="url-badge">
								<span class="url-icon">🌐</span>
								<span class="url-text">linkbio.com/{auth.user.username}</span>
							</div>
						{/if}
					</div>
					<div class="topbar-right">
						<button class="user-btn" onclick={handleLogout}>
							<span class="user-avatar">
								{auth.user?.email?.charAt(0).toUpperCase()}
							</span>
							<span class="user-name">{auth.user?.email}</span>
							<span class="logout-icon">↗</span>
						</button>
					</div>
				</header>

				<main class="content-area">
					{@render children()}
				</main>
			</div>
		</div>
	{/if}
{/if}

<style>
	/* Loading Screen */
	.loading-screen {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-bg);
	}

	.loading-spinner {
		width: 32px;
		height: 32px;
		border: 3px solid var(--color-border);
		border-top-color: var(--color-primary);
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	/* App Shell */
	.app-shell {
		display: flex;
		min-height: 100vh;
		background: var(--color-bg);
	}

	/* Sidebar */
	.sidebar {
		width: 260px;
		background: var(--color-bg-secondary);
		border-right: 1px solid var(--color-separator);
		display: flex;
		flex-direction: column;
		position: fixed;
		top: 0;
		left: 0;
		bottom: 0;
		z-index: 100;
	}

	.sidebar-brand {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-5);
		border-bottom: 1px solid var(--color-separator);
	}

	.brand-icon {
		width: 36px;
		height: 36px;
		background: linear-gradient(135deg, #007aff 0%, #5856d6 100%);
		border-radius: var(--radius-md);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.25rem;
	}

	.brand-text {
		font-size: var(--text-lg);
		font-weight: 700;
		color: var(--color-text);
	}

	.sidebar-nav {
		flex: 1;
		padding: var(--space-3);
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}

	.nav-link {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-3) var(--space-4);
		border-radius: var(--radius-md);
		color: var(--color-text);
		font-size: var(--text-sm);
		font-weight: 500;
		transition: all var(--transition-base);
	}

	.nav-link:hover {
		background: var(--color-bg);
		opacity: 1;
	}

	.nav-link.active {
		background: var(--color-primary-light);
		color: var(--color-primary);
	}

	.nav-icon {
		font-size: 1.125rem;
		width: 24px;
		text-align: center;
	}

	.sidebar-footer {
		padding: var(--space-3);
		border-top: 1px solid var(--color-separator);
	}

	.preview-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-2);
		width: 100%;
		padding: var(--space-3);
		background: var(--color-bg);
		border-radius: var(--radius-md);
		color: var(--color-text);
		font-size: var(--text-sm);
		font-weight: 500;
		transition: all var(--transition-base);
	}

	.preview-btn:hover {
		background: var(--color-border);
		opacity: 1;
	}

	/* Main Area */
	.main-area {
		flex: 1;
		margin-left: 260px;
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	/* Topbar */
	.topbar {
		background: var(--color-bg-secondary);
		border-bottom: 1px solid var(--color-separator);
		padding: var(--space-3) var(--space-5);
		display: flex;
		align-items: center;
		justify-content: space-between;
		position: sticky;
		top: 0;
		z-index: 50;
	}

	.url-badge {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2) var(--space-3);
		background: var(--color-bg);
		border-radius: var(--radius-full);
		font-size: var(--text-xs);
	}

	.url-icon {
		font-size: 0.875rem;
	}

	.url-text {
		font-family: var(--font-mono);
		color: var(--color-text-secondary);
	}

	.user-btn {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2);
		background: transparent;
		border-radius: var(--radius-full);
	}

	.user-btn:hover {
		background: var(--color-bg);
	}

	.user-btn:active {
		transform: none;
	}

	.user-avatar {
		width: 28px;
		height: 28px;
		background: linear-gradient(135deg, #ff9500 0%, #ff3b30 100%);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: var(--text-xs);
		font-weight: 600;
		color: white;
	}

	.user-name {
		font-size: var(--text-sm);
		color: var(--color-text-secondary);
		max-width: 150px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.logout-icon {
		color: var(--color-text-tertiary);
		font-size: 0.75rem;
	}

	/* Content */
	.content-area {
		flex: 1;
		padding: var(--space-5);
	}

	/* Responsive */
	@media (max-width: 768px) {
		.sidebar {
			width: 72px;
		}

		.brand-text,
		.nav-label,
		.preview-btn span:last-child {
			display: none;
		}

		.sidebar-brand {
			justify-content: center;
			padding: var(--space-4);
		}

		.nav-link {
			justify-content: center;
			padding: var(--space-3);
		}

		.preview-btn {
			padding: var(--space-3);
		}

		.main-area {
			margin-left: 72px;
		}

		.url-badge {
			display: none;
		}

		.user-name {
			display: none;
		}
	}
</style>
