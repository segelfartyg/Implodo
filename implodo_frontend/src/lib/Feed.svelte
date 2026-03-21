<script>
  import { fly } from 'svelte/transition';
  import { cubicOut, cubicIn } from 'svelte/easing';

  /** Combined scale + fly transition for individual cards */
  function popIn(node, { delay = 0, duration = 350 } = {}) {
    return {
      delay,
      duration,
      easing: cubicOut,
      css: (t) => `
        opacity: ${t};
        transform: translateY(${(1 - t) * 28}px) scale(${0.88 + t * 0.12});
      `
    };
  }

  /** Swipe-out when layout changes */
  function swipeOut(node, { duration = 180 } = {}) {
    return {
      duration,
      easing: cubicIn,
      css: (t) => `
        opacity: ${t};
        transform: translateY(${(1 - t) * -16}px) scale(${0.94 + t * 0.06});
      `
    };
  }
  import { fetchImages } from './api.js';
  import PostCard from './PostCard.svelte';

  let images = $state([]);
  let loading = $state(true);
  let error = $state('');
  let gridMode = $state(false);
  let showGenerated = $state(false);

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
  {#if !loading && !error && images.length > 0}
    <div class="toolbar">
      <button
        class="tool-btn"
        class:active={showGenerated}
        onclick={() => showGenerated = !showGenerated}
        title="Toggle between original and generated"
      >
        {showGenerated ? 'Generated' : 'Original'}
      </button>
      <div class="layout-toggle">
        <button
          class="tool-btn icon-btn"
          class:active={!gridMode}
          onclick={() => gridMode = false}
          title="List view"
          aria-label="List view"
        >
          <svg viewBox="0 0 20 20" fill="currentColor" width="16" height="16">
            <rect x="2" y="3" width="16" height="3" rx="1"/>
            <rect x="2" y="8.5" width="16" height="3" rx="1"/>
            <rect x="2" y="14" width="16" height="3" rx="1"/>
          </svg>
        </button>
        <button
          class="tool-btn icon-btn"
          class:active={gridMode}
          onclick={() => gridMode = true}
          title="Grid view"
          aria-label="Grid view"
        >
          <svg viewBox="0 0 20 20" fill="currentColor" width="16" height="16">
            <rect x="2" y="2" width="7" height="7" rx="1"/>
            <rect x="11" y="2" width="7" height="7" rx="1"/>
            <rect x="2" y="11" width="7" height="7" rx="1"/>
            <rect x="11" y="11" width="7" height="7" rx="1"/>
          </svg>
        </button>
      </div>
    </div>
  {/if}

  {#if loading}
    <div class="cards" class:cards-grid={gridMode}>
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
    {#key gridMode}
      <div
        class="cards"
        class:cards-grid={gridMode}
        out:swipeOut
        in:fly={{ y: 32, duration: 200, opacity: 1 }}
      >
        {#each images as image, i (image.id)}
          <div in:popIn={{ delay: i * 55, duration: 380 }}>
            <PostCard {image} {showGenerated} />
          </div>
        {/each}
      </div>
    {/key}
  {/if}
</section>

<style>
  .feed {
    width: 100%;
    max-width: 960px;
    margin: 0 auto;
    padding: 16px;
  }

  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    max-width: 600px;
    margin: 0 auto 12px;
  }

  .layout-toggle {
    display: flex;
    gap: 4px;
  }

  .tool-btn {
    padding: 5px 12px;
    border-radius: 6px;
    border: 1px solid var(--border);
    background: transparent;
    color: var(--text-muted);
    font-size: 0.8rem;
    cursor: pointer;
    transition: border-color 0.2s, color 0.2s, background 0.2s;
  }

  .tool-btn:hover {
    border-color: var(--border-hover);
    color: var(--text);
  }

  .tool-btn.active {
    border-color: var(--accent);
    color: var(--text);
    background: var(--surface-raised);
  }

  .icon-btn {
    padding: 5px 8px;
    display: flex;
    align-items: center;
  }

  .cards {
    display: flex;
    flex-direction: column;
    gap: 16px;
    max-width: 600px;
    margin: 0 auto;
  }

  .cards-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    max-width: 960px;
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
