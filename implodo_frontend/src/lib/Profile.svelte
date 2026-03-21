<script>
  import { onMount } from 'svelte';
  import { user } from './auth.js';
  import { fetchProfile } from './api.js';

  let bio = $state('');
  let saved = $state(false);

  let {onBack} = $props();

  onMount(async () => {
    const profile = await fetchProfile();
    bio = profile.bio ?? '';
  });

  function save() {
    saved = true;
  }

</script>

<div class="profile-page">
  <button class="back-btn" onclick={onBack}>← Back</button>

  <div class="card">
    <h1 class="name">{$user?.name ?? ''}</h1>
    <p class="email">{$user?.email ?? ''}</p>

    <div class="bio-section">
      <label for="bio">About me</label>
      <textarea
        id="bio"
        bind:value={bio}
        placeholder="Write something about yourself…"
        rows="6"
      ></textarea>
      <button class="btn-primary save-btn" onclick={save}>
        {saved ? 'Saved!' : 'Save'}
      </button>
    </div>
  </div>
</div>

<style>
  .profile-page {
    max-width: 560px;
    margin: 48px auto;
    padding: 0 20px;
  }

  .back-btn {
    background: transparent;
    border: none;
    color: var(--text-muted);
    font-size: 0.9rem;
    cursor: pointer;
    padding: 0;
    margin-bottom: 24px;
    transition: color 0.2s;
  }

  .back-btn:hover {
    color: var(--text);
  }

  .card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 32px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
  }

  .avatar {
    width: 72px;
    height: 72px;
    border-radius: 50%;
    background: var(--accent);
    color: #fff;
    font-size: 1.5rem;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 8px;
  }

  .name {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 700;
    color: var(--text);
  }

  .email {
    margin: 0;
    font-size: 0.85rem;
    color: var(--text-muted);
  }

  .bio-section {
    width: 100%;
    margin-top: 24px;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  label {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  textarea {
    width: 100%;
    background: var(--surface-raised);
    border: 1px solid var(--border);
    border-radius: 8px;
    color: var(--text);
    font-size: 0.95rem;
    font-family: inherit;
    padding: 12px;
    resize: vertical;
    transition: border-color 0.2s;
    box-sizing: border-box;
  }

  textarea:focus {
    outline: none;
    border-color: var(--border-hover);
  }

  .save-btn {
    align-self: flex-end;
  }
</style>
