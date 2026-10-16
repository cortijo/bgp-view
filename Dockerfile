# Estágio 1: dependências antes do código (favorece cache de build)
FROM golang:1.22-alpine AS deps
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

# Estágio 2: compilação
FROM deps AS build
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/bgpadvert ./cmd/bgpadvert

# Estágio 3: imagem mínima, não-root, somente leitura
FROM alpine:3.20
RUN adduser -D -u 10001 bgp && mkdir -p /data && chown bgp:bgp /data
COPY --from=build /out/bgpadvert /usr/local/bin/bgpadvert
USER bgp
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O /dev/null http://127.0.0.1:8080/health || exit 1
ENTRYPOINT ["/usr/local/bin/bgpadvert"]
