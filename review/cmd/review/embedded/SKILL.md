---
name: review
version: 1
description: >
  A structured code-review artifact that explains git diff changes with ordered
  comments anchored to files and diff line ranges (GitHub-like).
---

# Review (skill)

This skill defines a **machine-readable review artifact** that an agent can produce after making changes.

## When to use

Use `review` when you want to:
- summarize *what changed* in a commit/PR in a review-friendly way
- attach explanations to specific parts of the diff (file + diff line range)
- provide an ordered walkthrough (top-to-bottom narrative)

## Output format: `review.v1.json`

The agent should output a single JSON file (UTF-8) that follows this structure.

### Top-level shape

```jsonc
{
  "$schema": "https://example.invalid/review/v1",
  "version": 1,
  "title": "Short title of the change",
  "createdAt": "2025-12-24T12:34:56Z",

  "base": { "ref": "main", "commit": "<sha>" },
  "head": { "ref": "feature-branch", "commit": "<sha>" },

  "diff": {
    "tool": "git",
    "command": "git diff --no-color <base>..<head> > patch.diff",
    "patch": "patch.diff"
  },

  "comments": [
    {
      "id": "c1",
      "order": 1,
      "severity": "note",
      "message": "Explain what changed and why.",
      "anchor": {
        "file": "path/to/file.ts",
        "diffLineStart": 42,
        "diffLineEnd": 55
      }
    }
  ]
}
```

### Fields

- `version` (number): must be `1`.
- `title` (string): short human-readable title.
- `createdAt` (string): ISO-8601 timestamp.
- `base`, `head` (object): optional metadata for context.
- `diff` (object):
  - `patch` (string, required): relative path to a local unified diff file (e.g. `patch.diff`).
  - `command` (string, optional): how it was produced.
  - `unified` (string, optional): inline unified diff text (fallback; useful when not using a separate patch file).
- `comments` (array, required): ordered review comments.

### Comment fields

- `id` (string): unique within the file.
- `order` (number): explicit ordering for the walkthrough. The viewer sorts by `order` then `id`.
- `severity` (string): one of `note | nit | suggestion | warning`.
- `message` (string): Markdown-like text (viewer may render as plain text).
- `anchor` (object): where the comment points.

### Anchor fields (GitHub-like)

Anchors reference **diff line numbers**, not source line numbers.

- `file` (string): file path as it appears after `+++ b/<path>`.
- `diffLineStart` (number): 1-based line number within that file's rendered diff.
- `diffLineEnd` (number): inclusive end (>= start). Use the same as start for a single line.

Notes:
- `diffLineStart` counts only the lines in that file section (starting at the first `@@ ... @@` hunk line). The viewer treats the file diff as a continuous list of lines (including context/add/remove lines).
- If the file appears multiple times (rename/copy), anchor to the final `b/<path>`.

## Guidance for agents

1. Generate the diff with no color codes and save it to a patch file next to the JSON:

```bash
git diff --no-color > patch.diff
```

2. Create `review.v1.json` that references the patch:

```jsonc
{
  "version": 1,
  "diff": { "patch": "patch.diff" },
  "comments": [/* ... */]
}
```

3. Write comments in the order you want the human to read them.
4. Prefer small, targeted anchors (110 lines).
5. If you want to comment on a whole file, anchor to the first 15 diff lines.

## Example

See `examples/review.sample.json`.
