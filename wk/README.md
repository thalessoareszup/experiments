# wk

Ferramenta para gerenciar workflows sequenciais em agentes de IA

## Instalação

```bash
go build -o wk ./cmd/wk
```

Mova `wk` para algum lugar do seu PATH.

## Uso com Agentes de IA

Para usar o `wk` com agentes de IA (como `stk-codegen`), informe o agente sobre a ferramenta e executar `wk onboard`. Por exemplo:

```
You have access to `wk`, a workflow management tool. YOU MUST use `wk` to start any work I give you later.

Run `wk onboard` to learn how to use it.
```

O agente então deverá executar `wk onboard`, e será instruído a usar os comandos `wk start`, `wk status`, e `wk next` para executar o workflow passo a passo.

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

## Usar como SKILL no Claude Code

O `wk` pode ser usado como uma SKILL (habilidade) no Claude Code para agentes autônomos estruturarem seu trabalho em workflows sequenciais.

### Setup

1. Copie a pasta `skill/` para sua pasta `.claude/skills/wk`:

```bash
mkdir -p ~/.claude/skills
cp -r ./skill ~/.claude/skills/wk
```

Ou para uso compartilhado em projeto (dentro do repositório):

```bash
mkdir -p .claude/skills
cp -r ./skill .claude/skills/wk
```

2. Reinicie o Claude Code para que a skill seja reconhecida.

3. Defina seu workflow em `$HOME/.local/wk/workflow.yaml` (veja exemplo acima na seção "Uso com Agentes de IA").

### Como Funciona

Quando um agente tem acesso à skill `wk-workflow-manager`:

1. O agente reconhece que há uma ferramenta de workflow disponível
2. Pode usar comandos como:
   - `wk start` - inicia o workflow
   - `wk status` - verifica o status atual
   - `wk next` - passa para o próximo passo
   - `wk say "mensagem"` - adiciona notas ao passo atual

3. O agente executa o trabalho passo a passo de forma controlada e rastreável

### Exemplo de Uso

Você pode instruir o Claude Code ou outro agente assim:

```
Você tem acesso à skill `wk-workflow-manager` para estruturar seu trabalho.
Use-a para executar o workflow definido em sua configuração.
Comece com `wk start` e execute cada passo sequencialmente usando `wk status`,
`wk say`, e `wk next`.
```

Para mais detalhes, consulte [skill/SKILL.md](skill/SKILL.md).
