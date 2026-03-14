import { describe, expect, it } from 'vitest';
import { groupTestsByFunction } from './testUtils';
import type { TestDetail } from '../types';

describe('groupTestsByFunction', () => {
  const createTest = (
    functionName: string,
    testFile: string,
    testName: string
  ): TestDetail => ({
    functionName,
    testFile,
    testName,
    content: `test("${testName}", () => {})`,
    lineRange: { start: 1, end: 3 },
    coveredLines: { start: 1, end: 10 },
  });

  it('returns empty map for empty array', () => {
    const result = groupTestsByFunction([]);
    expect(result.size).toBe(0);
  });

  it('groups single test by function name', () => {
    const tests: TestDetail[] = [
      createTest('helper', 'utils.test.ts', 'should work'),
    ];

    const result = groupTestsByFunction(tests);

    expect(result.size).toBe(1);
    expect(result.get('helper')).toHaveLength(1);
    expect(result.get('helper')?.[0].testName).toBe('should work');
  });

  it('groups multiple tests by same function', () => {
    const tests: TestDetail[] = [
      createTest('helper', 'utils.test.ts', 'should work'),
      createTest('helper', 'utils.test.ts', 'should handle errors'),
      createTest('helper', 'utils.test.ts', 'should return correct value'),
    ];

    const result = groupTestsByFunction(tests);

    expect(result.size).toBe(1);
    expect(result.get('helper')).toHaveLength(3);
    expect(result.get('helper')?.map(t => t.testName)).toEqual([
      'should work',
      'should handle errors',
      'should return correct value',
    ]);
  });

  it('groups tests by different functions', () => {
    const tests: TestDetail[] = [
      createTest('helper', 'utils.test.ts', 'should work'),
      createTest('formatDate', 'date.test.ts', 'should format correctly'),
      createTest('parseJSON', 'json.test.ts', 'should parse valid JSON'),
    ];

    const result = groupTestsByFunction(tests);

    expect(result.size).toBe(3);
    expect(result.get('helper')).toHaveLength(1);
    expect(result.get('formatDate')).toHaveLength(1);
    expect(result.get('parseJSON')).toHaveLength(1);
  });

  it('handles mixed functions with multiple tests', () => {
    const tests: TestDetail[] = [
      createTest('helper', 'utils.test.ts', 'should work'),
      createTest('helper', 'utils.test.ts', 'should handle errors'),
      createTest('formatDate', 'date.test.ts', 'should format correctly'),
      createTest('formatDate', 'date.test.ts', 'should handle invalid dates'),
      createTest('parseJSON', 'json.test.ts', 'should parse valid JSON'),
    ];

    const result = groupTestsByFunction(tests);

    expect(result.size).toBe(3);
    expect(result.get('helper')).toHaveLength(2);
    expect(result.get('formatDate')).toHaveLength(2);
    expect(result.get('parseJSON')).toHaveLength(1);
  });

  it('trims whitespace from function names', () => {
    const tests: TestDetail[] = [
      { ...createTest('  helper  ', 'utils.test.ts', 'test1') },
      { ...createTest('helper', 'utils.test.ts', 'test2') },
    ];

    const result = groupTestsByFunction(tests);

    // After trimming, '  helper  ' becomes 'helper', so they should be grouped together
    expect(result.size).toBe(1);
    expect(result.get('helper')).toHaveLength(2);
  });

  it('preserves test order within groups', () => {
    const tests: TestDetail[] = [
      createTest('helper', 'utils.test.ts', 'first'),
      createTest('helper', 'utils.test.ts', 'second'),
      createTest('helper', 'utils.test.ts', 'third'),
    ];

    const result = groupTestsByFunction(tests);

    const helperTests = result.get('helper');
    expect(helperTests?.[0].testName).toBe('first');
    expect(helperTests?.[1].testName).toBe('second');
    expect(helperTests?.[2].testName).toBe('third');
  });

  it('handles tests from different test files', () => {
    const tests: TestDetail[] = [
      createTest('helper', 'utils.test.ts', 'should work'),
      createTest('helper', 'utils.spec.ts', 'should also work'),
      createTest('helper', 'utils.integration.test.ts', 'integration test'),
    ];

    const result = groupTestsByFunction(tests);

    expect(result.size).toBe(1);
    expect(result.get('helper')).toHaveLength(3);
  });

  it('handles function names with special characters', () => {
    const tests: TestDetail[] = [
      createTest('get_user_name', 'api.test.ts', 'test1'),
      createTest('getUserName', 'api.test.ts', 'test2'),
      createTest('$helper', 'utils.test.ts', 'test3'),
      createTest('日本語関数', 'utils.test.ts', 'test4'),
    ];

    const result = groupTestsByFunction(tests);

    expect(result.size).toBe(4);
    expect(result.get('get_user_name')).toHaveLength(1);
    expect(result.get('getUserName')).toHaveLength(1);
    expect(result.get('$helper')).toHaveLength(1);
    expect(result.get('日本語関数')).toHaveLength(1);
  });

  it('handles empty function name', () => {
    const tests: TestDetail[] = [
      createTest('', 'utils.test.ts', 'anonymous test'),
    ];

    const result = groupTestsByFunction(tests);

    expect(result.size).toBe(1);
    expect(result.get('')).toHaveLength(1);
  });

  it('handles large number of tests', () => {
    const tests: TestDetail[] = [];
    for (let i = 0; i < 100; i++) {
      tests.push(createTest(`func${i % 10}`, 'test.ts', `test${i}`));
    }

    const result = groupTestsByFunction(tests);

    expect(result.size).toBe(10);
    for (let i = 0; i < 10; i++) {
      expect(result.get(`func${i}`)).toHaveLength(10);
    }
  });

  it('preserves all test detail properties', () => {
    const testWithDetails: TestDetail = {
      functionName: 'helper',
      testFile: 'utils.test.ts',
      testName: 'should work',
      content: 'test content',
      lineRange: { start: 10, end: 20 },
      coveredLines: { start: 5, end: 15 },
      comment: 'Test comment',
      inputData: 'input',
      inputLines: { start: 1, end: 3 },
      expectedOutput: 'output',
      outputLines: { start: 4, end: 6 },
    };

    const result = groupTestsByFunction([testWithDetails]);

    const groupedTests = result.get('helper');
    expect(groupedTests).toHaveLength(1);
    expect(groupedTests?.[0]).toEqual(testWithDetails);
  });
});
