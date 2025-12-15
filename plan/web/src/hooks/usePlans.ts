import { useState, useEffect, useCallback } from 'react';
import type { Plan, Step } from '../types';
import { useWebSocket } from './useWebSocket';

const API_BASE = '/api';

export function usePlans() {
  const [plans, setPlans] = useState<Plan[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchPlans = useCallback(async () => {
    try {
      const response = await fetch(`${API_BASE}/plans`);
      if (!response.ok) throw new Error('Failed to fetch plans');
      const data = await response.json();
      setPlans(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  }, []);

  // SSE handlers
  const handlePlanCreated = useCallback((plan: Plan) => {
    setPlans(prev => [plan, ...prev]);
  }, []);

  const handlePlanUpdated = useCallback((updatedPlan: Plan) => {
    setPlans(prev => prev.map(p =>
      p.id === updatedPlan.id ? { ...p, ...updatedPlan } : p
    ));
  }, []);

  const handlePlanDeleted = useCallback((id: string) => {
    setPlans(prev => prev.filter(p => p.id !== id));
  }, []);

  const handleStepCreated = useCallback((step: Step) => {
    setPlans(prev => prev.map(p => {
      if (p.id === step.plan_id) {
        return {
          ...p,
          steps: [...(p.steps || []), step].sort((a, b) => a.step_order - b.step_order)
        };
      }
      return p;
    }));
  }, []);

  const handleStepUpdated = useCallback((updatedStep: Step) => {
    setPlans(prev => prev.map(p => {
      if (p.id === updatedStep.plan_id) {
        return {
          ...p,
          steps: (p.steps || []).map(s =>
            s.id === updatedStep.id ? updatedStep : s
          )
        };
      }
      return p;
    }));
  }, []);

  const handleStepDeleted = useCallback((id: string) => {
    setPlans(prev => prev.map(p => ({
      ...p,
      steps: (p.steps || []).filter(s => s.id !== id)
    })));
  }, []);

  const { connectionState } = useWebSocket(`${API_BASE}/ws`, {
    onPlanCreated: handlePlanCreated,
    onPlanUpdated: handlePlanUpdated,
    onPlanDeleted: handlePlanDeleted,
    onStepCreated: handleStepCreated,
    onStepUpdated: handleStepUpdated,
    onStepDeleted: handleStepDeleted,
  });

  useEffect(() => {
    fetchPlans();
  }, [fetchPlans]);

  // Refetch on reconnection to ensure consistency
  useEffect(() => {
    if (connectionState === 'connected') {
      fetchPlans();
    }
  }, [connectionState, fetchPlans]);

  return { plans, loading, error, connected: connectionState === 'connected', refetch: fetchPlans };
}
