import type { Step } from '../types';
import { StepItem } from './StepItem';

interface StepListProps {
  steps: Step[];
}

export function StepList({ steps }: StepListProps) {
  if (steps.length === 0) {
    return (
      <div style={{ color: '#9ca3af', fontSize: '13px', padding: '8px 12px' }}>
        No steps defined
      </div>
    );
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
      {steps.map((step) => (
        <StepItem key={step.id} step={step} />
      ))}
    </div>
  );
}
