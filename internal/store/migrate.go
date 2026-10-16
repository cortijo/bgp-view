package store

import (
	"fmt"
	"sort"
)

// schemaMigrations define as migrations versionadas aplicadas no boot.
// Ordem numérica crescente; reversões quando possível.
var schemaMigrations = []struct {
	version string
	sql     string
}{
	{"0001", `
CREATE TABLE IF NOT EXISTS schema_migrations (
  version    TEXT PRIMARY KEY,
  applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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
  nome        TEXT PRIMARY KEY,
  comunities  TEXT NOT NULL,
  criado_em    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
`},
}

func (s *Store) migrate() error {
	migrations := append([]struct {
		version string
		sql     string
	}{}, schemaMigrations...)
	sort.Slice(migrations, func(i, j int) bool { return migrations[i].version < migrations[j].version })

	for _, m := range migrations {
		var exists int
		err := s.db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, m.version).Scan(&exists)
		if err != nil {
			// tabela ainda não existe: cria e segue
			if _, cerr := s.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
				version TEXT PRIMARY KEY,
				applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			)`); cerr != nil {
				return fmt.Errorf("criar schema_migrations: %w", cerr)
			}
		}
		if exists > 0 {
			continue
		}
		if _, err := s.db.Exec(m.sql); err != nil {
			return fmt.Errorf("migracao %s: %w", m.version, err)
		}
		if _, err := s.db.Exec(`INSERT INTO schema_migrations(version) VALUES (?)`, m.version); err != nil {
			return fmt.Errorf("registrar migracao %s: %w", m.version, err)
		}
	}
	return nil
}
