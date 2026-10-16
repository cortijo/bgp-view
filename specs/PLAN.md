# Plano de Implementação

> Cada fase é uma meta verificável. Tarefas da primeira fase estão em `./tasks/`.

## Fase 1 — Fundação (sequencial)

- **Objetivo:** Repo Git inicial, scaffolding Go, Docker Compose funcional, migrations e schema SQLite, CI mínimo.
- **Specs envolvidas:** 000, 001, 002.
- **Dependências:** nenhuma (primeira fase).
- **Paralelismo:** nenhum — sequencial por definição.
- **Papel:** implementador (guardião só confirma specs).
- **Critério de saída:** `docker compose up` sobe o container, `/health` responde `200`, migrations aplicam sem erro, PR com CI verde e 1 review.

## Fase 2 — Sessão BGP (parcialmente paralela)

- **Objetivo:** Estabelecer sessão eBGP com MD5, worker de reconciliação, estado degraded e métricas de sessão.
- **Specs envolvidas:** 003 (e dependências 000/001).
- **Dependências:** Fase 1 concluída.
- **Paralelismo:** API de sessão (`GET /v1/bgp/session`) e worker podem correr em paralelo com a camada de store de prefixos.
- **Papel:** implementador; revisor valida testes de integração.
- **Critério de saída:** Sessão com Huawei estabelecida; boot reconcilia estado; `/ready` 503 em degraded; `verify:quick` verde.

## Fase 3 — Anúncio de Prefixos (paralela)

- **Objetivo:** CRUD de prefixos, validação, idempotência, aplicação no GoBGP, auditoria.
- **Specs envolvidas:** 004.
- **Dependências:** Fase 2.
- **Paralelismo:** API de prefixos, store e worker podem correr em paralelo.
- **Papel:** implementador (backend); implementador (testes) em paralelo.
- **Critério de saída:** CRUD completo funcional, alteração refletida ≤5s, auditoria registrada, `verify:quick` verde, cobertura ≥70%.

## Fase 4 — Communities e Perfis (paralela após 004)

- **Objetivo:** CRUD de perfis, expansão de perfil, reconciliação em cascata, seed de perfis padrão.
- **Specs envolvidas:** 005.
- **Dependências:** Fase 3.
- **Paralelismo:** API de perfis e lógica de expansão podem correr em paralelo.
- **Papel:** implementador.
- **Critério de saída:** Perfis CRUD + expansão ≤5s + seed + `verify:quick` verde.

## Fase 5 — CLI, Observabilidade e Hardening

- **Objetivo:** CLI leve, `/metrics` Prometheus, expurgo de logs, backup diário do SQLite, `gosec`/`govulncheck` no CI.
- **Specs envolvidas:** 001, 002.
- **Dependências:** Fases 2–4.
- **Paralelismo:** CLI, métricas e backup podem correr em paralelo.
- **Papel:** implementador; revisor valida documentação e CHANGELOG.
- **Critério de saída:** CLI funcional, `/metrics` expondo métricas, backup diário rodando, CI com `govulncheck`+`gosec`, docs em `docs/` atualizadas.

## Fase 6 — Release MVP

- **Objetivo:** Tag `v0.1.0`, CHANGELOG completo, deploy manual no VPS (docker compose), smoke test de produção.
- **Specs envolvidas:** todas.
- **Dependências:** Fases 1–5.
- **Paralelismo:** nenhum — sequencial (release).
- **Papel:** implementador (build), revisor (CHANGELOG), guardião (validação final de specs).
- **Critério de saída:** Tag publicada no GitHub privado, `docker compose up -d` no VPS com sessão `Established`, smoke test de anúncio/reflexão ≤5s, CHANGELOG completo.

## Backlog (pós-MVP)

- IPv6, múltiplos vizinhos, dashboard web, rate-limit por cliente, webhooks, RPKI/IRR, HA, deploy automatizado (GH Actions + GHCR).
