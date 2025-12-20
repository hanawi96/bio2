# Tailwind CSS + CSS Variables Setup

## ✅ Setup Complete!

Bạn giờ có thể dùng cả 2 cách:

### 1. Tailwind Utility Classes (sử dụng CSS Variables)
```svelte
<div class="bg-primary text-white p-4 rounded-lg shadow-md">
  <h1 class="text-base font-semibold">Hello Tailwind</h1>
  <p class="text-sm text-secondary">Using CSS variables</p>
</div>
```

### 2. CSS Variables (như trước)
```svelte
<style>
  .custom {
    background: var(--color-primary);
    padding: var(--space-4);
    border-radius: var(--radius-lg);
  }
</style>
```

### 3. Kết hợp cả 2
```svelte
<div class="flex items-center gap-3 p-4" style="background: var(--color-bg-secondary)">
  <button class="bg-primary text-white px-4 py-2 rounded-md hover:opacity-90">
    Click me
  </button>
</div>
```

## Tailwind Classes đã map với CSS Variables:

### Colors
- `bg-primary` → `var(--color-primary)`
- `text-primary` → `var(--color-primary)`
- `bg-success` → `var(--color-success)`
- `bg-danger` → `var(--color-danger)`
- `bg-bg` → `var(--color-bg)`
- `text-secondary` → `var(--color-text-secondary)`

### Spacing
- `p-4` → `padding: var(--space-4)` (16px)
- `m-3` → `margin: var(--space-3)` (12px)
- `gap-2` → `gap: var(--space-2)` (8px)

### Border Radius
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
