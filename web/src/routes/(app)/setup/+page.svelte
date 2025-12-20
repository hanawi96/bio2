<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/api/client';
	import { getAuth, setUsername } from '$lib/stores/auth.svelte';
	import { Check, X, Loader2 } from 'lucide-svelte';

	const authState = getAuth();

	let username = $state('');
	let error = $state('');
	let checking = $state(false);
	let available = $state<boolean | null>(null);
	let saving = $state(false);

	let debounceTimer: ReturnType<typeof setTimeout>;

	$effect(() => {
		if (!authState.loading && !authState.needsSetup && authState.isLoggedIn) {
			goto('/dashboard');
		}
	});

	function handleInput() {
		error = '';
		available = null;

		clearTimeout(debounceTimer);

		const value = username.toLowerCase().replace(/[^a-z0-9_]/g, '');
		if (value !== username) {
			username = value;
		}

		if (username.length >= 3) {
			checking = true;
			debounceTimer = setTimeout(async () => {
				try {
					const result = await auth.checkUsername(username);
					available = result.available;
				} catch {
					available = null;
				} finally {
					checking = false;
				}
			}, 300);
		}
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		
		if (username.length < 3) {
			error = 'Username phải có ít nhất 3 ký tự';
			return;
		}

		if (available === false) {
			error = 'Username này đã được sử dụng';
			return;
		}

		saving = true;
		error = '';

		try {
			await setUsername(username);
			goto('/dashboard');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Có lỗi xảy ra';
		} finally {
			saving = false;
		}
	}
</script>

<main>
	<div class="setup-card">
		<div class="header">
			<h1>Chọn username của bạn</h1>
			<p>Đây sẽ là địa chỉ trang bio của bạn</p>
		</div>

		<form onsubmit={handleSubmit}>
			<div class="input-wrapper">
				<span class="prefix">linkbio.com/</span>
				<input
					type="text"
					bind:value={username}
					oninput={handleInput}
					placeholder="username"
					maxlength="30"
					autocomplete="off"
					autocapitalize="off"
				/>
				{#if checking}
					<span class="status checking"><Loader2 size={16} class="spin" /></span>
				{:else if available === true}
					<span class="status available"><Check size={16} /></span>
				{:else if available === false}
					<span class="status taken"><X size={16} /></span>
				{/if}
			</div>

			<p class="hint">
				Chỉ dùng chữ thường, số và dấu gạch dưới. Tối thiểu 3 ký tự.
			</p>

			{#if error}
				<p class="error">{error}</p>
			{/if}

			{#if available === false}
				<p class="error">Username này đã được sử dụng</p>
			{/if}

			<button 
				type="submit" 
				class="btn-submit"
				disabled={saving || username.length < 3 || available === false}
			>
				{saving ? 'Đang lưu...' : 'Tiếp tục'}
			</button>
		</form>
	</div>
</main>

<style>
	main {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	}

	.setup-card {
		background: white;
		border-radius: 16px;
		padding: 2.5rem;
		width: 100%;
		max-width: 440px;
		box-shadow: 0 20px 60px rgba(0,0,0,0.15);
	}

	.header {
		text-align: center;
		margin-bottom: 2rem;
	}

	h1 {
		font-size: 1.5rem;
		font-weight: 700;
		color: #1a1a1a;
		margin-bottom: 0.5rem;
	}

	.header p {
		color: #666;
		font-size: 0.9rem;
	}

	.input-wrapper {
		display: flex;
		align-items: center;
		border: 2px solid #e5e7eb;
		border-radius: 12px;
		padding: 0 1rem;
		transition: border-color 0.2s;
		background: #f9fafb;
	}

	.input-wrapper:focus-within {
		border-color: #667eea;
		background: white;
	}

	.prefix {
		color: #666;
		font-size: 0.95rem;
		white-space: nowrap;
	}

	.input-wrapper input {
		flex: 1;
		border: none;
		background: transparent;
		padding: 1rem 0.5rem;
		font-size: 1rem;
		font-weight: 500;
		min-width: 0;
	}

	.input-wrapper input:focus {
		outline: none;
	}

	.status {
		width: 24px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.status.checking {
		color: #666;
	}

	:global(.spin) {
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.status.available {
		color: #22c55e;
	}

	.status.taken {
		color: #ef4444;
	}

	.hint {
		font-size: 0.8rem;
		color: #999;
		margin-top: 0.75rem;
		text-align: center;
	}

	.error {
		color: #ef4444;
		font-size: 0.875rem;
		margin-top: 0.75rem;
		text-align: center;
	}

	.btn-submit {
		width: 100%;
		padding: 1rem;
		margin-top: 1.5rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
		font-size: 1rem;
		font-weight: 600;
		border-radius: 12px;
		transition: opacity 0.2s, transform 0.2s;
	}

	.btn-submit:hover:not(:disabled) {
		opacity: 0.9;
		transform: translateY(-1px);
	}

	.btn-submit:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
