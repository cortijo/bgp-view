package audit

import (
	"encoding/json"
	"log/slog"
	"time"
)

// Evento é uma entrada do log de auditoria (retenção 90 dias, expurgo diário).
type Evento struct {
	Timestamp  time.Time `json:"timestamp"`
	Acao       string    `json:"acao"`
	Prefixo    string    `json:"prefixo,omitempty"`
	Comunities []string  `json:"comunities,omitempty"`
	Usuario    string    `json:"usuario"`
	Detalhe    string    `json:"detalhe,omitempty"`
}

// Registrar grava o evento no log estruturado de auditoria.
func Registrar(e Evento) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	if e.Usuario == "" {
		e.Usuario = "api-key"
	}
	b, _ := json.Marshal(e)
	slog.Info("auditoria", "evento", string(b))
}
