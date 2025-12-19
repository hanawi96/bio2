#!/bin/bash
# Seed database with initial data

set -e

DB_URL="${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/linkbio?sslmode=disable}"

echo "Seeding database..."

psql "$DB_URL" << 'EOF'
-- Insert plans
INSERT INTO plans (code, name) VALUES 
  ('FREE', 'Free Plan'),
  ('PRO', 'Pro Plan')
ON CONFLICT (code) DO NOTHING;

-- Insert default theme preset
INSERT INTO theme_presets (key, name, tier, visibility, config) VALUES 
  ('theme_a', 'Theme A (Free)', 'free', 'public', '{
    "meta": {
      "schemaVersion": 1,
      "id": "theme_a",
      "name": "Theme A (Free)"
    },
    "page": {
      "mode": "light",
      "layout": {
        "textAlign": "center",
        "baseFontSize": "M",
        "pagePadding": 16,
        "blockGap": 12
      },
      "defaults": {
        "linkGroup": {
          "textAlign": "center",
          "fontSize": "M",
          "radius": 16
        }
      }
    },
    "background": {
      "kind": "color",
      "color": "#0B0F19"
    }
  }')
ON CONFLICT (key) DO NOTHING;

-- Insert system domain for testing
INSERT INTO domains (hostname, status, is_system) VALUES 
  ('localhost:5173', 'active', true)
ON CONFLICT (hostname) DO NOTHING;

-- Insert test user
INSERT INTO users (email, password_hash, display_name) VALUES 
  ('test@example.com', '$2a$10$dummy', 'Test User')
ON CONFLICT (email) DO NOTHING;

-- Get IDs for foreign keys
DO $$
DECLARE
  v_user_id BIGINT;
  v_theme_id BIGINT;
  v_domain_id BIGINT;
  v_page_id BIGINT;
  v_group_id BIGINT;
BEGIN
  SELECT id INTO v_user_id FROM users WHERE email = 'test@example.com';
  SELECT id INTO v_theme_id FROM theme_presets WHERE key = 'theme_a';
  SELECT id INTO v_domain_id FROM domains WHERE hostname = 'localhost:5173';

  -- Insert test page
  INSERT INTO bio_pages (user_id, locale, title, status, theme_preset_id, theme_mode, settings)
  VALUES (v_user_id, 'vi', 'My Bio', 'published', v_theme_id, 'light', '{
    "header": {
      "name": "Yen",
      "bio": "Xin chào 👋",
      "verified": true,
      "social": [
        {"type": "facebook", "url": "https://facebook.com/yendev96"}
      ]
    },
    "cover": {
      "kind": "color",
      "color": "#0B0F19"
    }
  }')
  RETURNING id INTO v_page_id;

  -- Insert page route
  INSERT INTO page_routes (page_id, domain_id, path, is_current)
  VALUES (v_page_id, v_domain_id, '/', true);

  -- Insert link group
  INSERT INTO link_groups (page_id, title, layout_type, layout_config, style_override)
  VALUES (v_page_id, 'Follow me', 'cards', '{"gap": 12, "columns": 1}', '{"textAlign": "left", "fontSize": "L", "radius": 20}')
  RETURNING id INTO v_group_id;

  -- Insert links
  INSERT INTO links (group_id, title, url, sort_key) VALUES 
    (v_group_id, 'Instagram', 'https://instagram.com/yendev96', 'a0'),
    (v_group_id, 'TikTok', 'https://tiktok.com/@yendev96', 'a1');

  -- Insert blocks
  INSERT INTO blocks (page_id, type, sort_key, ref_id, content, is_visible) VALUES 
    (v_page_id, 'link_group', 'a0', v_group_id, '{}', true),
    (v_page_id, 'text', 'a1', NULL, '{"text": "My products", "variant": "heading"}', true);

  -- Insert publish cache
  INSERT INTO page_publish_cache (page_id, compiled_json) VALUES 
    (v_page_id, '{
      "version": 1,
      "page": {
        "id": ' || v_page_id || ',
        "locale": "vi",
        "settings": {
          "header": {
            "name": "Yen",
            "bio": "Xin chào 👋",
            "verified": true,
            "social": [{"type": "facebook", "url": "https://facebook.com/yendev96"}]
          },
          "cover": {"kind": "color", "color": "#0B0F19"}
        }
      },
      "theme": {
        "preset_key": "theme_a",
        "mode": "light",
        "compiled": {
          "page": {
            "layout": {"textAlign": "center", "baseFontSize": "M", "pagePadding": 16, "blockGap": 12},
            "defaults": {"linkGroup": {"textAlign": "center", "fontSize": "M", "radius": 16}}
          },
          "background": {"kind": "color", "color": "#0B0F19"}
        }
      },
      "blocks": [
        {
          "type": "link_group",
          "group": {
            "id": ' || v_group_id || ',
            "title": "Follow me",
            "layout_type": "cards",
            "layout_config": {"gap": 12, "columns": 1},
            "final_style": {"textAlign": "left", "fontSize": "L", "radius": 20},
            "links": [
              {"title": "Instagram", "url": "https://instagram.com/yendev96", "icon_url": null},
              {"title": "TikTok", "url": "https://tiktok.com/@yendev96", "icon_url": null}
            ]
          }
        },
        {
          "type": "text",
          "content": {"text": "My products", "variant": "heading"}
        }
      ]
    }');
END $$;

EOF

echo "Seed completed!"
