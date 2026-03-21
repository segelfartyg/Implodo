<script>
  import { getToken } from './auth.js';
  import ImageCompare from './ImageCompare.svelte';

  /** @type {{ onUploaded?: () => void }} */
  let { onUploaded } = $props();

  let file = $state(null);
  let preview = $state('');
  let status = $state('idle'); // 'idle' | 'uploading' | 'done' | 'error'
  let error = $state('');
  let result = $state(null); // { original_url, generated_url }
  let dragging = $state(false);

  let inputEl = $state();

  /** @param {Event} e */
  function onFileChange(e) {
    setFile(/** @type {HTMLInputElement} */ (e.target).files[0]);
  }

  /** @param {DragEvent} e */
  function onDrop(e) {
    dragging = false;
    const dropped = e.dataTransfer?.files[0];
    if (dropped) setFile(dropped);
  }

  /** @param {File} f */
  function setFile(f) {
    if (!f || !f.type.startsWith('image/')) {
      error = 'Please select an image file.';
      return;
    }
    file = f;
    error = '';
    result = null;
    status = 'idle';
    const reader = new FileReader();
    reader.onload = (e) => (preview = /** @type {string} */ (e.target?.result));
    reader.readAsDataURL(f);
  }

  async function upload() {
    if (!file) return;
    status = 'uploading';
    error = '';

    const form = new FormData();
    form.append('image', file);

    try {
      const res = await fetch('/api/upload', {
        method: 'POST',
        headers: { Authorization: `Bearer ${getToken()}` },
        body: form,
      });

      const body = await res.json();
      if (!res.ok) throw new Error(body.error ?? 'Upload failed');

      result = body; // { original_url, generated_url } — public GCS URLs
      status = 'done';
      onUploaded?.();
      file = null;
      preview = '';
    } catch (e) {
      error = e.message;
      status = 'error';
    }
  }

  function reset() {
    file = null;
    preview = '';
    result = null;
    status = 'idle';
    error = '';
  }
</script>

<div class="upload-card">
  <h2 class="upload-title">New post</h2>

  {#if status === 'done' && result}
    <div class="result">
      <p class="result-label">Original vs generated</p>
      <ImageCompare beforeUrl={result.original_url} afterUrl={result.generated_url} />
      <button class="btn-secondary" onclick={reset}>Upload another</button>
    </div>
  {:else}
    <!-- Drop zone -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      class="dropzone"
      class:dragging
      class:has-preview={!!preview}
      ondragover={(e) => { e.preventDefault(); dragging = true; }}
      ondragleave={() => (dragging = false)}
      ondrop={(e) => { e.preventDefault(); onDrop(e); }}
      onclick={() => inputEl.click()}
      onkeydown={(e) => e.key === 'Enter' && inputEl.click()}
      role="button"
      tabindex="0"
      aria-label="Select image"
    >
      {#if preview}
        <img class="preview-img" src={preview} alt="Preview" />
      {:else}
        <div class="dropzone-prompt">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round"
              d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" />
          </svg>
          <span>Drop an image here or click to browse</span>
        </div>
      {/if}
    </div>

    <input
      bind:this={inputEl}
      type="file"
      accept="image/jpeg,image/png,image/gif,image/webp"
      onchange={onFileChange}
      class="file-input"
      aria-hidden="true"
      tabindex="-1"
    />

    {#if error}
      <p class="error">{error}</p>
    {/if}

    <div class="actions">
      {#if preview}
        <button class="btn-ghost" onclick={reset} disabled={status === 'uploading'}>Clear</button>
      {/if}
      <button
        class="btn-primary"
        onclick={upload}
        disabled={!file || status === 'uploading'}
      >
        {#if status === 'uploading'}
          <span class="spinner" aria-hidden="true"></span>
          Generating…
        {:else}
          Upload & generate
        {/if}
      </button>
    </div>
  {/if}
</div>

<style>
  .upload-card {
    max-width: 560px;
    margin: 24px auto 0;
    padding: 20px;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 12px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .upload-title {
    font-size: 1rem;
    font-weight: 700;
    color: var(--text);
  }

  /* --- Drop zone --- */
  .dropzone {
    border: 1.5px dashed var(--border-hover);
    border-radius: 8px;
    min-height: 160px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: border-color 0.2s, background 0.2s;
    overflow: hidden;
    outline: none;
  }

  .dropzone:hover,
  .dropzone:focus-visible {
    border-color: var(--accent);
    background: rgba(124, 92, 252, 0.04);
  }

  .dropzone.dragging {
    border-color: var(--accent);
    background: rgba(124, 92, 252, 0.08);
  }

  .dropzone.has-preview {
    border-style: solid;
    min-height: unset;
  }

  .dropzone-prompt {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
    color: var(--text-muted);
    font-size: 0.875rem;
    padding: 32px 20px;
    text-align: center;
  }

  .dropzone-prompt svg {
    width: 32px;
    height: 32px;
    opacity: 0.5;
  }

  .preview-img {
    width: 100%;
    max-height: 360px;
    object-fit: contain;
    display: block;
  }

  /* --- Hidden file input --- */
  .file-input {
    display: none;
  }

  /* --- Actions --- */
  .actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }

  .btn-primary {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 18px;
    border-radius: 7px;
    border: none;
    background: var(--accent);
    color: #fff;
    font-weight: 600;
    font-size: 0.875rem;
    cursor: pointer;
    transition: opacity 0.2s;
  }

  .btn-primary:disabled {
    opacity: 0.45;
    cursor: not-allowed;
  }

  .btn-ghost {
    padding: 8px 14px;
    border-radius: 7px;
    border: 1px solid var(--border);
    background: transparent;
    color: var(--text-muted);
    font-size: 0.875rem;
    cursor: pointer;
    transition: border-color 0.2s, color 0.2s;
  }

  .btn-ghost:hover:not(:disabled) {
    border-color: var(--border-hover);
    color: var(--text);
  }

  .btn-ghost:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  /* --- Spinner --- */
  .spinner {
    width: 14px;
    height: 14px;
    border: 2px solid rgba(255,255,255,0.3);
    border-top-color: #fff;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
    flex-shrink: 0;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* --- Error --- */
  .error {
    font-size: 0.85rem;
    color: var(--like-color);
  }

  /* --- Result --- */
  .result {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .result-label {
    font-size: 0.8rem;
    color: var(--text-muted);
  }

  .btn-secondary {
    align-self: flex-end;
    padding: 7px 14px;
    border-radius: 7px;
    border: 1px solid var(--border);
    background: transparent;
    color: var(--text);
    font-size: 0.85rem;
    cursor: pointer;
    transition: border-color 0.2s;
  }

  .btn-secondary:hover {
    border-color: var(--border-hover);
  }
</style>
