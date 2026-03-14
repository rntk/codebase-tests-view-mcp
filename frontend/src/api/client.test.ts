import { afterEach, describe, expect, it, vi } from 'vitest';

import { toggleCommentResolved, updateComment } from './client';

describe('comment API response handling', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('treats updateComment empty 200 responses as success without Content-Length', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(null, { status: 200 })
    );

    await expect(updateComment('file.ts', 'comment-1', 'updated')).resolves.toBeUndefined();
  });

  it('treats toggleCommentResolved empty 200 responses as success without Content-Length', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(null, { status: 200 })
    );

    await expect(toggleCommentResolved('file.ts', 'comment-1')).resolves.toBeUndefined();
  });

  it('still parses JSON bodies for successful requests', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ comments: [] }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    );

    const { getComments } = await import('./client');
    await expect(getComments('file.ts')).resolves.toEqual({ comments: [] });
  });
});
