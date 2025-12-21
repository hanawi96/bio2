<script lang="ts">
	import { onMount, untrack } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getAuth } from '$lib/stores/auth.svelte';
	import { getAppearance, loadAppearance, selectPreset, selectCustomTheme, updateSetting, saveAppearance, resetAppearance, resetToPresetDefaults, deleteCustomTheme, buildGradientString } from '$lib/stores/appearance.svelte';
	import { bio } from '$lib/api/client';
	import { Palette, Image, User, Link, Droplets, Type, Save, AlertTriangle, Sun, Moon, AlignLeft, AlignCenter, AlignRight, Settings, RefreshCw, Instagram, Music, Facebook, Twitter, Youtube, Linkedin, Github, Globe, X } from 'lucide-svelte';

	const auth = getAuth();
	const appearance = getAppearance();
	let pageId = $state<number | null>(null);
	let bioData = $state<any>(null);
	let activeSection = $state<string>('theme');
	let activeHeaderTab = $state<string>('avatar'); // Sub-tab for header settings
	let activeBlockTab = $state<string>('basic'); // Sub-tab for block settings: 'basic' or 'advanced'
	let showDebug = $state(false);
	
	// Custom Gradient Builder State
	let showCustomGradientBuilder = $state(false);
	let customGradientDraft = $state({
		type: 'linear' as 'linear' | 'radial',
		angle: 135,
		fromColor: '#667eea',
		toColor: '#764ba2'
	});
	
	// Derived state for live preview - auto updates when draft changes
	let liveGradientPreview = $derived(buildGradientString(customGradientDraft));
	
	// Track previous values to avoid unnecessary updates
	let prevGradientString = '';
	
	// Auto-apply custom gradient when draft changes
	$effect(() => {
		// Track draft changes
		const draft = customGradientDraft;
		const isOpen = showCustomGradientBuilder;
		
		if (isOpen) {
			// Validate colors
			const hexColorRegex = /^#[0-9A-Fa-f]{6}$/;
			const validColors = hexColorRegex.test(draft.fromColor) && 
			                    hexColorRegex.test(draft.toColor);
			
			if (validColors) {
				const gradientString = buildGradientString(draft);
				
				// Only update if gradient actually changed (avoid infinite loop)
				if (gradientString !== prevGradientString) {
					prevGradientString = gradientString;
					
					updateSetting('background.customGradient', {
						type: draft.type,
						angle: draft.angle,
						fromColor: draft.fromColor,
						toColor: draft.toColor,
						enabled: true
					});
					updateSetting('background.gradient', gradientString);
				}
			}
		}
	});

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

		const data = await bio.get();
		bioData = data;
		if (data?.page) {
			pageId = data.page.id;
			await loadAppearance(data.page.id);
		}
	});

	async function handleSave() { await saveAppearance(); }

	// Handle cover image upload
	async function handleCoverUpload(e: Event) {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;

		// Validate file size (max 2MB)
		if (file.size > 2 * 1024 * 1024) {
			alert('File quá lớn! Vui lòng chọn ảnh dưới 2MB.');
			return;
		}

		// Validate file type
		if (!file.type.startsWith('image/')) {
			alert('Vui lòng chọn file ảnh!');
			return;
		}

		// TODO: Upload to server and get URL
		// For now, use local preview
		const reader = new FileReader();
		reader.onload = (e) => {
			const dataUrl = e.target?.result as string;
			updateSetting('header.cover.imageUrl', dataUrl);
			// TODO: Upload to server and save asset ID
			// updateSetting('header.cover.imageAssetId', assetId);
		};
		reader.readAsDataURL(file);
	}

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
	const linkStyles = [{ value: 'filled', label: 'Filled' }, { value: 'outline', label: 'Outline' }, { value: 'elevated', label: 'Elevated' }, { value: 'shadow', label: 'Shadow' }];
	const headerStyles = [{ value: 'classic', label: 'Classic' }, { value: 'minimal', label: 'Minimal' }];
	const avatarSizes = [{ value: 'S', label: 'S' }, { value: 'M', label: 'M' }, { value: 'L', label: 'L' }];
	const avatarShapes = [{ value: 'circle', label: 'Circle' }, { value: 'rounded', label: 'Rounded' }, { value: 'square', label: 'Square' }];
	const backgroundTypes = [{ value: 'solid', label: 'Solid' }, { value: 'gradient', label: 'Gradient' }, { value: 'image', label: 'Image' }];
	const presetColors = ['#007aff', '#5856d6', '#af52de', '#ff2d55', '#ff3b30', '#ff9500', '#ffcc00', '#34c759', '#00c7be', '#30b0c7', '#1c1c1e', '#8e8e93'];
	
	const colorSchemes = [
		{ name: 'Neutral', primary: '#1c1c1e', bg: '#f2f2f7', text: '#1c1c1e', textSecondary: '#8e8e93' },
		{ name: 'Bold', primary: '#ff3b30', bg: '#ffffff', text: '#1c1c1e', textSecondary: '#8e8e93' },
		{ name: 'Ocean', primary: '#007aff', bg: '#e3f2fd', text: '#1c1c1e', textSecondary: '#f0f9ff' },
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

	// Helper: normalize color to hex format for comparison
	// Converts rgba(255, 192, 203, 0.75) → #ffc0cb
	// Keeps hex colors as-is
	function normalizeColor(color: string | undefined): string {
		if (!color) return '#000000';
		
		// Already hex format
		if (color.startsWith('#')) {
			return color.toLowerCase();
		}
		
		// Parse rgba/rgb format
		const rgbaMatch = color.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)(?:,\s*[\d.]+)?\)/);
		if (rgbaMatch) {
			const r = parseInt(rgbaMatch[1]);
			const g = parseInt(rgbaMatch[2]);
			const b = parseInt(rgbaMatch[3]);
			return `#${r.toString(16).padStart(2, '0')}${g.toString(16).padStart(2, '0')}${b.toString(16).padStart(2, '0')}`;
		}
		
		// Fallback
		return color.toLowerCase();
	}

	// Helper: get actual block background color (same logic as getLinkBg but return color value)
	function getActualBlockBgColor(): string {
		if (!appearance.settings.links.showBackground) {
			return 'transparent';
		}
		
		const s = appearance.settings.links.style;
		const colors = appearance.settings.colors;
		
		// Return the actual color based on style
		if (s === 'filled') return colors.primary;
		if (s === 'soft') return colors.primary; // soft uses primary with opacity, but we show primary
		return colors.cardBackground;
	}

	// Helper: get actual block text color (same logic as getLinkColor but return color value)
	function getActualBlockTextColor(): string {
		// Nếu có custom textColor, dùng nó
		if (appearance.settings.links.textColor) {
			return appearance.settings.links.textColor;
		}
		
		// Auto based on style
		const s = appearance.settings.links.style;
		if (s === 'filled') return '#ffffff';
		if (s === 'soft' || s === 'outline') return appearance.settings.colors.primary;
		return appearance.settings.colors.text;
	}

	// Lấy theme preview styles - đồng bộ 100% với theme config
	function getThemePreviewStyles(preset: any) {
		console.log(`\n=== [${preset?.name}] Theme Preview Styles ===`);
		
		const config = preset?.config;
		if (!config) {
			console.log(`[${preset?.name}] No config found, using defaults`);
			return {
				background: '#f2f2f7',
				textColor: '#1c1c1e',
				textMuted: '#8e8e93',
				cardBg: 'rgba(255, 255, 255, 0.85)',
				cardRadius: '8px',
				cardShadow: '0 2px 6px rgba(0, 0, 0, 0.12)',
				cardBorder: 'none',
				avatarBg: 'linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.7))',
				avatarBorder: '2px solid rgba(255, 255, 255, 0.95)',
				socialBg: 'rgba(255, 255, 255, 0.75)'
			};
		}
		
		console.log(`[${preset.name}] Full config:`, config);
		
		// 1. Background
		let background = '#f2f2f7';
		if (config.background?.gradient) {
			background = config.background.gradient;
			console.log(`[${preset.name}] Background: gradient =`, background);
		} else if (config.background?.color) {
			background = config.background.color;
			console.log(`[${preset.name}] Background: color =`, background);
		}
		
		// 2. Text colors từ semantic
		const textColor = config.semantic?.color?.text?.default || '#1c1c1e';
		const textMuted = config.semantic?.color?.text?.muted || '#8e8e93';
		console.log(`[${preset.name}] Text colors:`, { textColor, textMuted });
		
		// 3. Card background từ recipes hoặc semantic - PRIORITY ORDER
		let cardBg = 'rgba(255, 255, 255, 0.85)'; // fallback
		
		console.log(`[${preset.name}] Checking recipes.linkItem.base:`, config.recipes?.linkItem?.base);
		
		// Ưu tiên 1: recipes.linkItem.base.background (chính xác nhất)
		if (config.recipes?.linkItem?.base?.background) {
			cardBg = config.recipes.linkItem.base.background;
			console.log(`[${preset.name}] ✅ Using recipes.linkItem.base.background:`, cardBg);
		} 
		// Ưu tiên 2: semantic.color.surface.card (backup)
		else if (config.semantic?.color?.surface?.card) {
			cardBg = config.semantic.color.surface.card;
			console.log(`[${preset.name}] ✅ Using semantic.color.surface.card:`, cardBg);
		} else {
			console.log(`[${preset.name}] ⚠️ Using fallback cardBg:`, cardBg);
		}
		
		// 4. Card border radius từ recipes
		let cardRadius = '8px';
		if (config.recipes?.linkItem?.base?.borderRadius) {
			const radius = config.recipes.linkItem.base.borderRadius;
			cardRadius = typeof radius === 'number' ? `${radius}px` : radius;
			console.log(`[${preset.name}] Card radius:`, cardRadius);
		}
		
		// 5. Card shadow từ recipes
		let cardShadow = '0 2px 6px rgba(0, 0, 0, 0.12)';
		if (config.recipes?.linkItem?.base?.shadow) {
			cardShadow = config.recipes.linkItem.base.shadow;
			console.log(`[${preset.name}] Card shadow:`, cardShadow);
		}
		
		// 6. Card border từ recipes (nếu có)
		let cardBorder = 'none';
		if (config.recipes?.linkItem?.base?.border) {
			cardBorder = config.recipes.linkItem.base.border;
			console.log(`[${preset.name}] Card border:`, cardBorder);
		}
		
		// 7. Avatar & social colors - tự động dựa trên background
		const isDark = background.includes('#1c1c1e') || background.includes('#0f172a') || background.includes('#2c2c2e');
		const avatarBg = isDark 
			? 'linear-gradient(135deg, rgba(255, 255, 255, 0.2), rgba(255, 255, 255, 0.1))'
			: 'linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.7))';
		const avatarBorder = isDark
			? '2px solid rgba(255, 255, 255, 0.3)'
			: '2px solid rgba(255, 255, 255, 0.95)';
		const socialBg = isDark
			? 'rgba(255, 255, 255, 0.15)'
			: 'rgba(255, 255, 255, 0.75)';
		
		const result = {
			background,
			textColor,
			textMuted,
			cardBg,
			cardRadius,
			cardShadow,
			cardBorder,
			avatarBg,
			avatarBorder,
			socialBg
		};
		
		console.log(`[${preset.name}] 🎨 Final styles:`, result);
		console.log(`=== End [${preset.name}] ===\n`);
		
		return result;
	}
	
	// Legacy functions for backward compatibility
	function getPresetBackground(preset: any): string {
		return getThemePreviewStyles(preset).background;
	}
	
	function getPresetTextColor(preset: any): string {
		return getThemePreviewStyles(preset).textColor;
	}
	
	function getPresetCardColor(preset: any): string {
		return getThemePreviewStyles(preset).cardBg;
	}
	
	function getPresetColor(preset: any): string {
		return getPresetBackground(preset);
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
		// Shadow style có shadow đen cố định
		if (appearance.settings.links.style === 'shadow') {
			return '4px 4px 0 rgba(0, 0, 0, 0.9)';
		}
		
		// Elevated style có shadow màu primary
		if (appearance.settings.links.style === 'elevated') {
			return `4px 4px 0 ${appearance.settings.colors.primary}`;
		}
		
		// Nếu bật custom shadow, dùng custom values
		if (appearance.settings.links.showShadow) {
			const blur = appearance.settings.links.shadowBlur ?? 8;
			const offsetX = appearance.settings.links.shadowOffsetX ?? 0;
			const offsetY = appearance.settings.links.shadowOffsetY ?? 2;
			const color = appearance.settings.links.shadowColor || '#000000';
			const opacity = appearance.settings.links.shadowOpacity ?? 0.12;
			
			// Convert hex color to rgba with opacity
			const r = parseInt(color.slice(1, 3), 16);
			const g = parseInt(color.slice(3, 5), 16);
			const b = parseInt(color.slice(5, 7), 16);
			
			return `${offsetX}px ${offsetY}px ${blur}px rgba(${r}, ${g}, ${b}, ${opacity})`;
		}
		
		return 'none';
	}

	function getLinkBg(): string {
		if (!appearance.settings.links.showBackground) {
			return 'transparent';
		}
		
		const s = appearance.settings.links.style;
		const colors = appearance.settings.colors;
		if (s === 'filled') return colors.primary;
		if (s === 'elevated') return 'transparent';
		if (s === 'shadow') return colors.primary;
		if (s === 'outline') return 'transparent';
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
		if (s === 'shadow') return '#fff';
		if (s === 'elevated' || s === 'outline') return appearance.settings.colors.primary;
		return appearance.settings.colors.text;
	}

	// Generate style for block demo in style selector
	// Helper: Get text color for a specific style type (used by block demo)
	function getTextColorForStyle(styleType: string): string {
		// Check custom textColor first
		if (appearance.settings.links.textColor) {
			return appearance.settings.links.textColor;
		}
		
		// Auto based on style type - match getLinkColor() logic
		if (styleType === 'filled' || styleType === 'shadow') {
			// Dark background (primary color) → white text for contrast
			return '#fff';
		}
		if (styleType === 'elevated' || styleType === 'outline') {
			// Transparent background → use primary color for text
			return appearance.settings.colors.primary;
		}
		return appearance.settings.colors.text;
	}

	// Generate style for block demo in style selector  
	function getBlockDemoStyle(styleType: string): string {
		const colors = appearance.settings.colors;
		const styles: string[] = [];
		const textColor = getTextColorForStyle(styleType);
		
		if (styleType === 'filled') {
			styles.push(`background: ${colors.primary}`);
			styles.push(`color: ${textColor}`);
		} else if (styleType === 'outline') {
			styles.push('background: transparent');
			styles.push(`border: 2px solid ${colors.primary}`);
			styles.push(`color: ${textColor}`);
		} else if (styleType === 'elevated') {
			styles.push('background: transparent');
			styles.push(`border: 3px solid ${colors.primary}`);
			styles.push(`color: ${textColor}`);
			styles.push(`box-shadow: 4px 4px 0 ${colors.primary}`);
		} else if (styleType === 'shadow') {
			styles.push(`background: ${colors.primary}`);
			styles.push(`color: ${textColor}`);
			styles.push('box-shadow: 4px 4px 0 rgba(0, 0, 0, 0.9)');
		}
		
		return styles.join('; ');
	}

	function getLinkBorder(): string {
		if (appearance.settings.links.showBorder) {
			const width = appearance.settings.links.borderWidth || 1;
			const color = appearance.settings.links.borderColor || '#e5e5ea';
			return `${width}px solid ${color}`;
		}
		
		if (appearance.settings.links.style === 'elevated') {
			return `3px solid ${appearance.settings.colors.primary}`;
		}
		
		if (appearance.settings.links.style === 'outline') {
			return `1px solid ${appearance.settings.colors.primary}`;
		}
		
		return 'none';
	}

	// Header style configurations
	function getHeaderContainerStyle(): string {
		const style = appearance.settings.header.style;
		const alignItems = 'center'; // Always center for consistency
		
		switch (style) {
			case 'classic':
				// Avatar overlap cover, content below
				// Calculate padding based on avatar size
				const baseSize = getAvatarSize();
				const avatarSize = Math.round(baseSize * 1.3);
				const topPadding = Math.round(avatarSize / 2) + 12; // Half avatar + gap
				return `padding: ${topPadding}px 16px 16px 16px; gap: 12px; align-items: ${alignItems}; margin-top: 0;`;
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
		const borderColor = appearance.settings.header.avatarBorderColor || '#ffffff';
		const borderWidth = appearance.settings.header.avatarBorderWidth || 3;
		
		let size = baseSize;
		if (style === 'minimal') size = Math.round(baseSize * 1.2); // Larger for minimal
		if (style === 'classic') size = Math.round(baseSize * 1.3);
		
		// Classic style: avatar overlap cover (centered) - positioned at bottom of cover
		if (style === 'classic') {
			const coverHeight = 160;
			const topPosition = coverHeight - (size / 2);
			return `width: ${size}px; height: ${size}px; border-radius: ${radius}; border: ${borderWidth}px solid ${borderColor}; top: ${topPosition}px; box-shadow: 0 4px 12px rgba(0,0,0,0.15);`;
		}
		
		// Minimal style: clean avatar with border and subtle shadow
		if (style === 'minimal') {
			return `width: ${size}px; height: ${size}px; border-radius: ${radius}; border: ${borderWidth}px solid ${borderColor}; box-shadow: 0 2px 12px rgba(0,0,0,0.08);`;
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
		
		// Always center align for consistency
		const alignSelf = 'center';
		const textAlign = 'center';
		
		return `color: ${textColor}; font-size: ${fontSize}; max-width: ${maxWidth}; text-align: ${textAlign}; align-self: ${alignSelf}; line-height: 1.5;`;
	}

	function getCoverStyle(): string {
		const cover = appearance.settings.header.cover;
		if (cover.type === 'image' && cover.imageUrl) {
			return `background-image: url('${cover.imageUrl}'); background-size: cover; background-position: center;`;
		}
		return `background: ${cover.color};`;
	}

	function getSocialIconStyle(): string {
		const iconColor = appearance.settings.header.socialIconsColor || appearance.settings.colors.textSecondary;
		const showBg = appearance.settings.header.socialIconsBg ?? true;
		
		if (showBg) {
			return `background: rgba(0,0,0,0.05); color: ${iconColor};`;
		}
		return `background: transparent; color: ${iconColor};`;
	}

	// Helper: parse gradient string to extract type, angle, and colors
	// Example: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)" → { type: 'linear', angle: 135, fromColor: '#667eea', toColor: '#764ba2' }
	// Also handles 3+ colors: "linear-gradient(135deg, #fce7f3 0%, #dbeafe 50%, #fef3c7 100%)" → takes first and last
	function parseGradientString(gradientStr: string): { type: 'linear' | 'radial'; angle: number; fromColor: string; toColor: string } | null {
		console.log('[parseGradientString] Input:', gradientStr);
		
		if (!gradientStr) {
			console.log('[parseGradientString] Empty string, returning null');
			return null;
		}
		
		// Parse linear gradient - match angle and all color stops
		const linearMatch = gradientStr.match(/linear-gradient\((\d+)deg,\s*(.+)\)/);
		if (linearMatch) {
			const angle = parseInt(linearMatch[1]);
			const colorStops = linearMatch[2];
			
			// Extract all colors from color stops
			const colorMatches = colorStops.match(/([#\w]+)\s+[\d.]+%/g);
			if (colorMatches && colorMatches.length >= 2) {
				// Extract just the color hex from first and last stops
				const firstColor = colorMatches[0].match(/([#\w]+)/)?.[1];
				const lastColor = colorMatches[colorMatches.length - 1].match(/([#\w]+)/)?.[1];
				
				if (firstColor && lastColor) {
					const result = {
						type: 'linear' as const,
						angle: angle,
						fromColor: firstColor,
						toColor: lastColor
					};
					console.log('[parseGradientString] Linear match found:', result);
					return result;
				}
			}
		}
		
		// Parse radial gradient
		const radialMatch = gradientStr.match(/radial-gradient\(circle,\s*(.+)\)/);
		if (radialMatch) {
			const colorStops = radialMatch[1];
			const colorMatches = colorStops.match(/([#\w]+)\s+[\d.]+%/g);
			
			if (colorMatches && colorMatches.length >= 2) {
				const firstColor = colorMatches[0].match(/([#\w]+)/)?.[1];
				const lastColor = colorMatches[colorMatches.length - 1].match(/([#\w]+)/)?.[1];
				
				if (firstColor && lastColor) {
					const result = {
						type: 'radial' as const,
						angle: 0,
						fromColor: firstColor,
						toColor: lastColor
					};
					console.log('[parseGradientString] Radial match found:', result);
					return result;
				}
			}
		}
		
		console.log('[parseGradientString] No match found, returning null');
		return null;
	}

	// Custom Gradient Functions
	function openCustomGradientBuilder() {
		console.log('[openCustomGradientBuilder] Called');
		const bg = appearance.settings.background;
		console.log('[openCustomGradientBuilder] Background settings:', {
			type: bg.type,
			gradient: bg.gradient,
			customGradient: bg.customGradient
		});
		
		if (bg.customGradient?.enabled) {
			// Load existing custom gradient config
			console.log('[openCustomGradientBuilder] Loading from customGradient config');
			customGradientDraft = {
				type: bg.customGradient.type,
				angle: bg.customGradient.angle,
				fromColor: bg.customGradient.fromColor,
				toColor: bg.customGradient.toColor
			};
		} else if (bg.gradient) {
			// Try to parse current gradient string
			console.log('[openCustomGradientBuilder] Trying to parse gradient string');
			const parsed = parseGradientString(bg.gradient);
			if (parsed) {
				console.log('[openCustomGradientBuilder] Using parsed gradient');
				customGradientDraft = parsed;
			} else {
				// Use default values if parsing fails
				console.log('[openCustomGradientBuilder] Parsing failed, using defaults');
				customGradientDraft = {
					type: 'linear',
					angle: 135,
					fromColor: '#667eea',
					toColor: '#764ba2'
				};
			}
		} else {
			// Use default values
			console.log('[openCustomGradientBuilder] No gradient found, using defaults');
			customGradientDraft = {
				type: 'linear',
				angle: 135,
				fromColor: '#667eea',
				toColor: '#764ba2'
			};
		}
		
		console.log('[openCustomGradientBuilder] Final customGradientDraft:', customGradientDraft);
		showCustomGradientBuilder = true;
	}

	function closeCustomGradientBuilder() {
		showCustomGradientBuilder = false;
	}

	function selectPresetGradient(presetGradient: string) {
		updateSetting('background.gradient', presetGradient);
		updateSetting('background.customGradient.enabled', false);
	}

	function isValidHexColor(color: string): boolean {
		return /^#[0-9A-Fa-f]{6}$/.test(color);
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
				<button class="btn-primary btn-save" class:has-changes={appearance.dirty} onclick={handleSave} disabled={!appearance.dirty || appearance.saving}>
					{#if appearance.saving}
						<span class="spinner"></span>
					{:else}
						<Save size={16} />
						{#if appearance.dirty}
							<span class="change-indicator"></span>
						{/if}
					{/if}
					<span>{appearance.saving ? 'Đang lưu...' : 'Lưu tất cả'}</span>
				</button>
			</div>
		</header>

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
							<div class="theme-grid">
								<!-- My Theme Card (if exists) -->
								{#if appearance.customTheme}
									{@const basePreset = appearance.presets.find(p => p.id === appearance.customTheme?.based_on_preset_id)}
									{@const styles = getThemePreviewStyles(basePreset)}
									<div class="theme-card-wrapper" class:active={appearance.isUsingCustom}>
										<button 
											class="theme-card-phone"
											onclick={() => selectCustomTheme()}
											aria-label="Select My Theme"
										>
											<!-- Phone Preview -->
											<div class="phone-preview" style="background: {styles.background}">
												<!-- Custom Badge -->
												<div class="custom-badge-overlay">
													<svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
														<path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
													</svg>
													<span>MY THEME</span>
												</div>
												
												<!-- Mini Bio Preview -->
												<div class="mini-bio-preview">
													<div class="mini-avatar" style="background: {styles.avatarBg}; border: {styles.avatarBorder}"></div>
													<div class="mini-name" style="color: {styles.textColor}">Username</div>
													<div class="mini-bio" style="color: {styles.textMuted}">Your bio text here</div>
													<div class="mini-socials">
														<div class="mini-social-icon" style="background: {styles.socialBg}"></div>
														<div class="mini-social-icon" style="background: {styles.socialBg}"></div>
														<div class="mini-social-icon" style="background: {styles.socialBg}"></div>
													</div>
												</div>
												
												<!-- Mini Link Blocks -->
												<div class="mini-links">
													<div class="mini-link" style="background: {styles.cardBg}; border-radius: {styles.cardRadius}; box-shadow: {styles.cardShadow}; border: {styles.cardBorder}"></div>
													<div class="mini-link" style="background: {styles.cardBg}; border-radius: {styles.cardRadius}; box-shadow: {styles.cardShadow}; border: {styles.cardBorder}"></div>
													<div class="mini-link" style="background: {styles.cardBg}; border-radius: {styles.cardRadius}; box-shadow: {styles.cardShadow}; border: {styles.cardBorder}"></div>
												</div>
											</div>
										</button>
										<div class="theme-card-footer">
											<span class="theme-name">My Theme</span>
											<button 
												class="btn-delete-mini" 
												onclick={(e) => {
													e.stopPropagation();
													if (confirm('Xóa My Theme? Bạn sẽ quay về preset gần nhất.')) {
														deleteCustomTheme();
													}
												}}
												title="Xóa"
												aria-label="Delete My Theme"
											>
												<X size={12} />
											</button>
										</div>
									</div>
								{/if}
								
								<!-- Preset Themes -->
								{#each appearance.presets as preset}
									{@const styles = getThemePreviewStyles(preset)}
									<div class="theme-card-wrapper" class:active={!appearance.isUsingCustom && appearance.selectedPresetId === preset.id}>
										<button class="theme-card-phone" onclick={() => selectPreset(preset.id)}>
											<!-- Phone Preview -->
											<div class="phone-preview" style="background: {styles.background}">
												<!-- Mini Bio Preview -->
												<div class="mini-bio-preview">
													<div class="mini-avatar" style="background: {styles.avatarBg}; border: {styles.avatarBorder}"></div>
													<div class="mini-name" style="color: {styles.textColor}">Username</div>
													<div class="mini-bio" style="color: {styles.textMuted}">Your bio text here</div>
													<div class="mini-socials">
														<div class="mini-social-icon" style="background: {styles.socialBg}"></div>
														<div class="mini-social-icon" style="background: {styles.socialBg}"></div>
														<div class="mini-social-icon" style="background: {styles.socialBg}"></div>
													</div>
												</div>
												
												<!-- Mini Link Blocks -->
												<div class="mini-links">
													<div class="mini-link" style="background: {styles.cardBg}; border-radius: {styles.cardRadius}; box-shadow: {styles.cardShadow}; border: {styles.cardBorder}"></div>
													<div class="mini-link" style="background: {styles.cardBg}; border-radius: {styles.cardRadius}; box-shadow: {styles.cardShadow}; border: {styles.cardBorder}"></div>
													<div class="mini-link" style="background: {styles.cardBg}; border-radius: {styles.cardRadius}; box-shadow: {styles.cardShadow}; border: {styles.cardBorder}"></div>
												</div>
											</div>
										</button>
										<div class="theme-card-footer">
											<span class="theme-name">{preset.name}</span>
											{#if preset.tier === 'pro'}<span class="badge pro">PRO</span>{/if}
										</div>
									</div>
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
									<div class="color-presets-row">
										<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
											<div class="gradient-icon"></div>
										</button>
										<input type="color" class="hidden-native-picker" value={appearance.settings.background.color} onchange={(e) => updateSetting('background.color', e.currentTarget.value)} />
										
										{#each presetColors as c}
											<button 
												class="color-preset-circle" 
												class:active={appearance.settings.background.color === c} 
												style="background:{c}; {c === '#ffffff' ? 'border: 2px solid #e5e5ea' : ''}" 
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
								<!-- Gradient Presets Grid -->
								{#if true}
									{@const currentGradient = appearance.settings.background.gradient}
									{@const hasCurrentInPresets = gradientPresets.includes(currentGradient)}
									{@const isCustomEnabled = appearance.settings.background.customGradient?.enabled}
									<div class="gradient-presets-grid">
										{#each gradientPresets as g, i}
											<button 
												class="preview-item" 
												style="background:{g}" 
												class:active={currentGradient === g && !isCustomEnabled}
												title="Gradient {i+1}" 
												onclick={() => selectPresetGradient(g)}
											></button>
										{/each}
										
										<!-- Custom Gradient Button -->
										<button 
											class="preview-item"
											class:active={isCustomEnabled || (!hasCurrentInPresets && currentGradient)}
											onclick={openCustomGradientBuilder}
											title="Custom Gradient"
											style="background: {isCustomEnabled || (!hasCurrentInPresets && currentGradient) ? currentGradient : 'linear-gradient(135deg, #e0e0e0 0%, #f5f5f5 100%)'}"
										>
											{#if !(isCustomEnabled || (!hasCurrentInPresets && currentGradient))}
												<div class="custom-icon">
													<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
														<line x1="12" y1="5" x2="12" y2="19"></line>
														<line x1="5" y1="12" x2="19" y2="12"></line>
													</svg>
												</div>
											{/if}
										</button>
									</div>
								{/if}
								
								<!-- Custom Gradient Builder -->
								{#if showCustomGradientBuilder}
									<div class="custom-gradient-builder">
										<div class="flex justify-between items-center mb-4 pb-3 border-b border-separator">
											<h3 class="text-base font-semibold">Custom Gradient</h3>
											<button class="btn-icon" onclick={closeCustomGradientBuilder} title="Close" aria-label="Close gradient builder">
												<X size={16} />
											</button>
										</div>
										
										<!-- Gradient Type -->
										<div class="flex justify-between items-center mb-4">
											<span class="setting-label">Type</span>
											<div class="option-group compact">
												<button 
													class="option-btn small"
													class:active={customGradientDraft.type === 'linear'}
													onclick={() => customGradientDraft.type = 'linear'}
													aria-label="Linear gradient"
												>
													Linear
												</button>
												<button 
													class="option-btn small"
													class:active={customGradientDraft.type === 'radial'}
													onclick={() => customGradientDraft.type = 'radial'}
													aria-label="Radial gradient"
												>
													Radial
												</button>
											</div>
										</div>
										
										<!-- Angle Slider (Linear only) -->
										{#if customGradientDraft.type === 'linear'}
											<div class="mb-4">
												<div class="flex justify-between items-center mb-2">
													<span class="slider-label">Angle</span>
													<span class="slider-value">{customGradientDraft.angle}°</span>
												</div>
												<input 
													type="range" 
													class="slider" 
													min="0" 
													max="360" 
													step="15"
													bind:value={customGradientDraft.angle}
													aria-label="Gradient angle"
												/>
											</div>
										{/if}
										
										<!-- Colors Row -->
										<div class="grid grid-cols-2 gap-3 my-4">
											<!-- From Color -->
											<div>
												<div class="flex justify-between items-center mb-2">
													<span class="setting-label">From Color</span>
													<div class="flex gap-2">
														<input 
															type="color" 
															class="color-picker-trigger"
															bind:value={customGradientDraft.fromColor}
															aria-label="Pick from color"
														/>
														<input 
															type="text" 
															class="color-hex-input"
															class:error={!isValidHexColor(customGradientDraft.fromColor)}
															bind:value={customGradientDraft.fromColor}
															placeholder="#667eea"
															aria-label="From color hex code"
														/>
													</div>
												</div>
											</div>
											
											<!-- To Color -->
											<div>
												<div class="flex justify-between items-center mb-2">
													<span class="setting-label">To Color</span>
													<div class="flex gap-2">
														<input 
															type="color" 
															class="color-picker-trigger"
															bind:value={customGradientDraft.toColor}
															aria-label="Pick to color"
														/>
														<input 
															type="text" 
															class="color-hex-input"
															class:error={!isValidHexColor(customGradientDraft.toColor)}
															bind:value={customGradientDraft.toColor}
															placeholder="#764ba2"
															aria-label="To color hex code"
														/>
													</div>
												</div>
											</div>
										</div>
										
										<!-- Live Preview -->
										<div class="my-4">
											<div class="text-sm font-medium text-text-secondary mb-2">Preview</div>
											<div 
												class="w-full h-[120px] rounded-md border border-separator shadow-inner"
												style="background:{liveGradientPreview}"
												role="img"
												aria-label="Gradient preview"
											></div>
										</div>
									</div>
								{/if}
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
								{#if appearance.settings.header.style === 'classic'}
									<button class="sub-tab" class:active={activeHeaderTab === 'cover'} onclick={() => activeHeaderTab = 'cover'}>
										<Image size={16} />
										<span>Cover</span>
									</button>
								{/if}
								<button class="sub-tab" class:active={activeHeaderTab === 'social'} onclick={() => activeHeaderTab = 'social'}>
									<Instagram size={16} />
									<span>Social</span>
								</button>
							</div>

							<!-- Tab Content -->
							<div class="sub-tab-content">
								{#if activeHeaderTab === 'avatar'}
									<!-- Avatar Settings -->
									
									<!-- Avatar Shape with Visual Preview -->
									<div class="setting-section">
										<span class="setting-label">Hình dạng Avatar</span>
										<p class="form-hint">Chọn kiểu hiển thị cho ảnh đại diện của bạn</p>
										
										<div class="avatar-shape-grid">
											{#each avatarShapes as shape}
												<button 
													class="avatar-shape-card" 
													class:active={appearance.settings.header.avatarShape === shape.value}
													onclick={() => updateSetting('header.avatarShape', shape.value)}
												>
													<div class="avatar-preview-wrapper">
														<div class="avatar-preview {shape.value}">
															<div class="avatar-demo-image"></div>
														</div>
													</div>
													<span class="shape-label">{shape.label}</span>
												</button>
											{/each}
										</div>
									</div>

									<!-- Avatar Size -->
									<div class="setting-section">
										<span class="setting-label">Kích thước Avatar</span>
										<div class="avatar-size-options">
											{#each avatarSizes as size}
												<button 
													class="avatar-size-btn" 
													class:active={appearance.settings.header.avatarSize === size.value}
													onclick={() => updateSetting('header.avatarSize', size.value)}
												>
													<div class="size-preview {size.value.toLowerCase()} {appearance.settings.header.avatarShape}">
														<div class="size-demo-image"></div>
													</div>
													<span class="size-label">{size.label}</span>
												</button>
											{/each}
										</div>
									</div>
									
									<!-- Avatar Border Width Slider -->
									<div class="setting-row" style="margin-top: var(--space-4)">
										<span class="setting-label">Độ dày viền Avatar</span>
									</div>
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
									<div class="color-setting-group" style="margin-top: var(--space-2)">
										<div class="color-presets-row">
											<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
												<div class="gradient-icon"></div>
											</button>
											<input type="color" class="hidden-native-picker" value={appearance.settings.header.avatarBorderColor || '#ffffff'} onchange={(e) => updateSetting('header.avatarBorderColor', e.currentTarget.value)} />
											
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || '#ffffff') === '#000000'} style="background:#000000" title="#000000" onclick={() => updateSetting('header.avatarBorderColor', '#000000')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || '#ffffff') === '#ffffff'} style="background:#ffffff; border: 2px solid #e5e5ea" title="#ffffff" onclick={() => updateSetting('header.avatarBorderColor', '#ffffff')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || '#ffffff') === '#1c1c1e'} style="background:#1c1c1e" title="#1c1c1e" onclick={() => updateSetting('header.avatarBorderColor', '#1c1c1e')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || '#ffffff') === '#8e8e93'} style="background:#8e8e93" title="#8e8e93" onclick={() => updateSetting('header.avatarBorderColor', '#8e8e93')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || '#ffffff') === '#ff3b30'} style="background:#ff3b30" title="#ff3b30" onclick={() => updateSetting('header.avatarBorderColor', '#ff3b30')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || '#ffffff') === '#ff9500'} style="background:#ff9500" title="#ff9500" onclick={() => updateSetting('header.avatarBorderColor', '#ff9500')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || '#ffffff') === '#007aff'} style="background:#007aff" title="#007aff" onclick={() => updateSetting('header.avatarBorderColor', '#007aff')}></button>
											<button class="color-preset-circle" class:active={(appearance.settings.header.avatarBorderColor || '#ffffff') === '#5856d6'} style="background:#5856d6" title="#5856d6" onclick={() => updateSetting('header.avatarBorderColor', '#5856d6')}></button>
											
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
											
											<!-- Theme default color button -->
											{#if !appearance.settings.header.nameColor}
												<button class="color-preset-circle active" style="background:{appearance.settings.colors.text}" title="Theme default: {appearance.settings.colors.text}" onclick={() => updateSetting('header.nameColor', '')}></button>
											{:else}
												<button class="color-preset-circle" style="background:{appearance.settings.colors.text}" title="Theme default: {appearance.settings.colors.text}" onclick={() => updateSetting('header.nameColor', '')}></button>
											{/if}
											
											<button class="color-preset-circle" class:active={appearance.settings.header.nameColor === '#000000'} style="background:#000000" title="#000000" onclick={() => updateSetting('header.nameColor', '#000000')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.header.nameColor === '#ffffff'} style="background:#ffffff; border: 2px solid #e5e5ea" title="#ffffff" onclick={() => updateSetting('header.nameColor', '#ffffff')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.header.nameColor === '#1c1c1e'} style="background:#1c1c1e" title="#1c1c1e" onclick={() => updateSetting('header.nameColor', '#1c1c1e')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.header.nameColor === '#8e8e93'} style="background:#8e8e93" title="#8e8e93" onclick={() => updateSetting('header.nameColor', '#8e8e93')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.header.nameColor === '#ff3b30'} style="background:#ff3b30" title="#ff3b30" onclick={() => updateSetting('header.nameColor', '#ff3b30')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.header.nameColor === '#ff9500'} style="background:#ff9500" title="#ff9500" onclick={() => updateSetting('header.nameColor', '#ff9500')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.header.nameColor === '#007aff'} style="background:#007aff" title="#007aff" onclick={() => updateSetting('header.nameColor', '#007aff')}></button>
											<button class="color-preset-circle" class:active={appearance.settings.header.nameColor === '#5856d6'} style="background:#5856d6" title="#5856d6" onclick={() => updateSetting('header.nameColor', '#5856d6')}></button>
											
											{#if appearance.settings.header.nameColor}
												<button class="color-preset-circle reset-circle" title="Reset về mặc định" onclick={() => updateSetting('header.nameColor', '')}>
													<RefreshCw size={18} />
												</button>
											{/if}
										</div>
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
											<span class="setting-label">Màu Bio</span>
										</div>
										<div class="color-setting-group" style="margin-top: var(--space-2)">
											<div class="color-presets-row">
												<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
													<div class="gradient-icon"></div>
												</button>
												<input type="color" class="hidden-native-picker" value={appearance.settings.header.bioColor || appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('header.bioColor', e.currentTarget.value)} />
												
												<!-- Theme default color button -->
												{#if !appearance.settings.header.bioColor}
													<button class="color-preset-circle active" style="background:{appearance.settings.colors.textSecondary}" title="Theme default: {appearance.settings.colors.textSecondary}" onclick={() => updateSetting('header.bioColor', '')}></button>
												{:else}
													<button class="color-preset-circle" style="background:{appearance.settings.colors.textSecondary}" title="Theme default: {appearance.settings.colors.textSecondary}" onclick={() => updateSetting('header.bioColor', '')}></button>
												{/if}
												
												<button class="color-preset-circle" class:active={appearance.settings.header.bioColor === '#636366'} style="background:#636366" title="#636366" onclick={() => updateSetting('header.bioColor', '#636366')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.bioColor === '#8e8e93'} style="background:#8e8e93" title="#8e8e93" onclick={() => updateSetting('header.bioColor', '#8e8e93')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.bioColor === '#aeaeb2'} style="background:#aeaeb2" title="#aeaeb2" onclick={() => updateSetting('header.bioColor', '#aeaeb2')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.bioColor === '#c7c7cc'} style="background:#c7c7cc" title="#c7c7cc" onclick={() => updateSetting('header.bioColor', '#c7c7cc')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.bioColor === '#1c1c1e'} style="background:#1c1c1e" title="#1c1c1e" onclick={() => updateSetting('header.bioColor', '#1c1c1e')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.bioColor === '#3a3a3c'} style="background:#3a3a3c" title="#3a3a3c" onclick={() => updateSetting('header.bioColor', '#3a3a3c')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.bioColor === '#48484a'} style="background:#48484a" title="#48484a" onclick={() => updateSetting('header.bioColor', '#48484a')}></button>
												
												{#if appearance.settings.header.bioColor}
													<button class="color-preset-circle reset-circle" title="Reset về mặc định" onclick={() => updateSetting('header.bioColor', '')}>
														<RefreshCw size={18} />
													</button>
												{/if}
											</div>
										</div>
									{/if}

								{:else if activeHeaderTab === 'cover'}
									<!-- Cover Settings -->
									{#if appearance.settings.header.style !== 'classic'}
										<div class="info-box">
											<AlertTriangle size={16} />
											<span>Cover chỉ hiển thị khi Header Style là "Classic"</span>
										</div>
									{/if}
									
									<!-- Cover Type -->
									<div class="setting-row">
										<span class="setting-label">Loại Cover</span>
										<div class="button-group-inline">
											<button 
												class="btn-option" 
												class:active={appearance.settings.header.cover.type === 'color'}
												onclick={() => updateSetting('header.cover.type', 'color')}
											>
												<Droplets size={14} />
												<span>Màu</span>
											</button>
											<button 
												class="btn-option" 
												class:active={appearance.settings.header.cover.type === 'image'}
												onclick={() => updateSetting('header.cover.type', 'image')}
											>
												<Image size={14} />
												<span>Ảnh</span>
											</button>
										</div>
									</div>

									{#if appearance.settings.header.cover.type === 'color'}
										<!-- Color Picker -->
										<div class="setting-row">
											<span class="setting-label">Màu Cover</span>
										</div>
										<div class="color-setting-group">
											<div class="color-presets-row">
												<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
													<div class="gradient-icon"></div>
												</button>
												<input type="color" class="hidden-native-picker" value={appearance.settings.header.cover.color} onchange={(e) => updateSetting('header.cover.color', e.currentTarget.value)} />
												
												{#each ['#c2185b', '#d32f2f', '#7b1fa2', '#512da8', '#303f9f', '#1976d2', '#0288d1', '#0097a7', '#00796b', '#388e3c', '#689f38', '#afb42b', '#fbc02d', '#ffa000', '#f57c00', '#e64a19', '#5d4037', '#616161', '#455a64', '#000000'] as color}
													<button 
														class="color-preset-circle" 
														class:active={appearance.settings.header.cover.color === color}
														style="background:{color}{color === '#000000' ? '' : ''}; {color === '#ffffff' ? 'border: 2px solid #e5e5ea' : ''}" 
														title={color}
														onclick={() => updateSetting('header.cover.color', color)}
													></button>
												{/each}
											</div>
										</div>
									{:else}
										<!-- Image Upload -->
										<div class="setting-row">
											<span class="setting-label">Ảnh Cover</span>
										</div>
										<div class="upload-area">
											{#if appearance.settings.header.cover.imageUrl}
												<div class="uploaded-image">
													<img src={appearance.settings.header.cover.imageUrl} alt="Cover" />
													<button class="remove-image" onclick={() => { updateSetting('header.cover.imageUrl', null); updateSetting('header.cover.imageAssetId', null); }}>
														<X size={16} />
													</button>
												</div>
											{:else}
												<label class="upload-label">
													<input type="file" accept="image/*" onchange={handleCoverUpload} style="display: none;" />
													<Image size={24} />
													<span>Click để upload ảnh cover</span>
													<span class="upload-hint">Khuyến nghị: 1200x400px, tối đa 2MB</span>
												</label>
											{/if}
										</div>
									{/if}

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
										
										<!-- Social Icons Style Settings -->
										<div class="toggle-row" style="margin-top: var(--space-4)">
											<span>Hiển thị Background</span>
											<button class="toggle" class:active={appearance.settings.header.socialIconsBg} aria-label="Toggle" onclick={() => updateSetting('header.socialIconsBg', !appearance.settings.header.socialIconsBg)}></button>
										</div>
										
										<!-- Social Icons Color Picker -->
										<div class="setting-row" style="margin-top: var(--space-4)">
											<span class="setting-label">Màu Icon Social</span>
										</div>
										<div class="color-setting-group" style="margin-top: var(--space-2)">
											<div class="color-presets-row">
												<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
													<div class="gradient-icon"></div>
												</button>
												<input type="color" class="hidden-native-picker" value={appearance.settings.header.socialIconsColor || appearance.settings.colors.textSecondary} onchange={(e) => updateSetting('header.socialIconsColor', e.currentTarget.value)} />
												
												<!-- Theme default color button -->
												{#if !appearance.settings.header.socialIconsColor}
													<button class="color-preset-circle active" style="background:{appearance.settings.colors.textSecondary}" title="Theme default: {appearance.settings.colors.textSecondary}" onclick={() => updateSetting('header.socialIconsColor', '')}></button>
												{:else}
													<button class="color-preset-circle" style="background:{appearance.settings.colors.textSecondary}" title="Theme default: {appearance.settings.colors.textSecondary}" onclick={() => updateSetting('header.socialIconsColor', '')}></button>
												{/if}
												
												<button class="color-preset-circle" class:active={appearance.settings.header.socialIconsColor === '#000000'} style="background:#000000" title="#000000" onclick={() => updateSetting('header.socialIconsColor', '#000000')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.socialIconsColor === '#ffffff'} style="background:#ffffff; border: 2px solid #e5e5ea" title="#ffffff" onclick={() => updateSetting('header.socialIconsColor', '#ffffff')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.socialIconsColor === '#8e8e93'} style="background:#8e8e93" title="#8e8e93" onclick={() => updateSetting('header.socialIconsColor', '#8e8e93')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.socialIconsColor === '#ff3b30'} style="background:#ff3b30" title="#ff3b30" onclick={() => updateSetting('header.socialIconsColor', '#ff3b30')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.socialIconsColor === '#ff9500'} style="background:#ff9500" title="#ff9500" onclick={() => updateSetting('header.socialIconsColor', '#ff9500')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.socialIconsColor === '#007aff'} style="background:#007aff" title="#007aff" onclick={() => updateSetting('header.socialIconsColor', '#007aff')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.socialIconsColor === '#5856d6'} style="background:#5856d6" title="#5856d6" onclick={() => updateSetting('header.socialIconsColor', '#5856d6')}></button>
												<button class="color-preset-circle" class:active={appearance.settings.header.socialIconsColor === '#34c759'} style="background:#34c759" title="#34c759" onclick={() => updateSetting('header.socialIconsColor', '#34c759')}></button>
												
												{#if appearance.settings.header.socialIconsColor}
													<button class="color-preset-circle reset-circle" title="Reset về mặc định" onclick={() => updateSetting('header.socialIconsColor', '')}>
														<RefreshCw size={18} />
													</button>
												{/if}
											</div>
										</div>
									{/if}
								{/if}
							</div>
						</div>
					</div>
				{/if}

				{#if activeSection === 'blocks'}
					<!-- Block Style Selection -->
					<div class="card">
						<div class="card-header"><h2>Block Style</h2></div>
						<div class="card-body">
							<div class="block-style-grid">
								{#each linkStyles as s}
									<button 
										class="block-style-card" 
										class:active={appearance.settings.links.style === s.value} 
										onclick={() => updateSetting('links.style', s.value)}
									>
										<div class="block-style-preview" style="background:{getBgStyle()}">
											<div class="block-style-demo" style="{getBlockDemoStyle(s.value)}">
												<span>Link</span>
											</div>
										</div>
										<span class="block-style-label">{s.label}</span>
									</button>
								{/each}
							</div>
						</div>
					</div>

					<!-- Block Settings with Tabs -->
					<div class="card">
						<div class="card-header"><h2>Cài đặt Block</h2></div>
						<div class="card-body">
							<!-- Sub-tabs Navigation -->
							<div class="sub-tabs">
								<button class="sub-tab" class:active={activeBlockTab === 'basic'} onclick={() => activeBlockTab = 'basic'}>
									<Droplets size={16} />
									<span>Cơ bản</span>
								</button>
								<button class="sub-tab" class:active={activeBlockTab === 'advanced'} onclick={() => activeBlockTab = 'advanced'}>
									<Settings size={16} />
									<span>Tùy chỉnh nâng cao</span>
								</button>
							</div>

							<!-- Tab Content -->
							<div class="sub-tab-content">
								{#if activeBlockTab === 'basic'}
									<!-- Block Colors - Basic Tab -->
									<p class="form-hint" style="margin-bottom: var(--space-4)">Tùy chỉnh màu sắc cho các block link</p>
									
									<!-- Background Color -->
									<div class="color-setting-group">
										<label class="setting-label">Background</label>
										{#if true}
											{@const currentBgColor = normalizeColor(getActualBlockBgColor())}
											{@const staticPresets = ['#ffffff', '#f9fafb', '#fce7f3', '#dbeafe', '#fef3c7', '#000000', '#1c1c1e', '#2c2c2e']}
											{@const hasCurrentInPresets = staticPresets.includes(currentBgColor)}
											<div class="color-presets-row">
												<!-- Color Picker Button -->
												<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
													<div class="gradient-icon"></div>
												</button>
												<input type="color" class="hidden-native-picker" value={getActualBlockBgColor()} onchange={(e) => {
													const style = appearance.settings.links.style;
													if (style === 'filled' || style === 'soft') {
														updateSetting('colors.primary', e.currentTarget.value);
												} else {
													updateSetting('colors.cardBackground', e.currentTarget.value);
												}
											}} />
											
											<!-- Current Color Button (if not in presets) -->
											{#if !hasCurrentInPresets && currentBgColor !== 'transparent'}
												<button 
													class="color-preset-circle active current-color" 
													style="background:{currentBgColor}; {currentBgColor === '#ffffff' ? 'border: 2px solid #e5e5ea' : ''}"
													title="Current: {currentBgColor}"
													onclick={() => {
														const style = appearance.settings.links.style;
														if (style === 'filled' || style === 'soft') {
															updateSetting('colors.primary', currentBgColor);
														} else {
															updateSetting('colors.cardBackground', currentBgColor);
														}
													}}
												>
													<span class="current-indicator">●</span>
												</button>
											{/if}
											
											<!-- Static Preset Colors -->
											<button class="color-preset-circle" class:active={currentBgColor === '#ffffff'} style="background:#ffffff; border: 2px solid #e5e5ea" title="#ffffff" onclick={() => {
												const style = appearance.settings.links.style;
												if (style === 'filled' || style === 'soft') {
													updateSetting('colors.primary', '#ffffff');
												} else {
													updateSetting('colors.cardBackground', '#ffffff');
												}
											}}></button>
											<button class="color-preset-circle" class:active={currentBgColor === '#f9fafb'} style="background:#f9fafb" title="#f9fafb" onclick={() => {
												const style = appearance.settings.links.style;
												if (style === 'filled' || style === 'soft') {
													updateSetting('colors.primary', '#f9fafb');
												} else {
													updateSetting('colors.cardBackground', '#f9fafb');
												}
											}}></button>
											<button class="color-preset-circle" class:active={currentBgColor === '#fce7f3'} style="background:#fce7f3" title="#fce7f3 (pink)" onclick={() => {
												const style = appearance.settings.links.style;
												if (style === 'filled' || style === 'soft') {
													updateSetting('colors.primary', '#fce7f3');
												} else {
													updateSetting('colors.cardBackground', '#fce7f3');
												}
											}}></button>
											<button class="color-preset-circle" class:active={currentBgColor === '#dbeafe'} style="background:#dbeafe" title="#dbeafe (blue)" onclick={() => {
												const style = appearance.settings.links.style;
												if (style === 'filled' || style === 'soft') {
													updateSetting('colors.primary', '#dbeafe');
												} else {
													updateSetting('colors.cardBackground', '#dbeafe');
												}
											}}></button>
											<button class="color-preset-circle" class:active={currentBgColor === '#fef3c7'} style="background:#fef3c7" title="#fef3c7 (yellow)" onclick={() => {
												const style = appearance.settings.links.style;
												if (style === 'filled' || style === 'soft') {
													updateSetting('colors.primary', '#fef3c7');
												} else {
													updateSetting('colors.cardBackground', '#fef3c7');
												}
											}}></button>
											<button class="color-preset-circle" class:active={currentBgColor === '#000000'} style="background:#000000" title="#000000" onclick={() => {
												const style = appearance.settings.links.style;
												if (style === 'filled' || style === 'soft') {
													updateSetting('colors.primary', '#000000');
												} else {
													updateSetting('colors.cardBackground', '#000000');
												}
											}}></button>
											<button class="color-preset-circle" class:active={currentBgColor === '#1c1c1e'} style="background:#1c1c1e" title="#1c1c1e" onclick={() => {
												const style = appearance.settings.links.style;
												if (style === 'filled' || style === 'soft') {
													updateSetting('colors.primary', '#1c1c1e');
												} else {
													updateSetting('colors.cardBackground', '#1c1c1e');
												}
											}}></button>
											<button class="color-preset-circle" class:active={currentBgColor === '#2c2c2e'} style="background:#2c2c2e" title="#2c2c2e" onclick={() => {
												const style = appearance.settings.links.style;
												if (style === 'filled' || style === 'soft') {
													updateSetting('colors.primary', '#2c2c2e');
												} else {
													updateSetting('colors.cardBackground', '#2c2c2e');
												}
											}}></button>
										</div>
										{/if}
									</div>

									<!-- Text Color -->
									<div class="color-setting-group">
										<label class="setting-label">Text</label>
										{#if true}
											{@const currentTextColor = normalizeColor(getActualBlockTextColor())}
											{@const staticPresets = ['#000000', '#ffffff', '#1c1c1e', '#8e8e93', '#ff3b30', '#ff6b6b', '#ffa5a5', '#a855f7', '#5856d6']}
											{@const hasCurrentInPresets = staticPresets.includes(currentTextColor)}
											<div class="color-presets-row">
												<!-- Color Picker Button -->
												<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
													<div class="gradient-icon"></div>
												</button>
												<input type="color" class="hidden-native-picker" value={getActualBlockTextColor()} onchange={(e) => updateSetting('links.textColor', e.currentTarget.value)} />
												
												<!-- Current Color Button (if not in presets) -->
												{#if !hasCurrentInPresets}
													<button 
														class="color-preset-circle active current-color" 
														style="background:{currentTextColor}; {currentTextColor === '#ffffff' ? 'border: 2px solid #e5e5ea' : ''}"
														title="Current: {currentTextColor}"
														onclick={() => updateSetting('links.textColor', currentTextColor)}
													>
														<span class="current-indicator">●</span>
													</button>
												{/if}
												
												<!-- Static Preset Colors -->
												<button class="color-preset-circle" class:active={currentTextColor === '#000000'} style="background:#000000" title="#000000" onclick={() => updateSetting('links.textColor', '#000000')}></button>
												<button class="color-preset-circle" class:active={currentTextColor === '#ffffff'} style="background:#ffffff; border: 2px solid #e5e5ea" title="#ffffff" onclick={() => updateSetting('links.textColor', '#ffffff')}></button>
												<button class="color-preset-circle" class:active={currentTextColor === '#1c1c1e'} style="background:#1c1c1e" title="#1c1c1e" onclick={() => updateSetting('links.textColor', '#1c1c1e')}></button>
												<button class="color-preset-circle" class:active={currentTextColor === '#8e8e93'} style="background:#8e8e93" title="#8e8e93" onclick={() => updateSetting('links.textColor', '#8e8e93')}></button>
												<button class="color-preset-circle" class:active={currentTextColor === '#ff3b30'} style="background:#ff3b30" title="#ff3b30" onclick={() => updateSetting('links.textColor', '#ff3b30')}></button>
												<button class="color-preset-circle" class:active={currentTextColor === '#ff6b6b'} style="background:#ff6b6b" title="#ff6b6b" onclick={() => updateSetting('links.textColor', '#ff6b6b')}></button>
												<button class="color-preset-circle" class:active={currentTextColor === '#ffa5a5'} style="background:#ffa5a5" title="#ffa5a5" onclick={() => updateSetting('links.textColor', '#ffa5a5')}></button>
												<button class="color-preset-circle" class:active={currentTextColor === '#a855f7'} style="background:#a855f7" title="#a855f7" onclick={() => updateSetting('links.textColor', '#a855f7')}></button>
												<button class="color-preset-circle" class:active={currentTextColor === '#5856d6'} style="background:#5856d6" title="#5856d6" onclick={() => updateSetting('links.textColor', '#5856d6')}></button>
											</div>
										{/if}
									</div>

									<!-- Border Radius Slider -->
									<div class="slider-group" style="margin-top: var(--space-4)">
										<div class="slider-header">
											<span class="slider-label">Border Radius</span>
											<span class="slider-value">{appearance.settings.links.borderRadius}px</span>
										</div>
										<input type="range" class="slider" min="0" max="24" step="2" value={appearance.settings.links.borderRadius} oninput={(e) => updateSetting('links.borderRadius', +e.currentTarget.value)}/>
									</div>

								{:else if activeBlockTab === 'advanced'}
									<!-- Appearance Settings - Advanced Tab -->
									<p class="form-hint" style="margin-bottom: var(--space-4)">Tùy chỉnh chi tiết giao diện block</p>
									
									<!-- Toggle Show Background -->
									<div class="toggle-row" style="margin-bottom: var(--space-4)">
										<div>
											<span class="setting-label">Hiển thị Background</span>
											<p class="form-hint" style="margin-top: 4px">Bật/tắt màu nền cho block</p>
										</div>
										<button 
											class="toggle" 
											class:active={appearance.settings.links.showBackground} 
											aria-label="Toggle background" 
											onclick={() => updateSetting('links.showBackground', !appearance.settings.links.showBackground)}
										></button>
									</div>
									
									<!-- Toggle Show Border -->
									<div class="toggle-row" style="margin-bottom: var(--space-4)">
										<div>
											<span class="setting-label">Hiển thị Border</span>
											<p class="form-hint" style="margin-top: 4px">Bật/tắt viền cho block</p>
										</div>
										<button 
											class="toggle" 
											class:active={appearance.settings.links.showBorder} 
											aria-label="Toggle border" 
											onclick={() => updateSetting('links.showBorder', !appearance.settings.links.showBorder)}
										></button>
									</div>
									
									<!-- Border Settings (conditional) -->
									{#if appearance.settings.links.showBorder}
										<div class="border-settings" style="margin-bottom: var(--space-4); padding: var(--space-3); background: var(--color-bg); border-radius: var(--radius-md);">
											<!-- Border Color -->
											<div class="color-setting-group" style="margin-bottom: var(--space-3);">
												<label class="setting-label" style="margin-bottom: var(--space-2); display: block;">Màu Border</label>
												<div class="color-presets-row">
													<!-- Color Picker Button -->
													<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
														<div class="gradient-icon"></div>
													</button>
													<input type="color" class="hidden-native-picker" value={appearance.settings.links.borderColor || '#e5e5ea'} onchange={(e) => updateSetting('links.borderColor', e.currentTarget.value)} />
													
													<!-- Preset Colors -->
													<button class="color-preset-circle" class:active={(appearance.settings.links.borderColor || '#e5e5ea') === '#e5e5ea'} style="background:#e5e5ea" title="#e5e5ea" onclick={() => updateSetting('links.borderColor', '#e5e5ea')}></button>
													<button class="color-preset-circle" class:active={(appearance.settings.links.borderColor || '#e5e5ea') === '#d1d1d6'} style="background:#d1d1d6" title="#d1d1d6" onclick={() => updateSetting('links.borderColor', '#d1d1d6')}></button>
													<button class="color-preset-circle" class:active={(appearance.settings.links.borderColor || '#e5e5ea') === '#c7c7cc'} style="background:#c7c7cc" title="#c7c7cc" onclick={() => updateSetting('links.borderColor', '#c7c7cc')}></button>
													<button class="color-preset-circle" class:active={(appearance.settings.links.borderColor || '#e5e5ea') === '#8e8e93'} style="background:#8e8e93" title="#8e8e93" onclick={() => updateSetting('links.borderColor', '#8e8e93')}></button>
													<button class="color-preset-circle" class:active={(appearance.settings.links.borderColor || '#e5e5ea') === '#000000'} style="background:#000000" title="#000000" onclick={() => updateSetting('links.borderColor', '#000000')}></button>
													<button class="color-preset-circle" class:active={(appearance.settings.links.borderColor || '#e5e5ea') === '#007aff'} style="background:#007aff" title="#007aff" onclick={() => updateSetting('links.borderColor', '#007aff')}></button>
													<button class="color-preset-circle" class:active={(appearance.settings.links.borderColor || '#e5e5ea') === '#5856d6'} style="background:#5856d6" title="#5856d6" onclick={() => updateSetting('links.borderColor', '#5856d6')}></button>
												</div>
											</div>
											
											<!-- Border Width Slider -->
											<div class="slider-group">
												<div class="slider-header">
													<span class="slider-label">Độ dày Border</span>
													<span class="slider-value">{appearance.settings.links.borderWidth}px</span>
												</div>
												<input 
													type="range" 
													class="slider" 
													style="background: white;"
													min="1" 
													max="8" 
													step="1" 
													value={appearance.settings.links.borderWidth} 
													oninput={(e) => updateSetting('links.borderWidth', +e.currentTarget.value)}
												/>
											</div>
										</div>
									{/if}
									
									<!-- Toggle Show Shadow -->
									<div class="toggle-row" style="margin-bottom: var(--space-4)">
										<div>
											<span class="setting-label">Hiển thị Shadow</span>
											<p class="form-hint" style="margin-top: 4px">Bật/tắt bóng đổ cho block</p>
										</div>
										<button 
											class="toggle" 
											class:active={appearance.settings.links.showShadow} 
											aria-label="Toggle shadow" 
											onclick={() => updateSetting('links.showShadow', !appearance.settings.links.showShadow)}
										></button>
									</div>
									
									<!-- Shadow Settings (conditional) -->
									{#if appearance.settings.links.showShadow}
										<div class="shadow-settings" style="margin-bottom: var(--space-4); padding: var(--space-3); background: var(--color-bg); border-radius: var(--radius-md);">
											<!-- Shadow Color -->
											<div class="color-setting-group" style="margin-bottom: var(--space-3);">
												<label class="setting-label" style="margin-bottom: var(--space-2); display: block;">Màu Shadow</label>
												<div class="color-presets-row">
													<!-- Color Picker Button -->
													<button class="color-preset-circle gradient-picker" title="Chọn màu tùy chỉnh" onclick={(e) => { e.currentTarget.nextElementSibling?.click(); }}>
														<div class="gradient-icon"></div>
													</button>
													<input type="color" class="hidden-native-picker" value={appearance.settings.links.shadowColor || '#000000'} onchange={(e) => updateSetting('links.shadowColor', e.currentTarget.value)} />
													
													<!-- Preset Colors -->
													<button class="color-preset-circle" class:active={(appearance.settings.links.shadowColor || '#000000') === '#000000'} style="background:#000000" title="#000000" onclick={() => updateSetting('links.shadowColor', '#000000')}></button>
													<button class="color-preset-circle" class:active={(appearance.settings.links.shadowColor || '#000000') === '#1c1c1e'} style="background:#1c1c1e" title="#1c1c1e" onclick={() => updateSetting('links.shadowColor', '#1c1c1e')}></button>
													<button class="color-preset-circle" class:active={(appearance.settings.links.shadowColor || '#000000') === '#3a3a3c'} style="background:#3a3a3c" title="#3a3a3c" onclick={() => updateSetting('links.shadowColor', '#3a3a3c')}></button>
													<button class="color-preset-circle" class:active={(appearance.settings.links.shadowColor || '#000000') === '#636366'} style="background:#636366" title="#636366" onclick={() => updateSetting('links.shadowColor', '#636366')}></button>
													<button class="color-preset-circle" class:active={(appearance.settings.links.shadowColor || '#000000') === '#007aff'} style="background:#007aff" title="#007aff" onclick={() => updateSetting('links.shadowColor', '#007aff')}></button>
													<button class="color-preset-circle" class:active={(appearance.settings.links.shadowColor || '#000000') === '#5856d6'} style="background:#5856d6" title="#5856d6" onclick={() => updateSetting('links.shadowColor', '#5856d6')}></button>
												</div>
											</div>
											
											<!-- Shadow Blur -->
											<div class="slider-group" style="margin-bottom: var(--space-3);">
												<div class="slider-header">
													<span class="slider-label">Độ mờ (Blur)</span>
													<span class="slider-value">{appearance.settings.links.shadowBlur}px</span>
												</div>
												<input 
													type="range" 
													class="slider" 
													style="background: white;"
													min="0" 
													max="32" 
													step="2" 
													value={appearance.settings.links.shadowBlur} 
													oninput={(e) => updateSetting('links.shadowBlur', +e.currentTarget.value)}
												/>
											</div>
											
											<!-- Shadow Offset X -->
											<div class="slider-group" style="margin-bottom: var(--space-3);">
												<div class="slider-header">
													<span class="slider-label">Offset X</span>
													<span class="slider-value">{appearance.settings.links.shadowOffsetX}px</span>
												</div>
												<input 
													type="range" 
													class="slider" 
													style="background: white;"
													min="-16" 
													max="16" 
													step="1" 
													value={appearance.settings.links.shadowOffsetX} 
													oninput={(e) => updateSetting('links.shadowOffsetX', +e.currentTarget.value)}
												/>
											</div>
											
											<!-- Shadow Offset Y -->
											<div class="slider-group" style="margin-bottom: var(--space-3);">
												<div class="slider-header">
													<span class="slider-label">Offset Y</span>
													<span class="slider-value">{appearance.settings.links.shadowOffsetY}px</span>
												</div>
												<input 
													type="range" 
													class="slider" 
													style="background: white;"
													min="-16" 
													max="16" 
													step="1" 
													value={appearance.settings.links.shadowOffsetY} 
													oninput={(e) => updateSetting('links.shadowOffsetY', +e.currentTarget.value)}
												/>
											</div>
											
											<!-- Shadow Opacity -->
											<div class="slider-group">
												<div class="slider-header">
													<span class="slider-label">Độ trong suốt</span>
													<span class="slider-value">{Math.round(appearance.settings.links.shadowOpacity * 100)}%</span>
												</div>
												<input 
													type="range" 
													class="slider" 
													style="background: white;"
													min="0" 
													max="100" 
													step="5" 
													value={appearance.settings.links.shadowOpacity * 100} 
													oninput={(e) => updateSetting('links.shadowOpacity', +e.currentTarget.value / 100)}
												/>
											</div>
										</div>
									{/if}
									
									<!-- Padding & Gap Sliders - Moved to end -->
									<div class="slider-group">
										<div class="slider-header">
											<span class="slider-label">Padding</span>
											<span class="slider-value">{appearance.settings.links.padding}px</span>
										</div>
										<input type="range" class="slider" min="8" max="24" step="2" value={appearance.settings.links.padding} oninput={(e) => updateSetting('links.padding', +e.currentTarget.value)}/>
									</div>
									
									<div class="slider-group">
										<div class="slider-header">
											<span class="slider-label">Gap</span>
											<span class="slider-value">{appearance.settings.links.gap}px</span>
										</div>
										<input type="range" class="slider" min="4" max="24" step="2" value={appearance.settings.links.gap} oninput={(e) => updateSetting('links.gap', +e.currentTarget.value)}/>
									</div>
								{/if}
							</div>
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
						<strong>Page:</strong>
						<div class="debug-item">Mode: {appearance.settings.page.mode}</div>
						<div class="debug-item">Text Align: {appearance.settings.page.layout.textAlign}</div>
						<div class="debug-item">Base Font Size: {appearance.settings.page.layout.baseFontSize}</div>
						<div class="debug-item">Max Width: {appearance.settings.page.layout.maxWidth}px</div>
						<div class="debug-item">Page Padding: {appearance.settings.page.layout.pagePadding}px</div>
						<div class="debug-item">Block Gap: {appearance.settings.page.layout.blockGap}px</div>
					</div>
					<div class="debug-section">
						<strong>Background:</strong>
						<div class="debug-item">Type: {appearance.settings.background.type}</div>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.background.color}"></span> Color: {appearance.settings.background.color}</div>
						<div class="debug-item">Gradient: {appearance.settings.background.gradient || 'none'}</div>
						{#if appearance.settings.background.imageUrl}
							<div class="debug-item">Image URL: {appearance.settings.background.imageUrl.substring(0, 40)}...</div>
						{/if}
						<div class="debug-item">Blur: {appearance.settings.background.effects.blur}px</div>
						<div class="debug-item">Dim: {appearance.settings.background.effects.dim}</div>
						{#if appearance.settings.background.customGradient?.enabled}
							<div class="debug-item">Custom Gradient: {appearance.settings.background.customGradient.type} {appearance.settings.background.customGradient.angle}° ({appearance.settings.background.customGradient.fromColor} → {appearance.settings.background.customGradient.toColor})</div>
						{/if}
					</div>
					<div class="debug-section">
						<strong>Colors (Semantic):</strong>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.colors.primary}"></span> Primary: {appearance.settings.colors.primary}</div>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.colors.text}"></span> Text: {appearance.settings.colors.text}</div>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.colors.textSecondary}"></span> Text Secondary: {appearance.settings.colors.textSecondary}</div>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.colors.background}"></span> Background: {appearance.settings.colors.background}</div>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.colors.cardBackground}"></span> Card BG: {appearance.settings.colors.cardBackground}</div>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.colors.border}"></span> Border: {appearance.settings.colors.border}</div>
					</div>
					<div class="debug-section">
						<strong>Links/Blocks:</strong>
						<div class="debug-item">Style: {appearance.settings.links.style}</div>
						<div class="debug-item">Text Align: {appearance.settings.links.textAlign}</div>
						<div class="debug-item">Font Size: {appearance.settings.links.fontSize}</div>
						<div class="debug-item">Radius: {appearance.settings.links.borderRadius}px</div>
						<div class="debug-item">Padding: {appearance.settings.links.padding}px</div>
						<div class="debug-item">Gap: {appearance.settings.links.gap}px</div>
						<div class="debug-item"><span class="color-dot" style="background:{getActualBlockTextColor()}"></span> Block Text Color: {getActualBlockTextColor()} {appearance.settings.links.textColor ? '(custom)' : '(auto)'}</div>
						<div class="debug-item"><span class="color-dot" style="background:{getActualBlockBgColor()}"></span> Block BG Color: {getActualBlockBgColor()}</div>
						<div class="debug-item">Show Background: {appearance.settings.links.showBackground}</div>
						<div class="debug-item">Show Border: {appearance.settings.links.showBorder}</div>
						{#if appearance.settings.links.showBorder}
							<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.links.borderColor}"></span> Border Color: {appearance.settings.links.borderColor}</div>
							<div class="debug-item">Border Width: {appearance.settings.links.borderWidth}px</div>
						{/if}
						<div class="debug-item">Show Shadow: {appearance.settings.links.showShadow}</div>
						{#if appearance.settings.links.showShadow}
							<div class="debug-item">Shadow Blur: {appearance.settings.links.shadowBlur}px</div>
							<div class="debug-item">Shadow Offset: {appearance.settings.links.shadowOffsetX}px, {appearance.settings.links.shadowOffsetY}px</div>
							<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.links.shadowColor}"></span> Shadow Color: {appearance.settings.links.shadowColor}</div>
							<div class="debug-item">Shadow Opacity: {appearance.settings.links.shadowOpacity}</div>
						{/if}
					</div>
					<div class="debug-section">
						<strong>Header:</strong>
						<div class="debug-item">Style: {appearance.settings.header.style}</div>
						<div class="debug-item">Show Avatar: {appearance.settings.header.showAvatar}</div>
						<div class="debug-item">Avatar Size: {appearance.settings.header.avatarSize}</div>
						<div class="debug-item">Avatar Shape: {appearance.settings.header.avatarShape}</div>
						{#if appearance.settings.header.avatarBorderColor}
							<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.header.avatarBorderColor}"></span> Avatar Border: {appearance.settings.header.avatarBorderWidth || 3}px {appearance.settings.header.avatarBorderColor}</div>
						{/if}
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.header.nameColor || appearance.settings.colors.text}"></span> Name Color: {appearance.settings.header.nameColor || '(auto)'}</div>
						<div class="debug-item">Show Bio: {appearance.settings.header.showBio}</div>
						<div class="debug-item">Bio Size: {appearance.settings.header.bioSize}</div>
						<div class="debug-item">Bio Align: {appearance.settings.header.bioAlign}</div>
						<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.header.bioColor || appearance.settings.colors.textSecondary}"></span> Bio Color: {appearance.settings.header.bioColor || '(auto)'}</div>
						<div class="debug-item">Show Socials: {appearance.settings.header.showSocials}</div>
						{#if appearance.settings.header.showSocials}
							<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.header.socialIconsColor || appearance.settings.colors.textSecondary}"></span> Social Icons Color: {appearance.settings.header.socialIconsColor || '(auto)'}</div>
							<div class="debug-item">Social Icons BG: {appearance.settings.header.socialIconsBg}</div>
						{/if}
						<div class="debug-item">Cover Type: {appearance.settings.header.cover.type}</div>
						{#if appearance.settings.header.cover.type === 'color'}
							<div class="debug-item"><span class="color-dot" style="background:{appearance.settings.header.cover.color}"></span> Cover Color: {appearance.settings.header.cover.color}</div>
						{:else if appearance.settings.header.cover.imageUrl}
							<div class="debug-item">Cover Image: {appearance.settings.header.cover.imageUrl.substring(0, 40)}...</div>
						{/if}
					</div>
					<div class="debug-section">
						<strong>Computed Styles (Preview):</strong>
						<div class="debug-item">BG Style: {getBgStyle().substring(0, 50)}...</div>
						<div class="debug-item">Link BG: {getLinkBg()}</div>
						<div class="debug-item">Link Color: {getLinkColor()}</div>
						<div class="debug-item">Link Border: {getLinkBorder()}</div>
						<div class="debug-item">Shadow: {getShadow()}</div>
					</div>
				</div>
				{/if}
				
				<div class="preview-frame">
					<div class="preview-content" style="background:{getBgStyle()}; text-align:center">
						{#if appearance.settings.background.effects.dim > 0}<div class="preview-overlay" style="background:rgba(0,0,0,{appearance.settings.background.effects.dim})"></div>{/if}
						
						{#if appearance.settings.header.style === 'classic'}
							<!-- Classic: Cover + Avatar Overlap Center -->
							<div class="preview-cover-classic" style="{getCoverStyle()}"></div>
							<div class="preview-avatar-classic" style="{getHeaderAvatarStyle()}">{auth.user?.email?.charAt(0).toUpperCase() || 'U'}</div>
							<div class="preview-header-classic" style="{getHeaderContainerStyle()}">
								<div class="preview-name" style="{getHeaderNameStyle()}">
									{#if auth.user?.display_name}
										{auth.user.display_name}
									{:else}
										@{auth.user?.username || 'username'}
									{/if}
								</div>
								{#if appearance.settings.header.showSocials}
									<div class="preview-social-icons">
										{#key `${appearance.settings.header.socialIconsColor}-${appearance.settings.header.socialIconsBg}`}
											<span class="preview-social-icon" style="{getSocialIconStyle()}"><Instagram size={16} /></span>
											<span class="preview-social-icon" style="{getSocialIconStyle()}"><Music size={16} /></span>
										{/key}
									</div>
								{/if}
								{#if appearance.settings.header.showBio}
									<div class="preview-bio" style="{getHeaderBioStyle()}">
										{#if bioData?.page?.settings && typeof bioData.page.settings === 'object'}
											{@const settings = bioData.page.settings as any}
											{settings.bio || 'Your bio here'}
										{:else}
											Your bio here
										{/if}
									</div>
								{/if}
							</div>
							
						{:else}
							<!-- Minimal: No Cover, Vertical Center Layout -->
							<div class="preview-header-minimal" style="padding: {appearance.settings.page.layout.pagePadding}px {appearance.settings.page.layout.pagePadding}px 0;">
								<div class="preview-avatar-minimal" style="{getHeaderAvatarStyle()}">{auth.user?.email?.charAt(0).toUpperCase() || 'U'}</div>
								<div class="preview-name" style="{getHeaderNameStyle()}">
									{#if auth.user?.display_name}
										{auth.user.display_name}
									{:else}
										@{auth.user?.username || 'username'}
									{/if}
								</div>
								{#if appearance.settings.header.showSocials}
									<div class="preview-social-icons">
										{#key `${appearance.settings.header.socialIconsColor}-${appearance.settings.header.socialIconsBg}`}
											<span class="preview-social-icon" style="{getSocialIconStyle()}"><Instagram size={16} /></span>
											<span class="preview-social-icon" style="{getSocialIconStyle()}"><Music size={16} /></span>
										{/key}
									</div>
								{/if}
								{#if appearance.settings.header.showBio}
									<div class="preview-bio" style="{getHeaderBioStyle()}">
										{#if bioData?.page?.settings && typeof bioData.page.settings === 'object'}
											{@const settings = bioData.page.settings as any}
											{settings.bio || 'Your bio here'}
										{:else}
											Your bio here
										{/if}
									</div>
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
	.btn-save { 
		display: flex; 
		align-items: center; 
		gap: var(--space-2); 
		position: relative;
		transition: all 0.3s ease;
	}
	.btn-save.has-changes {
		animation: pulse-glow 2s ease-in-out infinite;
	}
	.btn-save .spinner { 
		width: 16px; 
		height: 16px; 
		border: 2px solid rgba(255,255,255,0.3); 
		border-top-color: white; 
		border-radius: 50%; 
		animation: spin 0.8s linear infinite; 
	}
	.change-indicator {
		position: absolute;
		top: -4px;
		right: -4px;
		width: 10px;
		height: 10px;
		background: #ff9500;
		border: 2px solid white;
		border-radius: 50%;
		animation: pulse-dot 2s ease-in-out infinite;
	}
	
	@keyframes pulse-glow {
		0%, 100% { box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1); }
		50% { box-shadow: 0 2px 16px rgba(255, 149, 0, 0.4); }
	}
	
	@keyframes pulse-dot {
		0%, 100% { transform: scale(1); opacity: 1; }
		50% { transform: scale(1.2); opacity: 0.8; }
	}

	.layout { display: grid; grid-template-columns: 160px 1fr 360px; gap: var(--space-5); }
	.sidebar { display: flex; flex-direction: column; gap: var(--space-1); position: sticky; top: 80px; height: fit-content; }
	.nav-item { display: flex; align-items: center; gap: var(--space-3); padding: var(--space-3); background: transparent; border-radius: var(--radius-md); text-align: left; font-size: var(--text-sm); font-weight: 500; }
	.nav-item:hover { background: var(--color-bg); }
	.nav-item.active { background: var(--color-primary-light); color: var(--color-primary); }
	.nav-icon { display: flex; align-items: center; justify-content: center; width: 20px; }
	.content { display: flex; flex-direction: column; gap: var(--space-4); }

	.theme-card { position: relative; background: transparent; padding: 0; text-align: center; }
	.theme-card .badge { position: absolute; top: 4px; right: 4px; }
	
	/* Theme preview button - consistent with preset cards */
	.theme-preview-btn {
		width: 100%;
		background: transparent;
		border: none;
		padding: 0;
		cursor: pointer;
		display: block;
		text-align: center;
	}
	
	/* My Theme Card Styling - đồng bộ với preset cards */
	.my-theme-card {
		position: relative;
	}
	
	.my-theme-card .preview-item {
		position: relative;
		border: 2px solid transparent;
		transition: all 0.2s ease;
	}
	
	.my-theme-card .preview-item.active {
		border-color: var(--color-primary);
		box-shadow: 0 0 0 3px var(--color-primary-light);
	}
	
	.my-theme-card:hover .preview-item {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	}
	
	/* Custom badge - đẹp và nổi bật */
	.custom-badge {
		position: absolute;
		top: 8px;
		left: 8px;
		display: flex;
		align-items: center;
		gap: 4px;
		background: linear-gradient(135deg, #007aff 0%, #5856d6 100%);
		color: white;
		padding: 4px 10px;
		border-radius: 12px;
		font-size: 10px;
		font-weight: 700;
		letter-spacing: 0.5px;
		box-shadow: 0 2px 8px rgba(0, 122, 255, 0.4);
		z-index: 2;
	}
	
	.custom-badge svg {
		width: 12px;
		height: 12px;
		flex-shrink: 0;
	}
	
	/* Theme label - đồng bộ với preset cards */
	.theme-label {
		padding: 8px 0;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 4px;
	}
	
	.my-theme-card.active .theme-label {
		color: var(--color-primary);
		font-weight: 600;
	}
	
	/* Delete button - floating style */
	.btn-delete-theme {
		position: absolute;
		top: 8px;
		right: 8px;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 24px;
		height: 24px;
		background: rgba(255, 255, 255, 0.95);
		backdrop-filter: blur(8px);
		border: 1px solid rgba(0, 0, 0, 0.1);
		border-radius: 6px;
		color: var(--color-text-secondary);
		cursor: pointer;
		transition: all 0.2s ease;
		z-index: 3;
		opacity: 0;
	}
	
	.my-theme-card:hover .btn-delete-theme {
		opacity: 1;
	}
	
	.btn-delete-theme:hover {
		background: var(--color-danger);
		border-color: var(--color-danger);
		color: white;
		transform: scale(1.1);
	}

	.color-scheme-card { background: transparent; padding: var(--space-3); text-align: center; border: 2px solid transparent; border-radius: var(--radius-md); cursor: pointer; transition: all 0.2s; }
	.color-scheme-card:hover { border-color: var(--color-separator); background: var(--color-bg); }
	.color-scheme-preview { display: flex; align-items: center; justify-content: center; gap: var(--space-2); margin-bottom: var(--space-2); }
	.color-dot-large { width: 48px; height: 48px; border-radius: 50%; border: 2px solid rgba(0,0,0,0.1); }
	.color-dots-small { display: flex; flex-direction: column; gap: 4px; }
	.color-dot-small { width: 20px; height: 20px; border-radius: 4px; border: 1px solid rgba(0,0,0,0.1); }

	.advanced-toggle { background: transparent; border: none; color: var(--color-text-secondary); font-size: var(--text-sm); font-weight: 500; cursor: pointer; padding: 0; display: flex; align-items: center; gap: var(--space-2); }
	.advanced-toggle:hover { color: var(--color-text); }

	/* Block Style Grid - Modern Context-Aware Design */
	.block-style-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 12px;
	}

	.block-style-card {
		background: transparent;
		border: 3px solid transparent;
		border-radius: 12px;
		padding: 0;
		cursor: pointer;
		transition: all 0.2s ease;
		overflow: hidden;
	}

	.block-style-card:hover {
		border-color: var(--color-separator);
		transform: translateY(-2px);
	}

	.block-style-card.active {
		border-color: var(--color-primary);
		box-shadow: 0 0 0 3px var(--color-primary-light);
	}

	.block-style-preview {
		width: 100%;
		height: 120px;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 16px;
		border-radius: 9px 9px 0 0;
		position: relative;
		overflow: hidden;
	}

	.block-style-demo {
		width: 100%;
		padding: 12px 16px;
		border-radius: 8px;
		font-size: 13px;
		font-weight: 500;
		text-align: center;
		transition: all 0.2s ease;
		position: relative;
		z-index: 1;
	}

	.block-style-label {
		display: block;
		padding: 10px;
		font-size: 13px;
		font-weight: 500;
		color: var(--color-text);
		text-align: center;
		background: var(--color-bg);
	}

	.block-style-card.active .block-style-label {
		color: var(--color-primary);
		font-weight: 600;
	}

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
		position: relative;
	}
	
	/* Add indicator to all active preset buttons */
	.color-preset-circle.active::after {
		content: '●';
		position: absolute;
		top: -6px;
		right: -6px;
		width: 16px;
		height: 16px;
		background: var(--color-primary);
		color: white;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 10px;
		line-height: 1;
		border: 2px solid white;
		box-shadow: 0 2px 4px rgba(0,0,0,0.2);
	}
	
	.color-preset-circle.current-color {
		border-color: var(--color-primary);
		box-shadow: 0 0 0 3px var(--color-primary-light);
	}
	
	.current-indicator {
		position: absolute;
		top: -6px;
		right: -6px;
		width: 16px;
		height: 16px;
		background: var(--color-primary);
		color: white;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 10px;
		line-height: 1;
		border: 2px solid white;
		box-shadow: 0 2px 4px rgba(0,0,0,0.2);
	}
	
	/* Gradient preview items */
	.preview-item {
		position: relative;
	}
	
	.preview-item.active::after {
		content: '●';
		position: absolute;
		top: -6px;
		right: -6px;
		width: 16px;
		height: 16px;
		background: var(--color-primary);
		color: white;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 10px;
		line-height: 1;
		border: 2px solid white;
		box-shadow: 0 2px 4px rgba(0,0,0,0.2);
		z-index: 1;
	}
	
	/* iOS-style Gradient Presets Grid */
	.gradient-presets-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
		gap: 12px;
		margin-bottom: var(--space-4);
	}
	
	.gradient-presets-grid .preview-item {
		height: 120px;
		border-radius: 20px;
		border: none;
		cursor: pointer;
		transition: all 0.2s ease;
		box-shadow: 0 2px 8px rgba(0,0,0,0.08);
		position: relative;
		overflow: visible;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.gradient-presets-grid .preview-item:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 16px rgba(0,0,0,0.15);
	}
	
	.gradient-presets-grid .preview-item.active {
		box-shadow: 0 0 0 4px var(--color-primary), 0 4px 20px rgba(0,0,0,0.2);
		transform: scale(1.02);
	}
	
	.gradient-presets-grid .custom-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;
		color: #999;
		transition: all 0.2s ease;
		position: relative;
		z-index: 1;
	}
	
	.gradient-presets-grid .preview-item:hover .custom-icon {
		color: #666;
		transform: scale(1.1);
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

	/* Info box */
	.info-box {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-3);
		background: rgba(255, 149, 0, 0.1);
		border: 1px solid rgba(255, 149, 0, 0.3);
		border-radius: var(--radius-md);
		color: #f59e0b;
		font-size: var(--text-sm);
		margin-bottom: var(--space-4);
	}

	/* Button group inline */
	.button-group-inline {
		display: flex;
		gap: var(--space-2);
	}

	.btn-option {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2) var(--space-3);
		background: var(--color-bg);
		border: 1px solid var(--color-separator);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		color: var(--color-text-secondary);
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-option:hover {
		background: var(--color-bg-secondary);
		border-color: var(--color-primary);
	}

	.btn-option.active {
		background: var(--color-primary);
		color: white;
		border-color: var(--color-primary);
	}

	/* Upload area */
	.upload-area {
		margin-top: var(--space-2);
	}

	.upload-label {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: var(--space-2);
		padding: var(--space-6);
		background: var(--color-bg);
		border: 2px dashed var(--color-separator);
		border-radius: var(--radius-lg);
		cursor: pointer;
		transition: all 0.2s;
		color: var(--color-text-secondary);
	}

	.upload-label:hover {
		border-color: var(--color-primary);
		background: var(--color-bg-secondary);
	}

	.upload-hint {
		font-size: var(--text-xs);
		color: var(--color-text-tertiary);
	}

	.uploaded-image {
		position: relative;
		border-radius: var(--radius-lg);
		overflow: hidden;
		aspect-ratio: 3 / 1;
	}

	.uploaded-image img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.remove-image {
		position: absolute;
		top: var(--space-2);
		right: var(--space-2);
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(0, 0, 0, 0.6);
		color: white;
		border-radius: 50%;
		cursor: pointer;
		transition: all 0.2s;
	}

	.remove-image:hover {
		background: rgba(255, 59, 48, 0.9);
	}

	/* Avatar Shape Grid - Visual Preview Cards */
	.setting-section {
		margin-bottom: var(--space-5);
	}

	.avatar-shape-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 12px;
		margin-top: var(--space-3);
	}

	.avatar-shape-card {
		background: var(--color-bg);
		border: 2px solid var(--color-separator);
		border-radius: 12px;
		padding: var(--space-4);
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-3);
	}

	.avatar-shape-card:hover {
		border-color: var(--color-primary);
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
	}

	.avatar-shape-card.active {
		border-color: var(--color-primary);
		border-width: 3px;
		background: var(--color-primary-light);
		box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.1);
	}

	.avatar-preview-wrapper {
		width: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
		padding: var(--space-2);
	}

	.avatar-preview {
		width: 80px;
		height: 80px;
		position: relative;
		overflow: hidden;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
		transition: all 0.2s ease;
	}

	.avatar-preview.circle {
		border-radius: 50%;
	}

	.avatar-preview.rounded {
		border-radius: 20px;
	}

	.avatar-preview.square {
		border-radius: 8px;
	}

	.avatar-demo-image {
		width: 100%;
		height: 100%;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		display: flex;
		align-items: center;
		justify-content: center;
		position: relative;
	}

	.avatar-demo-image::before {
		content: '';
		position: absolute;
		width: 30px;
		height: 30px;
		background: rgba(255, 255, 255, 0.3);
		border-radius: 50%;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -60%);
	}

	.avatar-demo-image::after {
		content: '';
		position: absolute;
		width: 50px;
		height: 35px;
		background: rgba(255, 255, 255, 0.3);
		border-radius: 50% 50% 0 0;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
	}

	.shape-label {
		font-size: var(--text-sm);
		font-weight: 600;
		color: var(--color-text);
		text-align: center;
	}

	.avatar-shape-card.active .shape-label {
		color: var(--color-primary);
	}

	/* Avatar Size Options */
	.avatar-size-options {
		display: flex;
		gap: 12px;
		margin-top: var(--space-3);
		justify-content: flex-start;
	}

	.avatar-size-btn {
		background: var(--color-bg);
		border: 2px solid var(--color-separator);
		border-radius: 12px;
		padding: var(--space-3);
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-2);
		min-width: 90px;
	}

	.avatar-size-btn:hover {
		border-color: var(--color-primary);
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
	}

	.avatar-size-btn.active {
		border-color: var(--color-primary);
		border-width: 3px;
		background: var(--color-primary-light);
		box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.1);
	}

	.size-preview {
		position: relative;
		overflow: hidden;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
		transition: all 0.2s ease;
	}

	.size-preview.s {
		width: 40px;
		height: 40px;
	}

	.size-preview.m {
		width: 56px;
		height: 56px;
	}

	.size-preview.l {
		width: 72px;
		height: 72px;
	}

	.size-preview.circle {
		border-radius: 50%;
	}

	.size-preview.rounded {
		border-radius: 16px;
	}

	.size-preview.square {
		border-radius: 6px;
	}

	.size-demo-image {
		width: 100%;
		height: 100%;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		display: flex;
		align-items: center;
		justify-content: center;
		position: relative;
	}

	.size-demo-image::before {
		content: '';
		position: absolute;
		width: 35%;
		height: 35%;
		background: rgba(255, 255, 255, 0.3);
		border-radius: 50%;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -60%);
	}

	.size-demo-image::after {
		content: '';
		position: absolute;
		width: 60%;
		height: 40%;
		background: rgba(255, 255, 255, 0.3);
		border-radius: 50% 50% 0 0;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
	}

	.size-label {
		font-size: var(--text-sm);
		font-weight: 600;
		color: var(--color-text);
		text-align: center;
	}

	.avatar-size-btn.active .size-label {
		color: var(--color-primary);
	}
</style>


