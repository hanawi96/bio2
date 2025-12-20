<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getAuth } from '$lib/stores/auth.svelte';
	import { getAppearance, loadAppearance, selectPreset, updateSetting, saveAppearance, resetAppearance, resetToPresetDefaults } from '$lib/stores/appearance.svelte';
	import { bio } from '$lib/api/client';
	import { Palette, Image, User, Link, Droplets, Type, Save, AlertTriangle, Sun, Moon, AlignLeft, AlignCenter, AlignRight, Settings, RefreshCw, Instagram, Music, Facebook, Twitter, Youtube, Linkedin, Github, Globe } from 'lucide-svelte';

	const auth = getAuth();
	const appearance = getAppearance();
	let pageId = $state<number | null>(null);
	let activeSection = $state<string>('theme');
	let activeHeaderTab = $state<string>('avatar'); // Sub-tab for header settings
	let showDebug = $state(false);

	const sections = [
		{ id: 'theme', label: 'Theme', icon: 'palette' },
		{ id: 'colors', label: 'Colors', icon: 'droplets' },
		{ id: 'background', label: 'Background', icon: 'image' },
		{ id: 'header', label: 'Profile', icon: 'user' },
		{ id: 'blocks', label: 'Blocks', icon: 'link' },
		{ id: 'typography', label: 'Typography', icon: 'type' }
	];

	const validSections = new Set(sections.map(s => s.id));

	// Đồng bộ activeSection với URL (hỗ trợ browser back/forward)
	$effect(() => {
		const urlTab = $page.url.searchParams.get('tab');
		if (urlTab && validSections.has(urlTab) && urlTab !== activeSection) {
			activeSection = urlTab;
		} else if (!urlTab && activeSection !== 'theme') {
			activeSection = 'theme';
		}
	});

	onMount(async () => {
		// Đọc tab từ URL search params khi load lần đầu
		const urlTab = $page.url.searchParams.get('tab');
		if (urlTab && validSections.has(urlTab)) {
			activeSection = urlTab;
		}

		const bioData = await bio.get();
		if (bioData?.page) {
			pageId = bioData.page.id;
			await loadAppearance(bioData.page.id);
		}
	});

	async function handleSave() { await saveAppearance(); }

	// Apply color scheme
	function applyColorScheme(scheme: typeof colorSchemes[0]) {
		updateSetting('colors.primary', scheme.primary);
		updateSetting('colors.background', scheme.bg);
		updateSetting('colors.text', scheme.text);
		updateSetting('colors.textSecondary', scheme.textSecondary);
	}

	// Hàm chuyển tab với URL update
	function switchTab(tabId: string) {
		if (!validSections.has(tabId)) return;
		
		activeSection = tabId;
		
		// Cập nhật URL mà không reload trang
		const url = new URL(window.location.href);
		url.searchParams.set('tab', tabId);
		goto(url.pathname + url.search, { replaceState: false, noScroll: true, keepFocus: true });
	}

	const fontSizes = [{ value: 'S', label: 'S' }, { value: 'M', label: 'M' }, { value: 'L', label: 'L' }, { value: 'XL', label: 'XL' }];
	const alignOptions = [{ value: 'left', icon: 'left' }, { value: 'center', icon: 'center' }, { value: 'right', icon: 'right' }];
	const linkStyles = [{ value: 'filled', label: 'Filled' }, { value: 'outline', label: 'Outline' }, { value: 'soft', label: 'Soft' }, { value: 'ghost', label: 'Ghost' }];
	const shadowOptions = [{ value: 'none', label: 'None' }, { value: 'sm', label: 'S' }, { value: 'md', label: 'M' }, { value: 'lg', label: 'L' }];
	const headerStyles = [{ value: 'classic', label: 'Classic' }, { value: 'minimal', label: 'Minimal' }];
	const avatarSizes = [{ value: 'S', label: 'S' }, { value: 'M', label: 'M' }, { value: 'L', label: 'L' }];
	const avatarShapes = [{ value: 'circle', label: 'Circle' }, { value: 'rounded', label: 'Rounded' }, { value: 'square', label: 'Square' }];
	const backgroundTypes = [{ value: 'solid', label: 'Solid' }, { value: 'gradient', label: 'Gradient' }, { value: 'image', label: 'Image' }];
	const presetColors = ['#007aff', '#5856d6', '#af52de', '#ff2d55', '#ff3b30', '#ff9500', '#ffcc00', '#34c759', '#00c7be', '#30b0c7', '#1c1c1e', '#8e8e93'];
	
	const colorSchemes = [
		{ name: 'Neutral', primary: '#1c1c1e', bg: '#f2f2f7', text: '#1c1c1e', textSecondary: '#8e8e93' },
		{ name: 'Bold', primary: '#ff3b30', bg: '#ffffff', text: '#1c1c1e', textSecondary: '#8e8e93' },
		{ name: 'Ocean', primary: '#007aff', bg: '#e3f2fd', text: '#1c1c1e', textSecondary: '#6b7280' },
		{ name: 'Sunset', primary: '#ff9500', bg: '#fff7ed', text: '#1c1c1e', textSecondary: '#78716c' },
		{ name: 'Forest', primary: '#34c759', bg: '#f0fdf4', text: '#1c1c1e', textSecondary: '#6b7280' },
		{ name: 'Purple', primary: '#af52de', bg: '#faf5ff', text: '#1c1c1e', textSecondary: '#78716c' }
	];
	
	const gradientPresets = [
		'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
		'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
		'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
		'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
		'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
		'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)'
	];

	// Lấy màu preview cho theme card
	function getPresetColor(preset: any): string {
		const config = preset.config;
		if (!config) return '#f2f2f7';
		
		// Ưu tiên lấy từ background
		if (config.background) {
			if (config.background.gradient) return config.background.gradient;
			if (config.background.color) return config.background.color;
		}
		
		// Fallback theo tên
		const colors: Record<string, string> = {
			'Light': '#f2f2f7',
			'Dark': '#1c1c1e',
			'Ocean': 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
			'Sunset': 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
			'Forest': 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
			'Aurora Glass': 'linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%)'
		};
		return colors[preset.name] || '#f2f2f7';
	}

	// Lấy background style cho preview - dùng settings đã computed
	function getBgStyle(): string {
		const bg = appearance.settings.background;
		if (bg.type === 'gradient' && bg.gradient) return bg.gradient;
		if (bg.type === 'image' && bg.imageUrl) return `url(${bg.imageUrl}) center/cover`;
		return bg.color;
	}

	function getAvatarSize(): number {
		return { S: 48, M: 64, L: 80 }[appearance.settings.header.avatarSize] || 64;
	}

	function getAvatarRadius(): string {
		return { circle: '50%', rounded: '16px', square: '8px' }[appearance.settings.header.avatarShape] || '50%';
	}

	function getShadow(): string {
		return {
			none: 'none',
			sm: '0 1px 2px rgba(0,0,0,0.08)',
			md: '0 2px 8px rgba(0,0,0,0.12)',
			lg: '0 4px 16px rgba(0,0,0,0.16)'
		}[appearance.settings.links.shadow] || 'none';
	}

	function getLinkBg(): string {
		const s = appearance.settings.links.style;
		const colors = appearance.settings.colors;
		if (s === 'filled') return colors.primary;
		if (s === 'soft') return `${colors.primary}15`;
		return colors.cardBackground;
	}

	function getLinkColor(): string {
		// Nếu có custom textColor, dùng nó
		if (appearance.settings.links.textColor) {
			return appearance.settings.links.textColor;
		}
		
		// Auto based on style
		const s = appearance.settings.links.style;
		if (s === 'filled') return '#fff';
		if (s === 'soft' || s === 'outline') return appearance.settings.colors.primary;
		return appearance.settings.colors.text;
	}

	function getLinkBorder(): string {
		if (appearance.settings.links.style === 'outline') {
			return `1px solid ${appearance.settings.colors.primary}`;
		}
		return 'none';
	}

	// Header style configurations
	function getHeaderContainerStyle(): string {
		const style = appearance.settings.header.style;
		const align = appearance.settings.header.align;
		
		// Map align to flexbox alignment
		const alignItems = align === 'left' ? 'flex-start' : align === 'right' ? 'flex-end' : 'center';
		
		switch (style) {
			case 'classic':
				// Avatar overlap cover, content below
				// Calculate padding based on avatar size
				const baseSize = getAvatarSize();
				const avatarSize = Math.round(baseSize * 1.3);
				const topPadding = Math.round(avatarSize / 2) + 12; // Half avatar + gap
				return `padding: ${topPadding}px 16px 16px 16px; gap: 12px; align-items: center; margin-top: 0;`;
			case 'minimal':
				// No cover, simple layout
				return `padding: 16px 0; gap: 8px; align-items: ${alignItems};`;
			default:
				return `padding: 16px 0; gap: 12px; align-items: ${alignItems};`;
		}
	}

	function getHeaderAvatarStyle(): string {
		const style = appearance.settings.header.style;
		const baseSize = getAvatarSize();
		const radius = getAvatarRadius();
		// Explicitly read these to ensure reactivity
		const borderColor = appearance.settings.header.avatarBorderColor || 'rgba(0,0,0,0.15)';
		const borderWidth = appearance.settings.header.avatarBorderWidth ?? 3;
		
		let size = baseSize;
		if (style === 'minimal') size = Math.round(baseSize * 1.2); // Larger for minimal
		if (style === 'classic') size = Math.round(baseSize * 1.3);
		
		// Classic style: avatar overlap cover (centered) - positioned at bottom of cover
		if (style === 'classic') {
			const coverHeight = 160;
			const topPosition = coverHeight - (size / 2);
			return `width: ${size}px; height: ${size}px; border-radius: ${radius}; border: ${borderWidth}px solid ${borderColor}; top: ${topPosition}px; box-shadow: 0 4px 12px rgba(0,0,0,0.15);`;
		}
		
		// Minimal style: clean avatar with subtle shadow
		if (style === 'minimal') {
			return `width: ${size}px; height: ${size}px; border-radius: ${radius}; box-shadow: 0 2px 12px rgba(0,0,0,0.08);`;
		}
		
		return `width: ${size}px; height: ${size}px; border-radius: ${radius};`;
	}

	function getHeaderNameStyle(): string {
		const style = appearance.settings.header.style;
		const nameSettings = appearance.settings.header;
		const color = nameSettings.nameColor || appearance.settings.colors.text;
		
		let fontSize = '1.125rem';
		let fontWeight = '600';
		let textColor = color;
		
		if (style === 'minimal') { 
			fontSize = '1.25rem'; // Larger for minimal style
			fontWeight = '700';
		}
		if (style === 'classic') { 
			fontSize = '1.25rem'; 
			fontWeight = '700'; 
		}
		
		return `color: ${textColor}; font-size: ${fontSize}; font-weight: ${fontWeight};`;
	}

	function getHeaderBioStyle(): string {
		const style = appearance.settings.header.style;
		const bioSettings = appearance.settings.header;
		const color = bioSettings.bioColor || appearance.settings.colors.textSecondary;
		
		// Map bio size to font-size
		const sizeMap: Record<string, string> = { S: '0.75rem', M: '0.875rem', L: '1rem', XL: '1.125rem' };
		let fontSize = sizeMap[bioSettings.bioSize] || '0.875rem';
		let textColor = color;
		
		// Adjust based on header style
		if (style === 'minimal') {
			const minimalSizeMap: Record<string, string> = { S: '0.8125rem', M: '0.875rem', L: '0.9375rem', XL: '1rem' };
			fontSize = minimalSizeMap[bioSettings.bioSize] || '0.875rem';
		}
		if (style === 'classic') {
			const classicSizeMap: Record<string, string> = { S: '0.875rem', M: '1rem', L: '1.125rem', XL: '1.25rem' };
			fontSize = classicSizeMap[bioSettings.bioSize] || '1rem';
		}
		
		const maxWidth = style === 'minimal' ? '280px' : style === 'classic' ? '320px' : '280px';
		
		// Map align to align-self for flexbox
		const alignSelfMap: Record<string, string> = {
			left: 'flex-start',
			center: 'center',
			right: 'flex-end'
		};
		const alignSelf = alignSelfMap[bioSettings.bioAlign] || 'center';
		const textAlign = style === 'minimal' ? 'center' : bioSettings.bioAlign;
		
		return `color: ${textColor}; font-size: ${fontSize}; max-width: ${maxWidth}; text-align: ${textAlign}; align-self: ${alignSelf}; line-height: 1.5;`;
	}
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
					{#if appearance.saving}<span class="spinner"></span>{:else}<Save size={16} />{/if}
					<span>{appearance.saving ? 'Đang lưu...' : 'Lưu tất cả'}</span>
				</button>
			</div>
		</header>

		{#if appearance.dirty}<div class="alert-banner warning"><AlertTriangle size={16} /> Bạn có thay đổi chưa lưu</div>{/if}

		<div class="layout">
			<nav class="sidebar">
				{#each sections as s}
					<button class="nav-item" class:active={activeSection === s.id} onclick={() => switchTab(s.id)}>
						<span class="nav-icon">
							{#if s.icon === 'palette'}<Palette size={18} />
							{:else if s.icon === 'droplets'}<Droplets size={18} />
							{:else if s.icon === 'image'}<Image size={18} />
							{:else if s.icon === 'user'}<User size={18} />
							{:else if s.icon === 'link'}<Link size={18} />
							{:else if s.icon === 'type'}<Type size={18} />
							{/if}
						</span>
						<span>{s.label}</span>
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
								<button class="option-btn" class:active={appearance.settings.page.mode === 'light'} onclick={() => updateSetting('page.mode', 'light')}><Sun size={16} /> Light</button>
								<button class="option-btn" class:active={appearance.settings.page.mode === 'dark'} onclick={() => updateSetting('page.mode', 'dark')}><Moon size={16} /> Dark</button>
							</div>
						</div>
					</div>
				{/if}

				{#if activeSection === 'colors'}
					<!-- Color Schemes -->
					<div class="card">
						<div class="card-header"><h2>Color Schemes</h2></div>
						<div class="card-body">
							<p class="form-hint" style="margin-bottom: var(--space-3)">Quick start with pre-made color combinations</p>
							<div class="grid-3">
								{#each colorSchemes as scheme}
									<button class="color-scheme-card" onclick={() => applyColorScheme(scheme)}>
										<div class="color-scheme-preview">
											<div class="color-dot-large" style="background:{scheme.primary}"></div>
											<div class="color-dots-small">
												<div class="color-dot-small" style="background:{scheme.bg}"></div>
												<div class="color-dot-small" style="background:{scheme.text}"></div>
											</div>
										</div>
										<span class="text-sm">{scheme.name}</span>
									</button>
								{/each}
							</div>
						</div>
					</div>

					<!-- Text Colors -->
					<div class="card">
						<div class="card-header"><h2>Text Colors</h2></div>
						<div class="card-body">
							<!-- Primary Text Color -->
							<div class="color-setting-group">
								<div class="color-setting-header">
									<span class="setting-label">Primary Text</span>
									<div class="color-value-display">
										<input type="color" class="color-picker-trigger" value={appearance.settings.colors.text} onchange={(e) => updateSetting('colors.text', e.currentTarget.value)} title="Pick custom color"/>
										<input type="text" class="color-hex-input" value={appearance.settings.colors.text} onchange={(e) => updateSetting('colors.text', e.currentTarget.value)} placeholder="#000000"/>
									</div>
								</div>
								<div class="color-presets-inline">
									{#each ['#000000', '#1c1c1e', '#2c2c2e', '#3a3a3c', '#48484a', '#636366', '#8e8e93', '#ffffff'] as c}
										<button 
											class="color-preset-btn" 
											class:active={appearance.settings.colors.text === c} 
											style="background:{c}" 
											title={c} 
											onclick={() => updateSetting('colors.text', c)}
										></button>
									{/each}
								</div>
							</div>

							<!-- Secondary Text Color -->
							<div class="color-setting-group">
								<div class="color-setting-header">
									<span class="setting-label">Secondary Text</span>
									<div class="color-value-display">
										<input type="color" class="color-picker-trigger" value={appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('colors.textSecondary', e.currentTarget.value)} title="Pick custom color"/>
										<input type="text" class="color-hex-input" value={appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('colors.textSecondary', e.currentTarget.value)} placeholder="#000000"/>
									</div>
								</div>
								<div class="color-presets-inline">
									{#each ['#636366', '#8e8e93', '#aeaeb2', '#c7c7cc', '#d1d1d6', '#e5e5ea', '#f2f2f7'] as c}
										<button 
											class="color-preset-btn" 
											class:active={appearance.settings.colors.textSecondary === c} 
											style="background:{c}" 
											title={c} 
											onclick={() => updateSetting('colors.textSecondary', c)}
										></button>
									{/each}
								</div>
							</div>
						</div>
					</div>

					<!-- Advanced Colors -->
					<div class="card">
						<div class="card-header"><h2>Advanced Color Settings</h2></div>
						<div class="card-body">
							<p class="form-hint" style="margin-bottom: var(--space-4)">Fine-tune colors for specific sections</p>

							<!-- Bio Text Color -->
							<div class="color-setting-group">
								<div class="color-setting-header">
									<span class="setting-label">Bio Text Color</span>
									<div class="color-value-display">
										<input type="color" class="color-picker-trigger" value={appearance.settings.header.bioColor || appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('header.bioColor', e.currentTarget.value)} title="Pick custom color"/>
										<input type="text" class="color-hex-input" value={appearance.settings.header.bioColor || appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('header.bioColor', e.currentTarget.value)} placeholder="Auto"/>
										<button class="btn-icon-reset" onclick={() => updateSetting('header.bioColor', '')} title="Reset to auto"><RefreshCw size={14} /></button>
									</div>
								</div>
								<div class="color-presets-inline">
									{#each ['#636366', '#8e8e93', '#aeaeb2', '#c7c7cc', '#1c1c1e', '#3a3a3c', '#48484a'] as c}
										<button 
											class="color-preset-btn" 
											class:active={(appearance.settings.header.bioColor || appearance.settings.colors.textSecondary) === c} 
											style="background:{c}" 
											title={c} 
											onclick={() => updateSetting('header.bioColor', c)}
										></button>
									{/each}
								</div>
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
							<div class="card-header"><h2>Background Color</h2></div>
							<div class="card-body">
								<div class="color-setting-group">
									<div class="color-setting-header">
										<span class="setting-label">Solid Color</span>
										<div class="color-value-display">
											<input type="color" class="color-picker-trigger" value={appearance.settings.background.color} onchange={(e) => updateSetting('background.color', e.currentTarget.value)} title="Pick custom color"/>
											<input type="text" class="color-hex-input" value={appearance.settings.background.color} onchange={(e) => updateSetting('background.color', e.currentTarget.value)} placeholder="#000000"/>
										</div>
									</div>
									<div class="color-presets-inline">
										{#each presetColors as c}
											<button 
												class="color-preset-btn" 
												class:active={appearance.settings.background.color === c} 
												style="background:{c}" 
												title={c} 
												onclick={() => updateSetting('background.color', c)}
											></button>
										{/each}
									</div>
								</div>
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
					<!-- Header Style -->
					<div class="card">
						<div class="card-header"><h2>Kiểu Header</h2></div>
						<div class="card-body">
							<div class="header-style-grid">
								{#each headerStyles as s}
									<button 
										class="header-style-card" 
										class:active={appearance.settings.header.style === s.value} 
										onclick={() => updateSetting('header.style', s.value)}
										title={s.label}
									>
										<div class="header-style-preview {s.value}">
											<div class="preview-avatar-mini"></div>
											<div class="preview-text-mini">
												<div class="preview-name-mini">Username</div>
												<div class="preview-bio-mini">Bio text here</div>
											</div>
										</div>
										<span class="header-style-label">{s.label}</span>
									</button>
								{/each}
							</div>
						</div>
					</div>

					<!-- Settings Card with Sub-tabs -->
					<div class="card">
						<div class="card-header"><h2>Cài đặt</h2></div>
						<div class="card-body">
							<!-- Sub-tabs Navigation -->
							<div class="sub-tabs">
								<button class="sub-tab" class:active={activeHeaderTab === 'avatar'} onclick={() => activeHeaderTab = 'avatar'}>
									<User size={16} />
									<span>Avatar</span>
								</button>
								<button class="sub-tab" class:active={activeHeaderTab === 'bio'} onclick={() => activeHeaderTab = 'bio'}>
									<Type size={16} />
									<span>Bio</span>
								</button>
								<button class="sub-tab" class:active={activeHeaderTab === 'cover'} onclick={() => activeHeaderTab = 'cover'}>
									<Image size={16} />
									<span>Cover</span>
								</button>
								<button class="sub-tab" class:active={activeHeaderTab === 'social'} onclick={() => activeHeaderTab = 'social'}>
									<Instagram size={16} />
									<span>Social</span>
								</button>
							</div>

							<!-- Tab Content -->
							<div class="sub-tab-content">
								{#if activeHeaderTab === 'avatar'}
									<!-- Avatar Settings -->
									<div class="setting-row">
										<span class="setting-label">Căn chỉnh Header</span>
										<div class="segmented-control">
											{#each alignOptions as a}
												<button class:active={appearance.settings.header.align === a.value} title={a.value} onclick={() => updateSetting('header.align', a.value)}>
													{#if a.icon === 'left'}<AlignLeft size={16} />
													{:else if a.icon === 'center'}<AlignCenter size={16} />
													{:else if a.icon === 'right'}<AlignRight size={16} />
													{/if}
												</button>
											{/each}
										</div>
									</div>
									<p class="form-hint" style="margin-bottom: var(--space-3)">Căn chỉnh avatar, username và social icons</p>

									<div class="setting-row">
										<span class="setting-label">Kích thước Avatar</span>
										<div class="option-group compact">
											{#each avatarSizes as s}
												<button class="option-btn small" class:active={appearance.settings.header.avatarSize === s.value} onclick={() => updateSetting('header.avatarSize', s.value)}>{s.label}</button>
											{/each}
										</div>
									</div>
									<div class="setting-row">
										<span class="setting-label">Hình dạng Avatar</span>
										<div class="option-group compact">
											{#each avatarShapes as s}
												<button class="option-btn small" class:active={appearance.settings.header.avatarShape === s.value} onclick={() => updateSetting('header.avatarShape', s.value)}>{s.label}</button>
											{/each}
										</div>
									</div>
									
									<!-- Name Color Picker -->
									<div class="setting-row" style="margin-top: var(--space-4)">
										<span class="setting-label">Màu tên hiển thị</span>
									</div>
									<div class="color-setting-group" style="margin-top: var(--space-2)">
										<div class="color-presets-row">
											<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
												<div class="gradient-icon"></div>
											</button>
											<input type="color" class="hidden-native-picker" value={appearance.settings.header.nameColor || appearance.settings.colors.text} onchange={(e) => updateSetting('header.nameColor', e.currentTarget.value)} />
											
											<button class="color-preset-circle" class:active={(appearance.settings.header.nameColor || appearance.settings.colors.text) === '#000000'} style="background:#000000" title="#000000" onclick={() => updateSetting('header.nameColor', '#000000')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.nameColor || appearance.settings.colors.text) === '#ffffff'} style="background:#ffffff; border: 2px solid #e5e5ea" title="#ffffff" onclick={() => updateSetting('header.nameColor', '#ffffff')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.nameColor || appearance.settings.colors.text) === '#1c1c1e'} style="background:#1c1c1e" title="#1c1c1e" onclick={() => updateSetting('header.nameColor', '#1c1c1e')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.nameColor || appearance.settings.colors.text) === '#8e8e93'} style="background:#8e8e93" title="#8e8e93" onclick={() => updateSetting('header.nameColor', '#8e8e93')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.nameColor || appearance.settings.colors.text) === '#ff3b30'} style="background:#ff3b30" title="#ff3b30" onclick={() => updateSetting('header.nameColor', '#ff3b30')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.nameColor || appearance.settings.colors.text) === '#ff9500'} style="background:#ff9500" title="#ff9500" onclick={() => updateSetting('header.nameColor', '#ff9500')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.nameColor || appearance.settings.colors.text) === '#007aff'} style="background:#007aff" title="#007aff" onclick={() => updateSetting('header.nameColor', '#007aff')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.nameColor || appearance.settings.colors.text) === '#5856d6'} style="background:#5856d6" title="#5856d6" onclick={() => updateSetting('header.nameColor', '#5856d6')}></button>
											
											{#if appearance.settings.header.nameColor}
												<button class="color-preset-circle reset-circle" title="Reset về mặc định" onclick={() => updateSetting('header.nameColor', '')}>
													<RefreshCw size={18} />
												</button>
											{/if}
										</div>
									</div>
									
									<!-- Avatar Border Width Slider -->
									<div class="setting-row" style="margin-top: var(--space-4)">
										<span class="setting-label">Độ dày viền Avatar</span>
									</div>
									<p class="form-hint" style="margin-bottom: var(--space-2)">Chỉ áp dụng cho kiểu header "Classic"</p>
									<div class="slider-group">
										<div class="slider-header">
											<span class="slider-label">Border Width</span>
											<span class="slider-value">{appearance.settings.header.avatarBorderWidth || 3}px</span>
										</div>
										<input 
											type="range" 
											class="slider" 
											min="0" 
											max="8" 
											step="1" 
											value={appearance.settings.header.avatarBorderWidth || 3} 
											oninput={(e) => updateSetting('header.avatarBorderWidth', +e.currentTarget.value)}
										/>
									</div>
									
									<!-- Avatar Border Color Picker -->
									<div class="setting-row" style="margin-top: var(--space-4)">
										<span class="setting-label">Màu viền Avatar</span>
									</div>
									<p class="form-hint" style="margin-bottom: var(--space-2)">Chỉ áp dụng cho kiểu header "Classic"</p>
									<div class="color-setting-group" style="margin-top: var(--space-2)">
										<div class="color-presets-row">
											<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
												<div class="gradient-icon"></div>
											</button>
											<input type="color" class="hidden-native-picker" value={appearance.settings.header.avatarBorderColor || '#000000'} onchange={(e) => updateSetting('header.avatarBorderColor', e.currentTarget.value)} />
											
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || 'rgba(0,0,0,0.15)') === '#000000'} style="background:#000000" title="#000000" onclick={() => updateSetting('header.avatarBorderColor', '#000000')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || 'rgba(0,0,0,0.15)') === '#ffffff'} style="background:#ffffff; border: 2px solid #e5e5ea" title="#ffffff" onclick={() => updateSetting('header.avatarBorderColor', '#ffffff')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || 'rgba(0,0,0,0.15)') === '#1c1c1e'} style="background:#1c1c1e" title="#1c1c1e" onclick={() => updateSetting('header.avatarBorderColor', '#1c1c1e')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || 'rgba(0,0,0,0.15)') === '#8e8e93'} style="background:#8e8e93" title="#8e8e93" onclick={() => updateSetting('header.avatarBorderColor', '#8e8e93')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || 'rgba(0,0,0,0.15)') === '#ff3b30'} style="background:#ff3b30" title="#ff3b30" onclick={() => updateSetting('header.avatarBorderColor', '#ff3b30')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || 'rgba(0,0,0,0.15)') === '#ff9500'} style="background:#ff9500" title="#ff9500" onclick={() => updateSetting('header.avatarBorderColor', '#ff9500')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || 'rgba(0,0,0,0.15)') === '#007aff'} style="background:#007aff" title="#007aff" onclick={() => updateSetting('header.avatarBorderColor', '#007aff')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || 'rgba(0,0,0,0.15)') === '#5856d6'} style="background:#5856d6" title="#5856d6" onclick={() => updateSetting('header.avatarBorderColor', '#5856d6')}></button>
											
											{#if appearance.settings.header.avatarBorderColor}
												<button class="color-preset-circle reset-circle" title="Reset về mặc định" onclick={() => updateSetting('header.avatarBorderColor', '')}>
													<RefreshCw size={18} />
												</button>
											{/if}
										</div>
									</div>

								{:else if activeHeaderTab === 'bio'}
									<!-- Bio Settings -->
									<div class="toggle-row">
										<span>Hiển thị Bio</span>
										<button class="toggle" class:active={appearance.settings.header.showBio} aria-label="Toggle" onclick={() => updateSetting('header.showBio', !appearance.settings.header.showBio)}></button>
									</div>
									{#if appearance.settings.header.showBio}
										<div class="setting-row">
											<span class="setting-label">Cỡ chữ Bio</span>
											<div class="option-group compact">
												{#each fontSizes as s}
													<button class="option-btn small" class:active={appearance.settings.header.bioSize === s.value} onclick={() => updateSetting('header.bioSize', s.value)}>{s.label}</button>
												{/each}
											</div>
										</div>
										<div class="setting-row">
											<span class="setting-label">Căn lề Bio</span>
											<div class="segmented-control">
												{#each alignOptions as a}
													<button class:active={appearance.settings.header.bioAlign === a.value} title={a.value} onclick={() => updateSetting('header.bioAlign', a.value)}>
														{#if a.icon === 'left'}<AlignLeft size={16} />
														{:else if a.icon === 'center'}<AlignCenter size={16} />
														{:else if a.icon === 'right'}<AlignRight size={16} />
														{/if}
													</button>
												{/each}
											</div>
										</div>
										<div class="setting-row">
											<span class="setting-label">Màu Bio</span>
										</div>
										<div class="color-setting-group" style="margin-top: var(--space-2)">
											<div class="color-setting-header">
												<span class="setting-label" style="font-size: var(--text-sm); color: var(--color-text-secondary)">Chọn màu cho bio text</span>
												<div class="color-value-display">
													<input type="color" class="color-picker-trigger" value={appearance.settings.header.bioColor || appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('header.bioColor', e.currentTarget.value)} title="Pick custom color"/>
													<input type="text" class="color-hex-input" value={appearance.settings.header.bioColor || appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('header.bioColor', e.currentTarget.value)} placeholder="Auto"/>
													<button class="btn-icon-reset" onclick={() => updateSetting('header.bioColor', '')} title="Reset về mặc định"><RefreshCw size={14} /></button>
												</div>
											</div>
											<div class="color-presets-inline">
												{#each ['#636366', '#8e8e93', '#aeaeb2', '#c7c7cc', '#1c1c1e', '#3a3a3c', '#48484a'] as c}
													<button 
														class="color-preset-btn" 
														class:active={(appearance.settings.header.bioColor || appearance.settings.colors.textSecondary) === c} 
														style="background:{c}" 
														title={c} 
														onclick={() => updateSetting('header.bioColor', c)}
													></button>
												{/each}
											</div>
										</div>
									{/if}

								{:else if activeHeaderTab === 'cover'}
									<!-- Cover Settings -->
									<p class="form-hint">Cài đặt cover chỉ áp dụng cho kiểu header "Classic"</p>
									<div class="setting-row" style="margin-top: var(--space-3)">
										<span class="setting-label">Màu Cover</span>
									</div>
									<div class="color-setting-group" style="margin-top: var(--space-2)">
										<div class="color-presets-row">
											<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
												<div class="gradient-icon"></div>
											</button>
											<input type="color" class="hidden-native-picker" value={appearance.settings.colors.primary} onchange={(e) => updateSetting('colors.primary', e.currentTarget.value)} />
											
											<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#000000'} style="background:#000000" title="#000000" onclick={() => updateSetting('colors.primary', '#000000')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#ffffff'} style="background:#ffffff; border: 2px solid #e5e5ea" title="#ffffff" onclick={() => updateSetting('colors.primary', '#ffffff')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#667eea'} style="background:#667eea" title="#667eea" onclick={() => updateSetting('colors.primary', '#667eea')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#5856d6'} style="background:#5856d6" title="#5856d6" onclick={() => updateSetting('colors.primary', '#5856d6')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#ff3b30'} style="background:#ff3b30" title="#ff3b30" onclick={() => updateSetting('colors.primary', '#ff3b30')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#ff9500'} style="background:#ff9500" title="#ff9500" onclick={() => updateSetting('colors.primary', '#ff9500')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#34c759'} style="background:#34c759" title="#34c759" onclick={() => updateSetting('colors.primary', '#34c759')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#00c7be'} style="background:#00c7be" title="#00c7be" onclick={() => updateSetting('colors.primary', '#00c7be')}></button>
										</div>
									</div>

								{:else if activeHeaderTab === 'social'}
									<!-- Social Icons Settings -->
									<div class="toggle-row">
										<span>Hiển thị Social Icons</span>
										<button class="toggle" class:active={appearance.settings.header.showSocials} aria-label="Toggle" onclick={() => updateSetting('header.showSocials', !appearance.settings.header.showSocials)}></button>
									</div>
									{#if appearance.settings.header.showSocials}
										<p class="form-hint" style="margin-bottom: var(--space-3)">Chọn các mạng xã hội để hiển thị trên profile</p>
										<div class="social-icons-grid">
											<button class="social-icon-btn"><span class="icon"><Facebook size={18} /></span><span class="label">Facebook</span></button>
											<button class="social-icon-btn active"><span class="icon"><Instagram size={18} /></span><span class="label">Instagram</span></button>
											<button class="social-icon-btn"><span class="icon"><Twitter size={18} /></span><span class="label">Twitter</span></button>
											<button class="social-icon-btn active"><span class="icon"><Music size={18} /></span><span class="label">TikTok</span></button>
											<button class="social-icon-btn"><span class="icon"><Youtube size={18} /></span><span class="label">YouTube</span></button>
											<button class="social-icon-btn"><span class="icon"><Linkedin size={18} /></span><span class="label">LinkedIn</span></button>
											<button class="social-icon-btn"><span class="icon"><Github size={18} /></span><span class="label">GitHub</span></button>
											<button class="social-icon-btn"><span class="icon"><Globe size={18} /></span><span class="label">Website</span></button>
										</div>
									{/if}
								{/if}
							</div>
						</div>
					</div>
				{/if}

				{#if activeSection === 'blocks'}
					<div class="card">
						<div class="card-header"><h2>Block Style</h2></div>
						<div class="card-body">
							<div class="grid-4">
								{#each linkStyles as s}
									<button class="link-style-card" class:active={appearance.settings.links.style === s.value} onclick={() => updateSetting('links.style', s.value)}>
										<div class="link-demo {s.value}">Block</div>
										<span class="text-sm">{s.label}</span>
									</button>
								{/each}
							</div>
						</div>
					</div>

					<!-- Link Block Colors -->
					<div class="card">
						<div class="card-header"><h2>Block Colors</h2></div>
						<div class="card-body">
							<p class="form-hint" style="margin-bottom: var(--space-4)">Tùy chỉnh màu sắc cho các block link</p>
							
							<!-- Background Color -->
							<div class="color-setting-group">
								<label class="setting-label">Background</label>
								<div class="color-presets-row">
									<!-- Color Picker Button -->
									<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
										<div class="gradient-icon"></div>
									</button>
									<input type="color" class="hidden-native-picker" value={appearance.settings.colors.primary} onchange={(e) => updateSetting('colors.primary', e.currentTarget.value)} />
									
									<!-- Preset Colors -->
									<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#000000'} style="background:#000000" title="#000000" onclick={() => updateSetting('colors.primary', '#000000')}></button>
									<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#ffffff'} style="background:#ffffff; border: 2px solid #e5e5ea" title="#ffffff" onclick={() => updateSetting('colors.primary', '#ffffff')}></button>
									<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#8e8e93'} style="background:#8e8e93" title="#8e8e93" onclick={() => updateSetting('colors.primary', '#8e8e93')}></button>
									<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#ff3b30'} style="background:#ff3b30" title="#ff3b30" onclick={() => updateSetting('colors.primary', '#ff3b30')}></button>
									<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#ff6b6b'} style="background:#ff6b6b" title="#ff6b6b" onclick={() => updateSetting('colors.primary', '#ff6b6b')}></button>
									<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#ffa5a5'} style="background:#ffa5a5" title="#ffa5a5" onclick={() => updateSetting('colors.primary', '#ffa5a5')}></button>
									<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#a855f7'} style="background:#a855f7" title="#a855f7" onclick={() => updateSetting('colors.primary', '#a855f7')}></button>
									<button class="color-preset-circle" class:active={appearance.settings.colors.primary === '#5856d6'} style="background:#5856d6" title="#5856d6" onclick={() => updateSetting('colors.primary', '#5856d6')}></button>
								</div>
							</div>

							<!-- Text Color -->
							<div class="color-setting-group">
								<label class="setting-label">Text</label>
								<div class="color-presets-row">
									<!-- Color Picker Button -->
									<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
										<div class="gradient-icon"></div>
									</button>
									<input type="color" class="hidden-native-picker" value={appearance.settings.links.textColor || getLinkColor()} onchange={(e) => updateSetting('links.textColor', e.currentTarget.value)} />
									
									<!-- Preset Colors -->
									<button class="color-preset-circle" class:active={(appearance.settings.links.textColor || getLinkColor()) === '#000000'} style="background:#000000" title="#000000" onclick={() => updateSetting('links.textColor', '#000000')}></button>
									<button class="color-preset-circle" class:active={(appearance.settings.links.textColor || getLinkColor()) === '#ffffff'} style="background:#ffffff; border: 2px solid #e5e5ea" title="#ffffff" onclick={() => updateSetting('links.textColor', '#ffffff')}></button>
									<button class="color-preset-circle" class:active={(appearance.settings.links.textColor || getLinkColor()) === '#8e8e93'} style="background:#8e8e93" title="#8e8e93" onclick={() => updateSetting('links.textColor', '#8e8e93')}></button>
									<button class="color-preset-circle" class:active={(appearance.settings.links.textColor || getLinkColor()) === '#ff3b30'} style="background:#ff3b30" title="#ff3b30" onclick={() => updateSetting('links.textColor', '#ff3b30')}></button>
									<button class="color-preset-circle" class:active={(appearance.settings.links.textColor || getLinkColor()) === '#ff6b6b'} style="background:#ff6b6b" title="#ff6b6b" onclick={() => updateSetting('links.textColor', '#ff6b6b')}></button>
									<button class="color-preset-circle" class:active={(appearance.settings.links.textColor || getLinkColor()) === '#ffa5a5'} style="background:#ffa5a5" title="#ffa5a5" onclick={() => updateSetting('links.textColor', '#ffa5a5')}></button>
									<button class="color-preset-circle" class:active={(appearance.settings.links.textColor || getLinkColor()) === '#a855f7'} style="background:#a855f7" title="#a855f7" onclick={() => updateSetting('links.textColor', '#a855f7')}></button>
									<button class="color-preset-circle" class:active={(appearance.settings.links.textColor || getLinkColor()) === '#5856d6'} style="background:#5856d6" title="#5856d6" onclick={() => updateSetting('links.textColor', '#5856d6')}></button>
								</div>
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



				{#if activeSection === 'typography'}
					<div class="card">
						<div class="card-header"><h2>Căn chỉnh Block</h2></div>
						<div class="card-body">
							<div class="setting-row"><span class="setting-label">Text trong Block</span><div class="segmented-control">{#each alignOptions as a}<button class:active={appearance.settings.links.textAlign === a.value} title={a.value} onclick={() => updateSetting('links.textAlign', a.value)}>
								{#if a.icon === 'left'}<AlignLeft size={16} />
								{:else if a.icon === 'center'}<AlignCenter size={16} />
								{:else if a.icon === 'right'}<AlignRight size={16} />
								{/if}
							</button>{/each}</div></div>
						</div>
					</div>
					<div class="card">
						<div class="card-header"><h2>Cỡ chữ</h2></div>
						<div class="card-body">
							<div class="setting-row"><span class="setting-label">Blocks</span><div class="option-group compact">{#each fontSizes as s}<button class="option-btn small" class:active={appearance.settings.links.fontSize === s.value} onclick={() => updateSetting('links.fontSize', s.value)}>{s.label}</button>{/each}</div></div>
						</div>
					</div>

				{/if}
			</main>

			<aside class="preview">
				<div class="preview-header">
					<span class="text-sm font-medium text-secondary">Preview</span>
					<button class="debug-toggle" onclick={() => showDebug = !showDebug} title="Toggle Debug"><Settings size={14} /></button>
				</div>
				
				<!-- Debug Panel -->
				{#if showDebug}
				<div class="debug-panel">
					<div class="debug-section">
						<strong>Theme:</strong> {appearance.selectedPreset?.name || 'Unknown'} (ID: {appearance.selectedPresetId})
						<button class="reset-btn" onclick={resetToPresetDefaults}><RefreshCw size={10} /> Reset</button>
					</div>
					<div class="debug-section">
						<strong>Custom Patch:</strong>
						<div class="debug-item">{Object.keys(appearance.customPatch).length === 0 ? '(empty - using pure preset)' : JSON.stringify(appearance.customPatch).substring(0, 100) + '...'}</div>
					</div>
					<div class="debug-section">
						<strong>Background:</strong>
						<div class="debug-item">Type: {appearance.settings.background.type}</div>
						<div class="debug-item">Color: {appearance.settings.background.color}</div>
						<div class="debug-item">Gradient: {appearance.settings.background.gradient || 'none'}</div>
						<div class="debug-item">Blur: {appearance.settings.background.effects.blur}px</div>
						<div class="debug-item">Dim: {appearance.settings.background.effects.dim}</div>
					</div>
					<div class="debug-section">
						<strong>Colors:</strong>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.colors.primary}"></span> Primary: {appearance.settings.colors.primary}</div>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.colors.text}"></span> Text: {appearance.settings.colors.text}</div>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.colors.textSecondary}"></span> Text Secondary: {appearance.settings.colors.textSecondary}</div>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.colors.cardBackground}"></span> Card BG: {appearance.settings.colors.cardBackground}</div>
					</div>
					<div class="debug-section">
						<strong>Links:</strong>
						<div class="debug-item">Style: {appearance.settings.links.style}</div>
						<div class="debug-item">Radius: {appearance.settings.links.borderRadius}px</div>
						<div class="debug-item">Padding: {appearance.settings.links.padding}px</div>
						<div class="debug-item">Gap: {appearance.settings.links.gap}px</div>
						<div class="debug-item">Shadow: {appearance.settings.links.shadow}</div>
					</div>
					<div class="debug-section">
						<strong>Layout:</strong>
						<div class="debug-item">Page Padding: {appearance.settings.page.layout.pagePadding}px</div>
						<div class="debug-item">Block Gap: {appearance.settings.page.layout.blockGap}px</div>
					</div>
					<div class="debug-section">
						<strong>Header:</strong>
						<div class="debug-item">Style: {appearance.settings.header.style}</div>
						<div class="debug-item">Show Avatar: {appearance.settings.header.showAvatar}</div>
						<div class="debug-item">Avatar Size: {appearance.settings.header.avatarSize}</div>
						<div class="debug-item">Avatar Shape: {appearance.settings.header.avatarShape}</div>
						<div class="debug-item">Show Bio: {appearance.settings.header.showBio}</div>
						<div class="debug-item">Bio Size: {appearance.settings.header.bioSize}</div>
						<div class="debug-item">Bio Align: {appearance.settings.header.bioAlign}</div>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.header.bioColor || appearance.settings.colors.textSecondary}"></span> Bio Color: {appearance.settings.header.bioColor || '(auto)'}</div>
					</div>
					<div class="debug-section">
						<strong>Computed Styles (Preview):</strong>
						<div class="debug-item">BG Style: {getBgStyle().substring(0, 50)}...</div>
						<div class="debug-item">Link BG: {getLinkBg()}</div>
						<div class="debug-item">Link Color: {getLinkColor()}</div>
						<div class="debug-item">Shadow: {getShadow()}</div>
					</div>
				</div>
				{/if}
				
				<div class="preview-frame">
					<div class="preview-content" style="background:{getBgStyle()}; text-align:center">
						{#if appearance.settings.background.effects.dim > 0}<div class="preview-overlay" style="background:rgba(0,0,0,{appearance.settings.background.effects.dim})"></div>{/if}
						
						{#if appearance.settings.header.style === 'classic'}
							<!-- Classic: Cover + Avatar Overlap Center -->
							<div class="preview-cover-classic" style="background: {appearance.settings.colors.primary};"></div>
							{#key `${appearance.settings.header.avatarBorderColor}-${appearance.settings.header.avatarBorderWidth}-${appearance.settings.header.avatarSize}-${appearance.settings.header.avatarShape}`}
								<div class="preview-avatar-classic" style="{getHeaderAvatarStyle()}">{auth.user?.email?.charAt(0).toUpperCase() || 'U'}</div>
							{/key}
							<div class="preview-header-classic" style="{getHeaderContainerStyle()}">
								<div class="preview-name" style="{getHeaderNameStyle()}">@{auth.user?.username || 'username'}</div>
								{#if appearance.settings.header.showSocials}
									<div class="preview-social-icons">
										<span class="preview-social-icon"><Instagram size={16} /></span>
										<span class="preview-social-icon"><Music size={16} /></span>
									</div>
								{/if}
								{#if appearance.settings.header.showBio}
									<div class="preview-bio" style="{getHeaderBioStyle()}">Your bio here</div>
								{/if}
							</div>
							
						{:else}
							<!-- Minimal: No Cover, Vertical Center Layout -->
							<div class="preview-header-minimal" style="padding: {appearance.settings.page.layout.pagePadding}px {appearance.settings.page.layout.pagePadding}px 0;">
								{#key `${appearance.settings.header.avatarSize}-${appearance.settings.header.avatarShape}`}
									<div class="preview-avatar-minimal" style="{getHeaderAvatarStyle()}">{auth.user?.email?.charAt(0).toUpperCase() || 'U'}</div>
								{/key}
								<div class="preview-name" style="{getHeaderNameStyle()}">@{auth.user?.username || 'username'}</div>
								{#if appearance.settings.header.showSocials}
									<div class="preview-social-icons">
										<span class="preview-social-icon"><Instagram size={16} /></span>
										<span class="preview-social-icon"><Music size={16} /></span>
									</div>
								{/if}
								{#if appearance.settings.header.showBio}
									<div class="preview-bio" style="{getHeaderBioStyle()}">Your bio here</div>
								{/if}
							</div>
						{/if}
						
						<div class="preview-links" style="gap:{appearance.settings.links.gap}px; padding: 0 {appearance.settings.page.layout.pagePadding}px {appearance.settings.page.layout.pagePadding}px;">
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
	.page-container { max-width: 1200px; margin: 0 auto; }
	.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: var(--space-4); }
	.page-header h1 { font-size: var(--text-xl); font-weight: 600; margin-bottom: var(--space-1); }
	.header-actions { display: flex; gap: var(--space-3); }
	.btn-save { display: flex; align-items: center; gap: var(--space-2); }
	.btn-save .spinner { width: 16px; height: 16px; border: 2px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 0.8s linear infinite; }
	.alert-banner { margin-bottom: var(--space-4); display: flex; align-items: center; gap: var(--space-2); }

	.layout { display: grid; grid-template-columns: 160px 1fr 360px; gap: var(--space-5); }
	.sidebar { display: flex; flex-direction: column; gap: var(--space-1); position: sticky; top: 80px; height: fit-content; }
	.nav-item { display: flex; align-items: center; gap: var(--space-3); padding: var(--space-3); background: transparent; border-radius: var(--radius-md); text-align: left; font-size: var(--text-sm); font-weight: 500; }
	.nav-item:hover { background: var(--color-bg); }
	.nav-item.active { background: var(--color-primary-light); color: var(--color-primary); }
	.nav-icon { display: flex; align-items: center; justify-content: center; width: 20px; }
	.content { display: flex; flex-direction: column; gap: var(--space-4); }

	.theme-card { position: relative; background: transparent; padding: 0; text-align: center; }
	.theme-card .badge { position: absolute; top: 4px; right: 4px; }

	.color-scheme-card { background: transparent; padding: var(--space-3); text-align: center; border: 2px solid transparent; border-radius: var(--radius-md); cursor: pointer; transition: all 0.2s; }
	.color-scheme-card:hover { border-color: var(--color-separator); background: var(--color-bg); }
	.color-scheme-preview { display: flex; align-items: center; justify-content: center; gap: var(--space-2); margin-bottom: var(--space-2); }
	.color-dot-large { width: 48px; height: 48px; border-radius: 50%; border: 2px solid rgba(0,0,0,0.1); }
	.color-dots-small { display: flex; flex-direction: column; gap: 4px; }
	.color-dot-small { width: 20px; height: 20px; border-radius: 4px; border: 1px solid rgba(0,0,0,0.1); }

	.advanced-toggle { background: transparent; border: none; color: var(--color-text-secondary); font-size: var(--text-sm); font-weight: 500; cursor: pointer; padding: 0; display: flex; align-items: center; gap: var(--space-2); }
	.advanced-toggle:hover { color: var(--color-text); }

	.link-style-card { background: transparent; padding: 0; text-align: center; }
	.link-style-card.active span { color: var(--color-primary); font-weight: 500; }
	.link-demo { padding: var(--space-3); border-radius: var(--radius-md); font-size: var(--text-sm); margin-bottom: var(--space-2); }
	.link-demo.filled { background: var(--color-primary); color: white; }
	.link-demo.outline { background: transparent; border: 1px solid var(--color-primary); color: var(--color-primary); }
	.link-demo.soft { background: var(--color-primary-light); color: var(--color-primary); }
	.link-demo.ghost { background: transparent; color: var(--color-text); }
	.link-style-card.active .link-demo { box-shadow: 0 0 0 2px var(--color-primary); }

	.preview { position: sticky; top: 80px; height: fit-content; }
	.preview-header { margin-bottom: var(--space-3); display: flex; justify-content: space-between; align-items: center; }
	.debug-toggle { background: var(--color-bg); border: 1px solid var(--color-separator); border-radius: var(--radius-sm); padding: 4px 8px; cursor: pointer; display: flex; align-items: center; justify-content: center; }
	.debug-toggle:hover { background: var(--color-bg-secondary); }
	
	.debug-panel { background: var(--color-bg); border: 1px solid var(--color-separator); border-radius: var(--radius-md); padding: var(--space-3); margin-bottom: var(--space-3); font-size: 11px; font-family: monospace; max-height: 400px; overflow-y: auto; }
	.debug-section { margin-bottom: var(--space-2); padding-bottom: var(--space-2); border-bottom: 1px solid var(--color-separator); }
	.debug-section:last-child { margin-bottom: 0; padding-bottom: 0; border-bottom: none; }
	.debug-section strong { display: block; margin-bottom: 4px; font-size: 12px; }
	.debug-item { color: var(--color-text-secondary); margin-left: 8px; line-height: 1.6; display: flex; align-items: center; gap: 6px; }
	.color-dot { width: 12px; height: 12px; border-radius: 2px; border: 1px solid rgba(0,0,0,0.1); flex-shrink: 0; }
	.reset-btn { margin-left: 8px; padding: 2px 8px; font-size: 10px; background: var(--color-danger); color: white; border: none; border-radius: 4px; cursor: pointer; display: inline-flex; align-items: center; gap: 4px; }
	.btn-reset-color { background: var(--color-bg-secondary); border: 1px solid var(--color-separator); border-radius: var(--radius-sm); padding: 4px 8px; cursor: pointer; color: var(--color-text-secondary); display: flex; align-items: center; justify-content: center; }
	.btn-reset-color:hover { background: var(--color-bg); color: var(--color-text); }

	/* Modern Color Settings */
	.color-setting-group {
		padding: var(--space-3);
		background: var(--color-bg);
		border-radius: var(--radius-md);
		margin-bottom: var(--space-3);
	}
	.color-setting-group:last-child { margin-bottom: 0; }
	
	.color-setting-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: var(--space-3);
	}
	
	.color-value-display {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}
	
	.color-picker-trigger {
		width: 40px;
		height: 40px;
		border: 2px solid var(--color-separator);
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: all 0.2s ease;
		padding: 0;
		background: transparent;
	}
	.color-picker-trigger:hover {
		border-color: var(--color-primary);
		transform: scale(1.05);
	}
	
	.color-hex-input {
		width: 100px;
		padding: 8px 12px;
		border: 1px solid var(--color-separator);
		border-radius: var(--radius-sm);
		font-size: var(--text-sm);
		font-family: 'SF Mono', 'Monaco', 'Courier New', monospace;
		text-transform: uppercase;
		background: var(--color-surface);
		color: var(--color-text);
		transition: all 0.2s ease;
	}
	.color-hex-input:focus {
		outline: none;
		border-color: var(--color-primary);
		box-shadow: 0 0 0 3px var(--color-primary-light);
	}
	.color-hex-input::placeholder {
		color: var(--color-text-secondary);
		text-transform: none;
	}
	
	.btn-icon-reset {
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: 1px solid var(--color-separator);
		border-radius: var(--radius-sm);
		cursor: pointer;
		color: var(--color-text-secondary);
		transition: all 0.2s ease;
	}
	.btn-icon-reset:hover {
		background: var(--color-bg-secondary);
		color: var(--color-text);
		border-color: var(--color-text-secondary);
	}
	
	.color-presets-inline {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
	}
	
	.color-preset-btn {
		width: 36px;
		height: 36px;
		border: 2px solid transparent;
		border-radius: var(--radius-sm);
		cursor: pointer;
		transition: all 0.2s ease;
		position: relative;
		padding: 0;
		background-clip: padding-box;
	}
	.color-preset-btn:hover {
		transform: scale(1.1);
		box-shadow: 0 2px 8px rgba(0,0,0,0.15);
	}
	.color-preset-btn.active {
		border-color: var(--color-primary);
		box-shadow: 0 0 0 2px var(--color-primary-light);
		transform: scale(1.05);
	}
	.color-preset-btn::after {
		content: '';
		position: absolute;
		inset: -2px;
		border-radius: var(--radius-sm);
		border: 1px solid rgba(0,0,0,0.1);
		pointer-events: none;
	}
	
	/* Color Presets Row - Circle Style */
	.color-presets-row {
		display: flex;
		align-items: center;
		gap: 12px;
		flex-wrap: wrap;
	}
	
	/* Header Style Cards */
	.header-style-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
		gap: 12px;
	}
	
	.header-style-card {
		background: transparent;
		border: 2px solid var(--color-separator);
		border-radius: 12px;
		padding: 0;
		cursor: pointer;
		transition: all 0.2s ease;
		overflow: hidden;
	}
	
	.header-style-card:hover {
		border-color: var(--color-primary);
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	}
	
	.header-style-card.active {
		border-color: var(--color-primary);
		border-width: 3px;
		box-shadow: 0 0 0 3px var(--color-primary-light);
	}
	
	.header-style-preview {
		height: 120px;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: flex-start;
		gap: 0;
		padding: 0;
		position: relative;
		overflow: hidden;
	}
	
	/* Classic Style - Cover + Avatar Overlap Center */
	.header-style-preview.classic {
		background: linear-gradient(to bottom, #c2185b 0%, #c2185b 50%, #b2dfdb 50%, #b2dfdb 100%);
	}
	
	.header-style-preview.classic::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		height: 60px;
		background: #c2185b;
		z-index: 0;
	}
	
	.header-style-preview.classic .preview-avatar-mini {
		position: absolute;
		top: 30px;
		left: 50%;
		transform: translateX(-50%);
		width: 48px;
		height: 48px;
		border: 3px solid white;
		box-shadow: 0 2px 8px rgba(0,0,0,0.2);
		z-index: 2;
	}
	
	.header-style-preview.classic .preview-text-mini {
		margin-top: 78px;
		align-items: center;
		z-index: 1;
	}
	
	.header-style-preview.classic .preview-name-mini,
	.header-style-preview.classic .preview-bio-mini {
		color: #1c1c1e;
	}
	
	.header-style-preview.classic .preview-bio-mini {
		color: #666;
	}
	
	/* Minimal Style - No Cover, Vertical Center Layout */
	.header-style-preview.minimal {
		background: #fafafa;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 6px;
		padding: 16px;
	}
	
	.header-style-preview.minimal .preview-avatar-mini {
		width: 36px;
		height: 36px;
		position: relative;
		box-shadow: 0 1px 6px rgba(0,0,0,0.08);
	}
	
	.header-style-preview.minimal .preview-text-mini {
		align-items: center;
		margin-top: 0;
	}
	
	.header-style-preview.minimal .preview-name-mini {
		color: #1c1c1e;
		font-size: 11px;
		font-weight: 700;
	}
	
	.header-style-preview.minimal .preview-bio-mini {
		color: #8e8e93;
		font-size: 8.5px;
		line-height: 1.3;
	}
	
	.preview-avatar-mini {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		background: linear-gradient(135deg, #ff9500, #ff3b30);
		border: 2px solid white;
		flex-shrink: 0;
		position: relative;
		z-index: 1;
	}
	
	.preview-text-mini {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 4px;
		position: relative;
		z-index: 1;
	}
	
	.preview-name-mini {
		font-size: 11px;
		font-weight: 600;
		color: white;
		white-space: nowrap;
	}
	
	.preview-bio-mini {
		font-size: 9px;
		color: rgba(255, 255, 255, 0.8);
		white-space: nowrap;
	}
	
	.header-style-label {
		display: block;
		padding: 8px;
		font-size: 13px;
		font-weight: 500;
		color: var(--color-text);
		text-align: center;
		background: var(--color-bg);
	}
	
	.header-style-card.active .header-style-label {
		color: var(--color-primary);
		font-weight: 600;
	}
	
	.color-preset-circle {
		width: 48px;
		height: 48px;
		border: 3px solid transparent;
		border-radius: 50%;
		cursor: pointer;
		transition: all 0.2s ease;
		position: relative;
		padding: 0;
		background-clip: padding-box;
		flex-shrink: 0;
	}
	
	.color-preset-circle:hover {
		transform: scale(1.1);
		box-shadow: 0 4px 12px rgba(0,0,0,0.2);
	}
	
	.color-preset-circle.active {
		border-color: var(--color-primary);
		box-shadow: 0 0 0 2px var(--color-primary-light);
		transform: scale(1.05);
	}
	
	.color-preset-circle.gradient-picker {
		background: conic-gradient(from 0deg, 
			#ff0000 0deg, 
			#ff9500 60deg, 
			#ffcc00 120deg, 
			#34c759 180deg, 
			#00c7be 240deg, 
			#007aff 300deg, 
			#5856d6 330deg,
			#af52de 360deg,
			#ff0000 360deg
		);
		position: relative;
		overflow: hidden;
	}
	
	.color-preset-circle.gradient-picker:hover {
		transform: scale(1.1) rotate(5deg);
	}
	
	.color-preset-circle.reset-circle {
		background: var(--color-bg);
		border: 2px solid var(--color-separator);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--color-text-secondary);
	}
	
	.color-preset-circle.reset-circle:hover {
		background: var(--color-bg-secondary);
		color: var(--color-text);
		border-color: var(--color-text-secondary);
	}
	
	.hidden-color-input {
		position: absolute;
		inset: 0;
		opacity: 0;
		cursor: pointer;
		width: 100%;
		height: 100%;
		border: none;
	}
	
	.gradient-icon {
		position: absolute;
		inset: 6px;
		background: white;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		pointer-events: none;
	}
	
	.gradient-icon::before {
		content: '';
		width: 24px;
		height: 24px;
		background: conic-gradient(from 0deg, 
			#ff0000 0deg, 
			#ff9500 60deg, 
			#ffcc00 120deg, 
			#34c759 180deg, 
			#00c7be 240deg, 
			#007aff 300deg, 
			#5856d6 330deg,
			#af52de 360deg,
			#ff0000 360deg
		);
		border-radius: 50%;
	}
	
	.preview-frame { 
		background: #1a1a1a; 
		border-radius: 32px; 
		padding: 8px; 
		box-shadow: var(--shadow-lg); 
		border: 2px solid rgba(0, 0, 0, 0.25);
	}
	.preview-content { position: relative; min-height: 560px; border-radius: 24px; overflow: hidden; display: flex; flex-direction: column; }
	.preview-overlay { position: absolute; inset: 0; pointer-events: none; z-index: 0; }
	
	/* ===== CLASSIC STYLE ===== */
	.preview-cover-classic {
		width: 100%;
		height: 160px;
		flex-shrink: 0;
	}
	
	.preview-avatar-classic {
		position: absolute;
		left: 50%;
		transform: translateX(-50%);
		background: linear-gradient(135deg, #ff9500, #ff3b30);
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 600;
		color: white;
		z-index: 10;
	}
	
	.preview-header-classic {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 8px;
		text-align: center;
	}
	
	/* ===== MINIMAL STYLE ===== */
	.preview-header-minimal {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 10px;
		text-align: center;
	}
	
	.preview-avatar-minimal {
		background: linear-gradient(135deg, #ff9500, #ff3b30);
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 600;
		color: white;
		flex-shrink: 0;
		margin-bottom: 4px;
	}
	
	.preview-text-minimal {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 6px;
		text-align: center;
	}
	
	/* ===== COMMON ELEMENTS ===== */
	.preview-name { 
		font-weight: 600;
		line-height: 1.2;
	}
	
	.preview-bio { 
		line-height: 1.4;
		opacity: 0.9;
	}
	
	.preview-social-icons { 
		display: flex; 
		gap: 8px;
		flex-wrap: wrap;
		justify-content: center;
		margin: 2px 0;
	}
	
	.preview-social-icon { 
		display: flex; 
		align-items: center; 
		justify-content: center; 
		width: 32px; 
		height: 32px; 
		border-radius: 50%; 
		background: rgba(0,0,0,0.05); 
		color: var(--color-text-secondary); 
		transition: all 0.2s ease;
		flex-shrink: 0;
	}
	
	.preview-social-icon:hover { 
		background: rgba(0,0,0,0.1); 
	}
	
	.preview-links { 
		display: flex; 
		flex-direction: column;
		margin-top: 16px;
	}
	
	.preview-link { 
		font-size: 0.9375rem; 
		font-weight: 500; 
	}
	
	/* Sub-tabs Styling */
	.sub-tabs {
		display: flex;
		gap: 8px;
		margin-bottom: var(--space-4);
		padding-bottom: var(--space-3);
		border-bottom: 1px solid var(--color-separator);
	}
	
	.sub-tab {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 8px 16px;
		background: transparent;
		border: none;
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-weight: 500;
		color: var(--color-text-secondary);
		cursor: pointer;
		transition: all 0.2s ease;
	}
	
	.sub-tab:hover {
		background: var(--color-bg);
		color: var(--color-text);
	}
	
	.sub-tab.active {
		background: var(--color-primary-light);
		color: var(--color-primary);
	}
	
	.sub-tab-content {
		animation: fadeIn 0.2s ease;
	}
	
	@keyframes fadeIn {
		from { opacity: 0; transform: translateY(-4px); }
		to { opacity: 1; transform: translateY(0); }
	}

	@media (max-width: 1024px) {
		.layout { grid-template-columns: 1fr; }
		.sidebar { flex-direction: row; overflow-x: auto; position: static; padding-bottom: var(--space-3); border-bottom: 1px solid var(--color-separator); }
		.nav-item { flex-shrink: 0; }
		.preview { display: none; }
		.page-header { flex-direction: column; align-items: flex-start; gap: var(--space-3); }
		.header-actions { width: 100%; justify-content: flex-end; }
		.sub-tabs { overflow-x: auto; }
		.sub-tab { flex-shrink: 0; }
	}
	
	/* Hidden native color picker */
	.hidden-native-picker {
		position: absolute;
		width: 0;
		height: 0;
		opacity: 0;
		pointer-events: none;
	}
</style>


