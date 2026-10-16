# 003 — Sessão BGP

Status: rascunho
Depende de: 000, 001

## Contexto

Estabelece e mantém uma sessão **eBGP com MD5** entre o sistema (GoBGP in-process) e o roteador Huawei configurado via env vars. É a fundação para todos os anúncios.

## Requisitos

- RF-01: Estabelecer sessão eBGP com MD5 com o Huawei (ASN local e remoto via env).
- RF-02: Manter keepalive/holdtime padrão (60s/180s) — configurável via env.
- RF-03: Reaplicar todo estado `ativo` do `store` ao subir (reconciliação de boot).
- RF-04: Entrar em estado `degraded` se a sessão cair; `/ready` retorna 503.
- RF-05: Reconciliar divergências (estado desejado ≠ estado aplicado) com backoff exponencial.
- RNF-01: Sessão deve ficar `Established` em ≤30s após boot (quando o Huawei responder).
- RNF-02: Estado da sessão exposto em `/metrics` (`bgp_session_state`).

## Comportamento

- **Boot:** inicia GoBGP → conecta ao Huawei → lê `prefixos` ativos do SQLite → aplica todos → estado `Established`.
- **Perda de sessão:** worker entra em `degraded`, enfileira intenções, tenta reconnect com backoff (1s, 2s, 4s... máx 60s).
- **Recuperação:** reconciliação completa (aplica divergências entre desejado e aplicado).
- **Timeout de boot:** se o Huawei não responder em 30s, sistema continua em `degraded` (não crasha).

## Interfaces e contratos

- Env vars obrigatórias: `BGP_LOCAL_ASN`, `BGP_NEIGHBOR_ADDR`, `BGP_NEIGHBOR_ASN`, `BGP_MD5`.
- Config via `./internal/bgp/`.
- `openapi.yaml` — endpoint `GET /v1/bgp/session` retorna estado.

## Dados

- Tabela `sessao_bgp` (somente leitura da config; estado é efêmero).
- Dono: `bgp` (estado), `store` (config em env).

## Critérios de aceite

- [ ] Sessão eBGP com MD5 estabelecida com Huawei.
- [ ] Boot reconcilia prefixos ativos do SQLite.
- [ ] `degraded` + `/ready` 503 em falha de sessão.
- [ ] Backoff exponencial configurável.
- [ ] `verify:quick` verde, cobertura do pacote `bgp` ≥70%.

## Fora de escopo

- Múltiplos vizinhos, iBGP, eBGP multihop, IPv6 (backlog).

## Decisões relacionadas

- `specs/decisions/0002-arquitetura-monolito-gobgp.md`
