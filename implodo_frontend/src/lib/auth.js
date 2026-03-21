import { writable } from 'svelte/store';

const TOKEN_KEY = 'implodo_jwt';

// --- PKCE helpers ---

function base64url(buffer) {
  return btoa(String.fromCharCode(...new Uint8Array(buffer)))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '');
}

function generateState() {
  const buf = new Uint8Array(16);
  crypto.getRandomValues(buf);
  return base64url(buf.buffer);
}

function generateCodeVerifier() {
  const buf = new Uint8Array(32);
  crypto.getRandomValues(buf);
  return base64url(buf.buffer);
}

async function generateCodeChallenge(verifier) {
  const data = new TextEncoder().encode(verifier);
  const hash = await crypto.subtle.digest('SHA-256', data);
  return base64url(hash);
}

// --- Token storage ---

export function getToken() {
  return localStorage.getItem(TOKEN_KEY);
}

function setToken(token) {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

function parseTokenPayload(token) {
  try {
    const [, payload] = token.split('.');
    return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')));
  } catch {
    return null;
  }
}

// --- Auth store ---

/** null = loading, false = logged out, { name, email, google_id } = logged in */
export const user = writable(null);

/**
 * Check for a valid stored JWT and populate the user store.
 * Call once on app startup (outside of the callback flow).
 */
export function initAuth() {
  const token = getToken();
  if (!token) {
    user.set(false);
    return false;
  }
  const payload = parseTokenPayload(token);
  if (!payload || Date.now() >= payload.exp * 1000) {
    clearToken();
    user.set(false);
    return false;
  }
  user.set({ name: payload.name, email: payload.email, google_id: payload.google_id });
  return true;
}

/**
 * Start the Google OAuth flow. Generates PKCE params, registers the session with
 * the backend, then redirects the browser to Google's consent page.
 */
export async function login() {
  const state = generateState();
  const verifier = generateCodeVerifier();
  const challenge = await generateCodeChallenge(verifier);

  sessionStorage.setItem('auth_state', state);
  sessionStorage.setItem('auth_verifier', verifier);

  const res = await fetch('/auth/start', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ state, code_challenge: challenge }),
  });

  if (!res.ok) throw new Error('Failed to start auth');

  const { auth_url } = await res.json();
  window.location.href = auth_url;
}

/**
 * Complete the auth flow after the backend redirects back here with ?state=<state>.
 * Polls /auth/token until the JWT is ready (max ~30 s), then stores it.
 */
export async function handleCallback(state) {
  const savedState = sessionStorage.getItem('auth_state');
  const verifier = sessionStorage.getItem('auth_verifier');

  if (!verifier || savedState !== state) {
    throw new Error('Invalid callback — state mismatch');
  }

  for (let i = 0; i < 15; i++) {
    const res = await fetch('/auth/token', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ state, code_verifier: verifier }),
    });

    if (res.status === 202) {
      await new Promise((r) => setTimeout(r, 2000));
      continue;
    }

    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      throw new Error(body.error ?? 'Token exchange failed');
    }

    const { token } = await res.json();
    console.log(token);
    setToken(token);
    sessionStorage.removeItem('auth_state');
    sessionStorage.removeItem('auth_verifier');

    const payload = parseTokenPayload(token);
    const u = { name: payload.name, email: payload.email, google_id: payload.google_id };
    user.set(u);
    return u;
  }

  throw new Error('Auth timed out — please try again');
}

export function logout() {
  clearToken();
  user.set(false);
}
