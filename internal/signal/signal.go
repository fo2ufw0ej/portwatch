// Package signal handles OS signal interception for graceful shutdown.
package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Handler listens for termination signals and cancels a context.
type Handler struct {
	ctx    context.Context
	cancel context.CancelFunc
	ch     chan os.Signal
}

// New creates a Handler that cancels the returned context on SIGINT or SIGTERM.
func New(parent context.Context) (*Handler, context.Context) {
	ctx, cancel := context.WithCancel(parent)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	h := &Handler{ctx: ctx, cancel: cancel, ch: ch}
	go h.listen()
	return h, ctx
}

func (h *Handler) listen() {
	select {
	case <-h.ch:
		h.cancel()
	case <-h.ctx.Done():
	}
}

// Stop releases resources held by the handler.
func (h *Handler) Stop() {
	signal.Stop(h.ch)
	h.cancel()
}
