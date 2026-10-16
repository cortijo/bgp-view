# 000 — Visão Geral

Status: aprovada
Depende de:

## Contexto

O operador de rede precisa divulgar prefixos via BGP para um roteador Huawei de borda, aplicando **BGP Communities** de forma controlada (engenharia de tráfego, mitigação, preferência por IXP/trânsito). Hoje isso é feito manualmente na CLI do Huawei, o que é lento, sujeito a erro e sem auditoria.

Este sistema automatiza esse fluxo: mantém uma **sessão eBGP com MD5** com o Huawei, aceita comandos via **API REST + CLI**, persiste o estado desejado em **SQLite** e aplica/reconcilia os anúncios no Huawei via **GoBGP** (in-process), tudo em container **Docker**.

## Usuários

| Persona | Perfil |
|---|---|
| Operador de rede / NOC | Único operador humano; usa API/CLI para criar, editar, remover anúncios e aplicar communities. |
| Administrador | Mesmo operador; configura API key, ASN, MD5 e parâmetros da sessão via `.env`/env vars. |

## Objetivos e sucesso

- Reduzir o tempo de aplicação de uma mudança de rota de minutos (CLI Huawei manual) para **segundos** (API/CLI → BGP refletido).
- Garantir **rastreabilidade total** de toda alteração de anúncio/community via log de auditoria estruturado (JSON).
- Eliminar erros de digitação na CLI do Huawei (sistema centraliza e padroniza).

**Métrica de sucesso:** 100% das alterações de anúncio aplicadas via sistema refletidas no Huawei em ≤5s (p95 <3s), com auditoria completa.

## Escopo

### No escopo (MVP)

- Sessão eBGP com MD5 com o Huawei (único vizinho).
- Anúncio, edição e remoção de prefixos IPv4.
- Aplicação de BGP Communities (cruas e perfis nomeados).
- API REST (JSON) + CLI leve + health endpoints.
- Persistência em SQLite (restaurada no boot do container).
- Auditoria em log estruturado (JSON) com retenção de 90 dias.
- Dockerização (compose local/prod), CI no GitHub Actions.
- Métricas Prometheus em `/metrics`.

### Fora de escopo (backlog explícito)

- Suporte a IPv6.
- Múltiplos vizinhos BGP simultâneos.
- Dashboard web de gestão.
- Rate-limit por cliente / multi-tenant.
- Webhooks de notificação.
- RPKI / IRR / validação de posse de prefixo.
- Alta disponibilidade (HA), réplica do SQLite.
- Deploy automatizado (somente manual via docker compose).

## Restrições

- Prazo alvo do MVP: **≤2 semanas** (meta, não bloqueante).
- Orçamento: **0** (uso interno, sem licença formal).
- Equipe: **1 dev + agentes Claude Code**.
- Compliance: sem exigências regulatórias (LGPD não se aplica — infra de rede interna).

## Requisitos não funcionais (verificáveis)

| ID | Requisito |
|---|---|
| RNF-01 | Até **500 prefixos** simultâneos anunciados. |
| RNF-02 | Alteração de anúncio refletida no Huawei em **≤5s** (p95 <3s). |
| RNF-03 | Autenticação por API key (`X-API-Key`); sem chave → HTTP 401. |
| RNF-04 | Logs de auditoria em JSON, retenção **90 dias**; logs de app, **14 dias**; expurgo diário. |
| RNF-05 | Cobertura de testes ≥ **70%**. |
| RNF-06 | `go vet` + `golangci-lint` + `gofmt` limpos em todo PR. |
| RNF-07 | `govulncheck` + `gosec` sem vulnerabilidades críticas em CI. |
| RNF-08 | Health checks: `/health` (liveness) e `/ready` (readiness) respondendo em ≤100ms. |
| RNF-09 | Métricas Prometheus em `/metrics` (prefixos ativos, estado da sessão BGP, latência da API). |
| RNF-10 | Container com limites de CPU/RAM configuráveis; raiz read-only; ports >1024. |
| RNF-11 | Backup diário consistente do SQLite (`.backup`), retenção 30 dias, restore testado mensalmente. |

## Fluxos principais (ponta a ponta)

1. **Anunciar prefixo:** operador envia `POST /v1/prefixes` (ou CLI `bgp announce`) → persiste em SQLite → worker BGP aplica no GoBGP → Huawei recebe UPDATE → reflete em ≤5s → auditoria registrada.
2. **Aplicar community/perfil:** `PATCH /v1/prefixes/{cidr}` com communities ou perfil → persiste → worker aplica → reflete → audita.
3. **Remover anúncio:** `DELETE /v1/prefixes/{cidr}` → persiste (marca inativo) → worker withdraw → audita.
4. **Falha do GoBGP:** worker entra em `degraded`, enfileira intenções, reconcilia ao recuperar; `/ready` retorna 503 durante degraded.

## Decisões relacionadas

- `specs/decisions/0001-decisoes-iniciais.md`
- `specs/decisions/0004-passthrough-de-prefixos.md` (provisória)
