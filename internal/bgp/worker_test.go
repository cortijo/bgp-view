package bgp

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/bgpadvert/bgpadvert/internal/store"
)

func novoWorkerDeTeste(t *testing.T) *Worker {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "bgp.db"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NovoWorker(st)
}

func TestEstadoInicialDegradado(t *testing.T) {
	w := novoWorkerDeTeste(t)
	if w.Estado != Degraded {
		t.Errorf("estado inicial = %q, esperado Degraded", w.Estado)
	}
	if w.Pronto() {
		t.Error("worker deveria nascer degraded (nao pronto)")
	}
}

func TestSetaEstadoEstabelecido(t *testing.T) {
	w := novoWorkerDeTeste(t)
	w.SetaEstado(Established)
	if !w.Pronto() {
		t.Error("worker deveria ficar pronto apos Established")
	}
}

func TestBackoffExponencialComTeto(t *testing.T) {
	w := novoWorkerDeTeste(t)
	w.backoff = 30 * time.Second
	if d := w.ProximoBackoff(); d != 30*time.Second {
		t.Errorf("primeiro backoff = %v", d)
	}
	// 30s * 2 = 60s (atinge o teto)
	if d := w.ProximoBackoff(); d != 60*time.Second {
		t.Errorf("segundo backoff = %v", d)
	}
	// permanece no teto de 60s
	if d := w.ProximoBackoff(); d != 60*time.Second {
		t.Errorf("backoff acima do teto: %v", d)
	}
}

func TestReconciliarListaAtivos(t *testing.T) {
	w := novoWorkerDeTeste(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := w.Reconciliar(ctx); err != nil {
		t.Errorf("Reconciliar: %v", err)
	}
}

func TestReconciliarCancelado(t *testing.T) {
	w := novoWorkerDeTeste(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := w.Reconciliar(ctx); err == nil {
		t.Error("esperava erro de contexto cancelado")
	}
}

func TestSinalizarReconexaoNaoBloqueia(t *testing.T) {
	w := novoWorkerDeTeste(t)
	for i := 0; i < 3; i++ {
		w.SinalizarReconexao() // não deve bloquear com canal cheio
	}
}
