import type { ListFilesResponse, FileResponse, TestsResponse } from '../types';

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
