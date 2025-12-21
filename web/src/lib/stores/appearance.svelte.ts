import { themes, pages, type ThemePreset, type Page, type ThemeCustom } from '$lib/api/client';

// Appearance settings state
let loading = $state(true);
let saving = $state(false);
let presets = $state<ThemePreset[]>([]);
let customTheme = $state<ThemeCustom | null>(null);
let currentPage = $state<Page | null>(null);
let selectedPresetId = $state<number>(1);
let originalPresetId = $state<number>(1);
let customPatch = $state<Record<string, any>>({});
let originalPatch = $state<Record<string, any>>({});
let dirty = $state(false);
let isUsingCustom = $state(false); // Track if currently using custom theme

// Default appearance values (fallback khi preset không có)
const defaultAppearance = {
	page: {
		layout: {
			maxWidth: 520,
			pagePadding: 16,
			blockGap: 12,
			textAlign: 'center' as 'left' | 'center' | 'right',
			baseFontSize: 'M' as 'S' | 'M' | 'L' | 'XL'
		},
		mode: 'light' as 'light' | 'dark'
	},
	background: {
		type: 'solid' as 'solid' | 'gradient' | 'image',
		color: '#f2f2f7',
		gradient: '',
		imageUrl: '',
		effects: { blur: 0, dim: 0 },
		customGradient: {
			enabled: false,
			type: 'linear' as 'linear' | 'radial',
			angle: 135,
			fromColor: '#667eea',
			toColor: '#764ba2'
		}
	},
	header: {
		style: 'minimal' as 'minimal' | 'classic',
		showAvatar: true,
		showBio: true,
		showSocials: true,
		avatarSize: 'M' as 'S' | 'M' | 'L',
		avatarShape: 'circle' as 'circle' | 'rounded' | 'square',
		avatarBorderColor: '' as string, // empty = use default white
		avatarBorderWidth: 3, // border width in pixels
		nameColor: '' as string, // empty = use text from colors
		bioSize: 'M' as 'S' | 'M' | 'L' | 'XL',
		bioAlign: 'center' as 'left' | 'center' | 'right',
		bioColor: '' as string, // empty = use textSecondary from colors
		socialIconsColor: '' as string, // empty = use textSecondary
		socialIconsBg: true, // show background by default
		cover: {
			type: 'color' as 'color' | 'image',
			color: '#c2185b',
			imageUrl: null as string | null,
			imageAssetId: null as number | null
		}
	},
	links: {
		style: 'filled' as 'filled' | 'outline' | 'elevated' | 'shadow',
		borderRadius: 12,
		shadow: 'sm' as 'none' | 'sm' | 'md' | 'lg', // kept for backward compatibility
		padding: 16,
		gap: 12,
		textAlign: 'center' as 'left' | 'center' | 'right',
		fontSize: 'M' as 'S' | 'M' | 'L',
		textColor: '' as string, // empty = auto based on style
		showBackground: true, // toggle to show/hide block background
		showBorder: false, // toggle to show/hide border
		borderColor: '#e5e5ea', // border color when showBorder is true
		borderWidth: 1, // border width in pixels
		showShadow: false, // toggle to show/hide custom shadow
		shadowBlur: 8, // shadow blur radius in pixels
		shadowOffsetX: 0, // shadow horizontal offset
		shadowOffsetY: 2, // shadow vertical offset
		shadowColor: '#000000', // shadow color
		shadowOpacity: 0.12 // shadow opacity (0-1)
	},
	colors: {
		primary: '#007aff',
		text: '#1c1c1e',
		textSecondary: '#8e8e93',
		background: '#f2f2f7',
		cardBackground: '#ffffff',
		border: 'rgba(60, 60, 67, 0.1)'
	}
};

export type AppearanceSettings = typeof defaultAppearance;

// Build gradient CSS string from config
export function buildGradientString(config: {
	type: 'linear' | 'radial';
	angle: number;
	fromColor: string;
	toColor: string;
}): string {
	if (config.type === 'radial') {
		return `radial-gradient(circle, ${config.fromColor} 0%, ${config.toColor} 100%)`;
	}
	return `linear-gradient(${config.angle}deg, ${config.fromColor} 0%, ${config.toColor} 100%)`;
}

// Lấy config từ preset theo ID
function getPresetConfigById(presetId: number): Record<string, any> {
	const preset = presets.find(p => p.id === presetId);
	if (!preset?.config) return {};
	return preset.config as Record<string, any>;
}

// Helper: parse CSS size value to pixels (vd: "16px" -> 16, "1rem" -> 16)
function parseSize(value: any): number | undefined {
	if (typeof value === 'number') return value;
	if (typeof value === 'string') {
		// Lấy số đầu tiên trong string (bỏ qua phần sau như "1rem 1.25rem")
		const match = value.match(/^([\d.]+)/);
		if (!match) return undefined;
		
		const num = parseFloat(match[1]);
		if (isNaN(num)) return undefined;
		
		// Convert rem to px (1rem = 16px)
		if (value.includes('rem')) {
			return Math.round(num * 16);
		}
		// px hoặc số thuần
		return Math.round(num);
	}
	return undefined;
}

// Helper: parse border string to extract width and color
// Example: "1px solid #e5e5ea" -> { width: 1, color: '#e5e5ea' }
function parseBorder(borderValue: string | undefined): { width: number; color: string } | null {
	if (!borderValue || borderValue === 'none') return null;
	
	// Match pattern: "1px solid #e5e5ea" or "2px solid rgba(0,0,0,0.1)"
	const match = borderValue.match(/(\d+)px\s+solid\s+(#[0-9a-fA-F]{6}|rgba?\([^)]+\))/);
	if (match) {
		return {
			width: parseInt(match[1]),
			color: match[2].startsWith('#') ? match[2] : '#e5e5ea' // fallback for rgba
		};
	}
	return null;
}

// Helper: convert shadow level to custom shadow values
function shadowLevelToCustom(level: string): { blur: number; offsetX: number; offsetY: number; opacity: number } | null {
	const shadowMap: Record<string, { blur: number; offsetX: number; offsetY: number; opacity: number }> = {
		'sm': { blur: 2, offsetX: 0, offsetY: 1, opacity: 0.08 },
		'md': { blur: 8, offsetX: 0, offsetY: 2, opacity: 0.12 },
		'lg': { blur: 16, offsetX: 0, offsetY: 4, opacity: 0.16 }
	};
	
	return shadowMap[level] || null;
}

// Helper: parse shadow CSS string to extract values
// Example: "0 2px 8px rgba(0,0,0,0.12)" -> { blur: 8, offsetX: 0, offsetY: 2, opacity: 0.12 }
// Example with spread: "0 4px 6px -1px rgba(0,0,0,0.1)" -> { blur: 6, offsetX: 0, offsetY: 4, opacity: 0.1 }
function parseShadow(shadowValue: string | undefined): { blur: number; offsetX: number; offsetY: number; opacity: number; color: string } | null {
	if (!shadowValue || shadowValue === 'none') return null;
	
	// Match pattern with optional spread: "0 4px 6px -1px rgba(0,0,0,0.1)" or "0px 4px 6px -1px rgba(0,0,0,0.1)"
	// Note: px is optional for each value
	const matchWithSpread = shadowValue.match(/([-\d]+)(?:px)?\s+([-\d]+)(?:px)?\s+([-\d]+)(?:px)?\s+([-\d]+)(?:px)?\s+rgba?\((\d+),\s*(\d+),\s*(\d+),\s*([\d.]+)\)/);
	if (matchWithSpread) {
		return {
			offsetX: parseInt(matchWithSpread[1]),
			offsetY: parseInt(matchWithSpread[2]),
			blur: parseInt(matchWithSpread[3]),
			// spread (matchWithSpread[4]) is ignored for now - we don't have UI for it
			color: `#${parseInt(matchWithSpread[5]).toString(16).padStart(2, '0')}${parseInt(matchWithSpread[6]).toString(16).padStart(2, '0')}${parseInt(matchWithSpread[7]).toString(16).padStart(2, '0')}`,
			opacity: parseFloat(matchWithSpread[8])
		};
	}
	
	// Match pattern without spread: "0 2px 8px rgba(0,0,0,0.12)" or "0px 2px 8px rgba(0,0,0,0.12)"
	const match = shadowValue.match(/([-\d]+)(?:px)?\s+([-\d]+)(?:px)?\s+([-\d]+)(?:px)?\s+rgba?\((\d+),\s*(\d+),\s*(\d+),\s*([\d.]+)\)/);
	if (match) {
		return {
			offsetX: parseInt(match[1]),
			offsetY: parseInt(match[2]),
			blur: parseInt(match[3]),
			color: `#${parseInt(match[4]).toString(16).padStart(2, '0')}${parseInt(match[5]).toString(16).padStart(2, '0')}${parseInt(match[6]).toString(16).padStart(2, '0')}`,
			opacity: parseFloat(match[7])
		};
	}
	return null;
}

// Helper: map shadow string to level
function mapShadowToLevel(shadow: string | undefined): 'none' | 'sm' | 'md' | 'lg' {
	if (!shadow || shadow === 'none') return 'none';
	if (shadow.includes('20px') || shadow.includes('24px') || shadow.includes('16px')) return 'lg';
	if (shadow.includes('8px') || shadow.includes('12px') || shadow.includes('6px')) return 'md';
	return 'sm';
}

// Chuyển đổi preset config sang appearance settings format
function mapPresetToSettings(presetConfig: Record<string, any>): Partial<AppearanceSettings> {
	const result: Partial<AppearanceSettings> = {};

	// Map background
	if (presetConfig.background) {
		const bg = presetConfig.background;
		result.background = {
			type: bg.type || (bg.gradient ? 'gradient' : 'solid'),
			color: bg.color || '#f2f2f7',
			gradient: bg.gradient || '',
			imageUrl: bg.imageUrl || '',
			effects: {
				blur: bg.effects?.blur ?? 0,
				dim: bg.effects?.dim ?? 0
			},
			customGradient: {
				enabled: false,
				type: 'linear',
				angle: 135,
				fromColor: '#667eea',
				toColor: '#764ba2'
			}
		};
	}

	// Map page layout
	if (presetConfig.page?.layout) {
		const layout = presetConfig.page.layout;
		result.page = {
			layout: {
				maxWidth: layout.maxWidth ?? 520,
				pagePadding: layout.pagePadding ?? 16,
				blockGap: layout.blockGap ?? 12,
				textAlign: layout.textAlign ?? 'center',
				baseFontSize: layout.baseFontSize ?? 'M'
			},
			mode: presetConfig.page?.mode ?? 'light'
		};
	}

	// Map colors từ semantic
	if (presetConfig.semantic?.color) {
		const sc = presetConfig.semantic.color;
		result.colors = {
			primary: sc.primary || '#007aff',
			text: sc.text?.default || '#1c1c1e',
			textSecondary: sc.text?.muted || '#8e8e93',
			background: sc.surface?.page || '#f2f2f7',
			cardBackground: sc.surface?.card || '#ffffff',
			border: sc.border?.default || 'rgba(60, 60, 67, 0.1)'
		};
	}

	// Map link styles từ recipes.linkItem
	if (presetConfig.recipes?.linkItem?.base) {
		const base = presetConfig.recipes.linkItem.base;
		const linkDefaults = presetConfig.page?.defaults?.linkGroup || {};

		result.links = {
			style: 'filled',
			borderRadius: parseSize(base.borderRadius) ?? 12,
			shadow: mapShadowToLevel(base.shadow),
			padding: parseSize(base.padding) ?? 16,
			gap: parseSize(linkDefaults.gap) ?? 12,
			textAlign: linkDefaults.textAlign || 'center',
			fontSize: linkDefaults.fontSize || 'M',
			textColor: linkDefaults.textColor || '', // Map textColor from theme preset
			showBackground: true,
			showBorder: false,
			borderColor: '#e5e5ea',
			borderWidth: 1,
			showShadow: false,
			shadowBlur: 8,
			shadowOffsetX: 0,
			shadowOffsetY: 2,
			shadowColor: '#000000',
			shadowOpacity: 0.12
		};

		// Nếu có color trong linkItem, cập nhật colors.cardBackground
		if (base.background && result.colors) {
			result.colors.cardBackground = base.background;
		}
		if (base.color && result.colors) {
			// Nếu link có màu text riêng, có thể dùng cho text chính
			// Nhưng thường ta giữ nguyên từ semantic
		}
		
		// ===== DETECTION LOGIC: Auto-enable toggles based on preset values =====
		
		// 1. SHADOW DETECTION
		if (base.shadow) {
			const shadowValue = base.shadow;
			
			// Nếu có shadow string (CSS format) → parse và enable
			const parsedShadow = parseShadow(shadowValue);
			if (parsedShadow) {
				result.links.showShadow = true;
				result.links.shadowBlur = parsedShadow.blur;
				result.links.shadowOffsetX = parsedShadow.offsetX;
				result.links.shadowOffsetY = parsedShadow.offsetY;
				result.links.shadowColor = parsedShadow.color;
				result.links.shadowOpacity = parsedShadow.opacity;
			} else {
				// Nếu là shadow level ('sm', 'md', 'lg') → convert và enable
				const shadowLevel = mapShadowToLevel(shadowValue);
				if (shadowLevel && shadowLevel !== 'none') {
					const customShadow = shadowLevelToCustom(shadowLevel);
					if (customShadow) {
						result.links.showShadow = true;
						result.links.shadowBlur = customShadow.blur;
						result.links.shadowOffsetX = customShadow.offsetX;
						result.links.shadowOffsetY = customShadow.offsetY;
						result.links.shadowOpacity = customShadow.opacity;
						result.links.shadowColor = '#000000';
					}
				}
			}
		}
		
		// 2. BORDER DETECTION
		if (base.border) {
			const parsedBorder = parseBorder(base.border);
			if (parsedBorder) {
				result.links.showBorder = true;
				result.links.borderWidth = parsedBorder.width;
				result.links.borderColor = parsedBorder.color;
			}
		}
		
		// 3. BACKGROUND DETECTION
		if (base.background) {
			// Nếu background không phải transparent hoặc none → enable
			if (base.background !== 'transparent' && base.background !== 'none') {
				result.links.showBackground = true;
			}
		}
	}

	// Map từ page.defaults.linkGroup (override nếu có)
	if (presetConfig.page?.defaults?.linkGroup) {
		const lg = presetConfig.page.defaults.linkGroup;
		if (!result.links) result.links = { ...defaultAppearance.links };
		if (lg.textAlign) result.links.textAlign = lg.textAlign;
		if (lg.fontSize) result.links.fontSize = lg.fontSize;
		if (lg.radius) result.links.borderRadius = parseSize(lg.radius) ?? result.links.borderRadius;
		if (lg.shadow) result.links.shadow = lg.shadow as any;
	}

	return result;
}

// Tính toán settings cuối cùng: default → preset → customPatch
function computeSettings(): AppearanceSettings {
	// Đọc selectedPresetId trực tiếp để đảm bảo reactivity
	const currentPresetId = selectedPresetId;
	const currentPatch = customPatch;
	
	const presetConfig = getPresetConfigById(currentPresetId);
	const presetSettings = mapPresetToSettings(presetConfig);
	
	// Merge: default → preset → customPatch
	return deepMerge(
		deepMerge(defaultAppearance, presetSettings),
		currentPatch
	) as AppearanceSettings;
}

export function getAppearance() {
	return {
		get loading() { return loading; },
		get saving() { return saving; },
		get presets() { return presets; },
		get customTheme() { return customTheme; },
		get currentPage() { return currentPage; },
		get selectedPresetId() { return selectedPresetId; },
		get customPatch() { return customPatch; },
		get dirty() { return dirty; },
		get isUsingCustom() { return isUsingCustom; },
		get settings(): AppearanceSettings {
			return computeSettings();
		},
		// Lấy preset hiện tại
		get selectedPreset(): ThemePreset | undefined {
			return presets.find(p => p.id === selectedPresetId);
		}
	};
}

export async function loadAppearance(pageId: number) {
	loading = true;
	try {
		const [presetsData, customThemeData, pageData] = await Promise.all([
			themes.listPresets(),
			themes.getUserCustom(),
			pages.get(pageId)
		]);
		
		presets = presetsData;
		customTheme = customThemeData;
		currentPage = pageData;
		
		// Determine if using custom theme
		if (pageData.theme_custom_id && customThemeData) {
			isUsingCustom = true;
			selectedPresetId = customThemeData.based_on_preset_id;
			customPatch = (customThemeData.patch as Record<string, any>) || {};
			originalPatch = { ...customPatch };
		} else {
			isUsingCustom = false;
			selectedPresetId = pageData.theme_preset_id;
			originalPresetId = pageData.theme_preset_id;
			customPatch = {};
			originalPatch = {};
		}
		
		dirty = false;
	} catch (e) {
		console.error('Failed to load appearance:', e);
	} finally {
		loading = false;
	}
}

// Chọn preset - reset customPatch để dùng preset mặc định
export function selectPreset(presetId: number) {
	if (selectedPresetId === presetId && !isUsingCustom) return;
	
	selectedPresetId = presetId;
	isUsingCustom = false;
	customPatch = {};
	checkDirty();
}

// Select custom theme
export function selectCustomTheme() {
	if (!customTheme) return;
	
	isUsingCustom = true;
	selectedPresetId = customTheme.based_on_preset_id;
	customPatch = (customTheme.patch as Record<string, any>) || {};
	checkDirty();
}

// Update setting - lưu vào customPatch
export function updateSetting(path: string, value: any) {
	const keys = path.split('.');
	const newPatch = { ...customPatch };
	
	let current: any = newPatch;
	for (let i = 0; i < keys.length - 1; i++) {
		if (!current[keys[i]]) {
			current[keys[i]] = {};
		} else {
			current[keys[i]] = { ...current[keys[i]] };
		}
		current = current[keys[i]];
	}
	current[keys[keys.length - 1]] = value;
	
	customPatch = newPatch;
	
	// When user makes changes, they're now using custom theme
	if (Object.keys(newPatch).length > 0) {
		isUsingCustom = true;
	}
	
	checkDirty();
}

function checkDirty() {
	// Check if using custom theme changed
	const wasUsingCustom = currentPage?.theme_custom_id != null;
	if (isUsingCustom !== wasUsingCustom) {
		dirty = true;
		return;
	}
	
	// If using custom, check if patch changed
	if (isUsingCustom) {
		dirty = JSON.stringify(customPatch) !== JSON.stringify(originalPatch);
		return;
	}
	
	// If using preset, check if preset changed
	if (selectedPresetId !== originalPresetId) {
		dirty = true;
		return;
	}
	
	dirty = false;
}

// Reset to original values (undo changes since last save)
export function resetAppearance() {
	if (currentPage?.theme_custom_id && customTheme) {
		isUsingCustom = true;
		selectedPresetId = customTheme.based_on_preset_id;
		customPatch = { ...originalPatch };
	} else {
		isUsingCustom = false;
		selectedPresetId = originalPresetId;
		customPatch = {};
	}
	dirty = false;
}

// Reset to pure preset defaults (clear all customizations)
export function resetToPresetDefaults() {
	customPatch = {};
	checkDirty();
}

// Save all changes to server
export async function saveAppearance(): Promise<boolean> {
	if (!currentPage || !dirty) return true;
	
	saving = true;
	try {
		let customThemeId: number | null = null;
		
		// Only create custom theme when using custom with changes
		if (isUsingCustom && Object.keys(customPatch).length > 0) {
			const result = await themes.createCustom(selectedPresetId, customPatch);
			customThemeId = result.id;
			customTheme = result;
		} else {
			// Using preset only → set custom_id = null
			customThemeId = null;
		}
		
		// Update page with theme settings
		await pages.save(currentPage.id, {
			page: {
				...currentPage,
				theme_preset_id: selectedPresetId,
				theme_custom_id: customThemeId || undefined,
				settings: {}
			}
		});
		
		// Update original values after successful save
		originalPresetId = selectedPresetId;
		originalPatch = isUsingCustom ? { ...customPatch } : {};
		
		// Reload to get fresh data
		await loadAppearance(currentPage.id);
		
		dirty = false;
		return true;
	} catch (e) {
		console.error('Failed to save appearance:', e);
		return false;
	} finally {
		saving = false;
	}
}

// Delete custom theme
export async function deleteCustomTheme(): Promise<boolean> {
	if (!customTheme || !currentPage) return false;
	
	saving = true;
	try {
		await themes.deleteCustom(customTheme.id);
		
		// Switch to preset
		await pages.save(currentPage.id, {
			page: {
				...currentPage,
				theme_preset_id: selectedPresetId,
				theme_custom_id: undefined
			}
		});
		
		// Reload
		await loadAppearance(currentPage.id);
		
		return true;
	} catch (e) {
		console.error('Failed to delete custom theme:', e);
		return false;
	} finally {
		saving = false;
	}
}

function deepMerge(target: any, source: any): any {
	if (!source) return target;
	const result = { ...target };
	for (const key of Object.keys(source)) {
		if (source[key] && typeof source[key] === 'object' && !Array.isArray(source[key])) {
			result[key] = deepMerge(target[key] || {}, source[key]);
		} else {
			result[key] = source[key];
		}
	}
	return result;
}
