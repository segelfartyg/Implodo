<script>
  import ImageCompare from './ImageCompare.svelte';

  /** @type {{ id: string, original_url: string, generated_url: string }} */
  export let image;
  export let showGenerated = false;
</script>

<article class="post-card">
  <ImageCompare beforeUrl={image.original_url} afterUrl={image.generated_url} showAfter={showGenerated} />
</article>

<style>
  .post-card {
    background: var(--surface);
    border-radius: 12px;
    overflow: hidden;
    border: 1px solid var(--border);
    transition: border-color 0.25s, transform 0.25s, box-shadow 0.25s;
    will-change: transform;
  }

  .post-card:hover {
    border-color: var(--border-hover);
    transform: translateY(-3px);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5), 0 0 0 1px var(--border-hover);
  }
</style>
