# 0004 — Pass-through livre de prefixos (decisão provisória)

Status: provisória (revisar se houver necessidade de allowlist/ROA/IRR ou se ocorrer anúncio indevido)
Data: 2026-10-16

## Contexto

Na Fase 1 o usuário escolheu **"Qualquer prefixo (pass-through)"**: a API aceita qualquer CIDR válido, sem validação de posse/allowlist/ROA/IRR. Isso permite anúncio arbitrário (risco de prefix hijack) e transfere toda a proteção para a autenticação da API (API key única).

## Decisão

- Pass-through livre: aceitar qualquer CIDR IPv4 válido (com validação de sintaxe).
- Proteção baseada apenas em **API key única** (`X-API-Key`), sem rate-limit no MVP.

## Alternativas consideradas

- **Allowlist/ROA/IRR:** mais seguro, mas adiciona dependência de bases externas — deixado para backlog.
- **Espelho das rotas do Huawei:** não atende ao objetivo de inventário próprio.

## Consequências

- Risco de anúncio indevido se a API key vazar — mitigar com troca periódica da chave e rede interna.
- **Gatilho para reabrir:** ocorrer anúncio indevido em produção, ou adicionar RPKI/IRR/allowlist ao backlog.
