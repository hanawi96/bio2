-- Link-in-Bio Builder — Initial Migration (Postgres) (Updated)
-- One-shot schema creation. Run once on an empty database.
-- Updated decisions:
-- - Plans: Free/Pro. Custom domain + password page are Pro-only (enforced in app).
-- - Custom domain: one domain maps to a single page at root '/' (enforced in app/infra).
-- - Theme marketplace: preset themes can be free/pro and have visibility/author metadata.
-- - Modes enabled: light/dark/compact (stored on page).
-- - GDPR minimal: deleting a user cascades delete pages and user-upload assets.

BEGIN;

-- =========================
-- 0) USERS & AUTH
-- =========================

CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,

  email TEXT UNIQUE,
  password_hash TEXT, -- bcrypt/argon2id (nullable nếu OAuth-only)

  display_name TEXT,

  is_active BOOLEAN NOT NULL DEFAULT TRUE,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_user_auth_identifier CHECK (
    (email IS NOT NULL) OR (password_hash IS NOT NULL)
  )
);

CREATE TABLE oauth_accounts (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  provider TEXT NOT NULL,           -- 'google' | 'facebook'
  provider_user_id TEXT NOT NULL,

  access_token TEXT NULL,
  refresh_token TEXT NULL,
  token_expires_at TIMESTAMPTZ NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT uq_oauth_provider_user UNIQUE (provider, provider_user_id),
  CONSTRAINT uq_user_provider UNIQUE (user_id, provider)
);

CREATE INDEX idx_oauth_user ON oauth_accounts(user_id);

-- =========================
-- 0.5) PLANS & SUBSCRIPTIONS (Free/Pro)
-- =========================

CREATE TABLE plans (
  id BIGSERIAL PRIMARY KEY,
  code TEXT NOT NULL UNIQUE, -- 'FREE' | 'PRO'
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE subscriptions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  plan_id BIGINT NOT NULL REFERENCES plans(id),

  status TEXT NOT NULL, -- active|canceled|past_due
  current_period_end TIMESTAMPTZ NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_subscription_status CHECK (status IN ('active','canceled','past_due'))
);

CREATE INDEX idx_subscriptions_user ON subscriptions(user_id);

-- =========================
-- 1) THEME SYSTEM + MARKETPLACE
-- =========================

CREATE TABLE theme_presets (
  id BIGSERIAL PRIMARY KEY,
  key TEXT NOT NULL UNIQUE,     -- 'theme_a', 'theme_b', ...
  name TEXT NOT NULL,

  tier TEXT NOT NULL DEFAULT 'free',          -- free|pro
  visibility TEXT NOT NULL DEFAULT 'public',  -- public|unlisted|private
  is_official BOOLEAN NOT NULL DEFAULT TRUE,
  author_user_id BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,

  -- Full theme JSON: meta/tokens/semantic/recipes/page/background/modes
  config JSONB NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_theme_tier CHECK (tier IN ('free','pro')),
  CONSTRAINT chk_theme_visibility CHECK (visibility IN ('public','unlisted','private'))
);

CREATE INDEX idx_theme_presets_tier ON theme_presets(tier);

CREATE TABLE themes_custom (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  based_on_preset_id BIGINT NOT NULL REFERENCES theme_presets(id),
  name TEXT NULL,

  patch JSONB NOT NULL,
  compiled_config JSONB NULL,

  hash TEXT NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT uq_custom_theme_user_hash UNIQUE (user_id, hash)
);

CREATE INDEX idx_custom_themes_user ON themes_custom(user_id);

-- =========================
-- 2) ASSETS (GDPR: cascade delete user assets)
-- =========================

CREATE TABLE assets (
  id BIGSERIAL PRIMARY KEY,

  user_id BIGINT NULL REFERENCES users(id) ON DELETE CASCADE,
  scope TEXT NOT NULL,   -- system_preset | user_upload
  type TEXT NOT NULL,    -- image (future: video)

  provider TEXT NOT NULL,
  storage_key TEXT NOT NULL,
  url TEXT NULL,

  mime_type TEXT NULL,
  size_bytes BIGINT NULL,
  width INT NULL,
  height INT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_asset_scope CHECK (scope IN ('system_preset','user_upload'))
);

CREATE INDEX idx_assets_user ON assets(user_id, created_at);

-- Optional avatar on user
ALTER TABLE users
  ADD COLUMN avatar_asset_id BIGINT NULL;

ALTER TABLE users
  ADD CONSTRAINT fk_users_avatar_asset
  FOREIGN KEY (avatar_asset_id) REFERENCES assets(id) ON DELETE SET NULL;

-- =========================
-- 3) PAGES (draft/publish + locale + password + settings + mode)
-- =========================

CREATE TABLE bio_pages (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  locale TEXT NOT NULL DEFAULT 'vi',
  title TEXT NULL,

  status TEXT NOT NULL DEFAULT 'draft', -- draft|published

  access_type TEXT NOT NULL DEFAULT 'public', -- public|password
  password_hash TEXT NULL,
  password_updated_at TIMESTAMPTZ NULL,

  theme_preset_id BIGINT NOT NULL REFERENCES theme_presets(id),
  theme_custom_id BIGINT NULL REFERENCES themes_custom(id),

  theme_mode TEXT NOT NULL DEFAULT 'light', -- light|dark|compact

  settings JSONB NOT NULL DEFAULT '{}'::JSONB, -- header, cover, social, etc

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_locale CHECK (locale IN ('vi','en')),
  CONSTRAINT chk_access_type CHECK (access_type IN ('public','password')),
  CONSTRAINT chk_theme_mode CHECK (theme_mode IN ('light','dark','compact'))
);

CREATE INDEX idx_pages_user ON bio_pages(user_id, created_at);
CREATE INDEX idx_pages_status ON bio_pages(status, updated_at);

-- =========================
-- 4) DOMAINS + ROUTES (system domain supports many pages; custom domain single page at '/')
-- =========================

CREATE TABLE domains (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,

  hostname TEXT NOT NULL UNIQUE,
  status TEXT NOT NULL DEFAULT 'pending', -- pending|active|disabled
  is_system BOOLEAN NOT NULL DEFAULT FALSE,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_domain_status CHECK (status IN ('pending','active','disabled'))
);

CREATE INDEX idx_domains_user ON domains(user_id);

CREATE TABLE page_routes (
  id BIGSERIAL PRIMARY KEY,

  page_id BIGINT NOT NULL REFERENCES bio_pages(id) ON DELETE CASCADE,
  domain_id BIGINT NOT NULL REFERENCES domains(id) ON DELETE CASCADE,

  path TEXT NOT NULL,                   -- system domain: '/yendev96'; custom domain: '/' (enforced in app)
  is_current BOOLEAN NOT NULL DEFAULT TRUE,

  redirect_to_route_id BIGINT NULL REFERENCES page_routes(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_domain_current_path
  ON page_routes(domain_id, path)
  WHERE is_current = TRUE;

CREATE INDEX idx_route_lookup ON page_routes(domain_id, path);
CREATE INDEX idx_route_current_by_page ON page_routes(page_id, is_current);

-- =========================
-- 5) LINK GROUPS + LINKS (style by group)
-- =========================

CREATE TABLE link_groups (
  id BIGSERIAL PRIMARY KEY,
  page_id BIGINT NOT NULL REFERENCES bio_pages(id) ON DELETE CASCADE,

  title TEXT NULL,

  layout_type TEXT NOT NULL DEFAULT 'list',
  layout_config JSONB NOT NULL DEFAULT '{}'::JSONB,

  style_override JSONB NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_layout_type CHECK (layout_type IN ('list','cards','grid'))
);

CREATE INDEX idx_groups_page ON link_groups(page_id);

CREATE TABLE links (
  id BIGSERIAL PRIMARY KEY,
  group_id BIGINT NOT NULL REFERENCES link_groups(id) ON DELETE CASCADE,

  title TEXT NOT NULL,
  url TEXT NOT NULL,

  icon_asset_id BIGINT NULL REFERENCES assets(id) ON DELETE SET NULL,

  sort_key TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_links_group_sort ON links(group_id, sort_key);
CREATE INDEX idx_links_group_active ON links(group_id, is_active);

-- =========================
-- 6) BLOCKS (1-column layout, no block style override)
-- =========================

CREATE TABLE blocks (
  id BIGSERIAL PRIMARY KEY,
  page_id BIGINT NOT NULL REFERENCES bio_pages(id) ON DELETE CASCADE,

  type TEXT NOT NULL,     -- link_group|text|image|product|spacer|embed|social_row|form|...
  sort_key TEXT NOT NULL,

  ref_id BIGINT NULL REFERENCES link_groups(id) ON DELETE SET NULL,
  content JSONB NOT NULL DEFAULT '{}'::JSONB,

  is_visible BOOLEAN NOT NULL DEFAULT TRUE,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_block_ref_link_group CHECK (
    (type = 'link_group' AND ref_id IS NOT NULL) OR
    (type <> 'link_group' AND ref_id IS NULL)
  )
);

CREATE INDEX idx_blocks_page_sort ON blocks(page_id, sort_key);
CREATE INDEX idx_blocks_page_type ON blocks(page_id, type);

-- =========================
-- 7) PUBLISH CACHE (DB + CDN)
-- =========================

CREATE TABLE page_publish_cache (
  page_id BIGINT PRIMARY KEY REFERENCES bio_pages(id) ON DELETE CASCADE,
  compiled_json JSONB NOT NULL,
  published_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =========================
-- 8) PASSWORD SESSIONS (remember 7 days)
-- =========================

CREATE TABLE page_access_sessions (
  id BIGSERIAL PRIMARY KEY,
  page_id BIGINT NOT NULL REFERENCES bio_pages(id) ON DELETE CASCADE,

  token_hash TEXT NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT uq_page_token UNIQUE (page_id, token_hash)
);

CREATE INDEX idx_page_access_exp ON page_access_sessions(page_id, expires_at);

COMMIT;

-- =========================
-- REVIEW NOTES
-- =========================
-- - Pro-only gates (custom domain, password pages, pro themes) should be enforced in application logic using subscriptions.
-- - Custom domain single-page-at-root is enforced in app/infra by only issuing path='/' for non-system domains.
-- - User deletion cascades pages, routes, groups, links, blocks, publish cache, sessions, and user-owned assets (GDPR minimum).
