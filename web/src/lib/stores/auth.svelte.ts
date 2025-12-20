import { auth, type User } from '$lib/api/client';

let user = $state<User | null>(null);
let loading = $state(true);

export function getAuth() {
	return {
		get user() { return user; },
		get loading() { return loading; },
		get isLoggedIn() { return !!user; },
		get needsSetup() { return !!user && !user.username; }
	};
}

export async function checkAuth() {
	loading = true;
	try {
		user = await auth.me();
	} catch {
		user = null;
	} finally {
		loading = false;
	}
}

export async function login(email: string, password: string) {
	const result = await auth.login(email, password);
	user = result.user;
	return result;
}

export async function register(email: string, password: string) {
	const result = await auth.register(email, password);
	user = result.user;
	return result;
}

export async function setUsername(username: string) {
	user = await auth.setUsername(username);
	return user;
}

export async function logout() {
	await auth.logout();
	user = null;
}
