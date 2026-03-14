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
  comments?: Comment[];
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

export interface SourceReference {
  sourceFile: string;
  functionName: string;
  coveredLines: LineRange;
  testName: string;
  comment?: string;
  lineRange: LineRange;
  inputLines?: LineRange;
  outputLines?: LineRange;
}

export interface TestFileResponse {
  testFile: string;
  sources: SourceReference[];
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
  contextLines: number;
}

export interface ExportContextResponse {
  sourceFile: string;
  codeContext: CodeContextBlock[];
  tests?: TestDetail[];
  formatted: string;
}

// ==================== OVERVIEW TYPES ====================

export interface FunctionSummary {
  functionName: string;
  sourceFile: string;
  testCount: number;
  tests: TestDetail[];
}

export interface OverviewResponse {
  totalTests: number;
  totalFunctions: number;
  totalSourceFiles: number;
  totalTestFiles: number;
  functions: FunctionSummary[];
  testsBySourceFile: Record<string, TestDetail[]>;
}

export interface MetadataTestIssue {
  testFile: string;
  testName: string;
  isAbsolute: boolean;
  message: string;
}

export interface MetadataIssue {
  sourceFile: string;
  sourceValid: boolean;
  sourceIsAbsolute: boolean;
  sourceMessage?: string;
  commentsCount: number;
  invalidTestIssues: MetadataTestIssue[];
}

export interface MetadataIssuesResponse {
  issues: MetadataIssue[];
}

// ==================== SEARCH TYPES ====================

export type SearchResultType = 'file' | 'function' | 'test';

export interface SearchResult {
  type: SearchResultType;
  title: string;
  subtitle: string;
  path: string;
  line: number;
  relevance: number;
  matchedText: string;
}

export interface SearchResponse {
  query: string;
  results: SearchResult[];
}
