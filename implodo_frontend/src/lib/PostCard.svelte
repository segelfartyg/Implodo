<script>
  import ImageCompare from './ImageCompare.svelte';

  /** @type {{ id: string, original_url: string, generated_url: string }} */
  let { image } = $props();
</script>

<article class="post-card">
  <ImageCompare beforeUrl={image.original_url} afterUrl={image.generated_url} />
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
</style>
