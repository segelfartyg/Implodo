<script>
  import Feed from './lib/Feed.svelte';
  import ProfileBadge from './lib/ProfileBadge.svelte';
</script>

<div class="app">
  <header class="topbar">
    <span class="wordmark">Implodo</span>
    <ProfileBadge username="segelfartyg" />
  </header>

  <main>
    <Feed />
  </main>
</div>

<style>
  .app {
    min-height: 100vh;
    background: var(--bg);
    color: var(--text);
  }

  .topbar {
    position: sticky;
    top: 0;
    z-index: 100;
    height: 56px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
    background: var(--topbar-bg);
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    border-bottom: 1px solid var(--border);
  }

  .wordmark {
    font-size: 1.25rem;
    font-weight: 800;
    letter-spacing: -0.5px;
    color: var(--text);
  }

  main {
    padding-top: 8px;
  }
</style>
