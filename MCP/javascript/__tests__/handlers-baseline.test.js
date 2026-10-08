import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { get_all } from '../registry.js';
import { load_api_config } from '../config/config.js';
import { models } from '../models/models.js';

/**
 * Post-fix characterization: handlers no longer TDZ-crash on `url`.
 * With env config set they reach fetch(); placeholder `/api/unknown` may
 * fail the request, but the result is always a string (JSON or error text).
 */
const TDZ_URL_MESSAGE = /Cannot access 'url' before initialization/;

describe('handler no-crash after url declaration fix', () => {
  const prevBaseURL = process.env.API_BASE_URL;
  const prevBearer = process.env.API_BEARER_TOKEN;

  beforeEach(() => {
    process.env.API_BASE_URL = 'https://example.com';
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

  it('get_all returns a string without ReferenceError TDZ', async () => {
    const result = await get_all();
    expect(typeof result).toBe('string');
    expect(result).not.toMatch(TDZ_URL_MESSAGE);
  });

  it('load_api_config returns a string without ReferenceError TDZ', async () => {
    const result = await load_api_config();
    expect(typeof result).toBe('string');
    expect(result).not.toMatch(TDZ_URL_MESSAGE);
  });

  it('models returns a string without ReferenceError TDZ', async () => {
    const result = await models();
    expect(typeof result).toBe('string');
    expect(result).not.toMatch(TDZ_URL_MESSAGE);
  });
});
