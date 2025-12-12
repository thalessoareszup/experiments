# wk

Ferramenta para gerenciar workflows sequenciais em agentes de IA

## Instalação

```bash
go build -o wk ./cmd/wk
```

Mova `wk` para algum lugar do seu PATH.

## Uso com Agentes de IA

Para usar o `wk` com agentes de IA (como `stk-codegen`), informe o agente sobre a ferramenta e peça para executar:

```bash
wk onboard
```

O agente então será instruído a usar os comandos `wk start`, `wk status`, e `wk next` para executar o workflow passo a passo.

Defina o workflow em `$HOME/.local/wk/workflow.yaml`. Exemplo:

```yaml
name: Implementar feature
steps:
   - id: explore
     name: Explorar projeto
     description: Explore o projeto para descobrir qual a linguagem e ferramenta de build, depois focando nos pacotes de domínio para entender as principais lógicas de negócio.
   - id: compile
     name: Compilar projeto
     description: Tente compilar o projeto baseado na ferramenta de build e documentação existente
     requires-confirmation: true
```

## Interface Web

Para monitorar e confirmar etapas:

```bash
go run ./cmd/web
```

Acesse em `http://localhost:8080`
