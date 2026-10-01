import { defineConfig } from './web/node_modules/@playwright/test/index';

export default defineConfig({ testDir: './tests/e2e', timeout: 30_000, use: { headless: true } });
