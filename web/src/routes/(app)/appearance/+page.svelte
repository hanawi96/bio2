<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getAuth } from '$lib/stores/auth.svelte';
	import { getAppearance, loadAppearance, selectPreset, updateSetting, saveAppearance, resetAppearance, resetToPresetDefaults } from '$lib/stores/appearance.svelte';
	import { bio } from '$lib/api/client';
	import { Palette, Image, User, Link, Type, Save, AlertTriangle, Sun, Moon, Camera, AlignLeft, AlignCenter, AlignRight, Settings, RefreshCw, Instagram, Music, Facebook, Twitter, Youtube, Linkedin, Github, Globe } from 'lucide-svelte';

	const auth = getAuth();
	const appearance = getAppearance();
	let pageId = $state<number | null>(null);
	let activeSection = $state<string>('theme');
	let showDebug = $state(false);

	const sections = [
		{ id: 'theme', label: 'Theme', icon: 'palette' },
		{ id: 'background', label: 'Background', icon: 'image' },
		{ id: 'header', label: 'Header', icon: 'user' },
		{ id: 'links', label: 'Links', icon: 'link' },
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
		const cardBg = appearance.settings.colors.cardBackground;
		
		switch (style) {
			case 'minimal':
				return 'padding: 8px 0; gap: 6px;';
			case 'centered':
				return 'padding: 24px 0; gap: 16px;';
			case 'card':
				return `padding: 20px; gap: 12px; background: ${cardBg}; border-radius: 16px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); margin: 0 8px;`;
			default:
				return 'padding: 16px 0; gap: 12px;';
		}
	}

	function getHeaderAvatarStyle(): string {
		const style = appearance.settings.header.style;
		const baseSize = getAvatarSize();
		const radius = getAvatarRadius();
		
		let size = baseSize;
		if (style === 'minimal') size = Math.round(baseSize * 0.75);
		if (style === 'centered') size = Math.round(baseSize * 1.2);
		
		return `width: ${size}px; height: ${size}px; border-radius: ${radius};`;
	}

	function getHeaderNameStyle(): string {
		const style = appearance.settings.header.style;
		const color = appearance.settings.colors.text;
		
		let fontSize = '1.125rem';
		let fontWeight = '600';
		
		if (style === 'minimal') { fontSize = '1rem'; }
		if (style === 'centered') { fontSize = '1.5rem'; fontWeight = '700'; }
		if (style === 'card') { fontSize = '1.25rem'; }
		
		return `color: ${color}; font-size: ${fontSize}; font-weight: ${fontWeight};`;
	}

	function getHeaderBioStyle(): string {
		const style = appearance.settings.header.style;
		const bioSettings = appearance.settings.header;
		const color = bioSettings.bioColor || appearance.settings.colors.textSecondary;
		
		// Map bio size to font-size
		const sizeMap: Record<string, string> = { S: '0.75rem', M: '0.875rem', L: '1rem', XL: '1.125rem' };
		let fontSize = sizeMap[bioSettings.bioSize] || '0.875rem';
		
		// Adjust based on header style
		if (style === 'minimal') {
			const minimalSizeMap: Record<string, string> = { S: '0.625rem', M: '0.75rem', L: '0.875rem', XL: '1rem' };
			fontSize = minimalSizeMap[bioSettings.bioSize] || '0.75rem';
		}
		if (style === 'centered') {
			const centeredSizeMap: Record<string, string> = { S: '0.875rem', M: '1rem', L: '1.125rem', XL: '1.25rem' };
			fontSize = centeredSizeMap[bioSettings.bioSize] || '1rem';
		}
		
		const maxWidth = style === 'minimal' ? '240px' : style === 'centered' ? '320px' : '280px';
		
		// Map align to align-self for flexbox
		const alignSelfMap: Record<string, string> = {
			left: 'flex-start',
			center: 'center',
			right: 'flex-end'
		};
		const alignSelf = alignSelfMap[bioSettings.bioAlign] || 'center';
		
		return `color: ${color}; font-size: ${fontSize}; max-width: ${maxWidth}; text-align: ${bioSettings.bioAlign}; align-self: ${alignSelf};`;
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
					<!-- Profile Info -->
					<div class="card">
						<div class="card-header"><h2>Thông tin cá nhân</h2></div>
						<div class="card-body">
							<div class="form-group">
								<label class="form-label">Avatar</label>
								<div class="avatar-upload">
									<div class="avatar-upload-preview">
										{#if auth.user?.avatar_asset_id}
											<img src="/api/assets/{auth.user.avatar_asset_id}" alt="Avatar" />
										{:else}
											{auth.user?.email?.charAt(0).toUpperCase() || 'U'}
										{/if}
									</div>
									<label class="avatar-upload-btn" title="Upload avatar">
										<Camera size={16} />
										<input type="file" accept="image/*" />
									</label>
								</div>
								<p class="form-hint">Nhấn vào biểu tượng camera để thay đổi avatar</p>
							</div>
							<div class="form-group">
								<label class="form-label">Tên hiển thị</label>
								<input type="text" class="profile-input" placeholder="Nhập tên của bạn" value={auth.user?.display_name || ''} />
							</div>
							<div class="form-group">
								<label class="form-label">Username</label>
								<div class="input-with-prefix">
									<span class="input-prefix">@</span>
									<input type="text" placeholder="username" value={auth.user?.username || ''} readonly />
								</div>
								<p class="form-hint">Username không thể thay đổi tại đây</p>
							</div>
							<div class="form-group">
								<label class="form-label">Bio</label>
								<textarea class="profile-textarea" placeholder="Giới thiệu ngắn về bạn..." maxlength="150"></textarea>
								<p class="char-count">0/150</p>
							</div>
						</div>
					</div>

					<!-- Header Style -->
					<div class="card">
						<div class="card-header"><h2>Kiểu Header</h2></div>
						<div class="card-body">
							<div class="option-group">
								{#each headerStyles as s}
									<button class="option-btn" class:active={appearance.settings.header.style === s.value} onclick={() => updateSetting('header.style', s.value)}>{s.label}</button>
								{/each}
							</div>
						</div>
					</div>

					<!-- Avatar Settings -->
					<div class="card">
						<div class="card-header"><h2>Cài đặt Avatar</h2></div>
						<div class="card-body">
							<div class="toggle-row">
								<span>Hiển thị Avatar</span>
								<button class="toggle" class:active={appearance.settings.header.showAvatar} aria-label="Toggle" onclick={() => updateSetting('header.showAvatar', !appearance.settings.header.showAvatar)}></button>
							</div>
							{#if appearance.settings.header.showAvatar}
								<div class="setting-row">
									<span class="setting-label">Kích thước</span>
									<div class="option-group compact">
										{#each avatarSizes as s}
											<button class="option-btn small" class:active={appearance.settings.header.avatarSize === s.value} onclick={() => updateSetting('header.avatarSize', s.value)}>{s.label}</button>
										{/each}
									</div>
								</div>
								<div class="setting-row">
									<span class="setting-label">Hình dạng</span>
									<div class="option-group compact">
										{#each avatarShapes as s}
											<button class="option-btn small" class:active={appearance.settings.header.avatarShape === s.value} onclick={() => updateSetting('header.avatarShape', s.value)}>{s.label}</button>
										{/each}
									</div>
								</div>
							{/if}
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
									<div class="color-picker-row">
										<input type="color" class="color-input small" value={appearance.settings.header.bioColor || appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('header.bioColor', e.currentTarget.value)}/>
										<input type="text" class="color-text small" value={appearance.settings.header.bioColor || appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('header.bioColor', e.currentTarget.value)} placeholder="Auto"/>
										<button class="btn-reset-color" onclick={() => updateSetting('header.bioColor', '')} title="Reset về mặc định"><RefreshCw size={12} /></button>
									</div>
								</div>
							{/if}
						</div>
					</div>

					<!-- Social Icons -->
					<div class="card">
						<div class="card-header"><h2>Social Icons</h2></div>
						<div class="card-body">
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
						</div>
					</div>
				{/if}

				{#if activeSection === 'links'}
					<div class="card">
						<div class="card-header"><h2>Link Background Color</h2></div>
						<div class="card-body">
							<div class="color-picker-row"><input type="color" class="color-input" value={appearance.settings.colors.primary} onchange={(e) => updateSetting('colors.primary', e.currentTarget.value)}/><input type="text" class="color-text" value={appearance.settings.colors.primary} onchange={(e) => updateSetting('colors.primary', e.currentTarget.value)}/></div>
							<div class="color-presets">{#each presetColors as c}<button class="color-preset" style="background:{c}" class:active={appearance.settings.colors.primary === c} title={c} onclick={() => updateSetting('colors.primary', c)}></button>{/each}</div>
						</div>
					</div>
					<div class="card">
						<div class="card-header"><h2>Link Text Color</h2></div>
						<div class="card-body">
							<div class="color-picker-row">
								<input type="color" class="color-input" value={appearance.settings.links.textColor || getLinkColor()} onchange={(e) => updateSetting('links.textColor', e.currentTarget.value)}/>
								<input type="text" class="color-text" value={appearance.settings.links.textColor || getLinkColor()} onchange={(e) => updateSetting('links.textColor', e.currentTarget.value)} placeholder="Auto"/>
								<button class="btn-reset-color" onclick={() => updateSetting('links.textColor', '')} title="Reset về mặc định"><RefreshCw size={12} /></button>
							</div>
							<div class="color-presets">{#each presetColors as c}<button class="color-preset" style="background:{c}" class:active={appearance.settings.links.textColor === c} title={c} onclick={() => updateSetting('links.textColor', c)}></button>{/each}</div>
						</div>
					</div>
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



				{#if activeSection === 'typography'}
					<div class="card">
						<div class="card-header"><h2>Màu chữ</h2></div>
						<div class="card-body">
							<div class="setting-row"><span class="setting-label">Primary</span><div class="color-picker-row"><input type="color" class="color-input small" value={appearance.settings.colors.text} onchange={(e) => updateSetting('colors.text', e.currentTarget.value)}/><input type="text" class="color-text small" value={appearance.settings.colors.text} onchange={(e) => updateSetting('colors.text', e.currentTarget.value)}/></div></div>
							<div class="setting-row"><span class="setting-label">Secondary</span><div class="color-picker-row"><input type="color" class="color-input small" value={appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('colors.textSecondary', e.currentTarget.value)}/><input type="text" class="color-text small" value={appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('colors.textSecondary', e.currentTarget.value)}/></div></div>
						</div>
					</div>
					<div class="card">
						<div class="card-header"><h2>Căn chỉnh Link</h2></div>
						<div class="card-body">
							<div class="setting-row"><span class="setting-label">Text trong Link</span><div class="segmented-control">{#each alignOptions as a}<button class:active={appearance.settings.links.textAlign === a.value} title={a.value} onclick={() => updateSetting('links.textAlign', a.value)}>
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
							<div class="setting-row"><span class="setting-label">Links</span><div class="option-group compact">{#each fontSizes as s}<button class="option-btn small" class:active={appearance.settings.links.fontSize === s.value} onclick={() => updateSetting('links.fontSize', s.value)}>{s.label}</button>{/each}</div></div>
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
					<div class="preview-content" style="background:{getBgStyle()}; padding:{appearance.settings.page.layout.pagePadding}px; text-align:center">
						{#if appearance.settings.background.effects.dim > 0}<div class="preview-overlay" style="background:rgba(0,0,0,{appearance.settings.background.effects.dim})"></div>{/if}
						
						<!-- Header Section with dynamic style -->
						<div class="preview-header-section" style="{getHeaderContainerStyle()}">
							{#if appearance.settings.header.showAvatar}
								<div class="preview-avatar" style="{getHeaderAvatarStyle()}">{auth.user?.email?.charAt(0).toUpperCase() || 'U'}</div>
							{/if}
							<div class="preview-name" style="{getHeaderNameStyle()}">@{auth.user?.username || 'username'}</div>
							{#if appearance.settings.header.showBio}
								<div class="preview-bio" style="{getHeaderBioStyle()}">Your bio here</div>
							{/if}
							<!-- Social Icons Preview -->
							{#if appearance.settings.header.showSocials}
							<div class="preview-social-icons">
								<span class="preview-social-icon"><Instagram size={16} /></span>
								<span class="preview-social-icon"><Music size={16} /></span>
							</div>
							{/if}
						</div>
						
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
	.page-container { max-width: 1200px; margin: 0 auto; }
	.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: var(--space-4); }
	.page-header h1 { font-size: var(--text-xl); font-weight: 600; margin-bottom: var(--space-1); }
	.header-actions { display: flex; gap: var(--space-3); }
	.btn-save { display: flex; align-items: center; gap: var(--space-2); }
	.btn-save .spinner { width: 16px; height: 16px; border: 2px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 0.8s linear infinite; }
	.alert-banner { margin-bottom: var(--space-4); display: flex; align-items: center; gap: var(--space-2); }

	.layout { display: grid; grid-template-columns: 160px 1fr 300px; gap: var(--space-5); }
	.sidebar { display: flex; flex-direction: column; gap: var(--space-1); position: sticky; top: 80px; height: fit-content; }
	.nav-item { display: flex; align-items: center; gap: var(--space-3); padding: var(--space-3); background: transparent; border-radius: var(--radius-md); text-align: left; font-size: var(--text-sm); font-weight: 500; }
	.nav-item:hover { background: var(--color-bg); }
	.nav-item.active { background: var(--color-primary-light); color: var(--color-primary); }
	.nav-icon { display: flex; align-items: center; justify-content: center; width: 20px; }
	.content { display: flex; flex-direction: column; gap: var(--space-4); }

	.theme-card { position: relative; background: transparent; padding: 0; text-align: center; }
	.theme-card .badge { position: absolute; top: 4px; right: 4px; }

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
	
	.preview-frame { background: var(--color-bg-secondary); border-radius: var(--radius-xl); padding: var(--space-3); box-shadow: var(--shadow-lg); }
	.preview-content { position: relative; min-height: 480px; border-radius: var(--radius-lg); overflow: hidden; display: flex; flex-direction: column; align-items: center; padding-top: var(--space-4); }
	.preview-overlay { position: absolute; inset: 0; pointer-events: none; }
	.preview-header-section { position: relative; z-index: 1; display: flex; flex-direction: column; align-items: center; width: 100%; transition: all 0.2s ease; }
	.preview-avatar { background: linear-gradient(135deg, #ff9500, #ff3b30); display: flex; align-items: center; justify-content: center; font-size: 1.5rem; font-weight: 600; color: white; transition: all 0.2s ease; }
	.preview-name { transition: all 0.2s ease; }
	.preview-bio { transition: all 0.2s ease; }
	.preview-social-icons { display: flex; gap: var(--space-2); margin-top: var(--space-2); }
	.preview-social-icon { display: flex; align-items: center; justify-content: center; width: 32px; height: 32px; border-radius: 50%; background: rgba(0,0,0,0.05); color: var(--color-text-secondary); transition: all 0.2s ease; }
	.preview-social-icon:hover { background: rgba(0,0,0,0.1); }
	.preview-links { position: relative; z-index: 1; width: 100%; display: flex; flex-direction: column; }
	.preview-link { font-size: var(--text-sm); font-weight: 500; }

	@media (max-width: 1024px) {
		.layout { grid-template-columns: 1fr; }
		.sidebar { flex-direction: row; overflow-x: auto; position: static; padding-bottom: var(--space-3); border-bottom: 1px solid var(--color-separator); }
		.nav-item { flex-shrink: 0; }
		.preview { display: none; }
		.page-header { flex-direction: column; align-items: flex-start; gap: var(--space-3); }
		.header-actions { width: 100%; justify-content: flex-end; }
	}
</style>
