import { getToken } from './auth.js';

function authHeaders() {
  return { Authorization: `Bearer ${getToken()}` };
}

/**
 * Fetch all image pairs for the authenticated user.
 * @returns {Promise<Array<{ id: string, original_url: string, generated_url: string }>>}
 */
export async function fetchImages() {
  const res = await fetch('/api/images', { headers: authHeaders() });
  if (!res.ok) throw new Error('Failed to load images');
  return res.json();
}
