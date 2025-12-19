// ============================================
// Draft payload types (editor) - api_contract.md section 6.1
// ============================================

export interface DraftPayload {
  page: DraftPage;
  blocks: DraftBlock[];
  link_groups: Record<string, DraftGroup>;
  links: Record<string, DraftLink[]>;
}

export interface DraftPage {
  id: number;
  title: string;
  locale: 'vi' | 'en';
  status: 'draft' | 'published';
  theme: DraftTheme;
  access_type: 'public' | 'password';
  settings: DraftSettings;
}

export interface DraftTheme {
  preset_key: string;
  custom_id: number | null;
  mode: 'light' | 'dark' | 'compact';
}

export interface DraftSettings {
  header?: DraftHeader;
  cover?: DraftCover;
}

export interface DraftHeader {
  avatar_asset_id: number | null;
  name: string;
  bio: string;
  verified: boolean;
  social: SocialLink[];
}

export interface DraftCover {
  kind: 'color' | 'image';
  color: string | null;
  asset_id: number | null;
}

export interface DraftBlock {
  id: number;
  type: string;
  sort_key: string;
  ref_id: number | null;
  content: Record<string, unknown>;
  is_visible: boolean;
}

export interface DraftGroup {
  id: number;
  title: string;
  layout_type: 'list' | 'cards' | 'grid';
  layout_config: Record<string, unknown>;
  style_override: Record<string, unknown>;
}

export interface DraftLink {
  id: number;
  title: string;
  url: string;
  icon_asset_id: number | null;
  sort_key: string;
  is_active: boolean;
}

// ============================================
// Compiled payload types (public render) - api_contract.md section 6.2
// ============================================

export interface CompiledPayload {
  version: number;
  page: CompiledPage;
  theme: CompiledTheme;
  blocks: CompiledBlock[];
}

export interface CompiledPage {
  id: number;
  locale: 'vi' | 'en';
  settings: PageSettings;
}

export interface PageSettings {
  header?: HeaderSettings;
  cover?: CoverSettings;
}

export interface HeaderSettings {
  avatar_asset_id?: number | null;
  name?: string;
  bio?: string;
  verified?: boolean;
  social?: SocialLink[];
}

export interface SocialLink {
  type: string;
  url: string;
}

export interface CoverSettings {
  kind: 'color' | 'image';
  color?: string | null;
  asset_id?: number | null;
}

export interface CompiledTheme {
  preset_key: string;
  mode: 'light' | 'dark' | 'compact';
  compiled: ThemeCompiled;
}

export interface ThemeCompiled {
  page?: {
    layout?: {
      textAlign?: string;
      baseFontSize?: string;
      pagePadding?: number;
      blockGap?: number;
    };
    defaults?: {
      linkGroup?: LinkGroupStyle;
    };
  };
  background?: {
    kind: 'color' | 'image';
    color?: string;
    asset_id?: number;
    effects?: {
      blur?: number;
      dim?: number;
      overlayColor?: string | null;
    };
  };
}

export interface LinkGroupStyle {
  textAlign?: string;
  fontSize?: string;
  radius?: number;
}

// Block types
export type CompiledBlock = LinkGroupBlock | TextBlock | ProductBlock;

export interface LinkGroupBlock {
  type: 'link_group';
  group: CompiledLinkGroup;
}

export interface CompiledLinkGroup {
  id: number;
  title?: string;
  layout_type: 'list' | 'cards' | 'grid';
  layout_config: Record<string, unknown>;
  final_style: LinkGroupStyle;
  links: CompiledLink[];
}

export interface CompiledLink {
  title: string;
  url: string;
  icon_url?: string | null;
}

export interface TextBlock {
  type: 'text';
  content: {
    text: string;
    variant: 'body' | 'heading' | 'caption';
  };
}

export interface ProductBlock {
  type: 'product';
  content: {
    title: string;
    price: number;
    currency: string;
    url: string;
    image_url?: string | null;
  };
}
