import type { PageServerLoad } from './$types';

const API_BASE = 'http://localhost:8080';

export const load: PageServerLoad = async () => {
  try {
    const res = await fetch(`${API_BASE}/api/pages`);
    if (res.ok) {
      const data = await res.json();
      return { pages: data.pages || [] };
    }
  } catch {
    // ignore
  }
  return { pages: [] };
};
