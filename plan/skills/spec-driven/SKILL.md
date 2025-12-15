---
name: spec-driven
description: |
  Implement features using spec-driven development (SDD).
  Use this skill for features that benefit from structured specifications before implementation.
  The workflow follows: Specify → Plan → Tasks → Implement, with progress tracked via the plan CLI.
---

# Spec-Driven Development Skill

This skill implements a structured workflow where **specifications are the source of truth** and code serves the spec.

## When to Use

- New features with unclear requirements
- Complex features requiring architectural decisions
- Features that need stakeholder alignment
- Work that will be handed off or reviewed

## Workflow Overview

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌───────────┐
│ SPECIFY │ → │  PLAN   │ → │  TASKS  │ → │ IMPLEMENT │
└─────────┘    └─────────┘    └─────────┘    └───────────┘
   spec.md       plan.md       tasks.md        code + tests
```

## Getting Started

When the user requests a feature using spec-driven development:

```bash
# 1. Create the plan and track workflow
PLAN_ID=$(plan start --title "SDD: <feature-name>" | jq -r '.id')

# 2. Create spec directory structure
mkdir -p specs/<feature-branch>/{contracts,research}

# 3. Add workflow steps
plan step --plan $PLAN_ID --title "Specify: Define requirements and acceptance criteria"
plan step --plan $PLAN_ID --title "Plan: Design technical implementation"
plan step --plan $PLAN_ID --title "Tasks: Break down into executable work items"
plan step --plan $PLAN_ID --title "Implement: Build with test-driven development"
plan step --plan $PLAN_ID --title "Validate: Verify against specification"
```

---

## Phase 1: SPECIFY

**Goal:** Transform vague ideas into structured requirements.

```bash
plan progress --step <specify-step-id> --status in_progress
```

Create `specs/000X-<feature>/spec.md`:

```markdown
# Feature: <Feature Name>

## Overview
<One paragraph describing what this feature does and why it matters>

## User Stories

### Story 1: <Primary use case>
As a <user type>,
I want to <action>,
So that <benefit>.

**Acceptance Criteria:**
- [ ] <Criterion 1>
- [ ] <Criterion 2>
- [ ] <Criterion 3>

### Story 2: <Secondary use case>
...

## Constraints
- <Technical constraint>
- <Business constraint>
- <Security/compliance requirement>

## Out of Scope
- <What this feature explicitly does NOT do>

## Open Questions
- [ ] <Question needing clarification>
```

**Completion check:** All user stories have acceptance criteria, constraints are documented.

```bash
plan complete --step <specify-step-id>
```

---

## Phase 2: PLAN

**Goal:** Convert specification into technical implementation strategy.

```bash
plan progress --step <plan-step-id> --status in_progress
```

Create `specs/<feature>/plan.md`:

```markdown
# Implementation Plan: <Feature Name>

## Technical Approach
<High-level description of how the feature will be built>

## Architecture Decisions

### Decision 1: <Topic>
- **Context:** <Why this decision is needed>
- **Options considered:**
  1. <Option A> - <pros/cons>
  2. <Option B> - <pros/cons>
- **Decision:** <Chosen option>
- **Rationale:** <Why this option was chosen>

## Data Model

### Entity: <EntityName>
| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| ... | ... | ... |

## API Contracts

### Endpoint: <METHOD /path>
- **Request:** `{ field: type }`
- **Response:** `{ field: type }`
- **Errors:** 400 (validation), 404 (not found), 500 (server)

## Dependencies
- <Library/service> - <purpose>

## Risks & Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| <Risk> | High/Med/Low | <How to address> |
```

Optional: Create `specs/<feature>/research.md` for technology analysis.

**Completion check:** Architecture decisions made, data model defined, API contracts specified.

```bash
plan complete --step <plan-step-id>
```

---

## Phase 3: TASKS

**Goal:** Break down the plan into executable, parallelizable work items.

```bash
plan progress --step <tasks-step-id> --status in_progress
```

Create `specs/<feature>/tasks.md`:

```markdown
# Tasks: <Feature Name>

## Task Dependency Graph
```
[1] ─┬─→ [2] ─→ [4]
     └─→ [3] ─┘
```

## Tasks

### Task 1: <Foundation task>
- **Description:** <What needs to be done>
- **Files:** `path/to/file.ts`, `path/to/other.ts`
- **Acceptance:** <How to verify completion>
- **Parallelizable:** No (blocks other tasks)

### Task 2: <Implementation task>
- **Description:** <What needs to be done>
- **Depends on:** Task 1
- **Files:** `path/to/file.ts`
- **Acceptance:** <How to verify completion>
- **Parallelizable:** Yes (with Task 3)

### Task 3: <Another implementation task>
- **Description:** <What needs to be done>
- **Depends on:** Task 1
- **Files:** `path/to/file.ts`
- **Acceptance:** <How to verify completion>
- **Parallelizable:** Yes (with Task 2)

### Task 4: <Integration task>
- **Description:** <What needs to be done>
- **Depends on:** Task 2, Task 3
- **Files:** `path/to/file.ts`
- **Acceptance:** <How to verify completion>
- **Parallelizable:** No
```

**Completion check:** Each task has clear acceptance criteria and file targets.

```bash
plan complete --step <tasks-step-id>
```

---

## Phase 4: IMPLEMENT

**Goal:** Generate working code using test-driven development.

```bash
plan progress --step <implement-step-id> --status in_progress
```

For each task in `tasks.md`:

```bash
# Create sub-plan for implementation
IMPL_PLAN=$(plan start --title "Implement: <task-name>" --parent $PLAN_ID | jq -r '.id')

# Add TDD steps
plan step --plan $IMPL_PLAN --title "Write failing test"
plan step --plan $IMPL_PLAN --title "Implement minimal code to pass"
plan step --plan $IMPL_PLAN --title "Refactor and verify"
```

**Implementation principles:**
1. **Test first:** Write failing tests that verify acceptance criteria
2. **Minimal code:** Implement just enough to pass tests
3. **Spec compliance:** Code must satisfy the specification
4. **Refactor:** Clean up while keeping tests green

```bash
plan complete --step <implement-step-id>
```

---

## Phase 5: VALIDATE

**Goal:** Verify implementation against the original specification.

```bash
plan progress --step <validate-step-id> --status in_progress
```

Create `specs/<feature>/validation.md`:

```markdown
# Validation Report: <Feature Name>

## Specification Compliance

### User Story 1: <Story name>
| Acceptance Criterion | Status | Notes |
|---------------------|--------|-------|
| Criterion 1 | ✅ Pass | |
| Criterion 2 | ✅ Pass | |
| Criterion 3 | ⚠️ Partial | <Explanation> |

### User Story 2: <Story name>
...

## Test Coverage
- Unit tests: X passing
- Integration tests: X passing
- E2E tests: X passing

## Outstanding Items
- [ ] <Any remaining work>
```

**Completion check:** All acceptance criteria verified, tests passing.

```bash
plan complete --step <validate-step-id>
plan complete --plan $PLAN_ID
```

---

## File Structure

```
specs/<feature-branch>/
├── spec.md           # Requirements and acceptance criteria
├── plan.md           # Technical implementation strategy
├── tasks.md          # Executable work items
├── validation.md     # Final compliance report
├── research.md       # (optional) Technology analysis
└── contracts/        # (optional) Detailed API specs
    └── api.md
```

---

## Quick Reference

| Phase | Input | Output | Plan Step |
|-------|-------|--------|-----------|
| Specify | User request | spec.md | "Specify: Define requirements" |
| Plan | spec.md | plan.md | "Plan: Design implementation" |
| Tasks | plan.md | tasks.md | "Tasks: Break down work" |
| Implement | tasks.md | Code + Tests | "Implement: Build with TDD" |
| Validate | All specs | validation.md | "Validate: Verify against spec" |

---

## Example Usage

```bash
# User says: "Add user authentication with OAuth"

# Start SDD workflow
PLAN_ID=$(plan start --title "SDD: OAuth Authentication" | jq -r '.id')

# Set up workflow steps
STEP1=$(plan step --plan $PLAN_ID --title "Specify: Define auth requirements" | jq -r '.id')
STEP2=$(plan step --plan $PLAN_ID --title "Plan: Design OAuth integration" | jq -r '.id')
STEP3=$(plan step --plan $PLAN_ID --title "Tasks: Break down auth work" | jq -r '.id')
STEP4=$(plan step --plan $PLAN_ID --title "Implement: Build auth system" | jq -r '.id')
STEP5=$(plan step --plan $PLAN_ID --title "Validate: Verify auth flows" | jq -r '.id')

# Create spec directory
mkdir -p specs/oauth-auth

# Begin specify phase
plan progress --step $STEP1 --status in_progress
# ... create spec.md ...
plan complete --step $STEP1

# Continue through phases...
```
