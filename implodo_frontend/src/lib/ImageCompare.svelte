<script>
  export let beforeUrl; // revealed on the right as you drag left
  export let afterUrl;  // visible by default on the left

  let position = 50; // divider % from left
  let dragging = false;
  let container;

  /** @param {number} v @param {number} lo @param {number} hi */
  function clamp(v, lo, hi) {
    return Math.max(lo, Math.min(hi, v));
  }

  /** @param {PointerEvent} e */
  function onPointerDown(e) {
    dragging = true;
    container.setPointerCapture(e.pointerId);
  }

  /** @param {PointerEvent} e */
  function onPointerMove(e) {
    if (!dragging) return;
    const rect = container.getBoundingClientRect();
    position = clamp(((e.clientX - rect.left) / rect.width) * 100, 0, 100);
  }

  function onPointerUp() {
    dragging = false;
  }
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
  class="compare"
  bind:this={container}
  on:pointerdown={onPointerDown}
  on:pointermove={onPointerMove}
  on:pointerup={onPointerUp}
  on:pointercancel={onPointerUp}
>
  <!-- Back image — always fully visible -->
  <img class="img" src={beforeUrl} alt="Before" draggable="false" />

  <!-- Front image — clipped to the left of the divider via clip-path -->
  <img
    class="img"
    src={afterUrl}
    alt="After"
    draggable="false"
    style="clip-path: inset(0 {100 - position}% 0 0)"
  />

  <!-- Divider line + handle -->
  <div class="divider" style="left: {position}%" class:dragging>
    <div class="handle">
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="M8 5l-5 7 5 7M16 5l5 7-5 7"
              stroke="currentColor" stroke-width="2.5"
              stroke-linecap="round" stroke-linejoin="round" fill="none"/>
      </svg>
    </div>
  </div>
</div>

<style>
  .compare {
    position: relative;
    width: 100%;
    aspect-ratio: 1 / 1;
    overflow: hidden;
    cursor: col-resize;
    user-select: none;
    touch-action: pan-y;
    background: var(--surface-raised);
  }

  .img {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    pointer-events: none;
    draggable: false;
  }

  .divider {
    position: absolute;
    top: 0;
    bottom: 0;
    width: 2px;
    background: rgba(255, 255, 255, 0.85);
    transform: translateX(-50%);
    pointer-events: none;
    box-shadow: 0 0 8px rgba(0, 0, 0, 0.5);
  }

  .handle {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 38px;
    height: 38px;
    background: #fff;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.45);
    color: #222;
    transition: box-shadow 0.15s;
  }

  .divider.dragging .handle {
    box-shadow: 0 2px 16px rgba(0, 0, 0, 0.6);
  }

  .handle svg {
    width: 18px;
    height: 18px;
  }
</style>
