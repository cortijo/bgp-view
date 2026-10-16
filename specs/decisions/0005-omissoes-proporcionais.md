# 0005 — Escopo pulado: mensageria, cache, multi-tenancy, compliance (decisões de omissão)

Status: aceita
Data: 2026-10-16

## Contexto

Durante a entrevista alguns tópicos das fases foram deliberadamente pulados por inadequação ao porte do projeto (regra de proporcionalidade da skill).

## Decisão

Registrar explicitamente as omissões:

- **Mensageria/eventos:** pulada — monólito local, chamadas de função bastam.
- **Cache:** pulado — SQLite local com volume, sem latência de rede.
- **Multi-tenancy / rate-limit por cliente:** pulado — 1 operador.
- **Compliance/LGPD:** pulado — infra de rede interna, sem dados pessoais.
- **Contratos AsyncAPI:** não se aplica — não há sistema de eventos.

## Alternativas consideradas

- Introduzir NATS/Valkey desde o início: superdimensionado.

## Consequências

- Se houver mais de 1 consumidor da API, HA ou dados pessoais, reavaliar — abrir nova ADR.
