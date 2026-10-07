import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import fs from 'fs';
import os from 'os';
import path from 'path';
import { getConfig } from '../config/config.js';

describe('getConfig', () => {
  const prevBaseURL = process.env.API_BASE_URL;
  const prevBearer = process.env.API_BEARER_TOKEN;
  let homedirSpy;

  beforeEach(() => {
    delete process.env.API_BASE_URL;
    delete process.env.API_BEARER_TOKEN;
    // Isolate from a real ~/.api/config.json on the developer machine
    homedirSpy = vi.spyOn(os, 'homedir').mockReturnValue(
      path.join(os.tmpdir(), `wo-002-config-test-${process.pid}`)
    );
  });

  afterEach(() => {
    homedirSpy.mockRestore();
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

  it('returns baseURL and bearerToken from environment variables', () => {
    process.env.API_BASE_URL = 'https://api.example.com';
    process.env.API_BEARER_TOKEN = 'env-bearer-token';

    const config = getConfig();

    expect(config).toEqual({
      baseURL: 'https://api.example.com',
      bearerToken: 'env-bearer-token',
    });
  });

  it('throws when env vars are unset and no config file exists', () => {
    const missingHome = path.join(os.tmpdir(), `wo-002-missing-${process.pid}-${Date.now()}`);
    homedirSpy.mockReturnValue(missingHome);
    expect(fs.existsSync(path.join(missingHome, '.api', 'config.json'))).toBe(false);

    expect(() => getConfig()).toThrowError(/Configuration not found/);
  });
});
