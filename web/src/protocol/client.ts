export type Diagnostic = { name: string; status: 'pass' | 'warn' | 'fail'; detail: string; action?: string };
export type Project = { id: string; name: string; root: string; gitCommonDir?: string; capabilities?: Record<string, string> };
export type Runtime = { id: string; origin: string; protocolVersion: number; state: string; startedAt: string };
export type Agent = { id: string; provider: string; projectId: string; sessionId?: string; displayName: string; parentAgentId?: string; state: string; stateSource: string; currentTool?: { id?: string; name: string; state: string; outcome?: string }; task?: string; lastObservedAt: string; priorHistoryUnavailable?: boolean; staleAfterRestart?: boolean };
export type ProviderStatus = { provider: string; setupState: string; capabilities?: Record<string, { status: string; reason?: string }>; detail?: string; action?: string; adapterVersion?: string };
export type Health = { status: string; project: Project; runtime: Runtime; checks: Diagnostic[]; agents?: Agent[]; providerStatuses?: ProviderStatus[] };

export async function getHealth(): Promise<Health> {
  const response = await fetch('/api/v1/health');
  if (!response.ok) throw new Error('health request failed');
  return response.json() as Promise<Health>;
}
