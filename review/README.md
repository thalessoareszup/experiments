# review

A tiny viewer for rendering an agent-authored commentary of `git diff` changes.

## Install

```bash
make install
```

Builds the CLI and copies it to `~/.local/bin/review`.

## Usage

### 1. Add the skill to your project

```bash
review skill
```

Creates `skill/review/SKILL.md` which teaches the agent how to generate reviews.

### 2. Ask the agent to review its changes

The agent will:
- Generate a unified diff: `git diff --no-color > patch.diff`
- Create a `review.json` with comments anchored to diff lines

### 3. View the review

```bash
review serve --review path/to/review.json
```

Opens `http://localhost:6767` with the review loaded.

## Example

```bash
make serve
```
