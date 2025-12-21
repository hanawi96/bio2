-- Seed Theme Presets - Full Featured Themes
-- Run this after the initial migration

-- Add username column to users if not exists
DO $ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'username') THEN
        ALTER TABLE users ADD COLUMN username TEXT UNIQUE;
    END IF;
END $;

INSERT INTO theme_presets (key, name, tier, visibility, is_official, config) VALUES
('light', 'Light', 'free', 'public', true, '{
  "meta": {
    "id": "preset.light",
    "name": "Light",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Clean minimal light theme"
  },
  "semantic": {
    "color": {
      "primary": "#007aff",
      "text": { "default": "#1c1c1e", "muted": "#8e8e93" },
      "surface": { "page": "#f2f2f7", "card": "#ffffff" },
      "border": { "default": "rgba(60, 60, 67, 0.1)" }
    }
  },
  "recipes": {
    "linkItem": {
      "base": {
        "background": "#ffffff",
        "borderRadius": "12px",
        "padding": "16px",
        "shadow": "0 1px 3px rgba(0,0,0,0.08)",
        "color": "#1c1c1e"
      }
    }
  },
  "page": {
    "layout": { "maxWidth": 520, "pagePadding": 16, "blockGap": 12, "textAlign": "center", "baseFontSize": "M" },
    "defaults": {
      "linkGroup": { "textAlign": "center", "fontSize": "M", "radius": "12px", "shadow": "sm" }
    }
  },
  "background": {
    "type": "solid",
    "color": "#f2f2f7",
    "effects": { "blur": 0, "dim": 0 }
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
  "semantic": {
    "color": {
      "primary": "#0a84ff",
      "text": { "default": "#ffffff", "muted": "#98989d" },
      "surface": { "page": "#1c1c1e", "card": "#2c2c2e" },
      "border": { "default": "rgba(255, 255, 255, 0.1)" }
    }
  },
  "recipes": {
    "linkItem": {
      "base": {
        "background": "#2c2c2e",
        "borderRadius": "12px",
        "padding": "16px",
        "shadow": "0 2px 8px rgba(0,0,0,0.3)",
        "color": "#ffffff"
      }
    }
  },
  "page": {
    "layout": { "maxWidth": 520, "pagePadding": 16, "blockGap": 12, "textAlign": "center", "baseFontSize": "M" },
    "defaults": {
      "linkGroup": { "textAlign": "center", "fontSize": "M", "radius": "12px", "shadow": "md" }
    }
  },
  "background": {
    "type": "solid",
    "color": "#1c1c1e",
    "effects": { "blur": 0, "dim": 0 }
  }
}'::jsonb),

('ocean', 'Ocean', 'free', 'public', true, '{
  "meta": {
    "id": "preset.ocean",
    "name": "Ocean",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Beautiful ocean gradient with glass cards"
  },
  "semantic": {
    "color": {
      "primary": "#667eea",
      "text": { "default": "#1f2937", "muted": "#f0f9ff" },
      "surface": { "page": "transparent", "card": "rgba(255, 255, 255, 0.85)" },
      "border": { "default": "rgba(255, 255, 255, 0.3)" }
    }
  },
  "recipes": {
    "linkItem": {
      "base": {
        "background": "rgba(255, 255, 255, 0.85)",
        "borderRadius": "16px",
        "padding": "16px",
        "shadow": "0 4px 12px rgba(0,0,0,0.15)",
        "color": "#1f2937",
        "backdropFilter": "blur(8px)"
      }
    }
  },
  "page": {
    "layout": { "maxWidth": 480, "pagePadding": 20, "blockGap": 14, "textAlign": "center", "baseFontSize": "M" },
    "defaults": {
      "linkGroup": { "textAlign": "center", "fontSize": "M", "radius": "16px", "shadow": "lg" }
    }
  },
  "background": {
    "type": "gradient",
    "gradient": "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
    "effects": { "blur": 0, "dim": 0 }
  }
}'::jsonb),

('sunset', 'Sunset', 'pro', 'public', true, '{
  "meta": {
    "id": "preset.sunset",
    "name": "Sunset",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Warm sunset gradient with soft cards"
  },
  "semantic": {
    "color": {
      "primary": "#ff6b6b",
      "text": { "default": "#2d3436", "muted": "#fff5f5" },
      "surface": { "page": "transparent", "card": "rgba(255, 255, 255, 0.9)" },
      "border": { "default": "rgba(255, 255, 255, 0.4)" }
    }
  },
  "recipes": {
    "linkItem": {
      "base": {
        "background": "rgba(255, 255, 255, 0.9)",
        "borderRadius": "20px",
        "padding": "18px",
        "shadow": "0 8px 24px rgba(250, 112, 154, 0.25)",
        "color": "#2d3436"
      }
    }
  },
  "page": {
    "layout": { "maxWidth": 480, "pagePadding": 24, "blockGap": 16, "textAlign": "center", "baseFontSize": "M" },
    "defaults": {
      "linkGroup": { "textAlign": "center", "fontSize": "M", "radius": "20px", "shadow": "lg" }
    }
  },
  "background": {
    "type": "gradient",
    "gradient": "linear-gradient(135deg, #fa709a 0%, #fee140 100%)",
    "effects": { "blur": 0, "dim": 0 }
  }
}'::jsonb),

('forest', 'Forest', 'pro', 'public', true, '{
  "meta": {
    "id": "preset.forest",
    "name": "Forest",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Fresh forest gradient with nature vibes"
  },
  "semantic": {
    "color": {
      "primary": "#00b894",
      "text": { "default": "#2d3436", "muted": "#f0fdf4" },
      "surface": { "page": "transparent", "card": "rgba(255, 255, 255, 0.88)" },
      "border": { "default": "rgba(255, 255, 255, 0.35)" }
    }
  },
  "recipes": {
    "linkItem": {
      "base": {
        "background": "rgba(255, 255, 255, 0.88)",
        "borderRadius": "14px",
        "padding": "16px",
        "shadow": "0 6px 20px rgba(67, 233, 123, 0.2)",
        "color": "#2d3436"
      }
    }
  },
  "page": {
    "layout": { "maxWidth": 480, "pagePadding": 20, "blockGap": 14, "textAlign": "center", "baseFontSize": "M" },
    "defaults": {
      "linkGroup": { "textAlign": "center", "fontSize": "M", "radius": "14px", "shadow": "md" }
    }
  },
  "background": {
    "type": "gradient",
    "gradient": "linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)",
    "effects": { "blur": 0, "dim": 0 }
  }
}'::jsonb),

('midnight', 'Midnight', 'pro', 'public', true, '{
  "meta": {
    "id": "preset.midnight",
    "name": "Midnight",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Deep midnight blue with neon accents"
  },
  "semantic": {
    "color": {
      "primary": "#818cf8",
      "text": { "default": "#f1f5f9", "muted": "#94a3b8" },
      "surface": { "page": "#0f172a", "card": "rgba(30, 41, 59, 0.8)" },
      "border": { "default": "rgba(129, 140, 248, 0.2)" }
    }
  },
  "recipes": {
    "linkItem": {
      "base": {
        "background": "rgba(30, 41, 59, 0.8)",
        "borderRadius": "12px",
        "padding": "16px",
        "shadow": "0 0 20px rgba(129, 140, 248, 0.15)",
        "color": "#f1f5f9",
        "border": "1px solid rgba(129, 140, 248, 0.2)"
      }
    }
  },
  "page": {
    "layout": { "maxWidth": 480, "pagePadding": 20, "blockGap": 14, "textAlign": "center", "baseFontSize": "M" },
    "defaults": {
      "linkGroup": { "textAlign": "center", "fontSize": "M", "radius": "12px", "shadow": "lg" }
    }
  },
  "background": {
    "type": "solid",
    "color": "#0f172a",
    "effects": { "blur": 0, "dim": 0 }
  }
}'::jsonb),

('candy', 'Candy', 'free', 'public', true, '{
  "meta": {
    "id": "preset.candy",
    "name": "Candy",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Sweet pastel candy colors"
  },
  "semantic": {
    "color": {
      "primary": "#ec4899",
      "text": { "default": "#831843", "muted": "#be185d" },
      "surface": { "page": "transparent", "card": "rgba(255, 255, 255, 0.75)" },
      "border": { "default": "rgba(236, 72, 153, 0.2)" }
    }
  },
  "recipes": {
    "linkItem": {
      "base": {
        "background": "rgba(255, 255, 255, 0.75)",
        "borderRadius": "24px",
        "padding": "18px",
        "shadow": "0 4px 16px rgba(236, 72, 153, 0.2)",
        "color": "#831843"
      }
    }
  },
  "page": {
    "layout": { "maxWidth": 460, "pagePadding": 24, "blockGap": 16, "textAlign": "center", "baseFontSize": "M" },
    "defaults": {
      "linkGroup": { "textAlign": "center", "fontSize": "M", "radius": "24px", "shadow": "md" }
    }
  },
  "background": {
    "type": "gradient",
    "gradient": "linear-gradient(135deg, #fce7f3 0%, #dbeafe 50%, #fef3c7 100%)",
    "effects": { "blur": 0, "dim": 0 }
  }
}'::jsonb)

ON CONFLICT (key) DO UPDATE SET
  name = EXCLUDED.name,
  tier = EXCLUDED.tier,
  config = EXCLUDED.config,
  updated_at = NOW();

-- Seed plans if not exists
INSERT INTO plans (code, name) VALUES
('FREE', 'Free Plan'),
('PRO', 'Pro Plan')
ON CONFLICT (code) DO NOTHING;

-- Seed system domain if not exists
INSERT INTO domains (hostname, status, is_system) VALUES
('localhost', 'active', true)
ON CONFLICT (hostname) DO NOTHING;
