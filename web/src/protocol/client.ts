export type Diagnostic = { name: string; status: 'pass' | 'warn' | 'fail'; detail: string; action?: string };
export type Project = { id: string; name: string; root: string };
export type Runtime = { id: string; origin: string; protocolVersion: number; state: string; startedAt: string };
export type Health = { status: string; project: Project; runtime: Runtime; checks: Diagnostic[] };

export async function getHealth(): Promise<Health> {
  const response = await fetch('/api/v1/health');
  if (!response.ok) throw new Error('health request failed');
  return response.json() as Promise<Health>;
}
