import type { Step } from '../types';
import { StatusBadge } from './StatusBadge';

interface StepItemProps {
  step: Step;
}

export function StepItem({ step }: StepItemProps) {
  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'center',
        gap: '12px',
        padding: '8px 12px',
        borderLeft: '2px solid #e5e7eb',
        marginLeft: '8px',
      }}
    >
      <span
        style={{
          color: '#9ca3af',
          fontSize: '12px',
          minWidth: '20px',
        }}
      >
        {step.step_order}.
      </span>

      <div style={{ flex: 1 }}>
        <div style={{ fontWeight: 500, fontSize: '14px' }}>{step.title}</div>
        {step.description && (
          <div style={{ color: '#6b7280', fontSize: '12px', marginTop: '2px' }}>
            {step.description}
          </div>
        )}
      </div>

      <StatusBadge status={step.status} />

      {step.status === 'in_progress' && step.progress > 0 && (
        <div
          style={{
            width: '60px',
            height: '6px',
            backgroundColor: '#e5e7eb',
            borderRadius: '3px',
            overflow: 'hidden',
          }}
        >
          <div
            style={{
              width: `${step.progress}%`,
              height: '100%',
              backgroundColor: '#2563eb',
              transition: 'width 0.3s ease',
            }}
          />
        </div>
      )}
    </div>
  );
}
