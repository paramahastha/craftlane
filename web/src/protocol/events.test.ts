import { beforeEach, describe, expect, it, vi } from 'vitest';
import { createEventClient } from './events';

class FakeSocket extends EventTarget {
  sent: string[] = [];
  send(value: string) { this.sent.push(value); }
  close() { this.dispatchEvent(new Event('close')); }
}

describe('event client', () => {
  beforeEach(() => {
    const values = new Map<string, string>();
    vi.stubGlobal('localStorage', {
      getItem: (key: string) => values.get(key) ?? null,
      setItem: (key: string, value: string) => values.set(key, value),
      removeItem: (key: string) => values.delete(key),
      clear: () => values.clear(),
    });
  });

  it('persists snapshot cursors and resumes from them', () => {
    localStorage.setItem('craftlane.events.v1.cursor', '7');
    const socket = new FakeSocket();
    const states: string[] = [];
    const stop = createEventClient(vi.fn(), state => states.push(state), () => socket as unknown as WebSocket);
    socket.dispatchEvent(new Event('open'));
    expect(JSON.parse(socket.sent[0])).toMatchObject({ protocolVersion: 1, lastSequence: 7 });
    socket.dispatchEvent(new MessageEvent('message', { data: JSON.stringify({ type: 'snapshot', cursor: 9, snapshot: { sequence: 9, project: {}, runtime: {}, diagnostics: [] } }) }));
    expect(localStorage.getItem('craftlane.events.v1.cursor')).toBe('9');
    expect(states).toContain('connected');
    stop();
  });

  it('suppresses duplicate and older event cursors and clears continuity on resync', () => {
    const socket = new FakeSocket();
    const stop = createEventClient(vi.fn(), vi.fn(), () => socket as unknown as WebSocket);
    const event = { schemaVersion: 1, eventId: 'e1', sequence: 3, timestamp: '', projectId: 'p', source: { name: 'core' }, type: 'project.changed', payload: {} };
    socket.dispatchEvent(new MessageEvent('message', { data: JSON.stringify({ type: 'event', event }) }));
    socket.dispatchEvent(new MessageEvent('message', { data: JSON.stringify({ type: 'event', event }) }));
    expect(localStorage.getItem('craftlane.events.v1.cursor')).toBe('3');
    socket.dispatchEvent(new MessageEvent('message', { data: JSON.stringify({ type: 'resync_required', cursor: 10 }) }));
    expect(localStorage.getItem('craftlane.events.v1.cursor')).toBeNull();
    stop();
  });
});
