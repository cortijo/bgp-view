#!/usr/bin/env bash
# scripts/verify.sh — verificação em dois níveis (specs/001-arquitetura.md)
#   quick: lint dos alterados + testes relacionados desde merge-base com main
#   full : gofmt + vet + golangci-lint + todos os testes + build (para na 1ª falha)
# Uso: scripts/verify.sh [quick|full] [-- <testes explícitos...>]
set -uo pipefail
export PATH="$PATH:/usr/local/go/bin"
MODE="${1:-quick}"; shift || true
# Separa os testes explícitos: ./scripts/verify.sh quick -- <go test args...>
while [ "${1:-}" = "--" ]; do shift; done
START=$(date +%s)
FAIL=0
step() { printf '  %-46s %s %s\n' "$1" "$2" "${3:-}"; }

base() {
  local b
  b=$(git rev-parse --verify main >/dev/null 2>&1 && echo main || echo HEAD)
  git merge-base HEAD "$b" 2>/dev/null || echo ""
}

changed_files() {
  # Sem Git (repo ainda não inicializado) ou sem commits: considera todo o tree.
  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    find cmd internal -name '*.go' 2>/dev/null
    return
  fi
  local b; b=$(base)
  if [ -z "$b" ] || [ "$(git rev-parse HEAD)" = "$(git rev-parse "$b")" ]; then
    find cmd internal -name '*.go' 2>/dev/null
  else
    git diff --name-only --diff-filter=ACMR "$b" -- '*.go' 2>/dev/null
  fi
}

echo "== verify:$MODE =="

# 1) gofmt
if [ "$MODE" = "quick" ]; then
  files=$(changed_files)
  if [ -z "$files" ]; then step "gofmt (arquivos alterados)" "✓" "sem alterações"
  else
    out=$(echo "$files" | xargs gofmt -l 2>/dev/null)
    if [ -z "$out" ]; then step "gofmt (arquivos alterados)" "✓" ""; else step "gofmt" "✗" "desformatados: $out"; FAIL=1; fi
  fi
else
  out=$(gofmt -l . 2>/dev/null)
  if [ -z "$out" ]; then step "gofmt" "✓" ""; else step "gofmt" "✗" "desformatados: $out"; FAIL=1; fi
fi

# 2) go vet
if [ "$MODE" = "full" ]; then
  if go vet ./... >/tmp/vet.out 2>&1; then step "go vet" "✓" ""; else step "go vet" "✗" "$(tail -3 /tmp/vet.out)"; echo "== RESULTADO: ✗ reprovado =="; exit 1; fi
else
  if go vet ./internal/... ./cmd/... >/tmp/vet.out 2>&1; then step "go vet" "✓" ""; else step "go vet" "✗" "$(tail -3 /tmp/vet.out)"; FAIL=1; fi
fi

# 3) golangci-lint
if command -v golangci-lint >/dev/null 2>&1; then
  if [ "$MODE" = "quick" ]; then
    files=$(changed_files)
    if [ -z "$files" ]; then
      step "golangci-lint (alterados)" "pulado" "sem alterações"
    else
      # golangci-lint exige arquivos de um mesmo diretório por execução.
      dirs=$(echo "$files" | xargs -n1 dirname 2>/dev/null | sort -u)
      lint_ok=1
      lint_msg=""
      for d in $dirs; do
        if ! golangci-lint run --fast "$d"/... >/tmp/lint.out 2>&1; then
          lint_ok=0
          lint_msg="$d: $(tail -2 /tmp/lint.out | head -1)"
          break
        fi
      done
      if [ "$lint_ok" = "1" ]; then step "golangci-lint (alterados)" "✓" "$dirs"
      else step "golangci-lint" "✗" "$lint_msg"; FAIL=1; fi
    fi
  else
    if golangci-lint run >/tmp/lint.out 2>&1; then step "golangci-lint" "✓" ""; else step "golangci-lint" "✗" "$(tail -3 /tmp/lint.out)"; echo "== RESULTADO: ✗ reprovado =="; exit 1; fi
  fi
else
  step "golangci-lint" "pulado" "não instalado"
fi

# 4) testes
explicit=("$@")
if [ "${#explicit[@]}" -gt 0 ]; then
  if go test "${explicit[@]}" >/tmp/test.out 2>&1; then step "testes (explícitos)" "✓" "${explicit[*]}"; else step "testes (explícitos)" "✗" "$(tail -5 /tmp/test.out)"; FAIL=1; fi
elif [ "$MODE" = "quick" ]; then
  files=$(changed_files); pkgs=""
  for f in $files; do
    d=$(dirname "$f")
    case "$d" in internal/*|cmd/*) pkgs="$pkgs ./$d" ;; esac
  done
  pkgs=$(echo "$pkgs" | tr ' ' '\n' | sort -u | grep -v '^$' | tr '\n' ' ')
  if [ -z "$pkgs" ]; then step "testes (relacionados)" "pulado" "sem pacotes alterados"
  elif go test $pkgs >/tmp/test.out 2>&1; then step "testes (relacionados)" "✓" "$pkgs"
  else step "testes (relacionados)" "✗" "$(tail -5 /tmp/test.out)"; FAIL=1; fi
else
  if go test ./... -coverprofile=/tmp/coverage.out >/tmp/test.out 2>&1; then
    cov=$(go tool cover -func=/tmp/coverage.out | tail -1 | awk '{print $NF}' | tr -d '%')
    step "testes (todos)" "✓" "cobertura ${cov}%"
    awk -v c="$cov" 'BEGIN{exit !(c+0 >= 70)}' || { step "cobertura >=70%" "✗" "${cov}%"; FAIL=1; }
  else
    step "testes (todos)" "✗" "$(tail -5 /tmp/test.out)"; echo "== RESULTADO: ✗ reprovado =="; exit 1
  fi
fi

# 5) build (nunca no quick)
if [ "$MODE" = "full" ]; then
  if go build ./... >/tmp/build.out 2>&1; then step "go build" "✓" ""; else step "go build" "✗" "$(tail -3 /tmp/build.out)"; echo "== RESULTADO: ✗ reprovado =="; exit 1; fi
else
  step "go build" "pulado" "quick não roda build"
fi

END=$(date +%s)
echo "-- total: $((END-START))s"
if [ "$FAIL" -eq 0 ]; then echo "== RESULTADO: ✓ aprovado =="; else echo "== RESULTADO: ✗ reprovado =="; exit 1; fi
