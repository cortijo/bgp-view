# 001 — Arquitetura

Status: aprovada
Depende de: 000

## Contexto

Define os componentes, responsabilidades, dono dos dados, integrações externas, autenticação, contratos e o resumo da stack/qualidade definidos nas Fases 3 e 4.

## Estilo arquitetural

**Monólito enxuto** — um único processo Go contém API REST, CLI e worker BGP. Comunicação interna por chamadas de função, sem mensageria nem cache (decisões puladas na Fase 2 por desnecessárias em escala local).

### Subpacotes

```
cmd/bgpadvert/    → main, wiring, flags/env
internal/api/     → handlers HTTP (chi router), middleware auth/audit
internal/bgp/     → GoBGP client wrapper, worker de reconciliação
internal/store/   → SQLite (prefixos, comunities_perfil, sessao_bgp)
internal/audit/   → log de auditoria estruturado (JSON/slog)
internal/cli/     → CLI (usa a mesma API)
```

| Componente | Responsabilidade | Não é dono de |
|---|---|---|
| `api` | Expor endpoints, validar entrada, autenticar, orquestrar chamadas | Persistência, estado BGP |
| `bgp` | Sessão eBGP com Huawei via GoBGP, aplicar/withdraw prefixos, reconciliação | Regras de negócio de negócio |
| `store` | Persistir estado desejado (SQLite), migrations no boot | Estado efêmero BGP |
| `audit` | Registrar cada operação em JSON (retenção 90d) | — |
| `cli` | Interface de terminal (delega à API) | — |

## Entidades de dados e dono

| Entidade | Tabela | Dono | Descrição |
|---|---|---|---|
| Prefixo | `prefixos` | `store` | CIDR, communities (JSON), ativo, criado_em, atualizado_em |
| Perfil de community | `comunities_perfil` | `store` | nome → lista de communities (ex.: `rtbh`, `prepend2`, `prefer-ix`) |
| Sessão BGP | `sessao_bgp` | `store` (config) / `bgp` (estado) | vizinho, ASN local/remoto, MD5 (via env), estado |
| Auditoria | `auditoria` | `audit` | ação, prefixo, communities, usuário, timestamp |

- **Estado desejado:** fonte de verdade é o `store` (SQLite).
- **Estado aplicado:** GoBGP/Huawei é derivado; o worker reconcilia divergência.

## Integração externa

- **Única:** **GoBGP** via gRPC local (in-process, mesma imagem).
  - Falha de conexão/aplicação → worker marca `degraded`, enfileira intenções, **reconciliação** ao recuperar.
  - Sem retry infinito bloqueante; backoff exponencial com limite.
  - `/ready` retorna **503** enquanto `degraded`.
- **Risco:** se o processo GoBGP cair e a sessão com o Huawei cair, anúncios são perdidos no Huawei — mitigado pelo reconcile no boot (reaplica todo estado `ativo` do SQLite).

## Autenticação e autorização

- **API key única** via header `X-API-Key`, valor vindo de env var obrigatória (`BGP_API_KEY`, `${BGP_API_KEY:?}` no compose).
- Papel único: `admin`. Sem usuários, sem papéis.
- Bootstrap: a chave é definida no `.env` no primeiro boot (sem tabela de usuários).
- Falta/chave errada → `401`; sem rate-limit no MVP (backlog).

## Interfaces do usuário final

- API REST (JSON) — port `8080` (configurável, >1024).
- CLI (`bgpadvert`) — delega à API.
- Health: `GET /health`, `GET /ready`; métricas: `GET /metrics`.

## Contratos

- **`openapi.yaml` (OpenAPI 3.1)** na raiz — contrato da API REST. *(Incluído na árvore; decisão provisória registrada em ADR 0003.)*
- Sem `asyncapi.yaml` — não há mensageria (decisão pulada).

## Stack (resumo da Fase 3 — ver ADR 0001)

- Go 1.22+, `net/http` + **chi**, **GoBGP** embutida, **modernc.org/sqlite** (CGO-free).
- Persistência SQLite em volume Docker nomeado; migrations **goose** no boot.
- Dev = prod via Docker Compose local; deploy final em Linux (VPS); build local (sem registry); GitHub privado; sem licença.

## Qualidade e operação (resumo da Fase 4 — ver ADR 0001)

- Testes: unidade + integração (SQLite em memória) + e2e stub GoBGP; cobertura ≥70%.
- Qualidade: `gofmt` + `go vet` + `golangci-lint` + pre-commit hook.
- CI: GitHub Actions (lint+testes+build em PR); deploy manual via docker compose.
- Ambientes: dev local + prod (VPS); sem staging.
- Segredos: `.env.example` versionado, `.env` fora do Git, `${VAR:?}` obrigatório.
- Observabilidade: `slog` JSON + `/health` + `/ready` + `/metrics` (Prometheus).
- Retenção: auditoria 90d, app 14d, expurgo diário.
- Segurança: `govulncheck` + `gosec` + threat model (brute-force de API key, prefix hijack, vazamento de MD5 BGP, DoS por flooding de rotas).
- Versionamento: SemVer + tags + Keep a Changelog.
- Backup: diário do SQLite (`.backup`), retenção 30d, restore mensal.

## Testes de verificação (dois níveis)

- `scripts/verify.sh` com modos `quick` e `full`, expostos como `verify:quick` / `verify:full` no `Makefile`:
  - **quick:** lint dos arquivos alterados + testes relacionados (`git merge-base HEAD main`); roda durante dev e ao concluir tarefa; código de saída ≠0 se falhar.
  - **full:** lint + vet + todos os testes + build; roda sob pedido, antes de PR e no CI.
- Sem banco de testes isolado por template (SQLite em memória nos testes de integração).

## Decisões relacionadas

- `specs/decisions/0001-decisoes-iniciais.md`
- `specs/decisions/0002-arquitetura-monolito-gobgp.md`
- `specs/decisions/0003-sem-contrato-openapi.md` (provisória — revista, contrato incluído)
