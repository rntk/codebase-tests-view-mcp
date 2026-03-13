import type {
  ListFilesResponse,
  FileResponse,
  TestsResponse,
  SuggestionsResponse,
  TestFileResponse,
  CommentsResponse,
  CommentResponse,
  CommentRequest,
  ExportContextRequest,
  ExportContextResponse,
  OverviewResponse,
  MetadataIssuesResponse
} from '../types';

const API_BASE = '/api';

export async function listFiles(path: string = '.'): Promise<ListFilesResponse> {
  const response = await fetch(`${API_BASE}/files?path=${encodeURIComponent(path)}`);
  if (!response.ok) {
    throw new Error(`Failed to list files: ${response.statusText}`);
  }
  return response.json();
}

export async function getFileContent(path: string): Promise<FileResponse> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}`);
  if (!response.ok) {
    throw new Error(`Failed to get file: ${response.statusText}`);
  }
  return response.json();
}

export async function getRelatedTests(path: string): Promise<TestsResponse> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/tests`);
  if (!response.ok) {
    throw new Error(`Failed to get tests: ${response.statusText}`);
  }
  return response.json();
}

export async function getSuggestions(path: string): Promise<SuggestionsResponse> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/suggestions`);
  if (!response.ok) {
    throw new Error(`Failed to get suggestions: ${response.statusText}`);
  }
  return response.json();
}

export async function getSourceReferences(path: string): Promise<TestFileResponse> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/sources`);
  if (!response.ok) {
    throw new Error(`Failed to get source references: ${response.statusText}`);
  }
  return response.json();
}

// ==================== COMMENT API ====================

export async function getComments(path: string): Promise<CommentsResponse> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/comments`);
  if (!response.ok) {
    throw new Error(`Failed to get comments: ${response.statusText}`);
  }
  return response.json();
}

export async function createComment(path: string, request: CommentRequest): Promise<CommentResponse> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/comments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  });
  if (!response.ok) {
    throw new Error(`Failed to create comment: ${response.statusText}`);
  }
  return response.json();
}

export async function updateComment(path: string, commentId: string, content: string): Promise<void> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/comments/${encodeURIComponent(commentId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ content }),
  });
  if (!response.ok) {
    throw new Error(`Failed to update comment: ${response.statusText}`);
  }
}

export async function deleteComment(path: string, commentId: string): Promise<void> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/comments/${encodeURIComponent(commentId)}`, {
    method: 'DELETE',
  });
  if (!response.ok) {
    throw new Error(`Failed to delete comment: ${response.statusText}`);
  }
}

export async function toggleCommentResolved(path: string, commentId: string): Promise<void> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/comments/${encodeURIComponent(commentId)}/resolved`, {
    method: 'PATCH',
  });
  if (!response.ok) {
    throw new Error(`Failed to toggle comment resolved: ${response.statusText}`);
  }
}

// ==================== EXPORT API ====================

export async function exportContext(
  path: string,
  request: ExportContextRequest
): Promise<ExportContextResponse> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/export`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  });
  if (!response.ok) {
    throw new Error(`Failed to export context: ${response.statusText}`);
  }
  return response.json();
}

// ==================== OVERVIEW API ====================

export async function getProjectOverview(): Promise<OverviewResponse> {
  const response = await fetch(`${API_BASE}/overview`);
  if (!response.ok) {
    throw new Error(`Failed to get project overview: ${response.statusText}`);
  }
  return response.json();
}

export async function getMetadataIssues(): Promise<MetadataIssuesResponse> {
  const response = await fetch(`${API_BASE}/metadata/issues`);
  if (!response.ok) {
    throw new Error(`Failed to get metadata issues: ${response.statusText}`);
  }
  return response.json();
}

export async function updateSourcePath(oldPath: string, newPath: string): Promise<void> {
  const response = await fetch(`${API_BASE}/metadata/source-path`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ oldPath, newPath }),
  });
  if (!response.ok) {
    throw new Error(await response.text() || `Failed to update source path: ${response.statusText}`);
  }
}

export async function updateTestPath(sourceFile: string, testFile: string, testName: string, newTestFile: string): Promise<void> {
  const response = await fetch(`${API_BASE}/metadata/test-path`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ sourceFile, testFile, testName, newTestFile }),
  });
  if (!response.ok) {
    throw new Error(await response.text() || `Failed to update test path: ${response.statusText}`);
  }
}

export async function deleteSourcePath(path: string): Promise<void> {
  const response = await fetch(`${API_BASE}/metadata/source-path`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ path }),
  });
  if (!response.ok) {
    throw new Error(await response.text() || `Failed to delete source path: ${response.statusText}`);
  }
}

export async function deleteTestPath(sourceFile: string, testFile: string, testName: string): Promise<void> {
  const response = await fetch(`${API_BASE}/metadata/test-path`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ sourceFile, testFile, testName }),
  });
  if (!response.ok) {
    throw new Error(await response.text() || `Failed to delete test path: ${response.statusText}`);
  }
}
