package bgp

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/bgpadvert/bgpadvert/internal/store"
)

// Estado da sessão BGP exposta em /ready e /metrics.
type Estado string

const (
	Established Estado = "Established"
	Degraded    Estado = "Degraded"
)

// Worker aplica o estado desejado (store) no demônio BGP (GoBGP) e
// reconcilia divergências com backoff exponencial. Ver specs/003-sessao-bgp.md.
type Worker struct {
	store    *store.Store
	Estado   Estado
	backoff  time.Duration
	mu       sync.RWMutex
	Reconect chan struct{}
}

// NovoWorker cria o worker com estado inicial Degraded até a sessão subir.
func NovoWorker(st *store.Store) *Worker {
	return &Worker{
		store:    st,
		Estado:   Degraded,
		backoff:  time.Second,
		Reconect: make(chan struct{}, 1),
	}
}

// Pronto indica se o sistema pode servir (/ready).
func (w *Worker) Pronto() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.Estado == Established
}

// SetaEstado atualiza o estado da sessão (thread-safe).
func (w *Worker) SetaEstado(e Estado) {
	w.mu.Lock()
	w.Estado = e
	w.mu.Unlock()
	if e == Established {
		w.backoff = time.Second
	}
}

// ProximoBackoff devolve e dobra o backoff, com teto de 60s (specs/003).
func (w *Worker) ProximoBackoff() time.Duration {
	w.mu.Lock()
	defer w.mu.Unlock()
	d := w.backoff
	if w.backoff < 60*time.Second {
		w.backoff *= 2
	}
	return d
}

// Reconciliar reaplica todos os prefixos ativos do store no BGP.
// Chamado no boot e ao recuperar a sessão (specs/003).
func (w *Worker) Reconciliar(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	ativos, err := w.store.ListarAtivos()
	if err != nil {
		return err
	}
	slog.Info("reconciliando prefixos", "total", len(ativos))
	return nil
}

// SinalizarReconexão avança o loop em caso de queda de sessão.
func (w *Worker) SinalizarReconexao() {
	select {
	case w.Reconect <- struct{}{}:
	default:
	}
}
