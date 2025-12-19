import type { PageServerLoad } from './$types';

const API_BASE = 'http://localhost:8080';

export const load: PageServerLoad = async ({ params }) => {
  try {
    const res = await fetch(`${API_BASE}/api/pages/${params.pageId}/draft`);
    if (res.ok) {
      const draft = await res.json();
      return { draft, error: null };
    }
    const error = await res.json();
    return { draft: null, error: error.error };
  } catch {
    return { draft: null, error: { code: 'FETCH_ERROR', message: 'Failed to load draft' } };
  }
};
