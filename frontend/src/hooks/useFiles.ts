import { useState, useEffect } from 'react';
import { listFiles } from '../api/client';
import type { FileEntry } from '../types';

export function useFiles(path: string) {
  const [files, setFiles] = useState<FileEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    setError(null);

    listFiles(path)
      .then(res => setFiles(res.files))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, [path]);

  return { files, loading, error };
}
