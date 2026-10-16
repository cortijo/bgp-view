-- 0001: schema inicial (aplicada no boot via internal/store/migrate.go)
-- Mantida aqui como documentação da migração versionada.
CREATE TABLE IF NOT EXISTS sessao_bgp (
  id              INTEGER PRIMARY KEY CHECK (id = 1),
  neighbor_addr   TEXT NOT NULL,
  neighbor_asn    INTEGER NOT NULL,
  local_asn       INTEGER NOT NULL,
  md5_ref         TEXT NOT NULL,
  estado          TEXT NOT NULL DEFAULT 'Unknown',
  atualizado_em   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comunities_perfil (
  nome          TEXT PRIMARY KEY,
  comunities    TEXT NOT NULL,
  criado_em     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS prefixos (
  cidr         TEXT PRIMARY KEY,
  comunities   TEXT NOT NULL DEFAULT '[]',
  perfil       TEXT REFERENCES comunities_perfil(nome) ON DELETE SET NULL,
  ativo        BOOLEAN NOT NULL DEFAULT 1,
  criado_em     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  atualizado_em TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS auditoria (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  acao       TEXT NOT NULL,
  prefixo    TEXT,
  comunities TEXT,
  usuario    TEXT,
  detalhe    TEXT,
  criado_em   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_prefixos_ativo ON prefixos(ativo);
CREATE INDEX IF NOT EXISTS idx_auditoria_prefixo ON auditoria(prefixo);
CREATE INDEX IF NOT EXISTS idx_auditoria_criado ON auditoria(criado_em);
