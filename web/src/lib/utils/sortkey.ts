const chars = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';

export function generateSortKey(prev: string, next: string): string {
	if (!prev && !next) return 'U';

	if (!prev) {
		// Insert before next
		return decrementKey(next);
	}

	if (!next) {
		// Insert after prev
		return incrementKey(prev);
	}

	// Insert between
	return midKey(prev, next);
}

function incrementKey(key: string): string {
	if (!key) return 'U';

	const arr = key.split('');
	for (let i = arr.length - 1; i >= 0; i--) {
		const idx = chars.indexOf(arr[i]);
		if (idx < chars.length - 1) {
			arr[i] = chars[idx + 1];
			return arr.join('');
		}
		arr[i] = '0';
	}
	return arr.join('') + '0';
}

function decrementKey(key: string): string {
	if (!key) return 'U';

	const arr = key.split('');
	for (let i = arr.length - 1; i >= 0; i--) {
		const idx = chars.indexOf(arr[i]);
		if (idx > 0) {
			arr[i] = chars[idx - 1];
			return arr.join('');
		}
	}
	return '0' + key;
}

function midKey(prev: string, next: string): string {
	const maxLen = Math.max(prev.length, next.length);
	prev = prev.padEnd(maxLen, '0');
	next = next.padEnd(maxLen, '0');

	const result: string[] = [];
	for (let i = 0; i < maxLen; i++) {
		const p = chars.indexOf(prev[i]);
		const n = chars.indexOf(next[i]);
		const mid = Math.floor((p + n) / 2);
		result.push(chars[mid]);
	}

	const res = result.join('');
	if (res === prev) {
		return prev + 'U';
	}
	return res;
}
