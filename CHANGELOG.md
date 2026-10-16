# Changelog

Todo formato notável deste projeto é documentado neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/),
e este projeto adere a [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [Não lançado]

### Adicionado
- **Tarefa 003:** verificação em dois níveis (`scripts/verify.sh quick|full`, `make verify\:quick` / `verify\:full`), `.golangci.yml` e hook de pre-commit — `verify:full` aprovado (gofmt, vet, lint, testes, cobertura 71,2%, build).
- **Tarefa 002:** `Dockerfile` multi-stage (imagem mínima, non-root, healthcheck), `docker-compose.yml` (read-only, tmpfs, `${VAR:?}` para segredos, limites CPU/RAM, porta >1024, ulimits), `.env.example` versionado, CI GitHub Actions (gofmt, vet, testes com corte de 70%, build, imagem) — validado com container no ar (`/health` 200, healthcheck `healthy`).
- **Tarefa 001:** scaffold Go (módulo, `cmd/bgpadvert`, `internal/{api,bgp,store,audit,cli,config}`), schema SQLite com migrations versionadas (`0001`), middleware de API key e worker BGP com reconciliação/backoff — `go build`, `go vet`, `gofmt` e testes verdes (cobertura 70.2%).
- Estrutura inicial de specs (`./specs/`), registro de decisões (`./specs/decisions/`) e plano (`specs/PLAN.md`).
- Definições de papéis de agente (guardião, implementador, revisor) e regras de trabalho (`AGENTS.md`).
- Contrato `openapi.yaml` (OpenAPI 3.1) da API de anúncios BGP.
- Compose e `.env.example` para containerização (GoBGP in-process + SQLite).
- Fila de tarefas inicial da Fase 1 (`./tasks/`).
