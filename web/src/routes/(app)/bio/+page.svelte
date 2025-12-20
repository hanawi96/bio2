<script lang="ts">
	import { onMount } from 'svelte';
	import { bio, type BioData, type BlockWithGroup, type Link } from '$lib/api/client';
	import { getAuth } from '$lib/stores/auth.svelte';
	import { Eye, Link as LinkIcon, FileText, Image, ShoppingBag, Minus, GripVertical, Eye as EyeIcon, EyeOff, Trash2, Pencil, X, Plus } from 'lucide-svelte';

	const auth = getAuth();

	let bioData = $state<BioData | null>(null);
	let loading = $state(true);
	let saving = $state(false);
	let showAddMenu = $state(false);

	// Edit states
	let editingLinkId = $state<number | null>(null);
	let editTitle = $state('');
	let editUrl = $state('');

	// Add link form
	let addingToGroupId = $state<number | null>(null);
	let newTitle = $state('');
	let newUrl = $state('');

	onMount(async () => {
		await loadBio();
	});

	async function loadBio() {
		loading = true;
		try {
			bioData = await bio.get();
		} catch (err) {
			console.error(err);
		} finally {
			loading = false;
		}
	}

	function toggleAddMenu() {
		showAddMenu = !showAddMenu;
	}

	async function addLinkGroup() {
		showAddMenu = false;
		saving = true;
		try {
			await bio.addBlock('link_group');
			await loadBio();
		} catch (err) {
			console.error(err);
		} finally {
			saving = false;
		}
	}

	async function deleteBlock(blockId: number) {
		if (!confirm('Xóa block này?')) return;
		try {
			await bio.deleteBlock(blockId);
			await loadBio();
		} catch (err) {
			console.error(err);
		}
	}

	async function toggleBlockVisibility(block: BlockWithGroup) {
		try {
			await bio.updateBlock(block.id, { is_visible: !block.is_visible });
			await loadBio();
		} catch (err) {
			console.error(err);
		}
	}

	// Link operations
	function startAddLink(groupId: number) {
		addingToGroupId = groupId;
		newTitle = '';
		newUrl = '';
	}

	function cancelAddLink() {
		addingToGroupId = null;
		newTitle = '';
		newUrl = '';
	}

	async function addLink() {
		if (!addingToGroupId || !newTitle.trim() || !newUrl.trim()) return;

		saving = true;
		try {
			await bio.addLink(addingToGroupId, newTitle.trim(), newUrl.trim());
			await loadBio();
			cancelAddLink();
		} catch (err) {
			console.error(err);
		} finally {
			saving = false;
		}
	}

	function startEditLink(link: Link) {
		editingLinkId = link.id;
		editTitle = link.title;
		editUrl = link.url;
	}

	function cancelEditLink() {
		editingLinkId = null;
		editTitle = '';
		editUrl = '';
	}

	async function saveEditLink() {
		if (!editingLinkId || !editTitle.trim() || !editUrl.trim()) return;

		saving = true;
		try {
			await bio.updateLink(editingLinkId, {
				title: editTitle.trim(),
				url: editUrl.trim()
			});
			await loadBio();
			cancelEditLink();
		} catch (err) {
			console.error(err);
		} finally {
			saving = false;
		}
	}

	async function deleteLink(linkId: number) {
		if (!confirm('Xóa link này?')) return;
		try {
			await bio.deleteLink(linkId);
			await loadBio();
		} catch (err) {
			console.error(err);
		}
	}

	async function toggleLinkActive(link: Link) {
		try {
			await bio.updateLink(link.id, { is_active: !link.is_active });
			await loadBio();
		} catch (err) {
			console.error(err);
		}
	}

	// Computed
	const allLinks = $derived(
		bioData?.blocks
			.filter((b) => b.type === 'link_group' && b.group)
			.flatMap((b) => b.group?.links || [])
			.filter((l) => l.is_active) || []
	);

	const blockTypes = [
		{ type: 'link_group', icon: 'link', label: 'Link Group', desc: 'Nhóm các liên kết', available: true },
		{ type: 'text', icon: 'text', label: 'Text', desc: 'Văn bản, tiêu đề', available: false },
		{ type: 'image', icon: 'image', label: 'Image', desc: 'Hình ảnh', available: false },
		{ type: 'product', icon: 'product', label: 'Product', desc: 'Sản phẩm', available: false },
		{ type: 'spacer', icon: 'spacer', label: 'Spacer', desc: 'Khoảng cách', available: false }
	];
</script>

<div class="bio-page">
	<!-- Header -->
	<div class="page-header">
		<div class="header-info">
			<h1>My Bio</h1>
			<p>Quản lý nội dung trang bio của bạn</p>
		</div>
		{#if auth.user?.username}
			<a href="/{auth.user.username}" target="_blank" class="preview-link">
				<Eye size={16} /> Xem trang
			</a>
		{/if}
	</div>

	<div class="bio-layout">
		<!-- Phone Preview -->
		<div class="preview-col">
			<div class="phone">
				<div class="phone-notch"></div>
				<div class="phone-screen">
					{#if allLinks.length === 0 && !loading}
						<div class="screen-empty">
							<p>Chưa có nội dung</p>
						</div>
					{:else}
						<div class="screen-header">
							<div class="avatar"></div>
							<div class="username">@{auth.user?.username}</div>
						</div>
						<div class="screen-links">
							{#each allLinks as link}
								<div class="screen-link">{link.title}</div>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		</div>

		<!-- Editor -->
		<div class="editor-col">
			<!-- Add Block -->
			<div class="add-wrapper">
				<button class="add-btn" onclick={toggleAddMenu}>
					<span class="add-icon">+</span>
					<span>Thêm block</span>
				</button>

				{#if showAddMenu}
					<div class="add-dropdown">
						{#each blockTypes as bt}
							<button
								class="dropdown-item"
								onclick={() => (bt.type === 'link_group' ? addLinkGroup() : null)}
								disabled={!bt.available}
							>
								<span class="item-icon">
									{#if bt.icon === 'link'}<LinkIcon size={20} />
									{:else if bt.icon === 'text'}<FileText size={20} />
									{:else if bt.icon === 'image'}<Image size={20} />
									{:else if bt.icon === 'product'}<ShoppingBag size={20} />
									{:else if bt.icon === 'spacer'}<Minus size={20} />
									{/if}
								</span>
								<div class="item-info">
									<span class="item-label">{bt.label}</span>
									<span class="item-desc">{bt.desc}</span>
								</div>
								{#if !bt.available}
									<span class="soon-tag">Soon</span>
								{/if}
							</button>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Blocks -->
			{#if loading}
				<div class="loading-box">
					<div class="spinner"></div>
				</div>
			{:else if !bioData?.blocks.length}
				<div class="empty-box">
					<div class="empty-icon"><FileText size={40} /></div>
					<h3>Bắt đầu tạo Bio</h3>
					<p>Thêm Link Group để bắt đầu</p>
				</div>
			{:else}
				<div class="blocks">
					{#each bioData.blocks as block (block.id)}
						<div class="block-card" class:dimmed={!block.is_visible}>
							<!-- Block Header -->
							<div class="block-head">
								<div class="block-head-left">
									<button class="drag-btn"><GripVertical size={16} /></button>
									<span class="block-type">
										{#if block.type === 'link_group'}<LinkIcon size={14} /> Link Group{:else}<FileText size={14} /> {block.type}{/if}
										{#if block.group?.links.length}
											<span class="count">({block.group.links.length})</span>
										{/if}
									</span>
								</div>
								<div class="block-head-right">
									<button
										class="icon-btn"
										class:active={block.is_visible}
										onclick={() => toggleBlockVisibility(block)}
									>
										{#if block.is_visible}<EyeIcon size={16} />{:else}<EyeOff size={16} />{/if}
									</button>
									<button class="icon-btn danger" onclick={() => deleteBlock(block.id)}><Trash2 size={16} /></button>
								</div>
							</div>

							<!-- Block Body -->
							{#if block.type === 'link_group' && block.group}
								<div class="block-body">
									<!-- Links -->
									{#if block.group.links.length > 0}
										<div class="link-list">
											{#each block.group.links as link (link.id)}
												<div class="link-item" class:dimmed={!link.is_active}>
													{#if editingLinkId === link.id}
														<div class="link-edit">
															<input type="text" bind:value={editTitle} placeholder="Tiêu đề" />
															<input type="url" bind:value={editUrl} placeholder="URL" />
															<div class="edit-btns">
																<button class="btn-sm secondary" onclick={cancelEditLink}>Hủy</button>
																<button class="btn-sm primary" onclick={saveEditLink} disabled={saving}>Lưu</button>
															</div>
														</div>
													{:else}
														<div class="link-row">
															<div class="link-data">
																<span class="link-title">{link.title}</span>
																<span class="link-url">{link.url}</span>
															</div>
															<div class="link-btns">
																<button
																	class="toggle-pill"
																	class:on={link.is_active}
																	onclick={() => toggleLinkActive(link)}
																>
																	<span class="toggle-dot"></span>
																</button>
																<button class="mini-btn" onclick={() => startEditLink(link)}><Pencil size={12} /></button>
																<button class="mini-btn danger" onclick={() => deleteLink(link.id)}><X size={12} /></button>
															</div>
														</div>
													{/if}
												</div>
											{/each}
										</div>
									{/if}

									<!-- Add Link -->
									{#if addingToGroupId === block.group.id}
										<div class="add-link-box">
											<input type="text" bind:value={newTitle} placeholder="Tiêu đề link" />
											<input type="url" bind:value={newUrl} placeholder="https://..." />
											<div class="add-link-btns">
												<button class="btn-sm secondary" onclick={cancelAddLink}>Hủy</button>
												<button
													class="btn-sm primary"
													onclick={addLink}
													disabled={saving || !newTitle.trim() || !newUrl.trim()}
												>
													Thêm
												</button>
											</div>
										</div>
									{:else}
										<button class="add-link-btn" onclick={() => block.group && startAddLink(block.group.id)}>
											+ Thêm link
										</button>
									{/if}
								</div>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	.bio-page {
		max-width: 1100px;
		margin: 0 auto;
	}

	/* Header */
	.page-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: var(--space-5);
	}

	.page-header h1 {
		font-size: var(--text-xl);
		font-weight: 700;
		color: var(--color-text);
		margin-bottom: var(--space-1);
	}

	.page-header p {
		font-size: var(--text-sm);
		color: var(--color-text-secondary);
	}

	.preview-link {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2) var(--space-4);
		background: var(--color-bg-secondary);
		border-radius: var(--radius-full);
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--color-text);
		box-shadow: var(--shadow-sm);
		transition: all var(--transition-base);
	}

	.preview-link:hover {
		box-shadow: var(--shadow-md);
		opacity: 1;
	}

	/* Layout */
	.bio-layout {
		display: grid;
		grid-template-columns: 280px 1fr;
		gap: var(--space-5);
		align-items: start;
	}

	/* Phone Preview */
	.preview-col {
		position: sticky;
		top: 100px;
	}

	.phone {
		background: #1c1c1e;
		border-radius: 36px;
		padding: 12px;
		box-shadow: var(--shadow-lg);
	}

	.phone-notch {
		width: 90px;
		height: 28px;
		background: #1c1c1e;
		border-radius: 14px;
		margin: 0 auto 8px;
	}

	.phone-screen {
		background: #fff;
		border-radius: 28px;
		min-height: 460px;
		padding: var(--space-5) var(--space-4);
	}

	.screen-empty {
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-text-secondary);
		font-size: var(--text-sm);
	}

	.screen-header {
		text-align: center;
		margin-bottom: var(--space-5);
	}

	.avatar {
		width: 72px;
		height: 72px;
		background: linear-gradient(135deg, #007aff 0%, #5856d6 100%);
		border-radius: 50%;
		margin: 0 auto var(--space-3);
	}

	.username {
		font-size: var(--text-sm);
		font-weight: 600;
		color: var(--color-text);
	}

	.screen-links {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}

	.screen-link {
		background: #1c1c1e;
		color: #fff;
		padding: var(--space-3);
		border-radius: var(--radius-md);
		text-align: center;
		font-size: var(--text-sm);
		font-weight: 500;
	}

	/* Editor */
	.editor-col {
		min-height: 400px;
	}

	/* Add Block */
	.add-wrapper {
		position: relative;
		margin-bottom: var(--space-4);
	}

	.add-btn {
		width: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-2);
		padding: var(--space-4);
		background: var(--color-bg-secondary);
		border: 2px dashed var(--color-border);
		border-radius: var(--radius-lg);
		color: var(--color-text-secondary);
		font-size: var(--text-sm);
		font-weight: 500;
		transition: all var(--transition-base);
	}

	.add-btn:hover {
		border-color: var(--color-primary);
		color: var(--color-primary);
	}

	.add-btn:active {
		transform: none;
	}

	.add-icon {
		font-size: 1.25rem;
		font-weight: 300;
	}

	.add-dropdown {
		position: absolute;
		top: 100%;
		left: 0;
		right: 0;
		background: var(--color-bg-secondary);
		border-radius: var(--radius-lg);
		box-shadow: var(--shadow-lg);
		z-index: 100;
		margin-top: var(--space-2);
		overflow: hidden;
	}

	.dropdown-item {
		width: 100%;
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-3) var(--space-4);
		background: none;
		text-align: left;
		border-bottom: 1px solid var(--color-separator);
		border-radius: 0;
	}

	.dropdown-item:last-child {
		border-bottom: none;
	}

	.dropdown-item:hover:not(:disabled) {
		background: var(--color-bg);
	}

	.dropdown-item:disabled {
		opacity: 0.5;
	}

	.dropdown-item:active {
		transform: none;
	}

	.item-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-text-secondary);
	}

	.item-info {
		flex: 1;
	}

	.item-label {
		display: block;
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--color-text);
	}

	.item-desc {
		font-size: var(--text-xs);
		color: var(--color-text-secondary);
	}

	.soon-tag {
		font-size: 0.625rem;
		background: var(--color-primary-light);
		color: var(--color-primary);
		padding: 2px 6px;
		border-radius: var(--radius-full);
		font-weight: 600;
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
		text-align: center;
		padding: var(--space-8);
		background: var(--color-bg-secondary);
		border-radius: var(--radius-lg);
		border: 2px dashed var(--color-border);
	}

	.empty-icon {
		color: var(--color-text-tertiary);
		margin-bottom: var(--space-3);
		display: flex;
		justify-content: center;
	}

	.empty-box h3 {
		font-size: var(--text-base);
		font-weight: 600;
		color: var(--color-text);
		margin-bottom: var(--space-1);
	}

	.empty-box p {
		font-size: var(--text-sm);
		color: var(--color-text-secondary);
	}

	/* Blocks */
	.blocks {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}

	.block-card {
		background: var(--color-bg-secondary);
		border-radius: var(--radius-lg);
		overflow: hidden;
		box-shadow: var(--shadow-card);
		transition: all var(--transition-base);
	}

	.block-card.dimmed {
		opacity: 0.6;
	}

	.block-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: var(--space-3) var(--space-4);
		background: var(--color-bg);
		border-bottom: 1px solid var(--color-separator);
	}

	.block-head-left {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}

	.drag-btn {
		background: none;
		color: var(--color-text-tertiary);
		padding: var(--space-1);
		cursor: grab;
		display: flex;
		align-items: center;
	}

	.drag-btn:hover {
		color: var(--color-text-secondary);
	}

	.drag-btn:active {
		transform: none;
	}

	.block-type {
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--color-text);
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}

	.count {
		color: var(--color-text-secondary);
		font-weight: 400;
	}

	.block-head-right {
		display: flex;
		gap: var(--space-1);
	}

	.icon-btn {
		width: 32px;
		height: 32px;
		background: transparent;
		border-radius: var(--radius-sm);
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0;
		color: var(--color-text-secondary);
	}

	.icon-btn:hover {
		background: var(--color-bg-secondary);
	}

	.icon-btn.danger:hover {
		background: rgba(255, 59, 48, 0.1);
	}

	.icon-btn:active {
		transform: none;
	}

	.block-body {
		padding: var(--space-4);
	}

	/* Link List */
	.link-list {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
		margin-bottom: var(--space-3);
	}

	.link-item {
		background: var(--color-bg);
		border-radius: var(--radius-md);
		padding: var(--space-3);
		transition: all var(--transition-fast);
	}

	.link-item.dimmed {
		opacity: 0.5;
	}

	.link-row {
		display: flex;
		align-items: center;
		gap: var(--space-3);
	}

	.link-data {
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

	.link-btns {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}

	/* Toggle Pill */
	.toggle-pill {
		width: 44px;
		height: 26px;
		background: var(--color-text-tertiary);
		border-radius: var(--radius-full);
		position: relative;
		padding: 0;
		transition: background var(--transition-base);
	}

	.toggle-pill.on {
		background: var(--color-success);
	}

	.toggle-pill:active {
		transform: none;
	}

	.toggle-dot {
		position: absolute;
		top: 2px;
		left: 2px;
		width: 22px;
		height: 22px;
		background: #fff;
		border-radius: 50%;
		box-shadow: var(--shadow-sm);
		transition: transform var(--transition-spring);
	}

	.toggle-pill.on .toggle-dot {
		transform: translateX(18px);
	}

	.mini-btn {
		width: 28px;
		height: 28px;
		background: transparent;
		border-radius: var(--radius-sm);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-text-secondary);
		padding: 0;
	}

	.mini-btn:hover {
		background: var(--color-bg-secondary);
	}

	.mini-btn.danger:hover {
		background: rgba(255, 59, 48, 0.1);
		color: var(--color-danger);
	}

	.mini-btn:active {
		transform: none;
	}

	/* Link Edit */
	.link-edit {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}

	.link-edit input {
		padding: var(--space-2) var(--space-3);
		font-size: var(--text-sm);
	}

	.edit-btns {
		display: flex;
		gap: var(--space-2);
		justify-content: flex-end;
	}

	/* Add Link */
	.add-link-btn {
		width: 100%;
		padding: var(--space-3);
		background: none;
		border: 1px dashed var(--color-border);
		border-radius: var(--radius-md);
		color: var(--color-text-secondary);
		font-size: var(--text-sm);
	}

	.add-link-btn:hover {
		border-color: var(--color-primary);
		color: var(--color-primary);
	}

	.add-link-btn:active {
		transform: none;
	}

	.add-link-box {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
		padding: var(--space-3);
		background: var(--color-bg);
		border-radius: var(--radius-md);
	}

	.add-link-box input {
		padding: var(--space-2) var(--space-3);
		font-size: var(--text-sm);
	}

	.add-link-btns {
		display: flex;
		gap: var(--space-2);
		justify-content: flex-end;
	}

	/* Buttons */
	.btn-sm {
		padding: var(--space-2) var(--space-3);
		font-size: var(--text-xs);
		border-radius: var(--radius-sm);
	}

	.btn-sm.primary {
		background: var(--color-primary);
		color: #fff;
	}

	.btn-sm.primary:hover:not(:disabled) {
		background: #0066d6;
	}

	.btn-sm.secondary {
		background: var(--color-bg-secondary);
		color: var(--color-text);
	}

	.btn-sm.secondary:hover {
		background: var(--color-border);
	}

	/* Responsive */
	@media (max-width: 900px) {
		.bio-layout {
			grid-template-columns: 1fr;
		}

		.preview-col {
			display: none;
		}
	}
</style>
