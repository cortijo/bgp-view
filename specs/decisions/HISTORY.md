# Histórico de Decisões (HISTORY)

> Índice das ADRs em `./specs/decisions/`. Consulte antes de propor mudança de stack, arquitetura ou processo.

| ADR | Título | Status | Conteúdo |
|---|---|---|---|
| `0001-decisoes-iniciais.md` | Decisões Iniciais | aceita | Produto (Fase 1), stack (Fase 3), qualidade/operacão (Fase 4) em bloco. |
| `0002-arquitetura-monolito-gobgp.md` | Arquitetura: Monólito enxuto com GoBGP in-process | aceita | Monólito Go, GoBGP in-process, SQLite como fonte do estado desejado, reconciliação/degraded. |
| `0003-sem-contrato-openapi.md` | Contrato OpenAPI (provisória) | provisória | Incluir openapi.yaml (revisão do "sem contrato" da Fase 2); gatilho: múltiplos consumidores/dashboard. |
| `0004-passthrough-de-prefixos.md` | Pass-through livre de prefixos (provisória) | provisória | Sem validação de posse/ROA/IRR; proteção só por API key; gatilho: anúncio indevido ou backlog RPKI. |
| `0005-omissoes-proporcionais.md` | Escopo pulado: mensageria, cache, multi-tenancy, compliance | aceita | Registros das omissões por proporcionalidade ao porte. |

## Política

- Uma decisão aceita só é reaberta diante de um fato novo: criar nova ADR substituindo a anterior e marcá-la como "Substituída por nnnn".
- Toda ADR criada ou substituída deve ser registrada nesta tabela.
