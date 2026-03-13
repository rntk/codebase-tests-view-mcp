import { useState, useEffect, useRef, useMemo } from 'react';
import { getRelatedTests } from '../api/client';
import type { TestDetail, FileEntry } from '../types';

/**
 * Limits the concurrency of async operations.
 * @param items - Array of items to process
 * @param concurrency - Maximum number of concurrent operations
 * @param processor - Async function to process each item
 */
async function limitConcurrency<T, R>(
  items: T[],
  concurrency: number,
  processor: (item: T) => Promise<R>
): Promise<R[]> {
  const results: R[] = [];
  const executing: Promise<void>[] = [];

  for (const item of items) {
    const promise = processor(item).then((result) => {
      results.push(result);
      executing.splice(executing.indexOf(promise), 1);
    });
    executing.push(promise);

    if (executing.length >= concurrency) {
      await Promise.race(executing);
    }
  }

  await Promise.all(executing);
  return results;
}

/**
 * Hook to load tests for all files in a directory.
 * Aggregates tests from all source files in the given path.
 * Limits concurrent API requests to avoid overwhelming the backend.
 */
export function useDirectoryTests(path: string | null, files: FileEntry[]) {
  const [tests, setTests] = useState<TestDetail[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Use a ref to track the current request to avoid stale updates
  const requestIdRef = useRef(0);

  // Memoize file paths to avoid unnecessary re-fetches when array reference changes
  const filePaths = useMemo(
    () => files.filter((f) => !f.isDir).map((f) => f.path).sort().join(','),
    [files]
  );

  useEffect(() => {
    if (!path || !filePaths) {
      setTests([]);
      return;
    }

    const sourceFiles = files.filter((f) => !f.isDir).map((f) => f.path);
    if (sourceFiles.length === 0) {
      setTests([]);
      return;
    }

    // Increment request ID to track this specific request
    const currentRequestId = ++requestIdRef.current;
    setLoading(true);
    setError(null);

    // Fetch tests with limited concurrency (max 5 parallel requests)
    limitConcurrency(sourceFiles, 5, (filePath) =>
      getRelatedTests(filePath).catch(() => ({ tests: [] }))
    )
      .then((results) => {
        // Only update state if this is still the current request
        if (currentRequestId !== requestIdRef.current) {
          return;
        }

        const allTests: TestDetail[] = [];
        results.forEach((res) => {
          if (res.tests) {
            allTests.push(...res.tests);
          }
        });
        setTests(allTests);
      })
      .catch((err) => {
        if (currentRequestId !== requestIdRef.current) {
          return;
        }
        setError(err.message);
      })
      .finally(() => {
        if (currentRequestId === requestIdRef.current) {
          setLoading(false);
        }
      });
  }, [path, filePaths, files]);

  return { tests, loading, error };
}