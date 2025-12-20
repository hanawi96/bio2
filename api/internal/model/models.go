package model

import (
	"encoding/json"
	"time"
)

// User
type User struct {
	ID            int64     `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	Username      *string   `json:"username"`
	DisplayName   *string   `json:"display_name"`
	AvatarAssetID *int64    `json:"avatar_asset_id"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Plan
type Plan struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Subscription
type Subscription struct {
	ID               int64      `json:"id"`
	UserID           int64      `json:"user_id"`
	PlanID           int64      `json:"plan_id"`
	Status           string     `json:"status"`
	CurrentPeriodEnd *time.Time `json:"current_period_end"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ThemePreset
type ThemePreset struct {
	ID           int64           `json:"id"`
	Key          string          `json:"key"`
	Name         string          `json:"name"`
	Tier         string          `json:"tier"`
	Visibility   string          `json:"visibility"`
	IsOfficial   bool            `json:"is_official"`
	AuthorUserID *int64          `json:"author_user_id"`
	Config       json.RawMessage `json:"config"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// ThemeCustom
type ThemeCustom struct {
	ID              int64           `json:"id"`
	UserID          int64           `json:"user_id"`
	BasedOnPresetID int64           `json:"based_on_preset_id"`
	Name            *string         `json:"name"`
	Patch           json.RawMessage `json:"patch"`
	CompiledConfig  json.RawMessage `json:"compiled_config"`
	Hash            string          `json:"hash"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// Asset
type Asset struct {
	ID         int64     `json:"id"`
	UserID     *int64    `json:"user_id"`
	Scope      string    `json:"scope"`
	Type       string    `json:"type"`
	Provider   string    `json:"provider"`
	StorageKey string    `json:"storage_key"`
	URL        *string   `json:"url"`
	MimeType   *string   `json:"mime_type"`
	SizeBytes  *int64    `json:"size_bytes"`
	Width      *int      `json:"width"`
	Height     *int      `json:"height"`
	CreatedAt  time.Time `json:"created_at"`
}

// BioPage
type BioPage struct {
	ID              int64           `json:"id"`
	UserID          int64           `json:"user_id"`
	Locale          string          `json:"locale"`
	Title           *string         `json:"title"`
	Status          string          `json:"status"`
	AccessType      string          `json:"access_type"`
	PasswordHash    *string         `json:"-"`
	ThemePresetID   int64           `json:"theme_preset_id"`
	ThemeCustomID   *int64          `json:"theme_custom_id"`
	ThemeMode       string          `json:"theme_mode"`
	Settings        json.RawMessage `json:"settings"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// Domain
type Domain struct {
	ID        int64     `json:"id"`
	UserID    *int64    `json:"user_id"`
	Hostname  string    `json:"hostname"`
	Status    string    `json:"status"`
	IsSystem  bool      `json:"is_system"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PageRoute
type PageRoute struct {
	ID                int64     `json:"id"`
	PageID            int64     `json:"page_id"`
	DomainID          int64     `json:"domain_id"`
	Path              string    `json:"path"`
	IsCurrent         bool      `json:"is_current"`
	RedirectToRouteID *int64    `json:"redirect_to_route_id"`
	CreatedAt         time.Time `json:"created_at"`
}

// LinkGroup
type LinkGroup struct {
	ID            int64           `json:"id"`
	PageID        int64           `json:"page_id"`
	Title         *string         `json:"title"`
	LayoutType    string          `json:"layout_type"`
	LayoutConfig  json.RawMessage `json:"layout_config"`
	StyleOverride json.RawMessage `json:"style_override"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// Link
type Link struct {
	ID          int64     `json:"id"`
	GroupID     int64     `json:"group_id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	IconAssetID *int64    `json:"icon_asset_id"`
	SortKey     string    `json:"sort_key"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Block
type Block struct {
	ID        int64           `json:"id"`
	PageID    int64           `json:"page_id"`
	Type      string          `json:"type"`
	SortKey   string          `json:"sort_key"`
	RefID     *int64          `json:"ref_id"`
	Content   json.RawMessage `json:"content"`
	IsVisible bool            `json:"is_visible"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// PagePublishCache
type PagePublishCache struct {
	PageID       int64           `json:"page_id"`
	CompiledJSON json.RawMessage `json:"compiled_json"`
	PublishedAt  time.Time       `json:"published_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// PageAccessSession
type PageAccessSession struct {
	ID        int64     `json:"id"`
	PageID    int64     `json:"page_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
