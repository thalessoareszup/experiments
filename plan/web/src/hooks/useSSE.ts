import { useEffect, useState, useCallback } from 'react';
import type { Plan, Step } from '../types';

interface UseSSEOptions {
  onPlanCreated?: (plan: Plan) => void;
  onPlanUpdated?: (plan: Plan) => void;
  onPlanDeleted?: (id: string) => void;
  onStepCreated?: (step: Step) => void;
  onStepUpdated?: (step: Step) => void;
  onStepDeleted?: (id: string) => void;
}

export function useSSE(url: string, options: UseSSEOptions = {}) {
  const [connected, setConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const connect = useCallback(() => {
    const eventSource = new EventSource(url);

    eventSource.onopen = () => {
      setConnected(true);
      setError(null);
    };

    eventSource.onerror = () => {
      setConnected(false);
      setError('Connection lost. Reconnecting...');
    };

    eventSource.addEventListener('connected', () => {
      setConnected(true);
    });

    eventSource.addEventListener('plan:created', (e) => {
      const plan = JSON.parse(e.data) as Plan;
      options.onPlanCreated?.(plan);
    });

    eventSource.addEventListener('plan:updated', (e) => {
      const plan = JSON.parse(e.data) as Plan;
      options.onPlanUpdated?.(plan);
    });

    eventSource.addEventListener('plan:deleted', (e) => {
      const data = JSON.parse(e.data) as { id: string };
      options.onPlanDeleted?.(data.id);
    });

    eventSource.addEventListener('step:created', (e) => {
      const step = JSON.parse(e.data) as Step;
      options.onStepCreated?.(step);
    });

    eventSource.addEventListener('step:updated', (e) => {
      const step = JSON.parse(e.data) as Step;
      options.onStepUpdated?.(step);
    });

    eventSource.addEventListener('step:deleted', (e) => {
      const data = JSON.parse(e.data) as { id: string };
      options.onStepDeleted?.(data.id);
    });

    return eventSource;
  }, [url, options]);

  useEffect(() => {
    const eventSource = connect();
    return () => eventSource.close();
  }, [connect]);

  return { connected, error };
}
