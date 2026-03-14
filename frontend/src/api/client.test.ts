import { afterEach, describe, expect, it, vi } from 'vitest';

import {
  listFiles,
  getFileContent,
  getRelatedTests,
  getSuggestions,
  getSourceReferences,
  getComments,
  createComment,
  updateComment,
  deleteComment,
  toggleCommentResolved,
  exportContext,
  getProjectOverview,
  getMetadataIssues,
  updateSourcePath,
  updateTestPath,
  deleteSourcePath,
  deleteTestPath,
  search,
} from './client';
import type {
  ListFilesResponse,
  FileResponse,
  TestsResponse,
  SuggestionsResponse,
  TestFileResponse,
  CommentsResponse,
  CommentResponse,
  ExportContextResponse,
  OverviewResponse,
  MetadataIssuesResponse,
  SearchResponse,
} from '../types';

describe('API Client', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('handleResponse', () => {
    it('throws error on non-ok response', async () => {
      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response('Not Found', { status: 404, statusText: 'Not Found' })
      );

      await expect(listFiles()).rejects.toThrow('Not Found');
    });

    it('throws error with parsed error message from JSON', async () => {
      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify({ error: 'Custom error message', code: 'ERR_001' }), {
          status: 400,
          statusText: 'Bad Request',
        })
      );

      await expect(listFiles()).rejects.toThrow('Custom error message');
    });

    it('handles 204 No Content response', async () => {
      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(null, { status: 204 })
      );

      await expect(updateComment('file.ts', 'comment-1', 'updated')).resolves.toBeUndefined();
    });

    it('handles empty body with content-length 0', async () => {
      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response('', { status: 200, headers: { 'content-length': '0' } })
      );

      await expect(updateComment('file.ts', 'comment-1', 'updated')).resolves.toBeUndefined();
    });

    it('handles empty string response', async () => {
      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response('', { status: 200 })
      );

      await expect(updateComment('file.ts', 'comment-1', 'updated')).resolves.toBeUndefined();
    });

    it('handles whitespace-only response', async () => {
      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response('   ', { status: 200 })
      );

      await expect(updateComment('file.ts', 'comment-1', 'updated')).resolves.toBeUndefined();
    });
  });

  describe('listFiles', () => {
    it('fetches files with default path', async () => {
      const mockResponse: ListFilesResponse = {
        path: '.',
        files: [
          { name: 'file1.ts', path: 'file1.ts', isDir: false, modTime: '2024-01-01' },
          { name: 'dir1', path: 'dir1', isDir: true, modTime: '2024-01-01' },
        ],
      };

      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify(mockResponse), { status: 200 })
      );

      const result = await listFiles();
      expect(result).toEqual(mockResponse);
      expect(fetch).toHaveBeenCalledWith('/api/files?path=.');
    });

    it('fetches files with custom path', async () => {
      const mockResponse: ListFilesResponse = {
        path: 'src/components',
        files: [{ name: 'Button.tsx', path: 'src/components/Button.tsx', isDir: false, modTime: '2024-01-01' }],
      };

      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify(mockResponse), { status: 200 })
      );

      const result = await listFiles('src/components');
      expect(result).toEqual(mockResponse);
      expect(fetch).toHaveBeenCalledWith('/api/files?path=src%2Fcomponents');
    });

    it('encodes special characters in path', async () => {
      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify({ path: '', files: [] }), { status: 200 })
      );

      await listFiles('path with spaces & special chars');
      expect(fetch).toHaveBeenCalledWith('/api/files?path=path%20with%20spaces%20%26%20special%20chars');
    });
  });

  describe('getFileContent', () => {
    it('fetches file content', async () => {
      const mockResponse: FileResponse = {
        file: {
          path: 'src/App.tsx',
          name: 'App.tsx',
          content: 'export default function App() {}',
          size: 32,
          modTime: '2024-01-01',
          mimeType: 'text/typescript',
        },
      };

      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify(mockResponse), { status: 200 })
      );

      const result = await getFileContent('src/App.tsx');
      expect(result).toEqual(mockResponse);
      expect(fetch).toHaveBeenCalledWith('/api/files/src%2FApp.tsx');
    });
  });

  describe('getRelatedTests', () => {
    it('fetches related tests', async () => {
      const mockResponse: TestsResponse = {
        sourceFile: 'src/utils.ts',
        tests: [
          {
            functionName: 'helper',
            testFile: 'src/utils.test.ts',
            testName: 'should work',
            content: 'test("should work", () => {})',
            lineRange: { start: 1, end: 3 },
            coveredLines: { start: 5, end: 10 },
          },
        ],
      };

      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify(mockResponse), { status: 200 })
      );

      const result = await getRelatedTests('src/utils.ts');
      expect(result).toEqual(mockResponse);
      expect(fetch).toHaveBeenCalledWith('/api/files/src%2Futils.ts/tests');
    });
  });

  describe('getSuggestions', () => {
    it('fetches test suggestions', async () => {
      const mockResponse: SuggestionsResponse = {
        sourceFile: 'src/utils.ts',
        suggestions: [
          {
            sourceFile: 'src/utils.ts',
            targetLines: { start: 5, end: 10 },
            reason: 'Function not tested',
            suggestedName: 'should handle edge cases',
            testSkeleton: 'test("should handle edge cases", () => {})',
            priority: 'high',
          },
        ],
      };

      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify(mockResponse), { status: 200 })
      );

      const result = await getSuggestions('src/utils.ts');
      expect(result).toEqual(mockResponse);
    });
  });

  describe('getSourceReferences', () => {
    it('fetches source references', async () => {
      const mockResponse: TestFileResponse = {
        testFile: 'src/utils.test.ts',
        sources: [
          {
            sourceFile: 'src/utils.ts',
            functionName: 'helper',
            coveredLines: { start: 5, end: 10 },
            testName: 'should work',
            lineRange: { start: 1, end: 3 },
          },
        ],
      };

      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify(mockResponse), { status: 200 })
      );

      const result = await getSourceReferences('src/utils.test.ts');
      expect(result).toEqual(mockResponse);
    });
  });

  describe('comment API', () => {
    const mockComment = {
      id: 'comment-1',
      line: 10,
      content: 'Test comment',
      author: 'user',
      createdAt: '2024-01-01',
      updatedAt: '2024-01-01',
      resolved: false,
    };

    describe('getComments', () => {
      it('fetches comments', async () => {
        const mockResponse: CommentsResponse = {
          sourceFile: 'src/App.tsx',
          comments: [mockComment],
        };

        vi.spyOn(globalThis, 'fetch').mockResolvedValue(
          new Response(JSON.stringify(mockResponse), { status: 200 })
        );

        const result = await getComments('src/App.tsx');
        expect(result).toEqual(mockResponse);
      });
    });

    describe('createComment', () => {
      it('creates a comment', async () => {
        const mockResponse: CommentResponse = { comment: mockComment };

        vi.spyOn(globalThis, 'fetch').mockResolvedValue(
          new Response(JSON.stringify(mockResponse), { status: 201 })
        );

        const result = await createComment('src/App.tsx', {
          line: 10,
          content: 'Test comment',
          contextLines: { start: 7, end: 13 },
        });

        expect(result).toEqual(mockResponse);
        expect(fetch).toHaveBeenCalledWith('/api/files/src%2FApp.tsx/comments', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            line: 10,
            content: 'Test comment',
            contextLines: { start: 7, end: 13 },
          }),
        });
      });
    });

    describe('updateComment', () => {
      it('updates a comment', async () => {
        vi.spyOn(globalThis, 'fetch').mockResolvedValue(
          new Response(null, { status: 204 })
        );

        await expect(updateComment('file.ts', 'comment-1', 'updated')).resolves.toBeUndefined();
      });

      it('treats empty 200 responses as success without Content-Length', async () => {
        vi.spyOn(globalThis, 'fetch').mockResolvedValue(
          new Response(null, { status: 200 })
        );

        await expect(updateComment('file.ts', 'comment-1', 'updated')).resolves.toBeUndefined();
      });
    });

    describe('deleteComment', () => {
      it('deletes a comment', async () => {
        vi.spyOn(globalThis, 'fetch').mockResolvedValue(
          new Response(null, { status: 204 })
        );

        await expect(deleteComment('file.ts', 'comment-1')).resolves.toBeUndefined();
        expect(fetch).toHaveBeenCalledWith('/api/files/file.ts/comments/comment-1', {
          method: 'DELETE',
        });
      });
    });

    describe('toggleCommentResolved', () => {
      it('toggles resolved status', async () => {
        vi.spyOn(globalThis, 'fetch').mockResolvedValue(
          new Response(null, { status: 200, headers: { 'content-length': '0' } })
        );

        await expect(toggleCommentResolved('file.ts', 'comment-1')).resolves.toBeUndefined();
        expect(fetch).toHaveBeenCalledWith('/api/files/file.ts/comments/comment-1/resolved', {
          method: 'PATCH',
        });
      });
    });
  });

  describe('exportContext', () => {
    it('exports context with default options', async () => {
      const mockResponse: ExportContextResponse = {
        sourceFile: 'src/App.tsx',
        codeContext: [],
        formatted: 'formatted output',
      };

      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify(mockResponse), { status: 200 })
      );

      const result = await exportContext('src/App.tsx', {
        includeTests: true,
        includeSuggestions: false,
        contextLines: 5,
      });
      expect(result).toEqual(mockResponse);
    });

    it('exports context with custom options', async () => {
      const mockResponse: ExportContextResponse = {
        sourceFile: 'src/App.tsx',
        codeContext: [],
        tests: [],
        suggestions: [],
        formatted: 'formatted output',
      };

      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify(mockResponse), { status: 200 })
      );

      await exportContext('src/App.tsx', {
        includeTests: true,
        includeSuggestions: true,
        contextLines: 10,
      });

      expect(fetch).toHaveBeenCalledWith('/api/files/src%2FApp.tsx/export', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          includeTests: true,
          includeSuggestions: true,
          contextLines: 10,
        }),
      });
    });
  });

  describe('getProjectOverview', () => {
    it('fetches project overview', async () => {
      const mockResponse: OverviewResponse = {
        totalTests: 100,
        totalFunctions: 50,
        totalSourceFiles: 25,
        totalTestFiles: 20,
        functions: [],
        testsBySourceFile: {},
      };

      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify(mockResponse), { status: 200 })
      );

      const result = await getProjectOverview();
      expect(result).toEqual(mockResponse);
      expect(fetch).toHaveBeenCalledWith('/api/overview');
    });
  });

  describe('getMetadataIssues', () => {
    it('fetches metadata issues', async () => {
      const mockResponse: MetadataIssuesResponse = {
        issues: [
          {
            sourceFile: 'src/old.ts',
            sourceValid: false,
            sourceIsAbsolute: false,
            sourceMessage: 'File not found',
            suggestionsCount: 0,
            commentsCount: 0,
            invalidTestIssues: [],
          },
        ],
      };

      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify(mockResponse), { status: 200 })
      );

      const result = await getMetadataIssues();
      expect(result).toEqual(mockResponse);
      expect(fetch).toHaveBeenCalledWith('/api/metadata/issues');
    });
  });

  describe('updateSourcePath', () => {
    it('updates source path', async () => {
      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(null, { status: 204 })
      );

      await expect(updateSourcePath('old/path.ts', 'new/path.ts')).resolves.toBeUndefined();
      expect(fetch).toHaveBeenCalledWith('/api/metadata/source-path', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ oldPath: 'old/path.ts', newPath: 'new/path.ts' }),
      });
    });
  });

  describe('updateTestPath', () => {
    it('updates test path', async () => {
      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(null, { status: 204 })
      );

      await expect(updateTestPath('src.ts', 'old.test.ts', 'testName', 'new.test.ts')).resolves.toBeUndefined();
      expect(fetch).toHaveBeenCalledWith('/api/metadata/test-path', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          sourceFile: 'src.ts',
          testFile: 'old.test.ts',
          testName: 'testName',
          newTestFile: 'new.test.ts',
        }),
      });
    });
  });

  describe('deleteSourcePath', () => {
    it('deletes source path', async () => {
      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(null, { status: 204 })
      );

      await expect(deleteSourcePath('path/to/file.ts')).resolves.toBeUndefined();
      expect(fetch).toHaveBeenCalledWith('/api/metadata/source-path', {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ path: 'path/to/file.ts' }),
      });
    });
  });

  describe('deleteTestPath', () => {
    it('deletes test path', async () => {
      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(null, { status: 204 })
      );

      await expect(deleteTestPath('src.ts', 'test.ts', 'testName')).resolves.toBeUndefined();
      expect(fetch).toHaveBeenCalledWith('/api/metadata/test-path', {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          sourceFile: 'src.ts',
          testFile: 'test.ts',
          testName: 'testName',
        }),
      });
    });
  });

  describe('search', () => {
    it('searches with query', async () => {
      const mockResponse: SearchResponse = {
        query: 'test',
        results: [
          {
            type: 'file',
            title: 'test.ts',
            subtitle: 'Test file',
            path: 'src/test.ts',
            line: 1,
            relevance: 1.0,
            matchedText: 'test',
          },
        ],
      };

      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify(mockResponse), { status: 200 })
      );

      const result = await search('test');
      expect(result).toEqual(mockResponse);
      expect(fetch).toHaveBeenCalledWith('/api/search?q=test');
    });

    it('encodes special characters in search query', async () => {
      vi.spyOn(globalThis, 'fetch').mockResolvedValue(
        new Response(JSON.stringify({ query: '', results: [] }), { status: 200 })
      );

      await search('query with spaces & special chars');
      expect(fetch).toHaveBeenCalledWith('/api/search?q=query%20with%20spaces%20%26%20special%20chars');
    });
  });
});
