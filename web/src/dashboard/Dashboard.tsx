import { useEffect, useState } from 'react';
import { getHealth, type Health } from '../protocol/client';
import { createEventClient, type EventEnvelope } from '../protocol/events';

function applyEvent(current: Health, event: EventEnvelope): Health {
  if (event.type === 'project.changed' && event.payload.project) {
    return { ...current, project: event.payload.project as Health['project'] };
  }
  if (event.type === 'runtime.ready' || event.type === 'runtime.stopping') {
    return { ...current, runtime: { ...current.runtime, state: event.type === 'runtime.ready' ? 'ready' : 'stopping' } };
  }
  return current;
}

export function Dashboard() {
  const [health, setHealth] = useState<Health | null>(null);
  const [connection, setConnection] = useState<'connecting' | 'connected' | 'reconnecting' | 'degraded'>('connecting');
  useEffect(() => {
    getHealth().then(setHealth).catch(() => setConnection('degraded'));
    return createEventClient(snapshot => setHealth({ status: snapshot.diagnostics.some(check => check.status === 'fail') ? 'degraded' : 'ok', project: snapshot.project, runtime: snapshot.runtime, checks: snapshot.diagnostics }), setConnection, undefined, event => setHealth(current => current ? applyEvent(current, event) : current));
  }, []);
  return <main><h1>Craftlane</h1><p role="status">Connection: {connection}</p>{health ? <section aria-label="Project workspace">
    <h2>{health.project.name}</h2><dl><dt>Project root</dt><dd>{health.project.root}</dd><dt>Workspace address</dt><dd>{health.runtime.origin}</dd><dt>Health</dt><dd>{health.status}</dd><dt>Runtime</dt><dd>{health.runtime.state}</dd></dl>
    <h3>Diagnostics</h3><ul aria-label="Runtime diagnostics">{health.checks.map(check => <li key={check.name}><strong>{check.name}</strong>: {check.status} — {check.detail}{check.action && <>. Next: {check.action}</>}</li>)}</ul>
  </section> : <p>Runtime data unavailable</p>}</main>;
}
