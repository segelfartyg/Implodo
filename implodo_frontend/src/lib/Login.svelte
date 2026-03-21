<script>
  import { login } from './auth.js';

  /** @type {{ initialError?: string }} */
  let { initialError = '' } = $props();

  let loading = $state(false);
  let error = $state(initialError);

  async function handleLogin() {
    loading = true;
    error = '';
    try {
      await login();
    } catch (e) {
      error = e.message;
      loading = false;
    }
  }
</script>

<div class="login-page">
  <h1 class="wordmark animate-up" style="animation-delay: 0ms">Splodo</h1>
  <p class="tagline animate-up" style="animation-delay: 100ms">Sign in to continue</p>

  <button class="login-btn animate-up" style="animation-delay: 220ms" onclick={handleLogin} disabled={loading}>
    {#if loading}
      Redirecting…
    {:else}
      <svg class="google-icon" viewBox="0 0 24 24" aria-hidden="true">
        <path
          fill="#4285F4"
          d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
        />
        <path
          fill="#34A853"
          d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
        />
        <path
          fill="#FBBC05"
          d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z"
        />
        <path
          fill="#EA4335"
          d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
        />
      </svg>
      Sign in with Google
    {/if}
  </button>

  {#if error}
    <p class="error">{error}</p>
  {/if}
</div>

<style>
  .login-page {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 20px;
    color: var(--text);
    /* centered glow behind the wordmark */
    background-image: radial-gradient(ellipse 60% 40% at 50% 42%, rgba(224, 48, 48, 0.12), transparent 70%);
  }

  .wordmark {
    font-size: 3.5rem;
    font-weight: 800;
    letter-spacing: 0.02em;
    margin: 0;
    text-shadow:
      0 0 40px var(--accent-glow),
      0 0 80px rgba(224, 48, 48, 0.12);
  }

  .tagline {
    margin: -8px 0 0;
    color: var(--text-muted);
    font-size: 0.9rem;
    letter-spacing: 0.06em;
    text-transform: uppercase;
  }

  .login-btn {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-top: 12px;
    padding: 11px 22px;
    border-radius: 10px;
    border: 1px solid var(--border-hover);
    background: var(--surface-raised);
    color: var(--text);
    font-size: 0.95rem;
    cursor: pointer;
    transition: border-color 0.2s, background 0.2s, box-shadow 0.2s;
    box-shadow: var(--shadow-sm);
  }

  .login-btn:hover:not(:disabled) {
    border-color: var(--accent);
    background: var(--surface);
    box-shadow: 0 0 20px var(--accent-glow);
  }

  .login-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .google-icon {
    width: 18px;
    height: 18px;
    flex-shrink: 0;
  }

  .error {
    color: var(--like-color);
    font-size: 0.875rem;
    margin: 0;
  }

  /* ── Entrance animations ── */
  @keyframes fadeUp {
    from { opacity: 0; transform: translateY(18px); }
    to   { opacity: 1; transform: translateY(0); }
  }

  .animate-up {
    opacity: 0;
    animation: fadeUp 0.55s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  }

  /* Subtle wordmark glow pulse */
  @keyframes glowPulse {
    0%, 100% { text-shadow: 0 0 40px var(--accent-glow), 0 0 80px rgba(224, 48, 48, 0.12); }
    50%       { text-shadow: 0 0 60px rgba(224, 48, 48, 0.45), 0 0 120px rgba(224, 48, 48, 0.2); }
  }

  .wordmark {
    animation: fadeUp 0.55s cubic-bezier(0.16, 1, 0.3, 1) forwards,
               glowPulse 4s ease-in-out 0.6s infinite;
  }
</style>
