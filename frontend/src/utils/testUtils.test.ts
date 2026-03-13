import { filterItemsByLine } from './testUtils';
import type { LineRange } from '../types';

// Test item interface that matches the generic type T
interface TestItem {
  id: string;
  coveredLines?: LineRange;
}

describe('filterItemsByLine', () => {
  const createItem = (id: string, start: number, end: number): TestItem => ({
    id,
    coveredLines: { start, end },
  });

  const getCoveredLines = (item: TestItem): LineRange | undefined => item.coveredLines;

  describe('when selectedLine is null or undefined', () => {
    it('returns all items when selectedLine is null', () => {
      const items = [
        createItem('item1', 1, 10),
        createItem('item2', 15, 25),
        createItem('item3', 30, 40),
      ];

      const result = filterItemsByLine(items, null, getCoveredLines);

      expect(result).toEqual(items);
      expect(result.length).toBe(3);
    });

    it('returns all items when selectedLine is undefined', () => {
      const items = [
        createItem('item1', 1, 10),
        createItem('item2', 15, 25),
      ];

      const result = filterItemsByLine(items, undefined, getCoveredLines);

      expect(result).toEqual(items);
      expect(result.length).toBe(2);
    });
  });

  describe('when no items cover the selected line', () => {
    it('returns empty array', () => {
      const items = [
        createItem('item1', 1, 10),
        createItem('item2', 20, 30),
      ];

      const result = filterItemsByLine(items, 15, getCoveredLines);

      expect(result).toEqual([]);
      expect(result.length).toBe(0);
    });
  });

  describe('when items cover the selected line', () => {
    it('returns items that directly cover the selected line', () => {
      const items = [
        createItem('item1', 1, 10),
        createItem('item2', 5, 15),
        createItem('item3', 20, 30),
      ];

      const result = filterItemsByLine(items, 7, getCoveredLines);

      expect(result.length).toBe(2);
      expect(result.map(i => i.id)).toEqual(['item1', 'item2']);
    });

    it('calculates union range and returns overlapping items', () => {
      // item1 covers 1-10, item2 covers 5-15
      // Union range is 1-15
      // item3 covers 12-18, which overlaps with union range
      const items = [
        createItem('item1', 1, 10),
        createItem('item2', 5, 15),
        createItem('item3', 12, 18),
        createItem('item4', 25, 35),
      ];

      // Select line 7, which is covered by item1 and item2
      // Union range is 1-15
      // item3 (12-18) overlaps with union range
      const result = filterItemsByLine(items, 7, getCoveredLines);

      expect(result.length).toBe(3);
      expect(result.map(i => i.id)).toEqual(['item1', 'item2', 'item3']);
    });

    it('handles single item covering the line', () => {
      const items = [
        createItem('item1', 1, 10),
        createItem('item2', 20, 30),
      ];

      const result = filterItemsByLine(items, 5, getCoveredLines);

      expect(result.length).toBe(1);
      expect(result[0].id).toBe('item1');
    });

    it('handles item that exactly matches the selected line', () => {
      const items = [
        createItem('item1', 5, 5),
        createItem('item2', 10, 20),
      ];

      const result = filterItemsByLine(items, 5, getCoveredLines);

      expect(result.length).toBe(1);
      expect(result[0].id).toBe('item1');
    });

    it('handles item starting at the selected line', () => {
      const items = [
        createItem('item1', 10, 20),
        createItem('item2', 1, 5),
      ];

      const result = filterItemsByLine(items, 10, getCoveredLines);

      expect(result.length).toBe(1);
      expect(result[0].id).toBe('item1');
    });

    it('handles item ending at the selected line', () => {
      const items = [
        createItem('item1', 1, 10),
        createItem('item2', 15, 25),
      ];

      const result = filterItemsByLine(items, 10, getCoveredLines);

      expect(result.length).toBe(1);
      expect(result[0].id).toBe('item1');
    });
  });

  describe('edge cases', () => {
    it('handles empty items array', () => {
      const items: TestItem[] = [];

      const result = filterItemsByLine(items, 5, getCoveredLines);

      expect(result).toEqual([]);
      expect(result.length).toBe(0);
    });

    it('handles items with undefined coveredLines', () => {
      const items: TestItem[] = [
        { id: 'item1' },
        createItem('item2', 1, 10),
        { id: 'item3' },
      ];

      const result = filterItemsByLine(items, 5, getCoveredLines);

      expect(result.length).toBe(1);
      expect(result[0].id).toBe('item2');
    });

    it('handles all items with undefined coveredLines', () => {
      const items: TestItem[] = [
        { id: 'item1' },
        { id: 'item2' },
        { id: 'item3' },
      ];

      const result = filterItemsByLine(items, 5, getCoveredLines);

      expect(result).toEqual([]);
      expect(result.length).toBe(0);
    });

    it('handles selected line at boundary of union range', () => {
      const items = [
        createItem('item1', 1, 5),
        createItem('item2', 5, 10),
        createItem('item3', 10, 15),
      ];

      // Select line 5, covered by item1 and item2
      // Union range is 1-10
      // item3 (10-15) overlaps at boundary
      const result = filterItemsByLine(items, 5, getCoveredLines);

      expect(result.length).toBe(3);
      expect(result.map(i => i.id)).toEqual(['item1', 'item2', 'item3']);
    });

    it('handles non-contiguous ranges', () => {
      const items = [
        createItem('item1', 1, 5),
        createItem('item2', 10, 15),
        createItem('item3', 20, 25),
      ];

      // Select line 3, only item1 covers it
      const result = filterItemsByLine(items, 3, getCoveredLines);

      expect(result.length).toBe(1);
      expect(result[0].id).toBe('item1');
    });

    it('handles large line numbers', () => {
      const items = [
        createItem('item1', 1000, 2000),
        createItem('item2', 1500, 2500),
        createItem('item3', 3000, 4000),
      ];

      const result = filterItemsByLine(items, 1750, getCoveredLines);

      expect(result.length).toBe(2);
      expect(result.map(i => i.id)).toEqual(['item1', 'item2']);
    });

    it('handles items with same range', () => {
      const items = [
        createItem('item1', 1, 10),
        createItem('item2', 1, 10),
        createItem('item3', 1, 10),
      ];

      const result = filterItemsByLine(items, 5, getCoveredLines);

      expect(result.length).toBe(3);
      expect(result.map(i => i.id)).toEqual(['item1', 'item2', 'item3']);
    });

    it('handles single line ranges', () => {
      const items = [
        createItem('item1', 1, 1),
        createItem('item2', 2, 2),
        createItem('item3', 3, 3),
      ];

      const result = filterItemsByLine(items, 2, getCoveredLines);

      expect(result.length).toBe(1);
      expect(result[0].id).toBe('item2');
    });

    it('handles overlapping ranges with different sizes', () => {
      const items = [
        createItem('item1', 1, 100),
        createItem('item2', 50, 60),
        createItem('item3', 55, 57),
      ];

      const result = filterItemsByLine(items, 56, getCoveredLines);

      // All three items cover line 56
      expect(result.length).toBe(3);
      expect(result.map(i => i.id)).toEqual(['item1', 'item2', 'item3']);
    });
  });

  describe('union range calculation', () => {
    it('correctly calculates min start and max end', () => {
      const items = [
        createItem('item1', 10, 20),
        createItem('item2', 5, 15),
        createItem('item3', 12, 25),
      ];

      // Select line 13, covered by all three
      // Union range should be 5-25
      const result = filterItemsByLine(items, 13, getCoveredLines);

      expect(result.length).toBe(3);
    });

    it('includes items that partially overlap with union range', () => {
      const items = [
        createItem('item1', 10, 20),
        createItem('item2', 15, 25),
        createItem('item3', 22, 30),
        createItem('item4', 35, 45),
      ];

      // Select line 17, covered by item1 and item2
      // Union range is 10-25
      // item3 (22-30) overlaps with union range
      // item4 (35-45) does not overlap
      const result = filterItemsByLine(items, 17, getCoveredLines);

      expect(result.length).toBe(3);
      expect(result.map(i => i.id)).toEqual(['item1', 'item2', 'item3']);
    });
  });

  describe('generic type support', () => {
    it('works with custom object types', () => {
      interface CustomItem {
        name: string;
        range?: LineRange;
        metadata: string;
      }

      const items: CustomItem[] = [
        { name: 'test1', range: { start: 1, end: 10 }, metadata: 'meta1' },
        { name: 'test2', range: { start: 15, end: 25 }, metadata: 'meta2' },
      ];

      const result = filterItemsByLine(
        items,
        5,
        (item: CustomItem) => item.range
      );

      expect(result.length).toBe(1);
      expect(result[0].name).toBe('test1');
      expect(result[0].metadata).toBe('meta1');
    });
  });
});
