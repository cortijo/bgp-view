# Comandos do projeto (AGENTS.md)
# Observação: alvos com ':' (verify:quick/verify:full) precisam de escape '\:' no make.
.PHONY: run build test lint fmt verify-quick verify-full verify-quick verify-full

PATH := $(PATH):/usr/local/go/bin
export PATH

build:
	go build -o bin/bgpadvert ./cmd/bgpadvert

run:
	mkdir -p tmp
	BGP_DB_PATH=./tmp/bgpadvert.db go run ./cmd/bgpadvert

test:
	go test ./...

fmt:
	gofmt -w .

lint:
	golangci-lint run

# Nível rápido: durante desenvolvimento e ao concluir cada tarefa
verify-quick:
	@./scripts/verify.sh quick

# Nível completo: sob pedido, antes de PR e no CI (para na 1ª falha)
verify-full:
	@./scripts/verify.sh full

# Aliases documentados na spec: make verify\:quick / make verify\:full
verify\:quick: verify-quick
verify\:full: verify-full
