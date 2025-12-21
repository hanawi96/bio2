# Tailwind CSS + CSS Variables Setup

## ✅ Setup Complete!

Bạn giờ có thể dùng cả 2 cách:

### 🎯 CHIẾN LƯỢC: Tailwind + iOS Custom Styles

#### **Khi nào dùng Tailwind:**
- ✅ Layout (flex, grid, positioning)
- ✅ Spacing (margin, padding, gap)
- ✅ Responsive design
- ✅ Common utilities (text-center, hidden, etc.)

#### **Khi nào dùng CSS Variables + Custom Classes:**
- ✅ iOS-specific components (buttons, cards, inputs)
- ✅ Complex animations
- ✅ Component-specific styles
- ✅ Design system tokens

---

## 📚 EXAMPLES

### 1. Tailwind cho Layout
```svelte
<!-- ✅ GOOD: Dùng Tailwind cho layout -->
<div class="flex justify-between items-center mb-4">
  <h3 class="text-base font-semibold">Title</h3>
  <button class="p-2 rounded hover:bg-bg-tertiary">Close</button>
</div>

<div class="grid grid-cols-2 gap-3 my-4">
  <div>Column 1</div>
  <div>Column 2</div>
</div>
```

### 2. CSS Variables cho iOS Components
```svelte
<!-- ✅ GOOD: Dùng custom classes cho iOS components -->
<button class="btn-primary">Primary Button</button>
<button class="btn-secondary">Secondary Button</button>
<div class="card">
  <div class="card-header">Header</div>
  <div class="card-body">Body</div>
</div>
```

### 3. Kết hợp cả 2
```svelte
<!-- ✅ BEST: Kết hợp Tailwind layout + iOS components -->
<div class="flex gap-3 mt-4 pt-3 border-t border-separator">
  <button class="btn-secondary">Cancel</button>
  <button class="btn-primary">Save</button>
</div>

<div class="grid grid-cols-2 gap-3">
  <input type="text" class="color-hex-input" />
  <input type="text" class="color-hex-input" />
</div>
```

---

## 🎨 TAILWIND CLASSES ĐÃ MAP VỚI CSS VARIABLES

### Colors
- `bg-primary` → `var(--color-primary)`
- `text-primary` → `var(--color-primary)`
- `bg-success` → `var(--color-success)`
- `bg-danger` → `var(--color-danger)`
- `bg-bg` → `var(--color-bg)`
- `text-secondary` → `var(--color-text-secondary)`
- `border-separator` → `var(--color-separator)`

### Spacing
- `p-4` → `padding: var(--space-4)` (16px)
- `m-3` → `margin: var(--space-3)` (12px)
- `gap-2` → `gap: var(--space-2)` (8px)
- `mb-4` → `margin-bottom: var(--space-4)` (16px)

### Border Radius
- `rounded` → `var(--radius-md)` (12px)
- `rounded-lg` → `var(--radius-lg)` (16px)
- `rounded-full` → `var(--radius-full)` (9999px)

### Shadows
- `shadow` → `var(--shadow-md)`
- `shadow-lg` → `var(--shadow-lg)`
- `shadow-card` → `var(--shadow-card)`

---

## 📋 BEST PRACTICES

### ✅ DO:
```svelte
<!-- Layout với Tailwind -->
<div class="flex justify-between items-center mb-4">
  <span class="text-sm font-medium">Label</span>
  <button class="btn-icon">Icon</button>
</div>

<!-- Grid với Tailwind -->
<div class="grid grid-cols-2 gap-3">
  <div class="card">Card 1</div>
  <div class="card">Card 2</div>
</div>

<!-- Responsive với Tailwind -->
<div class="hidden md:block lg:flex">
  Responsive content
</div>
```

### ❌ DON'T:
```svelte
<!-- ❌ Không tạo custom classes cho layout đơn giản -->
<style>
  .my-flex-container {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
</style>

<!-- ❌ Không override Tailwind cho iOS components -->
<button class="bg-primary text-white px-4 py-2 rounded">
  <!-- Dùng .btn-primary thay vì -->
</button>
```

---

## 🚀 MIGRATION GUIDE

### Refactor từ custom CSS sang Tailwind:

**Before:**
```svelte
<style>
  .container {
    display: flex;
    gap: 12px;
    margin-top: 16px;
  }
</style>

<div class="container">
  <div>Item 1</div>
  <div>Item 2</div>
</div>
```

**After:**
```svelte
<div class="flex gap-3 mt-4">
  <div>Item 1</div>
  <div>Item 2</div>
</div>
```

---

## 💡 TIPS

1. **Dùng Tailwind cho 80% layout code**
2. **Giữ iOS custom classes cho buttons, cards, inputs**
3. **Kết hợp cả 2 khi cần**
4. **Không duplicate - nếu Tailwind có sẵn thì dùng**
5. **CSS Variables vẫn là source of truth cho design tokens**

---

## 🎯 RESULT

- ✅ Code ngắn gọn hơn 50%
- ✅ Consistent spacing & colors
- ✅ Responsive dễ dàng
- ✅ Giữ được iOS design system
- ✅ Best of both worlds!

- `rounded` → `var(--radius-md)` (12px)
- `rounded-lg` → `var(--radius-lg)` (16px)
- `rounded-xl` → `var(--radius-xl)` (20px)

### Shadows
- `shadow-sm` → `var(--shadow-sm)`
- `shadow-md` → `var(--shadow-md)`
- `shadow-lg` → `var(--shadow-lg)`

### Font Size
- `text-xs` → `var(--text-xs)`
- `text-sm` → `var(--text-sm)`
- `text-base` → `var(--text-base)`

## Ví dụ Component với Tailwind:

```svelte
<script lang="ts">
  let count = $state(0);
</script>

<div class="max-w-md mx-auto p-6 bg-bg-secondary rounded-xl shadow-card">
  <h2 class="text-base font-semibold text-primary mb-3">
    Counter Example
  </h2>
  
  <div class="flex items-center gap-4">
    <button 
      class="px-4 py-2 bg-primary text-white rounded-md hover:opacity-90 transition"
      onclick={() => count++}
    >
      Increment
    </button>
    
    <span class="text-2xl font-bold">{count}</span>
    
    <button 
      class="px-4 py-2 bg-danger text-white rounded-md hover:opacity-90 transition"
      onclick={() => count--}
    >
      Decrement
    </button>
  </div>
</div>
```

## Lợi ích:

1. ✅ **Nhanh hơn**: Dùng Tailwind classes thay vì viết CSS
2. ✅ **Nhất quán**: Vẫn dùng design system qua CSS variables
3. ✅ **Linh hoạt**: Có thể mix cả 2 cách
4. ✅ **Responsive**: Dùng `md:`, `lg:` prefixes của Tailwind
5. ✅ **Dark mode**: Dễ dàng thêm `dark:` variants sau này

## Tips:

- Dùng Tailwind cho layout, spacing, flex/grid
- Dùng CSS Variables cho colors, shadows để dễ theme
- Kết hợp cả 2 khi cần thiết
- Component phức tạp vẫn có thể dùng `<style>` scoped

Enjoy! 🎉
