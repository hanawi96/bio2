import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/private';

const API_BASE = env.API_URL || 'http://localhost:8080';

export const load: PageServerLoad = async ({ request, url }) => {
  const host = request.headers.get('host') || url.searchParams.get('host') || 'localhost';
  
  try {
    const response = await fetch(`${API_BASE}/r?host=${encodeURIComponent(host)}`, {
      headers: {
        'Host': host
      }
    });

    if (!response.ok) {
      const error = await response.json();
      return {
        error: error.error,
        compiled: null
      };
    }

    const compiled = await response.json();
    return {
      error: null,
      compiled
    };
  } catch (err) {
    return {
      error: { code: 'FETCH_ERROR', message: 'Failed to fetch page' },
      compiled: null
    };
  }
};
