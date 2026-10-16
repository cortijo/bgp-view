package store

import (
	"path/filepath"
	"testing"
)

func TestOpenMigraESchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() {
		if cerr := s.Close(); cerr != nil {
			t.Errorf("close: %v", cerr)
		}
	}()

	for _, tbl := range []string{"prefixos", "comunities_perfil", "sessao_bgp", "auditoria", "schema_migrations"} {
		var n int
		qerr := s.DB().QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, tbl).Scan(&n)
		if qerr != nil {
			t.Fatalf("consultar %s: %v", tbl, qerr)
		}
		if n != 1 {
			t.Errorf("tabela %s não criada", tbl)
		}
	}

	// idempotência: reabrir não reexecuta a migration
	if cerr := s.Close(); cerr != nil {
		t.Fatalf("fechar 1o store: %v", cerr)
	}
	s2, err := Open(path)
	if err != nil {
		t.Fatalf("reabrir: %v", err)
	}
	defer func() {
		if cerr := s2.Close(); cerr != nil {
			t.Errorf("close: %v", cerr)
		}
	}()
	var count int
	if err := s2.DB().QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&count); err != nil {
		t.Fatalf("contar migrations: %v", err)
	}
	if count != 1 {
		t.Errorf("esperava 1 migration aplicada, obtive %d", count)
	}
}

func TestPrefixosCRUD(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crud.db")
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() {
		if cerr := s.Close(); cerr != nil {
			t.Errorf("close: %v", cerr)
		}
	}()

	_, err = s.DB().Exec(`INSERT INTO prefixos(cidr, comunities) VALUES (?, ?)`, "203.0.113.0/24", `["65000:666"]`)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	var comunities string
	var ativo bool
	if err := s.DB().QueryRow(`SELECT comunities, ativo FROM prefixos WHERE cidr = ?`, "203.0.113.0/24").Scan(&comunities, &ativo); err != nil {
		t.Fatalf("select: %v", err)
	}
	if comunities != `["65000:666"]` || !ativo {
		t.Errorf("prefixo inconsistente: comunities=%s ativo=%v", comunities, ativo)
	}
}

func TestListarAtivos(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ativos.db")
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() {
		if cerr := s.Close(); cerr != nil {
			t.Errorf("close: %v", cerr)
		}
	}()

	_, err = s.DB().Exec(`INSERT INTO prefixos(cidr, comunities, ativo) VALUES (?, ?, 1)`, "198.51.100.0/24", `["65000:100"]`)
	if err != nil {
		t.Fatalf("insert ativo: %v", err)
	}
	_, err = s.DB().Exec(`INSERT INTO prefixos(cidr, comunities, ativo) VALUES (?, ?, 0)`, "203.0.113.0/24", `["65000:666"]`)
	if err != nil {
		t.Fatalf("insert inativo: %v", err)
	}

	ativos, err := s.ListarAtivos()
	if err != nil {
		t.Fatalf("ListarAtivos: %v", err)
	}
	if len(ativos) != 1 {
		t.Fatalf("esperava 1 ativo, obtive %d", len(ativos))
	}
	if ativos[0].Cidr != "198.51.100.0/24" {
		t.Errorf("cidr = %q", ativos[0].Cidr)
	}
	if len(ativos[0].Comunities) != 1 || ativos[0].Comunities[0] != "65000:100" {
		t.Errorf("comunities = %v", ativos[0].Comunities)
	}
}
