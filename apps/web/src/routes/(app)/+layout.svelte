<script>
  import '../../app.css';
</script>

<div class="app-layout">
  <nav class="app-nav">
    <a href="/app/pages" class="nav-link">Pages</a>
  </nav>
  <main class="app-main">
    <slot />
  </main>
</div>

<style>
  .app-layout {
    min-height: 100vh;
    background: #f5f5f5;
  }

  .app-nav {
    background: white;
    padding: 12px 24px;
    border-bottom: 1px solid #e0e0e0;
  }

  .nav-link {
    color: #333;
    text-decoration: none;
    font-weight: 500;
  }

  .app-main {
    padding: 24px;
    max-width: 1200px;
    margin: 0 auto;
  }
</style>
