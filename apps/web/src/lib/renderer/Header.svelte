<script lang="ts">
  import type { HeaderSettings } from '../../types/page';

  let { header }: { header: HeaderSettings } = $props();
</script>

<header class="header">
  {#if header.avatar_asset_id}
    <div class="avatar">
      <!-- TODO: resolve asset URL -->
      <div class="avatar-placeholder"></div>
    </div>
  {/if}

  {#if header.name}
    <h1 class="name">
      {header.name}
      {#if header.verified}
        <span class="verified" title="Verified">✓</span>
      {/if}
    </h1>
  {/if}

  {#if header.bio}
    <p class="bio">{header.bio}</p>
  {/if}

  {#if header.social && header.social.length > 0}
    <div class="social">
      {#each header.social as social}
        <a href={social.url} target="_blank" rel="noopener noreferrer" class="social-link">
          {social.type}
        </a>
      {/each}
    </div>
  {/if}
</header>

<style>
  .header {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 24px 0;
    gap: 8px;
  }

  .avatar {
    width: 96px;
    height: 96px;
    border-radius: 50%;
    overflow: hidden;
    margin-bottom: 8px;
  }

  .avatar-placeholder {
    width: 100%;
    height: 100%;
    background: rgba(255, 255, 255, 0.2);
  }

  .name {
    font-size: 1.5rem;
    font-weight: 600;
    margin: 0;
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .verified {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    background: #3b82f6;
    color: white;
    border-radius: 50%;
    font-size: 12px;
  }

  .bio {
    margin: 0;
    opacity: 0.8;
    max-width: 400px;
  }

  .social {
    display: flex;
    gap: 12px;
    margin-top: 8px;
  }

  .social-link {
    color: inherit;
    opacity: 0.7;
    text-decoration: none;
    font-size: 14px;
  }

  .social-link:hover {
    opacity: 1;
  }
</style>
