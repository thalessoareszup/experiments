This repository is a fork of [mini-swe-agent](https://github.com/SWE-agent/mini-swe-agent) that supports [Stackspot](https://stackspot.com/) agents via the `stk` model class.

### Prerequisites

You need to set the following environment variables:

- `STK_CLIENT_ID`: Your Stackspot Client ID
- `STK_CLIENT_SECRET`: Your Stackspot Client Secret
- `STK_REALM`: The Stackspot Realm (e.g., `stackspot`)
- `STK_AGENT_ID`: The Agent ID (Slug) you want to use (optional if passed via CLI)

### Usage

You can use the specialized `stk-mini` command

```bash
stk-mini --model <your-agent-id>
```

The `stk-mini` wrapper treats the `--model` argument as the `agent_id`.

