---
name: wk-workflow-manager
description: Manage and execute sequential AI workflows step-by-step. Use when you need to execute structured tasks in a controlled manner or when working with agents that require confirmation between steps.
---

# WK Workflow Manager

A skill for managing sequential workflows in AI agents. Use WK to execute tasks step-by-step with confirmation control and status tracking.

## When to Use WK

- You need to execute a sequence of tasks in a controlled manner
- The workflow requires confirmation or approval between steps
- You want to track progress through multiple stages of work
- You're working with an agent that should execute tasks sequentially

## Core Commands

### Initialize Your Workflow

Before using WK, define your workflow in `$HOME/.local/wk/workflow.yaml`:

```yaml
name: Your Workflow Name
steps:
  - id: step-id-1
    name: Step Name
    description: What this step does
    requires-confirmation: false
  - id: step-id-2
    name: Step Name
    description: What this step does
    requires-confirmation: true
```

Then run:

```bash
wk onboard
```

This will guide you through the workflow structure and available commands.

### Common Workflow Commands

**Start a new workflow execution:**
```bash
wk start
```

**Check current workflow status:**
```bash
wk status
```

**Move to the next step:**
```bash
wk next
```

**Add a message or note to the current step:**
```bash
wk say "Your message here"
```

## Workflow Configuration

### Basic Structure

Each step in your workflow has:
- `id`: Unique identifier for the step (lowercase, dashes)
- `name`: Human-readable name
- `description`: What the step does and any guidance for execution
- `requires-confirmation`: Optional boolean (default: false)

### Example Workflow

```yaml
name: Implement Feature
steps:
  - id: explore
    name: Explore Project
    description: |
      Explore the project structure to understand:
      - Programming language and framework
      - Build system and tooling
      - Key domain packages and business logic

  - id: design
    name: Design Solution
    description: Plan the implementation approach
    requires-confirmation: true

  - id: implement
    name: Implement Feature
    description: Write the actual code

  - id: test
    name: Test Implementation
    description: Write and run tests
    requires-confirmation: true

  - id: cleanup
    name: Cleanup and Documentation
    description: Polish code and update documentation
```

## Working With Workflows

### Step-by-Step Execution

1. **Start**: Run `wk start` to begin your workflow
2. **Work**: Execute the current step's tasks
3. **Update**: Use `wk say "message"` to track progress
4. **Next**: Run `wk next` to move to the next step
5. **Status**: Use `wk status` to check progress anytime

### With Confirmation Steps

When a step has `requires-confirmation: true`:
1. Complete the step's work
2. Run `wk status` to see the confirmation requirement
3. Get confirmation from the user or system
4. Run `wk next` once approved

## Monitoring Progress

Access the web interface to monitor and confirm steps:

```bash
go run ./cmd/web
```

Then visit `http://localhost:8080`

## Best Practices

1. **Clear descriptions**: Make each step's purpose explicit so agents understand what to do
2. **Logical grouping**: Group related tasks within a single step
3. **Confirmation strategically**: Use confirmation for critical milestones or when user input is needed
4. **Regular status checks**: Use `wk status` frequently to verify progress
5. **Meaningful messages**: Use `wk say` to document decisions and progress notes

## Integration with AI Agents

To use WK with Claude Code or other AI agents:

1. Copy the `workflow.yaml` file to `$HOME/.local/wk/workflow.yaml`
2. Run `wk onboard` once to initialize
3. Instruct the agent to use WK for task execution:

```
You have access to `wk`, a workflow management tool. Use it to structure your work.
Run `wk start` to begin, then execute each step sequentially using `wk status`,
`wk say`, and `wk next` commands.
```

The agent will then execute your workflow step-by-step with proper tracking.
