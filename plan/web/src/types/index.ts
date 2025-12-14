export type Status = 'pending' | 'in_progress' | 'completed' | 'failed';

export interface Plan {
  id: string;
  parent_id: string | null;
  title: string;
  description: string | null;
  status: Status;
  created_at: string;
  updated_at: string;
  steps?: Step[];
  children?: Plan[];
}

export interface Step {
  id: string;
  plan_id: string;
  title: string;
  description: string | null;
  status: Status;
  step_order: number;
  progress: number;
  created_at: string;
  updated_at: string;
}

export interface SSEEvent {
  type: string;
  data: Plan | Step | { id: string };
}
