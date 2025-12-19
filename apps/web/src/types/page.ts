// Compiled payload types (matches api_contract.md section 6.2)
// Re-export from shared contracts in production

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
