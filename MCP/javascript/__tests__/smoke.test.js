import { describe, it, expect } from 'vitest';
import { createGetAllTool } from '../registry.js';

describe('vitest smoke', () => {
  it('imports createGetAllTool via ESM and returns a valid tool shape', () => {
    const tool = createGetAllTool();

    expect(tool).toBeTypeOf('object');
    expect(tool.definition).toBeTypeOf('object');
    expect(tool.definition.name).toBe('get-all');
    expect(tool.handler).toBeTypeOf('function');
  });
});
