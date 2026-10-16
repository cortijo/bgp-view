# 0002 — Arquitetura: Monólito enxuto com GoBGP in-process

Status: aceita
Data: 2026-10-16

## Contexto

Precisamos manter uma sessão eBGP com o Huawei, aplicar/withdraw prefixos com communities e expor uma API. O porte do MVP (≤500 prefixos, 1 vizinho, 1 dev) não justifica microserviços.

## Decisão

- **Monólito enxuto em Go** com subpacotes `api/`, `bgp/`, `store/`, `audit/`, `cli/`.
- **GoBGP in-process** (biblioteca embutida, gRPC local) como único demônio BGP.
- Comunicação interna por chamadas de função; **sem mensageria, sem cache, sem event bus**.
- Falha do GoBGP → worker `degraded` + fila de reconciliação + backoff exponencial; `/ready` 503.
- Persistência em **SQLite** (fonte de verdade do estado desejado); estado aplicado é derivado e reconciliado.

## Alternativas consideradas

- **Dois serviços (API + daemon BGP sidecar):** mais complexidade de deploy e contrato interno para ganho irrelevante no porte atual.
- **Microserviços desde o início:** superdimensionado para MVP de 1 dev.
- **FRR / ExaBGP / BIRD:** FRR mais pesado; ExaBGP em manutenção reduzida e menos programático; BIRD focado em dumping, não anúncio dinâmico. GoBGP oferece API rica de communities + fácil de embutir.
- **SSH no Huawei (config push):** mais frágil, sem reconciliação nativa, sem estado BGP observável.

## Consequências

- Operação simples: um container, um binário.
- Escala limitada (1 vizinho, ~500 prefixos) — revisar se houver múltiplos vizinhos ou HA.
- Acoplamento GoBGP/Huawei depende de sessão única — mitigado por reconcile no boot.
