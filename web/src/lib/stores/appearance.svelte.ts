import { themes, pages, type ThemePreset, type Page } from '$lib/api/client';

// Appearance settings state
let loading = $state(true);
let saving = $state(false);
let presets = $state<ThemePreset[]>([]);
let currentPage = $state<Page | null>(null);
let selectedPresetId = $state<number>(1);
let originalPresetId = $state<number>(1);
let customPatch = $state<Record<string, any>>({});
let originalPatch = $state<Record<string, any>>({});
let dirty = $state(false);

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
		style: 'default' as 'default' | 'minimal' | 'centered' | 'card',
		showAvatar: true,
		showBio: true,
		showSocials: true,
		avatarSize: 'M' as 'S' | 'M' | 'L',
		avatarShape: 'circle' as 'circle' | 'rounded' | 'square',
		bioSize: 'M' as 'S' | 'M' | 'L' | 'XL',
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
		style: 'filled' as 'filled' | 'outline' | 'soft' | 'ghost',
		borderRadius: 12,
		shadow: 'sm' as 'none' | 'sm' | 'md' | 'lg',
		padding: 16,
		gap: 12,
		textAlign: 'center' as 'left' | 'center' | 'right',
		fontSize: 'M' as 'S' | 'M' | 'L',
		textColor: '' as string // empty = auto based on style
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
			textColor: ''
		};

		// Nếu có color trong linkItem, cập nhật colors.cardBackground
		if (base.background && result.colors) {
			result.colors.cardBackground = base.background;
		}
		if (base.color && result.colors) {
			// Nếu link có màu text riêng, có thể dùng cho text chính
			// Nhưng thường ta giữ nguyên từ semantic
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
		get currentPage() { return currentPage; },
		get selectedPresetId() { return selectedPresetId; },
		get customPatch() { return customPatch; },
		get dirty() { return dirty; },
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
		const [presetsData, pageData] = await Promise.all([
			themes.listPresets(),
			pages.get(pageId)
		]);
		
		presets = presetsData;
		currentPage = pageData;
		selectedPresetId = pageData.theme_preset_id;
		originalPresetId = pageData.theme_preset_id;
		
		// Load custom patch from page settings if exists
		if (pageData.settings && typeof pageData.settings === 'object') {
			const settings = pageData.settings as Record<string, any>;
			if (settings.appearance) {
				customPatch = { ...settings.appearance };
				originalPatch = { ...settings.appearance };
			} else {
				customPatch = {};
				originalPatch = {};
			}
		} else {
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
	if (selectedPresetId === presetId) return;
	
	selectedPresetId = presetId;
	// Reset custom patch khi đổi preset để preview hiển thị đúng preset
	customPatch = {};
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
	checkDirty();
}

function checkDirty() {
	// Check if preset changed
	if (selectedPresetId !== originalPresetId) {
		dirty = true;
		return;
	}
	// Check if patch changed
	dirty = JSON.stringify(customPatch) !== JSON.stringify(originalPatch);
}

// Reset to original values (undo changes since last save)
export function resetAppearance() {
	selectedPresetId = originalPresetId;
	customPatch = { ...originalPatch };
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
		const updatedSettings = {
			...(currentPage.settings || {}),
			appearance: customPatch
		};
		
		await pages.save(currentPage.id, {
			page: {
				...currentPage,
				theme_preset_id: selectedPresetId,
				settings: updatedSettings
			}
		});
		
		// Update original values after successful save
		originalPresetId = selectedPresetId;
		originalPatch = { ...customPatch };
		dirty = false;
		return true;
	} catch (e) {
		console.error('Failed to save appearance:', e);
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
