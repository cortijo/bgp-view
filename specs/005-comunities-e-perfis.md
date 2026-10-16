# 005 — Communities e Perfis

Status: rascunho
Depende de: 000, 001, 004

## Contexto

As BGP Communities permitem ao Huawei (e upstreams) manipular o tráfego (prepend, blackhole, preferência por IXP/trânsito). O sistema oferece duas formas de aplicar communities: **cruas** (operador digita `ASN:VALUE`) e **perfis nomeados** (abstração reutilizável).

## Requisitos

- RF-01: Aceitar community standard no formato `ASN:VALUE` (ex.: `65000:666`).
- RF-02: Aceitar extended communities (ex.: `rt:65000:100`, `lb:0`) com validação de sintaxe.
- RF-03: Gerenciar **perfis nomeados** (CRUD): `nome` → lista de communities.
- RF-04: Aplicar um perfil a um prefixo (expande para a lista de communities no worker).
- RF-05: Perfil alterado → re-aplicar a todos os prefixos que o referenciam (reconciliação).
- RF-06: Não permitir community malformada → `400`.
- RNF-01: Expansão de perfil refletida em ≤5s em todos os prefixos afetados.
- RNF-02: Até 50 perfis simultâneos (`BGP_MAX_PROFILES`).

## Comportamento

- **Perfil usado em prefixo:** ao alterar/remover perfil → lista de prefixos afetados → worker re-aplica cada um → auditoria por prefixo.
- **Remoção de perfil:** se referenciado por prefixo → `409` (remover referência antes).
- **Perfis padrão pré-carregados (seed):** `rtbh` (`65000:666`), `prepend2`, `prefer-ix` — valores default configuráveis via env na inicialização.

## Interfaces e contratos

- `openapi.yaml` — recursos `/v1/profiles`, `/v1/profiles/{nome}`, e campo `perfil` em `/v1/prefixes`.

## Dados

| Tabela | Campos |
|---|---|
| `comunities_perfil` | `nome` (PK), `comunities` (JSON array), `criado_em`, `atualizado_em` |

- Relacionamento: `prefixos.perfil` → `comunities_perfil.nome` (FK, 0..1).
- Dono: `store`. Retenção: indefinida.

## Critérios de aceite

- [ ] CRUD de perfis via API e CLI.
- [ ] Expansão de perfil aplicada em ≤5s aos prefixos referenciados.
- [ ] Validação de sintaxe de community (standard e extended).
- [ ] Seed de perfis padrão no boot.
- [ ] `verify:quick` verde, cobertura do pacote ≥70%.

## Fora de escopo

- Regras condicionais por prefix (ex.: com base em IGP cost), communities dinâmicas por rota, BGP large communities no MVP (backlog).

## Decisões relacionadas

- `specs/decisions/0001-decisoes-iniciais.md`
