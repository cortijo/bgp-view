# 0003 — Contrato OpenAPI (decisão provisória, revista)

Status: provisória (revisar se houver múltiplos consumidores automatizados da API, ou se dashboard web sair do backlog)
Data: 2026-10-16

## Contexto

Na entrevista (Fase 2) o usuário optou por **"sem contrato formal"**, em conflito com o padrão da skill, que exige `openapi.yaml` (OpenAPI 3.1) para toda aplicação que expõe API HTTP. Na Fase 6 o usuário pediu para seguir com a árvore que listava `openapi.yaml` como opcional, sem objeção — adotou-se **incluir** o contrato.

## Decisão

- Incluir **`openapi.yaml` (OpenAPI 3.1)** na raiz documentando a API REST.
- **Provisória:** enquanto houver 1 único consumidor (o próprio operador via CLI/scripts), pode-se tratá-lo como best-effort.

## Alternativas consideradas

- **Sem contrato (escolha original da Fase 2):** mais leve, mas quebra o padrão da skill e dificulta clientes automatizados.
- **Contrato gerado automaticamente dos handlers:** ideal, mas exige toolchain extra — considerar no backlog.

## Consequências

- O guardião de specs mantém `openapi.yaml` sincronizado.
- **Gatilho para reabrir:** se dashboard web, cliente SDK ou segundo consumidor automatizado forem adicionados, tornar obrigatório e exigir verificação de sincronismo no CI.
