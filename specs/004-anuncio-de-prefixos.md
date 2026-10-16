# 004 — Anúncio de Prefixos

Status: rascunho
Depende de: 000, 001, 003

## Contexto

Gerencia o ciclo de vida dos prefixos anunciados: criação, edição, remoção e listagem, com aplicação imediata via sessão BGP.

## Requisitos

- RF-01: Criar prefixo (CIDR IPv4 + lista de communities) via API/CLI.
- RF-02: Editar communities de um prefixo existente (substituição total).
- RF-03: Remover (withdraw) um prefixo de forma atômica.
- RF-04: Listar prefixos com filtro por `ativo`, `community` e paginação.
- RF-05: Pass-through livre: aceitar qualquer CIDR válido IPv4 (sem validação de posse/allowlist).
- RF-06: Validar sintaxe de CIDR (IPv4) e rejeitar inválidos com `400`.
- RF-07: Persistir estado em SQLite e restaurar no boot.
- RNF-01: Alteração refletida no Huawei em ≤5s (p95 <3s).
- RNF-02: Até 500 prefixos simultâneos.
- RNF-03: Operações idempotentes (criar prefixo já existente → `409`).

## Comportamento

- **Fluxo:** `POST /v1/prefixes` → valida CIDR → grava em `prefixos` → worker BGP aplica → auditoria → `201`.
- **Remoção:** `DELETE /v1/prefixes/{cidr}` → marca `ativo=false` → worker withdraw → auditoria → `204`.
- **Idempotência:** recriar CIDR já existente e ativo → `409 Conflict`.
- **Erro de sessão:** se GoBGP `degraded`, aceita a escrita no `store` mas responde `202 Accepted` com `queued: true`; reconcilia depois.
- **Remoção físico:** `DELETE` marca inativo (soft delete) — histórico mantido para auditoria.

## Interfaces e contratos

- `openapi.yaml` — recursos `/v1/prefixes`, `/v1/prefixes/{cidr}`.
- Env: `BGP_MAX_PREFIXES` (default 500) — exceder → `422`.

## Dados

| Campo | Tipo | Notas |
|---|---|---|
| `cidr` | TEXT PK | ex.: `203.0.113.0/24` |
| `comunities` | JSON (array) | lista de communities no formato `ASN:VALUE` ou extended |
| `perfil` | TEXT NULL | FK opcional para `comunities_perfil.nome` |
| `ativo` | BOOLEAN | soft delete |
| `criado_em` / `atualizado_em` | TIMESTAMP | UTC |

- Dono: `store`. Retenção: indefinida (auditoria cobre 90d; prefixos inativos mantidos).

## Critérios de aceite

- [ ] CRUD completo via API e CLI.
- [ ] Validação de CIDR e limite de 500 prefixos.
- [ ] Alteração refletida ≤5s em sessão ativa.
- [ ] Estado restaurado no boot.
- [ ] `verify:quick` verde, cobertura de `store`+`api` ≥70%.

## Fora de escopo

- IPv6, validação de posse/ROA/IRR, import em massa via CSV/YAML, bulk withdraw (backlog).

## Decisões relacionadas

- `specs/decisions/0004-passthrough-de-prefixos.md` (provisória)
