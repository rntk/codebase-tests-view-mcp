import { LineRange } from '../types';

/**
 * Filters items (tests) based on a selected line.
 * Instead of just returning items that strictly cover the line,
 * it calculates the union range of all items covering the line (the "common range"),
 * and returns all items that overlap with this common range.
 * This effectively shows tests for the "current function" if the tests cover the function.
 */
export function filterItemsByLine<T>(
    items: T[],
    selectedLine: number | null | undefined,
    getCoveredLines: (item: T) => LineRange | undefined
): T[] {
    if (selectedLine === null || selectedLine === undefined) {
        return items;
    }

    // 1. Find items that directly cover the selected line
    const directItems = items.filter(item => {
        const range = getCoveredLines(item);
        return range && selectedLine >= range.start && selectedLine <= range.end;
    });

    if (directItems.length === 0) {
        return [];
    }

    // 2. Calculate the union range of all direct items
    let minStart = Number.MAX_SAFE_INTEGER;
    let maxEnd = Number.MIN_SAFE_INTEGER;

    directItems.forEach(item => {
        const range = getCoveredLines(item);
        if (range) {
            if (range.start < minStart) minStart = range.start;
            if (range.end > maxEnd) maxEnd = range.end;
        }
    });

    // 3. Find all items that overlap with the union range
    return items.filter(item => {
        const range = getCoveredLines(item);
        if (!range) return false;

        // Check for overlap: max(start1, start2) <= min(end1, end2)
        return Math.max(range.start, minStart) <= Math.min(range.end, maxEnd);
    });
}
