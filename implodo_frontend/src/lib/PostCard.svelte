<script>
  import { toggleLike } from './api.js';
  import ImageCompare from './ImageCompare.svelte';

  export let post;

  let liked = post.likedByMe;
  let likeCount = post.likeCount;
  let liking = false;

  async function handleLike() {
    if (liking) return;
    liking = true;
    try {
      const result = await toggleLike(post.id, liked, likeCount);
      liked = result.likedByMe;
      likeCount = result.likeCount;
    } finally {
      liking = false;
    }
  }

  function formatCount(n) {
    if (n >= 1_000_000) return (n / 1_000_000).toFixed(1) + 'M';
    if (n >= 1_000) return (n / 1_000).toFixed(1) + 'K';
    return String(n);
  }
</script>

<article class="post-card">
  <ImageCompare afterUrl={post.imageUrl} beforeUrl={post.compareImageUrl} />

  <footer class="post-footer">
    <div class="post-author">
      <div class="avatar-small">
        {post.author.username[0].toUpperCase()}
      </div>
      <span class="username">@{post.author.username}</span>
    </div>

    <button
      class="like-btn"
      class:liked
      aria-label={liked ? 'Unlike' : 'Like'}
      aria-pressed={liked}
      on:click={handleLike}
      disabled={liking}
    >
      <svg class="heart-icon" viewBox="0 0 24 24" aria-hidden="true">
        {#if liked}
          <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5
                   2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09
                   C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5
                   c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
        {:else}
          <path d="M16.5 3c-1.74 0-3.41.81-4.5 2.09C10.91 3.81 9.24 3 7.5 3
                   5.42 3 3 5.42 3 8.5c0 3.78 3.4 6.86 8.55 11.54L12 21.35
                   l1.45-1.32C18.6 15.36 22 12.28 22 8.5 22 5.42 19.58 3 16.5 3z
                   M12 19.6C7.26 15.29 4 12.42 4 8.5 4 6 5.99 4 8.5 4
                   c1.54 0 3.04.99 3.57 2.36h1.87C14.46 4.99 15.96 4 17.5 4
                   20.01 4 22 6 22 8.5c0 3.92-3.26 6.79-8 11.1z"/>
        {/if}
      </svg>
      <span class="like-count">{formatCount(likeCount)}</span>
    </button>
  </footer>
</article>

<style>
  .post-card {
    background: var(--surface);
    border-radius: 12px;
    overflow: hidden;
    border: 1px solid var(--border);
    transition: border-color 0.2s;
  }

  .post-card:hover {
    border-color: var(--border-hover);
  }

  .post-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 14px;
    gap: 8px;
  }

  .post-author {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 0;
  }

  .avatar-small {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    background: var(--accent);
    color: #fff;
    font-size: 0.75rem;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .username {
    font-size: 0.82rem;
    color: var(--text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .like-btn {
    display: flex;
    align-items: center;
    gap: 5px;
    background: none;
    border: none;
    cursor: pointer;
    color: var(--text-muted);
    padding: 4px 6px;
    border-radius: 6px;
    transition: color 0.15s, background 0.15s;
    flex-shrink: 0;
  }

  .like-btn:hover:not(:disabled) {
    background: var(--surface-raised);
    color: var(--text);
  }

  .like-btn.liked {
    color: var(--like-color);
  }

  .like-btn:disabled {
    opacity: 0.6;
    cursor: default;
  }

  .heart-icon {
    width: 20px;
    height: 20px;
    fill: currentColor;
  }

  .like-count {
    font-size: 0.82rem;
    font-weight: 500;
    min-width: 24px;
  }
</style>
