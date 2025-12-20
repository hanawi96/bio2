<script lang="ts">
	import { getAuth } from '$lib/stores/auth.svelte';

	const auth = getAuth();
</script>

<main>
	<div class="hero">
		<h1>LinkBio</h1>
		<p>Create your link in bio page</p>

		{#if auth.loading}
			<p>Loading...</p>
		{:else if auth.isLoggedIn}
			<a href="/dashboard" class="btn-primary">Go to Dashboard</a>
		{:else}
			<div class="actions">
				<a href="/login" class="btn-primary">Login</a>
				<a href="/register" class="btn-secondary">Register</a>
			</div>
		{/if}
	</div>
</main>

<style>
	main {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.hero {
		text-align: center;
	}

	h1 {
		font-size: 3rem;
		margin-bottom: 0.5rem;
	}

	p {
		color: #666;
		margin-bottom: 2rem;
	}

	.actions {
		display: flex;
		gap: 1rem;
		justify-content: center;
	}

	.btn-primary, .btn-secondary {
		display: inline-block;
		padding: 0.75rem 1.5rem;
		border-radius: 8px;
		font-weight: 500;
	}

	.btn-primary {
		background: var(--color-primary);
		color: white;
	}

	.btn-secondary {
		background: var(--color-border);
		color: var(--color-text);
	}
</style>
