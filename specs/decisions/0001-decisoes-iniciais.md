# 0001 — Decisões Iniciais

Status: aceita
Data: 2026-10-16

## Contexto

Decisões iniciais de produto, stack, qualidade e operação definidas durante a entrevista de início de projeto (Fases 1, 3 e 4).

## Decisão

### Produto (Fase 1)
- API REST (JSON) + CLI leve para operação.
- Engenharia de tráfego/mitigação via communities (perfis nomeados + cruas).
- MVP ≤500 prefixos, alteração refletida em ≤5s (p95 <3s), MVP ≤2 semanas.
- Auditoria em log estruturado (JSON) com retenção 90 dias.
- Persistência: SQLite; IPv4 + IPv6 (IPv6 fica no backlog).
- Fluxo operacional: imperativo (cadastra → persiste → aplica → audita).

### Stack (Fase 3)
- Go 1.22+; stdlib net/http + chi; GoBGP embutida; modernc.org/sqlite (CGO-free).
- SQLite em volume Docker; dev = prod (Docker Compose local).
- Deploy final: Docker Compose em Linux (VPS); build local (sem registry); GitHub privado; sem licença formal + docs internas.

### Qualidade e operação (Fase 4)
- Testes unidade+integração+e2e stub, cobertura ≥70%.
- gofmt + go vet + golangci-lint + pre-commit hook.
- GH Actions CI (lint+testes+build) + deploy manual via docker compose; dev+prod (sem staging).
- `.env.example` versionado + `.env` fora do Git; `${VAR:?}` obrigatório para segredos.
- Observabilidade: slog JSON + `/health` + `/ready` + `/metrics` (Prometheus).
- Retenção: auditoria 90d, app 14d, expurgo diário.
- Segurança: govulncheck + gosec + threat model (brute-force, prefix hijack, MD5 leak, DoS).
- SemVer + Keep a Changelog; migrations goose no boot; backup diário SQLite + restore mensal (retenção 30d).

## Alternativas consideradas

- **Linguagem:** Python (FastAPI) — rejeitada por imagem maior e integração GoBGP menos nativa.
- **DB:** Postgres — rejeitado para MVP (complexidade extra; SQLite suficiente até 500 prefixos).
- **CI/CD:** deploy automatizado — adiado para backlog (1 dev, VPS manual).
- **Banco de testes isolado (template DB)** — não se aplica: testes de integração usam SQLite em memória.

## Consequências

- MVP rápido, imagem Docker pequena (binário estático), operação simples.
- Sem HA, sem registry, sem deploy automatizado — aceitável para uso interno; revisar se houver segunda máquina ou falha no VPS.
