# AGBG Advertise

Sistema de anúncios de rede via BGP: mantém uma sessão eBGP com MD5 com um roteador Huawei e divulga prefixos com BGP Communities (perfis nomeados ou cruas) para engenharia de tráfego, tudo em container Docker com auditoria estruturada.

## Fonte da verdade
As especificações em ./specs/ são a fonte da verdade. Toda mudança de comportamento começa pela spec.
Carregue cada spec sob demanda a partir do índice.

@specs/INDEX.md
@specs/PLAN.md
@specs/002-organizacao-do-trabalho.md

## Antes de propor mudanças
Consulte ./specs/decisions/ (comece pelo HISTORY.md). Não reabra uma decisão aceita sem um fato novo. Para mudá-la, proponha uma nova ADR.

## Fluxo de trabalho
- Pegue tarefas em ./tasks/ com status "aberta" e sem dependências pendentes.
- Marque a tarefa como "em-andamento" antes de começar.
- Altere apenas os arquivos listados na tarefa e as pastas do seu papel.
- Ao concluir, atenda a Definition of Done, remova o arquivo da tarefa e faça o commit da entrega citando o ID.

## Comandos
- Instalar: `go mod download` / `docker compose build`
- Rodar: `docker compose up -d` (API em `http://localhost:8080`)
- Testar: `make verify:quick` (durante dev) · `make verify:full` (antes de PR)
- Lint e formatação: `gofmt -w .` · `go vet ./...` · `golangci-lint run`

## Estrutura
- ./specs/: especificações (fonte da verdade)
- ./specs/decisions/: registro de decisões (ADRs) e HISTORY.md
- ./tasks/: fila de trabalho (fora do Git)
- ./tmp/: arquivos temporários (fora do Git)
- ./docs/: documentação derivada
- ./database/sqlite/: migrations (goose)
- ./openapi.yaml: contrato OpenAPI 3.1 da API
