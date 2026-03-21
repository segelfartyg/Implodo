<script>
  import { onDestroy } from 'svelte';
  import { fetchPosts } from './api.js';
  import PostCard from './PostCard.svelte';

  let posts = [];
  let page = 1;
  let hasMore = true;
  let loading = false;
  let error = null;

  // Sentinel element watched by IntersectionObserver.
  let sentinel;

  async function loadMore() {
    if (loading || !hasMore) return;
    loading = true;
    error = null;
    try {
      const result = await fetchPosts(page);
      posts = [...posts, ...result.posts];
      hasMore = result.hasMore;
      page = result.nextPage;
    } catch (e) {
      error = 'Failed to load posts. Tap to retry.';
    } finally {
      loading = false;
    }
  }

  // Kick off the first load immediately.
  loadMore();

  // Set up IntersectionObserver once the sentinel is in the DOM.
  let observer;
  $: if (sentinel) {
    observer?.disconnect();
    observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) loadMore();
      },
      { rootMargin: '300px' }
    );
    observer.observe(sentinel);
  }

  onDestroy(() => observer?.disconnect());
</script>

<section class="feed">
  <div class="grid">
    {#each posts as post (post.id)}
      <PostCard {post} />
    {/each}
  </div>

  <!-- Loading skeleton cards -->
  {#if loading}
    <div class="grid skeleton-grid">
      {#each Array(6) as _}
        <div class="skeleton-card">
          <div class="skeleton-img"></div>
          <div class="skeleton-footer">
            <div class="skeleton-line short"></div>
            <div class="skeleton-line tiny"></div>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  {#if error}
    <button class="retry-btn" on:click={loadMore}>{error}</button>
  {/if}

  {#if !hasMore && posts.length > 0}
    <p class="end-msg">You've seen everything ✓</p>
  {/if}

  <!-- Invisible sentinel that triggers the next page load -->
  <div bind:this={sentinel} class="sentinel" aria-hidden="true"></div>
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

  /* --- Skeleton --- */
  .skeleton-grid {
    margin-top: 0;
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

  .skeleton-footer {
    padding: 10px 14px;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .skeleton-line {
    height: 10px;
    border-radius: 4px;
    background: var(--skeleton-base);
    animation: shimmer 1.4s infinite;
  }

  .skeleton-line.short { width: 55%; }
  .skeleton-line.tiny  { width: 28%; }

  @keyframes shimmer {
    0%   { opacity: 0.5; }
    50%  { opacity: 1; }
    100% { opacity: 0.5; }
  }

  /* --- Misc --- */
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

  .end-msg {
    text-align: center;
    color: var(--text-muted);
    font-size: 0.85rem;
    padding: 32px 0;
  }

  .sentinel {
    height: 1px;
  }
</style>
