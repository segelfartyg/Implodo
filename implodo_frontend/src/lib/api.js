/**
 * API layer for Implodo.
 * All functions return a consistent shape so the backend can be swapped in
 * without touching components.
 *
 * Shape returned by fetchPosts:
 *   { posts: Post[], hasMore: boolean, nextPage: number }
 *
 * Post shape:
 *   { id, imageUrl, compareImageUrl, author: { id, username, avatarUrl }, likeCount, likedByMe, createdAt }
 */

const PAGE_SIZE = 10;
const TOTAL_MOCK_POSTS = 47; // finite mock set to demonstrate hasMore logic

/** @param {number} id */
function makeMockPost(id) {
  return {
    id,
    imageUrl: `https://picsum.photos/seed/${id}/600/600`,
    compareImageUrl: `https://picsum.photos/seed/${id + 100}/600/600`,
    author: {
      id: (id % 5) + 1,
      username: ['nova_lens', 'amber.frames', 'pixel_wolf', 'still_light', 'urban_flora'][id % 5],
      avatarUrl: null,
    },
    likeCount: Math.floor(Math.random() * 500),
    likedByMe: false,
    createdAt: new Date(Date.now() - id * 3_600_000).toISOString(),
  };
}

/**
 * Fetch a page of posts.
 * @param {number} page   1-based page number
 * @param {number} [size] items per page (defaults to PAGE_SIZE)
 * @returns {Promise<{ posts: object[], hasMore: boolean, nextPage: number }>}
 */
export async function fetchPosts(page = 1, size = PAGE_SIZE) {
  // Simulate network latency.
  await new Promise((r) => setTimeout(r, 600));

  const start = (page - 1) * size;
  const end = Math.min(start + size, TOTAL_MOCK_POSTS);
  const posts = Array.from({ length: end - start }, (_, i) => makeMockPost(start + i + 1));
  const hasMore = end < TOTAL_MOCK_POSTS;

  return { posts, hasMore, nextPage: page + 1 };
}

/**
 * Toggle like on a post.
 * @param {number|string} postId
 * @param {boolean} currentlyLiked
 * @returns {Promise<{ likedByMe: boolean, likeCount: number }>}
 */
/** @param {number|string} _postId @param {boolean} currentlyLiked @param {number} currentCount */
export async function toggleLike(_postId, currentlyLiked, currentCount) {
  await new Promise((r) => setTimeout(r, 200));
  return {
    likedByMe: !currentlyLiked,
    likeCount: currentlyLiked ? currentCount - 1 : currentCount + 1,
  };
}
