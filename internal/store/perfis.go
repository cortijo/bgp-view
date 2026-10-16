package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Perfil agrupa communities reutilizáveis (ex.: "rtbh", "prepend2").
type Perfil struct {
	Nome         string   `json:"nome"`
	Comunities   []string `json:"comunities"`
	CriadoEm     string   `json:"criado_em,omitempty"`
	AtualizadoEm string   `json:"atualizado_em,omitempty"`
}

// perfisPadrao é o seed aplicado no boot (specs/005).
var perfisPadrao = []Perfil{
	{Nome: "rtbh", Comunities: []string{"65000:666"}},
	{Nome: "prepend2", Comunities: []string{"65000:200"}},
	{Nome: "prefer-ix", Comunities: []string{"65000:300"}},
}

// ListarPerfis devolve todos os perfis cadastrados, ordenados por nome.
func (s *Store) ListarPerfis() ([]Perfil, error) {
	rows, err := s.db.Query(`SELECT nome, comunities, criado_em, atualizado_em
		FROM comunities_perfil ORDER BY nome`)
	if err != nil {
		return nil, fmt.Errorf("listar perfis: %w", err)
	}
	defer rows.Close() //nolint:errcheck // erro de fechamento refletido em rows.Err()

	var out []Perfil
	for rows.Next() {
		var p Perfil
		var comunities string
		if err := rows.Scan(&p.Nome, &comunities, &p.CriadoEm, &p.AtualizadoEm); err != nil {
			return nil, fmt.Errorf("scan perfil: %w", err)
		}
		if comunities != "" {
			if err := json.Unmarshal([]byte(comunities), &p.Comunities); err != nil {
				return nil, fmt.Errorf("decodificar comunities de %s: %w", p.Nome, err)
			}
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// CriarPerfil insere um perfil; retorna false se o nome já existe.
func (s *Store) CriarPerfil(nome string, comunities []string) (bool, error) {
	if nome == "" || len(comunities) == 0 {
		return false, fmt.Errorf("perfil exige nome e ao menos uma community")
	}
	b, err := json.Marshal(comunities)
	if err != nil {
		return false, fmt.Errorf("serializar comunities: %w", err)
	}
	res, err := s.db.Exec(`INSERT INTO comunities_perfil(nome, comunities) VALUES (?, ?)`, nome, string(b))
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// SeedPerfisPopula cria os perfis padrão que ainda não existem (idempotente).
func (s *Store) SeedPerfisPopula() error {
	for _, p := range perfisPadrao {
		b, err := json.Marshal(p.Comunities)
		if err != nil {
			return fmt.Errorf("serializar perfil %s: %w", p.Nome, err)
		}
		if _, err := s.db.Exec(`INSERT OR IGNORE INTO comunities_perfil(nome, comunities) VALUES (?, ?)`,
			p.Nome, string(b)); err != nil {
			return fmt.Errorf("seed perfil %s: %w", p.Nome, err)
		}
	}
	return nil
}

// PerfilExiste indica se há perfil com o nome informado.
func (s *Store) PerfilExiste(nome string) (bool, error) {
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM comunities_perfil WHERE nome = ?`, nome).Scan(&n)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return n > 0, nil
}
