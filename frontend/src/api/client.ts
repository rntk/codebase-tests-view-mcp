import type {
  ListFilesResponse,
  FileResponse,
  TestsResponse,
  TestFileResponse,
  CommentsResponse,
  CommentResponse,
  CommentRequest,
  ExportContextRequest,
  ExportContextResponse,
  OverviewResponse,
  MetadataIssuesResponse,
  SearchResponse
} from '../types';

const API_BASE = '/api';

interface APIError {
  error: string;
  code: string;
  details?: Record<string, unknown>;
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let errorMessage = response.statusText;
    try {
      const errorText = await response.text();
      if (errorText.trim() !== '') {
        const errorData: APIError = JSON.parse(errorText);
        errorMessage = errorData.error || errorMessage;
      }
    } catch (e) {
      console.error('Failed to parse error response:', e);
    }
    throw new Error(errorMessage);
  }
  
  if (response.status === 204 || response.headers.get('content-length') === '0') {
    return undefined as T;
  }

  const responseText = await response.text();
  if (responseText.trim() === '') {
    return undefined as T;
  }

  return JSON.parse(responseText) as T;
}

export async function listFiles(path: string = '.'): Promise<ListFilesResponse> {
  const response = await fetch(`${API_BASE}/files?path=${encodeURIComponent(path)}`);
  return handleResponse<ListFilesResponse>(response);
}

export async function getFileContent(path: string): Promise<FileResponse> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}`);
  return handleResponse<FileResponse>(response);
}

export async function getRelatedTests(path: string): Promise<TestsResponse> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/tests`);
  return handleResponse<TestsResponse>(response);
}

export async function getSourceReferences(path: string): Promise<TestFileResponse> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/sources`);
  return handleResponse<TestFileResponse>(response);
}

// ==================== COMMENT API ====================

export async function getComments(path: string): Promise<CommentsResponse> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/comments`);
  return handleResponse<CommentsResponse>(response);
}

export async function createComment(path: string, request: CommentRequest): Promise<CommentResponse> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/comments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  });
  return handleResponse<CommentResponse>(response);
}

export async function updateComment(path: string, commentId: string, content: string): Promise<void> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/comments/${encodeURIComponent(commentId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ content }),
  });
  // 204 No Content responses don't have a body
  if (response.status === 204) {
    return;
  }
  return handleResponse<void>(response);
}

export async function deleteComment(path: string, commentId: string): Promise<void> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/comments/${encodeURIComponent(commentId)}`, {
    method: 'DELETE',
  });
  // 204 No Content responses don't have a body
  if (response.status === 204) {
    return;
  }
  return handleResponse<void>(response);
}

export async function toggleCommentResolved(path: string, commentId: string): Promise<void> {
  const response = await fetch(`${API_BASE}/files/${encodeURIComponent(path)}/comments/${encodeURIComponent(commentId)}/resolved`, {
    method: 'PATCH',
  });
  // 200 OK with empty body
  if (response.status === 200 && response.headers.get('content-length') === '0') {
    return;
  }
  return handleResponse<void>(response);
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
  return handleResponse<ExportContextResponse>(response);
}

// ==================== OVERVIEW API ====================

export async function getProjectOverview(): Promise<OverviewResponse> {
  const response = await fetch(`${API_BASE}/overview`);
  return handleResponse<OverviewResponse>(response);
}

export async function getMetadataIssues(): Promise<MetadataIssuesResponse> {
  const response = await fetch(`${API_BASE}/metadata/issues`);
  return handleResponse<MetadataIssuesResponse>(response);
}

export async function updateSourcePath(oldPath: string, newPath: string): Promise<void> {
  const response = await fetch(`${API_BASE}/metadata/source-path`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ oldPath, newPath }),
  });
  // 204 No Content responses don't have a body
  if (response.status === 204) {
    return;
  }
  return handleResponse<void>(response);
}

export async function updateTestPath(sourceFile: string, testFile: string, testName: string, newTestFile: string): Promise<void> {
  const response = await fetch(`${API_BASE}/metadata/test-path`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ sourceFile, testFile, testName, newTestFile }),
  });
  // 204 No Content responses don't have a body
  if (response.status === 204) {
    return;
  }
  return handleResponse<void>(response);
}

export async function deleteSourcePath(path: string): Promise<void> {
  const response = await fetch(`${API_BASE}/metadata/source-path`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ path }),
  });
  // 204 No Content responses don't have a body
  if (response.status === 204) {
    return;
  }
  return handleResponse<void>(response);
}

export async function deleteTestPath(sourceFile: string, testFile: string, testName: string): Promise<void> {
  const response = await fetch(`${API_BASE}/metadata/test-path`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ sourceFile, testFile, testName }),
  });
  // 204 No Content responses don't have a body
  if (response.status === 204) {
    return;
  }
  return handleResponse<void>(response);
}

// ==================== SEARCH API ====================

export async function search(query: string): Promise<SearchResponse> {
  const response = await fetch(`${API_BASE}/search?q=${encodeURIComponent(query)}`);
  return handleResponse<SearchResponse>(response);
}
