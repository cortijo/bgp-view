# 002 — Organização do Trabalho

Status: aprovada
Depende de: 000

## Ferramenta de agente

**Claude Code** — gera `CLAUDE.md` (apenas `@AGENTS.md`), `AGENTS.md` e definições em `.claude/agents/`.

## Papéis

| Papéis | Responsabilidade | Pastas/arquivos que pode alterar |
|---|---|---|
| **guardiao-specs** | Guardião das specs e contratos; ADRs; planos | `./specs/`, `./specs/decisions/`, `openapi.yaml`, `AGENTS.md`, `PLAN.md` |
| **implementador** | Implementar tarefas (código, testes, Docker, CI) | `./cmd/`, `./internal/`, `./scripts/`, `./docker-compose.yml`, `.env.example`, `./database/` |
| **revisor** | Revisar PRs, executar verificação, redigir docs | `./docs/`, `./CHANGELOG.md`; **não escreve código** |

- **Somente o guardião** altera `./specs/`, `./specs/decisions/`, `openapi.yaml` e migrations (via aprovação).
- Implementador e revisor podem **propor** mudança de spec abrindo ADR em `./specs/decisions/` — mas a aplicação é do guardião.

## Paralelismo

- **Fase 1 (fundação):** sequencial (repo, scaffolding, Docker, migrations).
- **Fases 2+:** paralelizáveis — API, worker BGP, CLI e testes podem rodar em paralelo após as specs da fase serem aprovadas.
- **Sequenciais (nunca paralelas):** qualquer mudança de `openapi.yaml`, schema SQLite (`database/sqlite/`) e ADRs.

## Isolamento e merge

- Cada agente trabalha em um **worktree Git próprio** (`git worktree add`) por tarefa.
- Merge via **PR com rebase** antes do merge; `main` é trunk.
- Conflito de merge não resolvido → escalonar ao humano.

## Git

- **Modelo:** trunk-based, branch por tarefa: `feat/task-001-...`, `fix/task-002-...`.
- **Commits:** Conventional Commits (`feat:`, `fix:`, `chore:`, `docs:`) citando o ID da tarefa (ex.: `feat(api): ... (task 001)`).
- **PR:** obrigatório, 1 review aprovado (revisor), CI verde antes do merge.

## Fluxo das tarefas

1. Tarefas ficam em `./tasks/` (fora do Git), formato do template da skill.
2. Antes de começar: status `em-andamento` + preencher `responsavel`.
3. Nenhum agente pega tarefa `em-andamento`.
4. Ao concluir: atender DoD → **remover** o arquivo da tarefa → commit citando o ID.
5. Tarefa descartada também é removida; se representar decisão → nova ADR.

## Definition of Done (DoD)

Uma tarefa está concluída quando:

- [ ] Código implementado.
- [ ] Testes do nível acordado escritos e passando (`verify:quick` verde).
- [ ] `gofmt` + `go vet` + `golangci-lint` limpos.
- [ ] Spec/`openapi.yaml`/migrations atualizados, se mudaram.
- [ ] `./docs/` atualizado, quando aplicável.
- [ ] `CHANGELOG.md` recebeu resumo compactado (`## [Não lançado]`).
- [ ] Arquivo da tarefa removido.

## Escalonamento (parar e perguntar ao humano)

- Conflito com uma spec (código exige comportamento divergente da spec).
- Decisão irreversível (mudança de schema SQLite, mudança de contrato `openapi.yaml`, mudança na sessão BGP).
- **2 falhas** seguidas na mesma tarefa.
- Qualquer tentativa de reabrir ADR aceita sem fato novo.

## Verificação

- Durante dev e ao fim de cada tarefa: `make verify:quick`.
- Antes de PR e no CI: `make verify:full`.
- Ao terminar, sempre informar qual nível rodou e o resultado.
- Subagentes seguem a mesma regra: rápido por padrão, completo só quando delegado pedir.

## Decisões relacionadas

- `specs/decisions/0001-decisoes-iniciais.md`
