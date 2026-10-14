import type { Diagnostic, Project, Runtime } from './client';

export type EventEnvelope = { schemaVersion: number; eventId: string; sequence: number; timestamp: string; projectId: string; source: { name: string; adapterVersion?: string }; type: string; payload: Record<string, unknown> };
export type AgentState = 'working' | 'waiting' | 'waiting_permission' | 'idle' | 'completed' | 'failed' | 'unknown';
export type ToolActivity = { id?: string; name: string; state: string; outcome?: string };
export type Agent = { id: string; provider: string; projectId: string; sessionId?: string; displayName: string; parentAgentId?: string; state: AgentState; stateSource: 'reported' | 'inferred'; currentTool?: ToolActivity; task?: string; lastObservedAt: string; reportedAt?: string; priorHistoryUnavailable?: boolean; staleAfterRestart?: boolean };
export type ProviderStatus = { provider: string; setupState: string; capabilities?: Record<string, { status: string; reason?: string }>; detail?: string; action?: string; adapterVersion?: string };
export type Snapshot = { sequence: number; project: Project; runtime: Runtime; diagnostics: Diagnostic[]; agents?: Agent[]; providerStatuses?: ProviderStatus[] };
type Message = { type: 'snapshot'; cursor: number; snapshot: Snapshot } | { type: 'event'; event: EventEnvelope } | { type: 'resync_required'; cursor: number } | { type: 'heartbeat'; cursor: number } | { type: 'error'; message: string };
const cursorKey = 'craftlane.events.v1.cursor';

export function createEventClient(onSnapshot: (snapshot: Snapshot) => void, onState: (state: 'connected' | 'reconnecting' | 'degraded') => void, makeSocket: (url: string, protocol: string) => WebSocket = (url, protocol) => new WebSocket(url, protocol), onEvent: (event: EventEnvelope) => void = () => {}) {
  let socket: WebSocket | undefined;
  let stopped = false;
  let retry = 250;
  let cursor = Number(localStorage.getItem(cursorKey) || 0);
  const seen = new Set<string>();
  const connect = () => {
    if (stopped) return;
    const scheme = location.protocol === 'https:' ? 'wss:' : 'ws:';
    socket = makeSocket(`${scheme}//${location.host}/ws`, 'craftlane.events.v1');
    socket.addEventListener('open', () => {
      retry = 250;
      onState('connected');
      socket?.send(JSON.stringify({ type: 'subscribe', protocolVersion: 1, lastSequence: cursor }));
    });
    socket.addEventListener('message', ({ data }) => {
      const message = JSON.parse(String(data)) as Message;
      if (message.type === 'resync_required') {
        seen.clear(); cursor = 0; localStorage.removeItem(cursorKey);
      } else if (message.type === 'snapshot') {
        onSnapshot(message.snapshot); cursor = message.cursor; localStorage.setItem(cursorKey, String(cursor));
      } else if (message.type === 'event' && message.event.sequence > cursor && !seen.has(message.event.eventId)) {
        seen.add(message.event.eventId); onEvent(message.event); cursor = message.event.sequence; localStorage.setItem(cursorKey, String(cursor));
      } else if (message.type === 'error') onState('degraded');
    });
    socket.addEventListener('close', () => {
      if (stopped) return;
      onState('reconnecting'); window.setTimeout(connect, retry); retry = Math.min(retry * 2, 5000);
    });
    socket.addEventListener('error', () => onState('degraded'));
  };
  connect();
  return () => { stopped = true; socket?.close(); };
}
