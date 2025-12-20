<script lang="ts">
	import { goto } from '$app/navigation';
	import { register } from '$lib/stores/auth.svelte';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;

		try {
			await register(email, password);
			goto('/dashboard');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Registration failed';
		} finally {
			loading = false;
		}
	}
</script>

<main>
	<form onsubmit={handleSubmit}>
		<h1>Register</h1>

		{#if error}
			<p class="error">{error}</p>
		{/if}

		<div class="field">
			<label for="email">Email</label>
			<input type="email" id="email" bind:value={email} required />
		</div>

		<div class="field">
			<label for="password">Password</label>
			<input type="password" id="password" bind:value={password} required minlength="6" />
		</div>

		<button type="submit" class="btn-primary" disabled={loading}>
			{loading ? 'Loading...' : 'Register'}
		</button>

		<p class="link">
			Already have an account? <a href="/login">Login</a>
		</p>
	</form>
</main>

<style>
	main {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
	}

	form {
		width: 100%;
		max-width: 400px;
	}

	h1 {
		margin-bottom: 1.5rem;
	}

	.field {
		margin-bottom: 1rem;
	}

	label {
		display: block;
		margin-bottom: 0.25rem;
		font-size: 0.875rem;
		font-weight: 500;
	}

	button {
		width: 100%;
		padding: 0.75rem;
		margin-top: 0.5rem;
	}

	.error {
		color: #dc2626;
		margin-bottom: 1rem;
		padding: 0.5rem;
		background: #fef2f2;
		border-radius: 4px;
	}

	.link {
		margin-top: 1rem;
		text-align: center;
		font-size: 0.875rem;
	}
</style>
