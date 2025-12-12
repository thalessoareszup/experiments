Este repositório é um fork do [mini-swe-agent](https://github.com/SWE-agent/mini-swe-agent) que suporta agentes [Stackspot](https://stackspot.com/) através da _model class_ `stk`.

### Pré-requisitos

Você precisa configurar as seguintes variáveis de ambiente:

- `STK_CLIENT_ID`: Seu Client ID do Stackspot
- `STK_CLIENT_SECRET`: Seu Client Secret do Stackspot
- `STK_REALM`: O Realm do Stackspot (ex.: `stackspot`)

### Uso

Fazendo o setup usando uv:

```bash
uv venv
source .venv/bin/activate
uv pip install -e .
```

Você pode usar o comando especializado `mini-stk`

```bash
uv run mini-stk --agent-id <seu-agent-id>
```

## Loop de agente

```mermaid
flowchart TD
    Start([Tarefa]) --> Query["1. QUERY<br/>Consulta AI<br/>model.query(messages)"]

    Query --> Parse["2. PARSE<br/>Extrai comando bash<br/>parse_action()"]

    Parse --> Execute["3. EXECUTE<br/>Executa comando<br/>env.execute(action)"]

    Execute --> Observe["4. OBSERVE<br/>Adiciona resultado ao histórico<br/>add_message()"]

    Observe --> Check{"5. CHECK<br/>Tarefa completa?"}

    Check -->|Não| Query
    Check -->|Sim| End([Resultado Final])

    style Start fill:#90EE90
    style End fill:#FFB6C1
    style Query fill:#87CEEB
    style Parse fill:#FFD700
    style Execute fill:#DDA0DD
    style Observe fill:#98FB98
    style Check fill:#F0E68C
```
