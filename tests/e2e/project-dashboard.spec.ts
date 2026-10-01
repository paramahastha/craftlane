import playwright from '../../web/node_modules/@playwright/test/index';
import { execFile, execFileSync } from 'node:child_process';
import { mkdtempSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { createServer } from 'node:net';

const { expect, test } = playwright;

test('built binary serves a project-identifying dashboard', async ({ page }) => {
  const root = resolve(import.meta.dirname, '../..');
  const work = mkdtempSync(join(tmpdir(), 'craftlane-e2e-'));
  const project = join(work, 'demo-project');
  const state = join(work, 'state');
  const binary = join(work, process.platform === 'win32' ? 'craftlane.exe' : 'craftlane');
  execFileSync('git', ['init', project]);
  execFileSync('go', ['build', '-o', binary, './cmd/craftlane'], { cwd: root });
  const port = await new Promise<number>((resolvePort, reject) => {
    const server = createServer();
    server.once('error', reject);
    server.listen(0, '127.0.0.1', () => {
      const address = server.address();
      if (!address || typeof address === 'string') return reject(new Error('no port'));
      server.close(() => resolvePort(address.port));
    });
  });
  const env = { ...process.env, CRAFTLANE_STATE_DIR: state };
  await new Promise<void>((done, reject) => execFile(binary, ['start', '--no-open', '--port', String(port)], { cwd: project, env }, error => error ? reject(error) : done()));
  try {
    await page.goto(`http://127.0.0.1:${port}`);
    await expect(page.getByRole('heading', { name: 'demo-project' })).toBeVisible();
    await expect(page.getByText(project)).toBeVisible();
    await expect(page.getByText(`http://127.0.0.1:${port}`)).toBeVisible();
    await expect(page.getByText('ready')).toBeVisible();
  } finally {
    execFileSync(binary, ['stop'], { cwd: project, env });
  }
});
