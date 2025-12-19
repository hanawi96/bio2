import type { DraftPayload } from '../../types/draft';

export type EditorStatus = 'idle' | 'saving' | 'saved' | 'error' | 'publishing';

interface EditorSyncState {
  dirty: boolean;
  status: EditorStatus;
  lastSavedAt: string | null;
  lastError: { message: string } | null;
  requestId: number;
}

let state = $state<EditorSyncState>({
  dirty: false,
  status: 'idle',
  lastSavedAt: null,
  lastError: null,
  requestId: 0
});

let pendingSave: { pageId: number; draft: DraftPayload } | null = null;

export function markDirty() {
  state.dirty = true;
  if (state.status === 'saved') {
    state.status = 'idle';
  }
  state.lastError = null;
}

export async function saveDraft(pageId: number, draft: DraftPayload): Promise<boolean> {
  const currentRequestId = ++state.requestId;
  pendingSave = { pageId, draft };

  state.status = 'saving';
  state.lastError = null;

  try {
    const res = await fetch(`/api/pages/${pageId}/save`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(draft)
    });

    if (currentRequestId !== state.requestId) return false;

    if (res.ok) {
      const data = await res.json();
      state.dirty = false;
      state.status = 'saved';
      state.lastSavedAt = data.updated_at || new Date().toISOString();
      pendingSave = null;
      return true;
    } else {
      const error = await res.json();
      state.status = 'error';
      state.lastError = { message: error.error?.message || 'Save failed' };
      return false;
    }
  } catch {
    if (currentRequestId !== state.requestId) return false;
    state.status = 'error';
    state.lastError = { message: 'Network error' };
    return false;
  }
}

export async function publish(pageId: number, draft: DraftPayload): Promise<boolean> {
  const currentRequestId = ++state.requestId;

  state.status = 'publishing';
  state.lastError = null;

  try {
    // Save first
    const saveRes = await fetch(`/api/pages/${pageId}/save`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(draft)
    });

    if (!saveRes.ok) throw new Error('Save failed');

    // Then publish
    const publishRes = await fetch(`/api/pages/${pageId}/publish`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({})
    });

    if (currentRequestId !== state.requestId) return false;

    if (publishRes.ok) {
      state.dirty = false;
      state.status = 'saved';
      state.lastSavedAt = new Date().toISOString();
      return true;
    } else {
      const error = await publishRes.json();
      state.status = 'error';
      state.lastError = { message: error.error?.message || 'Publish failed' };
      return false;
    }
  } catch {
    if (currentRequestId !== state.requestId) return false;
    state.status = 'error';
    state.lastError = { message: 'Network error' };
    return false;
  }
}

export async function retrySave(): Promise<boolean> {
  if (!pendingSave) return false;
  return saveDraft(pendingSave.pageId, pendingSave.draft);
}

export function resetState() {
  state.dirty = false;
  state.status = 'idle';
  state.lastSavedAt = null;
  state.lastError = null;
  pendingSave = null;
}

export function showToast(_message: string, _type: 'success' | 'error' | 'info' = 'info') {
  // No-op for now - can implement later if needed
}

export function getEditorState() {
  return state;
}
