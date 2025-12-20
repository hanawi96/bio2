<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { checkAuth } from '$lib/stores/auth.svelte';

	let { children } = $props();

	onMount(() => {
		checkAuth();
	});
</script>

{@render children()}
