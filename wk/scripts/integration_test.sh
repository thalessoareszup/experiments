#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Build wk binary
cd "$ROOT_DIR"
go build -o wk ./cmd/wk

# Create temporary workspace and isolated HOME so wk uses test-local storage
TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

export HOME="$TMP_DIR/home"
mkdir -p "$HOME"

cp "$ROOT_DIR/wk" "$TMP_DIR/"
cd "$TMP_DIR"

WORKFLOW_DIR="$HOME/.local/wk"
WORKFLOW_FILE="$WORKFLOW_DIR/workflow.yaml"
DB_FILE="$WORKFLOW_DIR/runs.db"

mkdir -p "$WORKFLOW_DIR"

# Create a sample workflow with 3 steps
cat > "$WORKFLOW_FILE" <<'EOF'
name: Sample workflow
steps:
  - id: step1
    name: Step 1
    description: First step
  - id: step2
    name: Step 2
    description: Second step
  - id: step3
    name: Step 3
    description: Third step
EOF

echo "[integration] Using temp dir: $TMP_DIR"
echo "[integration] Using HOME: $HOME"
echo "[integration] Workflow file: $WORKFLOW_FILE"

if [[ ! -f "$WORKFLOW_FILE" ]]; then
  echo "[integration] ERROR: workflow.yaml was not created" >&2
  exit 1
fi

echo "[integration] workflow.yaml created successfully"

# 1) status before start should mention no runs
STATUS_BEFORE="$(./wk status || true)"
echo "$STATUS_BEFORE"
if ! grep -q "No runs found" <<<"$STATUS_BEFORE"; then
  echo "[integration] ERROR: initial status did not report 'No runs found'" >&2
  exit 1
fi

echo "[integration] initial status looks correct"

# 2) start should succeed and create SQLite DB
if ! OUTPUT="$(./wk start)"; then
  echo "[integration] ERROR: 'wk start' failed" >&2
  exit 1
fi

echo "$OUTPUT"

if [[ ! -f "$DB_FILE" ]]; then
  echo "[integration] ERROR: SQLite DB was not created by 'wk start'" >&2
  exit 1
fi

echo "[integration] SQLite DB created successfully at $DB_FILE"

# 3) status should show run #1 and step 1/3
STATUS_OUTPUT="$(./wk status)"
echo "$STATUS_OUTPUT"

if ! grep -q "Run #1" <<<"$STATUS_OUTPUT"; then
  echo "[integration] ERROR: status output does not mention Run #1" >&2
  exit 1
fi

if ! grep -q "Status: in_progress" <<<"$STATUS_OUTPUT"; then
  echo "[integration] ERROR: status output does not show Status: in_progress" >&2
  exit 1
fi

if ! grep -q "Step: 1/3" <<<"$STATUS_OUTPUT"; then
  echo "[integration] ERROR: status output does not show expected step 1/3" >&2
  exit 1
fi

echo "[integration] status after start looks correct"

# 4) next should advance to step 2/3
NEXT_OUTPUT1="$(./wk next)"
echo "$NEXT_OUTPUT1"

if ! grep -q "Run #1 advanced to step 2/3" <<<"$NEXT_OUTPUT1"; then
  echo "[integration] ERROR: next did not advance to step 2/3 as expected" >&2
  exit 1
fi

STATUS_OUTPUT2="$(./wk status)"
if ! grep -q "Step: 2/3" <<<"$STATUS_OUTPUT2"; then
  echo "[integration] ERROR: status after next does not show step 2/3" >&2
  exit 1
fi

echo "[integration] next and status after first advance look correct"

# 5) advance to last step
NEXT_OUTPUT2="$(./wk next)"  # should go to 3/3
echo "$NEXT_OUTPUT2"

if ! grep -q "Run #1 advanced to step 3/3" <<<"$NEXT_OUTPUT2"; then
  echo "[integration] ERROR: second next did not advance to step 3/3 as expected" >&2
  exit 1
fi

STATUS_OUTPUT3="$(./wk status)"
if ! grep -q "Step: 3/3" <<<"$STATUS_OUTPUT3"; then
  echo "[integration] ERROR: status after second next does not show step 3/3" >&2
  exit 1
fi

echo "[integration] status after reaching last step looks correct"

# 6) extra next: should report already at last step and completed
NEXT_OUTPUT3="$(./wk next)" || true
echo "$NEXT_OUTPUT3"

if ! grep -q "already at the last step and marked as completed" <<<"$NEXT_OUTPUT3"; then
  echo "[integration] ERROR: extra next did not report already at last step" >&2
  exit 1
fi

STATUS_OUTPUT4="$(./wk status)"
if ! grep -q "Status: completed" <<<"$STATUS_OUTPUT4"; then
  echo "[integration] ERROR: status does not show completed after extra next" >&2
  exit 1
fi

if ! grep -q "Step: 3/3" <<<"$STATUS_OUTPUT4"; then
  echo "[integration] ERROR: status after completion does not show step 3/3" >&2
  exit 1
fi

echo "[integration] next completion behavior looks correct"

echo "[integration] All checks passed"
