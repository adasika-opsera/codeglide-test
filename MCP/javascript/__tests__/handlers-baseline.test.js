import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { get_all } from '../registry.js';
import { load_api_config } from '../config/config.js';
import { models } from '../models/models.js';

/**
 * Characterization baseline: each handler references `url` before its `const`
 * declaration (TDZ). The ReferenceError is caught by the handler try/catch and
 * returned as a string — document that current broken behavior here.
 */
const TDZ_URL_MESSAGE = /Cannot access 'url' before initialization/;

describe('handler ReferenceError baseline', () => {
  const prevBaseURL = process.env.API_BASE_URL;
  const prevBearer = process.env.API_BEARER_TOKEN;

  beforeEach(() => {
    process.env.API_BASE_URL = 'https://api.example.com';
    process.env.API_BEARER_TOKEN = 'test-token';
  });

  afterEach(() => {
    if (prevBaseURL === undefined) {
      delete process.env.API_BASE_URL;
    } else {
      process.env.API_BASE_URL = prevBaseURL;
    }
    if (prevBearer === undefined) {
      delete process.env.API_BEARER_TOKEN;
    } else {
      process.env.API_BEARER_TOKEN = prevBearer;
    }
  });

  it('get_all fails with ReferenceError TDZ on url', async () => {
    const result = await get_all();
    expect(result).toMatch(TDZ_URL_MESSAGE);
    expect(result).toMatch(/^Request failed:/);
  });

  it('load_api_config fails with ReferenceError TDZ on url', async () => {
    const result = await load_api_config();
    expect(result).toMatch(TDZ_URL_MESSAGE);
    expect(result).toMatch(/^Request failed:/);
  });

  it('models fails with ReferenceError TDZ on url', async () => {
    const result = await models();
    expect(result).toMatch(TDZ_URL_MESSAGE);
    expect(result).toMatch(/^Request failed:/);
  });
});
