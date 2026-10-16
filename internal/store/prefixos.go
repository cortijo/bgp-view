package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Prefixo é o estado desejado de um anúncio (fonte de verdade).
type Prefixo struct {
	Cidr         string   `json:"cidr"`
	Comunities   []string `json:"comunities"`
	Perfil       string   `json:"perfil,omitempty"`
	Ativo        bool     `json:"ativo"`
	CriadoEm     string   `json:"criado_em"`
	AtualizadoEm string   `json:"atualizado_em"`
}

// ListarAtivos devolve os prefixos ativos (usado pelo worker de reconciliação).
func (s *Store) ListarAtivos() ([]Prefixo, error) {
	rows, err := s.db.Query(`SELECT cidr, comunities, COALESCE(perfil,''), ativo, criado_em, atualizado_em
		FROM prefixos WHERE ativo = 1 ORDER BY cidr`)
	if err != nil {
		return nil, fmt.Errorf("listar ativos: %w", err)
	}
	defer rows.Close() //nolint:errcheck // erro de fechamento já refletido em rows.Err()

	var out []Prefixo
	for rows.Next() {
		var p Prefixo
		var comunities, perfil sql.NullString
		if err := rows.Scan(&p.Cidr, &comunities, &perfil, &p.Ativo, &p.CriadoEm, &p.AtualizadoEm); err != nil {
			return nil, fmt.Errorf("scan prefixo: %w", err)
		}
		if comunities.Valid && comunities.String != "" {
			if err := json.Unmarshal([]byte(comunities.String), &p.Comunities); err != nil {
				return nil, fmt.Errorf("decodificar comunities de %s: %w", p.Cidr, err)
			}
		}
		p.Perfil = perfil.String
		out = append(out, p)
	}
	return out, rows.Err()
}
