// Draft payload types (matches api_contract.md section 6.1)

export interface DraftPayload {
  page: DraftPage;
  blocks: DraftBlock[];
  link_groups: Record<string, DraftGroup>;
  links: Record<string, DraftLink[]>;
}

export interface DraftPage {
  id: number;
  title: string;
  locale: string;
  status: string;
  theme: DraftTheme;
  access_type: string;
  settings: DraftSettings;
}

export interface DraftTheme {
  preset_key: string;
  custom_id: number | null;
  mode: string;
}

export interface DraftSettings {
  header: DraftHeader;
  cover?: DraftCover;
}

export interface DraftHeader {
  avatar_asset_id: number | null;
  name: string;
  bio: string;
  verified: boolean;
  social: SocialLink[];
}

export interface SocialLink {
  type: string;
  url: string;
}

export interface DraftCover {
  kind: string;
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
  layout_type: string;
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
