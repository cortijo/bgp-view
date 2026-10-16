# Índice de Especificações

> Fonte da verdade do projeto. Carregue cada spec sob demanda a partir desta tabela.

| Arquivo | Título | Escopo | Depende de | Status |
|---|---|---|---|---|
| `specs/000-visao-geral.md` | Visão Geral | Propósito, personas, escopo MVP/backlog, restrições, RNFs, fluxos | — | aprovada |
| `specs/001-arquitetura.md` | Arquitetura | Componentes, donos de dados, integrações, auth, contratos, stack e qualidade | 000 | aprovada |
| `specs/002-organizacao-do-trabalho.md` | Organização do Trabalho | Papéis, propriedade, paralelismo, Git, DoD, escalonamento | 000 | aprovada |
| `specs/003-sessao-bgp.md` | Sessão BGP | eBGP com MD5, reconciliação, estado degraded | 000, 001 | rascunho |
| `specs/004-anuncio-de-prefixos.md` | Anúncio de Prefixos | CRUD de prefixos, validação, idempotência, limites | 000, 001, 003 | rascunho |
| `specs/005-comunities-e-perfis.md` | Communities e Perfis | Community standard/extended, perfis nomeados, seed | 000, 001, 004 | rascunho |

## Contratos

| Contrato | Caminho | Status |
|---|---|---|
| OpenAPI 3.1 | `openapi.yaml` | rascunho (decisão provisória — ADR 0003) |

## Registro de decisões

Ver `specs/decisions/HISTORY.md`.
