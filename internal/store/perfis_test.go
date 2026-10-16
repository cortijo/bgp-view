package store

import (
	"path/filepath"
	"testing"
)

func novoStoreDeTeste(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "perfis.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		if cerr := s.Close(); cerr != nil {
			t.Errorf("close: %v", cerr)
		}
	})
	return s
}

func TestSeedPerfisPopula(t *testing.T) {
	s := novoStoreDeTeste(t)
	if err := s.SeedPerfisPopula(); err != nil {
		t.Fatalf("seed: %v", err)
	}
	perfis, err := s.ListarPerfis()
	if err != nil {
		t.Fatalf("listar: %v", err)
	}
	if len(perfis) != 3 {
		t.Errorf("esperava 3 perfis padrão, obtive %d", len(perfis))
	}
	if perfis[0].Nome != "prefer-ix" { // ordenado por nome (prepend2, prefer-ix, rtbh)
		t.Errorf("primeiro perfil = %q, esperado prefer-ix", perfis[0].Nome)
	}
	// idempotente
	if err := s.SeedPerfisPopula(); err != nil {
		t.Fatalf("seed 2a vez: %v", err)
	}
	if perfis, _ = s.ListarPerfis(); len(perfis) != 3 {
		t.Errorf("seed duplicou perfis: %d", len(perfis))
	}
}

func TestCriarPerfil(t *testing.T) {
	s := novoStoreDeTeste(t)
	ok, err := s.CriarPerfil("backup", []string{"65000:400"})
	if err != nil || !ok {
		t.Fatalf("criar perfil: ok=%v err=%v", ok, err)
	}
	existe, err := s.PerfilExiste("backup")
	if err != nil || !existe {
		t.Fatalf("PerfilExiste: existe=%v err=%v", existe, err)
	}
	if _, err := s.CriarPerfil("", nil); err == nil {
		t.Error("esperava erro para perfil invalido")
	}
}

func TestListarPerfisVazio(t *testing.T) {
	s := novoStoreDeTeste(t)
	perfis, err := s.ListarPerfis()
	if err != nil {
		t.Fatalf("listar: %v", err)
	}
	if len(perfis) != 0 {
		t.Errorf("esperava 0 perfis, obtive %d", len(perfis))
	}
}
