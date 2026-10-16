package audit

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestRegistrar(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })

	Registrar(Evento{Acao: "announce", Prefixo: "203.0.113.0/24", Comunities: []string{"65000:666"}})

	// O slog serializa o evento como string JSON escapada dentro do log externo.
	var externo map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &externo); err != nil {
		t.Fatalf("log externo nao e JSON: %v", err)
	}
	// "evento" chega como string JSON (dupla serializacao: slog -> string -> JSON externo).
	var crua string
	if err := json.Unmarshal(externo["evento"], &crua); err != nil {
		t.Fatalf("evento nao e string JSON: %v", err)
	}
	var evento Evento
	if err := json.Unmarshal([]byte(crua), &evento); err != nil {
		t.Fatalf("evento nao e JSON: %v", err)
	}

	if evento.Acao != "announce" {
		t.Errorf("acao = %q, esperado announce", evento.Acao)
	}
	if evento.Prefixo != "203.0.113.0/24" {
		t.Errorf("prefixo = %q", evento.Prefixo)
	}
	if evento.Usuario != "api-key" {
		t.Errorf("usuario default = %q", evento.Usuario)
	}
	if len(evento.Comunities) != 1 || evento.Comunities[0] != "65000:666" {
		t.Errorf("comunities = %v", evento.Comunities)
	}
	if evento.Timestamp.IsZero() {
		t.Error("timestamp ausente")
	}
}

func TestRegistrarUsuarioExplicito(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })

	Registrar(Evento{Acao: "withdraw", Usuario: "op1"})

	var externo map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &externo); err != nil {
		t.Fatalf("log externo nao e JSON: %v", err)
	}
	// "evento" chega como string JSON (dupla serializacao: slog -> string -> JSON externo).
	var crua string
	if err := json.Unmarshal(externo["evento"], &crua); err != nil {
		t.Fatalf("evento nao e string JSON: %v", err)
	}
	var evento Evento
	if err := json.Unmarshal([]byte(crua), &evento); err != nil {
		t.Fatalf("evento nao e JSON: %v", err)
	}
	if evento.Usuario != "op1" {
		t.Errorf("usuario = %q, esperado op1", evento.Usuario)
	}
}
