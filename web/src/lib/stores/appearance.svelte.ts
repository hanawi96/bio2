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

// Default appearance values
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
		effects: { blur: 0, dim: 0 }
	},
	header: {
		style: 'default' as 'default' | 'minimal' | 'centered' | 'card',
		showAvatar: true,
		showBio: true,
		avatarSize: 'M' as 'S' | 'M' | 'L',
		avatarShape: 'circle' as 'circle' | 'rounded' | 'square'
	},
	links: {
		style: 'filled' as 'filled' | 'outline' | 'soft' | 'ghost',
		borderRadius: 12,
		shadow: 'sm' as 'none' | 'sm' | 'md' | 'lg',
		padding: 16,
		gap: 12,
		textAlign: 'center' as 'left' | 'center' | 'right',
		fontSize: 'M' as 'S' | 'M' | 'L'
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
			return deepMerge(defaultAppearance, customPatch) as AppearanceSettings;
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

// Only update local state, no API call
export function selectPreset(presetId: number) {
	selectedPresetId = presetId;
	checkDirty();
}

// Only update local state, no API call
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

// Reset to original values
export function resetAppearance() {
	selectedPresetId = originalPresetId;
	customPatch = { ...originalPatch };
	dirty = false;
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
