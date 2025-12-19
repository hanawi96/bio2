// Block type enums

export const BlockTypes = {
  LINK_GROUP: 'link_group',
  TEXT: 'text',
  IMAGE: 'image',
  PRODUCT: 'product',
  SPACER: 'spacer',
  EMBED: 'embed',
  SOCIAL_ROW: 'social_row',
  FORM: 'form',
} as const;

export type BlockType = (typeof BlockTypes)[keyof typeof BlockTypes];

export const TextVariants = {
  BODY: 'body',
  HEADING: 'heading',
  CAPTION: 'caption',
} as const;

export type TextVariant = (typeof TextVariants)[keyof typeof TextVariants];

export const LayoutTypes = {
  LIST: 'list',
  CARDS: 'cards',
  GRID: 'grid',
} as const;

export type LayoutType = (typeof LayoutTypes)[keyof typeof LayoutTypes];
