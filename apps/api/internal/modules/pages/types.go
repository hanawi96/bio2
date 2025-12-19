package pages

import "time"

// CreatePageRequest for POST /api/pages
type CreatePageRequest struct {
	Title          string `json:"title"`
	Locale         string `json:"locale"`
	ThemePresetKey string `json:"theme_preset_key"`
}

// PageListItem for GET /api/pages response
type PageListItem struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Locale    string    `json:"locale"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DraftPayload matches api_contract.md section 6.1
type DraftPayload struct {
	Page       DraftPage              `json:"page"`
	Blocks     []DraftBlock           `json:"blocks"`
	LinkGroups map[string]DraftGroup  `json:"link_groups"`
	Links      map[string][]DraftLink `json:"links"`
}

type DraftPage struct {
	ID         int64         `json:"id"`
	Title      string        `json:"title"`
	Locale     string        `json:"locale"`
	Status     string        `json:"status"`
	Theme      DraftTheme    `json:"theme"`
	AccessType string        `json:"access_type"`
	Settings   DraftSettings `json:"settings"`
}

type DraftTheme struct {
	PresetKey string `json:"preset_key"`
	CustomID  *int64 `json:"custom_id"`
	Mode      string `json:"mode"`
}

type DraftSettings struct {
	Header *DraftHeader `json:"header,omitempty"`
	Cover  *DraftCover  `json:"cover,omitempty"`
}

type DraftHeader struct {
	AvatarAssetID *int64        `json:"avatar_asset_id"`
	Name          string        `json:"name"`
	Bio           string        `json:"bio"`
	Verified      bool          `json:"verified"`
	Social        []SocialLink  `json:"social"`
}

type SocialLink struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type DraftCover struct {
	Kind    string  `json:"kind"`
	Color   *string `json:"color"`
	AssetID *int64  `json:"asset_id"`
}

type DraftBlock struct {
	ID        int64                  `json:"id"`
	Type      string                 `json:"type"`
	SortKey   string                 `json:"sort_key"`
	RefID     *int64                 `json:"ref_id"`
	Content   map[string]interface{} `json:"content"`
	IsVisible bool                   `json:"is_visible"`
}

type DraftGroup struct {
	ID            int64                  `json:"id"`
	Title         string                 `json:"title"`
	LayoutType    string                 `json:"layout_type"`
	LayoutConfig  map[string]interface{} `json:"layout_config"`
	StyleOverride map[string]interface{} `json:"style_override"`
}

type DraftLink struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	IconAssetID *int64 `json:"icon_asset_id"`
	SortKey     string `json:"sort_key"`
	IsActive    bool   `json:"is_active"`
}

// Compiled payload types (for publish)
type CompiledPayload struct {
	Version int                    `json:"version"`
	Page    CompiledPage           `json:"page"`
	Theme   CompiledTheme          `json:"theme"`
	Blocks  []map[string]interface{} `json:"blocks"`
}

type CompiledPage struct {
	ID       int64                  `json:"id"`
	Locale   string                 `json:"locale"`
	Settings map[string]interface{} `json:"settings"`
}

type CompiledTheme struct {
	PresetKey string                 `json:"preset_key"`
	Mode      string                 `json:"mode"`
	Compiled  map[string]interface{} `json:"compiled"`
}
