<script>
  import { onMount } from 'svelte';
  import Feed from './lib/Feed.svelte';
  import Login from './lib/Login.svelte';
  import Profile from './lib/Profile.svelte';
  import ProfileBadge from './lib/ProfileBadge.svelte';
  import UploadForm from './lib/UploadForm.svelte';
  import { user, initAuth, handleCallback, logout } from './lib/auth.js';
  import logo from './assets/splodo_logo.png';

  // 'loading' | 'callback' | 'login' | 'feed' | 'upload' | 'profile'
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
    <span class="wordmark">Splodo</span>
    {#if view === 'callback'}
      <p class="splash-label">Signing you in…</p>
    {/if}
  </div>
{:else if view === 'login'}
  <Login initialError={callbackError} />
{:else}
  <div class="app">
    <header class="topbar">
      <span class="wordmark">Splodo
      </span>
      
      <div class="topbar-right">
        <button class="new-post-btn" onclick={() => view = 'upload'}>+ New Post</button>
        <ProfileBadge username={$user?.name ?? ''} onclick={() => view = 'profile'} />
        <button class="logout-btn" onclick={logout}>Sign out</button>
      </div>
    </header>

    <main>
      {#if view === 'profile'}
        <Profile onBack={() => view = 'feed'} />
      {:else if view === 'upload'}
        <UploadForm onUploaded={() => { feedKey++; view = 'feed'; }} onBack={() => view = 'feed'} />
      {:else}
        {#key feedKey}
          <Feed />
        {/key}
      {/if}
    </main>

    {#if view === 'feed'}
      <button class="fab" onclick={() => view = 'upload'} aria-label="New post">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
          <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
      </button>
    {/if}
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
    color: var(--text);
  }

  .splash-label {
    margin: 0;
    color: var(--muted);
    font-size: 0.9rem;
    letter-spacing: 0.03em;
  }

  .app {
    min-height: 100vh;
    color: var(--text);
  }

  @keyframes slideDown {
    from { transform: translateY(-100%); opacity: 0; }
    to   { transform: translateY(0);     opacity: 1; }
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(8px); }
    to   { opacity: 1; transform: translateY(0); }
  }

  .topbar {
    position: sticky;
    top: 0;
    z-index: 100;
    height: 58px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 24px;
    background: var(--topbar-bg);
    backdrop-filter: blur(16px) saturate(1.5);
    -webkit-backdrop-filter: blur(16px) saturate(1.5);
    border-bottom: 1px solid transparent;
    box-shadow:
      0 1px 0 0 var(--border),
      0 1px 32px rgba(0, 0, 0, 0.6);
    animation: slideDown 0.4s cubic-bezier(0.16, 1, 0.3, 1) both;
  }

  .topbar::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, var(--accent-glow), transparent);
    pointer-events: none;
  }

  .wordmark {
    height: 30px;
    font-size: 1.3rem;
    font-weight: 800;
    letter-spacing: 0.01em;
    color: var(--text);
    display: flex;
    align-items: center;
    text-shadow: 0 0 20px var(--accent-glow);
  }

  .topbar-right {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .new-post-btn {
    padding: 6px 14px;
    border-radius: 8px;
    border: 1px solid var(--accent);
    background: var(--accent-dim);
    color: var(--accent);
    font-size: 0.85rem;
    cursor: pointer;
    transition: background 0.2s, box-shadow 0.2s, color 0.2s;
    letter-spacing: 0.02em;
  }

  .new-post-btn:hover {
    background: var(--accent);
    color: #fff;
    box-shadow: 0 0 16px var(--accent-glow);
  }

  .fab {
    position: fixed;
    bottom: 28px;
    right: 28px;
    z-index: 200;
    width: 52px;
    height: 52px;
    border-radius: 50%;
    border: none;
    background: var(--accent);
    color: #fff;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 20px var(--accent-glow), 0 2px 8px rgba(0,0,0,0.5);
    transition: transform 0.2s, box-shadow 0.2s;
    animation: fabPop 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) both;
  }

  .fab:hover {
    transform: scale(1.1);
    box-shadow: 0 6px 28px var(--accent-glow), 0 2px 12px rgba(0,0,0,0.6);
  }

  .fab svg {
    width: 22px;
    height: 22px;
  }

  @keyframes fabPop {
    from { transform: scale(0); opacity: 0; }
    to   { transform: scale(1); opacity: 1; }
  }

  .logout-btn {
    padding: 5px 14px;
    border-radius: 6px;
    border: 1px solid var(--border);
    background: transparent;
    color: var(--text-muted);
    font-size: 0.8rem;
    cursor: pointer;
    transition: border-color 0.2s, color 0.2s, background 0.2s;
    letter-spacing: 0.02em;
  }

  .logout-btn:hover {
    border-color: var(--accent);
    color: var(--text);
    background: var(--accent-dim);
  }

  main {
    padding-top: 12px;
    animation: fadeIn 0.4s cubic-bezier(0.16, 1, 0.3, 1) 0.1s both;
  }
</style>
