-- Aurora Glass Theme - Modern Glassmorphism Theme
-- A beautiful, trendy theme with glass effects and aurora gradients

INSERT INTO theme_presets (key, name, tier, visibility, is_official, config) VALUES
('aurora-glass', 'Aurora Glass', 'free', 'public', true, '{
  "meta": {
    "id": "preset.aurora-glass",
    "name": "Aurora Glass",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Modern glassmorphism with aurora gradient background",
    "supports": {
      "modes": ["light", "dark"],
      "animation": true,
      "glass": true
    },
    "tier": "free",
    "visibility": "public",
    "preview": {
      "thumbnailUrl": null,
      "demoUrl": null
    },
    "contract": {
      "controls": [
        {
          "keyPath": "page.layout.textAlign",
          "type": "select",
          "options": ["left", "center", "right"],
          "label": "Text align (page)",
          "affects": ["page", "header"]
        },
        {
          "keyPath": "page.layout.baseFontSize",
          "type": "select",
          "options": ["S", "M", "L", "XL"],
          "label": "Font size (page)"
        },
        {
          "keyPath": "page.layout.pagePadding",
          "type": "slider",
          "min": 0,
          "max": 32,
          "step": 4,
          "label": "Page padding"
        },
        {
          "keyPath": "page.layout.blockGap",
          "type": "slider",
          "min": 8,
          "max": 32,
          "step": 4,
          "label": "Block spacing"
        },
        {
          "keyPath": "page.mode",
          "type": "select",
          "options": ["light", "dark"],
          "label": "Color mode"
        },
        {
          "keyPath": "page.defaults.linkGroup.textAlign",
          "type": "select",
          "options": ["left", "center", "right"],
          "label": "Link group text align"
        },
        {
          "keyPath": "page.defaults.linkGroup.fontSize",
          "type": "select",
          "options": ["S", "M", "L", "XL"],
          "label": "Link group font size"
        },
        {
          "keyPath": "background.effects.blur",
          "type": "slider",
          "min": 0,
          "max": 12,
          "step": 1,
          "label": "Background blur"
        },
        {
          "keyPath": "background.effects.dim",
          "type": "slider",
          "min": 0,
          "max": 0.8,
          "step": 0.05,
          "label": "Background dim"
        }
      ],
      "constraints": {
        "background.effects.dim": {"min": 0, "max": 0.8},
        "background.effects.blur": {"min": 0, "max": 12},
        "page.layout.pagePadding": {"min": 0, "max": 32},
        "page.layout.blockGap": {"min": 8, "max": 32}
      },
      "defaultsForMissing": {
        "page.layout.textAlign": "center",
        "page.layout.baseFontSize": "M"
      },
      "notes": [
        "Theme controls chỉ chỉnh ở mức theme/page defaults.",
        "Muốn link style riêng: tạo Link Group riêng và set group style_override."
      ]
    }
  },

  "tokens": {
    "color": {
      "white": "#ffffff",
      "black": "#000000",
      "transparent": "transparent",

      "gray": {
        "50": "#f9fafb",
        "100": "#f3f4f6",
        "200": "#e5e7eb",
        "300": "#d1d5db",
        "400": "#9ca3af",
        "500": "#6b7280",
        "600": "#4b5563",
        "700": "#374151",
        "800": "#1f2937",
        "900": "#111827",
        "950": "#030712"
      },

      "slate": {
        "50": "#f8fafc",
        "100": "#f1f5f9",
        "200": "#e2e8f0",
        "300": "#cbd5e1",
        "400": "#94a3b8",
        "500": "#64748b",
        "600": "#475569",
        "700": "#334155",
        "800": "#1e293b",
        "900": "#0f172a",
        "950": "#020617"
      },

      "violet": {
        "50": "#f5f3ff",
        "100": "#ede9fe",
        "200": "#ddd6fe",
        "300": "#c4b5fd",
        "400": "#a78bfa",
        "500": "#8b5cf6",
        "600": "#7c3aed",
        "700": "#6d28d9",
        "800": "#5b21b6",
        "900": "#4c1d95"
      },

      "fuchsia": {
        "50": "#fdf4ff",
        "100": "#fae8ff",
        "200": "#f5d0fe",
        "300": "#f0abfc",
        "400": "#e879f9",
        "500": "#d946ef",
        "600": "#c026d3",
        "700": "#a21caf",
        "800": "#86198f",
        "900": "#701a75"
      },

      "cyan": {
        "50": "#ecfeff",
        "100": "#cffafe",
        "200": "#a5f3fc",
        "300": "#67e8f9",
        "400": "#22d3ee",
        "500": "#06b6d4",
        "600": "#0891b2",
        "700": "#0e7490",
        "800": "#155e75",
        "900": "#164e63"
      },

      "emerald": {
        "400": "#34d399",
        "500": "#10b981"
      },

      "rose": {
        "400": "#fb7185",
        "500": "#f43f5e"
      },

      "amber": {
        "400": "#fbbf24",
        "500": "#f59e0b"
      }
    },

    "typography": {
      "fontFamily": {
        "sans": "Inter, -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, sans-serif",
        "display": "Inter, -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, sans-serif",
        "mono": "\"SF Mono\", \"Fira Code\", \"Fira Mono\", Menlo, monospace"
      },
      "fontSizeScale": {
        "xs": "0.75rem",
        "sm": "0.875rem",
        "base": "1rem",
        "lg": "1.125rem",
        "xl": "1.25rem",
        "2xl": "1.5rem",
        "3xl": "1.875rem",
        "4xl": "2.25rem"
      },
      "fontWeight": {
        "normal": 400,
        "medium": 500,
        "semibold": 600,
        "bold": 700
      },
      "lineHeight": {
        "tight": 1.25,
        "snug": 1.375,
        "normal": 1.5,
        "relaxed": 1.625
      },
      "letterSpacing": {
        "tight": "-0.025em",
        "normal": "0",
        "wide": "0.025em"
      }
    },

    "space": {
      "0": "0",
      "1": "0.25rem",
      "2": "0.5rem",
      "3": "0.75rem",
      "4": "1rem",
      "5": "1.25rem",
      "6": "1.5rem",
      "8": "2rem",
      "10": "2.5rem",
      "12": "3rem",
      "16": "4rem",
      "20": "5rem"
    },

    "size": {
      "avatar": {
        "sm": "2.5rem",
        "md": "4rem",
        "lg": "5rem",
        "xl": "6rem"
      },
      "icon": {
        "sm": "1rem",
        "md": "1.25rem",
        "lg": "1.5rem",
        "xl": "2rem"
      },
      "button": {
        "sm": "2rem",
        "md": "2.5rem",
        "lg": "3rem"
      }
    },

    "radius": {
      "none": "0",
      "sm": "0.25rem",
      "md": "0.5rem",
      "lg": "0.75rem",
      "xl": "1rem",
      "2xl": "1.5rem",
      "3xl": "2rem",
      "full": "9999px"
    },

    "elevation": {
      "none": "none",
      "xs": "0 1px 2px 0 rgba(0, 0, 0, 0.05)",
      "sm": "0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1)",
      "md": "0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1)",
      "lg": "0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)",
      "xl": "0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)",
      "2xl": "0 25px 50px -12px rgba(0, 0, 0, 0.25)",
      "inner": "inset 0 2px 4px 0 rgba(0, 0, 0, 0.05)",
      "glow": {
        "violet": "0 0 20px rgba(139, 92, 246, 0.35)",
        "fuchsia": "0 0 20px rgba(217, 70, 239, 0.35)",
        "cyan": "0 0 20px rgba(6, 182, 212, 0.35)"
      }
    },

    "motion": {
      "duration": {
        "fast": "150ms",
        "normal": "200ms",
        "slow": "300ms",
        "slower": "500ms"
      },
      "easing": {
        "default": "cubic-bezier(0.4, 0, 0.2, 1)",
        "in": "cubic-bezier(0.4, 0, 1, 1)",
        "out": "cubic-bezier(0, 0, 0.2, 1)",
        "inOut": "cubic-bezier(0.4, 0, 0.2, 1)",
        "bounce": "cubic-bezier(0.68, -0.55, 0.265, 1.55)"
      }
    },

    "zIndex": {
      "base": 0,
      "dropdown": 1000,
      "sticky": 1100,
      "modal": 1200,
      "popover": 1300,
      "tooltip": 1400
    },

    "breakpoint": {
      "sm": "640px",
      "md": "768px",
      "lg": "1024px",
      "xl": "1280px"
    },

    "glass": {
      "blur": {
        "sm": "8px",
        "md": "12px",
        "lg": "16px",
        "xl": "24px"
      },
      "opacity": {
        "light": 0.7,
        "medium": 0.5,
        "heavy": 0.3
      }
    }
  },

  "semantic": {
    "color": {
      "primary": "#8b5cf6",
      "primaryHover": "#7c3aed",
      "primaryMuted": "rgba(139, 92, 246, 0.15)",

      "secondary": "#06b6d4",
      "secondaryHover": "#0891b2",

      "accent": "#d946ef",
      "accentHover": "#c026d3",

      "text": {
        "default": "#1f2937",
        "muted": "#6b7280",
        "subtle": "#9ca3af",
        "invert": "#ffffff",
        "link": "#8b5cf6",
        "linkHover": "#7c3aed"
      },

      "border": {
        "default": "rgba(255, 255, 255, 0.2)",
        "subtle": "rgba(255, 255, 255, 0.1)",
        "strong": "rgba(255, 255, 255, 0.3)",
        "focus": "#8b5cf6"
      },

      "divider": "rgba(255, 255, 255, 0.15)",

      "surface": {
        "page": "transparent",
        "card": "rgba(255, 255, 255, 0.7)",
        "cardHover": "rgba(255, 255, 255, 0.85)",
        "overlay": "rgba(0, 0, 0, 0.5)",
        "glass": "rgba(255, 255, 255, 0.25)"
      },

      "status": {
        "success": "#10b981",
        "successBg": "rgba(16, 185, 129, 0.15)",
        "warning": "#f59e0b",
        "warningBg": "rgba(245, 158, 11, 0.15)",
        "danger": "#f43f5e",
        "dangerBg": "rgba(244, 63, 94, 0.15)",
        "info": "#06b6d4",
        "infoBg": "rgba(6, 182, 212, 0.15)"
      }
    },

    "typography": {
      "heading": {
        "fontFamily": "Inter, -apple-system, BlinkMacSystemFont, sans-serif",
        "fontWeight": 700,
        "lineHeight": 1.25,
        "letterSpacing": "-0.025em"
      },
      "body": {
        "fontFamily": "Inter, -apple-system, BlinkMacSystemFont, sans-serif",
        "fontWeight": 400,
        "lineHeight": 1.5,
        "letterSpacing": "0"
      },
      "caption": {
        "fontFamily": "Inter, -apple-system, BlinkMacSystemFont, sans-serif",
        "fontWeight": 400,
        "lineHeight": 1.5,
        "letterSpacing": "0"
      },
      "button": {
        "fontFamily": "Inter, -apple-system, BlinkMacSystemFont, sans-serif",
        "fontWeight": 600,
        "lineHeight": 1.25,
        "letterSpacing": "0"
      }
    },

    "surface": {
      "glass": {
        "background": "rgba(255, 255, 255, 0.25)",
        "backdropFilter": "blur(12px)",
        "border": "1px solid rgba(255, 255, 255, 0.2)",
        "shadow": "0 4px 6px -1px rgba(0, 0, 0, 0.1)"
      },
      "card": {
        "background": "rgba(255, 255, 255, 0.7)",
        "backdropFilter": "blur(12px)",
        "border": "1px solid rgba(255, 255, 255, 0.3)",
        "shadow": "0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1)"
      },
      "elevated": {
        "background": "rgba(255, 255, 255, 0.9)",
        "backdropFilter": "blur(16px)",
        "border": "1px solid rgba(255, 255, 255, 0.4)",
        "shadow": "0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)"
      }
    }
  },

  "recipes": {
    "linkItem": {
      "base": {
        "background": "rgba(255, 255, 255, 0.7)",
        "backdropFilter": "blur(12px)",
        "border": "1px solid rgba(255, 255, 255, 0.3)",
        "borderRadius": "1rem",
        "padding": "1rem 1.25rem",
        "shadow": "0 4px 6px -1px rgba(0, 0, 0, 0.1)",
        "transition": "all 200ms cubic-bezier(0.4, 0, 0.2, 1)",
        "color": "#1f2937",
        "fontSize": "1rem",
        "fontWeight": 500
      },
      "hover": {
        "background": "rgba(255, 255, 255, 0.85)",
        "transform": "translateY(-2px)",
        "shadow": "0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)"
      },
      "active": {
        "transform": "translateY(0)",
        "shadow": "0 4px 6px -1px rgba(0, 0, 0, 0.1)"
      },
      "variants": {
        "style": {
          "default": {},
          "outline": {
            "background": "transparent",
            "border": "2px solid rgba(255, 255, 255, 0.5)"
          },
          "solid": {
            "background": "#8b5cf6",
            "color": "#ffffff",
            "border": "none"
          },
          "ghost": {
            "background": "transparent",
            "border": "none",
            "shadow": "none"
          }
        }
      }
    },

    "linkGroup": {
      "base": {
        "display": "flex",
        "flexDirection": "column",
        "gap": "0.75rem",
        "width": "100%"
      },
      "variants": {
        "layout": {
          "list": {
            "columns": 1,
            "gap": "0.75rem"
          },
          "cards": {
            "columns": 1,
            "gap": "1rem",
            "cardStyle": "elevated"
          },
          "grid": {
            "columns": 2,
            "gap": "0.75rem"
          }
        }
      }
    },

    "header": {
      "base": {
        "display": "flex",
        "flexDirection": "column",
        "alignItems": "center",
        "gap": "1rem",
        "padding": "1.5rem 0",
        "textAlign": "center"
      },
      "avatar": {
        "size": "5rem",
        "borderRadius": "9999px",
        "border": "3px solid rgba(255, 255, 255, 0.5)",
        "shadow": "0 4px 6px -1px rgba(0, 0, 0, 0.1)"
      },
      "name": {
        "fontSize": "1.5rem",
        "fontWeight": 700,
        "color": "#1f2937",
        "lineHeight": 1.25
      },
      "bio": {
        "fontSize": "1rem",
        "fontWeight": 400,
        "color": "#6b7280",
        "lineHeight": 1.5,
        "maxWidth": "320px"
      }
    },

    "button": {
      "base": {
        "display": "inline-flex",
        "alignItems": "center",
        "justifyContent": "center",
        "gap": "0.5rem",
        "padding": "0.75rem 1.5rem",
        "borderRadius": "9999px",
        "fontWeight": 600,
        "fontSize": "0.875rem",
        "transition": "all 200ms cubic-bezier(0.4, 0, 0.2, 1)",
        "cursor": "pointer"
      },
      "variants": {
        "variant": {
          "primary": {
            "background": "#8b5cf6",
            "color": "#ffffff",
            "border": "none",
            "shadow": "0 0 20px rgba(139, 92, 246, 0.35)"
          },
          "secondary": {
            "background": "rgba(255, 255, 255, 0.7)",
            "color": "#1f2937",
            "border": "1px solid rgba(255, 255, 255, 0.3)",
            "backdropFilter": "blur(12px)"
          },
          "ghost": {
            "background": "transparent",
            "color": "#6b7280",
            "border": "none"
          }
        },
        "size": {
          "sm": {
            "padding": "0.5rem 1rem",
            "fontSize": "0.75rem"
          },
          "md": {
            "padding": "0.75rem 1.5rem",
            "fontSize": "0.875rem"
          },
          "lg": {
            "padding": "1rem 2rem",
            "fontSize": "1rem"
          }
        }
      }
    },

    "badge": {
      "base": {
        "display": "inline-flex",
        "alignItems": "center",
        "padding": "0.25rem 0.75rem",
        "borderRadius": "9999px",
        "fontSize": "0.75rem",
        "fontWeight": 500
      },
      "variants": {
        "variant": {
          "default": {
            "background": "rgba(255, 255, 255, 0.5)",
            "color": "#6b7280"
          },
          "primary": {
            "background": "rgba(139, 92, 246, 0.15)",
            "color": "#8b5cf6"
          },
          "success": {
            "background": "rgba(16, 185, 129, 0.15)",
            "color": "#10b981"
          }
        }
      }
    }
  },

  "page": {
    "layout": {
      "maxWidth": 480,
      "pagePadding": 20,
      "blockGap": 16,
      "textAlign": "center",
      "baseFontSize": "M",
      "contentAlign": "center",
      "safeArea": {
        "top": 0,
        "bottom": 0
      }
    },
    "defaults": {
      "linkGroup": {
        "textAlign": "center",
        "fontSize": "M",
        "radius": "1rem",
        "padding": "1rem",
        "shadow": "md",
        "layout": "list"
      },
      "textBlock": {
        "heading": {
          "fontSize": "1.5rem",
          "fontWeight": 700,
          "color": "#1f2937"
        },
        "body": {
          "fontSize": "1rem",
          "fontWeight": 400,
          "color": "#6b7280"
        },
        "caption": {
          "fontSize": "0.875rem",
          "fontWeight": 400,
          "color": "#9ca3af"
        }
      },
      "imageBlock": {
        "radius": "1rem",
        "shadow": "md",
        "objectFit": "cover"
      },
      "productBlock": {
        "radius": "1rem",
        "shadow": "md",
        "padding": "1rem",
        "background": "rgba(255, 255, 255, 0.7)"
      }
    },
    "mode": "light"
  },

  "background": {
    "type": "gradient",
    "gradient": "linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%)",
    "wallpaper": {
      "kind": "preset",
      "assetId": null
    },
    "effects": {
      "blur": 0,
      "dim": 0.0,
      "overlayColor": "rgba(0, 0, 0, 0)"
    }
  },

  "modes": {
    "dark": {
      "semantic": {
        "color": {
          "text": {
            "default": "#f9fafb",
            "muted": "#9ca3af",
            "subtle": "#6b7280",
            "invert": "#111827"
          },
          "surface": {
            "card": "rgba(31, 41, 55, 0.7)",
            "cardHover": "rgba(31, 41, 55, 0.85)",
            "glass": "rgba(31, 41, 55, 0.25)"
          },
          "border": {
            "default": "rgba(255, 255, 255, 0.1)",
            "subtle": "rgba(255, 255, 255, 0.05)",
            "strong": "rgba(255, 255, 255, 0.2)"
          }
        }
      },
      "recipes": {
        "linkItem": {
          "base": {
            "background": "rgba(31, 41, 55, 0.7)",
            "border": "1px solid rgba(255, 255, 255, 0.1)",
            "color": "#f9fafb"
          },
          "hover": {
            "background": "rgba(31, 41, 55, 0.85)"
          }
        },
        "header": {
          "name": {
            "color": "#f9fafb"
          },
          "bio": {
            "color": "#9ca3af"
          }
        }
      },
      "page": {
        "defaults": {
          "textBlock": {
            "heading": {
              "color": "#f9fafb"
            },
            "body": {
              "color": "#9ca3af"
            },
            "caption": {
              "color": "#6b7280"
            }
          },
          "productBlock": {
            "background": "rgba(31, 41, 55, 0.7)"
          }
        }
      },
      "background": {
        "gradient": "linear-gradient(135deg, #1e1b4b 0%, #312e81 50%, #4c1d95 100%)"
      }
    }
  }
}'::jsonb)

ON CONFLICT (key) DO UPDATE SET
  name = EXCLUDED.name,
  tier = EXCLUDED.tier,
  visibility = EXCLUDED.visibility,
  config = EXCLUDED.config,
  updated_at = NOW();
