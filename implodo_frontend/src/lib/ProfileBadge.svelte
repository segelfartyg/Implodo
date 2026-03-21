<script>
  // In the future this would receive user data from an auth store.
  // For now it accepts props with sensible defaults.
  export let username = 'You';
  export let avatarUrl = null;

  $: initials = username
    .split(/[\s._-]+/)
    .slice(0, 2)
    .map((w) => w[0]?.toUpperCase() ?? '')
    .join('');
</script>

<div class="profile-badge" title={username}>
  {#if avatarUrl}
    <img class="avatar" src={avatarUrl} alt={username} />
  {:else}
    <div class="avatar avatar-placeholder" aria-label={username}>
      {initials}
    </div>
  {/if}
  <span class="username">{username}</span>
</div>

<style>
  .profile-badge {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    padding: 4px 10px 4px 4px;
    border-radius: 24px;
    border: 1px solid var(--border);
    background: var(--surface);
    transition: border-color 0.2s, background 0.2s;
    user-select: none;
  }

  .profile-badge:hover {
    border-color: var(--border-hover);
    background: var(--surface-raised);
  }

  .avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    object-fit: cover;
    display: block;
  }

  .avatar-placeholder {
    background: var(--accent);
    color: #fff;
    font-size: 0.8rem;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .username {
    font-size: 0.85rem;
    font-weight: 500;
    color: var(--text);
    max-width: 120px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
