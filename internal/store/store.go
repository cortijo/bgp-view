package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Store encapsula o acesso ao SQLite (estado desejado dos anúncios).
type Store struct {
	db *sql.DB
}

// Open abre/cria o banco e aplica as migrations.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("abrir sqlite: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrar: %w", err)
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) DB() *sql.DB { return s.db }
