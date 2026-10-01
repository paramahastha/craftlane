import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Dashboard } from './Dashboard';

vi.mock('../protocol/events', () => ({ createEventClient: (_snapshot: unknown, state: (value: string) => void) => { state('connected'); return () => {}; } }));

describe('Dashboard', () => {
  beforeEach(() => vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, json: async () => ({
    status: 'ok', project: { id: 'p', name: 'demo', root: '/work/demo' },
    runtime: { id: 'r', origin: 'http://127.0.0.1:7331', protocolVersion: 1, state: 'ready', startedAt: new Date().toISOString() },
    checks: [{ name: 'git', status: 'pass', detail: 'Git is available' }],
  }) })));

  it('renders labelled project, root, origin, health, and diagnostics', async () => {
    render(<Dashboard />);
    expect(await screen.findByRole('heading', { name: 'demo' })).toBeVisible();
    expect(screen.getByText('/work/demo')).toBeVisible();
    expect(screen.getByText('http://127.0.0.1:7331')).toBeVisible();
    expect(screen.getByLabelText('Runtime diagnostics')).toHaveTextContent('git: pass');
    expect(screen.getByRole('status')).toHaveTextContent('connected');
  });
});
