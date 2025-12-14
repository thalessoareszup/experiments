import type { Plan } from '../types';
import { PlanCard } from './PlanCard';

interface PlanListProps {
  plans: Plan[];
}

export function PlanList({ plans }: PlanListProps) {
  if (plans.length === 0) {
    return (
      <div
        style={{
          textAlign: 'center',
          padding: '48px',
          color: '#6b7280',
        }}
      >
        <div style={{ fontSize: '48px', marginBottom: '16px' }}>
          <span role="img" aria-label="clipboard">
          </span>
        </div>
        <h3 style={{ margin: '0 0 8px', color: '#374151' }}>No plans yet</h3>
        <p style={{ margin: 0 }}>
          Use the CLI to create a plan: <code>plan start --title "My Plan"</code>
        </p>
      </div>
    );
  }

  return (
    <div>
      {plans.map((plan) => (
        <PlanCard key={plan.id} plan={plan} />
      ))}
    </div>
  );
}
