package rungroup

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type RunGroup struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel func()
}

func NewRunGroup(ctx context.Context) *RunGroup {
	wg := new(sync.WaitGroup)

	ctx, cancel := context.WithCancel(ctx)

	return &RunGroup{
		wg:  wg,
		ctx: ctx,
		cancel: func() {
			wg.Done()
			cancel()
		},
	}
}

func (rg *RunGroup) Run(fn func(ctx context.Context)) {
	rg.wg.Add(1)
	go func() {
		defer rg.cancel()
		fn(rg.ctx)
	}()
}

func (rg *RunGroup) Wait() {
	rg.wg.Wait()
}

func InterruptChannel(ctx context.Context) chan interface{} {
	systemCh := make(chan os.Signal, 1)
	signal.Notify(systemCh, os.Interrupt, syscall.SIGTERM)
	cancelCh := ctx.Done()

	ret := make(chan interface{}, 1)
	var closed bool
	close := func() {
		if !closed {
			ret <- true
			close(ret)
		}
	}

	go func() {
		select {
		case <-systemCh:
			close()
		case <-cancelCh:
			close()
		}
	}()

	return ret
}
