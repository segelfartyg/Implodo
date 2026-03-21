<script>
  import { fetchImages } from './api.js';
  import PostCard from './PostCard.svelte';

  let images = $state([]);
  let loading = $state(true);
  let error = $state('');

  async function load() {
    loading = true;
    error = '';
    try {
      images = await fetchImages();
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  load();
</script>

<section class="feed">
  {#if loading}
    <div class="grid">
      {#each Array(3) as _}
        <div class="skeleton-card">
          <div class="skeleton-img"></div>
        </div>
      {/each}
    </div>
  {:else if error}
    <button class="retry-btn" onclick={load}>{error} — tap to retry</button>
  {:else if images.length === 0}
    <p class="empty-msg">No posts yet. Upload an image above to get started.</p>
  {:else}
    <div class="grid">
      {#each images as image (image.id)}
        <PostCard {image} />
      {/each}
    </div>
  {/if}
</section>

<style>
  .feed {
    width: 100%;
    max-width: 960px;
    margin: 0 auto;
    padding: 16px;
  }

  .grid {
    display: flex;
    flex-direction: column;
    gap: 16px;
    max-width: 600px;
    margin: 0 auto;
  }

  .skeleton-card {
    background: var(--surface);
    border-radius: 12px;
    overflow: hidden;
    border: 1px solid var(--border);
  }

  .skeleton-img {
    width: 100%;
    aspect-ratio: 1 / 1;
    background: var(--skeleton-base);
    animation: shimmer 1.4s infinite;
  }

  @keyframes shimmer {
    0%   { opacity: 0.5; }
    50%  { opacity: 1; }
    100% { opacity: 0.5; }
  }

  .retry-btn {
    display: block;
    margin: 24px auto;
    padding: 10px 20px;
    background: var(--surface-raised);
    color: var(--text);
    border: 1px solid var(--border-hover);
    border-radius: 8px;
    cursor: pointer;
    font-size: 0.9rem;
  }

  .retry-btn:hover { background: var(--border); }

  .empty-msg {
    text-align: center;
    color: var(--text-muted);
    font-size: 0.875rem;
    padding: 48px 0;
  }
</style>
