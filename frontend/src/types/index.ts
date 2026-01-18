export interface FileEntry {
  name: string;
  path: string;
  isDir: boolean;
  size?: number;
  modTime: string;
}

export interface LineRange {
  start: number;
  end: number;
}

export interface TestReference {
  testFile: string;
  testName: string;
  comment?: string;
  lineRange: LineRange;
  coveredLines: LineRange;
  inputLines?: LineRange;
  outputLines?: LineRange;
}

export interface FileMetadata {
  tests?: TestReference[];
}

export interface FileContent {
  path: string;
  name: string;
  content: string;
  size: number;
  modTime: string;
  mimeType: string;
  metadata?: FileMetadata;
}

export interface TestDetail {
  testFile: string;
  testName: string;
  comment?: string;
  content: string;
  lineRange: LineRange;
  coveredLines: LineRange;
  inputData?: string;
  inputLines?: LineRange;
  expectedOutput?: string;
  outputLines?: LineRange;
}

export interface ListFilesResponse {
  path: string;
  files: FileEntry[];
}

export interface FileResponse {
  file: FileContent;
}

export interface TestsResponse {
  sourceFile: string;
  tests: TestDetail[];
}

export interface MindMapNode {
  id: string;
  label: string;
  children?: MindMapNode[];
}
