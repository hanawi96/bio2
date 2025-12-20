-- Seed Theme Presets and Additional Migrations
-- Run this after the initial migration

-- Add username column to users if not exists
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'username') THEN
        ALTER TABLE users ADD COLUMN username TEXT UNIQUE;
    END IF;
END $$;

INSERT INTO theme_presets (key, name, tier, visibility, is_official, config) VALUES
('light', 'Light', 'free', 'public', true, '{
  "meta": {
    "id": "preset.light",
    "name": "Light",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Clean light theme"
  },
  "tokens": {
    "color": {
      "gray": {"50": "#f9fafb", "100": "#f3f4f6", "200": "#e5e7eb", "300": "#d1d5db", "400": "#9ca3af", "500": "#6b7280", "600": "#4b5563", "700": "#374151", "800": "#1f2937", "900": "#111827"}
    }
  },
  "page": {
    "layout": {"maxWidth": 520, "pagePadding": 16, "blockGap": 12, "textAlign": "center", "baseFontSize": "M"}
  },
  "background": {
    "type": "solid",
    "color": "#f2f2f7"
  }
}'::jsonb),

('dark', 'Dark', 'free', 'public', true, '{
  "meta": {
    "id": "preset.dark",
    "name": "Dark",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Elegant dark theme"
  },
  "tokens": {
    "color": {
      "gray": {"50": "#111827", "100": "#1f2937", "200": "#374151", "300": "#4b5563", "400": "#6b7280", "500": "#9ca3af", "600": "#d1d5db", "700": "#e5e7eb", "800": "#f3f4f6", "900": "#f9fafb"}
    }
  },
  "page": {
    "layout": {"maxWidth": 520, "pagePadding": 16, "blockGap": 12, "textAlign": "center", "baseFontSize": "M"}
  },
  "background": {
    "type": "solid",
    "color": "#1c1c1e"
  }
}'::jsonb),

('ocean', 'Ocean', 'free', 'public', true, '{
  "meta": {
    "id": "preset.ocean",
    "name": "Ocean",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Beautiful ocean gradient"
  },
  "page": {
    "layout": {"maxWidth": 520, "pagePadding": 16, "blockGap": 12, "textAlign": "center", "baseFontSize": "M"}
  },
  "background": {
    "type": "gradient",
    "gradient": "linear-gradient(135deg, #667eea 0%, #764ba2 100%)"
  }
}'::jsonb),

('sunset', 'Sunset', 'pro', 'public', true, '{
  "meta": {
    "id": "preset.sunset",
    "name": "Sunset",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Warm sunset gradient"
  },
  "page": {
    "layout": {"maxWidth": 520, "pagePadding": 16, "blockGap": 12, "textAlign": "center", "baseFontSize": "M"}
  },
  "background": {
    "type": "gradient",
    "gradient": "linear-gradient(135deg, #fa709a 0%, #fee140 100%)"
  }
}'::jsonb),

('forest', 'Forest', 'pro', 'public', true, '{
  "meta": {
    "id": "preset.forest",
    "name": "Forest",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Fresh forest gradient"
  },
  "page": {
    "layout": {"maxWidth": 520, "pagePadding": 16, "blockGap": 12, "textAlign": "center", "baseFontSize": "M"}
  },
  "background": {
    "type": "gradient",
    "gradient": "linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)"
  }
}'::jsonb)

ON CONFLICT (key) DO UPDATE SET
  name = EXCLUDED.name,
  tier = EXCLUDED.tier,
  config = EXCLUDED.config,
  updated_at = NOW();

-- Also seed plans if not exists
INSERT INTO plans (code, name) VALUES
('FREE', 'Free Plan'),
('PRO', 'Pro Plan')
ON CONFLICT (code) DO NOTHING;

-- Seed system domain if not exists
INSERT INTO domains (hostname, status, is_system) VALUES
('localhost', 'active', true)
ON CONFLICT (hostname) DO NOTHING;
