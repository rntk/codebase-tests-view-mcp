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

export interface Comment {
  id: string;
  line: number;
  content: string;
  author?: string;
  createdAt: string;
  updatedAt: string;
  resolved: boolean;
  contextLines?: LineRange;
}

export interface TestReference {
  functionName: string;
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
  suggestions?: TestSuggestion[];
  comments?: Comment[];
}

export interface TestSuggestion {
  sourceFile: string;
  functionName?: string;
  targetLines: LineRange;
  reason: string;
  suggestedName: string;
  testSkeleton: string;
  priority: 'high' | 'medium' | 'low';
}

export interface SuggestionsResponse {
  sourceFile: string;
  suggestions: TestSuggestion[];
}

export interface FileContent {
  path: string;
  name: string;
  content: string;
  size: number;
  modTime: string;
  mimeType: string;
  metadata?: FileMetadata;
  coverageDepth?: CoverageDepth;
}

export interface CoverageDepth {
  [lineNumber: number]: string[]; // line number -> test names
}

export interface TestDetail {
  functionName: string;
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
  edgeLabel?: string;
}

export type LayoutType = 'horizontal' | 'radial' | 'clustered';

export interface MindMapTransform {
  x: number;
  y: number;
  scale: number;
}

// Comment-related types
export interface CommentRequest {
  line: number;
  content: string;
  contextLines?: LineRange;
}

export interface CommentResponse {
  comment: Comment;
}

export interface CommentsResponse {
  sourceFile: string;
  comments: Comment[];
}

export interface CodeContextBlock {
  lineRange: LineRange;
  code: string;
  comments: Comment[];
}

export interface ExportContextRequest {
  includeTests: boolean;
  includeSuggestions: boolean;
  contextLines: number;
}

export interface ExportContextResponse {
  sourceFile: string;
  codeContext: CodeContextBlock[];
  tests?: TestDetail[];
  suggestions?: TestSuggestion[];
  formatted: string;
}
