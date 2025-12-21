const API_URL = 'http://localhost:8080';

interface ApiResponse<T> {
	success: boolean;
	data?: T;
	error?: string;
}

async function request<T>(
	endpoint: string,
	options: RequestInit = {}
): Promise<T> {
	const res = await fetch(`${API_URL}${endpoint}`, {
		...options,
		credentials: 'include',
		headers: {
			'Content-Type': 'application/json',
			...options.headers
		}
	});

	const json: ApiResponse<T> = await res.json();

	if (!json.success) {
		throw new Error(json.error || 'Request failed');
	}

	return json.data as T;
}

// Auth
export const auth = {
	register: (email: string, password: string) =>
		request<{ user: User; token: string }>('/api/auth/register', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		}),

	login: (email: string, password: string) =>
		request<{ user: User; token: string }>('/api/auth/login', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		}),

	logout: () =>
		request('/api/auth/logout', { method: 'POST' }),

	me: () =>
		request<User>('/api/auth/me'),

	setUsername: (username: string) =>
		request<User>('/api/auth/username', {
			method: 'POST',
			body: JSON.stringify({ username })
		}),

	checkUsername: (username: string) =>
		request<{ available: boolean }>(`/api/auth/check-username?username=${encodeURIComponent(username)}`)
};

// Links (deprecated - use bio API)
export const links = {
	list: () => bio.get().then(d => d.blocks.flatMap(b => b.group?.links || [])),
	create: (title: string, url: string) => {
		// This is a simplified version - real implementation should use bio.addLink
		return request<Link>('/api/bio/links', {
			method: 'POST',
			body: JSON.stringify({ group_id: 0, title, url })
		});
	},
	update: (id: number, data: { title?: string; url?: string; is_active?: boolean }) =>
		bio.updateLink(id, data),
	delete: (id: number) => bio.deleteLink(id),
	reorder: (linkIds: number[]) =>
		request('/api/bio/blocks/reorder', {
			method: 'POST',
			body: JSON.stringify({ block_ids: linkIds })
		})
};

// Bio API (blocks + groups + links)
export const bio = {
	get: () =>
		request<BioData>('/api/bio'),

	addBlock: (type: string, content?: object) =>
		request<BlockWithGroup>('/api/bio/blocks', {
			method: 'POST',
			body: JSON.stringify({ type, content })
		}),

	updateBlock: (id: number, data: { content?: object; is_visible?: boolean }) =>
		request<Block>(`/api/bio/blocks/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		}),

	deleteBlock: (id: number) =>
		request(`/api/bio/blocks/${id}`, { method: 'DELETE' }),

	reorderBlocks: (blockIds: number[]) =>
		request('/api/bio/blocks/reorder', {
			method: 'POST',
			body: JSON.stringify({ block_ids: blockIds })
		}),

	addLink: (groupId: number, title: string, url: string) =>
		request<Link>('/api/bio/links', {
			method: 'POST',
			body: JSON.stringify({ group_id: groupId, title, url })
		}),

	updateLink: (id: number, data: { title?: string; url?: string; is_active?: boolean }) =>
		request<Link>(`/api/bio/links/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		}),

	deleteLink: (id: number) =>
		request(`/api/bio/links/${id}`, { method: 'DELETE' }),

	updateProfile: (displayName: string, bio: string) =>
		request('/api/bio/profile', {
			method: 'PUT',
			body: JSON.stringify({ display_name: displayName, bio })
		}),

	updateSocial: (social: SocialLinks) =>
		request('/api/bio/social', {
			method: 'PUT',
			body: JSON.stringify(social)
		})
};

// Pages
export const pages = {
	list: () =>
		request<Page[]>('/api/pages'),

	get: (id: number) =>
		request<Page>(`/api/pages/${id}`),

	create: (title: string, themePresetId = 1) =>
		request<Page>('/api/pages', {
			method: 'POST',
			body: JSON.stringify({ title, theme_preset_id: themePresetId })
		}),

	getDraft: (id: number) =>
		request<DraftData>(`/api/pages/${id}/draft`),

	save: (id: number, data: SaveRequest) =>
		request(`/api/pages/${id}/save`, {
			method: 'POST',
			body: JSON.stringify(data)
		}),

	publish: (id: number) =>
		request(`/api/pages/${id}/publish`, { method: 'POST' }),

	delete: (id: number) =>
		request(`/api/pages/${id}`, { method: 'DELETE' })
};

// Themes
export const themes = {
	listPresets: (tier?: string) =>
		request<ThemePreset[]>(`/api/themes/presets${tier ? `?tier=${tier}` : ''}`),

	getUserCustom: () =>
		request<ThemeCustom | null>('/api/themes/custom'),

	createCustom: (presetId: number, patch: object) =>
		request<ThemeCustom>('/api/themes/custom', {
			method: 'POST',
			body: JSON.stringify({ preset_id: presetId, patch })
		}),

	deleteCustom: (id: number) =>
		request<{ deleted: boolean }>(`/api/themes/custom/${id}`, {
			method: 'DELETE'
		})
};

// Types
export interface User {
	id: number;
	email: string;
	username?: string;
	display_name?: string;
	avatar_asset_id?: number;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface Link {
	id: number;
	group_id: number;
	title: string;
	url: string;
	sort_key: string;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface Page {
	id: number;
	user_id: number;
	locale: string;
	title?: string;
	status: string;
	access_type: string;
	theme_preset_id: number;
	theme_custom_id?: number;
	theme_mode: string;
	settings: object;
	created_at: string;
	updated_at: string;
}

export interface Block {
	id: number;
	page_id: number;
	type: string;
	sort_key: string;
	ref_id?: number;
	content: object;
	is_visible: boolean;
}

export interface LinkGroup {
	id: number;
	page_id: number;
	title?: string;
	layout_type: string;
	layout_config: object;
	style_override?: object;
}

export interface GroupWithLinks extends LinkGroup {
	links: Link[];
}

export interface BlockWithGroup extends Block {
	group?: GroupWithLinks;
}

export interface BioData {
	page: Page;
	blocks: BlockWithGroup[];
}

export interface SocialLinks {
	instagram?: string;
	facebook?: string;
	twitter?: string;
	tiktok?: string;
	youtube?: string;
	linkedin?: string;
	github?: string;
	website?: string;
}

export interface DraftData {
	page: Page;
	blocks: Block[];
	link_groups: LinkGroup[];
	links: Record<number, Link[]>;
}

export interface SaveRequest {
	page?: Partial<Page>;
	blocks?: SaveBlock[];
	link_groups?: SaveLinkGroup[];
	links?: SaveLink[];
}

export interface SaveBlock {
	id?: number;
	type: string;
	sort_key: string;
	ref_id?: number;
	content: object;
	is_visible: boolean;
	delete?: boolean;
}

export interface SaveLinkGroup {
	id?: number;
	title?: string;
	layout_type: string;
	layout_config?: object;
	style_override?: object;
	delete?: boolean;
}

export interface SaveLink {
	id?: number;
	group_id: number;
	title: string;
	url: string;
	sort_key: string;
	is_active: boolean;
	delete?: boolean;
}

export interface ThemePreset {
	id: number;
	key: string;
	name: string;
	tier: string;
	config: object;
}

export interface ThemeCustom {
	id: number;
	based_on_preset_id: number;
	patch: object;
	compiled_config?: object;
}
