import { useEffect, useState } from 'react';
import { getHealth, type Agent, type Health } from '../protocol/client';
import { createEventClient, type EventEnvelope } from '../protocol/events';

function applyEvent(current: Health, event: EventEnvelope): Health {
  if (event.type === 'project.changed' && event.payload.project) {
    return { ...current, project: event.payload.project as Health['project'] };
  }
  if (event.type === 'runtime.ready' || event.type === 'runtime.stopping') {
    return { ...current, runtime: { ...current.runtime, state: event.type === 'runtime.ready' ? 'ready' : 'stopping' } };
  }
  if (event.type === 'agent.updated' && event.payload.agent) {
    const agent = event.payload.agent as Agent;
    const agents = [...(current.agents || [])];
    const index = agents.findIndex(item => item.id === agent.id);
    if (index < 0) agents.push(agent); else agents[index] = agent;
    return { ...current, agents };
  }
  return current;
}

export function Dashboard() {
  const [health, setHealth] = useState<Health | null>(null);
  const [connection, setConnection] = useState<'connecting' | 'connected' | 'reconnecting' | 'degraded'>('connecting');
  useEffect(() => {
    getHealth().then(setHealth).catch(() => setConnection('degraded'));
    return createEventClient(snapshot => setHealth({ status: snapshot.diagnostics.some(check => check.status === 'fail') ? 'degraded' : 'ok', project: snapshot.project, runtime: snapshot.runtime, checks: snapshot.diagnostics, agents: snapshot.agents, providerStatuses: snapshot.providerStatuses }), setConnection, undefined, event => setHealth(current => current ? applyEvent(current, event) : current));
  }, []);
  return <main><h1>Craftlane</h1><p role="status">Connection: {connection}</p>{health ? <section aria-label="Project workspace">
    <h2>{health.project.name}</h2><dl><dt>Project root</dt><dd>{health.project.root}</dd><dt>Workspace address</dt><dd>{health.runtime.origin}</dd><dt>Health</dt><dd>{health.status}</dd><dt>Runtime</dt><dd>{health.runtime.state}</dd></dl>
    <h3>Provider coverage</h3><ul aria-label="Provider coverage">{(health.providerStatuses || []).map(status => <li key={status.provider}><strong>{status.provider}</strong>: {status.setupState} — {status.detail}{status.action && <>. Next: {status.action}</>}</li>)}</ul>
    <h3>Agents</h3>{health.agents?.length ? <ul aria-label="Observed agents">{health.agents.map(agent => <li key={agent.id}><strong>{agent.displayName}</strong> ({agent.provider}): {agent.state} ({agent.stateSource}){agent.currentTool && <> — {agent.currentTool.name} {agent.currentTool.state}</>}{agent.parentAgentId && <> — child of {agent.parentAgentId}</>}{agent.priorHistoryUnavailable && <> — prior history unavailable</>}{agent.staleAfterRestart && <> — stale after restart</>}</li>)}</ul> : <p>No observed agents yet.</p>}
    <h3>Diagnostics</h3><ul aria-label="Runtime diagnostics">{health.checks.map(check => <li key={check.name}><strong>{check.name}</strong>: {check.status} — {check.detail}{check.action && <>. Next: {check.action}</>}</li>)}</ul>
  </section> : <p>Runtime data unavailable</p>}</main>;
}
