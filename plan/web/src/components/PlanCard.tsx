import type { Plan } from '../types';
import { StatusBadge } from './StatusBadge';
import { StepList } from './StepList';

interface PlanCardProps {
  plan: Plan;
  depth?: number;
}

export function PlanCard({ plan, depth = 0 }: PlanCardProps) {
  const completedSteps = plan.steps?.filter(s => s.status === 'completed').length || 0;
  const totalSteps = plan.steps?.length || 0;

  return (
    <div
      style={{
        backgroundColor: '#fff',
        border: '1px solid #e5e7eb',
        borderRadius: '8px',
        marginLeft: depth * 24,
        marginBottom: '12px',
        overflow: 'hidden',
      }}
    >
      <div
        style={{
          padding: '16px',
          borderBottom: '1px solid #e5e7eb',
          display: 'flex',
          alignItems: 'flex-start',
          justifyContent: 'space-between',
          gap: '12px',
        }}
      >
        <div>
          <h3 style={{ margin: 0, fontSize: '16px', fontWeight: 600 }}>
            {plan.title}
          </h3>
          {plan.description && (
            <p style={{ margin: '4px 0 0', color: '#6b7280', fontSize: '14px' }}>
              {plan.description}
            </p>
          )}
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '12px',
              marginTop: '8px',
              fontSize: '12px',
              color: '#9ca3af',
            }}
          >
            <span>ID: {plan.id.slice(0, 8)}...</span>
            {totalSteps > 0 && (
              <span>
                {completedSteps}/{totalSteps} steps
              </span>
            )}
            <span>
              Updated: {new Date(plan.updated_at).toLocaleTimeString()}
            </span>
          </div>
        </div>
        <StatusBadge status={plan.status} />
      </div>

      {plan.steps && plan.steps.length > 0 && (
        <div style={{ padding: '12px 16px' }}>
          <StepList steps={plan.steps} />
        </div>
      )}

      {plan.children && plan.children.length > 0 && (
        <div style={{ padding: '0 16px 16px' }}>
          <div
            style={{
              fontSize: '12px',
              fontWeight: 500,
              color: '#6b7280',
              marginBottom: '8px',
              textTransform: 'uppercase',
            }}
          >
            Child Plans
          </div>
          {plan.children.map((child) => (
            <PlanCard key={child.id} plan={child} depth={depth + 1} />
          ))}
        </div>
      )}
    </div>
  );
}
