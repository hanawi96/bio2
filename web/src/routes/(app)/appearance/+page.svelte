<script lang="ts">
	import { onMount } from 'svelte';
	import { getAuth } from '$lib/stores/auth.svelte';
	import { getAppearance, loadAppearance, selectPreset, updateSetting, saveAppearance, resetAppearance } from '$lib/stores/appearance.svelte';
	import { bio } from '$lib/api/client';

	const auth = getAuth();
	const appearance = getAppearance();
	let pageId = $state<number | null>(null);
	let activeSection = $state<string>('theme');

	onMount(async () => {
		const bioData = await bio.get();
		if (bioData?.page) {
			pageId = bioData.page.id;
			await loadAppearance(bioData.page.id);
		}
	});

	async function handleSave() { await saveAppearance(); }

	const sections = [
		{ id: 'theme', label: 'Theme', icon: '🎨' },
		{ id: 'background', label: 'Background', icon: '🖼️' },
		{ id: 'header', label: 'Header', icon: '👤' },
		{ id: 'links', label: 'Links', icon: '🔗' },
		{ id: 'colors', label: 'Colors', icon: '🌈' },
		{ id: 'typography', label: 'Typography', icon: '✏️' }
	];

	const fontSizes = [{ value: 'S', label: 'S' }, { value: 'M', label: 'M' }, { value: 'L', label: 'L' }, { value: 'XL', label: 'XL' }];
	const alignOptions = [{ value: 'left', icon: '⬅️' }, { value: 'center', icon: '↔️' }, { value: 'right', icon: '➡️' }];
	const linkStyles = [{ value: 'filled', label: 'Filled' }, { value: 'outline', label: 'Outline' }, { value: 'soft', label: 'Soft' }, { value: 'ghost', label: 'Ghost' }];
	const shadowOptions = [{ value: 'none', label: 'None' }, { value: 'sm', label: 'S' }, { value: 'md', label: 'M' }, { value: 'lg', label: 'L' }];
	const headerStyles = [{ value: 'default', label: 'Default' }, { value: 'minimal', label: 'Minimal' }, { value: 'centered', label: 'Centered' }, { value: 'card', label: 'Card' }];
	const avatarSizes = [{ value: 'S', label: 'S' }, { value: 'M', label: 'M' }, { value: 'L', label: 'L' }];
	const avatarShapes = [{ value: 'circle', label: 'Circle' }, { value: 'rounded', label: 'Rounded' }, { value: 'square', label: 'Square' }];
	const backgroundTypes = [{ value: 'solid', label: 'Solid' }, { value: 'gradient', label: 'Gradient' }, { value: 'image', label: 'Image' }];
	const presetColors = ['#007aff', '#5856d6', '#af52de', '#ff2d55', '#ff3b30', '#ff9500', '#ffcc00', '#34c759', '#00c7be', '#30b0c7', '#1c1c1e', '#8e8e93'];
	const gradientPresets = [
		'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
		'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
		'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
		'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
		'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
		'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)'
	];

	function getPresetColor(preset: any): string {
		const colors: Record<string, string> = { 'Light': '#f2f2f7', 'Dark': '#1c1c1e', 'Ocean': 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', 'Sunset': 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)', 'Forest': 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' };
		return preset.config?.background?.color || colors[preset.name] || '#f2f2f7';
	}
	function getBgStyle(): string { const bg = appearance.settings.background; return bg.type === 'gradient' && bg.gradient ? bg.gradient : bg.type === 'image' && bg.imageUrl ? `url(${bg.imageUrl}) center/cover` : bg.color; }
	function getAvatarSize(): number { return { S: 48, M: 64, L: 80 }[appearance.settings.header.avatarSize] || 64; }
	function getAvatarRadius(): string { return { circle: '50%', rounded: '16px', square: '8px' }[appearance.settings.header.avatarShape] || '50%'; }
	function getShadow(): string { return { none: 'none', sm: '0 1px 2px rgba(0,0,0,0.08)', md: '0 2px 8px rgba(0,0,0,0.12)', lg: '0 4px 16px rgba(0,0,0,0.16)' }[appearance.settings.links.shadow] || 'none'; }
	function getLinkBg(): string { const s = appearance.settings.links.style; return s === 'filled' ? appearance.settings.colors.primary : s === 'soft' ? `${appearance.settings.colors.primary}15` : appearance.settings.colors.cardBackground; }
	function getLinkColor(): string { return appearance.settings.links.style === 'filled' ? '#fff' : appearance.settings.colors.text; }
	function getLinkBorder(): string { return appearance.settings.links.style === 'outline' ? `1px solid ${appearance.settings.colors.primary}` : 'none'; }
</script>

{#if appearance.loading}
	<div class="loading-container"><div class="loading-spinner"></div><p>Loading...</p></div>
{:else}
	<div class="page-container">
		<header class="page-header">
			<div><h1>Giao diện</h1><p class="text-secondary text-sm">Tùy chỉnh giao diện trang bio</p></div>
			<div class="header-actions">
				{#if appearance.dirty}<button class="btn-secondary" onclick={resetAppearance} disabled={appearance.saving}>Hủy</button>{/if}
				<button class="btn-primary btn-save" onclick={handleSave} disabled={!appearance.dirty || appearance.saving}>
					{#if appearance.saving}<span class="spinner"></span>{:else}💾{/if}
					<span>{appearance.saving ? 'Đang lưu...' : 'Lưu tất cả'}</span>
				</button>
			</div>
		</header>

		{#if appearance.dirty}<div class="alert-banner warning">⚠️ Bạn có thay đổi chưa lưu</div>{/if}

		<div class="layout">
			<nav class="sidebar">
				{#each sections as s}
					<button class="nav-item" class:active={activeSection === s.id} onclick={() => activeSection = s.id}>
						<span>{s.icon}</span><span>{s.label}</span>
					</button>
				{/each}
			</nav>

			<main class="content">
				{#if activeSection === 'theme'}
					<div class="card">
						<div class="card-header"><h2>Theme Presets</h2></div>
						<div class="card-body">
							<div class="grid-3">
								{#each appearance.presets as preset}
									<button class="theme-card" class:active={appearance.selectedPresetId === preset.id} onclick={() => selectPreset(preset.id)}>
										<div class="preview-item" style="background: {getPresetColor(preset)}" class:active={appearance.selectedPresetId === preset.id}></div>
										<span class="text-sm">{preset.name}</span>
										{#if preset.tier === 'pro'}<span class="badge pro">PRO</span>{/if}
									</button>
								{/each}
							</div>
						</div>
					</div>
					<div class="card">
						<div class="card-header"><h2>Page Mode</h2></div>
						<div class="card-body">
							<div class="option-group">
								<button class="option-btn" class:active={appearance.settings.page.mode === 'light'} onclick={() => updateSetting('page.mode', 'light')}>☀️ Light</button>
								<button class="option-btn" class:active={appearance.settings.page.mode === 'dark'} onclick={() => updateSetting('page.mode', 'dark')}>🌙 Dark</button>
							</div>
						</div>
					</div>
				{/if}

				{#if activeSection === 'background'}
					<div class="card">
						<div class="card-header"><h2>Background Type</h2></div>
						<div class="card-body">
							<div class="option-group">
								{#each backgroundTypes as t}<button class="option-btn" class:active={appearance.settings.background.type === t.value} onclick={() => updateSetting('background.type', t.value)}>{t.label}</button>{/each}
							</div>
						</div>
					</div>
					{#if appearance.settings.background.type === 'solid'}
						<div class="card">
							<div class="card-header"><h2>Color</h2></div>
							<div class="card-body">
								<div class="color-picker-row">
									<input type="color" class="color-input" value={appearance.settings.background.color} onchange={(e) => updateSetting('background.color', e.currentTarget.value)}/>
									<input type="text" class="color-text" value={appearance.settings.background.color} onchange={(e) => updateSetting('background.color', e.currentTarget.value)}/>
								</div>
								<div class="color-presets">{#each presetColors as c}<button class="color-preset" style="background:{c}" class:active={appearance.settings.background.color === c} title={c} onclick={() => updateSetting('background.color', c)}></button>{/each}</div>
							</div>
						</div>
					{/if}
					{#if appearance.settings.background.type === 'gradient'}
						<div class="card">
							<div class="card-header"><h2>Gradient</h2></div>
							<div class="card-body">
								<div class="grid-3">{#each gradientPresets as g, i}<button class="preview-item" style="background:{g}" class:active={appearance.settings.background.gradient === g} title="Gradient {i+1}" onclick={() => updateSetting('background.gradient', g)}></button>{/each}</div>
							</div>
						</div>
					{/if}
					<div class="card">
						<div class="card-header"><h2>Effects</h2></div>
						<div class="card-body">
							<div class="slider-group"><div class="slider-header"><span class="slider-label">Blur</span><span class="slider-value">{appearance.settings.background.effects.blur}px</span></div><input type="range" class="slider" min="0" max="20" value={appearance.settings.background.effects.blur} oninput={(e) => updateSetting('background.effects.blur', +e.currentTarget.value)}/></div>
							<div class="slider-group"><div class="slider-header"><span class="slider-label">Dim</span><span class="slider-value">{Math.round(appearance.settings.background.effects.dim * 100)}%</span></div><input type="range" class="slider" min="0" max="80" step="5" value={appearance.settings.background.effects.dim * 100} oninput={(e) => updateSetting('background.effects.dim', +e.currentTarget.value / 100)}/></div>
						</div>
					</div>
				{/if}

				{#if activeSection === 'header'}
					<div class="card">
						<div class="card-header"><h2>Header Style</h2></div>
						<div class="card-body"><div class="option-group">{#each headerStyles as s}<button class="option-btn" class:active={appearance.settings.header.style === s.value} onclick={() => updateSetting('header.style', s.value)}>{s.label}</button>{/each}</div></div>
					</div>
					<div class="card">
						<div class="card-header"><h2>Avatar</h2></div>
						<div class="card-body">
							<div class="toggle-row"><span>Show Avatar</span><button class="toggle" class:active={appearance.settings.header.showAvatar} aria-label="Toggle" onclick={() => updateSetting('header.showAvatar', !appearance.settings.header.showAvatar)}></button></div>
							{#if appearance.settings.header.showAvatar}
								<div class="setting-row"><span class="setting-label">Size</span><div class="option-group compact">{#each avatarSizes as s}<button class="option-btn small" class:active={appearance.settings.header.avatarSize === s.value} onclick={() => updateSetting('header.avatarSize', s.value)}>{s.label}</button>{/each}</div></div>
								<div class="setting-row"><span class="setting-label">Shape</span><div class="option-group compact">{#each avatarShapes as s}<button class="option-btn small" class:active={appearance.settings.header.avatarShape === s.value} onclick={() => updateSetting('header.avatarShape', s.value)}>{s.label}</button>{/each}</div></div>
							{/if}
							<div class="toggle-row"><span>Show Bio</span><button class="toggle" class:active={appearance.settings.header.showBio} aria-label="Toggle" onclick={() => updateSetting('header.showBio', !appearance.settings.header.showBio)}></button></div>
						</div>
					</div>
				{/if}

				{#if activeSection === 'links'}
					<div class="card">
						<div class="card-header"><h2>Link Style</h2></div>
						<div class="card-body">
							<div class="grid-4">
								{#each linkStyles as s}
									<button class="link-style-card" class:active={appearance.settings.links.style === s.value} onclick={() => updateSetting('links.style', s.value)}>
										<div class="link-demo {s.value}">Link</div>
										<span class="text-sm">{s.label}</span>
									</button>
								{/each}
							</div>
						</div>
					</div>
					<div class="card">
						<div class="card-header"><h2>Appearance</h2></div>
						<div class="card-body">
							<div class="slider-group"><div class="slider-header"><span class="slider-label">Border Radius</span><span class="slider-value">{appearance.settings.links.borderRadius}px</span></div><input type="range" class="slider" min="0" max="24" step="2" value={appearance.settings.links.borderRadius} oninput={(e) => updateSetting('links.borderRadius', +e.currentTarget.value)}/></div>
							<div class="slider-group"><div class="slider-header"><span class="slider-label">Padding</span><span class="slider-value">{appearance.settings.links.padding}px</span></div><input type="range" class="slider" min="8" max="24" step="2" value={appearance.settings.links.padding} oninput={(e) => updateSetting('links.padding', +e.currentTarget.value)}/></div>
							<div class="slider-group"><div class="slider-header"><span class="slider-label">Gap</span><span class="slider-value">{appearance.settings.links.gap}px</span></div><input type="range" class="slider" min="4" max="24" step="2" value={appearance.settings.links.gap} oninput={(e) => updateSetting('links.gap', +e.currentTarget.value)}/></div>
							<div class="setting-row"><span class="setting-label">Shadow</span><div class="option-group compact">{#each shadowOptions as s}<button class="option-btn small" class:active={appearance.settings.links.shadow === s.value} onclick={() => updateSetting('links.shadow', s.value)}>{s.label}</button>{/each}</div></div>
						</div>
					</div>
				{/if}

				{#if activeSection === 'colors'}
					<div class="card">
						<div class="card-header"><h2>Primary Color</h2></div>
						<div class="card-body">
							<div class="color-picker-row"><input type="color" class="color-input" value={appearance.settings.colors.primary} onchange={(e) => updateSetting('colors.primary', e.currentTarget.value)}/><input type="text" class="color-text" value={appearance.settings.colors.primary} onchange={(e) => updateSetting('colors.primary', e.currentTarget.value)}/></div>
							<div class="color-presets">{#each presetColors as c}<button class="color-preset" style="background:{c}" class:active={appearance.settings.colors.primary === c} title={c} onclick={() => updateSetting('colors.primary', c)}></button>{/each}</div>
						</div>
					</div>
					<div class="card">
						<div class="card-header"><h2>Text Colors</h2></div>
						<div class="card-body">
							<div class="setting-row"><span class="setting-label">Primary</span><div class="color-picker-row"><input type="color" class="color-input small" value={appearance.settings.colors.text} onchange={(e) => updateSetting('colors.text', e.currentTarget.value)}/><input type="text" class="color-text small" value={appearance.settings.colors.text} onchange={(e) => updateSetting('colors.text', e.currentTarget.value)}/></div></div>
							<div class="setting-row"><span class="setting-label">Secondary</span><div class="color-picker-row"><input type="color" class="color-input small" value={appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('colors.textSecondary', e.currentTarget.value)}/><input type="text" class="color-text small" value={appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('colors.textSecondary', e.currentTarget.value)}/></div></div>
						</div>
					</div>
					<div class="card">
						<div class="card-header"><h2>Surface Colors</h2></div>
						<div class="card-body">
							<div class="setting-row"><span class="setting-label">Page BG</span><div class="color-picker-row"><input type="color" class="color-input small" value={appearance.settings.colors.background} onchange={(e) => updateSetting('colors.background', e.currentTarget.value)}/><input type="text" class="color-text small" value={appearance.settings.colors.background} onchange={(e) => updateSetting('colors.background', e.currentTarget.value)}/></div></div>
							<div class="setting-row"><span class="setting-label">Card BG</span><div class="color-picker-row"><input type="color" class="color-input small" value={appearance.settings.colors.cardBackground} onchange={(e) => updateSetting('colors.cardBackground', e.currentTarget.value)}/><input type="text" class="color-text small" value={appearance.settings.colors.cardBackground} onchange={(e) => updateSetting('colors.cardBackground', e.currentTarget.value)}/></div></div>
						</div>
					</div>
				{/if}

				{#if activeSection === 'typography'}
					<div class="card">
						<div class="card-header"><h2>Alignment</h2></div>
						<div class="card-body">
							<div class="setting-row"><span class="setting-label">Page</span><div class="segmented-control">{#each alignOptions as a}<button class:active={appearance.settings.page.layout.textAlign === a.value} title={a.value} onclick={() => updateSetting('page.layout.textAlign', a.value)}>{a.icon}</button>{/each}</div></div>
							<div class="setting-row"><span class="setting-label">Links</span><div class="segmented-control">{#each alignOptions as a}<button class:active={appearance.settings.links.textAlign === a.value} title={a.value} onclick={() => updateSetting('links.textAlign', a.value)}>{a.icon}</button>{/each}</div></div>
						</div>
					</div>
					<div class="card">
						<div class="card-header"><h2>Font Size</h2></div>
						<div class="card-body">
							<div class="setting-row"><span class="setting-label">Base</span><div class="option-group compact">{#each fontSizes as s}<button class="option-btn small" class:active={appearance.settings.page.layout.baseFontSize === s.value} onclick={() => updateSetting('page.layout.baseFontSize', s.value)}>{s.label}</button>{/each}</div></div>
							<div class="setting-row"><span class="setting-label">Links</span><div class="option-group compact">{#each fontSizes as s}<button class="option-btn small" class:active={appearance.settings.links.fontSize === s.value} onclick={() => updateSetting('links.fontSize', s.value)}>{s.label}</button>{/each}</div></div>
						</div>
					</div>
					<div class="card">
						<div class="card-header"><h2>Spacing</h2></div>
						<div class="card-body">
							<div class="slider-group"><div class="slider-header"><span class="slider-label">Page Padding</span><span class="slider-value">{appearance.settings.page.layout.pagePadding}px</span></div><input type="range" class="slider" min="0" max="32" step="4" value={appearance.settings.page.layout.pagePadding} oninput={(e) => updateSetting('page.layout.pagePadding', +e.currentTarget.value)}/></div>
							<div class="slider-group"><div class="slider-header"><span class="slider-label">Block Gap</span><span class="slider-value">{appearance.settings.page.layout.blockGap}px</span></div><input type="range" class="slider" min="4" max="32" step="4" value={appearance.settings.page.layout.blockGap} oninput={(e) => updateSetting('page.layout.blockGap', +e.currentTarget.value)}/></div>
						</div>
					</div>
				{/if}
			</main>

			<aside class="preview">
				<div class="preview-header"><span class="text-sm font-medium text-secondary">Preview</span></div>
				<div class="preview-frame">
					<div class="preview-content" style="background:{getBgStyle()}; padding:{appearance.settings.page.layout.pagePadding}px; text-align:{appearance.settings.page.layout.textAlign}">
						{#if appearance.settings.background.effects.dim > 0}<div class="preview-overlay" style="background:rgba(0,0,0,{appearance.settings.background.effects.dim})"></div>{/if}
						{#if appearance.settings.header.showAvatar}<div class="preview-avatar" style="width:{getAvatarSize()}px; height:{getAvatarSize()}px; border-radius:{getAvatarRadius()}">{auth.user?.email?.charAt(0).toUpperCase() || 'U'}</div>{/if}
						<div class="preview-name" style="color:{appearance.settings.colors.text}">@{auth.user?.username || 'username'}</div>
						{#if appearance.settings.header.showBio}<div class="preview-bio" style="color:{appearance.settings.colors.textSecondary}">Your bio here</div>{/if}
						<div class="preview-links" style="gap:{appearance.settings.links.gap}px">
							{#each ['Website', 'Instagram', 'Twitter'] as link}
								<div class="preview-link" style="border-radius:{appearance.settings.links.borderRadius}px; padding:{appearance.settings.links.padding}px; text-align:{appearance.settings.links.textAlign}; box-shadow:{getShadow()}; background:{getLinkBg()}; color:{getLinkColor()}; border:{getLinkBorder()}">{link}</div>
							{/each}
						</div>
					</div>
				</div>
			</aside>
		</div>
	</div>
{/if}

<style>
	/* Page Layout - Only layout-specific styles */
	.page-container { max-width: 1200px; margin: 0 auto; }
	.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: var(--space-4); }
	.page-header h1 { font-size: var(--text-xl); font-weight: 600; margin-bottom: var(--space-1); }
	.header-actions { display: flex; gap: var(--space-3); }
	.btn-save { display: flex; align-items: center; gap: var(--space-2); }
	.btn-save .spinner { width: 16px; height: 16px; border: 2px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 0.8s linear infinite; }
	.alert-banner { margin-bottom: var(--space-4); }

	/* 3-Column Layout */
	.layout { display: grid; grid-template-columns: 160px 1fr 300px; gap: var(--space-5); }
	.sidebar { display: flex; flex-direction: column; gap: var(--space-1); position: sticky; top: 80px; height: fit-content; }
	.nav-item { display: flex; align-items: center; gap: var(--space-3); padding: var(--space-3); background: transparent; border-radius: var(--radius-md); text-align: left; font-size: var(--text-sm); font-weight: 500; }
	.nav-item:hover { background: var(--color-bg); }
	.nav-item.active { background: var(--color-primary-light); color: var(--color-primary); }
	.content { display: flex; flex-direction: column; gap: var(--space-4); }

	/* Theme Card */
	.theme-card { position: relative; background: transparent; padding: 0; text-align: center; }
	.theme-card .badge { position: absolute; top: 4px; right: 4px; }

	/* Link Style Card */
	.link-style-card { background: transparent; padding: 0; text-align: center; }
	.link-style-card.active span { color: var(--color-primary); font-weight: 500; }
	.link-demo { padding: var(--space-3); border-radius: var(--radius-md); font-size: var(--text-sm); margin-bottom: var(--space-2); }
	.link-demo.filled { background: var(--color-primary); color: white; }
	.link-demo.outline { background: transparent; border: 1px solid var(--color-primary); color: var(--color-primary); }
	.link-demo.soft { background: var(--color-primary-light); color: var(--color-primary); }
	.link-demo.ghost { background: transparent; color: var(--color-text); }
	.link-style-card.active .link-demo { box-shadow: 0 0 0 2px var(--color-primary); }

	/* Preview Panel */
	.preview { position: sticky; top: 80px; height: fit-content; }
	.preview-header { margin-bottom: var(--space-3); }
	.preview-frame { background: var(--color-bg-secondary); border-radius: var(--radius-xl); padding: var(--space-3); box-shadow: var(--shadow-lg); }
	.preview-content { position: relative; min-height: 480px; border-radius: var(--radius-lg); overflow: hidden; display: flex; flex-direction: column; align-items: center; padding-top: var(--space-8); }
	.preview-overlay { position: absolute; inset: 0; pointer-events: none; }
	.preview-avatar { position: relative; z-index: 1; background: linear-gradient(135deg, #ff9500, #ff3b30); display: flex; align-items: center; justify-content: center; font-size: 1.5rem; font-weight: 600; color: white; margin-bottom: var(--space-3); }
	.preview-name { position: relative; z-index: 1; font-size: var(--text-lg); font-weight: 600; margin-bottom: var(--space-1); }
	.preview-bio { position: relative; z-index: 1; font-size: var(--text-sm); margin-bottom: var(--space-5); }
	.preview-links { position: relative; z-index: 1; width: 100%; display: flex; flex-direction: column; }
	.preview-link { font-size: var(--text-sm); font-weight: 500; }

	/* Responsive */
	@media (max-width: 1024px) {
		.layout { grid-template-columns: 1fr; }
		.sidebar { flex-direction: row; overflow-x: auto; position: static; padding-bottom: var(--space-3); border-bottom: 1px solid var(--color-separator); }
		.nav-item { flex-shrink: 0; }
		.preview { display: none; }
		.page-header { flex-direction: column; align-items: flex-start; gap: var(--space-3); }
		.header-actions { width: 100%; justify-content: flex-end; }
	}
</style>
