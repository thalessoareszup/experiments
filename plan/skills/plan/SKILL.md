---
name: plan
description: Track and coordinate multi-step workflows using the plan CLI tool. Use this skill when working on complex tasks that benefit from progress tracking, or when coordinating with sub-agents.
---

# Plan Workflow Tracking Skill

This skill enables structured workflow tracking using the `plan` CLI tool. Use it for:
- Complex multi-step tasks
- Sub-agent coordination
- Progress visibility to users

## Quick Start

### Starting a New Plan
```bash
plan start --title "Your task description"
```
This outputs a JSON response with the plan ID. Note this ID for subsequent commands.

### Adding Steps
```bash
plan step --plan <plan-id> --title "Step 1: Research"
plan step --plan <plan-id> --title "Step 2: Implementation"
plan step --plan <plan-id> --title "Step 3: Testing"
```

### Updating Progress
```bash
# Mark step as in progress
plan progress --step <step-id> --status in_progress

# Update percentage completion
plan progress --step <step-id> --percent 50

# Complete a step
plan complete --step <step-id>
```

### Completing or Failing
```bash
# Complete the plan
plan complete --plan <plan-id>

# Or mark as failed with reason
plan fail --plan <plan-id> --reason "Encountered blocking issue"
```

## CLI Commands Reference

| Command | Description | Key Flags |
|---------|-------------|-----------|
| `plan start` | Start new plan | `--title`, `--parent`, `--description` |
| `plan step` | Add step to plan | `--plan`, `--title`, `--order` |
| `plan progress` | Update step progress | `--step`, `--percent`, `--status` |
| `plan complete` | Mark complete | `--step` or `--plan` |
| `plan fail` | Mark failed | `--step` or `--plan`, `--reason` |
| `plan query` | Query plans/steps | `--plan`, `--status`, `--children` |
| `plan serve` | Start webapp server | `--port` (default 8080) |

## Sub-Agent Coordination

When spawning sub-agents, pass the parent plan ID:

```bash
# Parent agent creates main plan
plan start --title "Main task"
# Returns: {"id": "parent-plan-id", ...}

# Sub-agent creates child plan
plan start --title "Subtask" --parent parent-plan-id
```

This creates a hierarchy visible in the webapp.

## Querying Status

```bash
# Get current plan status
plan query --plan <plan-id>

# Get all active plans
plan query --status in_progress

# Get plan with children
plan query --plan <plan-id> --children
```

## Environment Variables

- `PLAN_SESSION_ID` - Set automatically by `plan start`, used as default plan ID
- `PLAN_DB_PATH` - Override database location (default: `~/.local/plan/plan.db`)

## Best Practices

1. **Create plan at task start**: Begin with `plan start` to establish tracking
2. **Break down into steps**: Add 3-7 steps for visibility
3. **Update progress frequently**: Call `plan progress` at meaningful milestones
4. **Handle failures gracefully**: Use `plan fail` with descriptive reasons
5. **Complete plans explicitly**: Always call `plan complete` when done

## Example: Feature Implementation Workflow

```bash
# Start tracking
PLAN_ID=$(plan start --title "Implement user authentication" | jq -r '.id')

# Add steps
STEP1=$(plan step --plan $PLAN_ID --title "Design database schema" | jq -r '.id')
STEP2=$(plan step --plan $PLAN_ID --title "Implement API endpoints" | jq -r '.id')
STEP3=$(plan step --plan $PLAN_ID --title "Add frontend components" | jq -r '.id')
STEP4=$(plan step --plan $PLAN_ID --title "Write tests" | jq -r '.id')

# Work on step 1
plan progress --step $STEP1 --status in_progress
# ... do work ...
plan complete --step $STEP1

# Work on step 2
plan progress --step $STEP2 --status in_progress
plan progress --step $STEP2 --percent 50
# ... more work ...
plan complete --step $STEP2

# Continue with remaining steps...

# When done
plan complete --plan $PLAN_ID
```

## Viewing Progress

Start the webapp to visualize plans in real-time:

```bash
plan serve
# Open http://localhost:8080
```

The webapp shows:
- All active plans and their steps
- Real-time updates as progress changes
- Hierarchical view of parent-child plans
- Status indicators with color coding

## Creating Custom Workflows

Users can define skills that use the plan tool to follow specific workflows:

```markdown
---
name: spec-driven
description: Implement features following a spec-driven approach
---

# Spec-Driven Development

1. Start a plan with the feature name
2. Add steps:
   - Review specification
   - Design implementation
   - Write tests first
   - Implement feature
   - Verify against spec
3. Update progress as you work
4. Complete when spec requirements are met
```

This allows teams to standardize workflows while maintaining visibility into agent progress.
