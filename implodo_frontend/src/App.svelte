<script>
  import { onMount } from 'svelte';
  import Feed from './lib/Feed.svelte';
  import Login from './lib/Login.svelte';
  import ProfileBadge from './lib/ProfileBadge.svelte';
  import UploadForm from './lib/UploadForm.svelte';
  import { user, initAuth, handleCallback, logout } from './lib/auth.js';

  // 'loading' | 'callback' | 'login' | 'feed'
  let view = $state('loading');
  let callbackError = $state('');
  let feedKey = $state(0);

  onMount(async () => {
    const params = new URLSearchParams(window.location.search);
    const state = params.get('state');

    if (state) {
      // We've been redirected back from the backend after Google OAuth.
      view = 'callback';
      // Clean the URL before anything async so a refresh doesn't re-trigger the flow.
      history.replaceState(null, '', window.location.pathname);
      try {
        await handleCallback(state);
        view = 'feed';
      } catch (e) {
        callbackError = e.message;
        view = 'login';
      }
      return;
    }

    // Normal load — check for an existing valid JWT.
    initAuth();
  });

  $effect(() => {
    const u = $user;
    if (u === null) return; // still initialising
    if (view === 'loading') {
      view = u ? 'feed' : 'login';
    }
  });
</script>

{#if view === 'loading' || view === 'callback'}
  <div class="splash">
    <span class="wordmark">Implodo</span>
    {#if view === 'callback'}
      <p class="splash-label">Signing you in…</p>
    {/if}
  </div>
{:else if view === 'login'}
  <Login initialError={callbackError} />
{:else}
  <div class="app">
    <header class="topbar">
      <span class="wordmark">Implodo</span>
      <div class="topbar-right">
        <ProfileBadge username={$user?.name ?? ''} />
        <button class="logout-btn" onclick={logout}>Sign out</button>
      </div>
    </header>

    <main>
      <UploadForm onUploaded={() => feedKey++} />
      {#key feedKey}
        <Feed />
      {/key}
    </main>
  </div>
{/if}

<style>
  .splash {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    background: var(--bg);
    color: var(--text);
  }

  .splash-label {
    margin: 0;
    color: var(--muted);
    font-size: 0.9rem;
  }

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

  .topbar-right {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .logout-btn {
    padding: 5px 12px;
    border-radius: 6px;
    border: 1px solid var(--border);
    background: transparent;
    color: var(--muted);
    font-size: 0.8rem;
    cursor: pointer;
    transition: border-color 0.2s, color 0.2s;
  }

  .logout-btn:hover {
    border-color: var(--border-hover);
    color: var(--text);
  }

  main {
    padding-top: 8px;
  }
</style>
