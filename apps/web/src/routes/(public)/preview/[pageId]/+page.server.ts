import type { PageServerLoad } from './$types';

const API_BASE = 'http://localhost:8080';

export const load: PageServerLoad = async ({ params }) => {
  try {
    // Lấy compiled_json trực tiếp từ database qua API mới
    const res = await fetch(`${API_BASE}/api/pages/${params.pageId}/published`);
    
    if (res.ok) {
      const compiled = await res.json();
      return { error: null, compiled };
    }
    
    const error = await res.json();
    return { error: error.error, compiled: null };
  } catch {
    return { error: { code: 'FETCH_ERROR', message: 'Failed to fetch page' }, compiled: null };
  }
};
