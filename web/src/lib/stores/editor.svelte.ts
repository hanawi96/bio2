import { pages, type DraftData, type Block, type LinkGroup, type Link, type SaveRequest } from '$lib/api/client';
import { generateSortKey } from '$lib/utils/sortkey';

let draft = $state<DraftData | null>(null);
let saving = $state(false);
let dirty = $state(false);

// Pending changes
let pendingBlocks = $state<Map<number, Block>>(new Map());
let pendingGroups = $state<Map<number, LinkGroup>>(new Map());
let pendingLinks = $state<Map<number, Link>>(new Map());
let deletedBlocks = $state<Set<number>>(new Set());
let deletedGroups = $state<Set<number>>(new Set());
let deletedLinks = $state<Set<number>>(new Set());

let tempIdCounter = -1;

export function getEditor() {
	return {
		get draft() { return draft; },
		get saving() { return saving; },
		get dirty() { return dirty; },
		get blocks() { 
			if (!draft) return [];
			return draft.blocks.map(b => pendingBlocks.get(b.id) || b)
				.filter(b => !deletedBlocks.has(b.id))
				.sort((a, b) => a.sort_key.localeCompare(b.sort_key));
		},
		get linkGroups() {
			if (!draft) return [];
			return draft.link_groups.map(g => pendingGroups.get(g.id) || g)
				.filter(g => !deletedGroups.has(g.id));
		},
		getLinks(groupId: number) {
			if (!draft) return [];
			const links = draft.links[groupId] || [];
			return links.map(l => pendingLinks.get(l.id) || l)
				.filter(l => !deletedLinks.has(l.id))
				.sort((a, b) => a.sort_key.localeCompare(b.sort_key));
		}
	};
}

export async function loadDraft(pageId: number) {
	draft = await pages.getDraft(pageId);
	resetPending();
}

function resetPending() {
	pendingBlocks = new Map();
	pendingGroups = new Map();
	pendingLinks = new Map();
	deletedBlocks = new Set();
	deletedGroups = new Set();
	deletedLinks = new Set();
	dirty = false;
}

// Block operations
export function addBlock(type: string, afterBlockId?: number) {
	if (!draft) return;

	const blocks = getEditor().blocks;
	let sortKey: string;

	if (afterBlockId) {
		const idx = blocks.findIndex(b => b.id === afterBlockId);
		const prev = blocks[idx]?.sort_key || '';
		const next = blocks[idx + 1]?.sort_key || '';
		sortKey = generateSortKey(prev, next);
	} else {
		const last = blocks[blocks.length - 1]?.sort_key || '';
		sortKey = generateSortKey(last, '');
	}

	const newBlock: Block = {
		id: tempIdCounter--,
		page_id: draft.page.id,
		type,
		sort_key: sortKey,
		content: {},
		is_visible: true
	};

	draft.blocks = [...draft.blocks, newBlock];
	dirty = true;
}

export function updateBlock(id: number, updates: Partial<Block>) {
	if (!draft) return;
	const block = draft.blocks.find(b => b.id === id);
	if (block) {
		const updated = { ...block, ...updates };
		pendingBlocks.set(id, updated);
		dirty = true;
	}
}

export function deleteBlock(id: number) {
	deletedBlocks.add(id);
	dirty = true;
}

export function reorderBlock(id: number, newIndex: number) {
	if (!draft) return;
	const blocks = getEditor().blocks;
	const prev = blocks[newIndex - 1]?.sort_key || '';
	const next = blocks[newIndex + 1]?.sort_key || '';
	const sortKey = generateSortKey(prev, next);
	updateBlock(id, { sort_key: sortKey });
}

// Link Group operations
export function addLinkGroup(title?: string) {
	if (!draft) return null;

	const newGroup: LinkGroup = {
		id: tempIdCounter--,
		page_id: draft.page.id,
		title,
		layout_type: 'list',
		layout_config: {}
	};

	draft.link_groups = [...draft.link_groups, newGroup];
	draft.links[newGroup.id] = [];
	dirty = true;

	return newGroup;
}

export function updateLinkGroup(id: number, updates: Partial<LinkGroup>) {
	if (!draft) return;
	const group = draft.link_groups.find(g => g.id === id);
	if (group) {
		const updated = { ...group, ...updates };
		pendingGroups.set(id, updated);
		dirty = true;
	}
}

export function deleteLinkGroup(id: number) {
	deletedGroups.add(id);
	dirty = true;
}

// Link operations
export function addLink(groupId: number, title: string, url: string) {
	if (!draft) return;

	const links = getEditor().getLinks(groupId);
	const last = links[links.length - 1]?.sort_key || '';
	const sortKey = generateSortKey(last, '');

	const newLink: Link = {
		id: tempIdCounter--,
		group_id: groupId,
		title,
		url,
		sort_key: sortKey,
		is_active: true,
		created_at: new Date().toISOString(),
		updated_at: new Date().toISOString()
	};

	if (!draft.links[groupId]) {
		draft.links[groupId] = [];
	}
	draft.links[groupId] = [...draft.links[groupId], newLink];
	dirty = true;
}

export function updateLink(id: number, updates: Partial<Link>) {
	if (!draft) return;
	for (const groupLinks of Object.values(draft.links)) {
		const link = groupLinks.find(l => l.id === id);
		if (link) {
			const updated = { ...link, ...updates };
			pendingLinks.set(id, updated);
			dirty = true;
			break;
		}
	}
}

export function deleteLink(id: number) {
	deletedLinks.add(id);
	dirty = true;
}

// Save
export async function save() {
	if (!draft || !dirty) return;

	saving = true;
	try {
		const req: SaveRequest = {
			blocks: [],
			link_groups: [],
			links: []
		};

		// Collect block changes
		for (const block of draft.blocks) {
			if (deletedBlocks.has(block.id)) {
				if (block.id > 0) {
					req.blocks!.push({ ...block, delete: true });
				}
			} else if (pendingBlocks.has(block.id) || block.id < 0) {
				req.blocks!.push(pendingBlocks.get(block.id) || block);
			}
		}

		// Collect group changes
		for (const group of draft.link_groups) {
			if (deletedGroups.has(group.id)) {
				if (group.id > 0) {
					req.link_groups!.push({ ...group, delete: true });
				}
			} else if (pendingGroups.has(group.id) || group.id < 0) {
				req.link_groups!.push(pendingGroups.get(group.id) || group);
			}
		}

		// Collect link changes
		for (const groupLinks of Object.values(draft.links)) {
			for (const link of groupLinks) {
				if (deletedLinks.has(link.id)) {
					if (link.id > 0) {
						req.links!.push({ ...link, delete: true });
					}
				} else if (pendingLinks.has(link.id) || link.id < 0) {
					req.links!.push(pendingLinks.get(link.id) || link);
				}
			}
		}

		await pages.save(draft.page.id, req);

		// Reload draft to get real IDs
		await loadDraft(draft.page.id);
	} finally {
		saving = false;
	}
}

export async function publish() {
	if (!draft) return;
	await pages.publish(draft.page.id);
}
